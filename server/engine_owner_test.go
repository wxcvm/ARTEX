package server

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
)

func TestAcquireTaskOwnerReusesEpochAndRejectsConcurrentRuntime(t *testing.T) {
	for _, key := range []string{"ARTEX_TASK_BUDGET_TOKENS", "ARTEX_TASK_BUDGET_TURNS", "ARTEX_INTENT_BUDGET_TOKENS", "ARTEX_INTENT_BUDGET_TURNS"} {
		t.Setenv(key, "")
	}
	dsn, _, err := db.DSN()
	if err != nil {
		t.Skipf("PostgreSQL is not configured: %v", err)
	}
	pg, err := db.Open(dsn)
	if err != nil {
		t.Skipf("PostgreSQL is unavailable: %v", err)
	}
	task, err := pg.CreateTask("durable owner fixture", "lease lifecycle", nil, 0, 0)
	if err != nil {
		pg.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = pg.DeleteTask(task.ID)
		_ = pg.Close()
	})

	e := NewEngine(&Manager{pg: pg})
	handle := &Task{ID: strconv.FormatInt(task.ID, 10)}
	first, err := e.acquireTaskOwner(handle)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.acquireTaskOwner(handle); !errors.Is(err, db.ErrBudgetOwnerActive) {
		t.Fatalf("second runtime acquired an active task lease: %v", err)
	}
	if err := first.release(); err != nil {
		t.Fatal(err)
	}
	second, err := e.acquireTaskOwner(handle)
	if err != nil {
		t.Fatal(err)
	}
	if second.epoch != first.epoch+1 {
		t.Fatalf("owner epoch=%d after release; want %d", second.epoch, first.epoch+1)
	}
	if err := second.release(); err != nil {
		t.Fatal(err)
	}
}

func TestRunFailsClosedWhenDurableOwnerCannotBeAcquired(t *testing.T) {
	conn, err := sql.Open("pgx", "postgres://localhost/artex_test")
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	pg := &db.DB{DB: conn}
	m := &Manager{pg: pg, tasks: map[string]*Task{}}
	e := NewEngine(m)
	task := &Task{ID: "1", Store: pg.Exploration(1), notify: make(chan struct{}, 1)}
	e.Run(context.Background(), task)
	if e.Started(task.ID) {
		t.Fatal("task started without acquiring its durable execution owner")
	}
	e.runtimeMu.Lock()
	_, exists := e.runtimes[task.ID]
	e.runtimeMu.Unlock()
	if exists {
		t.Fatal("failed owner acquisition left a task runtime behind")
	}
}


func TestEngineRuntimeReleasesDurableOwnerAfterDrain(t *testing.T) {
	for _, key := range []string{"ARTEX_TASK_BUDGET_TOKENS", "ARTEX_TASK_BUDGET_TURNS", "ARTEX_INTENT_BUDGET_TOKENS", "ARTEX_INTENT_BUDGET_TURNS"} {
		t.Setenv(key, "")
	}
	dsn, _, err := db.DSN()
	if err != nil {
		t.Skipf("PostgreSQL is not configured: %v", err)
	}
	pg, err := db.Open(dsn)
	if err != nil {
		t.Skipf("PostgreSQL is unavailable: %v", err)
	}
	taskRow, err := pg.CreateTask("engine owner drain fixture", "release owner after runtime drain", nil, 0, 0)
	if err != nil {
		pg.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = pg.DeleteTask(taskRow.ID)
		_ = pg.Close()
	})

	taskID := strconv.FormatInt(taskRow.ID, 10)
	task := &Task{
		ID: taskID, ExpID: taskRow.ExplorationID, Store: pg.Exploration(taskRow.ExplorationID),
		Description: taskRow.Description, Goal: taskRow.Goal, notify: make(chan struct{}, 1),
	}
	m := &Manager{pg: pg, tasks: map[string]*Task{taskID: task}}
	e := NewEngine(m)
	e.Run(context.Background(), task)
	if !e.Started(taskID) {
		t.Fatal("task runtime did not start after acquiring its durable owner")
	}
	value, ok := e.taskOwners.Load(taskID)
	if !ok {
		t.Fatal("task runtime did not retain its durable owner")
	}
	owner := value.(*durableTaskOwner)
	stopped := make(chan struct{})
	go func() {
		e.StopTask(taskID)
		close(stopped)
	}()
	select {
	case <-stopped:
	case <-time.After(5 * time.Second):
		t.Fatal("task runtime did not drain and release its owner")
	}
	next, err := pg.AcquireBudgetOwner(taskRow.ID, owner.runID, taskOwnerLeaseTTL)
	if err != nil {
		t.Fatalf("owner was not released after task drain: %v", err)
	}
	if next != owner.epoch+1 {
		t.Fatalf("next owner epoch=%d; want %d", next, owner.epoch+1)
	}
	if err := pg.ReleaseBudgetOwner(taskRow.ID, owner.runID, next); err != nil {
		t.Fatal(err)
	}
}
