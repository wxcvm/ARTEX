package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/knowledge"
)

func knowledgeFixture(t *testing.T) (string, string) {
	t.Helper()
	hash := func(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
	src := knowledge.Source{Name: "knowledge.db", SHA256: strings.Repeat("a", 64)}
	id := hash(src.SHA256 + "\x00" + "demo")[:32]
	e := knowledge.Entry{ID: id, Code: "demo", Title: "中文标题", Content: "正文", ContentSHA256: hash("正文"), SourceVerified: true,
		Attachment: knowledge.Attachment{Path: "pocs/demo.txt", Status: "conflict", Content: "附件", SHA256: hash("附件")}}
	c := knowledge.Catalog{Version: 1, Source: src, Summary: knowledge.Summary{Total: 1, Conflicts: 1, WithContent: 1}, Entries: []knowledge.Entry{e}}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "catalog.json")
	if err := os.WriteFile(p, b, 0600); err != nil {
		t.Fatal(err)
	}
	return p, id
}
func TestKnowledgeRoutes(t *testing.T) {
	path, id := knowledgeFixture(t)
	t.Setenv("ARTEX_KNOWLEDGE_FILE", path)
	s := &Server{jwtKey: []byte("knowledge-test-only-32-character-key")}
	h := s.Handler()
	token, err := signJWT(s.jwtKey)
	if err != nil {
		t.Fatal(err)
	}
	request := func(path, token string) *httptest.ResponseRecorder {
		r := httptest.NewRequest("GET", path, nil)
		if token != "" {
			r.Header.Set("Authorization", "Bearer "+token)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}
	for _, url := range []string{"/api/knowledge/status", "/api/knowledge", "/api/knowledge/entries/" + id} {
		if w := request(url, ""); w.Code != 401 {
			t.Fatalf("unauth %s: %d", url, w.Code)
		}
		if w := request(url, token); w.Code != 200 || w.Header().Get("Cache-Control") != "no-store" {
			t.Fatalf("%s: %d %s", url, w.Code, w.Body)
		}
	}
	for _, q := range []string{"?limit=0", "?limit=101", "?offset=-1", "?offset=oops", "?conflicts=yes", "?q=" + strings.Repeat("x", 513)} {
		if w := request("/api/knowledge"+q, token); w.Code != 400 {
			t.Fatalf("invalid pagination: %d", w.Code)
		}
	}
	w := request("/api/knowledge?q=DEMO&conflicts=true&limit=1", token)
	if w.Code != 200 || strings.Contains(w.Body.String(), `"content"`) || !strings.Contains(w.Body.String(), `"total":1`) {
		t.Fatalf("list: %s", w.Body)
	}
	w = request("/api/knowledge?offset=1", token)
	if !strings.Contains(w.Body.String(), `"items":[]`) {
		t.Fatalf("page: %s", w.Body)
	}
	w = request("/api/knowledge/entries/"+id, token)
	if !strings.Contains(w.Body.String(), "正文") || !strings.Contains(w.Body.String(), "附件") {
		t.Fatal("lost conflicting bodies")
	}
	if w := request("/api/knowledge/entries/missing", token); w.Code != 404 {
		t.Fatal(w.Code)
	}
	t.Setenv("ARTEX_KNOWLEDGE_FILE", "")
	if w := request("/api/knowledge/status", token); w.Code != 200 || !strings.Contains(w.Body.String(), `"enabled":false`) {
		t.Fatal("disabled state")
	}
	if w := request("/api/knowledge", token); w.Code != 503 {
		t.Fatal("disabled query must fail")
	}
	t.Setenv("ARTEX_KNOWLEDGE_FILE", path)
	if err := os.WriteFile(path, []byte(`{broken`), 0600); err != nil {
		t.Fatal(err)
	}
	if w := request("/api/knowledge/status", token); w.Code != 503 || strings.Contains(w.Body.String(), path) {
		t.Fatal("corrupt must not look empty or leak path")
	}
}
func TestKnowledgeCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	(&Server{}).knowledgeList(w, r)
	if w.Code != 408 {
		t.Fatal(w.Code)
	}
}
