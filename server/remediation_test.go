package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
	"github.com/golang-jwt/jwt/v5"
)

type remediationTestStore struct {
	listCalls   int
	detailCalls int
	writeCalls  int
	filter      db.RemediationFilter
	id          string
	before      int64
	change      db.RemediationTransition
	err         error
}

func (s *remediationTestStore) ListFindingRemediations(_ context.Context, f db.RemediationFilter) ([]db.FindingRemediation, int64, error) {
	s.listCalls++
	s.filter = f
	return nil, 0, s.err
}

func (s *remediationTestStore) GetFindingRemediation(_ context.Context, id string, before int64) (db.FindingRemediationDetail, error) {
	s.detailCalls++
	s.id, s.before = id, before
	return db.FindingRemediationDetail{}, s.err
}

func (s *remediationTestStore) TransitionFindingRemediation(_ context.Context, id string, change db.RemediationTransition) (db.FindingRemediationDetail, error) {
	s.writeCalls++
	s.id, s.change = id, change
	return db.FindingRemediationDetail{}, s.err
}

func remediationTestHandler(t *testing.T) (http.Handler, []byte, string, *remediationTestStore) {
	t.Helper()
	store := &remediationTestStore{}
	key := []byte("test-only remediation JWT signing key")
	token, err := signJWT(key)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	h := remediationHTTP{
		store:  func(http.ResponseWriter) remediationStore { return store },
		jwtKey: key,
	}
	h.register(mux)
	s := &Server{jwtKey: key}
	return s.requireAuth(mux), key, token, store
}

func remediationRequest(h http.Handler, method, path, token, ct, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestRemediationHTTPRequiresNormalAuth(t *testing.T) {
	h, _, _, store := remediationTestHandler(t)
	t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN", "diagnostic-test-token")
	for _, route := range []struct{ method, path string }{
		{http.MethodGet, "/api/remediation"},
		{http.MethodGet, "/api/exploration/findings/42/remediation"},
		{http.MethodPost, "/api/exploration/findings/42/remediation/events"},
	} {
		t.Run(route.path, func(t *testing.T) {
			req := httptest.NewRequest(route.method, route.path, nil)
			req.Header.Set("X-Artex-Privacy-Token", "diagnostic-test-token")
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("bridge bypassed login: %d %s", w.Code, w.Body)
			}
		})
	}
	if store.listCalls+store.detailCalls+store.writeCalls != 0 {
		t.Fatalf("unauthenticated calls reached storage: %#v", store)
	}
}

func TestRemediationHTTPReadsAndBoundsFilters(t *testing.T) {
	h, _, token, store := remediationTestHandler(t)
	w := remediationRequest(h, "GET", "/api/remediation", token, "", "")
	if w.Code != http.StatusOK || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("list: %d %s", w.Code, w.Body)
	}
	var page remediationPage
	if err := json.Unmarshal(w.Body.Bytes(), &page); err != nil || page.Items == nil ||
		page.Page != 1 || page.PageSize != 25 || page.Total != 0 {
		t.Fatalf("expected a stable empty page, got %#v, error %v", page, err)
	}
	w = remediationRequest(h, "GET", "/api/remediation?stage=pending_fix&severity=high&q=hello&page=2&limit=30", token, "", "")
	if w.Code != http.StatusOK ||
		store.filter != (db.RemediationFilter{Stage: "pending_fix", Severity: "high", Query: "hello", Page: 2, Limit: 30}) {
		t.Fatalf("bad list response or filters: %d %s, %#v", w.Code, w.Body, store.filter)
	}
	w = remediationRequest(h, "GET", "/api/exploration/findings/42/remediation?before=10", token, "", "")
	if w.Code != http.StatusOK || store.id != "42" || store.before != 10 || store.writeCalls != 0 {
		t.Fatalf("bad read: %d %s, %#v", w.Code, w.Body, store)
	}
	for _, path := range []string{
		"/api/remediation?stage=fixed", "/api/remediation?severity=made_up",
		"/api/remediation?limit=101", "/api/remediation?page=0",
		"/api/remediation?q=" + strings.Repeat("x", 201),
		"/api/exploration/findings/42/remediation?before=-1",
		"/api/exploration/findings/42/remediation?context_task=10",
		"/api/exploration/findings/0/remediation",
	} {
		if w := remediationRequest(h, "GET", path, token, "", ""); w.Code != http.StatusBadRequest {
			t.Fatalf("%s: expected 400, got %d %s", path, w.Code, w.Body)
		}
	}
	if store.listCalls != 2 || store.detailCalls != 1 || store.writeCalls != 0 {
		t.Fatalf("invalid requests reached storage: %#v", store)
	}
}

func TestRemediationHTTPValidatesActionsAndUsesSignedActor(t *testing.T) {
	h, key, token, store := remediationTestHandler(t)
	path := "/api/exploration/findings/42/remediation/events"
	for _, c := range []struct {
		name, ct, body string
		status         int
	}{
		{"empty", "application/json", "null", 400},
		{"method not content type", "text/plain", `{"action":"review","expected_version":0,"note":"Reviewed"}`, 415},
		{"missing version", "application/json", `{"action":"review","note":"Reviewed"}`, 400},
		{"negative version", "application/json", `{"action":"review","expected_version":-1,"note":"Reviewed"}`, 400},
		{"wrong action fields", "application/json", `{"action":"review","expected_version":0,"note":"Reviewed","reference":"ticket"}`, 400},
		{"unknown action", "application/json", `{"action":"other","expected_version":0,"note":"Reviewed"}`, 400},
		{"spoofed actor", "application/json", `{"action":"review","expected_version":0,"note":"Reviewed","actor":"someone"}`, 400},
		{"extra raw field", "application/json", `{"action":"review","expected_version":0,"note":"Reviewed","attachment":"raw bytes not allowed"}`, 400},
		{"trailing value", "application/json", `{"action":"review","expected_version":0,"note":"Reviewed"} {}`, 400},
		{"missing submit reference", "application/json", `{"action":"submit","expected_version":0,"expected_evidence_version":2,"note":"Changed"}`, 400},
		{"no acceptance evidence version", "application/json", `{"action":"accept","expected_version":0,"note":"Checked"}`, 400},
		{"real calendar date", "application/json", `{"action":"update_plan","expected_version":0,"plan":"Update config","due_date":"2026-02-30"}`, 400},
		{"reference control", "application/json", `{"action":"submit","expected_version":0,"expected_evidence_version":2,"note":"Changed","reference":"work\u0000order"}`, 400},
		{"plain text bounds", "application/json", fmt.Sprintf(`{"action":"review","expected_version":0,"note":"%s"}`, strings.Repeat("测", 2001)), 400},
		{"request size bound", "application/json", fmt.Sprintf(`{"action":"review","expected_version":0,"note":"%s"}`, strings.Repeat("x", maxRemediationEventBody)), 413},
	} {
		t.Run(c.name, func(t *testing.T) {
			w := remediationRequest(h, "POST", path, token, c.ct, c.body)
			if w.Code != c.status {
				t.Fatalf("got %d want %d: %s", w.Code, c.status, w.Body)
			}
		})
	}
	withoutSubject, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}).SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	body := `{"action":"review","expected_version":0,"note":" Reviewed "}`
	if w := remediationRequest(h, "POST", path, withoutSubject, "application/json", body); w.Code != http.StatusUnauthorized {
		t.Fatalf("mutation without signed identity: %d %s", w.Code, w.Body)
	}
	if store.writeCalls != 0 {
		t.Fatalf("invalid writes reached storage: %#v", store)
	}
	w := remediationRequest(h, "POST", path, token, "application/json; charset=utf-8", body)
	if w.Code != http.StatusOK || store.writeCalls != 1 || store.id != "42" || store.change.Actor != "ARTEX" ||
		store.change.Note != "Reviewed" || store.change.ExpectedVersion == nil || *store.change.ExpectedVersion != 0 {
		t.Fatalf("valid event was not persisted as signed actor: %d %s %#v", w.Code, w.Body, store)
	}
	body = `{"action":"submit","expected_version":2,"expected_evidence_version":0,"note":"Changed","reference":"changes/45"}`
	w = remediationRequest(h, "POST", path, token, "application/json", body)
	if w.Code != http.StatusOK || store.change.ExpectedEvidenceVersion == nil || *store.change.ExpectedEvidenceVersion != 0 {
		t.Fatalf("expected explicit zero evidence version, got %d %s %#v", w.Code, w.Body, store.change)
	}
}

func TestRemediationHTTPDistinguishesConflictFromReadFailure(t *testing.T) {
	h, _, token, store := remediationTestHandler(t)
	post := "/api/exploration/findings/42/remediation/events"
	body := `{"action":"review","expected_version":0,"note":"Reviewed"}`
	for _, c := range []struct {
		err  error
		want int
	}{
		{db.ErrRemediationNotFound, 404},
		{db.ErrRemediationConflict, 409},
		{db.ErrRemediationReadOnly, 409},
		{db.ErrRemediationInvalid, 400},
		{fmt.Errorf("private database driver details"), 503},
	} {
		store.err = fmt.Errorf("private storage context: %w", c.err)
		w := remediationRequest(h, "POST", post, token, "application/json", body)
		if w.Code != c.want || strings.Contains(w.Body.String(), "private") {
			t.Fatalf("error mismatch or private details leaked: %d %s", w.Code, w.Body)
		}
	}
	w := remediationRequest(h, "GET", "/api/remediation", token, "", "")
	if w.Code != http.StatusServiceUnavailable || strings.Contains(w.Body.String(), `"items"`) {
		t.Fatalf("DB error rendered as healthy empty list: %d %s", w.Code, w.Body)
	}
}
