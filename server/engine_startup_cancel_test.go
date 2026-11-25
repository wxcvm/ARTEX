package server

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
)

func TestPlannerLoopCancellationInterruptsDebounce(t *testing.T) {
	e := NewEngine(nil)
	e.debounce = 5 * time.Second
	task := &Task{ID: "debounce", notify: make(chan struct{}, 1)}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		e.plannerLoop(ctx, task)
	}()
	task.Notify()
	// Allow the planner to enter its debounce wait, then cancel while its
	// normal five-second timer is still pending.
	time.Sleep(30 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("planner did not stop promptly when cancelled during debounce")
	}
}

func TestRunFailsClosedWhenCrashRecoveryCannotResetIntents(t *testing.T) {
	conn, err := sql.Open("pgx", "postgres://localhost/artex_test")
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	m := &Manager{tasks: map[string]*Task{}}
	e := NewEngine(m)
	task := &Task{ID: "1", Store: (&db.DB{DB: conn}).Exploration(1), notify: make(chan struct{}, 1)}
	e.Run(context.Background(), task)
	if e.Started(task.ID) {
		t.Fatal("task started even though crash recovery could not reset stale intents")
	}
	e.runtimeMu.Lock()
	_, exists := e.runtimes[task.ID]
	e.runtimeMu.Unlock()
	if exists {
		t.Fatal("failed crash recovery left a task runtime behind")
	}
}
