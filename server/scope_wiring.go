package server

import (
	"log"
	"os"
	"strings"

	pgdb "github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/guard"
)

// scopeModeFromEnv selects how the task-scope gate (guard/scope.go) is applied.
//
// ARTEX_SCOPE_MODE:
//
//	off     - do not consult the scope at all (the pre-scope behavior)
//	warn    - record an out-of-scope violation in the audit log, run the call
//	enforce - refuse the call
//
// The default is warn. It cannot break a task because it never refuses, and it
// starts producing, on real work, the evidence needed to decide whether enforce
// is safe. Anything unrecognized also falls back to warn with a log line, so a
// typo cannot silently disable the gate or silently start blocking.
func scopeModeFromEnv() guard.ScopeMode {
	raw, present := os.LookupEnv("ARTEX_SCOPE_MODE")
	if !present || strings.TrimSpace(raw) == "" {
		return guard.ScopeWarn
	}
	mode, valid := guard.ParseScopeMode(raw)
	if !valid {
		log.Printf("[scope] ARTEX_SCOPE_MODE=%q is not one of off|warn|enforce; using warn", raw)
		return guard.ScopeWarn
	}
	return mode
}

// taskScopeSource adapts a task's task_scope rows to guard.ScopeEntry. The
// planner can widen the scope mid-run through add_task_scope, so the guard reads
// this per call (with a short TTL) instead of snapshotting it here.
//
// A nil asset store or a non-positive task id yields a nil source, which leaves
// the guard in its inert configuration rather than pretending to enforce a scope
// that was never registered.
func taskScopeSource(assets *pgdb.AssetStore, taskID int64) guard.ScopeSource {
	if assets == nil || taskID <= 0 {
		return nil
	}
	return guard.ScopeSourceFunc(func() ([]guard.ScopeEntry, error) {
		rows, err := assets.ListTaskScope(taskID)
		if err != nil {
			return nil, err
		}
		out := make([]guard.ScopeEntry, 0, len(rows))
		for _, row := range rows {
			out = append(out, guard.ScopeEntry{Kind: row.Kind, Value: scopeEntryValue(row)})
		}
		return out, nil
	})
}

// scopeEntryValue picks the one column a task_scope row carries for its kind:
// domain/root_domain/subdomain live in Domain, cidr/ip in Net, and the free-text
// kinds in Value. Company rows carry only a reference, so their resolved name is
// the last resort — it has no host shape and guard treats it as non-enforceable.
func scopeEntryValue(row pgdb.TaskScope) string {
	for _, candidate := range []string{row.Domain, row.Net, row.Value} {
		if v := strings.TrimSpace(candidate); v != "" {
			return v
		}
	}
	return strings.TrimSpace(row.CompanyName)
}
