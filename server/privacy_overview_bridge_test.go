package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrivacyOverviewBridgeScope(t *testing.T) {
	t.Setenv("ARTEX_CONFIG", filepath.Join(t.TempDir(), "absent.json"))
	t.Setenv("ARTEX_PG_DSN", "")
	t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN", "listener-only-test")
	t.Setenv("ARTEX_PRIVACY_OVERVIEW_TOKEN", "overview-only-test")
	s := &Server{jwtKey: []byte("01234567890123456789012345678901")}
	s.SetBoundListenAddress("127.0.0.1:8787")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/privacy/bridge-overview", s.privacyBridgeOverview)
	mux.HandleFunc("GET /api/privacy/bridge-status", s.privacyBridgeStatus)
	mux.HandleFunc("GET /api/privacy/overview", s.privacyOverview)
	mux.HandleFunc("POST /api/privacy/model-keys/migrate", s.migrateModelKeys)
	h := s.requireAuth(mux)
	request := func(method, path, header, token string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(method, path, nil)
		if header != "" { r.Header.Set(header, token) }
		h.ServeHTTP(w, r)
		return w
	}
	const path = "/api/privacy/bridge-overview"
	const header = "X-Artex-Privacy-Overview-Token"
	for _, tc := range []struct{ method, path, header, token string }{
		{"GET", path, "", ""},
		{"GET", path, header, "wrong"},
		{"GET", path, "X-Artex-Privacy-Token", "listener-only-test"},
		{"POST", path, header, "overview-only-test"},
		{"GET", "/api/privacy/overview", header, "overview-only-test"},
		{"GET", "/api/privacy/bridge-status", header, "overview-only-test"},
		{"POST", "/api/privacy/model-keys/migrate", header, "overview-only-test"},
	} {
		if got := request(tc.method, tc.path, tc.header, tc.token); got.Code != http.StatusUnauthorized {
			t.Fatalf("unexpected access to %s %s: %d", tc.method, tc.path, got.Code)
		}
	}
	got := request("GET", path, header, "overview-only-test")
	if got.Code != http.StatusOK || got.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("overview status=%d", got.Code)
	}
	var body struct {
		Service string `json:"service"`
		SchemaVersion int `json:"schema_version"`
		Overview PrivacyOverview `json:"overview"`
	}
	if err := json.Unmarshal(got.Body.Bytes(), &body); err != nil { t.Fatal(err) }
	if body.Service != "artex" || body.SchemaVersion != 1 || body.Overview.Listener.Scope != "loopback" ||
		body.Overview.ModelKeys.State != "unavailable" {
		t.Fatalf("invalid overview: %+v", body)
	}
	for _, value := range []string{"127.0.0.1", "overview-only-test", "absent.json"} {
		if strings.Contains(got.Body.String(), value) { t.Fatal("raw data in overview") }
	}
	t.Setenv("ARTEX_PRIVACY_OVERVIEW_TOKEN", "")
	if got := request("GET", path, header, "overview-only-test"); got.Code != http.StatusUnauthorized {
		t.Fatal("disabled bridge accepted token")
	}
}
