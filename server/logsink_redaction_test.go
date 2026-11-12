package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogSinkRedactsBeforeFanout(t *testing.T) {
	var terminal bytes.Buffer
	sink := &logSinkT{cap: 10, subs: map[chan LogLine]struct{}{}, out: &terminal,
		dbCh: make(chan dbWriteReq, 1)}
	ch, unsubscribe := sink.subscribe()
	defer unsubscribe()
	input := []byte("[llm] request failed: Authorization: Bearer synthetic-secret-123\n")
	want := "[llm] request failed: Authorization: Bearer ***"
	if n, err := sink.Write(input); err != nil || n != len(input) {
		t.Fatalf("Write = %d, %v", n, err)
	}
	lines, _ := sink.recent(0, 10)
	if terminal.String() != want+"\n" || len(lines) != 1 || lines[0].Text != want {
		t.Fatal("terminal or ring received unredacted data")
	}
	if lines[0].Tag != "llm" || lines[0].Level != "error" {
		t.Fatal("diagnostic metadata was lost")
	}
	select {
	case got := <-ch:
		if got.Text != want {
			t.Fatal("subscriber received unredacted data")
		}
	default:
		t.Fatal("missing subscriber delivery")
	}
	select {
	case got := <-sink.dbCh:
		if got.ll.Text != want {
			t.Fatal("database queue received unredacted data")
		}
	default:
		t.Fatal("missing database delivery")
	}
}

func TestJSONErrorRedactsCredentials(t *testing.T) {
	w := httptest.NewRecorder()
	writeErr(w, http.StatusBadGateway, "provider failed: https://user:synthetic-password@example.test/v1?api_key=synthetic-key")
	var got map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if w.Code != http.StatusBadGateway || strings.Contains(got["error"], "synthetic-") ||
		!strings.Contains(got["error"], "provider failed: https://***@example.test/v1?api_key=***") {
		t.Fatal("JSON error did not preserve safe context")
	}
}
