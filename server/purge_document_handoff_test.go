package server

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestPurgeDocumentHandoffUpload(t *testing.T) {
	s, f, request := trafficEvidenceServer(t)
	finding, err := s.m.pg.GetFinding(f.FindingID)
	if err != nil || finding == nil || finding.TaskID == nil {
		t.Fatal("missing task", err)
	}
	task := strconv.FormatInt(*finding.TaskID, 10)
	other, err := s.m.CreateTask("handoff isolation fixture", "documentation", nil, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	otherID, _ := strconv.ParseInt(other.ID, 10, 64)
	t.Cleanup(func() { s.m.pg.DeleteTask(otherID) })

	token, err := signJWT(s.jwtKey)
	if err != nil {
		t.Fatal(err)
	}
	handler := s.Handler()
	const base = "/api/knowledge/documents"
	type uploaded struct {
		code   int
		reused bool
		doc    db.KnowledgeDocument
		raw    string
	}
	upload := func(target, filename, content string) uploaded {
		var body bytes.Buffer
		form := multipart.NewWriter(&body)
		file, err := form.CreateFormFile("file", filename)
		if err != nil {
			panic(err)
		}
		if _, err = file.Write([]byte(content)); err != nil {
			panic(err)
		}
		if err = form.Close(); err != nil {
			panic(err)
		}
		path := base + "/upload"
		if target != "" {
			path += "?task=" + target
		}
		r := httptest.NewRequest(http.MethodPost, path, &body)
		r.Header.Set("Authorization", "Bearer "+token)
		r.Header.Set("Content-Type", form.FormDataContentType())
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		out := uploaded{code: w.Code, reused: w.Header().Get("X-ARTEX-Knowledge-Upload") == "reused", raw: w.Body.String()}
		if w.Code == http.StatusOK || w.Code == http.StatusCreated {
			if err := json.Unmarshal(w.Body.Bytes(), &out.doc); err != nil {
				panic(err)
			}
		}
		return out
	}
	name := func(text string) string {
		hash := sha256.Sum256([]byte(text))
		return fmt.Sprintf("purge-document-handoff-%x.md", hash[:8])
	}
	content := "# 文档交付\n\n附录已核对。\n"
	filename := name(content)
	// Optional real file exported by a local Purge preview. CI uses only
	// the fixture above; neither path reads any active ARTEX data.
	if path := os.Getenv("PURGE_DOCUMENT_HANDOFF_TEST_FILE"); path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		content, filename = string(raw), filepath.Base(path)
		if name(content) != filename {
			t.Fatal("unexpected handoff fixture name")
		}
	}

	if got := upload("", filename, content); got.code != 400 {
		t.Fatal("global upload accepted", got)
	}
	if got := upload(task, filename, content+"changed"); got.code != 400 {
		t.Fatal("invalid hash accepted", got)
	}

	const n = 6
	results := make(chan uploaded, n)
	for i := 0; i < n; i++ {
		go func() { results <- upload(task, filename, content) }()
	}
	created, reused := 0, 0
	var original db.KnowledgeDocument
	for i := 0; i < n; i++ {
		got := <-results
		if got.code == http.StatusCreated && !got.reused {
			created++
		} else if got.code == http.StatusOK && got.reused {
			reused++
		} else {
			t.Fatal(got)
		}
		if i == 0 {
			original = got.doc
		}
		if got.doc.ID == "" || got.doc.ID != original.ID || got.doc.TaskID != task || got.doc.Content != content {
			t.Fatal("retries did not return the same document", got)
		}
	}
	if created != 1 || reused != n-1 {
		t.Fatalf("created=%d reused=%d", created, reused)
	}
	if got := upload(task, strings.TrimSuffix(filename, ".md")+" (1).MD", "\uFEFF"+strings.ReplaceAll(content, "\n", "\r\n")); got.code != 200 || got.doc.ID != original.ID {
		t.Fatal("browser copy or newline variation was not reused", got)
	}
	if got := upload(other.ID, filename, content); got.code != 201 || got.doc.ID == original.ID || got.doc.TaskID != other.ID {
		t.Fatal("separate task was incorrectly deduplicated", got)
	}
	path := base + "/" + original.ID + "?task=" + task
	if w := request("PUT", path, `{"title":"manual edit","content":"retain this edited content"}`); w.Code != 200 {
		t.Fatal("could not set up receiver conflict", w.Code, w.Body)
	}
	if got := upload(task, filename, content); got.code != 409 {
		t.Fatal("conflict was not reported", got)
	}
	w := request("GET", path, "")
	if w.Code != 200 || !strings.Contains(w.Body.String(), "retain this edited content") {
		t.Fatal("the existing edit was overwritten", w.Code, w.Body)
	}

	legacy := "# Legacy browser copy\n"
	legacyName := name(legacy)
	legacyDoc, err := s.m.pg.SaveKnowledgeDocument(context.Background(), *finding.TaskID, "",
		db.KnowledgeDocument{Title: "Legacy", Content: legacy, SourceFilename: strings.TrimSuffix(legacyName, ".md") + " (4).md"})
	if err != nil {
		t.Fatal(err)
	}
	if got := upload(task, legacyName, legacy); got.code != 200 || got.doc.ID != legacyDoc.ID {
		t.Fatal("legacy copy not recognized", got)
	}

	for i := 0; i < 2; i++ {
		if got := upload(task, "ordinary.md", "repeatable original upload"); got.code != 201 || got.reused {
			t.Fatal("ordinary upload changed behaviour", got)
		}
	}
}
