package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestProbeProxyEgressFallsBackAndValidatesIP(t *testing.T) {
	echo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/bad":
			_, _ = w.Write([]byte("not an IP address"))
		case "/fail":
			w.WriteHeader(http.StatusBadGateway)
		default:
			_, _ = w.Write([]byte("2001:db8::1\n"))
		}
	}))
	defer echo.Close()
	ip, _, err := probeProxyEgress(context.Background(), echo.Client(), []string{echo.URL + "/fail", echo.URL + "/bad", echo.URL + "/good"})
	if err != nil || ip != "2001:db8::1" {
		t.Fatalf("fallback returned ip=%q err=%v", ip, err)
	}
}

func TestProbeProxyEgressActuallyUsesConfiguredProxy(t *testing.T) {
	seen := make(chan string, 1)
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen <- r.URL.Host
		_, _ = w.Write([]byte("198.51.100.9"))
	}))
	defer proxy.Close()
	u, err := url.Parse(proxy.URL)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(u)}}
	ip, _, err := probeProxyEgress(context.Background(), client, []string{"http://unresolvable.invalid/echo"})
	if err != nil || ip != "198.51.100.9" {
		t.Fatalf("configured proxy failed: ip=%q err=%v", ip, err)
	}
	select {
	case host := <-seen:
		if host != "unresolvable.invalid" {
			t.Fatalf("request was not proxied: host=%q", host)
		}
	default:
		t.Fatal("proxy received no request")
	}
}

func TestGlobalProxyProbeRejectsMissingOrInvalidProxy(t *testing.T) {
	s := &Server{m: &Manager{}}
	for _, tc := range []struct {
		name, payload string
	}{
		{"empty", "{}"},
		{"unsupported scheme", `{"global_proxy":"ftp://localhost:42"}`},
		{"missing host", `{"global_proxy":"socks5://"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/settings/global-proxy/test", strings.NewReader(tc.payload))
			w := httptest.NewRecorder()
			s.testGlobalProxy(w, r)
			var got struct{ OK bool `json:"ok"` }
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil || w.Code != http.StatusOK || got.OK {
				t.Fatalf("status=%d payload=%s decode=%v", w.Code, w.Body.String(), err)
			}
		})
	}
}
