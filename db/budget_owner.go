package db

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"
)

var (
	ErrBudgetOwnerActive = errors.New("another budget run owner has an unexpired lease")
	ErrBudgetLeaseLost   = errors.New("budget run owner lease expired or superseded")
)

// A budget owner is identified by a monotonically increasing database epoch.
// Every writer must present that epoch and hold an unexpired server-clock lease.
// A new process may claim an expired run; old owners then fail every write.
func budgetOwnerValid(want, actual int64, active bool) error {
	if want <= 0 || want != actual || !active {
		return ErrBudgetLeaseLost
	}
	return nil
}

func budgetLeaseMillis(ttl time.Duration) (int64, error) {
	if ttl < time.Second || ttl > time.Hour {
		return 0, errors.New("budget owner lease must be between 1 second and 1 hour")
	}
	millis := ttl.Milliseconds()
	if millis <= 0 {
		return 0, errors.New("invalid budget owner lease duration")
	}
	return millis, nil
}

// AcquireBudgetOwner obtains exclusive write authority for an existing run.
// Lease evaluation uses PostgreSQL's clock, not the process clock. Expired
// owners are fenced by incrementing owner_epoch under the budget-row lock.
// Any unresolved requests from the previous owner become UNKNOWN, retain
// their held amounts and force recovery_required before new dispatch.
func (d *DB) AcquireBudgetOwner(taskID int64, runID string, ttl time.Duration) (int64, error) {
	if err := validateBudgetIdentity(taskID, runID); err != nil {
		return 0, err
	}
	ms, err := budgetLeaseMillis(ttl)
	if err != nil {
		return 0, err
	}
	tx, err := d.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck
	var epoch int64
	var leaseActive bool
	var state string
	err = tx.QueryRow(`SELECT owner_epoch,COALESCE(owner_until > clock_timestamp(),false),state
FROM task_budget_runs WHERE task_id=$1 AND run_id=$2 FOR UPDATE`, taskID, runID).
		Scan(&epoch, &leaseActive, &state)
	if err != nil {
		return 0, err
	}
	if state == "closed" {
		return 0, errors.New("cannot acquire a closed budget run")
	}
	if leaseActive {
		return 0, ErrBudgetOwnerActive
	}
	if epoch == math.MaxInt64 {
		return 0, errors.New("budget owner epoch exhausted")
	}
	epoch++
	if _, err = tx.Exec(`UPDATE task_budget_runs
SET owner_epoch=$3,owner_until=clock_timestamp()+($4::bigint * interval '1 millisecond'),
updated_at=now() WHERE task_id=$1 AND run_id=$2`, taskID, runID, epoch, ms); err != nil {
		return 0, err
	}
	if _, err = tx.Exec(`UPDATE task_budget_calls SET state='unknown'
WHERE task_id=$1 AND run_id=$2 AND state='reserved'`, taskID, runID); err != nil {
		return 0, err
	}
	var unresolved bool
	if err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM task_budget_calls
WHERE task_id=$1 AND run_id=$2 AND state='unknown')`, taskID, runID).
		Scan(&unresolved); err != nil {
		return 0, err
	}
	if unresolved {
		if _, err = tx.Exec(`UPDATE task_budget_runs SET state='recovery_required',updated_at=now()
WHERE task_id=$1 AND run_id=$2`, taskID, runID); err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return epoch, nil
}

// RenewBudgetOwner extends the lease only for the current, still-live epoch.
// A caller that missed its deadline cannot resurrect its stale authority.
func (d *DB) RenewBudgetOwner(taskID int64, runID string, epoch int64, ttl time.Duration) error {
	if err := validateBudgetIdentity(taskID, runID); err != nil {
		return err
	}
	ms, err := budgetLeaseMillis(ttl)
	if err != nil {
		return err
	}
	if epoch <= 0 {
		return ErrBudgetLeaseLost
	}
	result, err := d.Exec(`UPDATE task_budget_runs
SET owner_until=clock_timestamp()+($4::bigint * interval '1 millisecond'),
updated_at=now()
WHERE task_id=$1 AND run_id=$2 AND owner_epoch=$3
AND owner_until > clock_timestamp() AND state <> 'closed'`, taskID, runID, epoch, ms)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrBudgetLeaseLost
	}
	return nil
}

// ReleaseBudgetOwner revokes this epoch's write authority during graceful
// shutdown. Any still-reserved calls are quarantined rather than forgotten.
// It never resets or deletes counters, and old writers cannot settle afterward.
func (d *DB) ReleaseBudgetOwner(taskID int64, runID string, epoch int64) error {
	if err := validateBudgetIdentity(taskID, runID); err != nil {
		return err
	}
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck
	var actual int64
	var active bool
	var state string
	err = tx.QueryRow(`SELECT owner_epoch,COALESCE(owner_until > clock_timestamp(),false),state
FROM task_budget_runs WHERE task_id=$1 AND run_id=$2 FOR UPDATE`, taskID, runID).
		Scan(&actual, &active, &state)
	if err != nil {
		return err
	}
	if err := budgetOwnerValid(epoch, actual, active); err != nil {
		return err
	}
	if state == "closed" {
		return fmt.Errorf("cannot release a closed budget run")
	}
	if _, err = tx.Exec(`UPDATE task_budget_calls SET state='unknown'
WHERE task_id=$1 AND run_id=$2 AND state='reserved'`, taskID, runID); err != nil {
		return err
	}
	var unresolved bool
	if err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM task_budget_calls
WHERE task_id=$1 AND run_id=$2 AND state='unknown')`, taskID, runID).
		Scan(&unresolved); err != nil {
		return err
	}
	if unresolved {
		state = "recovery_required"
	}
	if _, err = tx.Exec(`UPDATE task_budget_runs SET owner_until=NULL,state=$3,updated_at=now()
WHERE task_id=$1 AND run_id=$2`, taskID, runID, state); err != nil {
		return err
	}
	return tx.Commit()
}

// ReadUnclosedBudgetRun returns the unique active or recovery-required epoch.
// sql.ErrNoRows means no run exists; an unexpected DB failure is never treated
// as permission to create a fresh budget epoch.
func (d *DB) ReadUnclosedBudgetRun(taskID int64) (BudgetRun, error) {
	if taskID <= 0 {
		return BudgetRun{}, errors.New("invalid budget task identity")
	}
	return scanBudgetRun(d.QueryRow(`SELECT `+budgetRunColumns+`
FROM task_budget_runs WHERE task_id=$1 AND state <> 'closed'`, taskID))
}

// Keep a compile-time reference to sql.ErrNoRows for callers using this API.
// (The method returns the standard database/sql error unchanged.)
var _ = sql.ErrNoRows
