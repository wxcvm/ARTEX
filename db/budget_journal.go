package db

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
)

// BudgetRun is one durable accounting epoch. A process restart must reuse the
// same run ID; an explicit terminal revival must start a new one.
type BudgetRun struct {
	TaskID, TokenLimit, TurnLimit int64
	RunID                        string
	SpentTokens, SpentTurns       int64
	HeldTokens, HeldTurns         int64
	State                        string
}

// BudgetCallReservation identifies an LLM attempt before it is sent. The
// caller reuses CallID when retrying the reservation write for the same HTTP attempt;
// a new outbound HTTP retry gets a new CallID. IntentID=0 means planner.
type BudgetCallReservation struct {
	TaskID, IntentID, Tokens int64
	OwnerEpoch              int64
	RunID, CallID            string
}

func validateBudgetIdentity(taskID int64, runID string) error {
	if taskID <= 0 || len(runID) < 1 || len(runID) > 128 || strings.TrimSpace(runID) != runID {
		return fmt.Errorf("invalid budget task/run identity")
	}
	return nil
}

func budgetJournalAdd(a, b int64) (int64, error) {
	if a < 0 || b < 0 || a > math.MaxInt64-b {
		return 0, errors.New("budget journal counter overflow or negative value")
	}
	return a + b, nil
}

func budgetJournalFits(limit, spent, held, requested int64) bool {
	if limit == 0 {
		return true // zero is the existing unlimited-budget convention
	}
	// A zero-token reservation must not reopen dispatch after exhaustion.
	// Unknown provider usage may exceed its reservation, so an exhausted
	// dimension must reject even an estimated-zero new attempt.
	return spent < limit && held < limit-spent && requested <= limit-spent-held
}

func scanBudgetRun(row *sql.Row) (BudgetRun, error) {
	var r BudgetRun
	err := row.Scan(&r.TaskID, &r.RunID, &r.TokenLimit, &r.TurnLimit,
		&r.SpentTokens, &r.SpentTurns, &r.HeldTokens, &r.HeldTurns, &r.State)
	return r, err
}

const budgetRunColumns = "task_id,run_id,token_limit,turn_limit,spent_tokens,spent_turns,held_tokens,held_turns,state"

// OpenBudgetRun inserts a run exactly once and never replaces its counters or
// limits on re-open. Engine uses the unclosed row as the durable task-owner epoch;
// process restarts must reuse it rather than minting a fresh run ID.
func (d *DB) OpenBudgetRun(taskID int64, runID string, tokenLimit, turnLimit int64) (BudgetRun, error) {
	if err := validateBudgetIdentity(taskID, runID); err != nil {
		return BudgetRun{}, err
	}
	if tokenLimit < 0 || turnLimit < 0 {
		return BudgetRun{}, errors.New("budget limits must not be negative")
	}
	tx, err := d.Begin()
	if err != nil {
		return BudgetRun{}, err
	}
	defer tx.Rollback() //nolint:errcheck
	_, err = tx.Exec(`INSERT INTO task_budget_runs(task_id,run_id,token_limit,turn_limit)
VALUES ($1,$2,$3,$4) ON CONFLICT (task_id,run_id) DO NOTHING`, taskID, runID, tokenLimit, turnLimit)
	if err != nil {
		return BudgetRun{}, err
	}
	var r BudgetRun
	err = tx.QueryRow(`SELECT `+budgetRunColumns+` FROM task_budget_runs
WHERE task_id=$1 AND run_id=$2 FOR UPDATE`, taskID, runID).
		Scan(&r.TaskID, &r.RunID, &r.TokenLimit, &r.TurnLimit,
			&r.SpentTokens, &r.SpentTurns, &r.HeldTokens, &r.HeldTurns, &r.State)
	if err != nil {
		return BudgetRun{}, err
	}
	if r.TokenLimit != tokenLimit || r.TurnLimit != turnLimit {
		return BudgetRun{}, fmt.Errorf("budget run %q already exists with different limits", runID)
	}
	if err := tx.Commit(); err != nil {
		return BudgetRun{}, err
	}
	return r, nil
}

// ReadBudgetRun reads the committed run even after reopening a DB connection.
// Unknown/reserved calls remain charged to HeldTokens/HeldTurns.
func (d *DB) ReadBudgetRun(taskID int64, runID string) (BudgetRun, error) {
	if err := validateBudgetIdentity(taskID, runID); err != nil {
		return BudgetRun{}, err
	}
	return scanBudgetRun(d.QueryRow(`SELECT `+budgetRunColumns+` FROM task_budget_runs
WHERE task_id=$1 AND run_id=$2`, taskID, runID))
}

// ReserveBudgetCall holds budget before any provider request. On an exact
// replay with the same call ID and arguments it returns created=false without
// charging a second time, even if that call has already been settled.
// Concurrent planner/worker reservations serialize on the budget-run row.
func (d *DB) ReserveBudgetCall(c BudgetCallReservation) (created bool, err error) {
	if err := validateBudgetIdentity(c.TaskID, c.RunID); err != nil {
		return false, err
	}
	if c.CallID == "" || len(c.CallID) > 128 || c.Tokens < 0 || c.IntentID < 0 {
		return false, errors.New("invalid budget call reservation")
	}
	tx, err := d.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback() //nolint:errcheck
	var limitTokens, limitTurns, spentTokens, spentTurns, heldTokens, heldTurns, ownerEpoch int64
	var leaseActive bool
	var state string
	err = tx.QueryRow(`SELECT token_limit,turn_limit,spent_tokens,spent_turns,held_tokens,held_turns,state,
owner_epoch,COALESCE(owner_until > clock_timestamp(),false)
FROM task_budget_runs WHERE task_id=$1 AND run_id=$2 FOR UPDATE`, c.TaskID, c.RunID).
		Scan(&limitTokens, &limitTurns, &spentTokens, &spentTurns, &heldTokens, &heldTurns, &state,
			&ownerEpoch, &leaseActive)
	if err != nil {
		return false, err
	}
	if err := budgetOwnerValid(c.OwnerEpoch, ownerEpoch, leaseActive); err != nil {
		return false, err
	}
	var oldTokens, oldIntent int64
	err = tx.QueryRow(`SELECT reserved_tokens,intent_id FROM task_budget_calls
WHERE task_id=$1 AND run_id=$2 AND call_id=$3`, c.TaskID, c.RunID, c.CallID).
		Scan(&oldTokens, &oldIntent)
	switch {
	case err == nil:
		if oldTokens != c.Tokens || oldIntent != c.IntentID {
			return false, errors.New("budget call ID replayed with different reservation")
		}
		return false, tx.Commit()
	case !errors.Is(err, sql.ErrNoRows):
		return false, err
	}
	if state != "active" {
		return false, fmt.Errorf("budget run is %s; new calls are blocked", state)
	}
	if !budgetJournalFits(limitTokens, spentTokens, heldTokens, c.Tokens) ||
		!budgetJournalFits(limitTurns, spentTurns, heldTurns, 1) {
		return false, errors.New("insufficient budget for a new call reservation")
	}
	nextHeldTokens, err := budgetJournalAdd(heldTokens, c.Tokens)
	if err != nil {
		return false, err
	}
	nextHeldTurns, err := budgetJournalAdd(heldTurns, 1)
	if err != nil {
		return false, err
	}
	if _, err = tx.Exec(`INSERT INTO task_budget_calls(task_id,run_id,call_id,intent_id,reserved_tokens)
VALUES ($1,$2,$3,$4,$5)`, c.TaskID, c.RunID, c.CallID, c.IntentID, c.Tokens); err != nil {
		return false, err
	}
	if _, err = tx.Exec(`UPDATE task_budget_runs SET held_tokens=$3,held_turns=$4,updated_at=now()
WHERE task_id=$1 AND run_id=$2`, c.TaskID, c.RunID, nextHeldTokens, nextHeldTurns); err != nil {
		return false, err
	}
	return true, tx.Commit()
}

// SettleBudgetCall confirms the measured token cost after the model attempt.
// It charges exactly one turn (including a zero-token failure). Actual tokens
// can exceed the reservation; this blocks future reservations if over limit.
// Failed/unknown calls must NEVER be settled as zero without evidence.
func (d *DB) SettleBudgetCall(taskID int64, runID, callID string, actualTokens, ownerEpoch int64) (changed bool, err error) {
	if err := validateBudgetIdentity(taskID, runID); err != nil {
		return false, err
	}
	if callID == "" || actualTokens < 0 {
		return false, errors.New("invalid budget call settlement")
	}
	tx, err := d.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback() //nolint:errcheck
	var spentTokens, spentTurns, heldTokens, heldTurns, actualEpoch int64
	var leaseActive bool
	err = tx.QueryRow(`SELECT spent_tokens,spent_turns,held_tokens,held_turns,
owner_epoch,COALESCE(owner_until > clock_timestamp(),false)
FROM task_budget_runs WHERE task_id=$1 AND run_id=$2 FOR UPDATE`, taskID, runID).
		Scan(&spentTokens, &spentTurns, &heldTokens, &heldTurns, &actualEpoch, &leaseActive)
	if err != nil {
		return false, err
	}
	if err := budgetOwnerValid(ownerEpoch, actualEpoch, leaseActive); err != nil {
		return false, err
	}
	var reserved int64
	var previous sql.NullInt64
	var state string
	err = tx.QueryRow(`SELECT reserved_tokens,actual_tokens,state FROM task_budget_calls
WHERE task_id=$1 AND run_id=$2 AND call_id=$3 FOR UPDATE`, taskID, runID, callID).
		Scan(&reserved, &previous, &state)
	if err != nil {
		return false, err
	}
	if state == "settled" {
		if !previous.Valid || previous.Int64 != actualTokens {
			return false, errors.New("settled call cannot be replayed with changed token count")
		}
		return false, tx.Commit()
	}
	if state != "reserved" {
		return false, errors.New("unknown budget call cannot be auto-settled")
	}
	if heldTokens < reserved || heldTurns < 1 {
		return false, errors.New("budget journal held-count inconsistency")
	}
	nextTokens, err := budgetJournalAdd(spentTokens, actualTokens)
	if err != nil {
		return false, err
	}
	nextTurns, err := budgetJournalAdd(spentTurns, 1)
	if err != nil {
		return false, err
	}
	if _, err = tx.Exec(`UPDATE task_budget_calls SET state='settled',actual_tokens=$4,settled_at=now()
WHERE task_id=$1 AND run_id=$2 AND call_id=$3`, taskID, runID, callID, actualTokens); err != nil {
		return false, err
	}
	if _, err = tx.Exec(`UPDATE task_budget_runs SET spent_tokens=$3,spent_turns=$4,
held_tokens=$5,held_turns=$6,updated_at=now() WHERE task_id=$1 AND run_id=$2`,
		taskID, runID, nextTokens, nextTurns, heldTokens-reserved, heldTurns-1); err != nil {
		return false, err
	}
	return true, tx.Commit()
}

// QuarantineBudgetRun marks all outstanding calls unknown and blocks new ones.
// This is a recovery action for a single owner after ensuring no old process
// can still settle; it intentionally DOES NOT release held reservations.
func (d *DB) QuarantineBudgetRun(taskID int64, runID string, ownerEpoch int64) (int64, error) {
	if err := validateBudgetIdentity(taskID, runID); err != nil {
		return 0, err
	}
	tx, err := d.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck
	var state string
	var actualEpoch int64
	var leaseActive bool
	if err := tx.QueryRow(`SELECT state,owner_epoch,COALESCE(owner_until > clock_timestamp(),false)
FROM task_budget_runs WHERE task_id=$1 AND run_id=$2 FOR UPDATE`, taskID, runID).
		Scan(&state, &actualEpoch, &leaseActive); err != nil {
		return 0, err
	}
	if err := budgetOwnerValid(ownerEpoch, actualEpoch, leaseActive); err != nil {
		return 0, err
	}
	if state == "closed" {
		return 0, errors.New("cannot quarantine a closed budget run")
	}
	result, err := tx.Exec(`UPDATE task_budget_calls SET state='unknown'
WHERE task_id=$1 AND run_id=$2 AND state='reserved'`, taskID, runID)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	var unresolved bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM task_budget_calls WHERE
task_id=$1 AND run_id=$2 AND state='unknown')`, taskID, runID).Scan(&unresolved)
	if err != nil {
		return 0, err
	}
	if unresolved {
		_, err = tx.Exec(`UPDATE task_budget_runs SET state='recovery_required',updated_at=now()
WHERE task_id=$1 AND run_id=$2`, taskID, runID)
		if err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return affected, nil
}

// CloseBudgetRun is the explicit terminal boundary for one accounting epoch.
// It is idempotent after successful closure. An unresolved reservation or a
// recovery-required run cannot be discarded to mint a fresh budget.
func (d *DB) CloseBudgetRun(taskID int64, runID string, ownerEpoch int64) error {
	if err := validateBudgetIdentity(taskID, runID); err != nil {
		return err
	}
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck
	var heldTokens, heldTurns, actualEpoch int64
	var state string
	var leaseActive bool
	if err := tx.QueryRow(`SELECT held_tokens,held_turns,state,owner_epoch,
COALESCE(owner_until > clock_timestamp(),false)
FROM task_budget_runs WHERE task_id=$1 AND run_id=$2 FOR UPDATE`, taskID, runID).
		Scan(&heldTokens, &heldTurns, &state, &actualEpoch, &leaseActive); err != nil {
		return err
	}
	if state == "closed" {
		// A previously committed close is idempotent; no old owner is given
		// renewed write authority by this acknowledgement.
		return tx.Commit()
	}
	if err := budgetOwnerValid(ownerEpoch, actualEpoch, leaseActive); err != nil {
		return err
	}
	if state != "active" || heldTokens != 0 || heldTurns != 0 {
		return errors.New("cannot close a run with unknown or outstanding budget calls")
	}
	var unresolved bool
	if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM task_budget_calls
WHERE task_id=$1 AND run_id=$2 AND state <> 'settled')`, taskID, runID).
		Scan(&unresolved); err != nil {
		return err
	}
	if unresolved {
		return errors.New("cannot close a run with unresolved calls")
	}
	if _, err := tx.Exec(`UPDATE task_budget_runs SET state='closed',updated_at=now()
WHERE task_id=$1 AND run_id=$2`, taskID, runID); err != nil {
		return err
	}
	return tx.Commit()
}

// MarkBudgetCallUnknown records an ambiguous provider outcome without discarding
// the held reservation. The current owner may continue only within the remaining
// budget; graceful release or process takeover promotes the run to recovery_required.
func (d *DB) MarkBudgetCallUnknown(taskID int64, runID, callID string, ownerEpoch int64) (changed bool, err error) {
	if err := validateBudgetIdentity(taskID, runID); err != nil {
		return false, err
	}
	if callID == "" || len(callID) > 128 {
		return false, errors.New("invalid budget call identity")
	}
	tx, err := d.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback() //nolint:errcheck

	var actualEpoch int64
	var leaseActive bool
	var runState string
	if err := tx.QueryRow(`SELECT owner_epoch,COALESCE(owner_until > clock_timestamp(),false),state
FROM task_budget_runs WHERE task_id=$1 AND run_id=$2 FOR UPDATE`, taskID, runID).
		Scan(&actualEpoch, &leaseActive, &runState); err != nil {
		return false, err
	}
	if err := budgetOwnerValid(ownerEpoch, actualEpoch, leaseActive); err != nil {
		return false, err
	}
	if runState == "closed" {
		return false, errors.New("cannot mark a call unknown in a closed budget run")
	}

	var state string
	if err := tx.QueryRow(`SELECT state FROM task_budget_calls
WHERE task_id=$1 AND run_id=$2 AND call_id=$3 FOR UPDATE`, taskID, runID, callID).Scan(&state); err != nil {
		return false, err
	}
	switch state {
	case "unknown":
		return false, tx.Commit()
	case "reserved":
		if _, err := tx.Exec(`UPDATE task_budget_calls SET state='unknown'
WHERE task_id=$1 AND run_id=$2 AND call_id=$3 AND state='reserved'`, taskID, runID, callID); err != nil {
			return false, err
		}
		return true, tx.Commit()
	default:
		return false, fmt.Errorf("cannot mark budget call in state %q unknown", state)
	}
}

