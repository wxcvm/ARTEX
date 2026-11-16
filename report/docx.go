package report

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"strings"
)

const DOCXContentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

const wordNS = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"

// DOCX packages the existing task report as an editable Word document.
// The report text remains the source of truth: no model calls or new dependencies.
func DOCX(markdown string) ([]byte, error) {
	var body bytes.Buffer
	body.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><w:document xmlns:w=\"" + wordNS + "\"><w:body>")
	inCode := false
	for _, line := range strings.Split(strings.ReplaceAll(markdown, "\r\n", "\n"), "\n") {
		style := "Normal"
		if strings.HasPrefix(strings.TrimSpace(line), "~~~") || strings.HasPrefix(strings.TrimSpace(line), strings.Repeat(string(rune(96)), 3)) {
			inCode = !inCode
			style = "Code"
		} else if inCode {
			style = "Code"
		} else {
			switch {
			case strings.HasPrefix(line, "# "):
				style, line = "Title", strings.TrimPrefix(line, "# ")
			case strings.HasPrefix(line, "## "):
				style, line = "Heading1", strings.TrimPrefix(line, "## ")
			case strings.HasPrefix(line, "### "):
				style, line = "Heading2", strings.TrimPrefix(line, "### ")
			}
		}
		body.WriteString("<w:p><w:pPr><w:pStyle w:val=\"" + style + "\"/></w:pPr><w:r><w:t xml:space=\"preserve\">")
		if err := xml.EscapeText(&body, []byte(line)); err != nil {
			return nil, err
		}
		body.WriteString("</w:t></w:r></w:p>")
	}
	body.WriteString("<w:sectPr><w:pgSz w:w=\"11906\" w:h=\"16838\"/><w:pgMar w:top=\"1440\" w:right=\"1440\" w:bottom=\"1440\" w:left=\"1440\" w:header=\"720\" w:footer=\"720\" w:gutter=\"0\"/></w:sectPr></w:body></w:document>")

	var out bytes.Buffer
	z := zip.NewWriter(&out)
	parts := []struct{ name, data string }{
		{"[Content_Types].xml", contentTypes},
		{"_rels/.rels", packageRels},
		{"word/document.xml", body.String()},
		{"word/_rels/document.xml.rels", documentRels},
		{"word/styles.xml", wordStyles},
	}
	for _, part := range parts {
		w, err := z.Create(part.name)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write([]byte(part.data)); err != nil {
			return nil, err
		}
	}
	if err := z.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

const contentTypes = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" +
	"<Types xmlns=\"http://schemas.openxmlformats.org/package/2006/content-types\">" +
	"<Default Extension=\"rels\" ContentType=\"application/vnd.openxmlformats-package.relationships+xml\"/>" +
	"<Default Extension=\"xml\" ContentType=\"application/xml\"/>" +
	"<Override PartName=\"/word/document.xml\" ContentType=\"application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml\"/>" +
	"<Override PartName=\"/word/styles.xml\" ContentType=\"application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml\"/></Types>"

const packageRels = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" +
	"<Relationships xmlns=\"http://schemas.openxmlformats.org/package/2006/relationships\">" +
	"<Relationship Id=\"rId1\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument\" Target=\"word/document.xml\"/></Relationships>"

const documentRels = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" +
	"<Relationships xmlns=\"http://schemas.openxmlformats.org/package/2006/relationships\">" +
	"<Relationship Id=\"rId1\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles\" Target=\"styles.xml\"/></Relationships>"

const wordStyles = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><w:styles xmlns:w=\"" + wordNS + "\">" +
	"<w:style w:type=\"paragraph\" w:default=\"1\" w:styleId=\"Normal\"><w:name w:val=\"Normal\"/><w:pPr><w:spacing w:after=\"120\"/></w:pPr><w:rPr><w:rFonts w:ascii=\"Calibri\" w:hAnsi=\"Calibri\" w:eastAsia=\"Microsoft YaHei\"/><w:sz w:val=\"22\"/></w:rPr></w:style>" +
	"<w:style w:type=\"paragraph\" w:styleId=\"Title\"><w:name w:val=\"Title\"/><w:basedOn w:val=\"Normal\"/><w:pPr><w:spacing w:after=\"240\"/></w:pPr><w:rPr><w:b/><w:sz w:val=\"40\"/></w:rPr></w:style>" +
	"<w:style w:type=\"paragraph\" w:styleId=\"Heading1\"><w:name w:val=\"heading 1\"/><w:basedOn w:val=\"Normal\"/><w:pPr><w:keepNext/><w:spacing w:before=\"240\" w:after=\"120\"/><w:outlineLvl w:val=\"0\"/></w:pPr><w:rPr><w:b/><w:sz w:val=\"30\"/></w:rPr></w:style>" +
	"<w:style w:type=\"paragraph\" w:styleId=\"Heading2\"><w:name w:val=\"heading 2\"/><w:basedOn w:val=\"Normal\"/><w:pPr><w:keepNext/><w:spacing w:before=\"180\" w:after=\"80\"/><w:outlineLvl w:val=\"1\"/></w:pPr><w:rPr><w:b/><w:sz w:val=\"26\"/></w:rPr></w:style>" +
	"<w:style w:type=\"paragraph\" w:styleId=\"Code\"><w:name w:val=\"Code\"/><w:basedOn w:val=\"Normal\"/><w:pPr><w:spacing w:after=\"0\"/></w:pPr><w:rPr><w:rFonts w:ascii=\"Consolas\" w:hAnsi=\"Consolas\" w:eastAsia=\"Microsoft YaHei\"/><w:sz w:val=\"20\"/></w:rPr></w:style></w:styles>"
