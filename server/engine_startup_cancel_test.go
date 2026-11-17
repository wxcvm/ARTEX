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

func TestRunStartupFrontierRetryIsStoppedWithTask(t *testing.T) {
	// A closed sql.DB forces the startup HasActiveIntent query to fail without
	// requiring a live PostgreSQL server or a provider call.
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
	task := &Task{
		ID:     "1",
		Store:  pg.Exploration(1),
		notify: make(chan struct{}, 1),
	}
	m.tasks[task.ID] = task
	e.Run(context.Background(), task)
	if !e.Started(task.ID) {
		t.Fatal("task runtime was not registered")
	}
	e.runtimeMu.Lock()
	rt := e.runtimes[task.ID]
	e.runtimeMu.Unlock()
	if rt == nil {
		t.Fatal("task runtime was not recorded")
	}

	stopped := make(chan struct{})
	go func() {
		e.StopTask(task.ID)
		close(stopped)
	}()
	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("task shutdown waited for the startup retry timer")
	}
	if rt.ctx.Err() == nil || e.Started(task.ID) {
		t.Fatal("shutdown did not cancel and release the task runtime")
	}
	e.runtimeMu.Lock()
	_, retained := e.runtimes[task.ID]
	e.runtimeMu.Unlock()
	if retained {
		t.Fatal("shutdown retained a task runtime after stopping startup retry")
	}
}
