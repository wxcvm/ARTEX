package kb

import (
	"crypto/sha256"
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
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

type knowledgeTestRow struct {
	id                          int64
	title, description, content string
}

func addKnowledgeTestRows(t *testing.T, dir string, rows []knowledgeTestRow, fts bool) {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(dir, "knowledge.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, row := range rows {
		if _, err := db.Exec(`INSERT INTO poc (id, title, description, content) VALUES (?, ?, ?, ?)`, row.id, row.title, row.description, row.content); err != nil {
			t.Fatal(err)
		}
		if fts {
			if _, err := db.Exec(`INSERT INTO poc_fts (poc_id, title, content) VALUES (?, ?, ?)`, row.id, row.title, row.content); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func knowledgeHitIDs(hits []Hit) []int64 {
	ids := make([]int64, len(hits))
	for i, hit := range hits {
		ids[i] = hit.ID
	}
	return ids
}

func TestSearchLikeLiteralMetacharacters(t *testing.T) {
	dir := testKnowledgeDir(t, false)
	addKnowledgeTestRows(t, dir, []knowledgeTestRow{
		{id: 2, title: "a_b"},
		{id: 3, title: "axb"},
		{id: 4, title: "100%done"},
		{id: 5, title: "100Xdone"},
		{id: 6, title: `path\value`},
		{id: 7, title: "path/value"},
		{id: 8, title: `combo\%done`},
		{id: 9, title: `combo\Xdone`},
		{id: 10, title: "description target", description: "desc_needle"},
		{id: 11, title: "description decoy", description: "descXneedle"},
		{id: 12, title: "content target", content: "body%needle"},
		{id: 13, title: "content decoy", content: "bodyXneedle"},
		{id: 14, title: `double\\slash`},
	}, false)
	k, err := Open(dir)
	if err != nil || k == nil {
		t.Fatalf("open fixture: handle=%v, err=%v", k, err)
	}
	defer k.Close()
	if k.HasFTS() {
		t.Fatal("literal-match fixture must exercise the LIKE fallback")
	}
	for _, tc := range []struct {
		query string
		ids   []int64
	}{
		{"a_b", []int64{2}},
		{"100%done", []int64{4}},
		{"%", []int64{4, 8, 12}},
		{"_", []int64{2, 10}},
		{`\`, []int64{6, 8, 9, 14}},
		{`path\value`, []int64{6}},
		{`combo\%done`, []int64{8}},
		{`\%`, []int64{8}},
		{`\\`, []int64{14}},
		{"desc_needle", []int64{10}},
		{"body%needle", []int64{12}},
	} {
		t.Run(tc.query, func(t *testing.T) {
			hits, degraded, err := k.Search(tc.query, 30)
			if err != nil || !degraded {
				t.Fatalf("literal search %q: degraded=%v, err=%v", tc.query, degraded, err)
			}
			if ids := knowledgeHitIDs(hits); !reflect.DeepEqual(ids, tc.ids) {
				t.Fatalf("literal search %q: IDs=%v, want %v", tc.query, ids, tc.ids)
			}
		})
	}
}

func knowledgeTestTree(t *testing.T, root string) []string {
	t.Helper()
	var paths []string
	if err := filepath.WalkDir(root, func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		paths = append(paths, rel)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return paths
}

func TestOpenReservedDirectoryNamesReadOnly(t *testing.T) {
	fixture, err := os.ReadFile(filepath.Join(testKnowledgeDir(t, false), "knowledge.db"))
	if err != nil {
		t.Fatal(err)
	}
	names := []string{"hash#segment", "percent%20segment", "with space", "中文目录", "空 格#hash%20"}
	if runtime.GOOS != "windows" {
		names = append(names, "query?segment")
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, name)
			if err := os.Mkdir(dir, 0700); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, "knowledge.db")
			if err := os.WriteFile(path, fixture, 0600); err != nil {
				t.Fatal(err)
			}
			beforeHash := sha256.Sum256(fixture)
			beforeTree := knowledgeTestTree(t, root)
			// Check even failed opens: a malformed URI must not create a sibling file.
			defer func() {
				after, err := os.ReadFile(path)
				if err != nil {
					t.Errorf("read original database after Open: %v", err)
				} else if sha256.Sum256(after) != beforeHash {
					t.Error("Open/search/get changed the source database bytes")
				}
				if afterTree := knowledgeTestTree(t, root); !reflect.DeepEqual(afterTree, beforeTree) {
					t.Errorf("Open created or removed files: before=%q, after=%q", beforeTree, afterTree)
				}
			}()
			k, err := Open(dir)
			if k != nil {
				defer k.Close()
			}
			if err != nil || k == nil {
				t.Fatalf("open reserved directory %q: handle=%v, err=%v", name, k, err)
			}
			if k.Path() != path {
				t.Errorf("Path=%q, want %q", k.Path(), path)
			}
			hits, degraded, err := k.Search("test service", 5)
			if err != nil || !degraded || !reflect.DeepEqual(knowledgeHitIDs(hits), []int64{1}) {
				t.Fatalf("search exact source database: hits=%v, degraded=%v, err=%v", knowledgeHitIDs(hits), degraded, err)
			}
			entry, err := k.Get(1)
			if err != nil || entry == nil || entry.Content != "test content" {
				t.Fatalf("get exact source database: entry=%#v, err=%v", entry, err)
			}
			if _, err := k.db.Exec(`UPDATE poc SET title = 'unexpected write' WHERE id = 1`); err == nil {
				t.Fatal("knowledge handle accepted a write to the read-only fixture")
			}
		})
	}
}

func TestOpenRejectsIncompleteSchema(t *testing.T) {
	columns := []string{"id INTEGER PRIMARY KEY", "title TEXT", "category TEXT", "kind TEXT", "tags TEXT", "verified INTEGER", "description TEXT", "usage TEXT", "content TEXT"}
	for _, missing := range []string{"poc", "id", "title", "category", "kind", "tags", "verified", "description", "usage", "content"} {
		t.Run(missing, func(t *testing.T) {
			dir := t.TempDir()
			db, err := sql.Open("sqlite", filepath.Join(dir, "knowledge.db"))
			if err != nil {
				t.Fatal(err)
			}
			statement := "CREATE TABLE unrelated (id INTEGER)"
			if missing != "poc" {
				var present []string
				for _, column := range columns {
					if !strings.HasPrefix(column, missing+" ") {
						present = append(present, column)
					}
				}
				statement = "CREATE TABLE poc (" + strings.Join(present, ", ") + ")"
			}
			_, execErr := db.Exec(statement)
			closeErr := db.Close()
			if execErr != nil || closeErr != nil {
				t.Fatalf("create incomplete fixture: exec=%v, close=%v", execErr, closeErr)
			}
			k, err := Open(dir)
			if k != nil {
				defer k.Close()
				t.Error("Open returned a usable handle for an incomplete schema")
			}
			if err == nil {
				t.Fatal("Open accepted an existing database with missing " + missing)
			}
			if !strings.Contains(strings.ToLower(err.Error()), missing) {
				t.Errorf("schema error %q does not identify missing %q", err, missing)
			}
		})
	}
}

func TestSearchStableOrder(t *testing.T) {
	for _, fts := range []bool{false, true} {
		t.Run(map[bool]string{false: "like", true: "fts"}[fts], func(t *testing.T) {
			dir := testKnowledgeDir(t, fts)
			var rows []knowledgeTestRow
			for _, id := range []int64{9, 2, 7, 3} {
				rows = append(rows, knowledgeTestRow{id: id, title: "stable result", description: "stable description", content: "stable content"})
			}
			addKnowledgeTestRows(t, dir, rows, fts)
			k, err := Open(dir)
			if err != nil || k == nil {
				t.Fatalf("open fixture: handle=%v, err=%v", k, err)
			}
			defer k.Close()
			// Expose reliance on table-scan order; ordered SELECTs are unaffected.
			k.db.SetMaxOpenConns(1)
			if _, err := k.db.Exec("PRAGMA reverse_unordered_selects = ON"); err != nil {
				t.Fatal(err)
			}
			if fts {
				ranks, err := k.db.Query(`SELECT rank FROM poc_fts WHERE poc_fts MATCH '"stable"'`)
				if err != nil {
					t.Fatal(err)
				}
				var firstRank float64
				count := 0
				for ranks.Next() {
					var rank float64
					if err := ranks.Scan(&rank); err != nil {
						ranks.Close()
						t.Fatal(err)
					}
					if count == 0 {
						firstRank = rank
					} else if rank != firstRank {
						ranks.Close()
						t.Fatalf("fixture must have tied FTS ranks: %g != %g", rank, firstRank)
					}
					count++
				}
				rowsErr := ranks.Err()
				ranks.Close()
				if rowsErr != nil || count != 4 {
					t.Fatalf("tied-rank fixture: rows=%d, err=%v", count, rowsErr)
				}
			}
			for attempt := 0; attempt < 3; attempt++ {
				for _, tc := range []struct {
					limit int
					ids   []int64
				}{{2, []int64{2, 3}}, {30, []int64{2, 3, 7, 9}}} {
					hits, degraded, err := k.Search("stable", tc.limit)
					if err != nil || degraded == fts {
						t.Fatalf("search: degraded=%v, err=%v", degraded, err)
					}
					if ids := knowledgeHitIDs(hits); !reflect.DeepEqual(ids, tc.ids) {
						t.Fatalf("attempt %d, limit %d: IDs=%v, want %v", attempt, tc.limit, ids, tc.ids)
					}
				}
			}
		})
	}
}
