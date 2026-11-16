package server

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestTaskReportRecordedToolUsage(t *testing.T){
	s,f,req:=trafficEvidenceServer(t)
	finding,err:=s.m.pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	task:=s.m.ResolveTask(fmt.Sprint(*finding.TaskID));if task==nil{t.Fatal("missing runtime task")}
	for _,a:=range []db.Activity{
		{Kind:"tool_use",Tool:"fixture-probe",ToolUseID:"one",Worker:"fixture",Detail:"PRIVATE_INPUT"},
		{Kind:"tool_result",ToolUseID:"one",Worker:"fixture",Detail:"PRIVATE_OUTPUT"},
		{Kind:"tool_use",Tool:"fixture-probe",ToolUseID:"two",Worker:"fixture"},
	}{if _,err:=task.Store.AppendActivity(a);err!=nil{t.Fatal(err)}}
	w:=req("GET",fmt.Sprintf("/api/report?task=%s&lang=en",task.ID),"")
	if w.Code!=200||!strings.Contains(w.Body.String(),"fixture-probe: calls 2; returned 1; errors 0; no result 1"){t.Fatal(w.Code,w.Body)}
	if strings.Contains(w.Body.String(),"PRIVATE_INPUT")||strings.Contains(w.Body.String(),"PRIVATE_OUTPUT"){t.Fatal("raw tool data leaked")}
}
