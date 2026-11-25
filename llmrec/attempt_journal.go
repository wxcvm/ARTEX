package llmrec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/Autumn-27/norma/llm"
)

var ErrAttemptJournal = errors.New("durable provider-attempt journal failed")

// AttemptJournal reserves and settles one row per HTTP request. Unknown outcomes
// retain their reservation; the database owner lease promotes them to
// recovery_required on graceful release or process takeover.
type AttemptJournal interface {
	Reserve(context.Context, int64) (string, error)
	Settle(context.Context, string, int64) error
	MarkUnknown(context.Context, string) error
}

type attemptJournalContextKey struct{}
type attemptScopeContextKey struct{}

// WithAttemptJournal enables durable accounting for every HTTP request made by
// the provider under ctx. Non-task callers without this hook keep existing behavior.
func WithAttemptJournal(ctx context.Context, journal AttemptJournal) context.Context {
	if journal == nil {
		return ctx
	}
	return context.WithValue(ctx, attemptJournalContextKey{}, journal)
}

// AttemptScope groups the HTTP retries performed inside one provider invocation.
// The transport reserves each RoundTrip; Recorder finishes the scope with the
// measured usage only after the provider has completed consuming its response.
type AttemptScope struct {
	journal AttemptJournal
	ctx     context.Context
	estimate int64

	mu       sync.Mutex
	pending  []*TransportAttempt
	finished bool
}

type TransportAttempt struct {
	scope          *AttemptScope
	callID         string
	reservedTokens int64
	once           sync.Once
}

// NewAttemptScope creates a fresh retry scope for one Recorder -> Provider call.
// Each empty-response retry gets its own scope and therefore its own reservations.
func NewAttemptScope(ctx context.Context, req llm.CompletionRequest) (context.Context, *AttemptScope) {
	if ctx == nil {
		return ctx, nil
	}
	journal, _ := ctx.Value(attemptJournalContextKey{}).(AttemptJournal)
	if journal == nil {
		return ctx, nil
	}
	scope := &AttemptScope{
		journal: journal,
		ctx: ctx,
		estimate: estimateCompletionTokens(req),
	}
	return context.WithValue(ctx, attemptScopeContextKey{}, scope), scope
}

// BeginTransportAttempt reserves budget before the HTTP request is sent. A nil
// result means no durable journal is attached to this request context.
func BeginTransportAttempt(ctx context.Context) (*TransportAttempt, error) {
	if ctx == nil {
		return nil, nil
	}
	scope, _ := ctx.Value(attemptScopeContextKey{}).(*AttemptScope)
	if scope == nil {
		return nil, nil
	}
	scope.mu.Lock()
	finished := scope.finished
	scope.mu.Unlock()
	if finished {
		return nil, fmt.Errorf("%w: provider request started after its attempt scope finished", ErrAttemptJournal)
	}
	callID, err := scope.journal.Reserve(ctx, scope.estimate)
	if err != nil {
		return nil, fmt.Errorf("%w: reserve provider request: %v", ErrAttemptJournal, err)
	}
	return &TransportAttempt{scope: scope, callID: callID, reservedTokens: scope.estimate}, nil
}

// Complete records the HTTP response boundary. Rejections known not to have
// produced a model response settle at zero tokens; transport failures and 408/5xx
// remain unknown. Successful responses stay reserved until Recorder sees whether
// the provider consumed the response completely.
func (a *TransportAttempt) Complete(status int, roundTripErr error) error {
	if a == nil || a.scope == nil {
		return nil
	}
	var result error
	a.once.Do(func() {
		switch {
		case roundTripErr != nil || status == 0 || status == 408 || status >= 500:
			if err := a.scope.journal.MarkUnknown(a.scope.ctx, a.callID); err != nil {
				result = fmt.Errorf("%w: retain ambiguous HTTP attempt: %v", ErrAttemptJournal, err)
			}
		case status >= 200 && status < 300:
			a.scope.mu.Lock()
			a.scope.pending = append(a.scope.pending, a)
			a.scope.mu.Unlock()
		default:
			if err := a.scope.journal.Settle(a.scope.ctx, a.callID, 0); err != nil {
				result = fmt.Errorf("%w: settle rejected HTTP attempt: %v", ErrAttemptJournal, err)
			}
		}
	})
	return result
}

// Finish settles the final completed HTTP response with measured input+output
// usage. Earlier 2xx responses imply an internal provider retry and stay unknown;
// an interrupted response also stays unknown. When a provider omits usage, the
// reserved request estimate is charged rather than silently treating it as free.
func (s *AttemptScope) Finish(usage llm.Usage, completed bool) error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	if s.finished {
		s.mu.Unlock()
		return nil
	}
	s.finished = true
	pending := append([]*TransportAttempt(nil), s.pending...)
	s.mu.Unlock()

	if len(pending) == 0 {
		return nil
	}
	var errs []error
	if !completed {
		for _, attempt := range pending {
			if err := s.journal.MarkUnknown(s.ctx, attempt.callID); err != nil {
				errs = append(errs, fmt.Errorf("%w: retain interrupted provider response: %v", ErrAttemptJournal, err))
			}
		}
		return errors.Join(errs...)
	}
	for i, attempt := range pending {
		if i < len(pending)-1 {
			if err := s.journal.MarkUnknown(s.ctx, attempt.callID); err != nil {
				errs = append(errs, fmt.Errorf("%w: retain superseded HTTP response: %v", ErrAttemptJournal, err))
			}
			continue
		}
		actualTokens := max(int64(usage.InputTokens), 0) + max(int64(usage.OutputTokens), 0)
		if actualTokens == 0 {
			actualTokens = attempt.reservedTokens
		}
		if err := s.journal.Settle(s.ctx, attempt.callID, actualTokens); err != nil {
			errs = append(errs, fmt.Errorf("%w: settle completed provider response: %v", ErrAttemptJournal, err))
		}
	}
	return errors.Join(errs...)
}

func estimateCompletionTokens(req llm.CompletionRequest) int64 {
	raw, err := json.Marshal(req)
	var estimate int64 = 1
	if err == nil {
		estimate = int64((len(raw) + 3) / 4)
	}
	if req.MaxTokens > 0 {
		estimate += int64(req.MaxTokens)
	}
	return max(estimate, 1)
}
