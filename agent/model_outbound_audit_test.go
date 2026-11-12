package agent

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestModelOutboundAuditRoundTrip(t *testing.T) {
	audit := newModelAudit(2)
	base := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("Authorization") != "Bearer test-secret" {
			t.Fatal("request header was changed")
		}
		return &http.Response{StatusCode: 429, Body: io.NopCloser(strings.NewReader("test-response-secret"))}, nil
	})
	transport := modelAuditTransport{base: base, audit: audit, proxyConfigured: true}
	req, err := http.NewRequest(http.MethodPost, "https://user:test-password@example.test:8443/private-path?api_key=test-query-secret", strings.NewReader("test-body-secret"))
	if err != nil { t.Fatal(err) }
	req.Header.Set("Authorization", "Bearer test-secret")
	resp, err := transport.RoundTrip(req)
	if err != nil { t.Fatal(err) }
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil || string(body) != "test-response-secret" || resp.StatusCode != 429 {
		t.Fatal("audit changed the response")
	}
	s := audit.snapshot()
	if s.Total != 1 || len(s.Events) != 1 { t.Fatal("missing audit event") }
	e := s.Events[0]
	if e.Host != "example.test" || e.Port != "8443" || e.Scheme != "https" || !e.ProxyConfigured ||
		e.Status != 429 || e.Result != "http_response" || e.Sequence != 1 {
		t.Fatalf("unexpected audit metadata: %+v", e)
	}
	encoded, err := json.Marshal(s)
	if err != nil { t.Fatal(err) }
	for _, secret := range []string{"test-secret", "test-password", "private-path", "test-query-secret", "test-body-secret", "test-response-secret"} {
		if strings.Contains(string(encoded), secret) { t.Fatal("audit retained sensitive request or response data") }
	}
}

func TestModelOutboundAuditErrorsAndCapacity(t *testing.T) {
	audit := newModelAudit(2)
	req, err := http.NewRequest(http.MethodGet, "http://[::1]/v1", nil)
	if err != nil { t.Fatal(err) }
	for _, tc := range []struct{ err error; result string }{
		{errors.New("private transport error"), "transport_error"},
		{context.Canceled, "canceled"},
		{context.DeadlineExceeded, "timeout"},
	} {
		tr := modelAuditTransport{base: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			return nil, tc.err
		}), audit: audit}
		if _, err := tr.RoundTrip(req); err != tc.err { t.Fatal("transport error changed") }
		s := audit.snapshot()
		e := s.Events[0]
		if e.Result != tc.result || e.Status != 0 || e.Host != "::1" || e.Port != "80" {
			t.Fatalf("unexpected error metadata: %+v", e)
		}
	}
	s := audit.snapshot()
	if s.Total != 3 || s.Capacity != 2 || len(s.Events) != 2 || s.Events[0].Sequence != 3 || s.Events[1].Sequence != 2 {
		t.Fatal("incorrect ring retention or ordering")
	}
	s.Events[0].Host = "changed"
	if audit.snapshot().Events[0].Host != "::1" { t.Fatal("snapshot aliases audit storage") }
}

func TestModelOutboundAuditRealHTTP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()
	audit := newModelAudit(8)
	client := &http.Client{Transport: modelAuditTransport{base: http.DefaultTransport, audit: audit}}
	resp, err := client.Get(srv.URL + "/private?token=secret")
	if err != nil { t.Fatal(err) }
	resp.Body.Close()
	s := audit.snapshot()
	if s.Total != 1 || s.Events[0].Status != http.StatusNoContent || s.Events[0].Host != "127.0.0.1" && s.Events[0].Host != "::1" {
		t.Fatalf("unexpected real HTTP audit: %+v", s)
	}
}

func TestModelOutboundAuditConcurrent(t *testing.T) {
	audit := newModelAudit(16)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			audit.add(ModelOutboundEvent{Host: "example.test"})
			_ = audit.snapshot()
		}()
	}
	wg.Wait()
	s := audit.snapshot()
	if s.Total != 100 || len(s.Events) != 16 || s.Events[0].Sequence != 100 {
		t.Fatalf("unexpected concurrent audit counts: %+v", s)
	}
}
