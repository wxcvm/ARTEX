package server

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestKnowledgeDocumentScopeAndCRUD(t *testing.T) {
	s,f,req:=trafficEvidenceServer(t)
	finding,err:=s.m.pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	taskID:=strconv.FormatInt(*finding.TaskID,10)
	other,err:=s.m.CreateTask("knowledge isolation fixture","local",nil,0,0)
	if err!=nil{t.Fatal(err)}
	otherID,_:=strconv.ParseInt(other.ID,10,64)
	defer s.m.pg.DeleteTask(otherID)
	const base="/api/knowledge/documents"
	create:=func(path,title,content string)db.KnowledgeDocument{
		t.Helper()
		raw,_:=json.Marshal(map[string]string{"title":title,"category":"fixture","content":content})
		w:=req("POST",path,string(raw))
		if w.Code!=200{t.Fatal(w.Code,w.Body)}
		var doc db.KnowledgeDocument
		if err:=json.Unmarshal(w.Body.Bytes(),&doc);err!=nil||doc.ID==""{t.Fatal(doc,err)}
		return doc
	}
	global:=create(base,"Global fixture","100% unicode 中文 unique-kn-test")
	defer s.m.pg.DeleteKnowledgeDocument(context.Background(),0,global.ID)
	local:=create(base+"?task="+taskID,"Task fixture","100% unique-kn-test local")
	if local.TaskID!=taskID||global.TaskID!=""{t.Fatal("scope",local,global)}
	if w:=req("GET",base+"/"+local.ID,"");w.Code!=404{t.Fatal("global leaked task doc",w.Code,w.Body)}
	if w:=req("GET",base+"/"+local.ID+"?task="+other.ID,"");w.Code!=404{t.Fatal("cross task leaked",w.Code,w.Body)}
	if w:=req("GET",base+"/"+global.ID+"?task="+taskID,"");w.Code!=404{t.Fatal("unexpected global visibility",w.Code)}
	if w:=req("GET",base+"/"+global.ID+"?task="+taskID+"&include_global=true","");w.Code!=200{t.Fatal(w.Code,w.Body)}
	if w:=req("PUT",base+"/"+global.ID+"?task="+taskID+"&include_global=true",`{"title":"wrong","content":"wrong"}`);w.Code!=404{t.Fatal("edited global from task",w.Code,w.Body)}
	if w:=req("DELETE",base+"/"+local.ID+"?task="+other.ID,"");w.Code!=404{t.Fatal("cross-task deletion",w.Code)}
	for _,tc:=range []struct{path string;want int}{
		{base+"?q=unique-kn-test",1},
		{base+"?task="+taskID+"&q=unique-kn-test",1},
		{base+"?task="+taskID+"&include_global=true&q=unique-kn-test",2},
		{base+"?task="+other.ID+"&q=unique-kn-test",0},
		{base+"?task="+taskID+"&q=%25&category=fixture",1},
		{base+"?task="+taskID+"&q=%25&category=missing",0},
	}{
		w:=req("GET",tc.path,"")
		var page db.KnowledgeDocumentPage
		if err:=json.Unmarshal(w.Body.Bytes(),&page);w.Code!=200||err!=nil||len(page.Items)!=tc.want{t.Fatal(tc.path,w.Code,page,err)}
	}
	if w:=req("PUT",base+"/"+local.ID+"?task="+taskID,`{"title":"Updated","content":"updated text","category":"changed"}`);w.Code!=200{t.Fatal(w.Code,w.Body)}
	if w:=req("GET",base+"/"+local.ID+"?task="+taskID,"");w.Code!=200||!strings.Contains(w.Body.String(),"updated text"){t.Fatal(w.Code,w.Body)}
	if w:=req("POST",base+"?task=invalid",`{"title":"x","content":"x"}`);w.Code!=400{t.Fatal("invalid task became global",w.Code)}
	if w:=req("DELETE",base+"/"+local.ID+"?task="+taskID,"");w.Code!=200{t.Fatal(w.Code,w.Body)}
	if w:=req("GET",base+"/"+local.ID+"?task="+taskID,"");w.Code!=404{t.Fatal(w.Code,w.Body)}

	var body bytes.Buffer
	form:=multipart.NewWriter(&body)
	file,err:=form.CreateFormFile("file","fixture.md")
	if err!=nil{t.Fatal(err)}
	if _,err=file.Write([]byte("uploaded 中文 text"));err!=nil{t.Fatal(err)}
	if err=form.WriteField("category","uploaded");err!=nil{t.Fatal(err)}
	if err=form.Close();err!=nil{t.Fatal(err)}
	token,err:=signJWT(s.jwtKey)
	if err!=nil{t.Fatal(err)}
	r:=httptest.NewRequest("POST",base+"/upload?task="+taskID,&body)
	r.Header.Set("Authorization","Bearer "+token)
	r.Header.Set("Content-Type",form.FormDataContentType())
	w:=httptest.NewRecorder()
	s.Handler().ServeHTTP(w,r)
	var uploaded db.KnowledgeDocument
	if err:=json.Unmarshal(w.Body.Bytes(),&uploaded);w.Code!=201||err!=nil||uploaded.SourceFilename!="fixture.md"||uploaded.Content!="uploaded 中文 text"||uploaded.TaskID!=taskID{t.Fatal(w.Code,uploaded,err)}
}

func TestKnowledgeDocumentArchiveRoundTrip(t *testing.T) {
	s,f,_:=trafficEvidenceServer(t)
	s.ctx=context.Background()
	finding,err:=s.m.pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal(err)}
	tid:=*finding.TaskID
	doc,err:=s.m.pg.SaveKnowledgeDocument(s.ctx,tid,"",db.KnowledgeDocument{Title:"archive knowledge",Category:"notes",Content:"retained 中文"})
	if err!=nil{t.Fatal(err)}
	if err=s.m.pg.RecordKnowledgeRetrieval(s.ctx,tid,"fixture","archived query",[]db.KnowledgeRetrievalHit{{ID:doc.ID,Scope:"task",Title:doc.Title}},false);err!=nil{t.Fatal(err)}
	if err=s.m.SetTaskPaused(strconv.FormatInt(tid,10),true);err!=nil{t.Fatal(err)}
	job,err:=s.m.pg.QueueTaskArchive(tid)
	if err!=nil{t.Fatal(err)}
	if _,err=s.m.pg.SaveKnowledgeDocument(s.ctx,tid,doc.ID,db.KnowledgeDocument{Title:"late",Content:"late"});err==nil{t.Fatal("archive allowed late edit")}
	if err=s.runOneTaskArchiveJob();err!=nil{t.Fatal(err)}
	if _,err=s.m.pg.GetKnowledgeDocument(s.ctx,tid,false,doc.ID);err==nil{t.Fatal("archived document retained hot")}
	if history,err:=s.m.pg.ListKnowledgeRetrievals(s.ctx,tid);err!=nil||len(history)!=0{t.Fatal("history retained hot",history,err)}
	if _,err=s.m.pg.QueueTaskArchiveRestore(job.ID);err!=nil{t.Fatal(err)}
	if err=s.runOneTaskArchiveJob();err!=nil{t.Fatal(err)}
	got,err:=s.m.pg.GetKnowledgeDocument(s.ctx,tid,false,doc.ID)
	if err!=nil||got.Content!=doc.Content||got.Category!=doc.Category||got.Title!=doc.Title{t.Fatal(got,err)}
	history,err:=s.m.pg.ListKnowledgeRetrievals(s.ctx,tid)
	if err!=nil||len(history)!=1||history[0].Query!="archived query"||len(history[0].Hits)!=1||history[0].Hits[0].ID!=doc.ID{t.Fatal("history not restored",history,err)}
}
