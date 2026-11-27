package db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrToolExecutionUnknown   = errors.New("tool execution outcome is unknown; automatic replay is blocked")
	ErrToolExecutionInProgress = errors.New("tool execution with this id is already in progress")
	ErrToolExecutionConflict  = errors.New("tool use id was reused with different tool input")
)

// sameToolResult compares JSON values rather than their byte representation.
// PostgreSQL JSONB normalizes whitespace and object key order on storage.
func sameToolResult(a, b []byte) bool {
	decode := func(raw []byte) (any, error) {
		decoder := json.NewDecoder(bytes.NewReader(raw))
		decoder.UseNumber()
		var value any
		if err := decoder.Decode(&value); err != nil {
			return nil, err
		}
		return value, nil
	}
	left, err := decode(a)
	if err != nil {
		return false
	}
	right, err := decode(b)
	if err != nil {
		return false
	}
	leftJSON, err := json.Marshal(left)
	if err != nil {
		return false
	}
	rightJSON, err := json.Marshal(right)
	if err != nil {
		return false
	}
	return bytes.Equal(leftJSON, rightJSON)
}

type ToolExecutionRequest struct {
	TaskID     int64
	RunID      string
	OwnerEpoch int64
	IntentID   int64
	SessionID  string
	ToolUseID  string
	ToolKey    string
	AgentKey   string
	InputHash  string
}

func validateToolExecutionRequest(req ToolExecutionRequest) error {
	if req.TaskID <= 0 || req.RunID == "" || req.OwnerEpoch <= 0 ||
		req.SessionID == "" || len(req.SessionID) > 256 ||
		req.ToolUseID == "" || len(req.ToolUseID) > 128 ||
		req.ToolKey == "" || len(req.ToolKey) > 256 ||
		len(req.InputHash) != 64 {
		return errors.New("invalid durable tool execution identity")
	}
	return nil
}

// BeginToolExecution claims one side-effecting tool call by its stable ToolUseID.
// A completed duplicate returns the stored result; running or unknown operations
// are never executed a second time.
func (d *DB) BeginToolExecution(req ToolExecutionRequest) (json.RawMessage, bool, error) {
	if err := validateToolExecutionRequest(req); err != nil {
		return nil, false, err
	}
	tx, err := d.Begin()
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback() //nolint:errcheck

	var actualEpoch int64
	var leaseActive bool
	var runState string
	if err := tx.QueryRow(`SELECT owner_epoch,COALESCE(owner_until > clock_timestamp(),false),state
FROM task_budget_runs WHERE task_id=$1 AND run_id=$2 FOR UPDATE`, req.TaskID, req.RunID).
		Scan(&actualEpoch, &leaseActive, &runState); err != nil {
		return nil, false, err
	}
	if err := budgetOwnerValid(req.OwnerEpoch, actualEpoch, leaseActive); err != nil {
		return nil, false, err
	}

	var existingRun, existingTool, existingHash, state string
	var existingEpoch int64
	var result []byte
	err = tx.QueryRow(`SELECT run_id,owner_epoch,tool_key,input_hash,state,result
FROM tool_executions WHERE task_id=$1 AND session_id=$2 AND tool_use_id=$3 FOR UPDATE`,
		req.TaskID, req.SessionID, req.ToolUseID).
		Scan(&existingRun, &existingEpoch, &existingTool, &existingHash, &state, &result)
	if err == nil {
		if existingTool != req.ToolKey || existingHash != req.InputHash {
			return nil, false, ErrToolExecutionConflict
		}
		switch state {
		case "done":
			if !json.Valid(result) {
				return nil, false, errors.New("completed tool execution has an invalid stored result")
			}
			if err := tx.Commit(); err != nil {
				return nil, false, err
			}
			return json.RawMessage(result), true, nil
		case "running":
			if existingRun == req.RunID && existingEpoch == req.OwnerEpoch {
				return nil, false, ErrToolExecutionInProgress
			}
			if _, err := tx.Exec(`UPDATE tool_executions SET state='unknown',
error='previous owner left the tool outcome unresolved',finished_at=now()
WHERE task_id=$1 AND session_id=$2 AND tool_use_id=$3 AND state='running'`,
				req.TaskID, req.SessionID, req.ToolUseID); err != nil {
				return nil, false, err
			}
			if _, err := tx.Exec(`UPDATE task_budget_runs SET state='recovery_required',updated_at=now()
WHERE task_id=$1 AND run_id=$2 AND state <> 'closed'`, req.TaskID, req.RunID); err != nil {
				return nil, false, err
			}
			if err := tx.Commit(); err != nil {
				return nil, false, err
			}
			return nil, false, ErrToolExecutionUnknown
		case "unknown":
			if _, err := tx.Exec(`UPDATE task_budget_runs SET state='recovery_required',updated_at=now()
WHERE task_id=$1 AND run_id=$2 AND state <> 'closed'`, req.TaskID, req.RunID); err != nil {
				return nil, false, err
			}
			if err := tx.Commit(); err != nil {
				return nil, false, err
			}
			return nil, false, ErrToolExecutionUnknown
		default:
			return nil, false, fmt.Errorf("unsupported tool execution state %q", state)
		}
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, false, err
	}
	if runState != "active" {
		return nil, false, errors.New("task budget run is not active; new tool side effects are blocked")
	}
	if _, err := tx.Exec(`INSERT INTO tool_executions
(task_id,run_id,owner_epoch,intent_id,session_id,tool_use_id,tool_key,agent_key,input_hash,state)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,'running')`,
		req.TaskID, req.RunID, req.OwnerEpoch, req.IntentID, req.SessionID,
		req.ToolUseID, req.ToolKey, req.AgentKey, req.InputHash); err != nil {
		return nil, false, err
	}
	if err := tx.Commit(); err != nil {
		return nil, false, err
	}
	return nil, false, nil
}

// CompleteToolExecution stores the exact tool result before it is returned to the
// agent. Replaying the same completion after a lost commit acknowledgement is safe.
func (d *DB) CompleteToolExecution(req ToolExecutionRequest, result json.RawMessage) error {
	if err := validateToolExecutionRequest(req); err != nil {
		return err
	}
	if !json.Valid(result) {
		return errors.New("tool execution result must be valid JSON")
	}
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	var actualEpoch int64
	var leaseActive bool
	if err := tx.QueryRow(`SELECT owner_epoch,COALESCE(owner_until > clock_timestamp(),false)
FROM task_budget_runs WHERE task_id=$1 AND run_id=$2 FOR UPDATE`, req.TaskID, req.RunID).
		Scan(&actualEpoch, &leaseActive); err != nil {
		return err
	}
	if err := budgetOwnerValid(req.OwnerEpoch, actualEpoch, leaseActive); err != nil {
		return err
	}

	var runID, toolKey, inputHash, state string
	var ownerEpoch int64
	var existing []byte
	if err := tx.QueryRow(`SELECT run_id,owner_epoch,tool_key,input_hash,state,result
FROM tool_executions WHERE task_id=$1 AND session_id=$2 AND tool_use_id=$3 FOR UPDATE`,
		req.TaskID, req.SessionID, req.ToolUseID).
		Scan(&runID, &ownerEpoch, &toolKey, &inputHash, &state, &existing); err != nil {
		return err
	}
	if toolKey != req.ToolKey || inputHash != req.InputHash {
		return ErrToolExecutionConflict
	}
	if state == "done" {
		if !sameToolResult(existing, result) {
			return ErrToolExecutionConflict
		}
		return tx.Commit()
	}
	if state != "running" || runID != req.RunID || ownerEpoch != req.OwnerEpoch {
		return ErrToolExecutionUnknown
	}
	if _, err := tx.Exec(`UPDATE tool_executions SET state='done',result=$4::jsonb,error='',
finished_at=now() WHERE task_id=$1 AND session_id=$2 AND tool_use_id=$3
AND state='running' AND run_id=$5 AND owner_epoch=$6`,
		req.TaskID, req.SessionID, req.ToolUseID, []byte(result), req.RunID, req.OwnerEpoch); err != nil {
		return err
	}
	return tx.Commit()
}

// MarkToolExecutionUnknown permanently blocks automatic replay of an operation
// whose side-effect outcome could not be confirmed. The current task run is moved
// to recovery_required so no further model or tool work can be dispatched.
func (d *DB) MarkToolExecutionUnknown(req ToolExecutionRequest, reason string) error {
	if err := validateToolExecutionRequest(req); err != nil {
		return err
	}
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	var actualEpoch int64
	var leaseActive bool
	if err := tx.QueryRow(`SELECT owner_epoch,COALESCE(owner_until > clock_timestamp(),false)
FROM task_budget_runs WHERE task_id=$1 AND run_id=$2 FOR UPDATE`, req.TaskID, req.RunID).
		Scan(&actualEpoch, &leaseActive); err != nil {
		return err
	}
	if err := budgetOwnerValid(req.OwnerEpoch, actualEpoch, leaseActive); err != nil {
		return err
	}
	var runID string
	var ownerEpoch int64
	var state string
	if err := tx.QueryRow(`SELECT run_id,owner_epoch,state FROM tool_executions
WHERE task_id=$1 AND session_id=$2 AND tool_use_id=$3 FOR UPDATE`,
		req.TaskID, req.SessionID, req.ToolUseID).
		Scan(&runID, &ownerEpoch, &state); err != nil {
		return err
	}
	if state == "done" {
		return tx.Commit()
	}
	if state == "unknown" {
		if _, err := tx.Exec(`UPDATE task_budget_runs SET state='recovery_required',updated_at=now()
WHERE task_id=$1 AND run_id=$2 AND state <> 'closed'`, req.TaskID, req.RunID); err != nil {
			return err
		}
		return tx.Commit()
	}
	if state != "running" || runID != req.RunID || ownerEpoch != req.OwnerEpoch {
		return ErrToolExecutionUnknown
	}
	reason = strings.TrimSpace(reason)
	if len(reason) > 2000 {
		reason = reason[:2000]
	}
	if reason == "" {
		reason = "tool execution ended without a durable result"
	}
	if _, err := tx.Exec(`UPDATE tool_executions SET state='unknown',error=$4,finished_at=now()
WHERE task_id=$1 AND session_id=$2 AND tool_use_id=$3 AND state='running'`,
		req.TaskID, req.SessionID, req.ToolUseID, reason); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE task_budget_runs SET state='recovery_required',updated_at=now()
WHERE task_id=$1 AND run_id=$2 AND state <> 'closed'`, req.TaskID, req.RunID); err != nil {
		return err
	}
	return tx.Commit()
}
