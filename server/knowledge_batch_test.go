package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/knowledgeindex"
)

func TestBackgroundKnowledgeBatchCancelAndResume(t *testing.T) {
	s,f,req:=trafficEvidenceServer(t)
	ctx,cancel:=context.WithCancel(context.Background())
	defer cancel()
	s.ctx=ctx
	pg:=s.m.pg
	finding,err:=pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	tid:=*finding.TaskID
	task:=strconv.FormatInt(tid,10)
	old,exists,err:=pg.GetSetting(knowledgeindex.SettingKey)
	if err!=nil{t.Fatal(err)}
	defer func(){if exists{_=pg.SetSetting(knowledgeindex.SettingKey,old)}else{_,_=pg.Exec("DELETE FROM settings WHERE key=$1",knowledgeindex.SettingKey)}}()
	var calls atomic.Int64
	started:=make(chan struct{})
	release:=make(chan struct{})
	mock:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		var in struct{Input []string `json:"input"`}
		if err:=json.NewDecoder(r.Body).Decode(&in);err!=nil{w.WriteHeader(400);return}
		if calls.Add(1)==1{
			close(started)
			select{case <-release:case <-r.Context().Done():return}
		}
		out:=[]any{}
		for i:=range in.Input{out=append(out,map[string]any{"index":i,"embedding":[]float32{1,0}})}
		_=json.NewEncoder(w).Encode(map[string]any{"data":out})
	}))
	defer mock.Close()
	defer close(release)
	c,err:=knowledgeindex.Save(pg,knowledgeindex.Config{Enabled:true,BaseURL:mock.URL,Model:"background-fixture"})
	if err!=nil{t.Fatal(err)}
	doc,err:=pg.SaveKnowledgeDocument(ctx,tid,"",db.KnowledgeDocument{Title:"background fixture",Content:"document for local model"})
	if err!=nil{t.Fatal(err)}
	global,err:=pg.SaveKnowledgeDocument(ctx,0,"",db.KnowledgeDocument{Title:"global fixture",Content:"not this task"})
	if err!=nil{t.Fatal(err)}
	defer pg.DeleteKnowledgeDocument(ctx,0,global.ID)
	path:="/api/knowledge/documents/index-batch?task="+task
	request:=func(method string)knowledgeBatchStatus{
		t.Helper()
		w:=req(method,path,"")
		var st knowledgeBatchStatus
		if err:=json.Unmarshal(w.Body.Bytes(),&st);err!=nil||(w.Code!=200&&w.Code!=202){t.Fatal(w.Code,w.Body,err)}
		return st
	}
	wait:=func(state string)knowledgeBatchStatus{
		t.Helper()
		deadline:=time.Now().Add(10*time.Second)
		for time.Now().Before(deadline){
			st:=request("GET")
			if st.State==state{return st}
			time.Sleep(10*time.Millisecond)
		}
		t.Fatalf("did not reach %s: %+v",state,request("GET"))
		return knowledgeBatchStatus{}
	}
	if request("POST").State!="running"{t.Fatal("job did not start")}
	select{case <-started:case <-time.After(5*time.Second):t.Fatal("embedding did not start")}
	if request("POST").State!="running"{t.Fatal("duplicate request replaced active job")}
	if st:=request("DELETE");st.State!="stopping"&&st.State!="cancelled"{t.Fatal(st)}
	if st:=wait("cancelled");st.Succeeded!=0||st.Failed!=0{t.Fatal(st)}
	if n,err:=pg.KnowledgeVectorCount(ctx,tid,doc.ID,c.Key());err!=nil||n!=0{t.Fatal(n,err)}
	request("POST")
	st:=wait("finished")
	if st.Succeeded!=1||st.Embedded!=1||st.Failed!=0{t.Fatal(st)}
	if n,err:=pg.KnowledgeVectorCount(ctx,0,global.ID,c.Key());err!=nil||n!=0{t.Fatal(n,err)}
	request("POST")
	if st:=wait("finished");st.Succeeded!=0||st.Embedded!=0{t.Fatal(st)}
	if calls.Load()!=2{t.Fatalf("embedding calls = %d",calls.Load())}
}
