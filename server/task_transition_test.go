package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/Autumn-27/artex/agent"
)

func transitionTask() *Task {
	return &Task{ID: "42", Status: "running", DeadlineAt: 100, notify: make(chan struct{}, 1)}
}

func TestTaskTransitionDeadlineUsesCommittedSnapshot(t *testing.T) {
	e, task := NewEngine(nil), transitionTask()
	if got := e.taskDeadline(task); got != 100 {
		t.Fatal(got)
	}
	task.lifecycleMu.Lock()
	task.DeadlineAt = 200
	task.lifecycleMu.Unlock()
	if got := e.taskDeadline(task); got != 200 {
		t.Fatalf("stale deadline: %d, want 200", got)
	}
	if action, _ := e.nextDeadlineAction(task, time.Unix(150, 0)); action != deadlineWait {
		t.Fatal("settled using old deadline")
	}
}

func TestTaskTransitionPauseWins(t *testing.T) {
	e, task := NewEngine(nil), transitionTask()
	execCtx := e.execContextFor(context.Background(), task.ID)
	committing, release := make(chan struct{}), make(chan struct{})
	paused, timed := make(chan error, 1), make(chan deadlineAction, 1)
	go func() {
		paused <- e.pauseTaskWithCommit(task, agent.AbortPausedByUser, func() error {
			close(committing)
			<-release
			task.lifecycleMu.Lock()
			task.Paused = true
			task.lifecycleMu.Unlock()
			return nil
		})
	}()
	<-committing // pause owns the transition; timeout must wait for its commit
	go func() { a, _ := e.nextDeadlineAction(task, time.Unix(101, 0)); timed <- a }()
	if !e.IsPaused(task.ID) || execCtx.Err() != nil {
		t.Fatal("barrier must precede persistence, cancellation must follow it")
	}
	close(release)
	if err := <-paused; err != nil {
		t.Fatal(err)
	}
	if a := <-timed; a != deadlineWait || e.isSettling(task.ID) {
		t.Fatal("timeout bypassed committed pause")
	}
	if execCtx.Err() == nil {
		t.Fatal("successful pause did not cancel execution")
	}
}

func TestTaskTransitionSettlementWinsAndControlsConflict(t *testing.T) {
	e, task := NewEngine(nil), transitionTask()
	if a, _ := e.nextDeadlineAction(task, time.Unix(101, 0)); a != deadlineSettle {
		t.Fatal(a)
	}
	if err := e.pauseTaskWithCommit(task, nil, func() error { t.Fatal("must not persist after settlement admission"); return nil }); !errors.Is(err, errTaskSettling) {
		t.Fatal(err)
	}
	if e.IsPaused(task.ID) {
		t.Fatal("conflict installed a pause barrier")
	}
	s := &Server{engine: e, m: &Manager{tasks: map[string]*Task{task.ID: task}}}
	req := httptest.NewRequest("POST", "/api/tasks/42/control", strings.NewReader(`{"action":"pause"}`))
	req.SetPathValue("id", task.ID)
	w := httptest.NewRecorder()
	s.control(w, req)
	if w.Code != 409 || !strings.Contains(w.Body.String(), errTaskSettling.Error()) {
		t.Fatalf("single: %d %s", w.Code, w.Body)
	}
	w = httptest.NewRecorder()
	s.controlTasksBatch(w, httptest.NewRequest("POST", "/api/tasks/control", strings.NewReader(`{"action":"pause","task_ids":["42","042","bad","99"]}`)))
	var result struct {
		Items []batchControlItem `json:"items"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 || len(result.Items) != 3 || result.Items[0].OK || result.Items[0].Error != errTaskSettling.Error() || result.Items[1].Error == "" || result.Items[2].Error == "" {
		t.Fatalf("batch: %d %s", w.Code, w.Body)
	}
}

func TestTaskTransitionFailedCommitRestoresBarrier(t *testing.T) {
	for _, previouslyPaused := range []bool{false, true} {
		t.Run(map[bool]string{false: "running", true: "queue_barrier"}[previouslyPaused], func(t *testing.T) {
			e, task := NewEngine(nil), transitionTask()
			ctx := e.execContextFor(context.Background(), task.ID)
			if previouslyPaused {
				e.paused.Store(task.ID, true)
			}
			want := errors.New("injected database commit failure")
			err := e.pauseTaskWithCommit(task, agent.AbortPausedByUser, func() error { return want })
			if !errors.Is(err, want) || e.IsPaused(task.ID) != previouslyPaused || task.lifecycleSnapshot().Paused || ctx.Err() != nil {
				t.Fatalf("half-committed pause: err=%v barrier=%v ctx=%v", err, e.IsPaused(task.ID), ctx.Err())
			}
		})
	}
}

func TestTaskTransitionRollbackPreservesDeletion(t *testing.T) {
	e, task := NewEngine(nil), transitionTask()
	err := e.pauseTaskWithCommit(task, nil, func() error {
		if !e.BeginDelete(task.ID) {
			t.Fatal("delete was not admitted")
		}
		return errors.New("commit failed")
	})
	if err == nil || !e.IsDeleting(task.ID) || !e.IsPaused(task.ID) {
		t.Fatal("rollback removed delete barrier")
	}
	if a, _ := e.nextDeadlineAction(task, time.Unix(101, 0)); a != deadlineStop {
		t.Fatal("deleted task settled")
	}
}

func TestTaskTransitionCoordinatorWaitsAcrossPauseAndResume(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		e, task := NewEngine(nil), transitionTask()
		task.DeadlineAt = time.Now().Add(time.Second).Unix()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var calls atomic.Int64
		done := make(chan struct{})
		go func() { e.coordinateDeadline(ctx, task, func(context.Context, *Task) { calls.Add(1) }); close(done) }()
		synctest.Wait()
		if err := e.pauseTaskWithCommit(task, nil, func() error {
			task.lifecycleMu.Lock()
			defer task.lifecycleMu.Unlock()
			task.Paused = true
			return nil
		}); err != nil {
			t.Fatal(err)
		}
		time.Sleep(5 * time.Second) // virtual time only; never wall-clock sleeping
		synctest.Wait()
		select {
		case <-done:
			t.Fatal("coordinator exited while paused")
		default:
		}
		if calls.Load() != 0 {
			t.Fatal("paused task settled")
		}
		task.transitionMu.Lock()
		task.lifecycleMu.Lock()
		task.Paused = false
		task.DeadlineAt = time.Now().Add(4 * time.Second).Unix()
		task.lifecycleMu.Unlock()
		task.transitionMu.Unlock()
		e.Resume(task)
		time.Sleep(3 * time.Second)
		synctest.Wait()
		if calls.Load() != 0 {
			t.Fatal("resumed using old deadline")
		}
		time.Sleep(3 * time.Second)
		synctest.Wait()
		<-done
		if calls.Load() != 1 {
			t.Fatalf("settle calls=%d", calls.Load())
		}
		e.coordinateDeadline(ctx, task, func(context.Context, *Task) { calls.Add(1) })
		if calls.Load() != 1 {
			t.Fatal("duplicate settlement")
		}
	})
}

func TestTaskTransitionCoordinatorCancelDeleteTerminalAndQueue(t *testing.T) {
	for _, state := range []string{"cancel", "delete", "done", "queued"} {
		t.Run(state, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				e, task := NewEngine(nil), transitionTask()
				task.DeadlineAt = time.Now().Add(time.Second).Unix()
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				switch state {
				case "cancel":
					cancel()
				case "delete":
					e.BeginDelete(task.ID)
				case "done":
					task.Status = "done"
				case "queued":
					task.Queued = true
					e.paused.Store(task.ID, true)
				}
				calls := 0
				e.coordinateDeadline(ctx, task, func(context.Context, *Task) { calls++ })
				want := 0
				if state == "queued" {
					want = 1
				}
				if calls != want {
					t.Fatalf("%s: calls=%d want=%d", state, calls, want)
				}
			})
		})
	}
}

func TestTaskTransitionDeleteAfterAdmissionPreventsSettlementWork(t *testing.T) {
	e, task := NewEngine(nil), transitionTask()
	if a, _ := e.nextDeadlineAction(task, time.Unix(101, 0)); a != deadlineSettle {
		t.Fatal(a)
	}
	e.BeginDelete(task.ID)
	e.settleTask(context.Background(), task) // nil manager would panic if planner/DB work started
	if task.lifecycleSnapshot().Status != "running" {
		t.Fatal("deletion changed terminal status")
	}
}
