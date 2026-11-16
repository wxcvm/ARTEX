package kb

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestExtractKnowledgeDocumentText(t *testing.T) {
	for _,name:=range []string{"notes.txt","notes.md","notes.csv","notes.json"} {
		text,err:=ExtractDocumentText(name,[]byte("\uFEFF中文 notes 100%"))
		if err!=nil||text!="中文 notes 100%"{t.Fatal(name,text,err)}
	}
	for _,tc:=range []struct{name string; raw []byte}{
		{"bad.txt",[]byte{0xff}},{"empty.md",[]byte(" \n")},{"nul.txt",[]byte{0}},
		{"big.md",bytes.Repeat([]byte("x"),DocumentTextLimit+1)},
		{"bad.pdf",[]byte("not supported")},{"bad.docx",[]byte("not zip")},
	}{
		if _,err:=ExtractDocumentText(tc.name,tc.raw);err==nil{t.Fatal("accepted",tc.name)}
	}
}

func TestExtractKnowledgeDOCX(t *testing.T) {
	var buf bytes.Buffer
	z:=zip.NewWriter(&buf)
	w,err:=z.Create("word/document.xml")
	if err!=nil{t.Fatal(err)}
	_,err=w.Write([]byte(`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>中文 &amp; notes</w:t><w:tab/><w:t>row</w:t></w:r></w:p><w:tbl><w:tr><w:tc><w:p><w:r><w:t>cell text</w:t></w:r></w:p></w:tc></w:tr></w:tbl></w:body></w:document>`))
	if err!=nil{t.Fatal(err)}
	w,err=z.Create("word/other.xml");if err!=nil{t.Fatal(err)}
	if _,err=w.Write([]byte("not main document"));err!=nil{t.Fatal(err)}
	if err=z.Close();err!=nil{t.Fatal(err)}
	text,err:=ExtractDocumentText("fixture.docx",buf.Bytes())
	if err!=nil||!strings.Contains(text,"中文 & notes\trow\ncell text")||strings.Contains(text,"not main document"){t.Fatal(text,err)}
}
