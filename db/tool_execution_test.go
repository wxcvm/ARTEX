package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func toolExecutionFixture(taskID, epoch int64, runID, sessionID, toolUseID, input string) ToolExecutionRequest {
	return ToolExecutionRequest{
		TaskID: taskID, RunID: runID, OwnerEpoch: epoch, IntentID: 42,
		SessionID: sessionID, ToolUseID: toolUseID, ToolKey: "insert_assets",
		AgentKey: "worker", InputHash: fmt.Sprintf("%064x", inputHashForTest(input)),
	}
}

func inputHashForTest(input string) uint64 {
	var hash uint64 = 1469598103934665603
	for i := 0; i < len(input); i++ {
		hash ^= uint64(input[i])
		hash *= 1099511628211
	}
	return hash
}

func TestToolExecutionJournalReplaysOnlyCompletedMatchingCalls(t *testing.T) {
	d, taskID := budgetJournalTestTask(t)
	const run = "tool-replay"
	if _, err := d.OpenBudgetRun(taskID, run, 100, 10); err != nil {
		t.Fatal(err)
	}
	epoch, err := d.AcquireBudgetOwner(taskID, run, 10*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	req := toolExecutionFixture(taskID, epoch, run, "session-1", "tool-use-1", "asset-a")
	if result, replay, err := d.BeginToolExecution(req); err != nil || replay || len(result) != 0 {
		t.Fatalf("first begin: result=%s replay=%v err=%v", result, replay, err)
	}
	stored := json.RawMessage(`{"content":[{"type":"text","text":"inserted"}]}`)
	if err := d.CompleteToolExecution(req, stored); err != nil {
		t.Fatalf("complete: %v", err)
	}
	reformatted := json.RawMessage(`{ "content" : [ { "text" : "inserted", "type" : "text" } ] }`)
	if err := d.CompleteToolExecution(req, reformatted); err != nil {
		t.Fatalf("idempotent complete replay with JSONB normalization: %v", err)
	}
	result, replay, err := d.BeginToolExecution(req)
	var got, want any
	if err == nil {
		err = json.Unmarshal(result, &got)
	}
	if err == nil {
		err = json.Unmarshal(stored, &want)
	}
	if err != nil || !replay || !reflect.DeepEqual(got, want) {
		t.Fatalf("completed call was not replayed semantically: result=%s replay=%v err=%v", result, replay, err)
	}
	conflict := req
	conflict.InputHash = fmt.Sprintf("%064x", inputHashForTest("asset-b"))
	if _, _, err := d.BeginToolExecution(conflict); !errors.Is(err, ErrToolExecutionConflict) {
		t.Fatalf("same tool_use_id with different input was not rejected: %v", err)
	}
}

func TestToolExecutionJournalQuarantinesUnknownOutcomeImmediately(t *testing.T) {
	d, taskID := budgetJournalTestTask(t)
	const run = "tool-unknown"
	if _, err := d.OpenBudgetRun(taskID, run, 100, 10); err != nil {
		t.Fatal(err)
	}
	epoch, err := d.AcquireBudgetOwner(taskID, run, 10*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	req := toolExecutionFixture(taskID, epoch, run, "session-2", "tool-use-2", "asset-c")
	if _, replay, err := d.BeginToolExecution(req); err != nil || replay {
		t.Fatalf("begin: replay=%v err=%v", replay, err)
	}
	if err := d.MarkToolExecutionUnknown(req, "side effect outcome not confirmed"); err != nil {
		t.Fatalf("mark unknown: %v", err)
	}
	if _, _, err := d.BeginToolExecution(req); !errors.Is(err, ErrToolExecutionUnknown) {
		t.Fatalf("unknown tool operation was replayed: %v", err)
	}
	if _, err := d.ReserveBudgetCall(BudgetCallReservation{
		TaskID: taskID, RunID: run, OwnerEpoch: epoch, CallID: "after-unknown", Tokens: 1,
	}); err == nil {
		t.Fatal("recovery-required run accepted a new provider call")
	}
	state, err := d.ReadBudgetRun(taskID, run)
	if err != nil || state.State != "recovery_required" {
		t.Fatalf("unknown tool outcome did not fence the run: %+v %v", state, err)
	}
}

func TestToolExecutionJournalTakeoverQuarantinesStaleRunningCall(t *testing.T) {
	d, taskID := budgetJournalTestTask(t)
	const run = "tool-takeover"
	if _, err := d.OpenBudgetRun(taskID, run, 100, 10); err != nil {
		t.Fatal(err)
	}
	oldEpoch, err := d.AcquireBudgetOwner(taskID, run, 10*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	req := toolExecutionFixture(taskID, oldEpoch, run, "session-3", "tool-use-3", "asset-d")
	if _, replay, err := d.BeginToolExecution(req); err != nil || replay {
		t.Fatalf("begin: replay=%v err=%v", replay, err)
	}
	if _, err := d.Exec(`UPDATE task_budget_runs SET owner_until=clock_timestamp()-interval '1 second'
WHERE task_id=$1 AND run_id=$2`, taskID, run); err != nil {
		t.Fatal(err)
	}
	newEpoch, err := d.AcquireBudgetOwner(taskID, run, 10*time.Minute)
	if err != nil || newEpoch <= oldEpoch {
		t.Fatalf("takeover epoch=%d old=%d err=%v", newEpoch, oldEpoch, err)
	}
	state, err := d.ReadBudgetRun(taskID, run)
	if err != nil || state.State != "recovery_required" {
		t.Fatalf("stale running tool did not require recovery: %+v %v", state, err)
	}
	replayReq := toolExecutionFixture(taskID, newEpoch, run, "session-3", "tool-use-3", "asset-d")
	if _, _, err := d.BeginToolExecution(replayReq); !errors.Is(err, ErrToolExecutionUnknown) {
		t.Fatalf("stale operation was automatically replayed: %v", err)
	}
}

func TestToolExecutionJournalGracefulReleaseQuarantinesRunningCall(t *testing.T) {
	d, taskID := budgetJournalTestTask(t)
	const run = "tool-graceful-release"
	if _, err := d.OpenBudgetRun(taskID, run, 100, 10); err != nil {
		t.Fatal(err)
	}
	epoch, err := d.AcquireBudgetOwner(taskID, run, 10*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	req := toolExecutionFixture(taskID, epoch, run, "session-4", "tool-use-4", "asset-e")
	if _, replay, err := d.BeginToolExecution(req); err != nil || replay {
		t.Fatalf("begin: replay=%v err=%v", replay, err)
	}
	if err := d.CloseBudgetRun(taskID, run, epoch); err == nil {
		t.Fatal("run closed while a tool side effect was still running")
	}
	if err := d.ReleaseBudgetOwner(taskID, run, epoch); err != nil {
		t.Fatalf("release: %v", err)
	}
	state, err := d.ReadBudgetRun(taskID, run)
	if err != nil || state.State != "recovery_required" {
		t.Fatalf("graceful release did not quarantine the running tool: %+v %v", state, err)
	}
	newEpoch, err := d.AcquireBudgetOwner(taskID, run, 10*time.Minute)
	if err != nil || newEpoch <= epoch {
		t.Fatalf("reacquire epoch=%d old=%d err=%v", newEpoch, epoch, err)
	}
	replayReq := toolExecutionFixture(taskID, newEpoch, run, "session-4", "tool-use-4", "asset-e")
	if _, _, err := d.BeginToolExecution(replayReq); !errors.Is(err, ErrToolExecutionUnknown) {
		t.Fatalf("released operation was automatically replayed: %v", err)
	}
}

