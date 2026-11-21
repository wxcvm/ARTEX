package db

import (
	"math"
	"sync"
	"testing"
)

func TestBudgetJournalArithmeticGuard(t *testing.T) {
	if !budgetJournalFits(0, math.MaxInt64, 0, 2) {
		t.Fatal("unlimited cap should not apply a finite limit")
	}
	if budgetJournalFits(10, 8, 2, 1) || budgetJournalFits(10, 11, 0, 0) {
		t.Fatal("reservation exceeded a finite budget")
	}
	if !budgetJournalFits(10, 5, 2, 3) {
		t.Fatal("exact remaining allowance was rejected")
	}
	if _, err := budgetJournalAdd(math.MaxInt64, 1); err == nil {
		t.Fatal("overflow was accepted")
	}
	if _, err := budgetJournalAdd(3, -1); err == nil {
		t.Fatal("negative journal entry was accepted")
	}
}

func budgetJournalTestTask(t *testing.T) (*DB, int64) {
	t.Helper()
	d, err := Open(testDSN(t))
	if err != nil {
		t.Skipf("postgres unavailable (%v) — skipping", err)
	}
	task, err := d.CreateTask("budget journal fixture", "local integration", nil, 0, 0)
	if err != nil {
		d.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := d.DeleteTask(task.ID); err != nil {
			t.Error(err)
		}
		d.Close()
	})
	return d, task.ID
}

func TestBudgetJournalReserveSettleAndReopen(t *testing.T) {
	d, taskID := budgetJournalTestTask(t)
	const run = "epoch-one"
	created, err := d.OpenBudgetRun(taskID, run, 20, 3)
	if err != nil || created.State != "active" {
		t.Fatalf("create budget run: %+v %v", created, err)
	}
	a := BudgetCallReservation{TaskID: taskID, RunID: run, CallID: "a", IntentID: 12, Tokens: 8}
	if ok, err := d.ReserveBudgetCall(a); err != nil || !ok {
		t.Fatalf("first reserve: %v %v", ok, err)
	}
	if ok, err := d.ReserveBudgetCall(a); err != nil || ok {
		t.Fatalf("replayed reserve charged twice: %v %v", ok, err)
	}
	if ok, err := d.ReserveBudgetCall(BudgetCallReservation{TaskID: taskID, RunID: run, CallID: "a", IntentID: 12, Tokens: 7}); err == nil || ok {
		t.Fatalf("changed reservation was accepted: %v %v", ok, err)
	}
	b := BudgetCallReservation{TaskID: taskID, RunID: run, CallID: "b", Tokens: 11}
	if ok, err := d.ReserveBudgetCall(b); err != nil || !ok {
		t.Fatalf("second reserve: %v %v", ok, err)
	}
	if ok, err := d.ReserveBudgetCall(BudgetCallReservation{TaskID: taskID, RunID: run, CallID: "c", Tokens: 2}); err == nil || ok {
		t.Fatalf("over-allocation accepted: %v %v", ok, err)
	}
	if ok, err := d.SettleBudgetCall(taskID, run, "a", 9); err != nil || !ok {
		t.Fatalf("settle a: %v %v", ok, err)
	}
	if ok, err := d.SettleBudgetCall(taskID, run, "a", 9); err != nil || ok {
		t.Fatalf("duplicate settlement charged twice: %v %v", ok, err)
	}
	if ok, err := d.SettleBudgetCall(taskID, run, "a", 10); err == nil || ok {
		t.Fatalf("changed settlement was accepted: %v %v", ok, err)
	}
	if ok, err := d.ReserveBudgetCall(BudgetCallReservation{TaskID: taskID, RunID: run, CallID: "c", Tokens: 1}); err == nil || ok {
		t.Fatalf("spent+held exceeded cap without rejection: %v %v", ok, err)
	}
	if ok, err := d.SettleBudgetCall(taskID, run, "b", 5); err != nil || !ok {
		t.Fatalf("settle b: %v %v", ok, err)
	}
	c := BudgetCallReservation{TaskID: taskID, RunID: run, CallID: "c", Tokens: 5}
	if ok, err := d.ReserveBudgetCall(c); err != nil || !ok {
		t.Fatalf("reserve after refund: %v %v", ok, err)
	}
	if ok, err := d.SettleBudgetCall(taskID, run, "c", 0); err != nil || !ok {
		t.Fatalf("zero-token attempt must still charge a turn: %v %v", ok, err)
	}
	st, err := d.ReadBudgetRun(taskID, run)
	if err != nil || st.SpentTokens != 14 || st.SpentTurns != 3 || st.HeldTokens != 0 || st.HeldTurns != 0 {
		t.Fatalf("unexpected budget journal: %+v %v", st, err)
	}
	reopened, err := d.OpenBudgetRun(taskID, run, 20, 3)
	if err != nil || reopened != st {
		t.Fatalf("run was reset by re-open: %+v %v", reopened, err)
	}
	if _, err := d.OpenBudgetRun(taskID, run, 100, 3); err == nil {
		t.Fatal("reopening an epoch with different limits was allowed")
	}
	d2, err := Open(testDSN(t))
	if err != nil {
		t.Fatal(err)
	}
	defer d2.Close()
	afterRestart, err := d2.ReadBudgetRun(taskID, run)
	if err != nil || afterRestart != st {
		t.Fatalf("fresh db handle did not restore committed counters: %+v %v", afterRestart, err)
	}
	if ok, err := d2.ReserveBudgetCall(BudgetCallReservation{TaskID: taskID, RunID: run, CallID: "d", Tokens: 0}); err == nil || ok {
		t.Fatalf("a fourth turn exceeded the persisted cap: %v %v", ok, err)
	}
}

func TestBudgetJournalQuarantinePreservesUnknownCost(t *testing.T) {
	d, taskID := budgetJournalTestTask(t)
	const run = "crash-window"
	if _, err := d.OpenBudgetRun(taskID, run, 100, 5); err != nil {
		t.Fatal(err)
	}
	c := BudgetCallReservation{TaskID: taskID, RunID: run, CallID: "provider-request", Tokens: 60}
	if ok, err := d.ReserveBudgetCall(c); err != nil || !ok {
		t.Fatalf("reserve: %v %v", ok, err)
	}
	if n, err := d.QuarantineBudgetRun(taskID, run); err != nil || n != 1 {
		t.Fatalf("quarantine: n=%d err=%v", n, err)
	}
	if n, err := d.QuarantineBudgetRun(taskID, run); err != nil || n != 0 {
		t.Fatalf("quarantine must be idempotent: n=%d err=%v", n, err)
	}
	st, err := d.ReadBudgetRun(taskID, run)
	if err != nil || st.State != "recovery_required" || st.HeldTokens != 60 || st.HeldTurns != 1 {
		t.Fatalf("crash released an unknown reservation: %+v %v", st, err)
	}
	if ok, err := d.ReserveBudgetCall(BudgetCallReservation{TaskID: taskID, RunID: run, CallID: "next", Tokens: 1}); err == nil || ok {
		t.Fatalf("new work admitted during recovery: %v %v", ok, err)
	}
	if ok, err := d.SettleBudgetCall(taskID, run, c.CallID, 0); err == nil || ok {
		t.Fatalf("unknown provider response was silently charged as zero: %v %v", ok, err)
	}
	if ok, err := d.ReserveBudgetCall(c); err != nil || ok {
		t.Fatalf("exact replay of held call should not charge twice: %v %v", ok, err)
	}
}

func TestBudgetJournalConcurrentReservationCannotOverbook(t *testing.T) {
	d, taskID := budgetJournalTestTask(t)
	const run = "parallel-workers"
	if _, err := d.OpenBudgetRun(taskID, run, 50, 1); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	gate := make(chan struct{})
	results := make(chan bool, 2)
	errs := make(chan error, 2)
	for _, id := range []string{"worker-1", "worker-2"} {
		wg.Add(1)
		go func(callID string) {
			defer wg.Done()
			<-gate
			ok, err := d.ReserveBudgetCall(BudgetCallReservation{TaskID: taskID, RunID: run, CallID: callID, Tokens: 40})
			results <- ok
			errs <- err
		}(id)
	}
	close(gate)
	wg.Wait()
	close(results)
	close(errs)
	wins := 0
	for ok := range results {
		if ok {
			wins++
		}
	}
	for err := range errs {
		// The losing call is expected to hit the budget limit.
		if err == nil {
			continue
		}
	}
	st, err := d.ReadBudgetRun(taskID, run)
	if err != nil || wins != 1 || st.HeldTokens != 40 || st.HeldTurns != 1 {
		t.Fatalf("parallel reservations were not serialized: wins=%d state=%+v err=%v", wins, st, err)
	}
}
