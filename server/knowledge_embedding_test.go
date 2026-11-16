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

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/knowledgeindex"
)

func TestLocalKnowledgeSemanticIndexLifecycle(t *testing.T){
	s,f,req:=trafficEvidenceServer(t)
	pg:=s.m.pg
	ctx:=context.Background()
	finding,err:=pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	tid:=*finding.TaskID
	task:=strconv.FormatInt(tid,10)
	other,err:=s.m.CreateTask("semantic scope fixture","local",nil,0,0)
	if err!=nil{t.Fatal(err)}
	otherID,_:=strconv.ParseInt(other.ID,10,64)
	defer pg.DeleteTask(otherID)
	old,exists,err:=pg.GetSetting(knowledgeindex.SettingKey)
	if err!=nil{t.Fatal(err)}
	defer func(){if exists{_ = pg.SetSetting(knowledgeindex.SettingKey,old)}else{_,_=pg.Exec("DELETE FROM settings WHERE key=$1",knowledgeindex.SettingKey)}}()
	var calls atomic.Int64
	mock:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		calls.Add(1)
		if r.URL.Path!="/v1/embeddings"{t.Errorf("unexpected path %s",r.URL.Path)}
		var in struct{Input []string `json:"input"`;Model string `json:"model"`}
		if err:=json.NewDecoder(r.Body).Decode(&in);err!=nil{t.Error(err);w.WriteHeader(400);return}
		out:=[]any{}
		for i,text:=range in.Input{
			v:=[]float32{1,0}
			if strings.Contains(text,"second topic"){v=[]float32{0,1}}
			out=append(out,map[string]any{"index":i,"embedding":v})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data":out})
	}))
	defer mock.Close()
	c:=knowledgeindex.Config{Enabled:true,BaseURL:mock.URL+"/v1",Model:"fixture"}
	saveConfig:=func(c knowledgeindex.Config){
		raw,_:=json.Marshal(c)
		if w:=req("PUT","/api/knowledge/embedding",string(raw));w.Code!=200{t.Fatal(w.Code,w.Body)}
	}
	saveConfig(c)
	save:=func(taskID int64,title,body string)db.KnowledgeDocument{
		t.Helper()
		doc,err:=pg.SaveKnowledgeDocument(ctx,taskID,"",db.KnowledgeDocument{Title:title,Content:body})
		if err!=nil{t.Fatal(err)}
		return doc
	}
	local:=save(tid,"Task local","nebula actual record")
	global:=save(0,"Global local","second topic")
	defer pg.DeleteKnowledgeDocument(ctx,0,global.ID)
	foreign:=save(otherID,"Foreign local","nebula foreign record")
	if _,attempted,err:=knowledgeindex.Search(ctx,pg,tid,"alias",5);err!=nil||attempted||calls.Load()!=0{t.Fatal("called model before index",attempted,err,calls.Load())}
	index:=func(doc db.KnowledgeDocument,scope string){
		t.Helper()
		path:="/api/knowledge/documents/"+doc.ID+"/index"+scope
		if w:=req("POST",path,"{}");w.Code!=200{t.Fatal(w.Code,w.Body)}
		if w:=req("GET",path,"");w.Code!=200||!strings.Contains(w.Body.String(),`"chunks":1`){t.Fatal(w.Code,w.Body)}
	}
	index(local,"?task="+task)
	index(global,"")
	index(foreign,"?task="+other.ID)
	hits,attempted,err:=knowledgeindex.Search(ctx,pg,tid,"alias",5)
	if err!=nil||!attempted||len(hits)!=2||hits[0].Document.ID!=local.ID||hits[1].Document.ID!=global.ID{t.Fatal("incorrect scoped ranking",hits,attempted,err)}
	// Actual agent tool returns a semantic hit without literal keyword overlap.
	tool:=agent.ScopedKnowledgeTools(pg,tid,"fixture")[0]
	res,err:=tool.Call(ctx,json.RawMessage(`{"q":"alias","limit":8}`),nil)
	if err!=nil||res.IsError||!strings.Contains(res.Flatten(),local.ID)||strings.Contains(res.Flatten(),foreign.ID)||!strings.Contains(res.Flatten(),`"match":"semantic"`){t.Fatal(err,res.Flatten())}
	res,err=tool.Call(ctx,json.RawMessage(`{"q":"nebula","limit":8}`),nil)
	if err!=nil||res.IsError||!strings.Contains(res.Flatten(),`"match":"hybrid"`){t.Fatal(err,res.Flatten())}
	// Editing invalidates the old document version. New content cannot use stale vectors.
	updated,err:=pg.SaveKnowledgeDocument(ctx,tid,local.ID,db.KnowledgeDocument{Title:local.Title,Content:"updated keyword"})
	if err!=nil{t.Fatal(err)}
	hits,_,err=knowledgeindex.Search(ctx,pg,tid,"alias",5)
	if err!=nil||len(hits)!=1||hits[0].Document.ID!=global.ID{t.Fatal("stale index remained",hits,err)}
	if err:=pg.ReplaceKnowledgeVectors(ctx,tid,local.ID,c.Key(),local.UpdatedAt,[]db.KnowledgeVector{{Embedding:[]float32{1,0}}});err!=db.ErrKnowledgeIndexStale{t.Fatal("stale build accepted",err)}
	index(updated,"?task="+task)
	// Disabled search does not contact the provider.
	disabled:=c;disabled.Enabled=false;saveConfig(disabled)
	before:=calls.Load()
	if _,attempted,err:=knowledgeindex.Search(ctx,pg,tid,"alias",5);err!=nil||attempted||calls.Load()!=before{t.Fatal("disabled call",attempted,err)}
	saveConfig(c)
	// Failure of the configured local model must preserve keyword results.
	mock.Close()
	res,err=tool.Call(ctx,json.RawMessage(`{"q":"updated keyword","limit":8}`),nil)
	if err!=nil||res.IsError||!strings.Contains(res.Flatten(),local.ID)||!strings.Contains(res.Flatten(),`"degraded":true`){t.Fatal("keyword fallback failed",err,res.Flatten())}
}
