package server

import (
	"encoding/json"
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
