// Package toolrunner provides a Unix-socket transport for Norma's built-in
// file and shell tools. It isolates execution only when the server itself is
// deployed in a separate process/container without control-plane credentials.
package toolrunner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Autumn-27/norma/tool"
)

type Request struct {
	Session        string          `json:"session"`
	Name           string          `json:"name"`
	Input          json.RawMessage `json:"input"`
	WorkingDir     string          `json:"working_dir"`
	OutputDir      string          `json:"output_dir"`
	Env            []string        `json:"env"`
	MaxOutputChars int             `json:"max_output_chars"`
	CACertPEM          string          `json:"ca_cert_pem,omitempty"`
	DisableBackground bool            `json:"disable_background,omitempty"`
	DisableInteractive bool           `json:"disable_interactive,omitempty"`
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
	digest := sha256.Sum256([]byte(session))
	m, err := tool.NewManager("artex-runner-" + hex.EncodeToString(digest[:16]))
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Session == "" {
		http.Error(w, "invalid session", http.StatusBadRequest)
		return
	}
	s.closeSession(req.Session)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) handleCall(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid tool request", http.StatusBadRequest)
		return
	}
	if req.Session == "" {
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
	// Preserve the controller's existing SDK opt-outs even though the
	// runner does not inherit control-plane environment variables.
	if req.DisableBackground || (req.DisableInteractive && strings.HasPrefix(req.Name, "shell_")) {
		tc.Tasks = nil
	}
	result, callErr := t.Call(r.Context(), req.Input, tc)
	reply := Response{Result: result}
	if callErr != nil {
		reply.Error = callErr.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(reply)
}

// Serve listens only on a filesystem socket: no TCP management endpoint is
// exposed. Socket mount ownership is enforced by the Docker deployment.
// remotePeerOnlyListener rejects commands running inside the tool container.
// In the opt-in Docker deployment the controller lives in a separate PID
// namespace, so its PID is not visible to the receiving namespace (PID=0).
// This is an extra control-plane check, not a substitute for per-task identity.
type remotePeerOnlyListener struct { net.Listener }

func (l remotePeerOnlyListener) Accept() (net.Conn, error) {
	for {
		conn, err := l.Listener.Accept()
		if err != nil {
			return nil, err
		}
		if peerOutsidePIDNamespace(conn) {
			return conn, nil
		}
		_ = conn.Close()
	}
}

func Serve(ctx context.Context, socket string) error {
	service := NewService()
	defer service.Close()
	return serveOnUnix(ctx, socket, service.Handler())
}

// serveOnUnix is shared by the legacy direct Runner and the isolated broker.
func serveOnUnix(ctx context.Context, socket string, handler http.Handler) error {
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
	if os.Getenv("ARTEX_RUNNER_REMOTE_PEER_ONLY") == "1" {
		if !peerGuardSupported() {
			return errors.New("remote-peer-only runner requires Linux SO_PEERCRED")
		}
		ln = remotePeerOnlyListener{Listener: ln}
	}
	if err := os.Chmod(socket, 0600); err != nil {
		return err
	}
	return serveHTTP(ctx, ln, handler)
}

func serveHTTP(ctx context.Context, ln net.Listener, handler http.Handler) error {
	server := &http.Server{Handler: handler, ReadHeaderTimeout: 5 * time.Second}
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
	err := server.Serve(ln)
	close(done)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
