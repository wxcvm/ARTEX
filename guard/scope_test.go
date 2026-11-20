package guard

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/Autumn-27/norma/hook"
)

// scopeWith builds a Guard whose source returns the given entries and never
// errors, installed in the given mode.
func scopeWith(mode ScopeMode, entries ...ScopeEntry) *Guard {
	g := New()
	g.UseScope(mode, ScopeSourceFunc(func() ([]ScopeEntry, error) { return entries, nil }))
	return g
}

// bash runs the PreToolUse hook for a shell command and reports whether it was
// blocked.
func bash(g *Guard, cmd string) bool {
	input, _ := json.Marshal(map[string]string{"command": cmd})
	blocked, _, _ := g.Hooks().PreToolUse(context.Background(), "Bash", input)
	return blocked
}

// With no source installed, or in ScopeOff, the gate is inert: this is the
// backward-compatibility guarantee for every task that declares no scope.
func TestScopeOffIsInert(t *testing.T) {
	entries := []ScopeEntry{{Kind: "cidr", Value: "10.0.0.0/24"}}

	noSource := New()
	if noSource.ScopeMode() != ScopeOff {
		t.Fatalf("a new Guard must default to ScopeOff, got %v", noSource.ScopeMode())
	}
	if bash(noSource, "curl http://198.51.100.7/") {
		t.Error("without a source nothing may be blocked")
	}

	off := scopeWith(ScopeOff, entries...)
	if bash(off, "curl http://198.51.100.7/") {
		t.Error("ScopeOff must not block")
	}
}

// A network scope admits exactly its own range and nothing else.
func TestCIDRScope(t *testing.T) {
	g := scopeWith(ScopeEnforce, ScopeEntry{Kind: "cidr", Value: "10.0.188.0/24"})

	for _, cmd := range []string{
		"curl http://10.0.188.5:3000/api",
		"nmap -sT -p 1-1000 10.0.188.200",
		"curl http://10.0.188.9/x",
	} {
		if bash(g, cmd) {
			t.Errorf("in-scope command was blocked: %q", cmd)
		}
	}
	for _, cmd := range []string{
		"curl http://10.0.189.5/",
		"nmap 192.0.2.10",
	} {
		if !bash(g, cmd) {
			t.Errorf("out-of-scope command was allowed: %q", cmd)
		}
	}
}

// A root_domain scope admits the apex and its subdomains, and must not admit a
// lookalike that merely ends with the same characters.
func TestRootDomainScopeBoundary(t *testing.T) {
	g := scopeWith(ScopeEnforce, ScopeEntry{Kind: "root_domain", Value: "example.com"})

	for _, cmd := range []string{
		"curl https://example.com/login",
		"curl https://a.b.example.com/x",
		"http://api.example.com:8443/v1",
	} {
		if bash(g, cmd) {
			t.Errorf("in-scope command was blocked: %q", cmd)
		}
	}
	// evil-example.com ends with "-example.com", not ".example.com"; a naive
	// strings.HasSuffix("example.com") check would wrongly admit it.
	for _, cmd := range []string{
		"curl https://evil-example.com/",
		"curl https://example.com.evil.net/",
	} {
		if !bash(g, cmd) {
			t.Errorf("out-of-scope lookalike was allowed: %q", cmd)
		}
	}
}

// A subdomain entry is an exact host: sibling hosts are out of scope.
func TestSubdomainScopeIsExact(t *testing.T) {
	g := scopeWith(ScopeEnforce, ScopeEntry{Kind: "subdomain", Value: "app.example.com"})
	if bash(g, "curl https://app.example.com/") {
		t.Error("exact host was blocked")
	}
	if !bash(g, "curl https://api.example.com/") {
		t.Error("sibling host must be out of scope")
	}
}

// Tooling has to keep working under a narrow scope: loopback, special-use names,
// and the infrastructure the agent downloads from are never targets.
func TestLocalAndInfrastructureNeverBlocked(t *testing.T) {
	g := scopeWith(ScopeEnforce, ScopeEntry{Kind: "domain", Value: "app.example.com"})

	for _, cmd := range []string{
		"curl http://127.0.0.1:8787/api/status",
		"curl http://localhost:5432/",
		"curl -L -o nuclei.zip https://github.com/projectdiscovery/nuclei/releases/download/v3.0.0/nuclei.zip",
		"go mod download proxy.golang.org",
		"curl https://nvd.nist.gov/vuln/detail/CVE-2024-0001",
		"curl http://[::1]:8080/",
	} {
		if bash(g, cmd) {
			t.Errorf("local/infrastructure command was blocked: %q", cmd)
		}
	}
}

// A scope lookup failure is distinct from an empty, successfully read scope.
func TestUnreadableScope(t *testing.T) {
	const internalError = "internal scope store details"
	for _, tc := range []struct {
		mode        ScopeMode
		wantBlocked bool
		wantActions []string
	}{
		{ScopeOff, false, []string{"allow"}},
		{ScopeWarn, false, []string{"warn", "allow"}},
		{ScopeEnforce, true, []string{"block"}},
	} {
		t.Run(tc.mode.String(), func(t *testing.T) {
			g := New()
			queried := false
			g.UseScope(tc.mode, ScopeSourceFunc(func() ([]ScopeEntry, error) {
				queried = true
				return nil, errors.New(internalError)
			}))

			input := json.RawMessage(`{"url":"https://app.example.com/"}`)
			blocked, msg, _ := g.Hooks().PreToolUse(context.Background(), "WebFetch", input)
			if blocked != tc.wantBlocked {
				t.Errorf("blocked = %v, want %v", blocked, tc.wantBlocked)
			}
			if queried != (tc.mode != ScopeOff) {
				t.Errorf("unexpected scope lookup state: queried = %v", queried)
			}
			if blocked && !strings.Contains(msg, "ARTEX platform control") {
				t.Errorf("block is missing platform framing: %q", msg)
			}

			audit := g.Audit()
			if len(audit) != len(tc.wantActions) {
				t.Fatalf("got audit %+v, want actions %v", audit, tc.wantActions)
			}
			for i, action := range tc.wantActions {
				entry := audit[i]
				if entry.Action != action {
					t.Errorf("audit[%d].Action = %q, want %q", i, entry.Action, action)
				}
				if action != "allow" && !strings.Contains(entry.Reason, "task scope could not be verified") {
					t.Errorf("audit[%d] is missing the lookup failure reason: %q", i, entry.Reason)
				}
				if strings.Contains(entry.Reason, internalError) {
					t.Errorf("audit[%d] exposes internal scope store details", i)
				}
			}
		})
	}
}

func TestScopeRefreshFailureBlocksUntilRecovery(t *testing.T) {
	g := New()
	unavailable := false
	g.UseScope(ScopeEnforce, ScopeSourceFunc(func() ([]ScopeEntry, error) {
		if unavailable {
			return nil, errors.New("scope store unavailable")
		}
		return []ScopeEntry{{Kind: "domain", Value: "app.example.com"}}, nil
	}))

	const cmd = "curl https://app.example.com/"
	if bash(g, cmd) {
		t.Fatal("call within the initial scope was blocked")
	}
	// Expire the successful snapshot without waiting for the real TTL.
	g.scopeCache.at = g.scopeCache.at.Add(-scopeCacheTTL)
	unavailable = true
	if !bash(g, cmd) {
		t.Error("call was allowed after the cached scope expired and refresh failed")
	}
	unavailable = false
	if bash(g, cmd) {
		t.Error("call within the verified scope was blocked after recovery")
	}
}

// Scope kinds with no host shape cannot decide a call. A company- or
// keyword-only scope must not block, because it cannot say what is outside.
func TestKeywordOnlyScopeIsNotEnforceable(t *testing.T) {
	g := scopeWith(ScopeEnforce,
		ScopeEntry{Kind: "company", Value: "Acme Holdings"},
		ScopeEntry{Kind: "keyword", Value: "fixture-keyword"},
	)
	if bash(g, "curl http://198.51.100.7/") {
		t.Error("a scope with no host or network boundary must not block")
	}
}

// Warn mode is the rollout mode: the violation is recorded, the call proceeds.
func TestWarnModeRecordsAndAllows(t *testing.T) {
	g := scopeWith(ScopeWarn, ScopeEntry{Kind: "cidr", Value: "10.0.188.0/24"})
	if bash(g, "curl http://198.51.100.7/") {
		t.Fatal("warn mode must not block")
	}
	var warns, allows int
	for _, entry := range g.Audit() {
		switch entry.Action {
		case "warn":
			warns++
			if !strings.Contains(entry.Reason, "198.51.100.7") {
				t.Errorf("warn reason must name the target, got %q", entry.Reason)
			}
		case "allow":
			allows++
		}
	}
	if warns != 1 {
		t.Errorf("got %d warn entries, want 1", warns)
	}
	if allows != 1 {
		t.Errorf("got %d allow entries, want 1 (the call still ran)", allows)
	}
}

// An enforced block must name the boundary and the offending target, and must
// carry the platform-control framing so the worker does not read it as a target
// defense it should try to evade.
func TestEnforcedBlockMessage(t *testing.T) {
	g := scopeWith(ScopeEnforce, ScopeEntry{Kind: "cidr", Value: "10.0.188.0/24"})
	input, _ := json.Marshal(map[string]string{"command": "curl http://198.51.100.7/"})
	blocked, msg, _ := g.Hooks().PreToolUse(context.Background(), "Bash", input)
	if !blocked {
		t.Fatal("out-of-scope call was allowed")
	}
	for _, want := range []string{"ARTEX platform control, not target defense", "198.51.100.7", "10.0.188.0/24", "add_task_scope"} {
		if !strings.Contains(msg, want) {
			t.Errorf("block message missing %q: %s", want, msg)
		}
	}
}

// A scope entry the DB stored in an unparseable shape is dropped, never treated
// as "allow everything": the call would otherwise be silently unrestricted.
func TestUnparseableEntryDoesNotWiden(t *testing.T) {
	g := scopeWith(ScopeEnforce,
		ScopeEntry{Kind: "cidr", Value: "not-a-cidr"},
		ScopeEntry{Kind: "domain", Value: "  "},
	)
	if bash(g, "curl http://198.51.100.7/") {
		t.Error("a scope with no parseable boundary must not block")
	}
}

// Extraction is the input to the gate, so its false-positive behavior is part of
// the contract: file names and module paths are not targets.
func TestTargetsFromInput(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{"url", `curl -k https://api.example.com/v1/login`, []string{"api.example.com"}},
		{"ipv4", `nmap -Pn -p 80 10.0.188.9`, []string{"10.0.188.9"}},
		{"hostport", `curl http://10.0.188.9:3000/x`, []string{"10.0.188.9"}},
		{"userinfo", `curl https://user:pass@example.com/`, []string{"example.com"}},
		{"files are not hosts", `hashcat -m 16500 jwt.txt rockyou.txt -r best64.rule`, nil},
		{"module path", `go build ./agent/...`, nil},
		{"bare label", `ls -la`, nil},
		{"nested", `curl https://a.b.example.com/x --resolve example.com:443:10.0.0.1`, []string{"a.b.example.com", "10.0.0.1", "example.com"}},
		{"version string", `kubectl v1.31.4 -o json`, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := TargetsFromInput([]byte(tc.input))
			if len(got) != len(tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Fatalf("got %v, want %v", got, tc.want)
				}
			}
		})
	}
}

// The compiled scope must be recomputed when the source changes, so a mid-task
// add_task_scope widening takes effect on the next call rather than after a
// process restart. The cache TTL is bypassed by installing a new source.
func TestScopeFollowsSourceChange(t *testing.T) {
	entries := []ScopeEntry{{Kind: "cidr", Value: "10.0.188.0/24"}}
	g := New()
	g.UseScope(ScopeEnforce, ScopeSourceFunc(func() ([]ScopeEntry, error) { return entries, nil }))
	if !bash(g, "curl http://10.9.9.9/") {
		t.Fatal("10.9.9.9 must start out of scope")
	}
	entries = append(entries, ScopeEntry{Kind: "cidr", Value: "10.9.0.0/16"})
	g.UseScope(ScopeEnforce, ScopeSourceFunc(func() ([]ScopeEntry, error) { return entries, nil }))
	if bash(g, "curl http://10.9.9.9/") {
		t.Error("widened scope must be honored on the next call")
	}
}

// An intercept-layer block must win over the scope gate so the recorded reason
// stays the one the operator wrote, not a scope message.
func TestInterceptBlockTakesPrecedence(t *testing.T) {
	g := scopeWith(ScopeEnforce, ScopeEntry{Kind: "cidr", Value: "10.0.188.0/24"})
	g.recordToolDecision("Bash", "x", hook.Result{Decision: "block", Message: "operator rule"})
	audit := g.Audit()
	if len(audit) != 1 || audit[0].Reason != "operator rule" {
		t.Fatalf("unexpected audit: %+v", audit)
	}
}

// Gating is on the reach surface, not on content: a local write that merely
// quotes an out-of-scope URL must not be refused, while the same URL reached
// through a shell is.
func TestOnlyReachToolsAreGated(t *testing.T) {
	g := scopeWith(ScopeEnforce, ScopeEntry{Kind: "cidr", Value: "10.0.188.0/24"})

	input, _ := json.Marshal(map[string]string{"path": "report.md", "content": "see http://198.51.100.7/admin"})
	if blocked, _, _ := g.Hooks().PreToolUse(context.Background(), "write", input); blocked {
		t.Error("a local file write must not be gated by a URL it quotes")
	}
	if blocked, _, _ := g.Hooks().PreToolUse(context.Background(), "WebFetch", input); !blocked {
		t.Error("WebFetch to an out-of-scope host must be gated")
	}
}

// ParseScopeMode is the configuration surface; an unknown value must be reported
// rather than silently mapped to a default.
func TestParseScopeMode(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want ScopeMode
		ok   bool
	}{
		{"", ScopeOff, false},
		{"off", ScopeOff, true},
		{" WARN ", ScopeWarn, true},
		{"enforce", ScopeEnforce, true},
		{"on", ScopeEnforce, true},
		{"yes", ScopeOff, false},
	} {
		got, ok := ParseScopeMode(tc.in)
		if ok != tc.ok || (ok && got != tc.want) {
			t.Errorf("ParseScopeMode(%q) = (%v,%v), want (%v,%v)", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}
