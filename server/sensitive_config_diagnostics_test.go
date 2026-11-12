package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/config"
)

func TestSensitiveConfigDiagnosticsAuthentication(t *testing.T) {
	t.Setenv("ARTEX_CONFIG", filepath.Join(t.TempDir(), "absent.json"))
	t.Setenv("ARTEX_PG_DSN", "postgres://user:synthetic-password@localhost/db")
	t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN", "synthetic-bridge-token")
	s := &Server{jwtKey: []byte("01234567890123456789012345678901")}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/privacy/sensitive-config", s.sensitiveConfigDiagnostics)
	handler := s.requireAuth(mux)
	request := func(jwt, bridge string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/privacy/sensitive-config", nil)
		if jwt != "" { r.Header.Set("Authorization", "Bearer "+jwt) }
		if bridge != "" { r.Header.Set("X-Artex-Privacy-Token", bridge) }
		handler.ServeHTTP(w, r)
		return w
	}
	if got := request("", ""); got.Code != http.StatusUnauthorized { t.Fatalf("anonymous status=%d", got.Code) }
	if got := request("", "synthetic-bridge-token"); got.Code != http.StatusUnauthorized { t.Fatalf("scoped token status=%d", got.Code) }
	token, err := signJWT(s.jwtKey)
	if err != nil { t.Fatal(err) }
	got := request(token, "")
	if got.Code != http.StatusOK || got.Header().Get("Cache-Control") != "no-store" { t.Fatalf("status=%d", got.Code) }
	var status config.SensitiveConfigStatus
	if err := json.Unmarshal(got.Body.Bytes(), &status); err != nil { t.Fatal(err) }
	if status.ProfilesState != "unavailable" || status.DatabaseSource != "environment" || status.FileState != "missing" { t.Fatalf("unexpected status: %+v", status) }
	if strings.Contains(got.Body.String(), "synthetic-") || strings.Contains(got.Body.String(), "localhost") { t.Fatal("diagnostics exposed values") }
}
