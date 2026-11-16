package report

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"image"
	"image/png"
	"io"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestDOCX(t *testing.T) {
	fence := strings.Repeat(string(rune(96)), 3)
	md := "# 安全测试报告\n## Findings\n### 1. Example\n#### Linked evidence\n" + fence + "\n# literal <tag> & value\n" + fence + "\n## Coverage\n"
	data, err := DOCX(md)
	if err != nil {
		t.Fatal(err)
	}
	z, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	parts := map[string]string{}
	for _, f := range z.File {
		r, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		b, err := io.ReadAll(r)
		r.Close()
		if err != nil {
			t.Fatal(err)
		}
		dec := xml.NewDecoder(bytes.NewReader(b))
		for {
			if _, err := dec.Token(); err == io.EOF {
				break
			} else if err != nil {
				t.Fatalf("%s invalid XML: %v", f.Name, err)
			}
		}
		parts[f.Name] = string(b)
	}
	for _, name := range []string{"[Content_Types].xml", "_rels/.rels", "word/document.xml", "word/_rels/document.xml.rels", "word/styles.xml"} {
		if parts[name] == "" {
			t.Fatalf("missing %s", name)
		}
	}
	for _, want := range []string{"安全测试报告", "Findings", "1. Example", "# literal &lt;tag&gt; &amp; value", "Coverage", "Linked evidence", "w:val=\"Heading1\"", "w:val=\"Heading3\"", "w:val=\"Code\""} {
		if !strings.Contains(parts["word/document.xml"], want) {
			t.Fatalf("missing %q", want)
		}
	}
}

func TestDOCXWithScreenshots(t *testing.T) {
	var img bytes.Buffer
	if err := png.Encode(&img, image.NewRGBA(image.Rect(0, 0, 12, 8))); err != nil {
		t.Fatal(err)
	}
	f := &db.DBFinding{ID: 42, Name: "Fixture", Screenshots: []db.FindingScreenshotData{
		{FindingScreenshot: db.FindingScreenshot{MIME: "image/png", Caption: "截图 <&>", BindingID: "91"}, Image: img.Bytes()},
	}}
	raw, err := DOCXWithScreenshots("# Fixture", []*db.DBFinding{f}, "zh-CN")
	if err != nil {
		t.Fatal(err)
	}
	z, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatal(err)
	}
	parts := map[string][]byte{}
	for _, entry := range z.File {
		r, err := entry.Open()
		if err != nil {
			t.Fatal(err)
		}
		data, err := io.ReadAll(r)
		r.Close()
		if err != nil {
			t.Fatal(err)
		}
		parts[entry.Name] = data
		if strings.HasSuffix(entry.Name, ".xml") || strings.HasSuffix(entry.Name, ".rels") {
			dec := xml.NewDecoder(bytes.NewReader(data))
			for {
				if _, err := dec.Token(); err == io.EOF {
					break
				} else if err != nil {
					t.Fatalf("invalid %s: %v", entry.Name, err)
				}
			}
		}
	}
	if !bytes.Equal(parts["word/media/screenshot-1.png"], img.Bytes()) {
		t.Fatal("image bytes not embedded")
	}
	for part, values := range map[string][]string{
		"word/document.xml": {"附录：截图证据", "截图 &lt;&amp;&gt;", "关联证据 #91", "r:embed=\"rId2\"", "wp:extent"},
		"word/_rels/document.xml.rels": {"rId2", "media/screenshot-1.png", "/image"},
		"[Content_Types].xml": {"image/png"},
	} {
		for _, want := range values {
			if !bytes.Contains(parts[part], []byte(want)) {
				t.Fatalf("%s missing %q", part, want)
			}
		}
	}
}
