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
	s,f,req:=trafficEvidenceServer(t)
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
	for _,tc:=range []struct{path string;want []string;exclude []string}{
		{"/api/knowledge/retrievals?task="+strconv.FormatInt(tid,10),[]string{key,global.ID,local.ID},[]string{foreign.ID}},
		{"/api/knowledge/retrievals",[]string{key,global.ID},[]string{local.ID,foreign.ID}},
		{"/api/knowledge/retrievals?task="+other.ID,nil,[]string{key,local.ID,global.ID}},
	}{
		w:=req("GET",tc.path,"")
		if w.Code!=200{t.Fatal(w.Code,w.Body)}
		for _,id:=range tc.want{if !strings.Contains(w.Body.String(),id){t.Fatal("missing history",id,w.Body)}}
		for _,id:=range tc.exclude{if strings.Contains(w.Body.String(),id){t.Fatal("history leaked",id,w.Body)}}
	}
	defer pg.Exec("DELETE FROM knowledge_retrievals WHERE task_id IS NULL AND query=$1",key)
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


func TestRankedKnowledgeDocumentSearch(t *testing.T){
	s,f,_:=trafficEvidenceServer(t)
	pg:=s.m.pg
	ctx:=context.Background()
	finding,err:=pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	tid:=*finding.TaskID
	key:="ranking-"+strconv.FormatInt(tid,10)
	save:=func(title,body string)db.KnowledgeDocument{
		t.Helper()
		doc,err:=pg.SaveKnowledgeDocument(ctx,tid,"",db.KnowledgeDocument{Title:title,Content:body})
		if err!=nil{t.Fatal(err)}
		return doc
	}
	exact:=save(key+" 部署指南","some reference")
	multi:=save("部署 notes",key+" intervening 指南")
	_ = multi
	partial:=save("unrelated",key)
	hits,err:=pg.SearchKnowledgeDocuments(ctx,tid,key+" 部署指南",10)
	if err!=nil||len(hits)!=1||hits[0].ID!=exact.ID{t.Fatal("phrase match",hits,err)}
	// Two terms can appear in separate fields and need not be adjacent.
	hits,err=pg.SearchKnowledgeDocuments(ctx,tid,key+" 部署",10)
	if err!=nil||len(hits)!=2||hits[0].ID!=exact.ID||hits[1].ID!=multi.ID{t.Fatal("ranked terms",hits,err)}
	for _,hit:=range hits{if hit.ID==partial.ID{t.Fatal("partial term matched")}}
}


// Exercise all 27 shipped references through the same scoped tools
// exposed to the agent, not merely the manifest or the database insert.
func TestCyberStrikeReferenceToolRetrieval(t *testing.T) {
	s,f,_:=trafficEvidenceServer(t)
	pg:=s.m.pg
	ctx:=context.Background()
	finding,err:=pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	tid:=*finding.TaskID
	if err:=pg.SeedCyberStrikeReferencePack();err!=nil{t.Fatal(err)}
	// A task-owned document mentions the same code but is not the
	// canonical code-titled reference. It must still be searchable.
	decoy,err:=pg.SaveKnowledgeDocument(ctx,tid,"",db.KnowledgeDocument{Title:"NIST CSF PR.AA-05 project notes",Content:"A project note about PR.AA-05, not the reference.",Category:"fixture"})
	if err!=nil{t.Fatal(err)}
	tools:=names(agent.ScopedKnowledgeTools(pg,tid,"cyberstrike-reference-fixture"))
	call:=func(name string,args map[string]any)string{
		t.Helper()
		raw,err:=json.Marshal(args)
		if err!=nil{t.Fatal(err)}
		tool:=tools[name]
		if tool==nil{t.Fatal("missing tool",name)}
		res,err:=tool.Call(ctx,raw,nil)
		if err!=nil||res.IsError{t.Fatal(name,err,res.Flatten())}
		return res.Flatten()
	}
	type hit struct{ DocumentID string `json:"document_id"`; Title string `json:"title"`; Scope string `json:"scope"`; SourceFilename string `json:"source_filename"` }
	search:=func(q string,limit int)[]hit{
		t.Helper()
		var out struct{Hits []hit `json:"hits"`}
		if err:=json.Unmarshal([]byte(call("kb_search",map[string]any{"q":q,"limit":limit})),&out);err!=nil{t.Fatal(err)}
		return out.Hits
	}
	for _,tc:=range []struct{q,code string}{
		{"ID.AM-01","ID.AM-01"},{"ID.RA-01","ID.RA-01"},{"pr.aa-05","PR.AA-05"},
		{"NIST CSF PR.AA-05","PR.AA-05"},{"PR.DS-11","PR.DS-11"},{"PR.PS-02","PR.PS-02"},
		{"DE.CM-01","DE.CM-01"},{"RS.AN","RS.AN"},{"RC.RP","RC.RP"},{"GV.SC","GV.SC"},
		{"最小权限","PR.AA-05"},{"软件生命周期","PR.PS-02"},
		{"GV.RR-02","GV.RR-02"},{"ID.AM-02","ID.AM-02"},{"ID.AM-03","ID.AM-03"},
		{"ID.IM-01","ID.IM-01"},{"PR.AA-03","PR.AA-03"},{"PR.DS-01","PR.DS-01"},
		{"DE.AE-03","DE.AE-03"},{"RS.MA-02","RS.MA-02"},{"RC.RP-05","RC.RP-05"},
		{"GV.RM-06","GV.RM-06"},{"GV.SC-07","GV.SC-07"},{"ID.RA-06","ID.RA-06"},
		{"ID.RA-07","ID.RA-07"},{"PR.DS-02","PR.DS-02"},{"PR.PS-04","PR.PS-04"},
		{"DE.CM-09","DE.CM-09"},{"RS.AN-06","RS.AN-06"},{"RC.CO-03","RC.CO-03"},
	}{
		hits:=search(tc.q,1)
		if len(hits)!=1||hits[0].Scope!="global"||!strings.HasPrefix(hits[0].Title,tc.code+" ")||hits[0].DocumentID==""||!strings.HasPrefix(hits[0].SourceFilename,".cyberstrike/skill/NIST/CSF_v2.0/"){
			t.Fatalf("query %q: %+v",tc.q,hits)
		}
		var doc struct{Content string `json:"content"`; Scope string `json:"scope"`; SourceFilename string `json:"source_filename"`}
		if err:=json.Unmarshal([]byte(call("kb_get",map[string]any{"document_id":hits[0].DocumentID})),&doc);err!=nil{t.Fatal(err)}
		if doc.Scope!="global"||doc.SourceFilename!=hits[0].SourceFilename||!strings.Contains(doc.Content,"9303b210ab82bd6362d7ad7fbd2b2c8f2e335450")||!strings.Contains(doc.Content,"Reference only"){
			t.Fatalf("missing reference provenance for %q: %+v",tc.q,doc)
		}
	}
	hits:=search("NIST CSF PR.AA-05",8)
	seenDecoy:=false
	for _,h:=range hits{if h.DocumentID==decoy.ID{seenDecoy=true}}
	if !seenDecoy{t.Fatal("ordinary task document no longer searchable")}
}
