package server

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/db"
)

const (
	taskOwnerLeaseTTL     = 2 * time.Minute
	taskOwnerRenewEvery   = 30 * time.Second
)

type durableTaskOwner struct {
	store  *db.DB
	taskID int64
	runID  string
	epoch  int64
	ttl    time.Duration

	mu       sync.Mutex
	released bool
}

func journalBudgetLimit(value int64) int64 {
	if value == agent.Unlimited {
		return 0
	}
	return value
}

// acquireTaskOwner reuses the task's unclosed accounting epoch across process
// restarts. A fresh epoch is created only when none exists; changing a persisted
// cap or finding unresolved calls fails closed instead of resetting the budget.
func (e *Engine) acquireTaskOwner(t *Task) (*durableTaskOwner, error) {
	if e.m == nil || e.m.pg == nil {
		return nil, nil
	}
	if t == nil {
		return nil, errors.New("task is required for durable ownership")
	}
	taskID, err := strconv.ParseInt(t.ID, 10, 64)
	if err != nil || taskID <= 0 {
		return nil, fmt.Errorf("invalid task id %q for durable ownership", t.ID)
	}
	taskBudget, _ := agent.BudgetsFromEnv()
	tokenLimit, turnLimit := journalBudgetLimit(taskBudget.Tokens), journalBudgetLimit(taskBudget.Turns)

	run, err := e.m.pg.ReadUnclosedBudgetRun(taskID)
	if errors.Is(err, sql.ErrNoRows) {
		runID, tokenErr := newIntentClaimToken()
		if tokenErr != nil {
			return nil, tokenErr
		}
		run, err = e.m.pg.OpenBudgetRun(taskID, runID, tokenLimit, turnLimit)
		// A second process may have inserted the unique live epoch after our
		// initial read. Re-read that winner rather than minting another run ID.
		if err != nil {
			if existing, readErr := e.m.pg.ReadUnclosedBudgetRun(taskID); readErr == nil {
				run, err = existing, nil
			}
		}
	}
	if err != nil {
		return nil, fmt.Errorf("read or open durable task run: %w", err)
	}
	if run.TokenLimit != tokenLimit || run.TurnLimit != turnLimit {
		return nil, fmt.Errorf("task %s durable budget limits differ from current configuration (stored tokens/turns=%d/%d, configured=%d/%d)",
			t.ID, run.TokenLimit, run.TurnLimit, tokenLimit, turnLimit)
	}
	if run.State != "active" {
		return nil, fmt.Errorf("task %s has an unresolved durable run (%s); recovery is required before execution", t.ID, run.State)
	}

	epoch, err := e.m.pg.AcquireBudgetOwner(taskID, run.RunID, taskOwnerLeaseTTL)
	if err != nil {
		return nil, fmt.Errorf("acquire durable task owner: %w", err)
	}
	// Takeover may have converted an in-flight reservation to UNKNOWN. Do not
	// dispatch under a new owner until that ambiguity has been explicitly resolved.
	after, err := e.m.pg.ReadBudgetRun(taskID, run.RunID)
	if err != nil {
		_ = e.m.pg.ReleaseBudgetOwner(taskID, run.RunID, epoch)
		return nil, fmt.Errorf("verify durable task run after ownership claim: %w", err)
	}
	if after.State != "active" {
		_ = e.m.pg.ReleaseBudgetOwner(taskID, run.RunID, epoch)
		return nil, fmt.Errorf("task %s durable run entered %s during owner takeover; recovery is required", t.ID, after.State)
	}
	return &durableTaskOwner{
		store: e.m.pg, taskID: taskID, runID: run.RunID,
		epoch: epoch, ttl: taskOwnerLeaseTTL,
	}, nil
}

func (o *durableTaskOwner) renew() error {
	if o == nil {
		return nil
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.released {
		return db.ErrBudgetLeaseLost
	}
	return o.store.RenewBudgetOwner(o.taskID, o.runID, o.epoch, o.ttl)
}

func (o *durableTaskOwner) release() error {
	if o == nil {
		return nil
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.released {
		return nil
	}
	err := o.store.ReleaseBudgetOwner(o.taskID, o.runID, o.epoch)
	if err == nil || errors.Is(err, db.ErrBudgetLeaseLost) || errors.Is(err, sql.ErrNoRows) {
		o.released = true
	}
	return err
}

// maintainTaskOwner renews the database lease for the entire task runtime. On
// shutdown it waits for planner/worker/deadline routines to drain before release,
// so in-flight work is not quarantined merely because graceful shutdown began.
func (e *Engine) maintainTaskOwner(rt *taskRuntime, t *Task, owner *durableTaskOwner) {
	ticker := time.NewTicker(taskOwnerRenewEvery)
	defer ticker.Stop()
	lost := false

	for {
		select {
		case <-rt.ctx.Done():
			goto drain
		case <-ticker.C:
			if err := owner.renew(); err != nil {
				log.Printf("[engine] task %s lost durable execution lease: %v; cancelling local runtime", t.ID, err)
				lost = true
				rt.cancel()
				e.cancelExec(t.ID, err)
				goto drain
			}
		}
	}

drain:
	drained := make(chan struct{})
	go func() {
		rt.loops.Wait()
		close(drained)
	}()
	for {
		select {
		case <-drained:
			if !lost {
				if err := owner.release(); err != nil {
					log.Printf("[engine] task %s durable owner release failed: %v", t.ID, err)
				}
			}
			e.taskOwners.Delete(t.ID)
			return
		case <-ticker.C:
			if !lost {
				if err := owner.renew(); err != nil {
					log.Printf("[engine] task %s lost durable execution lease while draining: %v", t.ID, err)
					lost = true
					rt.cancel()
					e.cancelExec(t.ID, err)
				}
			}
		}
	}
}
