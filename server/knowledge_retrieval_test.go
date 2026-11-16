package server

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/kb"
	actool "github.com/Autumn-27/norma/tool"
)

func TestScopedKnowledgeToolResolution(t *testing.T) {
	s,f,_:=trafficEvidenceServer(t)
	pg:=s.m.pg
	finding,err:=pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	tid:=*finding.TaskID
	other,err:=s.m.CreateTask("other knowledge fixture","local",nil,0,0)
	if err!=nil{t.Fatal(err)}
	otherID,_:=strconv.ParseInt(other.ID,10,64)
	defer pg.DeleteTask(otherID)
	oldResolve:=agent.ToolResolve
	oldKB:=agent.Knowledge()
	defer func(){agent.ToolResolve=oldResolve;agent.SetKnowledge(oldKB)}()
	// Exercise a deployment without the optional bundled SQLite file.
	agent.SetKnowledge(nil)
	wireTools(pg,nil)
	ctx:=context.Background()
	key:="scope-fixture-"+strconv.FormatInt(tid,10)
	save:=func(task int64,title,content string)db.KnowledgeDocument{
		t.Helper()
		doc,err:=pg.SaveKnowledgeDocument(ctx,task,"",db.KnowledgeDocument{Title:title,Content:content,Category:"fixture"})
		if err!=nil{t.Fatal(err)}
		return doc
	}
	global:=save(0,"Global "+key,"public reference "+key)
	defer pg.DeleteKnowledgeDocument(ctx,0,global.ID)
	local:=save(tid,"Task "+key,strings.Repeat("中文",6500)+" "+key)
	foreign:=save(otherID,"Foreign "+key,"foreign "+key)
	bind:=func(task int64)map[string]actool.CoreTool{
		return names(agent.ToolResolve(agent.WithRunInfo(ctx,agent.RunInfo{TaskID:task}),"mainagent",nil))
	}
	call:=func(tool actool.CoreTool, args map[string]any)(actool.Result,error){
		t.Helper()
		if tool==nil{t.Fatal("missing resolved knowledge tool")}
		raw,err:=json.Marshal(args)
		if err!=nil{t.Fatal(err)}
		return tool.Call(ctx,raw,nil)
	}
	taskTools:=bind(tid)
	chatTools:=bind(0)
	search:=func(tools map[string]actool.CoreTool,want []string,exclude []string){
		t.Helper()
		// A fabricated task_id cannot change the scope captured at assembly.
		res,err:=call(tools["kb_search"],map[string]any{"q":key,"limit":8,"task_id":otherID})
		if err!=nil||res.IsError{t.Fatal(err,res.Flatten())}
		for _,id:=range want{if !strings.Contains(res.Flatten(),id){t.Fatal("missing hit",id,res.Flatten())}}
		for _,id:=range exclude{if strings.Contains(res.Flatten(),id){t.Fatal("leaked hit",id,res.Flatten())}}
	}
	search(taskTools,[]string{global.ID,local.ID},[]string{foreign.ID})
	search(chatTools,[]string{global.ID},[]string{local.ID,foreign.ID})
	res,err:=call(taskTools["kb_get"],map[string]any{"document_id":local.ID})
	if err!=nil||res.IsError{t.Fatal(err,res.Flatten())}
	var page struct{Content string `json:"content"`;Next int `json:"next_offset"`;Scope string `json:"scope"`}
	if err=json.Unmarshal([]byte(res.Flatten()),&page);err!=nil||len([]rune(page.Content))!=12000||page.Next!=12000||page.Scope!="task"{t.Fatal("bad page",page.Next,err)}
	res,err=call(taskTools["kb_get"],map[string]any{"document_id":local.ID,"offset":page.Next})
	if err!=nil||res.IsError||!strings.Contains(res.Flatten(),key){t.Fatal("missing second page",err)}
	for _,tc:=range []struct{tools map[string]actool.CoreTool;id string}{
		{taskTools,foreign.ID},{chatTools,local.ID},
	}{
		res,err=call(tc.tools["kb_get"],map[string]any{"document_id":tc.id,"task_id":tid})
		if err!=nil||!res.IsError{t.Fatal("cross-scope read accepted",err,res.Flatten())}
	}
	// Legacy schema still exposes numeric id, but no longer requires it for
	// a document-only call. Existing enabled/bound settings are not rewritten.
	row,err:=pg.GetTool("kb_get")
	if err!=nil||row==nil{t.Fatal(err)}
	defer pg.UpdateTool(row.Key,row.Description,row.Schema,mustJSON(row.Agents),row.Enabled)
	legacy:=mustJSON(map[string]any{"type":"object","properties":map[string]any{"id":map[string]any{"type":"integer"}},"required":[]string{"id"}})
	if err=pg.UpdateTool(row.Key,"custom legacy description",legacy,mustJSON([]string{"mainagent"}),true);err!=nil{t.Fatal(err)}
	tool:=bind(tid)["kb_get"]
	if tool==nil||!strings.Contains(tool.Description(),"custom legacy description"){t.Fatal("lost custom description")}
	schema:=tool.InputSchema()
	props:=schema["properties"].(map[string]any)
	if props["id"]==nil||props["document_id"]==nil||schema["required"]!=nil{t.Fatal("legacy schema incompatible",schema)}
	if err=pg.UpdateTool(row.Key,row.Description,row.Schema,mustJSON([]string{"mainagent"}),false);err!=nil{t.Fatal(err)}
	if bind(tid)["kb_get"]!=nil{t.Fatal("disabled tool injected")}
	if err=pg.UpdateTool(row.Key,row.Description,row.Schema,mustJSON(row.Agents),row.Enabled);err!=nil{t.Fatal(err)}

	// Bundled entries retain their original numeric identifier contract.
	bundled,err:=kb.Open(filepath.Join("..","knowledge"))
	if err!=nil||bundled==nil{t.Fatal("open bundled",err)}
	defer bundled.Close()
	agent.SetKnowledge(bundled)
	tools:=bind(tid)
	res,err=call(tools["kb_search"],map[string]any{"q":"PostgreSQL","limit":8})
	if err!=nil||res.IsError{t.Fatal(err,res.Flatten())}
	var hits struct{Hits []struct{ID int64 `json:"id"`;Scope string `json:"scope"`} `json:"hits"`}
	if err=json.Unmarshal([]byte(res.Flatten()),&hits);err!=nil{t.Fatal(err)}
	var id int64
	for _,hit:=range hits.Hits{if hit.Scope=="bundled"&&hit.ID>0{id=hit.ID;break}}
	if id==0{t.Fatal("bundled result missing",res.Flatten())}
	res,err=call(tools["kb_get"],map[string]any{"id":id})
	if err!=nil||res.IsError||!strings.Contains(res.Flatten(),`"scope":"bundled"`){t.Fatal("legacy read failed",err,res.Flatten())}
}
