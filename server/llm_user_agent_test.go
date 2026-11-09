package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

// TestLLMTestValidatesUserAgent confirms the test endpoint uses the same
// header semantics as real model requests and blocks invalid header bytes.
func TestLLMTestValidatesUserAgent(t *testing.T) {
	var hits atomic.Int32
	var gotUA atomic.Value
	fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		gotUA.Store(r.Header.Get("User-Agent"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"m","type":"message","role":"assistant","model":"m","content":[{"type":"text","text":"OK"}],"stop_reason":"end_turn","usage":{"input_tokens":1,"output_tokens":1}}`))
	}))
	t.Cleanup(fake.Close)
	s := &Server{}
	post := func(ua string) *httptest.ResponseRecorder {
		body, _ := json.Marshal(map[string]any{"provider":"anthropic","model":"m","base_url":fake.URL,"api_key":"k","streaming":false,"user_agent":ua})
		req := httptest.NewRequest(http.MethodPost, "/api/llm/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		s.testLLM(rec, req)
		return rec
	}
	if rec := post("good/1.0\r\nX-Injected: 1"); rec.Code != 400 || hits.Load() != 0 {
		t.Fatalf("invalid UA: status %d hits %d body %s", rec.Code, hits.Load(), rec.Body.String())
	}
	const want = "ARTEX-Gateway/1.0"
	if rec := post("  "+want+"  "); rec.Code != 200 {
		t.Fatalf("valid UA: status %d body %s", rec.Code, rec.Body.String())
	}
	if got, _ := gotUA.Load().(string); hits.Load() != 1 || got != want {
		t.Fatalf("header = %q hits=%d want %q", got, hits.Load(), want)
	}
}
