// Package knowledge exposes read-only, provenance-preserving material previews.
package knowledge

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

const MaxCatalogBytes = 32 << 20
const MaxMaterialBytes = 2 << 20

type Source struct {
	Name   string `json:"name"`
	SHA256 string `json:"sha256"`
}
type Summary struct {
	Total       int `json:"total"`
	Conflicts   int `json:"conflicts"`
	Missing     int `json:"missing"`
	WithContent int `json:"with_content"`
}
type Attachment struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	SHA256  string `json:"sha256"`
	Status  string `json:"status"`
}
type Entry struct {
	ID             string     `json:"id"`
	Code           string     `json:"code"`
	Title          string     `json:"title"`
	Kind           string     `json:"kind"`
	Category       string     `json:"category"`
	Component      string     `json:"component"`
	Versions       string     `json:"versions"`
	CVE            string     `json:"cve"`
	Tags           string     `json:"tags"`
	Source         string     `json:"source"`
	SourceURL      string     `json:"source_url"`
	SourceVerified bool       `json:"source_verified"`
	Description    string     `json:"description"`
	Usage          string     `json:"usage"`
	Content        string     `json:"content"`
	ContentSHA256  string     `json:"content_sha256"`
	Attachment     Attachment `json:"attachment"`
}
type Catalog struct {
	Version int     `json:"version"`
	Source  Source  `json:"source"`
	Summary Summary `json:"summary"`
	Entries []Entry `json:"entries"`
}

// EntrySummary deliberately excludes both bodies and usage text.
type EntrySummary struct {
	ID               string `json:"id"`
	Code             string `json:"code"`
	Title            string `json:"title"`
	Kind             string `json:"kind"`
	Category         string `json:"category"`
	Component        string `json:"component"`
	CVE              string `json:"cve"`
	Tags             string `json:"tags"`
	SourceVerified   bool   `json:"source_verified"`
	AttachmentStatus string `json:"attachment_status"`
}

func digest(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
func validHash(s string) bool {
	b, err := hex.DecodeString(s)
	return err == nil && len(b) == 32 && s == strings.ToLower(s)
}
func safePath(s string) bool {
	if !strings.HasPrefix(s, "pocs/") || strings.ContainsAny(s, "\\\x00:") {
		return false
	}
	for _, p := range strings.Split(s, "/") {
		if p == "" || p == "." || p == ".." {
			return false
		}
	}
	return true
}
func Load(path string) (*Catalog, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read(f)
}
func Read(r io.Reader) (*Catalog, error) {
	b, err := io.ReadAll(io.LimitReader(r, MaxCatalogBytes+1))
	if err != nil {
		return nil, err
	}
	if len(b) > MaxCatalogBytes || !utf8.Valid(b) {
		return nil, errors.New("catalog exceeds 32 MiB or is not UTF-8")
	}
	var c Catalog
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	if c.Version != 1 || !validHash(c.Source.SHA256) || c.Source.Name == "" ||
		strings.ContainsAny(c.Source.Name, "/\\\x00:") || len(c.Entries) > 10000 {
		return nil, errors.New("invalid catalog version, source or entry count")
	}
	seen := make(map[string]bool, len(c.Entries))
	var actual Summary
	for _, e := range c.Entries {
		if e.Code == "" || strings.TrimSpace(e.Title) == "" || seen[e.Code] || e.ID != digest(c.Source.SHA256 + "\x00" + e.Code)[:32] {
			return nil, errors.New("invalid or duplicate entry identity")
		}
		seen[e.Code] = true
		if len(e.Content) > MaxMaterialBytes || len(e.Attachment.Content) > MaxMaterialBytes || e.ContentSHA256 != digest(e.Content) {
			return nil, errors.New("invalid material hash or size")
		}
		a := e.Attachment
		if a.Path != "" && !safePath(a.Path) {
			return nil, errors.New("invalid attachment path")
		}
		switch a.Status {
		case "none", "missing", "ambiguous":
			if a.Content != "" || a.SHA256 != "" || (a.Status == "none" && a.Path != "") {
				return nil, errors.New("unavailable attachment contains material")
			}
			if a.Status != "none" {
				actual.Missing++
			}
		case "match", "conflict":
			if a.Path == "" || a.SHA256 != digest(a.Content) || (a.Status == "match") != (a.Content == e.Content) {
				return nil, errors.New("attachment comparison or hash is inconsistent")
			}
			if a.Status == "conflict" {
				actual.Conflicts++
			}
		default:
			return nil, errors.New("unknown attachment status")
		}
		actual.Total++
		if e.Content != "" {
			actual.WithContent++
		}
	}
	if actual != c.Summary {
		return nil, errors.New("catalog summary does not match entries")
	}
	sort.Slice(c.Entries, func(i, j int) bool { return c.Entries[i].Code < c.Entries[j].Code })
	return &c, nil
}
func (c *Catalog) Search(q string, conflicts bool, offset, limit int) ([]EntrySummary, int) {
	q = strings.ToLower(strings.TrimSpace(q))
	items := make([]EntrySummary, 0)
	total := 0
	for _, e := range c.Entries {
		if conflicts && e.Attachment.Status != "conflict" {
			continue
		}
		text := strings.Join([]string{e.Code, e.Title, e.Kind, e.Category, e.Component, e.Versions, e.CVE, e.Tags, e.Description}, "\n")
		if q != "" && !strings.Contains(strings.ToLower(text), q) {
			continue
		}
		total++
		if total <= offset || len(items) >= limit {
			continue
		}
		items = append(items, EntrySummary{e.ID, e.Code, e.Title, e.Kind, e.Category, e.Component, e.CVE, e.Tags, e.SourceVerified, e.Attachment.Status})
	}
	return items, total
}
