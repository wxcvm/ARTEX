package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/llmrec"
)

type taskBudgetAttemptJournal struct {
	owner    *durableTaskOwner
	intentID int64
}

func (r *taskLLMRuntime) withAttemptJournal(ctx context.Context) (context.Context, error) {
	if r == nil || r.s == nil || r.s.m == nil || r.s.m.pg == nil {
		return ctx, nil
	}
	if r.s.engine == nil {
		return nil, &taskLLMError{taskID: r.taskID, cause: fmt.Errorf("%w: durable task execution owner is unavailable", llmrec.ErrAttemptJournal)}
	}
	value, ok := r.s.engine.taskOwners.Load(r.taskID)
	if !ok {
		return nil, &taskLLMError{taskID: r.taskID, cause: fmt.Errorf("%w: task %s has no durable execution owner", llmrec.ErrAttemptJournal, r.taskID)}
	}
	owner, ok := value.(*durableTaskOwner)
	if !ok || owner == nil {
		return nil, &taskLLMError{taskID: r.taskID, cause: fmt.Errorf("%w: task %s has an invalid durable execution owner", llmrec.ErrAttemptJournal, r.taskID)}
	}
	ctx = llmrec.WithTaskOwner(ctx, llmrec.TaskOwner{
		TaskID: owner.taskID, RunID: owner.runID, Epoch: owner.epoch,
	})
	return llmrec.WithAttemptJournal(ctx, &taskBudgetAttemptJournal{
		owner: owner, intentID: agent.RunInfoFrom(ctx).IntentID,
	}), nil
}

func (j *taskBudgetAttemptJournal) Reserve(ctx context.Context, estimatedTokens int64) (string, error) {
	if j == nil || j.owner == nil || j.owner.store == nil {
		return "", errors.New("durable task budget journal is unavailable")
	}
	if estimatedTokens < 0 {
		return "", errors.New("provider request estimate must not be negative")
	}
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate provider call id: %w", err)
	}
	callID := hex.EncodeToString(raw[:])
	reservation := db.BudgetCallReservation{
		TaskID: j.owner.taskID, RunID: j.owner.runID, OwnerEpoch: j.owner.epoch,
		CallID: callID, IntentID: j.intentID, Tokens: estimatedTokens,
	}
	_, err := j.owner.store.ReserveBudgetCall(reservation)
	if err != nil {
		// The commit acknowledgement can be lost after PostgreSQL committed.
		// Replaying the same key distinguishes that case from a failed write.
		_, err = j.owner.store.ReserveBudgetCall(reservation)
	}
	if err != nil {
		return "", err
	}
	return callID, nil
}

func (j *taskBudgetAttemptJournal) Settle(ctx context.Context, callID string, actualTokens int64) error {
	if j == nil || j.owner == nil || j.owner.store == nil {
		return errors.New("durable task budget journal is unavailable")
	}
	settle := func() error {
		_, err := j.owner.store.SettleBudgetCall(j.owner.taskID, j.owner.runID, callID, actualTokens, j.owner.epoch)
		return err
	}
	if err := settle(); err != nil {
		return settle() // idempotent replay confirms a commit whose acknowledgement was lost.
	}
	return nil
}

func (j *taskBudgetAttemptJournal) MarkUnknown(ctx context.Context, callID string) error {
	if j == nil || j.owner == nil || j.owner.store == nil {
		return errors.New("durable task budget journal is unavailable")
	}
	mark := func() error {
		_, err := j.owner.store.MarkBudgetCallUnknown(j.owner.taskID, j.owner.runID, callID, j.owner.epoch)
		return err
	}
	if err := mark(); err != nil {
		return mark() // unknown transitions are idempotent across a lost commit acknowledgement.
	}
	return nil
}
