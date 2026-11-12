package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Autumn-27/artex/agent"
)

func TestModelOutboundAuditAuthentication(t *testing.T) {
	const bridgeToken = "test-read-only-bridge-token"
	t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN", bridgeToken)
	s := &Server{jwtKey: []byte("01234567890123456789012345678901")}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/privacy/model-outbound", s.modelOutboundAudit)
	handler := s.requireAuth(mux)
	request := func(jwt, bridge string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/privacy/model-outbound", nil)
		if jwt != "" { r.Header.Set("Authorization", "Bearer "+jwt) }
		if bridge != "" { r.Header.Set("X-Artex-Privacy-Token", bridge) }
		handler.ServeHTTP(w, r)
		return w
	}
	if got := request("", ""); got.Code != http.StatusUnauthorized { t.Fatalf("anonymous status=%d", got.Code) }
	if got := request("", bridgeToken); got.Code != http.StatusUnauthorized { t.Fatalf("scoped token status=%d", got.Code) }
	token, err := signJWT(s.jwtKey)
	if err != nil { t.Fatal(err) }
	got := request(token, "")
	if got.Code != http.StatusOK || got.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("authenticated status=%d cache=%q", got.Code, got.Header().Get("Cache-Control"))
	}
	var snapshot agent.ModelOutboundSnapshot
	if err := json.Unmarshal(got.Body.Bytes(), &snapshot); err != nil { t.Fatal(err) }
	if snapshot.Capacity != 256 || snapshot.Events == nil { t.Fatal("invalid audit snapshot") }
}
