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
		return nil, errors.New("durable task execution owner is unavailable")
	}
	value, ok := r.s.engine.taskOwners.Load(r.taskID)
	if !ok {
		return nil, fmt.Errorf("task %s has no durable execution owner", r.taskID)
	}
	owner, ok := value.(*durableTaskOwner)
	if !ok || owner == nil {
		return nil, fmt.Errorf("task %s has an invalid durable execution owner", r.taskID)
	}
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
	created, err := j.owner.store.ReserveBudgetCall(db.BudgetCallReservation{
		TaskID: j.owner.taskID, RunID: j.owner.runID, OwnerEpoch: j.owner.epoch,
		CallID: callID, IntentID: j.intentID, Tokens: estimatedTokens,
	})
	if err != nil {
		return "", err
	}
	if !created {
		return "", errors.New("new provider call id already exists")
	}
	return callID, nil
}

func (j *taskBudgetAttemptJournal) Settle(ctx context.Context, callID string, actualTokens int64) error {
	if j == nil || j.owner == nil || j.owner.store == nil {
		return errors.New("durable task budget journal is unavailable")
	}
	if _, err := j.owner.store.SettleBudgetCall(j.owner.taskID, j.owner.runID, callID, actualTokens, j.owner.epoch); err != nil {
		return err
	}
	return nil
}

func (j *taskBudgetAttemptJournal) MarkUnknown(ctx context.Context, callID string) error {
	if j == nil || j.owner == nil || j.owner.store == nil {
		return errors.New("durable task budget journal is unavailable")
	}
	if _, err := j.owner.store.MarkBudgetCallUnknown(j.owner.taskID, j.owner.runID, callID, j.owner.epoch); err != nil {
		return err
	}
	return nil
}
