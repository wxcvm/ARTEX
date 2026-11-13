package server

import (
	"crypto/sha256"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/kb"
)

func TestOnboardingKnowledgeReadOnlyAndAuthentication(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "knowledge.db")
	fixture, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		"CREATE TABLE poc (id INTEGER PRIMARY KEY, title TEXT, category TEXT, kind TEXT, tags TEXT, verified INTEGER, description TEXT, usage TEXT, content TEXT)",
		"INSERT INTO poc VALUES (1, 'Windows reference', 'reference', 'markdown', 'sample', 1, 'sample', 'read only', '<script>never execute</script>')",
	} {
		if _, err := fixture.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	fixture.Close()
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	k, err := kb.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer k.Close()
	previous := agent.Knowledge()
	agent.SetKnowledge(k)
	defer agent.SetKnowledge(previous)
	s := &Server{jwtKey: []byte("onboarding-test-key-not-a-real-key")}
	handler := s.requireAuth(http.HandlerFunc(s.onboardingKnowledge))
	token, err := signJWT(s.jwtKey)
	if err != nil {
		t.Fatal(err)
	}
	for _, scenario := range []struct {
		path          string
		authenticated bool
		status        int
		contains      string
	}{
		{"/api/knowledge?q=Windows", false, 401, "Unauthorized"},
		{"/api/knowledge?q=Windows", true, 200, "Windows reference"},
		{"/api/knowledge?id=1", true, 200, "never execute"},
		{"/api/knowledge?id=999", true, 404, ""},
		{"/api/knowledge?id=-1", true, 400, ""},
		{"/api/knowledge", true, 400, ""},
	} {
		req := httptest.NewRequest(http.MethodGet, scenario.path, nil)
		if scenario.authenticated {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		if response.Code != scenario.status || !strings.Contains(response.Body.String(), scenario.contains) {
			t.Fatalf("%s: status %d, %s", scenario.path, response.Code, response.Body.String())
		}
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if sha256.Sum256(original) != sha256.Sum256(after) {
		t.Fatal("Read-only browsing changed the ledger")
	}
	agent.SetKnowledge(nil)
	response := httptest.NewRecorder()
	s.onboardingKnowledge(response, httptest.NewRequest(http.MethodGet, "/api/knowledge?q=test", nil))
	if response.Code != 503 {
		t.Fatal("Unavailable knowledge must not appear to be an empty success")
	}
}
