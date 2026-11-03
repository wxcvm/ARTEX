package llmpool

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Autumn-27/norma/llm"
)

// Cancelling a request or rejecting its payload must neither penalize nor heal
// the provider. In particular, intermittent cancellations must not erase the
// transient failures needed to trip the circuit breaker.
func TestCompleteNonFailoverErrorPreservesHealth(t *testing.T) {
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	expired, stop := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer stop()
	cases := []struct {
		name string
		ctx  context.Context
		err  error
	}{
		{"cancelled context", cancelled, errors.New("request interrupted")},
		{"expired context", expired, errors.New("request interrupted")},
		{"cancellation error", context.Background(), fmt.Errorf("provider: %w", context.Canceled)},
		{"deadline error", context.Background(), fmt.Errorf("provider: %w", context.DeadlineExceeded)},
		{"bad request", context.Background(), errors.New("openai: status 400: invalid request")},
		{"unprocessable request", context.Background(), errors.New("openai: status 422: invalid payload")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var forgotten atomic.Bool
			reg := NewRegistry(nil, func(int64) { forgotten.Store(true) })
			before := State{Fails: 2, Trips: 1, LastError: "status 503", LastAt: time.Now()}
			reg.Restore(1, before)
			a := &fakeProv{errs: []error{tc.err}}
			b := okProv("backup", "unused")
			p := New([]*Member{member(1, "primary", 10, a), member(2, "backup", 5, b)}, reg)

			_, _, _, err := p.Complete(tc.ctx, llm.CompletionRequest{})
			if err != tc.err {
				t.Fatalf("err = %v, want the original error %v", err, tc.err)
			}
			if a.calls != 1 || b.calls != 0 {
				t.Fatalf("calls: primary=%d backup=%d, want 1/0", a.calls, b.calls)
			}
			if after := reg.Get(1); after != before {
				t.Fatalf("non-success changed health: before=%+v after=%+v", before, after)
			}
			flushRegistry(t, reg)
			if forgotten.Load() {
				t.Fatal("non-success deleted persisted circuit-breaker state")
			}
		})
	}
}

// The last-resort probe of an open single-member chain must not clear its
// cooling-off window merely because the caller cancelled the probe.
func TestCompleteCancelledProbePreservesOpenBreaker(t *testing.T) {
	reg := NewRegistry(nil, nil)
	before := State{Trips: 2, OpenUntil: time.Now().Add(time.Hour), LastError: "status 402", LastAt: time.Now()}
	reg.Restore(1, before)
	a := &fakeProv{errs: []error{context.Canceled}}
	p := New([]*Member{member(1, "primary", RankActive, a)}, reg)
	_, _, _, err := p.Complete(context.Background(), llm.CompletionRequest{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	if a.calls != 1 {
		t.Fatalf("probe calls=%d, want 1", a.calls)
	}
	if after := reg.Get(1); after != before {
		t.Fatalf("cancelled probe changed health: before=%+v after=%+v", before, after)
	}
}

func TestCompleteSuccessClearsHealth(t *testing.T) {
	var forgotten atomic.Bool
	reg := NewRegistry(nil, func(int64) { forgotten.Store(true) })
	reg.Restore(1, State{Fails: 2, Trips: 1, LastError: "status 503", LastAt: time.Now()})
	p := New([]*Member{member(1, "primary", RankActive, okProv("primary", "recovered"))}, reg)
	_, _, _, err := p.Complete(context.Background(), llm.CompletionRequest{})
	if err != nil {
		t.Fatalf("successful completion returned %v", err)
	}
	flushRegistry(t, reg)
	if after := reg.Get(1); after != (State{}) || !forgotten.Load() {
		t.Fatalf("success did not clear health: state=%+v forgotten=%v", after, forgotten.Load())
	}
}
