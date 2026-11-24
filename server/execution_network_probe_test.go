package server

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func resetModeBProbeForTest() {
	modeBProbeState.Lock()
	modeBProbeState.snapshot = modeBProbeSnapshot{}
	modeBProbeState.Unlock()
}

func TestVerifyExecutionNetworkUsesIsolatedRunnerAndStoresOnlyComparison(t *testing.T) {
	for _, tc := range []struct {
		name, runnerStatus, wantVerification, wantFinal, wantActual string
		wantOK bool
	}{
		{"match", "exit_ip_match", "exit_ip_match", "matched_expected", "matches_expected", true},
		{"mismatch", "exit_ip_mismatch", "exit_ip_mismatch", "mismatch", "mismatch", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resetModeBProbeForTest()
			t.Cleanup(resetModeBProbeForTest)
			root := t.TempDir()
			socket := filepath.Join(root, "runner.sock")
			listener, err := net.Listen("unix", socket)
			if err != nil {
				t.Fatal(err)
			}
			server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/egress" {
					t.Errorf("unexpected runner request: %s %s", r.Method, r.URL.Path)
					http.NotFound(w, r)
					return
				}
				var input struct {
					ExpectedExitIP string `json:"expected_exit_ip"`
				}
				if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
					t.Errorf("decode runner request: %v", err)
				}
				if input.ExpectedExitIP != "1.1.1.1" {
					t.Errorf("expected normalized exit IP, got %q", input.ExpectedExitIP)
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]any{"ok": tc.wantOK, "status": tc.runnerStatus})
			})}
			go func() { _ = server.Serve(listener) }()
			t.Cleanup(func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_ = server.Shutdown(ctx)
			})
			t.Setenv("ARTEX_RUNNER_SOCKET", socket)
			t.Setenv("ARTEX_RUNNER_PROXY_ADDR", "http://mode-b-egress:7890")
			s := &Server{}
			req := httptest.NewRequest(http.MethodPost, "/api/privacy/execution-network/verify",
				strings.NewReader(`{"expected_exit_ip":" 1.1.1.1 "}`))
			w := httptest.NewRecorder()
			s.verifyExecutionNetwork(w, req)
			if w.Code != http.StatusOK {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "1.1.1.1") {
				t.Fatal("verification response exposed raw observed/expected IP")
			}
			got := applyModeBProbeSnapshot(executionNetworkOverview(os.Getenv))
			if got.Verification != tc.wantVerification || got.FinalExit != tc.wantFinal ||
				got.ActualEgress != tc.wantActual || got.ObservedAt == nil {
				t.Fatalf("unexpected verified status: %+v", got)
			}
			if got.FailClosed != "not_tested" || got.VerifyBeforeTask {
				t.Fatalf("an exit-IP comparison must not claim kill-switch or pre-task enforcement: %+v", got)
			}
			if len(got.CoveredProcesses) != 1 || got.CoveredProcesses[0] != "linux_docker_tool_executor_proxy_aware" {
				t.Fatalf("unexpected process coverage: %+v", got.CoveredProcesses)
			}
		})
	}
}

func TestVerifyExecutionNetworkRejectsInvalidOrUnconfiguredRequests(t *testing.T) {
	resetModeBProbeForTest()
	t.Cleanup(resetModeBProbeForTest)
	t.Setenv("ARTEX_RUNNER_SOCKET", "")
	t.Setenv("ARTEX_RUNNER_PROXY_ADDR", "")
	s := &Server{}
	for _, tc := range []struct {
		name, body string
		wantStatus int
	}{
		{"invalid ip", `{"expected_exit_ip":"127.0.0.1"}`, http.StatusBadRequest},
		{"private ip", `{"expected_exit_ip":"10.0.0.1"}`, http.StatusBadRequest},
		{"extra field", `{"expected_exit_ip":"1.1.1.1","url":"https://example.com"}`, http.StatusBadRequest},
		{"runner not configured", `{"expected_exit_ip":"1.1.1.1"}`, http.StatusConflict},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/privacy/execution-network/verify", strings.NewReader(tc.body))
			w := httptest.NewRecorder()
			s.verifyExecutionNetwork(w, req)
			if w.Code != tc.wantStatus {
				t.Fatalf("status=%d want=%d body=%s", w.Code, tc.wantStatus, w.Body.String())
			}
		})
	}
}
