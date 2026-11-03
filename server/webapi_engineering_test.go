package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/Autumn-27/artex/llmpool"
	"github.com/Autumn-27/artex/mcphttp"
	"github.com/Autumn-27/norma/llm"
)

// This suite exercises existing transport, pagination, control normalization and
// model-pool contracts. It does NOT run a planner or measure vulnerability discovery.
const webAPIEngineeringVersion = "webapi-engineering-v1"

type engineeringProvider struct {
	text  string
	calls int
}

func (p *engineeringProvider) Complete(context.Context, llm.CompletionRequest) (llm.Message, string, llm.Usage, error) {
	p.calls++
	return llm.Message{Content: []llm.ContentBlock{llm.TextBlock(p.text)}}, "end_turn", llm.Usage{}, nil
}
func (p *engineeringProvider) Stream(context.Context, llm.CompletionRequest) iter.Seq2[llm.StreamEvent, error] {
	return func(yield func(llm.StreamEvent, error) bool) {
		p.calls++
		yield(llm.StreamEvent{Type: llm.SETextDelta, Text: p.text}, nil)
	}
}

type engineeringEvidence struct {
	ID     string          `json:"id"`
	SHA256 string          `json:"sha256"`
	File   string          `json:"file,omitempty"`
	Body   json.RawMessage `json:"-"`
}

func TestWebAPIEngineering(t *testing.T) {
	for _, scenario := range []string{"normal", "pagination", "authentication", "rate_limit", "timeout", "malformed_json", "large_utf8", "cancel", "duplicates", "incomplete"} {
		t.Run(scenario, func(t *testing.T) {
			started := time.Now()
			var requests, toolCalls atomic.Int64
			var mu sync.Mutex
			var evidence []engineeringEvidence
			entered, release := make(chan struct{}), make(chan struct{})
			large := strings.Repeat("证据🧪é", 32768)
			fixture := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					return
				}
				var req struct {
					ID     int    `json:"id"`
					Method string `json:"method"`
					Params struct {
						Arguments struct {
							Page int `json:"pageIndex"`
						} `json:"arguments"`
					} `json:"params"`
				}
				if err := json.Unmarshal(body, &req); err != nil {
					t.Error(err)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				var response string
				status := 200
				switch req.Method {
				case "initialize":
					response = fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"protocolVersion":"2025-06-18","capabilities":{},"serverInfo":{"name":"fixture","version":"1"}}}`, req.ID)
				case "notifications/initialized":
					status = 202
				case "tools/call":
					toolCalls.Add(1)
					text := `{"list":[{"id":"1","evidence_ref":"fixture-1"}]}`
					switch scenario {
					case "authentication":
						status = 401
						response = `{"error":"fixture authentication failed"}`
					case "rate_limit":
						status = 429
						w.Header().Set("Retry-After", "1")
						response = `{"error":"fixture rate limited"}`
					case "malformed_json":
						response = `{"jsonrpc":broken`
					case "timeout", "cancel":
						close(entered)
						select {
						case <-r.Context().Done():
						case <-release:
						}
						status = 499
						response = `{"error":"request stopped"}`
					case "large_utf8":
						text = large
					case "pagination":
						if req.Params.Arguments.Page == 1 {
							text = `{"list":[{"id":"1","evidence_ref":"fixture-1"},{"id":"2","evidence_ref":"fixture-2"}]}`
						} else {
							text = `{"list":[{"id":"3","evidence_ref":"fixture-3"}]}`
						}
					case "duplicates":
						text = `{"list":[{"id":"1"},{"id":"01"},{"id":"2"},{"id":"1"}]}`
					case "incomplete":
						text = `{"status":"running","findings":[],"evidence_refs":["fixture-1"]}`
					}
					if response == "" {
						b, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": map[string]any{"content": []any{map[string]string{"type": "text", "text": text}}}})
						response = string(b)
					}
				default:
					status = 400
					response = `{"error":"unexpected method"}`
				}
				record, _ := json.Marshal(map[string]any{"request": json.RawMessage(body), "status": status, "response": response})
				sum := sha256.Sum256(record)
				mu.Lock()
				evidence = append(evidence, engineeringEvidence{ID: fmt.Sprintf("rpc-%d-%s", req.ID, req.Method), SHA256: hex.EncodeToString(sum[:]), Body: record})
				mu.Unlock()
				w.WriteHeader(status)
				_, _ = io.WriteString(w, response)
			}))
			defer fixture.Close()
			defer close(release)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			client, err := mcphttp.New(ctx, "engineering-fixture", fixture.URL, nil, false, 1)
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()
			provider := &engineeringProvider{text: `{"status":"incomplete","evidence_refs":["fixture-1"]}`}
			assessment := "not_applicable"
			switch scenario {
			case "pagination":
				items, truncated, err := (&Server{}).ssPageAll(ctx, client, "subdomain", nil, 2)
				if err != nil || truncated || len(items) != 3 || toolCalls.Load() != 2 || !strings.Contains(string(items[2]), "fixture-3") {
					t.Fatalf("pagination: items=%s truncated=%v err=%v", items, truncated, err)
				}
			case "cancel":
				go func() { <-entered; cancel() }()
				_, err := client.Call(ctx, "fixture", nil)
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("cancel: %v", err)
				}
				assessment = "incomplete"
			default:
				text, err := client.Call(ctx, "fixture", nil)
				switch scenario {
				case "authentication", "rate_limit":
					code := "401"
					if scenario == "rate_limit" {
						code = "429"
					}
					if err == nil || !strings.Contains(err.Error(), code) || toolCalls.Load() != 1 {
						t.Fatalf("expected %s without implicit retry: %v", code, err)
					}
					assessment = "environment_error"
				case "timeout":
					if !errors.Is(err, context.DeadlineExceeded) {
						t.Fatalf("deadline: %v", err)
					}
					assessment = "incomplete"
				case "malformed_json":
					if err == nil {
						t.Fatal("malformed JSON accepted")
					}
					assessment = "environment_error"
				case "large_utf8":
					if err != nil || !utf8.ValidString(text) || text != large {
						t.Fatalf("UTF-8 output lost: bytes=%d err=%v", len(text), err)
					}
				case "duplicates":
					if err != nil {
						t.Fatal(err)
					}
					var data struct {
						List []struct {
							ID string `json:"id"`
						} `json:"list"`
					}
					if err := json.Unmarshal([]byte(text), &data); err != nil {
						t.Fatal(err)
					}
					var ids []string
					for _, item := range data.List {
						ids = append(ids, item.ID)
					}
					unique := normalizeBatchTaskIDs(ids)
					if len(unique) != 2 || unique[0].id != "1" || unique[1].id != "2" {
						t.Fatalf("duplicate control identities: %+v", unique)
					}
				case "normal", "incomplete":
					if err != nil || !strings.Contains(text, "fixture-1") {
						t.Fatalf("evidence reference lost: %s %v", text, err)
					}
					if scenario == "normal" {
						provider.text = `{"status":"completed_no_findings","evidence_refs":["fixture-1"]}`
					}
					pool := llmpool.New([]*llmpool.Member{{ID: 1, Name: "mock", Model: "fixture-v1", Prov: provider}}, nil)
					msg, _, _, err := pool.Complete(ctx, llm.CompletionRequest{})
					if err != nil || len(msg.Content) != 1 {
						t.Fatalf("mock model: %v", err)
					}
					var result struct {
						Status string   `json:"status"`
						Refs   []string `json:"evidence_refs"`
					}
					if err := json.Unmarshal([]byte(msg.Content[0].Text), &result); err != nil {
						t.Fatal(err)
					}
					want := "incomplete"
					if scenario == "normal" {
						want = "completed_no_findings"
					}
					if result.Status != want || len(result.Refs) != 1 || !strings.Contains(text, result.Refs[0]) || provider.calls != 1 {
						t.Fatalf("model result/evidence contract: %+v", result)
					}
					assessment = result.Status // explicitly synthetic, not a real assessment
				}
			}
			fixture.Close() // join handlers before reading evidence (including cancel)
			mu.Lock()
			defer mu.Unlock()
			if int64(len(evidence)) != requests.Load() || len(evidence) < 3 {
				t.Fatalf("missing request evidence: %d/%d", len(evidence), requests.Load())
			}
			evidenceDir := os.Getenv("ARTEX_ENGINEERING_EVIDENCE_DIR")
			if evidenceDir != "" {
				dir, err := os.MkdirTemp(evidenceDir, scenario+"-")
				if err != nil {
					t.Fatal(err)
				}
				for i := range evidence {
					path := filepath.Join(dir, fmt.Sprintf("%d.json", i))
					if err := os.WriteFile(path, evidence[i].Body, 0600); err != nil {
						t.Fatal(err)
					}
					evidence[i].File = path
				}
			}
			metric, _ := json.Marshal(map[string]any{"case": scenario, "case_version": webAPIEngineeringVersion, "elapsed_ms": time.Since(started).Milliseconds(), "http_calls": requests.Load(), "tool_calls": toolCalls.Load(), "model_calls": provider.calls, "usage_kind": "simulated", "real_model_calls": 0, "assessment_status": assessment, "evidence_complete": true, "evidence": evidence})
			t.Log("ARTEX_ENGINEERING_METRIC " + string(metric))
		})
	}
}
