package kb

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// DocumentTextLimit bounds in-memory text extraction (including decompressed
// Word XML). It is not a limit on the number of documents in a library.
const DocumentTextLimit = 4 << 20

// ExtractDocumentText reads local documents without executing them or contacting
// any service. Only the main Word document part is read; macros are never run.
func ExtractDocumentText(filename string, raw []byte) (string,error) {
	var text string
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".txt",".md",".csv",".json":
		if len(raw)>DocumentTextLimit||!utf8.Valid(raw){return "",errors.New("text must be valid UTF-8 and at most 4 MiB")}
		text=strings.TrimPrefix(string(raw),"\uFEFF")
	case ".docx":
		z,err:=zip.NewReader(bytes.NewReader(raw),int64(len(raw)))
		if err!=nil{return "",err}
		found:=false
		for _,f:=range z.File {
			if f.Name!="word/document.xml"{continue}
			if found{return "",errors.New("duplicate Word document part")}
			found=true
			if f.UncompressedSize64>DocumentTextLimit{return "",errors.New("Word document text XML exceeds 4 MiB")}
			r,err:=f.Open()
			if err!=nil{return "",err}
			part,err:=io.ReadAll(io.LimitReader(r,DocumentTextLimit+1))
			closeErr:=r.Close()
			if err!=nil||closeErr!=nil{return "",errors.Join(err,closeErr)}
			if len(part)>DocumentTextLimit{return "",errors.New("Word document text XML exceeds 4 MiB")}
			dec:=xml.NewDecoder(bytes.NewReader(part))
			var b strings.Builder
			inText:=false
			for {
				tok,err:=dec.Token()
				if err==io.EOF{break}
				if err!=nil{return "",err}
				switch t:=tok.(type){
				case xml.StartElement:
					if t.Name.Space=="http://schemas.openxmlformats.org/wordprocessingml/2006/main" {
						if t.Name.Local=="t"{inText=true}
						if t.Name.Local=="tab"{b.WriteByte('\t')}
						if t.Name.Local=="br"{b.WriteByte('\n')}
					}
				case xml.CharData:
					if inText{b.Write([]byte(t))}
				case xml.EndElement:
					if t.Name.Space=="http://schemas.openxmlformats.org/wordprocessingml/2006/main"{
						if t.Name.Local=="t"{inText=false}
						if t.Name.Local=="p"{b.WriteByte('\n')}
					}
				}
			}
			text=b.String()
		}
		if !found{return "",errors.New("Word document has no main text part")}
	default:
		return "",errors.New("supported formats: TXT, Markdown, CSV, JSON, DOCX")
	}
	if strings.TrimSpace(text)==""||strings.ContainsRune(text,0){return "",errors.New("document is empty or contains NUL")}
	if len(text)>DocumentTextLimit{return "",errors.New("document text exceeds 4 MiB")}
	return text,nil
}
