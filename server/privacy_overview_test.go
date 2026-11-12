package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/config"
	"github.com/Autumn-27/artex/db"
)

func TestPrivacyOverviewAggregation(t *testing.T) {
	keys := summarizeModelKeys(db.ModelKeyStatus{
		KeyState: "configured", Total: 4, Empty: 1, PlaintextOnly: 1,
		EncryptedOnly: 1, Mixed: 1, Verified: 1, Failed: 1,
		Profiles: []db.ModelKeyProfile{{ID: "987654321", Storage: "encrypted_only"}},
	})
	cfg := summarizeConfig(config.SensitiveConfigStatus{
		FileState: "readable", ProfilesState: "partial", ProfilesChecked: 2,
		Findings: []config.ConfigFinding{
			{Code: "url_userinfo", ProfileID: "987654321"},
			{Code: "query_not_inspected"},
		},
	})
	out := summarizeOutbound(agent.ModelOutboundSnapshot{
		Total: 10, Capacity: 256, Events: []agent.ModelOutboundEvent{
			{Host: "private.example.test", Result: "http_response", Status: 200},
			{Result: "http_response", Status: 429, ProxyConfigured: true},
			{Result: "transport_error"}, {Result: "canceled"}, {Result: "timeout"},
		},
	})
	if keys.State != "complete" || keys.Total != 4 || keys.Failed != 1 || keys.PlaintextOnly != 1 {
		t.Fatalf("incorrect key summary: %+v", keys)
	}
	if cfg.ProfilesState != "partial" || cfg.Findings != 2 || cfg.IncompleteFindings != 1 {
		t.Fatalf("incorrect config summary: %+v", cfg)
	}
	if out.Total != 10 || out.Retained != 5 || out.HTTPResponses != 2 || out.HTTPNon2xx != 1 ||
		out.TransportFailures != 1 || out.Canceled != 1 || out.Timeouts != 1 || out.ProxyConfigured != 1 {
		t.Fatalf("incorrect outbound summary: %+v", out)
	}
	raw, err := json.Marshal(PrivacyOverview{ModelKeys: keys, Config: cfg, Outbound: out})
	if err != nil { t.Fatal(err) }
	for _, secret := range []string{"987654321", "private.example.test", "url_userinfo"} {
		if strings.Contains(string(raw), secret) { t.Fatal("overview exposed detail data") }
	}
}

func TestPrivacyOverviewAuthenticationAndUnavailableDatabase(t *testing.T) {
	t.Setenv("ARTEX_CONFIG", filepath.Join(t.TempDir(), "absent.json"))
	t.Setenv("ARTEX_PG_DSN", "")
	t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN", "synthetic-bridge-token")
	s := &Server{jwtKey: []byte("01234567890123456789012345678901")}
	s.SetBoundListenAddress("127.0.0.1:8787")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/privacy/overview", s.privacyOverview)
	handler := s.requireAuth(mux)
	request := func(jwt, bridge string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/privacy/overview", nil)
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
	var snap PrivacyOverview
	if err := json.Unmarshal(got.Body.Bytes(), &snap); err != nil { t.Fatal(err) }
	if snap.ModelKeys.State != "unavailable" || snap.ModelKeys.KeyState != "not_inspected" ||
		snap.Config.ProfilesState != "unavailable" || snap.Listener.Scope != "loopback" ||
		snap.Listener.BindingSource != "runtime_listener" || snap.Listener.PublicReachability != "not_tested" ||
		snap.CheckedAt.IsZero() {
		t.Fatalf("unexpected partial overview: %+v", snap)
	}
	if strings.Contains(got.Body.String(), "127.0.0.1") || strings.Contains(got.Body.String(), "absent.json") {
		t.Fatal("overview exposed raw address or file path")
	}
}
