package toolrunner

import (
	"context"
	"encoding/json"
	"net/http"
	"net"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestModeBProbeUsesOnlyExplicitProxy(t *testing.T) {
	seen := make(chan string, 1)
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen <- r.Method + " " + r.URL.Host + r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("1.1.1.1\n"))
	}))
	defer proxy.Close()
	proxyURL, err := url.Parse(proxy.URL)
	if err != nil {
		t.Fatal(err)
	}
	got, err := probeModeBEgress(context.Background(), proxyURL, "http://probe.invalid/ip")
	if err != nil || got != "1.1.1.1" {
		t.Fatalf("probe got=%q err=%v", got, err)
	}
	select {
	case request := <-seen:
		if request != "GET probe.invalid/ip" {
			t.Fatalf("probe did not use explicit HTTP proxy: %q", request)
		}
	default:
		t.Fatal("explicit proxy received no request")
	}
}

func TestModeBProbeNeverFallsBackToDirectWhenProxyIsDown(t *testing.T) {
	var directRequests int
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		directRequests++
		_, _ = w.Write([]byte("1.1.1.1"))
	}))
	defer target.Close()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	proxyAddress := listener.Addr().String()
	_ = listener.Close()
	proxyURL, err := url.Parse("http://" + proxyAddress)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := probeModeBEgress(context.Background(), proxyURL, target.URL); err == nil {
		t.Fatal("probe unexpectedly succeeded with its explicit proxy unavailable")
	}
	if directRequests != 0 {
		t.Fatalf("probe fell back to direct traffic: %d requests", directRequests)
	}
}

func TestModeBProxyURLRejectsAnyNonPrivateGateway(t *testing.T) {
	for _, raw := range []string{
		"", "http://127.0.0.1:7890", "http://host.docker.internal:7890",
		"http://mode-b-egress:7891", "socks5://mode-b-egress:7890",
		"http://user:pass@mode-b-egress:7890", "http://mode-b-egress:7890/path",
	} {
		if _, err := modeBProxyURL(raw); err == nil {
			t.Fatalf("accepted unapproved proxy URL %q", raw)
		}
	}
	if _, err := modeBProxyURL("http://mode-b-egress:7890"); err != nil {
		t.Fatalf("rejected configured private gateway: %v", err)
	}
}

func TestModeBEgressProbeHandlerRequiresModeBAndPublicExpectedIP(t *testing.T) {
	s := NewService()
	defer s.Close()
	handler := s.Handler()
	request := func(expected string) *httptest.ResponseRecorder {
		t.Helper()
		body, err := json.Marshal(EgressProbeRequest{ExpectedExitIP: expected})
		if err != nil {
			t.Fatal(err)
		}
		r := httptest.NewRequest(http.MethodPost, "/egress", strings.NewReader(string(body)))
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w
	}
	for _, invalid := range []string{"", "127.0.0.1", "10.0.0.1", "::1", "not-an-ip"} {
		if got := request(invalid); got.Code != http.StatusBadRequest {
			t.Fatalf("expected invalid exit %q to be rejected, got %d", invalid, got.Code)
		}
	}
	t.Setenv("ARTEX_RUNNER_PROXY_ADDR", "")
	if got := request("1.1.1.1"); got.Code != http.StatusServiceUnavailable {
		t.Fatalf("missing Mode B gateway should be unavailable, got %d", got.Code)
	}
	t.Setenv("ARTEX_RUNNER_PROXY_ADDR", "http://host.docker.internal:7890")
	if got := request("1.1.1.1"); got.Code != http.StatusServiceUnavailable {
		t.Fatalf("unapproved gateway should be unavailable, got %d", got.Code)
	}
}
