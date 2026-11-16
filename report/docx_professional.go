package report

import (
	"bytes"
	"encoding/xml"
	"strings"

	"github.com/Autumn-27/artex/db"
)

const professionalContentTypes = `<Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/><Override PartName="/word/settings.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.settings+xml"/>`
const professionalRelationships = `<Relationship Id="rIdFooter" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer" Target="footer1.xml"/><Relationship Id="rIdSettings" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings" Target="settings.xml"/>`
const professionalSettings = `<?xml version="1.0" encoding="UTF-8"?><w:settings xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:updateFields w:val="true"/></w:settings>`
const professionalFooter = `<?xml version="1.0" encoding="UTF-8"?><w:ftr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:t xml:space="preserve">ARTEX  |  </w:t></w:r><w:fldSimple w:instr=" PAGE "><w:r><w:t>1</w:t></w:r></w:fldSimple></w:p></w:ftr>`
const professionalStyles = `<w:style w:type="paragraph" w:styleId="CoverKicker"><w:name w:val="Cover kicker"/><w:basedOn w:val="Normal"/><w:pPr><w:spacing w:before="3000" w:after="240"/><w:jc w:val="center"/></w:pPr><w:rPr><w:b/><w:color w:val="28517B"/><w:sz w:val="28"/></w:rPr></w:style><w:style w:type="paragraph" w:styleId="CoverTitle"><w:name w:val="Cover title"/><w:basedOn w:val="Normal"/><w:pPr><w:spacing w:after="480"/><w:jc w:val="center"/></w:pPr><w:rPr><w:b/><w:color w:val="17324D"/><w:sz w:val="48"/></w:rPr></w:style><w:style w:type="paragraph" w:styleId="CoverNote"><w:name w:val="Cover note"/><w:basedOn w:val="Normal"/><w:pPr><w:jc w:val="center"/></w:pPr><w:rPr><w:color w:val="64748B"/><w:sz w:val="22"/></w:rPr></w:style>`

func writeProfessionalFrontMatter(body *bytes.Buffer, markdown, locale string, template *db.ReportTemplate) error {
	title := reportText(locale, "安全测试报告", "Security assessment report")
	for _, line := range strings.Split(markdown, "\n") {
		if strings.HasPrefix(line, "# ") {
			title = nz(strings.TrimSpace(strings.TrimPrefix(line, "# ")), title)
			break
		}
	}
	paragraph := func(style, text string) error {
		body.WriteString("<w:p><w:pPr><w:pStyle w:val=\"" + style + "\"/></w:pPr><w:r><w:t xml:space=\"preserve\">")
		if err := xml.EscapeText(body, []byte(text)); err != nil {
			return err
		}
		body.WriteString("</w:t></w:r></w:p>")
		return nil
	}
	brand, subtitle := "ARTEX", reportText(locale, "安全评估 · 证据与发现", "Security assessment · evidence and findings")
	if template != nil {
		brand = nz(template.Brand, brand)
		title = nz(template.CoverTitle, title)
		subtitle = nz(template.Subtitle, subtitle)
	}
	if err := paragraph("CoverKicker", brand); err != nil {
		return err
	}
	if err := paragraph("CoverTitle", title); err != nil {
		return err
	}
	if err := paragraph("CoverNote", subtitle); err != nil {
		return err
	}
	if template != nil {
		if m:=template.Metadata;m!=nil{
			for _,field:=range []struct{label,value string}{{reportMetadataText(locale,"项目","專案","Project"),m.Project},{reportMetadataText(locale,"版本","版本","Version"),m.Version}}{if field.value!=""{if err:=paragraph("CoverNote",field.label+": "+field.value);err!=nil{return err}}}
		}
		for _, field := range []struct{label,value string}{
			{reportText(locale,"委托方","Prepared for"),template.PreparedFor},
			{reportText(locale,"报告编制","Prepared by"),template.Author},
			{reportText(locale,"保密级别","Confidentiality"),template.Confidentiality},
		}{
			if field.value != "" { if err := paragraph("CoverNote",field.label+": "+field.value); err != nil {return err} }
		}
	}
	pageBreak := `<w:p><w:r><w:br w:type="page"/></w:r></w:p>`
	body.WriteString(pageBreak)
	if err := paragraph("Title", reportText(locale, "目录", "Contents")); err != nil {
		return err
	}
	body.WriteString(`<w:p><w:r><w:fldChar w:fldCharType="begin" w:dirty="true"/></w:r><w:r><w:instrText xml:space="preserve"> TOC \o "1-3" \h \z \u </w:instrText></w:r><w:r><w:fldChar w:fldCharType="separate"/></w:r><w:r><w:t>`)
	if err := xml.EscapeText(body, []byte(reportText(locale, "在 Word 中更新目录以显示页码。", "Update the table of contents in Word to display page numbers."))); err != nil {
		return err
	}
	body.WriteString(`</w:t></w:r><w:r><w:fldChar w:fldCharType="end"/></w:r></w:p>`)
	body.WriteString(pageBreak)
	return nil
}

func templateFooter(t *db.ReportTemplate) string {
	if t == nil || t.Footer == "" {return professionalFooter}
	var escaped bytes.Buffer
	_ = xml.EscapeText(&escaped, []byte(t.Footer))
	return strings.Replace(professionalFooter,"ARTEX  |  ",escaped.String()+"  |  ",1)
}

// Custom sections are literal user-authored paragraphs, not generated findings.
// The normal report body and evidence appendices are always retained.
func writeTemplateSections(body *bytes.Buffer,t *db.ReportTemplate,position string)error{
	if t==nil{return nil}
	paragraph:=func(style,text string)error{
		body.WriteString("<w:p><w:pPr><w:pStyle w:val=\""+style+"\"/></w:pPr><w:r><w:t xml:space=\"preserve\">")
		if err:=xml.EscapeText(body,[]byte(text));err!=nil{return err}
		body.WriteString("</w:t></w:r></w:p>")
		return nil
	}
	for _,section:=range t.Sections{
		if section.Position!=position{continue}
		if err:=paragraph("Heading1",section.Title);err!=nil{return err}
		for _,line:=range strings.Split(strings.ReplaceAll(section.Body,"\r\n","\n"),"\n"){
			if err:=paragraph("Normal",line);err!=nil{return err}
		}
	}
	return nil
}

// Personnel and standards are explicitly supplied by the report author.
// They are not an assertion of certification, test completion or compliance.
func writeTemplateMetadata(body *bytes.Buffer,t *db.ReportTemplate,locale string)error{
	if t==nil||t.Metadata==nil{return nil}
	m:=t.Metadata
	paragraph:=func(style,text string)error{
		body.WriteString("<w:p><w:pPr><w:pStyle w:val=\""+style+"\"/></w:pPr><w:r><w:t xml:space=\"preserve\">")
		if err:=xml.EscapeText(body,[]byte(text));err!=nil{return err}
		body.WriteString("</w:t></w:r></w:p>")
		return nil
	}
	if len(m.Personnel)>0{
		if err:=paragraph("Heading1",reportMetadataText(locale,"项目人员","專案人員","Project personnel"));err!=nil{return err}
		for _,p:=range m.Personnel{
			line:=p.Name
			if p.Role!=""{line+=" — "+p.Role}
			if p.Contact!=""{line+=" ("+p.Contact+")"}
			if err:=paragraph("Normal",line);err!=nil{return err}
		}
	}
	hasStandards:=false
	for _,s:=range m.Standards{
		if s==""{continue}
		if !hasStandards{
			if err:=paragraph("Heading1",reportMetadataText(locale,"参考标准","參考標準","Reference standards"));err!=nil{return err}
			if err:=paragraph("Normal",reportMetadataText(locale,"以下为报告作者指定的参考标准，不代表合规认证。","以下為報告作者指定的參考標準，不代表合規認證。","References specified by the report author; not a certification of compliance."));err!=nil{return err}
			hasStandards=true
		}
		if err:=paragraph("Normal",s);err!=nil{return err}
	}
	return nil
}

func reportMetadataText(locale,zh,tw,en string)string{if locale=="zh-TW"{return tw};return reportText(locale,zh,en)}
