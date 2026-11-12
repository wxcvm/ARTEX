package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestModelKeyDiagnosticsAuthentication(t *testing.T) {
	const bridgeToken = "test-bridge-token-only-for-listener"
	t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN", bridgeToken)
	s := &Server{jwtKey: []byte("01234567890123456789012345678901")}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/privacy/model-keys", s.modelKeyDiagnostics)
	handler := s.requireAuth(mux)
	request := func(jwt, bridge string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/privacy/model-keys", nil)
		if jwt != "" {
			r.Header.Set("Authorization", "Bearer "+jwt)
		}
		if bridge != "" {
			r.Header.Set("X-Artex-Privacy-Token", bridge)
		}
		handler.ServeHTTP(w, r)
		return w
	}
	if got := request("", ""); got.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous status=%d", got.Code)
	}
	if got := request("", bridgeToken); got.Code != http.StatusUnauthorized {
		t.Fatalf("scoped token accessed model diagnostics: %d", got.Code)
	}
	token, err := signJWT(s.jwtKey)
	if err != nil {
		t.Fatal(err)
	}
	got := request(token, "")
	if got.Code != http.StatusServiceUnavailable || got.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("authenticated unavailable database status=%d cache=%q", got.Code, got.Header().Get("Cache-Control"))
	}
}

func TestModelKeyMigrationAuthentication(t *testing.T) {
	const bridgeToken = "test-bridge-token-only-for-listener"
	t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN", bridgeToken)
	s := &Server{jwtKey: []byte("01234567890123456789012345678901")}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/privacy/model-keys/migrate", s.migrateModelKeys)
	handler := s.requireAuth(mux)
	request := func(jwt, bridge string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/privacy/model-keys/migrate", nil)
		if jwt != "" {
			r.Header.Set("Authorization", "Bearer "+jwt)
		}
		if bridge != "" {
			r.Header.Set("X-Artex-Privacy-Token", bridge)
		}
		handler.ServeHTTP(w, r)
		return w
	}
	if got := request("", ""); got.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous status=%d", got.Code)
	}
	if got := request("", bridgeToken); got.Code != http.StatusUnauthorized {
		t.Fatalf("scoped token accessed migration: %d", got.Code)
	}
	token, err := signJWT(s.jwtKey)
	if err != nil {
		t.Fatal(err)
	}
	got := request(token, "")
	if got.Code != http.StatusServiceUnavailable || got.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("authenticated unavailable database status=%d cache=%q", got.Code, got.Header().Get("Cache-Control"))
	}
}
