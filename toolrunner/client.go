package toolrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Autumn-27/norma/agentcore"
	"github.com/Autumn-27/norma/tool"
	"github.com/google/uuid"
)

// Client communicates only over the mounted private Unix socket.
type Client struct {
	http   *http.Client
	socket string
}

func NewClient(socket string) *Client {
	dialer := &net.Dialer{Timeout: 3 * time.Second}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, "unix", socket)
		},
	}
	return &Client{socket: socket, http: &http.Client{Transport: transport}}
}

func (c *Client) post(ctx context.Context, endpoint string, request any, response any) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://runner"+endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("isolated runner unavailable: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("isolated runner HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	if response == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(response)
}

func (c *Client) Call(ctx context.Context, req Request) (tool.Result, error) {
	var result Response
	if err := c.post(ctx, "/call", req, &result); err != nil {
		return tool.Result{}, err
	}
	if result.Error != "" {
		return result.Result, errors.New(result.Error)
	}
	return result.Result, nil
}

func (c *Client) CloseSession(ctx context.Context, session string) error {
	return c.post(ctx, "/close", struct {
		Session string `json:"session"`
	}{Session: session}, nil)
}

// runnerToolEnv only forwards transport configuration. The control plane may
// hold database passwords, provider keys and other secrets in its environment;
// an execution runner must never inherit extra BashEnv entries by accident.
// ARTEX's current proxyEnv uses these exact proxy/CA variables.
func runnerToolEnv(env []string) []string {
	allowed := map[string]bool{
		"HTTP_PROXY": true, "HTTPS_PROXY": true, "ALL_PROXY": true, "NO_PROXY": true,
		"http_proxy": true, "https_proxy": true, "all_proxy": true, "no_proxy": true,
		"SSL_CERT_FILE": true, "CURL_CA_BUNDLE": true,
		"REQUESTS_CA_BUNDLE": true, "NODE_EXTRA_CA_CERTS": true,
		"GIT_SSL_CAINFO": true,
	}
	var out []string
	for _, item := range env {
		key, _, found := strings.Cut(item, "=")
		if found && allowed[key] {
			out = append(out, item)
		}
	}
	return out
}

type remoteTool struct {
	tool.CoreTool // permissions/schema/name remain those of the original SDK tool
	client  *Client
	session string
	disableBackground bool
	disableInteractive bool
}

func (r remoteTool) Call(ctx context.Context, input json.RawMessage, tc *tool.ToolContext) (tool.Result, error) {
	request := Request{Session: r.session, Name: r.Name(), Input: input, DisableBackground: r.disableBackground, DisableInteractive: r.disableInteractive}
	if tc != nil {
		request.WorkingDir = tc.WorkingDir
		request.OutputDir = tc.OutputDir
		request.MaxOutputChars = tc.MaxOutputChars
		request.Env = runnerToolEnv(tc.Env)
	}
	// Preserve existing traffic recording through the Compose-only relay.
	if proxy := os.Getenv("ARTEX_RUNNER_PROXY_ADDR"); proxy != "" {
		for i, item := range request.Env {
			key, value, ok := strings.Cut(item, "=")
			if !ok || !strings.Contains(strings.ToUpper(key), "PROXY") {
				continue
			}
			u, err := url.Parse(value)
			if err == nil && (u.Hostname() == "127.0.0.1" || u.Hostname() == "localhost") && u.Port() == "8788" {
				request.Env[i] = key + "=" + proxy
			}
		}
	}
	// Transmit only the *public certificate* used to verify a recording proxy.
	// Never mount the proxy's CA private key or the main data directory.
	for _, item := range request.Env {
		if !strings.HasPrefix(item, "SSL_CERT_FILE=") {
			continue
		}
		path := strings.TrimPrefix(item, "SSL_CERT_FILE=")
		if path == "" {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return tool.Result{}, fmt.Errorf("read proxy public certificate: %w", err)
		}
		if len(data) > 256*1024 {
			return tool.Result{}, errors.New("proxy certificate too large")
		}
		if !bytes.Contains(data, []byte("-----BEGIN CERTIFICATE-----")) ||
			bytes.Contains(data, []byte("PRIVATE KEY")) {
			return tool.Result{}, errors.New("proxy certificate must contain public certificates only")
		}
		request.CACertPEM = string(data)
		break
	}
	return r.client.Call(ctx, request)
}

var runnerTools = map[string]bool{
	"Read": true, "Write": true, "Edit": true, "MultiEdit": true,
	"LS": true, "Glob": true, "Grep": true, "Bash": true,
	"shell_open": true, "shell_send": true, "shell_read": true,
	"shell_close": true, "shell_list": true,
	"TaskOutput": true, "TaskStop": true, "TaskList": true, "Monitor": true,
}

// Configure replaces the execution-capable SDK builtins with Unix-socket
// proxies only when ARTEX_RUNNER_SOCKET is explicitly set. The runner has no
// access to controller environment or keys: the caller forwards only per-tool
// proxy settings and public CA certificate. Missing runner => error on tool
// call, never a silent fallback to local execution.
//
// The cleanup closure releases remote background/PTY sessions at the end of a
// run. Other ARTEX domain tools and third-party MCP/skill processes are not
// isolated by this initial opt-in bridge.
func Configure(opts *agentcore.Options) func() {
	socket := os.Getenv("ARTEX_RUNNER_SOCKET")
	if socket == "" || opts == nil {
		return func() {}
	}
	if opts.SessionID == "" {
		opts.SessionID = uuid.NewString()
	}
	if opts.Tools == nil {
		opts.Tools = tool.DefaultTools()
	}
	c := NewClient(socket)
	backgroundDisabled := opts.DisableBackgroundTasks || tool.BackgroundTasksDisabled()
	interactiveDisabled := tool.InteractiveShellDisabled()
	for i, t := range opts.Tools {
		if runnerTools[t.Name()] {
			opts.Tools[i] = remoteTool{CoreTool: t, client: c, session: opts.SessionID, disableBackground: backgroundDisabled, disableInteractive: interactiveDisabled}
		}
	}
	// The controller must never create a local background task manager in
	// runner mode. Only expose remote management tools when previously enabled.
	opts.DisableBackgroundTasks = true
	if !backgroundDisabled {
		for _, t := range []tool.CoreTool{
			tool.NewTaskOutput(), tool.NewTaskStop(), tool.NewTaskList(), tool.NewMonitor(),
		} {
			opts.Tools = append(opts.Tools, remoteTool{CoreTool: t, client: c, session: opts.SessionID, disableInteractive: interactiveDisabled})
		}
	}
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = c.CloseSession(ctx, opts.SessionID)
	}
}
