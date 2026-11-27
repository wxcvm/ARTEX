package modeb

import (
    "os"
    "path/filepath"
    "strings"
    "testing"

    "gopkg.in/yaml.v3"
)

type composeProfile struct {
    Services map[string]struct {
        Networks []string          `yaml:"networks"`
        Environment map[string]string `yaml:"environment"`
        Ports []any `yaml:"ports"`
        ReadOnly bool `yaml:"read_only"`
        Privileged bool `yaml:"privileged"`
        CapDrop []string `yaml:"cap_drop"`
        CapAdd []string `yaml:"cap_add"`
        SecurityOpt []string `yaml:"security_opt"`
        Pid string `yaml:"pid"`
        Tmpfs []string `yaml:"tmpfs"`
        MemLimit string `yaml:"mem_limit"`
        MemswapLimit string `yaml:"memswap_limit"`
        PidsLimit string `yaml:"pids_limit"`
    } `yaml:"services"`
    Networks map[string]struct {
        Internal bool `yaml:"internal"`
        Driver string `yaml:"driver"`
        DriverOpts map[string]string `yaml:"driver_opts"`
    } `yaml:"networks"`
}

func TestModeBExecutorHasOnlyIsolatedNetwork(t *testing.T) {
    data, err := os.ReadFile(filepath.Join("..", "docker-compose.mode-b.yml"))
    if err != nil { t.Fatal(err) }
    var profile composeProfile
    if err := yaml.Unmarshal(data, &profile); err != nil { t.Fatal(err) }
    executionNet := profile.Networks["runner-execution"]
    if !executionNet.Internal || executionNet.Driver != "bridge" {
        t.Fatal("runner-execution must be an internal Docker bridge")
    }
    for _, key := range []string{
        "com.docker.network.bridge.gateway_mode_ipv4",
        "com.docker.network.bridge.gateway_mode_ipv6",
    } {
        if executionNet.DriverOpts[key] != "isolated" {
            t.Fatalf("%s must isolate host bridge gateway, got %q", key, executionNet.DriverOpts[key])
        }
    }
    executor, ok := profile.Services["tool-executor"]
    if !ok || len(executor.Networks) != 1 || executor.Networks[0] != "runner-execution" {
        t.Fatalf("executor must have only isolated network: %+v", executor.Networks)
    }
    gateway, ok := profile.Services["mode-b-egress"]
    if !ok || len(gateway.Networks) != 2 ||
        gateway.Networks[0] != "runner-execution" || gateway.Networks[1] != "mode-b-uplink" {
        t.Fatalf("only the gateway may bridge to outside: %+v", gateway.Networks)
    }
    if len(gateway.Ports) > 0 {
        t.Fatal("gateway proxy port must never be published on the host")
    }
    if _, ok := profile.Services["artex-proxy-relay"]; ok {
        t.Fatal("recording relay is an uncontrolled second egress")
    }
    for _, key := range []string{"HTTP_PROXY", "HTTPS_PROXY", "ALL_PROXY"} {
        if !strings.Contains(executor.Environment[key], "mode-b-egress:7890") {
            t.Errorf("%s must use the isolated gateway, got %q", key, executor.Environment[key])
        }
    }
    if got := executor.Environment["ARTEX_RUNNER_PROXY_ADDR"]; got != "http://mode-b-egress:7890" {
        t.Fatalf("executor egress probe must use the private gateway, got %q", got)
    }
    if got := profile.Services["artex"].Environment["ARTEX_RUNNER_PROXY_ADDR"]; got != "http://mode-b-egress:7890" {
        t.Fatalf("controller tool proxy rewrite must target gateway: %q", got)
    }
}

func TestModeBExampleHasChainedExitWithoutDirectFallback(t *testing.T) {
    data, err := os.ReadFile(filepath.Join("..", "network", "mode-b", "config.example.yaml"))
    if err != nil { t.Fatal(err) }
    var profile struct {
        Proxies []struct {
            Name string `yaml:"name"`
            DialerProxy string `yaml:"dialer-proxy"`
        } `yaml:"proxies"`
        Rules []string `yaml:"rules"`
    }
    if err := yaml.Unmarshal(data, &profile); err != nil { t.Fatal(err) }
    byName := map[string]string{}
    for _, p := range profile.Proxies { byName[p.Name] = p.DialerProxy }
    if byName["vps1"] != "airport-entry" {
        t.Fatalf("VPS1 must dial through the airport: %+v", byName)
    }
    if len(profile.Rules) != 1 || profile.Rules[0] != "MATCH,vps1" {
        t.Fatalf("unexpected fallback or bypass rules: %+v", profile.Rules)
    }
}

func TestRunnerComposeProfilesKeepExecutorConstrained(t *testing.T) {
	for _, tc := range []struct {
		name string
		file string
		services []string
	}{
		{"runner", "docker-compose.runner.yml", []string{"tool-runner", "tool-executor", "artex-proxy-relay"}},
		{"mode-b", "docker-compose.mode-b.yml", []string{"tool-runner", "tool-executor", "mode-b-egress"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("..", tc.file))
			if err != nil { t.Fatal(err) }
			var profile composeProfile
			if err := yaml.Unmarshal(data, &profile); err != nil { t.Fatal(err) }
			for _, name := range tc.services {
				svc, ok := profile.Services[name]
				if !ok { t.Fatalf("missing service %q", name) }
				if !svc.ReadOnly { t.Errorf("%s must have a read-only root filesystem", name) }
				if svc.Privileged { t.Errorf("%s must never run privileged", name) }
				if svc.Pid == "host" { t.Errorf("%s must not share the host PID namespace", name) }
				if !containsString(svc.CapDrop, "ALL") { t.Errorf("%s must drop all capabilities by default: %v", name, svc.CapDrop) }
				if !containsString(svc.SecurityOpt, "no-new-privileges:true") {
					t.Errorf("%s must enforce no-new-privileges: %v", name, svc.SecurityOpt)
				}
			}
			executor := profile.Services["tool-executor"]
			if !containsString(executor.Tmpfs, "/tmp:rw,nosuid,size=1g") {
				t.Errorf("tool-executor must use a bounded nosuid /tmp tmpfs: %v", executor.Tmpfs)
			}
			for _, key := range []string{"ARTEX_PG_DSN", "ARTEX_JWT_KEY_DIR", "ARTEX_LLM_SECRET_KEY", "OPENAI_API_KEY", "ANTHROPIC_API_KEY", "ARTEX_SETUP_TOKEN"} {
				if _, ok := executor.Environment[key]; ok {
					t.Errorf("tool-executor must not receive controller secret %s", key)
				}
			}
		})
	}
}

func containsString(items []string, want string) bool {
	for _, item := range items { if item == want { return true } }
	return false
}

func TestDefaultComposeHasResourceCeilings(t *testing.T) {
    data, err := os.ReadFile(filepath.Join("..", "docker-compose.yml"))
    if err != nil { t.Fatal(err) }
    var profile composeProfile
    if err := yaml.Unmarshal(data, &profile); err != nil { t.Fatal(err) }
    app, ok := profile.Services["artex"]
    if !ok || app.MemLimit == "" || app.MemswapLimit != app.MemLimit || app.PidsLimit == "" {
        t.Fatalf("default ARTEX container must have memory, no-swap and PID ceilings: %+v", app)
    }
}
