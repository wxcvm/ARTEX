package knowledge

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func fixture() Catalog {
	src := Source{"knowledge.db", strings.Repeat("a", 64)}
	body := "中文材料\n<script>not executed</script>"
	e := Entry{Code: "demo", Title: "中文标题", Content: body, ContentSHA256: digest(body), SourceVerified: true,
		Attachment: Attachment{Path: "pocs/demo.txt", Content: "附件", SHA256: digest("附件"), Status: "conflict"}}
	e.ID = digest(src.SHA256 + "\x00" + e.Code)[:32]
	return Catalog{1, src, Summary{1, 1, 0, 1}, []Entry{e}}
}
func encoded(t *testing.T, c Catalog) []byte {
	t.Helper()
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func TestReadAndSearch(t *testing.T) {
	c, err := Read(bytes.NewReader(encoded(t, fixture())))
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range []string{"中文", "DEMO", ""} {
		items, total := c.Search(q, true, 0, 20)
		if len(items) != 1 || total != 1 || !items[0].SourceVerified {
			t.Fatalf("search %q: %+v", q, items)
		}
		wire, _ := json.Marshal(items)
		if bytes.Contains(wire, []byte("<script>")) || bytes.Contains(wire, []byte(`"content"`)) {
			t.Fatal("list leaked body")
		}
	}
	items, total := c.Search("", false, 10, 20)
	if len(items) != 0 || total != 1 || items == nil {
		t.Fatal("pagination must return empty array and full count")
	}
	if c.Entries[0].Attachment.Content != "附件" {
		t.Fatal("conflict lost")
	}
}
func TestReadRejectsInconsistentCatalog(t *testing.T) {
	for name, change := range map[string]func(*Catalog){
		"version":          func(c *Catalog) { c.Version = 2 },
		"source path":      func(c *Catalog) { c.Source.Name = "/tmp/knowledge.db" },
		"summary":          func(c *Catalog) { c.Summary.Conflicts = 0 },
		"identity":         func(c *Catalog) { c.Entries[0].ID = "other" },
		"duplicate":        func(c *Catalog) { c.Entries = append(c.Entries, c.Entries[0]) },
		"body hash":        func(c *Catalog) { c.Entries[0].Content = "changed" },
		"attachment hash":  func(c *Catalog) { c.Entries[0].Attachment.SHA256 = "other" },
		"fake match":       func(c *Catalog) { c.Entries[0].Attachment.Status = "match" },
		"unsafe path":      func(c *Catalog) { c.Entries[0].Attachment.Path = "pocs/../outside" },
		"unavailable body": func(c *Catalog) { c.Entries[0].Attachment.Status = "missing" },
		"unknown status":   func(c *Catalog) { c.Entries[0].Attachment.Status = "verified" },
		"material size":    func(c *Catalog) { c.Entries[0].Content = strings.Repeat("x", MaxMaterialBytes+1) },
	} {
		t.Run(name, func(t *testing.T) {
			c := fixture()
			change(&c)
			if _, err := Read(bytes.NewReader(encoded(t, c))); err == nil {
				t.Fatal("accepted invalid catalog")
			}
		})
	}
	for _, b := range [][]byte{[]byte(`{broken`), {0xff}, []byte(strings.Repeat(" ", MaxCatalogBytes+1))} {
		if _, err := Read(bytes.NewReader(b)); err == nil {
			t.Fatal("accepted malformed catalog")
		}
	}
}
func TestReadMissingAndNone(t *testing.T) {
	for _, status := range []string{"none", "missing", "ambiguous"} {
		c := fixture()
		c.Entries[0].Attachment = Attachment{Status: status}
		c.Summary.Conflicts = 0
		if status != "none" {
			c.Summary.Missing = 1
		}
		if _, err := Read(bytes.NewReader(encoded(t, c))); err != nil {
			t.Fatal(err)
		}
	}
}
