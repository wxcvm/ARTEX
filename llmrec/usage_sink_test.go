package llmrec

import (
	"context"
	"testing"

	"github.com/Autumn-27/norma/llm"
)

func TestUsageSinkWithoutDatabase(t *testing.T) {
	var calls int
	var got llm.Usage
	ctx := WithUsageSink(context.Background(), func(u llm.Usage) { calls++; got = u })
	r := &Recorder{model: "test"}
	r.recordUsage(ctx, "1", 2, "worker", llm.Usage{InputTokens: 12, OutputTokens: 3}, 1, "ok")
	if calls != 1 || got.InputTokens != 12 || got.OutputTokens != 3 {
		t.Fatalf("usage sink: calls=%d usage=%+v", calls, got)
	}
}
