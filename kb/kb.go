// Package kb exposes the deployment knowledge base (knowledge/knowledge.db) to
// the agents as two read-only tools.
//
// The file is a SQLite database produced by the deployment package: table `poc`
// plus an FTS5 index `poc_fts`. Nothing here writes: the handle is opened
// read-only and every failure is returned to the caller instead of being
// swallowed, so a broken knowledge base shows up in the logs at startup.
package kb

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// KB is a read-only handle on knowledge.db.
type KB struct {
	db   *sql.DB
	path string
	fts  bool // poc_fts exists -> MATCH search is available
}

// Hit is one search result (no body: call Get for the full text).
type Hit struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Tags     string `json:"tags,omitempty"`
	Tier     string `json:"tier"` // ammo | playbook | negative | reference
	Verified bool   `json:"verified"`
	Snippet  string `json:"snippet,omitempty"`
}

// Entry is one full record.
type Entry struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Tags     string `json:"tags,omitempty"`
	Verified bool   `json:"verified"`
	Usage    string `json:"usage,omitempty"`
	Content  string `json:"content"`
}

// Open loads <dir>/knowledge.db. A missing file returns (nil, nil): the feature
// is optional, so callers only check for a non-nil handle.
func Open(dir string) (*KB, error) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return nil, nil
	}
	path := filepath.Join(dir, "knowledge.db")
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	dsn := fmt.Sprintf("file:%s?mode=ro&_pragma=busy_timeout(2000)", filepath.ToSlash(path))
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	k := &KB{db: db, path: path}
	k.fts = k.hasFTS()
	return k, nil
}

// Path is the file backing this handle.
func (k *KB) Path() string { return k.path }

// HasFTS reports whether the full-text index is present.
func (k *KB) HasFTS() bool { return k.fts }

// Close releases the handle. Safe on a nil receiver.
func (k *KB) Close() error {
	if k == nil || k.db == nil {
		return nil
	}
	return k.db.Close()
}

func (k *KB) hasFTS() bool {
	var n int
	if err := k.db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name = 'poc_fts'`).Scan(&n); err != nil {
		return false
	}
	return n > 0
}

// ftsQuery turns free text into an FTS5 AND query: every whitespace-separated
// term becomes its own phrase, so '1panel 验证码' means both must appear. Quotes
// are stripped so a model-supplied string cannot inject operators (NEAR, *, -)
// that would error out or scan the whole index.
func ftsQuery(q string) string { return joinTerms(q, " AND ") }

// orQuery is the looser form, used only when the AND query finds nothing.
func orQuery(q string) string { return joinTerms(q, " OR ") }

func joinTerms(q, op string) string {
	fields := strings.Fields(strings.ReplaceAll(q, `"`, " "))
	if len(fields) == 0 {
		return `""`
	}
	for i, f := range fields {
		fields[i] = `"` + f + `"`
	}
	return strings.Join(fields, op)
}

// tierOf maps a row onto the three tiers the prompts refer to: ammo (verified,
// directly reusable), playbook (methodology), negative (already proven dead) and
// reference (everything else).
func tierOf(verified bool, category, kind, title string) string {
	low := strings.ToLower(title + " " + category + " " + kind)
	for _, neg := range []string{"判死", "负面", "证否", "全否", "不成立", "negative"} {
		if strings.Contains(low, neg) {
			return "negative"
		}
	}
	if verified {
		return "ammo"
	}
	if category == "methodology" || kind == "markdown" {
		return "playbook"
	}
	return "reference"
}

// Search runs a full-text query: AND first, then the looser OR form, and only
// then a LIKE scan. degraded=true means FTS was unavailable or matched nothing,
// so the LIKE scan answered (slower, still correct).
func (k *KB) Search(q string, limit int) (hits []Hit, degraded bool, err error) {
	q = strings.TrimSpace(q)
	if q == "" {
		return nil, false, errors.New("q is required")
	}
	if limit <= 0 {
		limit = 8
	}
	if limit > 30 {
		limit = 30
	}
	if k.fts {
		for _, match := range []string{ftsQuery(q), orQuery(q)} {
			got, qerr := k.searchFTS(match, limit)
			if qerr != nil {
				continue // syntax/index problem -> try the next form, then LIKE
			}
			if len(got) > 0 {
				return got, false, nil
			}
		}
	}
	got, lerr := k.searchLike(q, limit)
	return got, true, lerr
}

func (k *KB) searchFTS(match string, limit int) ([]Hit, error) {
	rows, err := k.db.Query(
		`SELECT p.id, p.title, IFNULL(p.category,''), IFNULL(p.kind,''), IFNULL(p.tags,''),
		        IFNULL(p.verified,0), substr(IFNULL(p.content,''),1,240)
		   FROM poc_fts f JOIN poc p ON p.id = f.poc_id
		  WHERE poc_fts MATCH ? ORDER BY rank LIMIT ?`, match, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanHits(rows)
}

func (k *KB) searchLike(q string, limit int) ([]Hit, error) {
	like := "%" + q + "%"
	rows, err := k.db.Query(
		`SELECT id, title, IFNULL(category,''), IFNULL(kind,''), IFNULL(tags,''),
		        IFNULL(verified,0), substr(IFNULL(content,''),1,240)
		   FROM poc
		  WHERE title LIKE ? OR IFNULL(description,'') LIKE ? OR IFNULL(content,'') LIKE ?
		  LIMIT ?`, like, like, like, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanHits(rows)
}

func scanHits(rows *sql.Rows) ([]Hit, error) {
	var hits []Hit
	for rows.Next() {
		var (
			h        Hit
			verified int
		)
		if err := rows.Scan(&h.ID, &h.Title, &h.Category, &h.Kind, &h.Tags, &verified, &h.Snippet); err != nil {
			return nil, err
		}
		h.Verified = verified == 1
		h.Tier = tierOf(h.Verified, h.Category, h.Kind, h.Title)
		hits = append(hits, h)
	}
	return hits, rows.Err()
}

// Get returns one full record, or (nil, nil) when the id does not exist.
func (k *KB) Get(id int64) (*Entry, error) {
	var (
		e        Entry
		verified int
	)
	err := k.db.QueryRow(
		`SELECT id, title, IFNULL(category,''), IFNULL(kind,''), IFNULL(tags,''),
		        IFNULL(verified,0), IFNULL(usage,''), IFNULL(content,'')
		   FROM poc WHERE id = ?`, id).
		Scan(&e.ID, &e.Title, &e.Category, &e.Kind, &e.Tags, &verified, &e.Usage, &e.Content)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e.Verified = verified == 1
	return &e, nil
}
