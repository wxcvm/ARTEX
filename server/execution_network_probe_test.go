package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/toolrunner"
)

func resetModeBProbeForTest() {
	modeBProbeState.Lock()
	modeBProbeState.snapshot = modeBProbeSnapshot{}
	modeBProbeState.Unlock()
}

func modeBTestEnv(configured bool) func(string) string {
	return func(key string) string {
		if !configured {
			return ""
		}
		switch key {
		case "ARTEX_RUNNER_SOCKET":
			return "/run/artex-runner/runner.sock"
		case "ARTEX_RUNNER_PROXY_ADDR":
			return "http://mode-b-egress:7890"
		default:
			return ""
		}
	}
}

func TestVerifyExecutionNetworkRecordsOnlyExitComparison(t *testing.T) {
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
			req := httptest.NewRequest(http.MethodPost, "/api/privacy/execution-network/verify",
				strings.NewReader(`{"expected_exit_ip":" 1.1.1.1 "}`))
			w := httptest.NewRecorder()
			probe := func(_ context.Context, expected string) (toolrunner.EgressProbeResponse, error) {
				if expected != "1.1.1.1" {
					t.Fatalf("expected normalized exit IP, got %q", expected)
				}
				return toolrunner.EgressProbeResponse{OK: tc.wantOK, Status: tc.runnerStatus}, nil
			}
			verifyExecutionNetworkWith(w, req, modeBTestEnv(true), probe)
			if w.Code != http.StatusOK {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "1.1.1.1") {
				t.Fatal("verification response exposed raw observed/expected IP")
			}
			got := applyModeBProbeSnapshot(executionNetworkOverview(modeBTestEnv(true)))
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
	s := &Server{}
	for _, tc := range []struct {
		name, body string
		getenv func(string) string
		wantStatus int
	}{
		{"invalid ip", `{"expected_exit_ip":"127.0.0.1"}`, modeBTestEnv(true), http.StatusBadRequest},
		{"private ip", `{"expected_exit_ip":"10.0.0.1"}`, modeBTestEnv(true), http.StatusBadRequest},
		{"extra field", `{"expected_exit_ip":"1.1.1.1","url":"https://example.com"}`, modeBTestEnv(true), http.StatusBadRequest},
		{"runner not configured", `{"expected_exit_ip":"1.1.1.1"}`, modeBTestEnv(false), http.StatusConflict},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/privacy/execution-network/verify", strings.NewReader(tc.body))
			w := httptest.NewRecorder()
			called := false
			verifyExecutionNetworkWith(w, req, tc.getenv, func(context.Context, string) (toolrunner.EgressProbeResponse, error) {
				called = true
				return toolrunner.EgressProbeResponse{OK: true, Status: "exit_ip_match"}, nil
			})
			if w.Code != tc.wantStatus {
				t.Fatalf("status=%d want=%d body=%s", w.Code, tc.wantStatus, w.Body.String())
			}
			if called {
				t.Fatal("probe must not run for invalid or unconfigured requests")
			}
		})
	}
}

func TestVerifyExecutionNetworkReportsProbeFailure(t *testing.T) {
	resetModeBProbeForTest()
	t.Cleanup(resetModeBProbeForTest)
	req := httptest.NewRequest(http.MethodPost, "/api/privacy/execution-network/verify",
		strings.NewReader(`{"expected_exit_ip":"1.1.1.1"}`))
	w := httptest.NewRecorder()
	verifyExecutionNetworkWith(w, req, modeBTestEnv(true), func(context.Context, string) (toolrunner.EgressProbeResponse, error) {
		return toolrunner.EgressProbeResponse{}, errors.New("private runner detail")
	})
	if w.Code != http.StatusBadGateway || strings.Contains(w.Body.String(), "private runner detail") {
		t.Fatalf("probe failure leaked details or returned wrong status: %d %s", w.Code, w.Body.String())
	}
	got := applyModeBProbeSnapshot(executionNetworkOverview(modeBTestEnv(true)))
	if got.Verification != "probe_failed" || got.FinalExit != "unavailable" ||
		got.ActualEgress != "unavailable" || got.ObservedAt == nil {
		t.Fatalf("probe failure was not recorded: %+v", got)
	}
}
