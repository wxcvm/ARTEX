package llmrec

import (
	"context"
	"errors"
	"testing"

	"github.com/Autumn-27/norma/llm"
)

func TestProviderPreflightBlocksFirstCompleteWithoutChargingUsage(t *testing.T) {
	denied := errors.New("budget admission denied")
	provider := &emptyScriptProvider{}
	rec := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	metered := 0
	ctx := WithUsageSink(context.Background(), func(llm.Usage) { metered++ })
	preflightCalls := 0
	ctx = WithProviderAdmission(ctx, func(context.Context) error {
		preflightCalls++
		return denied
	})
	_, _, usage, err := rec.Complete(ctx, llm.CompletionRequest{})
	if !errors.Is(err, denied) || provider.completeCalls != 0 || metered != 0 ||
		preflightCalls != 1 || usage.InputTokens != 0 {
		t.Fatalf("denied complete reached provider or consumed budget: calls=%d meter=%d preflight=%d usage=%+v err=%v",
			provider.completeCalls, metered, preflightCalls, usage, err)
	}
}

func TestProviderPreflightBlocksCompleteRetryAndPreservesFirstUsage(t *testing.T) {
	denied := errors.New("no remaining turns")
	provider := &emptyScriptProvider{messages: []llm.Message{{}, {Content: []llm.ContentBlock{llm.TextBlock("must not send")}}}}
	rec := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	var usages []llm.Usage
	ctx := WithUsageSink(context.Background(), func(u llm.Usage) { usages = append(usages, u) })
	preflightCalls := 0
	ctx = WithProviderAdmission(ctx, func(context.Context) error {
		preflightCalls++
		if preflightCalls == 2 {
			return denied
		}
		return nil
	})
	_, _, usage, err := rec.Complete(ctx, llm.CompletionRequest{})
	if !errors.Is(err, denied) || provider.completeCalls != 1 || preflightCalls != 2 ||
		len(usages) != 1 || usage.InputTokens != 2 {
		t.Fatalf("complete retry bypassed preflight or lost usage: attempts=%d preflight=%d metered=%+v total=%+v err=%v",
			provider.completeCalls, preflightCalls, usages, usage, err)
	}
}

func TestProviderPreflightBlocksFirstStreamWithoutChargingUsage(t *testing.T) {
	denied := errors.New("database unavailable")
	provider := &emptyScriptProvider{streams: [][]llm.StreamEvent{emptyStream()}}
	rec := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	metered := 0
	ctx := WithUsageSink(context.Background(), func(llm.Usage) { metered++ })
	ctx = WithProviderAdmission(ctx, func(context.Context) error { return denied })
	n := 0
	for _, err := range rec.Stream(ctx, llm.CompletionRequest{}) {
		n++
		if !errors.Is(err, denied) {
			t.Fatalf("unexpected preflight stream error: %v", err)
		}
	}
	if n != 1 || provider.attempts != 0 || metered != 0 {
		t.Fatalf("denied streaming request reached provider: events=%d attempts=%d meter=%d", n, provider.attempts, metered)
	}
}

func TestProviderPreflightBlocksStreamRetryAndPreservesFirstUsage(t *testing.T) {
	denied := errors.New("owner lease lost")
	provider := &emptyScriptProvider{streams: [][]llm.StreamEvent{emptyStream(), emptyStream()}}
	rec := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	metered := 0
	ctx := WithUsageSink(context.Background(), func(llm.Usage) { metered++ })
	preflightCalls := 0
	ctx = WithProviderAdmission(ctx, func(context.Context) error {
		preflightCalls++
		if preflightCalls == 2 {
			return denied
		}
		return nil
	})
	events := 0
	for _, err := range rec.Stream(ctx, llm.CompletionRequest{}) {
		events++
		if !errors.Is(err, denied) {
			t.Fatalf("unexpected stream retry error: %v", err)
		}
	}
	if events != 1 || provider.attempts != 1 || metered != 1 || preflightCalls != 2 {
		t.Fatalf("stream retry escaped preflight: events=%d attempts=%d metered=%d preflight=%d",
			events, provider.attempts, metered, preflightCalls)
	}
}

func TestProviderPreflightApprovesAllThreeRetries(t *testing.T) {
	content := append(emptyStream()[:1], llm.StreamEvent{Type: llm.SETextDelta, Text: "answer"})
	provider := &emptyScriptProvider{streams: [][]llm.StreamEvent{emptyStream(), emptyStream(), content}}
	rec := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	n := 0
	ctx := WithProviderAdmission(context.Background(), func(context.Context) error {
		n++
		return nil
	})
	for _, err := range rec.Stream(ctx, llm.CompletionRequest{}) {
		if err != nil {
			t.Fatal(err)
		}
	}
	if provider.attempts != 3 || n != 3 {
		t.Fatalf("admission was not re-evaluated per provider request: attempts=%d admission=%d", provider.attempts, n)
	}
}

func TestProviderPreflightRespectsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	provider := &emptyScriptProvider{}
	rec := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	_, _, _, err := rec.Complete(ctx, llm.CompletionRequest{})
	if !errors.Is(err, context.Canceled) || provider.completeCalls != 0 {
		t.Fatalf("cancelled context reached provider: calls=%d err=%v", provider.completeCalls, err)
	}
}
