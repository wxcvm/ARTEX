package server

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClassifyListener(t *testing.T) {
	tests := []struct {
		addr string
		scope string
	}{
		{"127.0.0.1:8787", "loopback"},
		{"[::1]:8787", "loopback"},
		{"localhost:8787", "loopback"},
		{":8787", "all_interfaces"},
		{"0.0.0.0:8787", "all_interfaces"},
		{"[::]:8787", "all_interfaces"},
		{"192.168.1.10:8787", "specific_interface"},
		{"[2001:db8::5]:8787", "specific_interface"},
		{"unknown-host.example:8787", "unknown"},
		{"invalid", "unknown"},
	}
	for _, tc := range tests {
		t.Run(tc.addr, func(t *testing.T) {
			got := classifyListener(tc.addr)
			if got.Scope != tc.scope {
				t.Fatalf("scope=%q, want %q", got.Scope, tc.scope)
			}
			if got.PublicReachability != "not_tested" {
				t.Fatalf("cannot infer public reachability from listener: %+v", got)
			}
		})
	}
}

func TestExposureEndpointRequiresJWT(t *testing.T) {
	s := &Server{jwtKey: []byte("01234567890123456789012345678901")}
	s.SetListenAddress("0.0.0.0:8787")
	guarded := s.requireAuth(http.HandlerFunc(s.exposureStatus))

	withoutToken := httptest.NewRecorder()
	guarded.ServeHTTP(withoutToken, httptest.NewRequest(http.MethodGet, "/api/privacy/exposure", nil))
	if withoutToken.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated status=%d, want 401", withoutToken.Code)
	}

	token, err := signJWT(s.jwtKey)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/privacy/exposure", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	withToken := httptest.NewRecorder()
	guarded.ServeHTTP(withToken, request)
	if withToken.Code != http.StatusOK {
		t.Fatalf("authenticated status=%d, want 200", withToken.Code)
	}
	var info ExposureStatus
	if err := json.Unmarshal(withToken.Body.Bytes(), &info); err != nil {
		t.Fatal(err)
	}
	if info.Scope != "all_interfaces" || info.PublicReachability != "not_tested" {
		t.Fatalf("incorrect exposure assessment: %+v", info)
	}
}

func TestExposureRuntimeListener(t *testing.T) {
	for _, tc := range []struct {
		configured string
		bind       string
		scope      string
	}{
		{"localhost:0", "127.0.0.1:0", "loopback"},
		{"0.0.0.0:0", "0.0.0.0:0", "all_interfaces"},
	} {
		t.Run(tc.configured, func(t *testing.T) {
			listener, err := net.Listen("tcp4", tc.bind)
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			s := &Server{}
			s.SetListenAddress(tc.configured)
			s.SetBoundListenAddress(listener.Addr().String())
			w := httptest.NewRecorder()
			s.exposureStatus(w, httptest.NewRequest(http.MethodGet, "/api/privacy/exposure", nil))
			var info ExposureStatus
			if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil {
				t.Fatal(err)
			}
			if info.ConfiguredAddress != tc.configured || info.BoundAddress != listener.Addr().String() ||
				info.ListenAddress != info.BoundAddress || info.BindingSource != "runtime_listener" ||
				info.Scope != tc.scope || info.PublicReachability != "not_tested" {
				t.Fatalf("incorrect runtime exposure assessment: %+v", info)
			}
		})
	}
}

func TestExposureConfiguredFallback(t *testing.T) {
	s := &Server{}
	s.SetListenAddress("127.0.0.1:8787")
	w := httptest.NewRecorder()
	s.exposureStatus(w, httptest.NewRequest(http.MethodGet, "/api/privacy/exposure", nil))
	var info ExposureStatus
	if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil {
		t.Fatal(err)
	}
	if info.BindingSource != "configured" || info.BoundAddress != "" ||
		info.ConfiguredAddress != "127.0.0.1:8787" || info.Scope != "loopback" ||
		info.PublicReachability != "not_tested" {
		t.Fatalf("incorrect configured fallback: %+v", info)
	}
}

func TestPrivacyBridgeScopedToken(t *testing.T) {
	const token = "test-read-only-bridge-token-123456789"
	s := &Server{jwtKey: []byte("01234567890123456789012345678901")}
	s.SetListenAddress("0.0.0.0:8787")
	s.SetBoundListenAddress("0.0.0.0:8787")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/privacy/bridge-status", s.privacyBridgeStatus)
	mux.HandleFunc("GET /api/privacy/exposure", s.exposureStatus)
	handler := s.requireAuth(mux)
	request := func(method, path, credential string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(method, path, nil)
		if credential != "" {
			r.Header.Set("X-Artex-Privacy-Token", credential)
		}
		handler.ServeHTTP(w, r)
		return w
	}
	t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN", "")
	if got := request("GET", "/api/privacy/bridge-status", token); got.Code != http.StatusUnauthorized {
		t.Fatalf("unconfigured bridge status=%d", got.Code)
	}
	t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN", token)
	for _, wrong := range []string{"", "wrong-token"} {
		if got := request("GET", "/api/privacy/bridge-status", wrong); got.Code != http.StatusUnauthorized {
			t.Fatalf("invalid bridge credential status=%d", got.Code)
		}
	}
	got := request("GET", "/api/privacy/bridge-status", token)
	if got.Code != http.StatusOK || got.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("bridge response status=%d cache=%q", got.Code, got.Header().Get("Cache-Control"))
	}
	var info map[string]any
	if err := json.Unmarshal(got.Body.Bytes(), &info); err != nil {
		t.Fatal(err)
	}
	if len(info) != 4 || info["service"] != "artex" || info["scope"] != "all_interfaces" ||
		info["binding_source"] != "runtime_listener" || info["public_reachability"] != "not_tested" {
		t.Fatalf("unexpected or excessive bridge response: %+v", info)
	}
	for _, tc := range []struct{ method, path string }{
		{"GET", "/api/privacy/exposure"},
		{"POST", "/api/privacy/bridge-status"},
	} {
		if got := request(tc.method, tc.path, token); got.Code != http.StatusUnauthorized {
			t.Fatalf("scoped credential reached %s %s: %d", tc.method, tc.path, got.Code)
		}
	}
}
