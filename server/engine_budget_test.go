package server

import (
	"context"
	"testing"
)

func TestRuntimeBudgetDisabledAndRevival(t *testing.T) {
	for _, key := range []string{"ARTEX_TASK_BUDGET_TOKENS", "ARTEX_TASK_BUDGET_TURNS", "ARTEX_INTENT_BUDGET_TOKENS", "ARTEX_INTENT_BUDGET_TURNS"} {
		t.Setenv(key, "")
	}
	e, task := NewEngine(nil), transitionTask()
	ctx := context.Background()
	if e.budgetBlocksClaim(task.ID) || e.budgetStop(ctx, task) || e.budgetRoundCtx(ctx, task, false) != ctx {
		t.Fatal("disabled budget changed task behavior")
	}
	if got, meter := e.budgetIntent(ctx, task.ID); got != ctx || meter != nil {
		t.Fatal("disabled budget changed worker behavior")
	}
	t.Setenv("ARTEX_TASK_BUDGET_TURNS", "1")
	e.resetSettledRun(task.ID)
	l := e.budgetFor(task.ID)
	if l == nil || !e.budgetBlocksClaim(task.ID) {
		t.Fatal("exhausted configured budget did not block dispatch")
	}
	e.paused.Store(task.ID, true)
	if e.admitBudgetSettlement(ctx, task, l) || e.isSettling(task.ID) {
		t.Fatal("paused task admitted budget settlement")
	}
	e.paused.Delete(task.ID)
	if !e.admitBudgetSettlement(ctx, task, l) || !e.isSettling(task.ID) {
		t.Fatal("budget settlement was not admitted")
	}
	if e.admitBudgetSettlement(ctx, task, l) {
		t.Fatal("duplicate settlement")
	}
	e.resetSettledRun(task.ID)
	if e.isSettling(task.ID) || e.budgetFor(task.ID) == l {
		t.Fatal("terminal revival retained old settlement or meter")
	}
}
