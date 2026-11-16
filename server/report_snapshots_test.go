package server

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestTaskReportSnapshotArchiveRoundTrip(t *testing.T){
	s,f,req:=trafficEvidenceServer(t)
	s.ctx=context.Background()
	finding,err:=s.m.pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	tid:=*finding.TaskID
	task:=strconv.FormatInt(tid,10)
	base:="/api/report/snapshots"
	if w:=req("POST",base+"?format=md","{}");w.Code!=400{t.Fatal(w.Code,w.Body)}
	if w:=req("POST",base+"?task=active&format=md","{}");w.Code!=400{t.Fatal(w.Code,w.Body)}
	w:=req("POST","/api/report/templates",`{"name":"Snapshot fixture","brand":"Snapshot & Co","footer":"Original footer"}`)
	if w.Code!=200{t.Fatal(w.Code,w.Body)}
	var tmpl db.ReportTemplate
	if err:=json.Unmarshal(w.Body.Bytes(),&tmpl);err!=nil{t.Fatal(err)}
	defer s.m.pg.DeleteReportTemplate(t.Context(),tmpl.ID)
	w=req("POST",base+"?task="+task+"&format=docx&scope=all&layout=professional&template="+tmpl.ID+"&name=Version+one","{}")
	if w.Code!=201{t.Fatal(w.Code,w.Body)}
	var snap db.TaskReportSnapshot
	if err:=json.Unmarshal(w.Body.Bytes(),&snap);err!=nil||snap.ID==""||snap.TemplateName!="Snapshot fixture"{t.Fatal(snap,err)}
	path:=base+"/"+snap.ID+"?task="+task
	w=req("GET",path,"")
	if w.Code!=200||int64(w.Body.Len())!=snap.Size{t.Fatal(w.Code,w.Body.Len(),snap.Size)}
	original:=append([]byte(nil),w.Body.Bytes()...)
	if w=req("PUT","/api/report/templates/"+tmpl.ID,`{"name":"Changed","brand":"New brand"}`);w.Code!=200{t.Fatal(w.Code,w.Body)}
	if w=req("GET",path,"");w.Code!=200||!bytes.Equal(w.Body.Bytes(),original){t.Fatal("snapshot changed",w.Code)}
	other,err:=s.m.CreateTask("snapshot other","fixture",nil,0,0)
	if err!=nil{t.Fatal(err)}
	otherID,_:=strconv.ParseInt(other.ID,10,64)
	defer s.m.pg.DeleteTask(otherID)
	if w=req("GET",base+"/"+snap.ID+"?task="+other.ID,"");w.Code!=404{t.Fatal("cross task read",w.Code)}
	if w=req("DELETE",base+"/"+snap.ID+"?task="+other.ID,"");w.Code!=404{t.Fatal("cross task delete",w.Code)}
	if w=req("GET",base+"?task="+task,"");w.Code!=200||!strings.Contains(w.Body.String(),snap.ID){t.Fatal(w.Code,w.Body)}
	if w=req("GET","/api/report/library?q=Version+one&format=docx","");w.Code!=200||!strings.Contains(w.Body.String(),snap.ID)||!strings.Contains(w.Body.String(),`"task_id":"`+task+`"`){t.Fatal(w.Code,w.Body)}
	if w=req("GET","/api/report/library?q=Version+one&format=md","");w.Code!=200||strings.Contains(w.Body.String(),snap.ID){t.Fatal("format filter failed",w.Code,w.Body)}
	if err=s.m.SetTaskPaused(task,true);err!=nil{t.Fatal(err)}
	job,err:=s.m.pg.QueueTaskArchive(tid)
	if err!=nil{t.Fatal(err)}
	if _,err=s.m.pg.SaveTaskReportSnapshot(s.ctx,tid,snap,original);err==nil{t.Fatal("archive allowed late snapshot")}
	if err=s.runOneTaskArchiveJob();err!=nil{t.Fatal(err)}
	if _,_,err=s.m.pg.ReadTaskReportSnapshot(s.ctx,tid,snap.ID);err==nil{t.Fatal("snapshot retained hot")}
	if _,err=s.m.pg.QueueTaskArchiveRestore(job.ID);err!=nil{t.Fatal(err)}
	if err=s.runOneTaskArchiveJob();err!=nil{t.Fatal(err)}
	format,restored,err:=s.m.pg.ReadTaskReportSnapshot(s.ctx,tid,snap.ID)
	if err!=nil||format!="docx"||!bytes.Equal(restored,original){t.Fatal("snapshot restore",format,len(restored),err)}
	if err=s.m.pg.DeleteTaskReportSnapshot(s.ctx,tid,snap.ID);err!=nil{t.Fatal(err)}
	w=req("POST",base+"?task="+task+"&format=md&scope=confirmed&name=Markdown","{}")
	if w.Code!=201{t.Fatal(w.Code,w.Body)}
	var md db.TaskReportSnapshot
	if err:=json.Unmarshal(w.Body.Bytes(),&md);err!=nil||md.Format!="md"||md.Scope!="confirmed"{t.Fatal(md,err)}
	if w=req("GET",base+"/"+md.ID+"?task="+task,"");w.Code!=200||!strings.HasPrefix(w.Header().Get("Content-Type"),"text/markdown"){t.Fatal(w.Code,w.Header())}
	if w=req("DELETE",base+"/"+md.ID+"?task="+task,"");w.Code!=200{t.Fatal(w.Code,w.Body)}
}
