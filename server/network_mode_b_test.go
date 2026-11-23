package server

import "testing"

// The overview is intentionally not an egress probe. Even a fully configured
// chain must not be shown as verified, and node addresses must not be echoed.
func TestExecutionNetworkOverviewIsConfigurationOnly(t *testing.T) {
	cases := []struct {
		name string
		vars map[string]string
		mode, configuration, scope string
	}{
		{"default", nil, "off", "not_configured", "none"},
		{"socket_only", map[string]string{"ARTEX_RUNNER_SOCKET": "/run/artex-runner/runner.sock"}, "other", "partial", "runner_builtins_intended"},
		{"proxy_only", map[string]string{"ARTEX_RUNNER_PROXY_ADDR": "http://mode-b-egress:7890"}, "other", "partial", "none"},
		{"mode_b", map[string]string{
			"ARTEX_RUNNER_SOCKET": "/run/artex-runner/runner.sock",
			"ARTEX_RUNNER_PROXY_ADDR": "http://mode-b-egress:7890",
		}, "B", "configured_only", "runner_builtins_intended"},
		{"unknown_proxy", map[string]string{
			"ARTEX_RUNNER_SOCKET": "/run/artex-runner/runner.sock",
			"ARTEX_RUNNER_PROXY_ADDR": "http://untrusted.example:8080",
		}, "other", "partial", "runner_builtins_intended"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := executionNetworkOverview(func(key string) string { return tc.vars[key] })
			if got.Mode != tc.mode || got.Configuration != tc.configuration || got.Scope != tc.scope {
				t.Fatalf("unexpected summary: %+v", got)
			}
			if got.Verification != "not_tested" || got.FinalExit != "not_tested" {
				t.Fatalf("configured must not claim transit / final exit: %+v", got)
			}
		})
	}
}
