package llmrec

import (
	"context"
	"errors"
	"testing"

	"github.com/Autumn-27/norma/llm"
)

func retryMeterContext(t *testing.T) (context.Context, *[]llm.Usage) {
	t.Helper()
	var recorded []llm.Usage
	ctx := WithUsageSink(context.Background(), func(u llm.Usage) {
		recorded = append(recorded, u)
	})
	return ctx, &recorded
}

func TestStreamingEmptyRetriesMeterEveryProviderAttempt(t *testing.T) {
	content := append(emptyStream()[:1], llm.StreamEvent{Type: llm.SETextDelta, Text: "answer"})
	provider := &emptyScriptProvider{streams: [][]llm.StreamEvent{emptyStream(), emptyStream(), content}}
	rec := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	ctx, recorded := retryMeterContext(t)
	for _, err := range rec.Stream(ctx, llm.CompletionRequest{}) {
		if err != nil {
			t.Fatal(err)
		}
	}
	if provider.attempts != 3 || len(*recorded) != 3 {
		t.Fatalf("stream attempts=%d budget calls=%d, want 3 each", provider.attempts, len(*recorded))
	}
	var tokens int
	for _, usage := range *recorded {
		tokens += usage.InputTokens + usage.OutputTokens
	}
	if tokens != 2+1+2+1+2 {
		t.Fatalf("retries lost or duplicated streamed usage: tokens=%d recorded=%+v", tokens, *recorded)
	}
}

func TestStreamingEmptyExhaustionMetersEveryProviderAttempt(t *testing.T) {
	provider := &emptyScriptProvider{streams: [][]llm.StreamEvent{emptyStream(), emptyStream(), emptyStream()}}
	rec := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	ctx, recorded := retryMeterContext(t)
	foundError := false
	for _, err := range rec.Stream(ctx, llm.CompletionRequest{}) {
		if errors.Is(err, errEmptyResponse) {
			foundError = true
		} else if err != nil {
			t.Fatal(err)
		}
	}
	if !foundError || provider.attempts != 3 || len(*recorded) != 3 {
		t.Fatalf("exhausted stream: error=%v attempts=%d budget calls=%d", foundError, provider.attempts, len(*recorded))
	}
}

func TestCompleteEmptyRetriesMeterEveryProviderAttempt(t *testing.T) {
	provider := &emptyScriptProvider{
		messages: []llm.Message{{}, {}, {Content: []llm.ContentBlock{llm.TextBlock("answer")}}},
	}
	rec := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	ctx, recorded := retryMeterContext(t)
	msg, _, usage, err := rec.Complete(ctx, llm.CompletionRequest{})
	if err != nil || msg.Text() != "answer" || usage.InputTokens != 6 {
		t.Fatalf("completed retry result: msg=%q usage=%+v err=%v", msg.Text(), usage, err)
	}
	if provider.completeCalls != 3 || len(*recorded) != 3 {
		t.Fatalf("complete attempts=%d budget calls=%d, want 3", provider.completeCalls, len(*recorded))
	}
	sum := 0
	for _, attempt := range *recorded {
		sum += attempt.InputTokens + attempt.OutputTokens
	}
	if sum != 6 {
		t.Fatalf("per-attempt tokens=%d, expected aggregate 6: %+v", sum, *recorded)
	}
}

func TestSingleProviderAttemptMetersExactlyOnce(t *testing.T) {
	provider := &emptyScriptProvider{messages: []llm.Message{{Content: []llm.ContentBlock{llm.TextBlock("answer")}}}}
	rec := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	ctx, recorded := retryMeterContext(t)
	_, _, _, err := rec.Complete(ctx, llm.CompletionRequest{})
	if err != nil || provider.completeCalls != 1 || len(*recorded) != 1 {
		t.Fatalf("single completion: attempts=%d meter calls=%d err=%v", provider.completeCalls, len(*recorded), err)
	}
}
