package agent

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/Autumn-27/artex/llmrec"
	"github.com/Autumn-27/norma/llm"
)

type transportAttemptJournalFixture struct {
	mu         sync.Mutex
	next       int
	reserved   map[string]int64
	settled    map[string]int64
	unknown    map[string]bool
	reserveErr error
}

func (j *transportAttemptJournalFixture) Reserve(_ context.Context, estimate int64) (string, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.reserveErr != nil {
		return "", j.reserveErr
	}
	j.next++
	id := fmt.Sprintf("call-%d", j.next)
	if j.reserved == nil {
		j.reserved = map[string]int64{}
	}
	j.reserved[id] = estimate
	return id, nil
}

func (j *transportAttemptJournalFixture) Settle(_ context.Context, id string, tokens int64) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.settled == nil {
		j.settled = map[string]int64{}
	}
	j.settled[id] = tokens
	return nil
}

func (j *transportAttemptJournalFixture) MarkUnknown(_ context.Context, id string) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.unknown == nil {
		j.unknown = map[string]bool{}
	}
	j.unknown[id] = true
	return nil
}

type statusSequenceTransport struct {
	statuses []int
	calls    int
}

func (t *statusSequenceTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.calls >= len(t.statuses) {
		return nil, errors.New("unexpected extra HTTP attempt")
	}
	status := t.statuses[t.calls]
	t.calls++
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("fixture response")),
		Request:    req,
	}, nil
}

func TestQuotaAwareTransportJournalsEveryHTTPRetry(t *testing.T) {
	journal := &transportAttemptJournalFixture{}
	ctx, scope := llmrec.NewAttemptScope(
		llmrec.WithAttemptJournal(context.Background(), journal),
		llm.CompletionRequest{System: []string{"system prompt"}, MaxTokens: 128},
	)
	base := &statusSequenceTransport{statuses: []int{http.StatusBadGateway, http.StatusOK}}
	transport := quotaAwareTransport{base: base}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://provider.invalid/v1/messages", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []int{http.StatusBadGateway, http.StatusOK} {
		resp, err := transport.RoundTrip(req)
		if err != nil {
			t.Fatalf("RoundTrip: %v", err)
		}
		if resp.StatusCode != want {
			t.Fatalf("status=%d want=%d", resp.StatusCode, want)
		}
		_ = resp.Body.Close()
	}
	if err := scope.Finish(llm.Usage{InputTokens: 9, OutputTokens: 4}, true); err != nil {
		t.Fatal(err)
	}
	if base.calls != 2 || len(journal.reserved) != 2 {
		t.Fatalf("HTTP attempts were not individually reserved: calls=%d reservations=%v", base.calls, journal.reserved)
	}
	if !journal.unknown["call-1"] {
		t.Fatalf("superseded 502 attempt was not retained as unknown: %+v", journal)
	}
	if got := journal.settled["call-2"]; got != 13 {
		t.Fatalf("final response settled at %d tokens; want 13", got)
	}
}

func TestQuotaAwareTransportSettlesKnownRateLimitRejection(t *testing.T) {
	journal := &transportAttemptJournalFixture{}
	ctx, scope := llmrec.NewAttemptScope(
		llmrec.WithAttemptJournal(context.Background(), journal),
		llm.CompletionRequest{System: []string{"system prompt"}},
	)
	base := &statusSequenceTransport{statuses: []int{http.StatusTooManyRequests}}
	transport := quotaAwareTransport{base: base}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://provider.invalid/v1/messages", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if err := scope.Finish(llm.Usage{}, false); err != nil {
		t.Fatal(err)
	}
	if got := journal.settled["call-1"]; got != 0 {
		t.Fatalf("429 attempt settled at %d tokens; want zero", got)
	}
	if journal.unknown["call-1"] {
		t.Fatal("known 429 rejection was marked unknown")
	}
}

func TestQuotaAwareTransportDoesNotSendWhenReservationFails(t *testing.T) {
	journal := &transportAttemptJournalFixture{reserveErr: errors.New("database unavailable")}
	ctx, _ := llmrec.NewAttemptScope(
		llmrec.WithAttemptJournal(context.Background(), journal),
		llm.CompletionRequest{System: []string{"system prompt"}},
	)
	base := &statusSequenceTransport{statuses: []int{http.StatusOK}}
	transport := quotaAwareTransport{base: base}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://provider.invalid/v1/messages", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnprocessableEntity || base.calls != 0 {
		t.Fatalf("reservation failure did not fail closed: status=%d network_calls=%d", resp.StatusCode, base.calls)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "durable provider-attempt journal") {
		t.Fatalf("missing journal failure context: %s", body)
	}
}
