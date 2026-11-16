package server

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
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
