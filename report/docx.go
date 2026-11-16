package report

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"strings"

	"github.com/Autumn-27/artex/db"
)

const DOCXContentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

const wordNS = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"

// DOCX packages the existing task report as an editable Word document.
// The report text remains the source of truth: no model calls or new dependencies.
func DOCX(markdown string) ([]byte, error) {
	return DOCXWithScreenshots(markdown, nil, "en")
}

// DOCXWithScreenshots embeds local image bytes in a finding-labelled appendix.
// Report text cannot introduce an image reference or remote relationship.
func DOCXWithScreenshots(markdown string, findings []*db.DBFinding, locale string) ([]byte, error) {
	return DOCXWithLayout(markdown, findings, locale, "simple")
}

// DOCXWithLayout supports built-in simple and professional report layouts.
func DOCXWithLayout(markdown string, findings []*db.DBFinding, locale, layout string) ([]byte, error) {
	return DOCXWithTemplate(markdown, findings, locale, layout, nil)
}

// DOCXWithTemplate applies presentation settings without altering findings or evidence.
func DOCXWithTemplate(markdown string, findings []*db.DBFinding, locale, layout string, template *db.ReportTemplate) ([]byte, error) {
	if layout != "" && layout != "simple" && layout != "professional" {
		return nil, fmt.Errorf("unknown Word report layout: %s", layout)
	}
	professional := layout == "professional" || template != nil
	var body bytes.Buffer
	body.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><w:document xmlns:w=\"" + wordNS + "\" xmlns:r=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships\" xmlns:wp=\"http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing\" xmlns:a=\"http://schemas.openxmlformats.org/drawingml/2006/main\" xmlns:pic=\"http://schemas.openxmlformats.org/drawingml/2006/picture\"><w:body>")
	var logoMedia []docxMedia
	var logoRels string
	if professional && template != nil && len(template.Logo)>0 {
		cfg,format,err:=reportLogoConfig(template.Logo)
		if err!=nil{return nil,err}
		ext:="png";if format=="jpeg"{ext="jpg"}
		name:="report-logo."+ext
		cx,cy:=float64(cfg.Width)*9525,float64(cfg.Height)*9525
		scale:=math.Min(1,math.Min(2.5*914400/cx,1.0*914400/cy))
		writeDOCXImage(&body,int64(cx*scale),int64(cy*scale),0,name,"rIdLogo","Report logo")
		logoMedia=[]docxMedia{{"word/media/"+name,template.Logo}}
		logoRels="<Relationship Id=\"rIdLogo\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/image\" Target=\"media/"+name+"\"/>"
	}
	if professional {
		if err := writeProfessionalFrontMatter(&body, markdown, locale, template); err != nil {
			return nil, err
		}
		if err:=writeTemplateMetadata(&body,template,locale);err!=nil{return nil,err}
		if err:=writeTemplateSections(&body,template,"before");err!=nil{return nil,err}
	}
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
			case strings.HasPrefix(line, "#### "):
				style, line = "Heading3", strings.TrimPrefix(line, "#### ")
			}
		}
		body.WriteString("<w:p><w:pPr><w:pStyle w:val=\"" + style + "\"/></w:pPr><w:r><w:t xml:space=\"preserve\">")
		if err := xml.EscapeText(&body, []byte(line)); err != nil {
			return nil, err
		}
		body.WriteString("</w:t></w:r></w:p>")
	}
	if professional{if err:=writeTemplateSections(&body,template,"after");err!=nil{return nil,err}}
	media, relationships, err := appendDOCXScreenshots(&body, findings, locale)
	if err != nil {
		return nil, err
	}
	media=append(logoMedia,media...)
	relationships+=logoRels
	body.WriteString("<w:sectPr>")
	if professional {
		body.WriteString(`<w:footerReference w:type="default" r:id="rIdFooter"/>`)
	}
	body.WriteString("<w:pgSz w:w=\"11906\" w:h=\"16838\"/><w:pgMar w:top=\"1440\" w:right=\"1440\" w:bottom=\"1440\" w:left=\"1440\" w:header=\"720\" w:footer=\"720\" w:gutter=\"0\"/></w:sectPr></w:body></w:document>")

	var out bytes.Buffer
	z := zip.NewWriter(&out)
	types := strings.Replace(contentTypes, "</Types>", "<Default Extension=\"png\" ContentType=\"image/png\"/><Default Extension=\"jpg\" ContentType=\"image/jpeg\"/></Types>", 1)
	styles := wordStyles
	if professional {
		types = strings.Replace(types, "</Types>", professionalContentTypes+"</Types>", 1)
		relationships += professionalRelationships
		styles = strings.Replace(styles, "</w:styles>", professionalStyles+"</w:styles>", 1)
	}
	parts := []struct{ name, data string }{
		{"[Content_Types].xml", types},
		{"_rels/.rels", packageRels},
		{"word/document.xml", body.String()},
		{"word/_rels/document.xml.rels", strings.Replace(documentRels, "</Relationships>", relationships+"</Relationships>", 1)},
		{"word/styles.xml", styles},
	}
	if professional {
		parts = append(parts, struct{ name, data string }{"word/footer1.xml", templateFooter(template)}, struct{ name, data string }{"word/settings.xml", professionalSettings})
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
	for _, item := range media {
		w, err := z.Create(item.name)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(item.data); err != nil {
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
	"<w:style w:type=\"paragraph\" w:styleId=\"Heading3\"><w:name w:val=\"heading 3\"/><w:basedOn w:val=\"Normal\"/><w:pPr><w:keepNext/><w:spacing w:before=\"140\" w:after=\"80\"/><w:outlineLvl w:val=\"2\"/></w:pPr><w:rPr><w:b/><w:sz w:val=\"24\"/></w:rPr></w:style>" +
	"<w:style w:type=\"paragraph\" w:styleId=\"Code\"><w:name w:val=\"Code\"/><w:basedOn w:val=\"Normal\"/><w:pPr><w:spacing w:after=\"0\"/></w:pPr><w:rPr><w:rFonts w:ascii=\"Consolas\" w:hAnsi=\"Consolas\" w:eastAsia=\"Microsoft YaHei\"/><w:sz w:val=\"20\"/></w:rPr></w:style></w:styles>"

type docxMedia struct {
	name string
	data []byte
}

func appendDOCXScreenshots(body *bytes.Buffer, findings []*db.DBFinding, locale string) ([]docxMedia, string, error) {
	var media []docxMedia
	var rels strings.Builder
	items := append([]*db.DBFinding(nil), findings...)
	sortFindingsForExport(items)
	writeParagraph := func(style, text string) error {
		body.WriteString("<w:p><w:pPr><w:pStyle w:val=\"" + style + "\"/></w:pPr><w:r><w:t xml:space=\"preserve\">")
		if err := xml.EscapeText(body, []byte(text)); err != nil {
			return err
		}
		body.WriteString("</w:t></w:r></w:p>")
		return nil
	}
	for _, f := range items {
		if len(f.Screenshots) == 0 {
			continue
		}
		if len(media) == 0 {
			body.WriteString("<w:p><w:r><w:br w:type=\"page\"/></w:r></w:p>")
			if err := writeParagraph("Heading1", reportText(locale, "附录：截图证据", "Appendix: Screenshot evidence")); err != nil {
				return nil, "", err
			}
		}
		if err := writeParagraph("Heading2", fmt.Sprintf("%s #%d — %s", reportText(locale, "发现", "Finding"), f.ID, findingTitle(f))); err != nil {
			return nil, "", err
		}
		for i, shot := range f.Screenshots {
			cfg, format, err := image.DecodeConfig(bytes.NewReader(shot.Image))
			if err != nil || cfg.Width <= 0 || cfg.Height <= 0 || (format != "png" && format != "jpeg") {
				return nil, "", fmt.Errorf("invalid screenshot for finding %d", f.ID)
			}
			ext := "png"
			if format == "jpeg" {
				ext = "jpg"
			}
			if shot.MIME != "image/"+format {
				return nil, "", fmt.Errorf("screenshot format mismatch for finding %d", f.ID)
			}
			// Use native 96 dpi and scale down to fit the A4 printable area.
			cx, cy := float64(cfg.Width)*9525, float64(cfg.Height)*9525
			scale := math.Min(1, math.Min(5.9*914400/cx, 7.5*914400/cy))
			x, y := int64(cx*scale), int64(cy*scale)
			n := len(media)+1
			name := fmt.Sprintf("screenshot-%d.%s", n, ext)
			rid := fmt.Sprintf("rId%d", n+1)
			media = append(media, docxMedia{"word/media/"+name, shot.Image})
			fmt.Fprintf(&rels, "<Relationship Id=\"%s\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/image\" Target=\"media/%s\"/>", rid, name)
			writeDOCXImage(body,x,y,n,name,rid,fmt.Sprintf("Screenshot %d",n))
			caption := fmt.Sprintf("%s %d: %s", reportText(locale, "图", "Figure"), i+1, nz(shot.Caption, reportText(locale, "无图注", "No caption")))
			if shot.BindingID != "" {
				caption += fmt.Sprintf(" (%s #%s)", reportText(locale, "关联证据", "linked evidence"), shot.BindingID)
			}
			if err := writeParagraph("Normal", caption); err != nil {
				return nil, "", err
			}
		}
	}
	return media, rels.String(), nil
}

func writeDOCXImage(body *bytes.Buffer,x,y int64,id int,name,rid,label string){
	fmt.Fprintf(body, "<w:p><w:pPr><w:jc w:val=\"center\"/></w:pPr><w:r><w:drawing><wp:inline distT=\"0\" distB=\"0\" distL=\"0\" distR=\"0\"><wp:extent cx=\"%d\" cy=\"%d\"/><wp:docPr id=\"%d\" name=\"%s\"/><a:graphic><a:graphicData uri=\"http://schemas.openxmlformats.org/drawingml/2006/picture\"><pic:pic><pic:nvPicPr><pic:cNvPr id=\"%d\" name=\"%s\"/><pic:cNvPicPr/></pic:nvPicPr><pic:blipFill><a:blip r:embed=\"%s\"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill><pic:spPr><a:xfrm><a:off x=\"0\" y=\"0\"/><a:ext cx=\"%d\" cy=\"%d\"/></a:xfrm><a:prstGeom prst=\"rect\"><a:avLst/></a:prstGeom></pic:spPr></pic:pic></a:graphicData></a:graphic></wp:inline></w:drawing></w:r></w:p>", x, y, id, label, id, name, rid, x, y)
}
