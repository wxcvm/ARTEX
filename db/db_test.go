package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type schemaExecFunc func(context.Context, string, ...any) (sql.Result, error)

func (f schemaExecFunc) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return f(ctx, query, args...)
}

func TestApplySchemaRetriesOnlyDeadlocks(t *testing.T) {
	deadlockAttempts := 0
	delays := []time.Duration{}
	err := applySchemaWithRetry(t.Context(), schemaExecFunc(func(context.Context, string, ...any) (sql.Result, error) {
		deadlockAttempts++
		if deadlockAttempts < 3 {
			return nil, &pgconn.PgError{Code: "40P01", Message: "deadlock detected"}
		}
		return nil, nil
	}), func(_ context.Context, delay time.Duration) error {
		delays = append(delays, delay)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if deadlockAttempts != 3 {
		t.Fatalf("schema attempts=%d, want 3", deadlockAttempts)
	}
	if len(delays) != 2 || delays[0] != schemaDeadlockRetryDelays[0] || delays[1] != schemaDeadlockRetryDelays[1] {
		t.Fatalf("schema retry delays=%v", delays)
	}

	nonRetryable := errors.New("permission denied")
	nonRetryableAttempts := 0
	err = applySchemaWithRetry(t.Context(), schemaExecFunc(func(context.Context, string, ...any) (sql.Result, error) {
		nonRetryableAttempts++
		return nil, nonRetryable
	}), func(context.Context, time.Duration) error {
		t.Fatal("non-deadlock error must not be retried")
		return nil
	})
	if !errors.Is(err, nonRetryable) || nonRetryableAttempts != 1 {
		t.Fatalf("non-retryable result: attempts=%d err=%v", nonRetryableAttempts, err)
	}
}

// testDSN returns the configured DSN, skipping the test when neither the env var
// nor a config file supplies one (DSN no longer has a built-in default).
func testDSN(t *testing.T) string {
	t.Helper()
	dsn, _, err := DSN()
	if err != nil {
		t.Skipf("no database config (%v) — skipping", err)
	}
	return dsn
}

// TestOpenSeed opens the live dev PG, applies schema, seeds, and verifies the
// builtin agents + their variable catalog exist. Skips if PG is unreachable.
func TestOpenSeed(t *testing.T) {
	d, err := Open(testDSN(t))
	if err != nil {
		t.Skipf("postgres unavailable (%v) — skipping", err)
	}
	defer d.Close()

	wantAgents := len(builtinAgents)
	var agents int
	if err := d.QueryRow(`SELECT count(*) FROM agents WHERE builtin`).Scan(&agents); err != nil {
		t.Fatal(err)
	}
	if agents != wantAgents {
		t.Fatalf("want %d builtin agents, got %d", wantAgents, agents)
	}

	// planner should have at least its seeded catalog vars
	wantPlannerVars := 0
	for _, a := range builtinAgents {
		if a.key == "planner" {
			wantPlannerVars = len(a.vars)
			break
		}
	}
	var plannerVars int
	if err := d.QueryRow(`
SELECT count(*) FROM agent_prompt_vars v
JOIN agents a ON a.id = v.agent_id
WHERE a.key = 'planner'`).Scan(&plannerVars); err != nil {
		t.Fatal(err)
	}
	if plannerVars < wantPlannerVars {
		t.Fatalf("want at least %d planner vars, got %d", wantPlannerVars, plannerVars)
	}

	// re-open must be idempotent (no duplicate agents)
	d2, err := Open(testDSN(t))
	if err != nil {
		t.Fatal(err)
	}
	defer d2.Close()
	if err := d2.QueryRow(`SELECT count(*) FROM agents WHERE builtin`).Scan(&agents); err != nil {
		t.Fatal(err)
	}
	if agents != wantAgents {
		t.Fatalf("after re-open want %d agents, got %d", wantAgents, agents)
	}
}

// TestEnsureDatabaseConcurrentCreation uses a uniquely named disposable database.
// It never creates or drops the database named by the configured DSN itself.
func TestEnsureDatabaseConcurrentCreation(t *testing.T) {
	u, err := url.Parse(testDSN(t))
	if err != nil || (u.Scheme != "postgres" && u.Scheme != "postgresql") {
		t.Skip("requires a PostgreSQL URL")
	}
	adminURL := *u
	adminURL.Path = "/postgres"
	adminURL.RawPath = ""
	admin, err := sql.Open("pgx", adminURL.String())
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	if err := admin.PingContext(ctx); err != nil {
		t.Skipf("maintenance database unavailable: %v", err)
	}
	var canCreate bool
	if err := admin.QueryRowContext(ctx, `SELECT rolcreatedb OR rolsuper FROM pg_roles WHERE rolname=current_user`).Scan(&canCreate); err != nil {
		t.Fatal(err)
	}
	if !canCreate {
		t.Skip("test requires permission to create a disposable database")
	}

	name := fmt.Sprintf("artex_create_test_%d", time.Now().UnixNano())
	u.Path = "/" + name
	u.RawPath = ""
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if _, err := admin.ExecContext(cleanup, "DROP DATABASE IF EXISTS "+pgx.Identifier{name}.Sanitize()); err != nil {
			t.Errorf("drop disposable database: %v", err)
		}
	}()

	const workers = 6
	start := make(chan struct{})
	results := make(chan error, workers)
	for i := 0; i < workers; i++ {
		go func() {
			<-start
			results <- ensureDatabase(u.String())
		}()
	}
	close(start)
	for i := 0; i < workers; i++ {
		if err := <-results; err != nil {
			t.Errorf("concurrent database initialization: %v", err)
		}
	}
	var count int
	if err := admin.QueryRowContext(ctx, `SELECT count(*) FROM pg_database WHERE datname=$1`, name).Scan(&count); err != nil || count != 1 {
		t.Fatalf("disposable database count=%d, err=%v", count, err)
	}
}

func TestApplySchemaRetryWaitHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	attempts := 0
	err := applySchemaWithRetry(ctx, schemaExecFunc(func(context.Context, string, ...any) (sql.Result, error) {
		attempts++
		return nil, &pgconn.PgError{Code: "40P01", Message: "deadlock detected"}
	}), func(ctx context.Context, _ time.Duration) error {
		cancel()
		return waitForSchemaRetry(ctx, time.Hour)
	})
	if !errors.Is(err, context.Canceled) || attempts != 1 {
		t.Fatalf("retry cancellation: attempts=%d err=%v", attempts, err)
	}
}

func TestOpenContextRejectsExpiredDeadline(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := OpenContext(ctx, "postgres://127.0.0.1/artex"); !errors.Is(err, context.Canceled) {
		t.Fatalf("OpenContext with canceled context: %v", err)
	}
}
