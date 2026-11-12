package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestStoragePermissionsAuthentication(t *testing.T) {
	t.Setenv("ARTEX_CONFIG", filepath.Join(t.TempDir(), "private-config.json"))
	t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN", "listener-fixture")
	t.Setenv("ARTEX_PRIVACY_OVERVIEW_TOKEN", "overview-fixture")
	s := &Server{jwtKey: []byte("01234567890123456789012345678901")}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/privacy/storage-permissions", s.storagePermissions)
	h := s.requireAuth(mux)
	request := func(header, token string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/privacy/storage-permissions?path=caller-supplied", nil)
		if header != "" { r.Header.Set(header, token) }
		h.ServeHTTP(w, r)
		return w
	}
	for _, tc := range []struct{ header, token string }{
		{"", ""}, {"X-Artex-Privacy-Token", "listener-fixture"},
		{"X-Artex-Privacy-Overview-Token", "overview-fixture"},
	} {
		if got := request(tc.header, tc.token); got.Code != http.StatusUnauthorized { t.Fatalf("unauthorized status=%d", got.Code) }
	}
	token, err := signJWT(s.jwtKey)
	if err != nil { t.Fatal(err) }
	got := request("Authorization", "Bearer "+token)
	if got.Code != http.StatusOK || got.Header().Get("Cache-Control") != "no-store" { t.Fatalf("status=%d", got.Code) }
	var snap StoragePermissions
	if err := json.Unmarshal(got.Body.Bytes(), &snap); err != nil { t.Fatal(err) }
	if len(snap.Items) != 3 || snap.Items[0].Kind != "config_file" || snap.Items[1].Kind != "jwt_key" || snap.Items[2].Kind != "data_root" {
		t.Fatalf("unexpected fixed scope: %+v", snap)
	}
	if strings.Contains(got.Body.String(), "private-config") || strings.Contains(got.Body.String(), "caller-supplied") {
		t.Fatal("storage diagnostics exposed a path")
	}
}
