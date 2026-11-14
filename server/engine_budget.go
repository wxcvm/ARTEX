package server

import (
	"context"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/llmrec"
)

func (e *Engine) budgetFor(taskID string) *agent.BudgetLedger {
	if v, ok := e.budgets.Load(taskID); ok {
		return v.(*agent.BudgetLedger)
	}
	v, _ := e.budgets.LoadOrStore(taskID, agent.NewBudgetLedger(agent.BudgetsFromEnv()))
	return v.(*agent.BudgetLedger)
}

func (e *Engine) budgetBlocksClaim(taskID string) bool {
	l := e.budgetFor(taskID)
	return l != nil && !l.State().Dispatch
}

// budgetStop uses the same transition lock and deletion barrier as deadline
// settlement, so pause, deletion and the deadline coordinator cannot race it.
func (e *Engine) budgetStop(ctx context.Context, t *Task) bool {
	l := e.budgetFor(t.ID)
	if l == nil || l.State().Dispatch {
		return false
	}
	t.transitionMu.Lock()
	e.deleteMu.RLock()
	state := t.lifecycleSnapshot()
	admit := !e.IsDeleting(t.ID) && !isTerminalStatus(state.Status) &&
		!state.Paused && !state.Queued && !e.IsPaused(t.ID) && !e.isSettling(t.ID) && ctx.Err() == nil
	if admit && l.BeginWrapup() {
		e.settling.Store(t.ID, true)
	} else {
		admit = false
	}
	e.deleteMu.RUnlock()
	t.transitionMu.Unlock()
	if admit {
		e.settleTaskWithBudget(ctx, t, true)
	}
	return true
}

func (e *Engine) budgetRoundCtx(ctx context.Context, t *Task, wrapup bool) context.Context {
	l := e.budgetFor(t.ID)
	if l == nil {
		return ctx
	}
	ctx = agent.WithBudgetRound(ctx, agent.BudgetRound{Line: l.State().Line(), WrapUp: wrapup})
	return llmrec.WithUsageSink(ctx, l.Record)
}

func (e *Engine) budgetIntent(ctx context.Context, taskID string) (context.Context, *agent.IntentMeter) {
	l := e.budgetFor(taskID)
	if l == nil {
		return ctx, nil
	}
	m := l.StartIntent()
	return llmrec.WithUsageSink(agent.WithIntentMeter(ctx, m), m.Record), m
}
