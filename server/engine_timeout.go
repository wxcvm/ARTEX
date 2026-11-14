package server

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/db"
)

// Task-level timeout coordinator (see docs/task-timeout-and-settle-design.md §4/§8.5).
// One timer goroutine per task with a timeout; user pauses suspend admission and
// resume credits the persisted pause interval onto DeadlineAt. When due it drives
// the ordered settle sequence:
//   ① settling → ② workers stop claiming new intents / ③ planner drops normal notifies
//   ④ wait for running workers to drain (bounded by grace) → ⑤ one final planner round → ⑥ set the terminal state (guarded)

const (
	settleDrainGrace     = 90 * time.Second // upper bound for running workers to finish gracefully; beyond it, hard cancel
	deadlinePollInterval = 2 * time.Second  // poll interval while the deadline is unstamped / the LLM is not ready
	deadlineMaxSleep     = 30 * time.Second // longest single sleep (allows periodic terminal-state re-checks)
)

// ---------- settling state ----------

func (e *Engine) isSettling(taskID string) bool {
	v, _ := e.settling.Load(taskID)
	b, _ := v.(bool)
	return b
}

// ---------- in-flight count (worker.Execute + planner.Plan), used for drain ----------

func (e *Engine) inflightCounter(taskID string) *int64 {
	v, _ := e.inflight.LoadOrStore(taskID, new(int64))
	return v.(*int64)
}

// beginTaskOperation atomically registers a task-owned operation unless deletion
// has already installed its barrier. The delete handler can therefore wait for
// inflight==0 without a check-then-start race recreating files after cleanup.
func (e *Engine) beginTaskOperation(taskID string) bool {
	e.deleteMu.RLock()
	defer e.deleteMu.RUnlock()
	if e.IsDeleting(taskID) {
		return false
	}
	atomic.AddInt64(e.inflightCounter(taskID), 1)
	return true
}

func (e *Engine) decInflight(taskID string) { atomic.AddInt64(e.inflightCounter(taskID), -1) }
func (e *Engine) inflightCount(taskID string) int64 {
	return atomic.LoadInt64(e.inflightCounter(taskID))
}

// ---------- deadline ----------

// taskDeadline reads the committed lifecycle mirror, never a second clock cache.
// Resume and timeout revival can change it after the first run in this process.
// 0 = no timeout / not yet stamped.
func (e *Engine) taskDeadline(t *Task) int64 {
	return t.lifecycleSnapshot().DeadlineAt
}

// resetTimeoutRevival clears only the per-run timeout state after PostgreSQL has
// atomically committed timeout -> running and reset first_run_at/deadline_at. The
// configured TimeoutSeconds remains on Task, so the next real Planner/Worker run
// stamps a fresh full budget. coordStarted is reset because the coordinator that
// produced the timeout has already completed (or is in its final return path).
func (e *Engine) resetTimeoutRevival(taskID string) {
	e.resetSettledRun(taskID)
	e.stamped.Delete(taskID)
	e.coordStarted.Delete(taskID)
}

func (e *Engine) resetSettledRun(taskID string) {
	e.settling.Delete(taskID)
	e.budgets.Delete(taskID)
}

// stampFirstRun records first_run_at + deadline_at on the FIRST real run (LLM ready)
// of a timeout task, once per process. No-op when the task has no timeout.
func (e *Engine) stampFirstRun(t *Task) {
	if t.TimeoutSeconds <= 0 {
		return
	}
	if _, loaded := e.stamped.LoadOrStore(t.ID, true); loaded {
		return
	}
	dl, err := e.m.StampTaskFirstRun(t.ID)
	if err != nil {
		log.Printf("[deadline] task %s failed to stamp first_run: %v", t.ID, err)
		e.stamped.Delete(t.ID) // allow the next attempt to retry
		return
	}
	if dl > 0 {
		log.Printf("[deadline] task %s first run, deadline at %s", t.ID, time.Unix(dl, 0).Format("2006-01-02 15:04:05"))
	}
}

// clockCtx layers the task's TaskClock (absolute deadline) onto a run's context so
// worker/planner can clamp their wall-clock budget and pick per-run vs task-timeout
// wrap-up words. final marks the coordinator-driven terminal planner round.
func (e *Engine) clockCtx(base context.Context, t *Task, final bool) context.Context {
	dl := e.taskDeadline(t)
	if dl <= 0 && !final {
		return base // no timeout → unchanged behavior
	}
	return agent.WithTaskClock(base, agent.TaskClock{DeadlineUnix: dl, Final: final})
}

// ---------- coordinator ----------

// startDeadlineCoordinator launches the per-task deadline timer once (idempotent).
// Called from Run() and from the restart reload path, so non-active timeout tasks
// still get settled after their deadline even without live planner/worker loops.
func (e *Engine) startDeadlineCoordinator(ctx context.Context, t *Task) {
	if t == nil || t.TimeoutSeconds <= 0 {
		return
	}
	e.deleteMu.RLock()
	if e.IsDeleting(t.ID) {
		e.deleteMu.RUnlock()
		return
	}
	if _, loaded := e.coordStarted.LoadOrStore(t.ID, true); loaded {
		e.deleteMu.RUnlock()
		return
	}
	rt := e.registerTaskRoutines(ctx, t.ID, 1)
	e.deleteMu.RUnlock()
	runTaskRoutine(rt, func(loopCtx context.Context) { e.deadlineCoordinator(loopCtx, t) })
}

// deadlineAction is decided while holding the same task transition lock as a
// user pause and deadline updates. A deferred admission is NOT a terminal exit.
type deadlineAction uint8

const (
	deadlineWait deadlineAction = iota
	deadlineStop
	deadlineSettle
)

func (e *Engine) nextDeadlineAction(t *Task, now time.Time) (deadlineAction, time.Duration) {
	t.transitionMu.Lock()
	defer t.transitionMu.Unlock()
	e.deleteMu.RLock()
	defer e.deleteMu.RUnlock()
	state := t.lifecycleSnapshot()
	if e.IsDeleting(t.ID) || isTerminalStatus(state.Status) || e.isSettling(t.ID) {
		return deadlineStop, 0
	}
	// Queue-held time keeps its existing semantics. A persisted user pause or a
	// non-queue execution barrier (including admission in progress) defers settle.
	if state.Paused || (e.IsPaused(t.ID) && !state.Queued) || state.DeadlineAt <= 0 {
		return deadlineWait, deadlinePollInterval
	}
	if remaining := time.Unix(state.DeadlineAt, 0).Sub(now); remaining > 0 {
		return deadlineWait, min(remaining, deadlineMaxSleep)
	}
	// The transition lock makes check-and-claim atomic with pause commit/rollback.
	e.settling.Store(t.ID, true)
	return deadlineSettle, 0
}

// deadlineCoordinator uses production time/sleep; the callback seam lets the
// coordinator's actual wait/resume behavior be tested with a virtual clock and
// no model, database, or externally executed work.
func (e *Engine) deadlineCoordinator(ctx context.Context, t *Task) {
	e.coordinateDeadline(ctx, t, e.settleTask)
}

func (e *Engine) coordinateDeadline(ctx context.Context, t *Task, settle func(context.Context, *Task)) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		action, nap := e.nextDeadlineAction(t, time.Now())
		switch action {
		case deadlineStop:
			return
		case deadlineWait:
			if sleepCtx(ctx, nap) {
				return
			}
		case deadlineSettle:
			if ctx.Err() == nil {
				settle(ctx, t)
			}
			return
		}
	}
}

// settleTask runs the ordered settle sequence after nextDeadlineAction has
// exclusively claimed it. Do not call without that admission.
func (e *Engine) settleTask(ctx context.Context, t *Task) {
	e.settleTaskWithBudget(ctx, t, false)
}

func (e *Engine) settleTaskWithBudget(ctx context.Context, t *Task, budget bool) {
	if ctx.Err() != nil || e.IsDeleting(t.ID) {
		return
	}
	if budget {
		log.Printf("[budget] task %s reached its run budget, entering settlement", t.ID)
	} else {
		log.Printf("[deadline] task %s reached the timeout limit, entering the settle sequence", t.ID)
	}

	// ④ wait for running workers/planners to drain (running runs enter wrap-up on their own
	// via the clamped MaxDuration); if not drained within the grace period, hard-cancel the
	// task's exec ctx (the settling-aware branches classify it correctly).
	hardStop := time.Now().Add(settleDrainGrace)
	for e.inflightCount(t.ID) > 0 {
		if time.Now().After(hardStop) {
			log.Printf("[deadline] task %s drain timed out (%s), hard-cancelling running runs", t.ID, settleDrainGrace)
			e.cancelExec(t.ID, agent.AbortSettleDrainTimeout)
			_ = sleepCtx(ctx, 3*time.Second) // give the worker branches a moment to persist/classify
			break
		}
		if sleepCtx(ctx, 500*time.Millisecond) {
			return // whole engine shutting down
		}
	}

	// ⑤ one final planner round (task-timeout words, final goal judgment, no new intents).
	met := false
	if !budget {
		met = e.runFinalPlannerRound(ctx, t)
	} else if open, err := t.Store.HasOpenGoal(); err == nil && open {
		met = e.runFinalPlannerRoundWithBudget(ctx, t, true)
	}
	if ctx.Err() != nil {
		return // shutdown is not a completed final assessment
	}
	if !e.beginTaskOperation(t.ID) {
		return
	}
	defer e.decInflight(t.ID)

	// ⑥ set the terminal state (guarded): met → done (completed); otherwise timeout. If the
	// normal path already settled done first, the guard (SetTaskStatusGuarded) rejects the
	// overwrite, preserving the completed semantics.
	status := "timeout"
	if met && !budget {
		status = "done"
	}
	if budget {
		e.emitActivity(t, db.Activity{Worker: "system", Kind: "text", Summary: "Task run budget reached; final assessment ended by resource limit."})
	}
	won, err := e.m.SetTaskStatusGuarded(t.ID, status)
	switch {
	case err != nil:
		log.Printf("[deadline] task %s failed to settle the terminal state: %v", t.ID, err)
	case won:
		log.Printf("[deadline] task %s settle complete, terminal state=%s", t.ID, status)
	default:
		log.Printf("[deadline] task %s was already terminal when settling, keeping the original state", t.ID)
	}
}

// runFinalPlannerRound drives exactly ONE terminal planner round with the
// task-timeout planner words (final goal judgment; no new intents). Waits for the
// LLM to be ready (bounded by ctx) so a completable task isn't mis-judged timeout.
func (e *Engine) runFinalPlannerRound(ctx context.Context, t *Task) (met bool) {
	return e.runFinalPlannerRoundWithBudget(ctx, t, false)
}

func (e *Engine) runFinalPlannerRoundWithBudget(ctx context.Context, t *Task, budget bool) (met bool) {
	if e.IsDeleting(t.ID) {
		return false
	}
	planner, _ := e.snapshotFor(t)
	for planner == nil {
		if sleepCtx(ctx, deadlinePollInterval) {
			return false
		}
		if isTerminalStatus(e.m.TaskStatus(t.ID)) {
			return false
		}
		if e.IsDeleting(t.ID) {
			return false
		}
		planner, _ = e.snapshotFor(t)
	}
	// Independent ctx (not wired to execCancel, so a pause/hard cancel cannot interrupt this
	// final round), with Final injecting the task-timeout words.
	fctx := e.clockCtx(ctx, t, true)
	label := "task timeout settle"
	if budget {
		fctx = e.budgetRoundCtx(fctx, t, true)
		label = "task budget settle"
	}
	if !e.beginTaskOperation(t.ID) {
		return false
	}
	defer e.decInflight(t.ID)
	emit := func(r db.Activity) { e.emitActivity(t, r) }
	e.emitActivity(t, db.Activity{Worker: "planner", Kind: "round",
		Summary: fmt.Sprintf("%s · final judgment (round %d)", label, e.nextPlannerRound(t.ID))})
	tTaskID, _ := strconv.ParseInt(t.ID, 10, 64)
	e.BeginLLMCall(t.ID)
	met, reason, err := planner.Plan(fctx, tTaskID, e.m.assets, t.Store, t.Goal, t.drainTriggers(), emit)
	e.EndLLMCall(t.ID)
	if err != nil {
		log.Printf("[deadline] task %s final planning error: %v", t.ID, err)
	} else if met {
		log.Printf("[deadline] task %s final judgment: goal met: %s", t.ID, reason)
	}
	return met
}
