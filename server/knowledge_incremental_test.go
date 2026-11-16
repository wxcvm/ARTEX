package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/knowledgeindex"
)

func TestIncrementalKnowledgeIndexReuse(t *testing.T){
	s,f,req:=trafficEvidenceServer(t)
	pg:=s.m.pg
	ctx:=context.Background()
	finding,err:=pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	tid:=*finding.TaskID
	task:=strconv.FormatInt(tid,10)
	old,exists,err:=pg.GetSetting(knowledgeindex.SettingKey)
	if err!=nil{t.Fatal(err)}
	defer func(){if exists{_ = pg.SetSetting(knowledgeindex.SettingKey,old)}else{_,_=pg.Exec("DELETE FROM settings WHERE key=$1",knowledgeindex.SettingKey)}}()
	var inputs atomic.Int64
	mock:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		var in struct{Input []string `json:"input"`}
		if err:=json.NewDecoder(r.Body).Decode(&in);err!=nil{t.Error(err);w.WriteHeader(400);return}
		inputs.Add(int64(len(in.Input)))
		out:=[]any{}
		for i:=range in.Input{out=append(out,map[string]any{"index":i,"embedding":[]float32{1,0}})}
		_ = json.NewEncoder(w).Encode(map[string]any{"data":out})
	}))
	defer mock.Close()
	c:=knowledgeindex.Config{Enabled:true,BaseURL:mock.URL,Model:"incremental-fixture"}
	c,err=knowledgeindex.Save(pg,c)
	if err!=nil{t.Fatal(err)}
	body:=strings.Repeat("甲",800)+strings.Repeat("乙",100)
	doc,err:=pg.SaveKnowledgeDocument(ctx,tid,"",db.KnowledgeDocument{Title:"incremental fixture",Content:body})
	if err!=nil{t.Fatal(err)}
	global,err:=pg.SaveKnowledgeDocument(ctx,0,"",db.KnowledgeDocument{Title:"global pending fixture",Content:"global"})
	if err!=nil{t.Fatal(err)}
	defer pg.DeleteKnowledgeDocument(ctx,0,global.ID)
	pending:="/api/knowledge/documents/index-pending?task="+task
	checkPending:=func(want bool){
		t.Helper()
		w:=req("GET",pending,"")
		if w.Code!=200||strings.Contains(w.Body.String(),doc.ID)!=want||strings.Contains(w.Body.String(),global.ID){t.Fatal(w.Code,w.Body)}
	}
	checkPending(true)
	path:="/api/knowledge/documents/"+doc.ID+"/index?task="+task+"&incremental=true"
	build:=func()knowledgeindex.BuildResult{
		t.Helper()
		w:=req("POST",path,"{}")
		var result knowledgeindex.BuildResult
		if err:=json.Unmarshal(w.Body.Bytes(),&result);w.Code!=200||err!=nil{t.Fatal(w.Code,w.Body,err)}
		return result
	}
	first:=build()
	if first.Chunks!=2||first.Embedded!=2||first.Reused!=0||inputs.Load()!=2{t.Fatal(first,inputs.Load())}
	checkPending(false)
	second:=build()
	if !second.Skipped||second.Reused!=2||inputs.Load()!=2{t.Fatal(second,inputs.Load())}
	doc,err=pg.SaveKnowledgeDocument(ctx,tid,doc.ID,db.KnowledgeDocument{Title:doc.Title,Content:strings.Repeat("甲",800)+strings.Repeat("丙",100)})
	if err!=nil{t.Fatal(err)}
	checkPending(true)
	third:=build()
	if third.Chunks!=2||third.Reused!=1||third.Embedded!=1||inputs.Load()!=3{t.Fatal(third,inputs.Load())}
	checkPending(false)
	// Explicit rebuild ignores reuse even when content is unchanged.
	if n,err:=knowledgeindex.Build(ctx,pg,tid,doc.ID,c);err!=nil||n!=2||inputs.Load()!=5{t.Fatal(n,err,inputs.Load())}
}
