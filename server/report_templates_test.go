package server

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"image"
	"image/png"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestReusableReportTemplateExports(t *testing.T){
	s,f,req:=trafficEvidenceServer(t)
	finding,err:=s.m.pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing fixture task",err)}
	base:="/api/report/templates"
	if w:=req("POST",base,`{"name":""}`);w.Code!=400{t.Fatal(w.Code,w.Body)}
	w:=req("POST",base,`{"name":"Client fixture","brand":"Acme & Labs","cover_title":"Custom <report>","subtitle":"Assessment","prepared_for":"Client A","author":"Analyst","confidentiality":"Internal only","footer":"Acme & Co"}`)
	if w.Code!=200{t.Fatal(w.Code,w.Body)}
	var tmpl db.ReportTemplate
	if err:=json.Unmarshal(w.Body.Bytes(),&tmpl);err!=nil||tmpl.ID==""{t.Fatal(tmpl,err)}
	defer s.m.pg.DeleteReportTemplate(t.Context(),tmpl.ID)
	if w=req("GET",base,"");w.Code!=200||!strings.Contains(w.Body.String(),tmpl.ID){t.Fatal(w.Code,w.Body)}
	path:=base+"/"+tmpl.ID
	if w=req("PUT",path,`{"name":"Updated template","brand":"Acme & Labs","cover_title":"Custom <report>","subtitle":"Assessment","prepared_for":"Client A","author":"Analyst","confidentiality":"Internal only","footer":"Acme & Co"}`);w.Code!=200||!strings.Contains(w.Body.String(),"Updated template"){t.Fatal(w.Code,w.Body)}
	for _,export:=range []string{
		fmt.Sprintf("/api/report?task=%d&scope=all&format=docx&layout=professional",*finding.TaskID),
		fmt.Sprintf("/api/exploration/findings/export?scope=selected&ids=%d&format=docx&layout=professional",f.FindingID),
	}{
		w=req("GET",export+"&template="+tmpl.ID,"")
		if w.Code!=200{t.Fatal(export,w.Code,w.Body)}
		z,err:=zip.NewReader(bytes.NewReader(w.Body.Bytes()),int64(w.Body.Len()))
		if err!=nil{t.Fatal(err)}
		parts:=map[string]string{}
		for _,entry:=range z.File{
			if entry.Name!="word/document.xml"&&entry.Name!="word/footer1.xml"{continue}
			reader,err:=entry.Open()
			if err!=nil{t.Fatal(err)}
			raw,err:=io.ReadAll(reader)
			reader.Close()
			if err!=nil{t.Fatal(err)}
			dec:=xml.NewDecoder(bytes.NewReader(raw))
			for{if _,err:=dec.Token();err==io.EOF{break}else if err!=nil{t.Fatal(entry.Name,err)}}
			parts[entry.Name]=string(raw)
		}
		for _,want:=range []string{"Acme &amp; Labs","Custom &lt;report&gt;","Assessment","Client A","Analyst","Internal only","TOC"}{
			if !strings.Contains(parts["word/document.xml"],want){t.Fatal("missing",want,export)}
		}
		if !strings.Contains(parts["word/footer1.xml"],"Acme &amp; Co")||!strings.Contains(parts["word/footer1.xml"],"PAGE"){t.Fatal("footer missing")}
	}
	if w=req("DELETE",path,"");w.Code!=200{t.Fatal(w.Code,w.Body)}
	if w=req("GET",fmt.Sprintf("/api/report?task=%d&format=docx&template=%s",*finding.TaskID,tmpl.ID),"");w.Code!=404{t.Fatal("deleted template silently accepted",w.Code)}
}

func TestReportTemplateLogoLifecycle(t *testing.T){
	s,f,req:=trafficEvidenceServer(t)
	finding,err:=s.m.pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing fixture task",err)}
	w:=req("POST","/api/report/templates",`{"name":"Logo fixture"}`)
	var tmpl db.ReportTemplate
	if err:=json.Unmarshal(w.Body.Bytes(),&tmpl);w.Code!=200||err!=nil{t.Fatal(w.Code,w.Body,err)}
	defer s.m.pg.DeleteReportTemplate(t.Context(),tmpl.ID)
	path:="/api/report/templates/"+tmpl.ID
	if w=req("PUT",path+"/logo","not an image");w.Code!=400{t.Fatal(w.Code,w.Body)}
	var img bytes.Buffer
	if err:=png.Encode(&img,image.NewRGBA(image.Rect(0,0,24,12)));err!=nil{t.Fatal(err)}
	if w=req("PUT",path+"/logo",img.String());w.Code!=200||!strings.Contains(w.Body.String(),`"has_logo":true`){t.Fatal(w.Code,w.Body)}
	if w=req("PUT",path,`{"name":"Edited logo fixture"}`);w.Code!=200||!strings.Contains(w.Body.String(),`"has_logo":true`){t.Fatal("text edit removed logo",w.Code,w.Body)}
	export:=fmt.Sprintf("/api/report?task=%d&format=docx&template=%s",*finding.TaskID,tmpl.ID)
	w=req("GET",export,"")
	if w.Code!=200{t.Fatal(w.Code,w.Body)}
	z,err:=zip.NewReader(bytes.NewReader(w.Body.Bytes()),int64(w.Body.Len()))
	if err!=nil{t.Fatal(err)}
	parts:=map[string][]byte{}
	for _,entry:=range z.File{
		r,err:=entry.Open();if err!=nil{t.Fatal(err)}
		raw,err:=io.ReadAll(r);r.Close();if err!=nil{t.Fatal(err)}
		parts[entry.Name]=raw
	}
	if !bytes.Equal(parts["word/media/report-logo.png"],img.Bytes())||!bytes.Contains(parts["word/document.xml"],[]byte(`r:embed="rIdLogo"`))||!bytes.Contains(parts["word/_rels/document.xml.rels"],[]byte("media/report-logo.png")){t.Fatal("logo not embedded")}
	if w=req("DELETE",path+"/logo","");w.Code!=200||!strings.Contains(w.Body.String(),`"has_logo":false`){t.Fatal(w.Code,w.Body)}
	w=req("GET",export,"")
	if w.Code!=200{t.Fatal(w.Code,w.Body)}
	z,err=zip.NewReader(bytes.NewReader(w.Body.Bytes()),int64(w.Body.Len()))
	if err!=nil{t.Fatal(err)}
	for _,entry:=range z.File{if strings.Contains(entry.Name,"report-logo"){t.Fatal("removed logo still exported")}}
}

func TestReportTemplateCustomSections(t *testing.T){
	s,f,req:=trafficEvidenceServer(t)
	finding,err:=s.m.pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	base:="/api/report/templates"
	if w:=req("POST",base,`{"name":"bad","sections":[{"title":"bad","body":"text","position":"middle"}]}`);w.Code!=400{t.Fatal(w.Code,w.Body)}
	w:=req("POST",base,`{"name":"Section fixture","sections":[{"title":"Fixture methodology <&>","body":"Method paragraph 1\nMethod paragraph 2","position":"before"},{"title":"Fixture scope","body":"Scope paragraph","position":"before"},{"title":"Fixture conclusion","body":"Conclusion paragraph","position":"after"}]}`)
	var tmpl db.ReportTemplate
	if err:=json.Unmarshal(w.Body.Bytes(),&tmpl);w.Code!=200||err!=nil||len(tmpl.Sections)!=3{t.Fatal(w.Code,w.Body,err)}
	defer s.m.pg.DeleteReportTemplate(t.Context(),tmpl.ID)
	path:=base+"/"+tmpl.ID
	w=req("PUT",path,`{"name":"Edited section fixture"}`)
	if err:=json.Unmarshal(w.Body.Bytes(),&tmpl);w.Code!=200||err!=nil||len(tmpl.Sections)!=3{t.Fatal("omitted sections lost",w.Code,w.Body,err)}
	export:=fmt.Sprintf("/api/report?task=%d&format=docx&template=%s",*finding.TaskID,tmpl.ID)
	w=req("GET",export,"")
	if w.Code!=200{t.Fatal(w.Code,w.Body)}
	z,err:=zip.NewReader(bytes.NewReader(w.Body.Bytes()),int64(w.Body.Len()))
	if err!=nil{t.Fatal(err)}
	var document string
	for _,entry:=range z.File{if entry.Name=="word/document.xml"{r,err:=entry.Open();if err!=nil{t.Fatal(err)};raw,err:=io.ReadAll(r);r.Close();if err!=nil{t.Fatal(err)};document=string(raw);dec:=xml.NewDecoder(bytes.NewReader(raw));for{if _,err:=dec.Token();err==io.EOF{break}else if err!=nil{t.Fatal(err)}}}}
	prev:=-1
	for _,want:=range []string{"Fixture methodology &lt;&amp;&gt;","Method paragraph 1","Method paragraph 2","Fixture scope","Scope paragraph","任务目标","Fixture conclusion","Conclusion paragraph"}{
		pos:=strings.Index(document,want);if pos<=prev{t.Fatalf("missing or out of order %q",want)};prev=pos
	}
	w=req("PUT",path,`{"name":"Clear sections","sections":[]}`)
	if err:=json.Unmarshal(w.Body.Bytes(),&tmpl);w.Code!=200||err!=nil||len(tmpl.Sections)!=0{t.Fatal(w.Code,w.Body,err)}
}
