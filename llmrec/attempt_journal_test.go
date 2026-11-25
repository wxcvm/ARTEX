package llmrec

import (
	"context"
	"fmt"
	"iter"
	"sync"
	"testing"

	"github.com/Autumn-27/norma/llm"
)

type attemptJournalFixture struct {
	mu      sync.Mutex
	next    int
	settled map[string]int64
	unknown map[string]bool
}

func (j *attemptJournalFixture) Reserve(_ context.Context, _ int64) (string, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.next++
	return fmt.Sprintf("call-%d", j.next), nil
}

func (j *attemptJournalFixture) Settle(_ context.Context, id string, tokens int64) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.settled == nil {
		j.settled = map[string]int64{}
	}
	j.settled[id] = tokens
	return nil
}

func (j *attemptJournalFixture) MarkUnknown(_ context.Context, id string) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.unknown == nil {
		j.unknown = map[string]bool{}
	}
	j.unknown[id] = true
	return nil
}

type journaledProviderFixture struct{ journal *attemptJournalFixture }

func (p journaledProviderFixture) Stream(ctx context.Context, _ llm.CompletionRequest) iter.Seq2[llm.StreamEvent, error] {
	return func(yield func(llm.StreamEvent, error) bool) {
		attempt, err := BeginTransportAttempt(ctx)
		if err != nil {
			yield(llm.StreamEvent{}, err)
			return
		}
		if err := attempt.Complete(200, nil); err != nil {
			yield(llm.StreamEvent{}, err)
			return
		}
		if !yield(llm.StreamEvent{Type: llm.SEMessageStart, Usage: llm.Usage{InputTokens: 5}}, nil) {
			return
		}
		if !yield(llm.StreamEvent{Type: llm.SETextDelta, Text: "answer"}, nil) {
			return
		}
		if !yield(llm.StreamEvent{Type: llm.SEMessageDelta, Usage: llm.Usage{OutputTokens: 7}}, nil) {
			return
		}
		yield(llm.StreamEvent{Type: llm.SEMessageStop}, nil)
	}
}

func (p journaledProviderFixture) Complete(ctx context.Context, _ llm.CompletionRequest) (llm.Message, string, llm.Usage, error) {
	attempt, err := BeginTransportAttempt(ctx)
	if err != nil {
		return llm.Message{}, "", llm.Usage{}, err
	}
	if err := attempt.Complete(200, nil); err != nil {
		return llm.Message{}, "", llm.Usage{}, err
	}
	return llm.Message{Role: llm.RoleAssistant, Content: []llm.ContentBlock{llm.TextBlock("answer")}},
		"stop", llm.Usage{InputTokens: 11, OutputTokens: 3}, nil
}

func TestRecorderSettlesDurableHTTPAttemptsForStreamAndComplete(t *testing.T) {
	journal := &attemptJournalFixture{}
	recorder := Wrap(journaledProviderFixture{journal: journal}, nil, "fixture", "profile", "", "", func() bool { return false })
	ctx := WithAttemptJournal(context.Background(), journal)

	for _, err := range recorder.Stream(ctx, llm.CompletionRequest{System: []string{"system"}}) {
		if err != nil {
			t.Fatalf("Stream: %v", err)
		}
	}
	if got := journal.settled["call-1"]; got != 12 {
		t.Fatalf("stream settled at %d tokens; want 12", got)
	}

	_, _, usage, err := recorder.Complete(ctx, llm.CompletionRequest{System: []string{"system"}})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if usage.InputTokens != 11 || usage.OutputTokens != 3 {
		t.Fatalf("Complete usage=%+v", usage)
	}
	if got := journal.settled["call-2"]; got != 14 {
		t.Fatalf("complete settled at %d tokens; want 14", got)
	}
}
