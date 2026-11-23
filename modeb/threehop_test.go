package modeb

import (
    "os"
    "path/filepath"
    "strings"
    "testing"

    "gopkg.in/yaml.v3"
)

// The three-hop profile is opt-in. Validate the direction of the chain so
// the last node is really the exit and errors never pick a direct fallback.
func TestModeBVPS2ProfileHasNoFallbackOrCircularHop(t *testing.T) {
    data, err := os.ReadFile(filepath.Join("..", "network", "mode-b", "config.vps2.example.yaml"))
    if err != nil { t.Fatal(err) }
    var profile struct {
        Proxies []struct {
            Name string `yaml:"name"`
            Type string `yaml:"type"`
            Server string `yaml:"server"`
            Dialer string `yaml:"dialer-proxy"`
        } `yaml:"proxies"`
        Rules []string `yaml:"rules"`
        DNS struct {
            Nameservers []string `yaml:"nameserver"`
        } `yaml:"dns"`
    }
    if err := yaml.Unmarshal(data, &profile); err != nil { t.Fatal(err) }
    if len(profile.Proxies) != 3 {
        t.Fatalf("three explicit nodes are required, got %d", len(profile.Proxies))
    }
    nodes := map[string]struct { typ, server, dialer string }{}
    for _, p := range profile.Proxies {
        if p.Name == "" || strings.EqualFold(p.Name, "DIRECT") || strings.EqualFold(p.Dialer, "DIRECT") {
            t.Fatalf("unsupported direct/unnamed node: %+v", p)
        }
        if _, exists := nodes[p.Name]; exists { t.Fatalf("duplicate node: %s", p.Name) }
        nodes[p.Name] = struct{typ, server, dialer string}{p.Type, p.Server, p.Dialer}
    }
    if nodes["airport-entry"].typ != "socks5" || nodes["airport-entry"].dialer != "" {
        t.Fatalf("airport should be a SOCKS5 first hop: %+v", nodes["airport-entry"])
    }
    if nodes["vps1"].typ != "vless" || nodes["vps1"].dialer != "airport-entry" {
        t.Fatalf("VPS1 must dial via airport: %+v", nodes["vps1"])
    }
    if nodes["vps2"].typ != "vless" || nodes["vps2"].dialer != "vps1" {
        t.Fatalf("VPS2 must dial via VPS1: %+v", nodes["vps2"])
    }
    if nodes["vps1"].server == nodes["vps2"].server {
        t.Fatal("VPS2 template must not reuse the same server address as VPS1")
    }
    if len(profile.Rules) != 1 || profile.Rules[0] != "MATCH,vps2" {
        t.Fatalf("only the chained final exit is allowed: %+v", profile.Rules)
    }
    if len(profile.DNS.Nameservers) != 1 || !strings.HasSuffix(profile.DNS.Nameservers[0], "#vps2") {
        t.Fatalf("resolver must also be routed via VPS2: %+v", profile.DNS.Nameservers)
    }
}
