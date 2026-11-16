package report

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
	"strings"
	"testing"
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
