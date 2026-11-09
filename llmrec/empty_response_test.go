package llmrec

import (
	"context"
	"errors"
	"iter"
	"testing"

	"github.com/Autumn-27/norma/llm"
)

type emptyScriptProvider struct {
	attempts int
	streams [][]llm.StreamEvent
	completeCalls int
	messages []llm.Message
}

func (p *emptyScriptProvider) Stream(_ context.Context, _ llm.CompletionRequest) iter.Seq2[llm.StreamEvent, error] {
	index := p.attempts
	p.attempts++
	return func(yield func(llm.StreamEvent, error) bool) {
		if index >= len(p.streams) { return }
		for _, event := range p.streams[index] {
			if !yield(event, nil) { return }
		}
	}
}

func (p *emptyScriptProvider) Complete(_ context.Context, _ llm.CompletionRequest) (llm.Message, string, llm.Usage, error) {
	index := p.completeCalls
	p.completeCalls++
	if index >= len(p.messages) { return llm.Message{}, "stop", llm.Usage{InputTokens: 2}, nil }
	return p.messages[index], "stop", llm.Usage{InputTokens: 2}, nil
}

func emptyStream() []llm.StreamEvent {
	return []llm.StreamEvent{
		{Type: llm.SEMessageStart, Usage: llm.Usage{InputTokens: 2}},
		{Type: llm.SEMessageDelta, Usage: llm.Usage{OutputTokens: 1}},
		{Type: llm.SEMessageStop},
	}
}

func TestStreamRetriesOnlyEmptyCompletions(t *testing.T) {
	content := append(emptyStream()[:1], llm.StreamEvent{Type: llm.SETextDelta, Text: "answer"})
	provider := &emptyScriptProvider{streams: [][]llm.StreamEvent{emptyStream(), emptyStream(), content}}
	r := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	var events []llm.StreamEvent
	for ev, err := range r.Stream(context.Background(), llm.CompletionRequest{}) {
		if err != nil { t.Fatal(err) }
		events = append(events, ev)
	}
	if provider.attempts != 3 { t.Fatalf("attempts=%d", provider.attempts) }
	if len(events) != 2 || events[0].Type != llm.SEMessageStart ||
		events[1].Type != llm.SETextDelta || events[1].Text != "answer" {
		t.Fatalf("unexpected events: %+v", events)
	}
}

func TestStreamEmptyExhaustedReturnsError(t *testing.T) {
	provider := &emptyScriptProvider{streams: [][]llm.StreamEvent{emptyStream(), emptyStream(), emptyStream()}}
	r := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	count := 0
	for _, err := range r.Stream(context.Background(), llm.CompletionRequest{}) {
		count++
		if !errors.Is(err, errEmptyResponse) { t.Fatalf("event %d, err=%v", count, err) }
	}
	if count != 1 || provider.attempts != 3 { t.Fatalf("events=%d attempts=%d", count, provider.attempts) }
}

func TestStreamToolCallDoesNotRetry(t *testing.T) {
	one := append(emptyStream()[:1], llm.StreamEvent{Type: llm.SEToolUseStart})
	provider := &emptyScriptProvider{streams: [][]llm.StreamEvent{one}}
	r := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	for _, err := range r.Stream(context.Background(), llm.CompletionRequest{}) {
		if err != nil { t.Fatal(err) }
	}
	if provider.attempts != 1 { t.Fatalf("tool call retried: %d", provider.attempts) }
}

func TestStreamConsumerStopsWithoutRetry(t *testing.T) {
	one := append(emptyStream()[:1], llm.StreamEvent{Type: llm.SETextDelta, Text: "partial"})
	provider := &emptyScriptProvider{streams: [][]llm.StreamEvent{one}}
	r := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	for ev := range r.Stream(context.Background(), llm.CompletionRequest{}) {
		if ev.Type == llm.SETextDelta { break }
	}
	if provider.attempts != 1 { t.Fatalf("consumer stopped but retried: %d", provider.attempts) }
}

func TestCompleteRetriesEmptyAndTotalsUsage(t *testing.T) {
	provider := &emptyScriptProvider{messages: []llm.Message{{}, {}, {Content: []llm.ContentBlock{llm.TextBlock("answer")}}}}
	r := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	msg, _, usage, err := r.Complete(context.Background(), llm.CompletionRequest{})
	if err != nil || msg.Text() != "answer" || usage.InputTokens != 6 || provider.completeCalls != 3 {
		t.Fatalf("msg=%q usage=%+v attempts=%d err=%v", msg.Text(), usage, provider.completeCalls, err)
	}
}

func TestCompleteEmptyExhaustedReturnsError(t *testing.T) {
	provider := &emptyScriptProvider{}
	r := Wrap(provider, nil, "", "", "", "", func() bool { return false })
	_, _, usage, err := r.Complete(context.Background(), llm.CompletionRequest{})
	if !errors.Is(err, errEmptyResponse) || usage.InputTokens != 6 || provider.completeCalls != 3 {
		t.Fatalf("usage=%+v attempts=%d err=%v", usage, provider.completeCalls, err)
	}
}

func TestEmptyMessageRecognizesThinkingAndToolCalls(t *testing.T) {
	for _, msg := range []llm.Message{
		{Content: []llm.ContentBlock{llm.TextBlock("ok")}},
		{Content: []llm.ContentBlock{{Type: llm.BlockThinking, Thinking: "yes"}}},
		{Content: []llm.ContentBlock{{Type: llm.BlockToolUse}}},
	} {
		if emptyMessage(msg) { t.Fatalf("nonempty message marked empty: %+v", msg) }
	}
	if !emptyMessage(llm.Message{}) { t.Fatal("empty message not detected") }
}
