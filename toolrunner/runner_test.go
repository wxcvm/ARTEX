package toolrunner

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Autumn-27/norma/agentcore"
	"github.com/Autumn-27/norma/tool"
)

func TestUnixRunnerExecutesBuiltinsInWorkingDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-socket deployment is Linux-only")
	}
	root := t.TempDir()
	socket := filepath.Join(root, "tool.sock")
	workdir := filepath.Join(root, "workspace")
	if err := os.MkdirAll(workdir, 0755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- Serve(ctx, socket) }()
	for i := 0; i < 100; i++ {
		if _, err := os.Stat(socket); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	c := NewClient(socket)
	call := func(name string, input string) tool.Result {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		got, err := c.Call(ctx, Request{
			Session: "unit-test-1", Name: name, Input: json.RawMessage(input),
			WorkingDir: workdir,
		})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if got.IsError {
			t.Fatalf("%s returned error result: %s", name, got.Flatten())
		}
		return got
	}
	call("Write", `{"file_path":"fixture.txt","content":"hello from runner"}`)
	if got := call("Read", `{"file_path":"fixture.txt"}`).Flatten(); !strings.Contains(got, "hello from runner") {
		t.Fatalf("remote Read returned %q", got)
	}
	if got := call("Bash", `{"command":"printf runner-ok"}`).Flatten(); !strings.Contains(got, "runner-ok") {
		t.Fatalf("remote Bash returned %q", got)
	}
	if err := c.CloseSession(context.Background(), "unit-test-1"); err != nil {
		t.Fatal(err)
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("runner failed to stop")
	}
}

func TestConfigureLeavesLegacyModeAlone(t *testing.T) {
	t.Setenv("ARTEX_RUNNER_SOCKET", "")
	opts := agentcore.Options{Tools: tool.DefaultTools()}
	closeRunner := Configure(&opts)
	defer closeRunner()
	if opts.DisableBackgroundTasks {
		t.Fatal("legacy mode unexpectedly disabled background tasks")
	}
	for _, item := range opts.Tools {
		if _, ok := item.(remoteTool); ok {
			t.Fatal("legacy mode unexpectedly replaced a tool")
		}
	}
}

func TestConfiguredRunnerDoesNotFallbackToLocalBash(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-socket deployment is Linux-only")
	}
	t.Setenv("ARTEX_RUNNER_SOCKET", filepath.Join(t.TempDir(), "missing.sock"))
	opts := agentcore.Options{Tools: tool.DefaultTools()}
	closeRunner := Configure(&opts)
	defer closeRunner()
	if !opts.DisableBackgroundTasks {
		t.Fatal("SDK local background manager must be disabled in remote mode")
	}
	var remoteBash tool.CoreTool
	var remoteOutput tool.CoreTool
	for _, item := range opts.Tools {
		if item.Name() == "Bash" {
			remoteBash = item
		}
		if item.Name() == "TaskOutput" {
			remoteOutput = item
		}
	}
	if remoteBash == nil || remoteOutput == nil {
		t.Fatal("missing remote shell or background-management tool")
	}
	if _, ok := remoteBash.(remoteTool); !ok {
		t.Fatalf("Bash remains local: %T", remoteBash)
	}
	if _, ok := remoteOutput.(remoteTool); !ok {
		t.Fatalf("TaskOutput remains local: %T", remoteOutput)
	}
	result, err := remoteBash.Call(context.Background(), json.RawMessage(`{"command":"printf must-not-run"}`), &tool.ToolContext{WorkingDir: t.TempDir()})
	if err == nil || !strings.Contains(err.Error(), "runner unavailable") {
		t.Fatalf("missing runner must fail, never execute locally: result=%+v, err=%v", result, err)
	}
}

func TestRunnerRejectsPrivateCA(t *testing.T) {
	if _, err := installPublicCA(nil, "-----BEGIN PRIVATE KEY-----\nsecret\n-----END PRIVATE KEY-----"); err == nil {
		t.Fatal("private key should never be materialized in the runner")
	}
}

func TestRunnerRejectsUnknownTool(t *testing.T) {
	s := NewService()
	defer s.Close()
	if _, ok := s.tools["model_provider_key"]; ok {
		t.Fatal("control-plane tool was mistakenly exposed")
	}
	if _, ok := s.tools["Bash"]; !ok {
		t.Fatal("Bash must be available")
	}
}

func TestRunnerToolEnvForwardsOnlyProxyAndPublicCA(t *testing.T) {
	env := []string{
		"HTTP_PROXY=http://127.0.0.1:8788",
		"ALL_PROXY=socks5://proxy.example:1080",
		"SSL_CERT_FILE=/app/data/_ca/public.pem",
		"OPENAI_API_KEY=model-secret",
		"ARTEX_PG_DSN=postgres://user:password@postgres/db",
		"ARTEX_LLM_SECRET_KEY=internal-secret",
		"CUSTOM_API_TOKEN=credential",
		"PATH=/custom/bin",
	}
	got := runnerToolEnv(env)
	if len(got) != 3 {
		t.Fatalf("forwarded unexpected environment variables: %v", got)
	}
	for i, want := range env[:3] {
		if got[i] != want {
			t.Fatalf("env[%d]=%q, want %q", i, got[i], want)
		}
	}
}

func TestRunnerPreservesExistingSDKDisableFlags(t *testing.T) {
	t.Setenv("ARTEX_RUNNER_SOCKET", filepath.Join(t.TempDir(), "missing.sock"))
	t.Setenv("AGENT_CORE_DISABLE_BACKGROUND_TASKS", "1")
	t.Setenv("AGENT_CORE_DISABLE_INTERACTIVE_SHELL", "1")
	opts := agentcore.Options{Tools: tool.DefaultTools()}
	cleanup := Configure(&opts)
	defer cleanup()
	if !opts.DisableBackgroundTasks {
		t.Fatal("local SDK task manager was enabled")
	}
	for _, item := range opts.Tools {
		if item.Name() == "TaskOutput" || item.Name() == "TaskList" {
			t.Fatalf("background task tool exposed despite existing global opt-out: %s", item.Name())
		}
		if item.Name() == "Bash" {
			r, ok := item.(remoteTool)
			if !ok || !r.disableBackground || !r.disableInteractive {
				t.Fatalf("runner did not preserve SDK tool disabling flags: %+v", item)
			}
		}
	}
}

func TestRunnerServiceHonorsSDKDisableFlags(t *testing.T) {
	s := NewService()
	defer s.Close()
	root := t.TempDir()
	call := func(req Request) Response {
		t.Helper()
		payload, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, httptest.NewRequest("POST", "/call", strings.NewReader(string(payload))))
		if w.Code != 200 {
			t.Fatalf("unexpected runner HTTP %d: %s", w.Code, w.Body.String())
		}
		var result Response
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		return result
	}
	bash := call(Request{
		Session: "disable-test", Name: "Bash", WorkingDir: root, DisableBackground: true,
		Input: json.RawMessage(`{"command":"printf sync-only","run_in_background":true}`),
	})
	if bash.Result.IsError || !strings.Contains(bash.Result.Flatten(), "sync-only") {
		t.Fatalf("background opt-out must run synchronously: %+v", bash)
	}
	shell := call(Request{
		Session: "disable-test", Name: "shell_open", WorkingDir: root, DisableInteractive: true,
		Input: json.RawMessage(`{"command":"printf should-not-start"}`),
	})
	if !shell.Result.IsError {
		t.Fatalf("interactive Shell opt-out was bypassed: %+v", shell)
	}
}

func TestRemotePeerOnlyRejectsLocalSDKToolRPC(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("SO_PEERCRED only available on Linux")
	}
	t.Setenv("ARTEX_RUNNER_REMOTE_PEER_ONLY", "1")
	root := t.TempDir()
	socket := filepath.Join(root, "runner.sock")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- Serve(ctx, socket) }()
	for i := 0; i < 100; i++ {
		if _, err := os.Stat(socket); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	client := NewClient(socket)
	callCtx, callCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer callCancel()
	_, err := client.Call(callCtx, Request{
		Session: "local-attacker", Name: "Write", WorkingDir: root,
		Input: json.RawMessage(`{"file_path":"bypass.txt","content":"must-not-write"}`),
	})
	if err == nil {
		t.Fatal("local SDK tool bypassed remote-peer-only listener")
	}
	if _, err := os.Stat(filepath.Join(root, "bypass.txt")); !os.IsNotExist(err) {
		t.Fatalf("local caller wrote through control RPC, stat=%v", err)
	}
	cancel()
	select {
	case err := <-done:
		if err != nil { t.Fatal(err) }
	case <-time.After(5*time.Second):
		t.Fatal("guarded listener did not stop")
	}
}
