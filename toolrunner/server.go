// Package toolrunner provides a Unix-socket transport for Norma's built-in
// file and shell tools. It isolates execution only when the server itself is
// deployed in a separate process/container without control-plane credentials.
package toolrunner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/Autumn-27/norma/tool"
)

const maxRequestBytes = 2 << 20

var validSession = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

type Request struct {
	Session        string          `json:"session"`
	Name           string          `json:"name"`
	Input          json.RawMessage `json:"input"`
	WorkingDir     string          `json:"working_dir"`
	OutputDir      string          `json:"output_dir"`
	Env            []string        `json:"env"`
	MaxOutputChars int             `json:"max_output_chars"`
	CACertPEM      string          `json:"ca_cert_pem,omitempty"`
}

type Response struct {
	Result tool.Result `json:"result"`
	Error  string      `json:"error,omitempty"`
}

type Service struct {
	mu       sync.Mutex
	sessions map[string]*tool.Manager
	tools    map[string]tool.CoreTool
}

func NewService() *Service {
	s := &Service{sessions: make(map[string]*tool.Manager), tools: make(map[string]tool.CoreTool)}
	builtins := append(tool.DefaultTools(), tool.ShellSessionTools()...)
	builtins = append(builtins, tool.NewTaskOutput(), tool.NewTaskStop(), tool.NewTaskList(), tool.NewMonitor())
	for _, t := range builtins {
		s.tools[t.Name()] = t
	}
	return s
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /call", s.handleCall)
	mux.HandleFunc("POST /close", s.handleClose)
	return mux
}

func (s *Service) manager(session string) (*tool.Manager, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m := s.sessions[session]; m != nil {
		return m, nil
	}
	m, err := tool.NewManager("artex-runner-" + session)
	if err != nil {
		return nil, err
	}
	s.sessions[session] = m
	return m, nil
}

func (s *Service) closeSession(session string) {
	s.mu.Lock()
	m := s.sessions[session]
	delete(s.sessions, session)
	s.mu.Unlock()
	if m != nil {
		m.Cleanup()
	}
}

func (s *Service) Close() {
	s.mu.Lock()
	managers := s.sessions
	s.sessions = make(map[string]*tool.Manager)
	s.mu.Unlock()
	for _, m := range managers {
		m.Cleanup()
	}
}

func (s *Service) handleClose(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Session string `json:"session"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&req); err != nil || !validSession.MatchString(req.Session) {
		http.Error(w, "invalid session", http.StatusBadRequest)
		return
	}
	s.closeSession(req.Session)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) handleCall(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxRequestBytes)).Decode(&req); err != nil {
		http.Error(w, "invalid tool request", http.StatusBadRequest)
		return
	}
	if !validSession.MatchString(req.Session) {
		http.Error(w, "invalid session", http.StatusBadRequest)
		return
	}
	t, ok := s.tools[req.Name]
	if !ok {
		http.Error(w, "unavailable runner tool", http.StatusBadRequest)
		return
	}
	m, err := s.manager(req.Session)
	if err != nil {
		http.Error(w, "runner session unavailable", http.StatusInternalServerError)
		return
	}
	env := req.Env
	if req.CACertPEM != "" {
		env, err = installPublicCA(env, req.CACertPEM)
		if err != nil {
			http.Error(w, "invalid public CA", http.StatusBadRequest)
			return
		}
	}
	tc := &tool.ToolContext{
		WorkingDir:     req.WorkingDir,
		OutputDir:      req.OutputDir,
		MaxOutputChars: req.MaxOutputChars,
		Env:            env,
		Tasks:          m,
	}
	result, callErr := t.Call(r.Context(), req.Input, tc)
	reply := Response{Result: result}
	if callErr != nil {
		reply.Error = callErr.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reply); err != nil {
		_ = err
	}
}

// Serve listens only on a filesystem socket: no TCP management endpoint is
// exposed. Socket mount ownership is enforced by the Docker deployment.
func Serve(ctx context.Context, socket string) error {
	if socket == "" || !filepath.IsAbs(socket) {
		return errors.New("runner socket must be an absolute path")
	}
	if err := os.MkdirAll(filepath.Dir(socket), 0700); err != nil {
		return fmt.Errorf("create runner socket dir: %w", err)
	}
	if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove stale runner socket: %w", err)
	}
	ln, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}
	defer ln.Close()
	defer os.Remove(socket)
	if err := os.Chmod(socket, 0600); err != nil {
		return err
	}
	service := NewService()
	defer service.Close()
	server := &http.Server{Handler: service.Handler(), ReadHeaderTimeout: 5 * time.Second}
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = server.Shutdown(shutdownCtx)
		case <-done:
		}
	}()
	err = server.Serve(ln)
	close(done)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
