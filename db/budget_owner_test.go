package db

import (
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestBudgetOwnerFenceBlocksExpiredAndSupersededWriters(t *testing.T) {
	d, taskID := budgetJournalTestTask(t)
	const run = "owner-fencing"
	if _, err := d.OpenBudgetRun(taskID, run, 100, 5); err != nil {
		t.Fatal(err)
	}
	first, err := d.AcquireBudgetOwner(taskID, run, 5*time.Minute)
	if err != nil || first != 1 {
		t.Fatalf("first owner: epoch=%d err=%v", first, err)
	}
	if _, err := d.AcquireBudgetOwner(taskID, run, 5*time.Minute); !errors.Is(err, ErrBudgetOwnerActive) {
		t.Fatalf("a live owner was stolen: %v", err)
	}
	call := BudgetCallReservation{TaskID: taskID, RunID: run, CallID: "in-flight", Tokens: 30, OwnerEpoch: first}
	if ok, err := d.ReserveBudgetCall(call); err != nil || !ok {
		t.Fatalf("first owner could not reserve: ok=%v err=%v", ok, err)
	}
	// Simulate lease expiry without sleeping or relying on local clock skew.
	if _, err := d.Exec(`UPDATE task_budget_runs
SET owner_until=clock_timestamp()-interval '1 second' WHERE task_id=$1 AND run_id=$2`, taskID, run); err != nil {
		t.Fatal(err)
	}
	if err := d.RenewBudgetOwner(taskID, run, first, 5*time.Minute); !errors.Is(err, ErrBudgetLeaseLost) {
		t.Fatalf("expired owner renewed itself: %v", err)
	}
	if ok, err := d.SettleBudgetCall(taskID, run, call.CallID, 12, first); ok || !errors.Is(err, ErrBudgetLeaseLost) {
		t.Fatalf("expired owner settled a call: ok=%v err=%v", ok, err)
	}
	next, err := d.AcquireBudgetOwner(taskID, run, 5*time.Minute)
	if err != nil || next != first+1 {
		t.Fatalf("next owner: epoch=%d err=%v", next, err)
	}
	st, err := d.ReadBudgetRun(taskID, run)
	if err != nil || st.State != "recovery_required" || st.HeldTokens != 30 || st.HeldTurns != 1 {
		t.Fatalf("takeover erased an unknown provider call: %+v %v", st, err)
	}
	if ok, err := d.ReserveBudgetCall(BudgetCallReservation{
		TaskID: taskID, RunID: run, CallID: "fresh", Tokens: 1, OwnerEpoch: first,
	}); ok || !errors.Is(err, ErrBudgetLeaseLost) {
		t.Fatalf("superseded owner wrote through its fence: ok=%v err=%v", ok, err)
	}
	if ok, err := d.SettleBudgetCall(taskID, run, call.CallID, 0, next); ok || err == nil {
		t.Fatalf("new owner silently treated unknown response as free: ok=%v err=%v", ok, err)
	}
	if err := d.CloseBudgetRun(taskID, run, first); !errors.Is(err, ErrBudgetLeaseLost) {
		t.Fatalf("stale owner closed the epoch: %v", err)
	}
	if err := d.ReleaseBudgetOwner(taskID, run, next); err != nil {
		t.Fatalf("release quarantined owner: %v", err)
	}
	third, err := d.AcquireBudgetOwner(taskID, run, 5*time.Minute)
	if err != nil || third != next+1 {
		t.Fatalf("reacquire quarantined run: epoch=%d err=%v", third, err)
	}
	alive, err := d.ReadUnclosedBudgetRun(taskID)
	if err != nil || alive.State != "recovery_required" || alive.RunID != run {
		t.Fatalf("did not find the unclosed recovery epoch: %+v %v", alive, err)
	}
}

func TestBudgetOwnerGracefulReleaseKeepsCounters(t *testing.T) {
	d, taskID := budgetJournalTestTask(t)
	const run = "owner-graceful"
	if _, err := d.OpenBudgetRun(taskID, run, 100, 5); err != nil {
		t.Fatal(err)
	}
	first, err := d.AcquireBudgetOwner(taskID, run, 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	call := BudgetCallReservation{TaskID: taskID, RunID: run, CallID: "settled", Tokens: 20, OwnerEpoch: first}
	if ok, err := d.ReserveBudgetCall(call); !ok || err != nil {
		t.Fatalf("reserve: %v %v", ok, err)
	}
	if ok, err := d.SettleBudgetCall(taskID, run, call.CallID, 7, first); !ok || err != nil {
		t.Fatalf("settle: %v %v", ok, err)
	}
	if err := d.ReleaseBudgetOwner(taskID, run, first); err != nil {
		t.Fatal(err)
	}
	if err := d.RenewBudgetOwner(taskID, run, first, 5*time.Minute); !errors.Is(err, ErrBudgetLeaseLost) {
		t.Fatalf("released owner renewed: %v", err)
	}
	second, err := d.AcquireBudgetOwner(taskID, run, 5*time.Minute)
	if err != nil || second != first+1 {
		t.Fatalf("reclaim: epoch=%d err=%v", second, err)
	}
	st, err := d.ReadBudgetRun(taskID, run)
	if err != nil || st.State != "active" || st.SpentTokens != 7 || st.SpentTurns != 1 {
		t.Fatalf("graceful owner rotation reset budget: %+v %v", st, err)
	}
	if err := d.RenewBudgetOwner(taskID, run, second, 5*time.Minute); err != nil {
		t.Fatalf("current owner renew: %v", err)
	}
	if err := d.CloseBudgetRun(taskID, run, second); err != nil {
		t.Fatalf("final close: %v", err)
	}
	if _, err := d.ReadUnclosedBudgetRun(taskID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("closed run was returned as active: %v", err)
	}
}

func TestBudgetOwnerConcurrentClaimsHaveOneWinner(t *testing.T) {
	d, taskID := budgetJournalTestTask(t)
	const run = "owner-claim-race"
	if _, err := d.OpenBudgetRun(taskID, run, 100, 5); err != nil {
		t.Fatal(err)
	}
	secondDB, err := Open(testDSN(t))
	if err != nil {
		t.Fatal(err)
	}
	defer secondDB.Close()

	var wg sync.WaitGroup
	ready := make(chan struct{})
	out := make(chan error, 2)
	for _, client := range []*DB{d, secondDB} {
		wg.Add(1)
		go func(client *DB) {
			defer wg.Done()
			<-ready
			_, err := client.AcquireBudgetOwner(taskID, run, 5*time.Minute)
			out <- err
		}(client)
	}
	close(ready)
	wg.Wait()
	close(out)
	wins, alreadyOwned := 0, 0
	for err := range out {
		switch {
		case err == nil:
			wins++
		case errors.Is(err, ErrBudgetOwnerActive):
			alreadyOwned++
		default:
			t.Fatalf("unexpected concurrent claim failure: %v", err)
		}
	}
	if wins != 1 || alreadyOwned != 1 {
		t.Fatalf("concurrent owner fencing failed: wins=%d busy=%d", wins, alreadyOwned)
	}
}
