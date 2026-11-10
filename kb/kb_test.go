package kb

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func testKnowledgeDir(t *testing.T, fts bool) string {
	t.Helper()
	dir := t.TempDir()
	db, err := sql.Open("sqlite", filepath.Join(dir, "knowledge.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, statement := range []string{
		`CREATE TABLE poc (id INTEGER PRIMARY KEY, title TEXT, category TEXT, kind TEXT, tags TEXT, verified INTEGER, description TEXT, usage TEXT, content TEXT)`,
		`INSERT INTO poc (id, title, category, kind, tags, verified, description, usage, content) VALUES (1, 'test service', 'methodology', 'markdown', 'sample', 0, 'test description', 'review', 'test content')`,
	} {
		if _, err := db.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	if fts {
		for _, statement := range []string{
			`CREATE VIRTUAL TABLE poc_fts USING fts5(poc_id UNINDEXED, title, content)`,
			`INSERT INTO poc_fts (poc_id, title, content) VALUES (1, 'test service', 'test content')`,
		} {
			if _, err := db.Exec(statement); err != nil {
				t.Fatal(err)
			}
		}
	}
	return dir
}

func TestOpenAndSearch(t *testing.T) {
	for _, fts := range []bool{false, true} {
		t.Run(map[bool]string{false: "like", true: "fts"}[fts], func(t *testing.T) {
			k, err := Open(testKnowledgeDir(t, fts))
			if err != nil {
				t.Fatal(err)
			}
			defer k.Close()
			if k.HasFTS() != fts {
				t.Fatalf("HasFTS=%v, want %v", k.HasFTS(), fts)
			}
			hits, degraded, err := k.Search("test", 5)
			if err != nil {
				t.Fatal(err)
			}
			if len(hits) != 1 || hits[0].ID != 1 {
				t.Fatalf("unexpected hits: %#v", hits)
			}
			if degraded == fts {
				t.Fatalf("degraded=%v, want %v", degraded, !fts)
			}
			entry, err := k.Get(1)
			if err != nil || entry == nil || entry.Content != "test content" {
				t.Fatalf("unexpected entry: %#v, err=%v", entry, err)
			}
		})
	}
}

func TestMissingDatabase(t *testing.T) {
	k, err := Open(t.TempDir())
	if err != nil || k != nil {
		t.Fatalf("missing database: handle=%v, err=%v", k, err)
	}
}

func TestRepositoryKnowledge(t *testing.T) {
	k, err := Open(filepath.Join("..", "knowledge"))
	if err != nil || k == nil {
		t.Fatalf("open repository knowledge.db: handle=%v, err=%v", k, err)
	}
	defer k.Close()
	hits, _, err := k.Search("1panel", 8)
	if err != nil || len(hits) == 0 {
		t.Fatalf("search repository knowledge.db: hits=%d, err=%v", len(hits), err)
	}
	entry, err := k.Get(hits[0].ID)
	if err != nil || entry == nil || entry.Title == "" {
		t.Fatalf("get repository knowledge entry: entry=%v, err=%v", entry, err)
	}
}
