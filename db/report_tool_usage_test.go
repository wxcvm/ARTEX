package db

import "testing"

func TestReportToolUsageScopesAndPairsResults(t *testing.T){
	d,err:=Open(testDSN(t));if err!=nil{t.Skipf("postgres unavailable: %v",err)};defer d.Close()
	exp,err:=d.CreateExploration("report usage","fixture");if err!=nil{t.Fatal(err)}
	defer d.Exec("DELETE FROM explorations WHERE id=$1",exp)
	other,err:=d.CreateExploration("other usage","fixture");if err!=nil{t.Fatal(err)}
	defer d.Exec("DELETE FROM explorations WHERE id=$1",other)
	add:=func(id int64,kind,tool,call,worker string,failed bool){
		t.Helper()
		if _,err:=d.Exploration(id).AppendActivity(Activity{Kind:kind,Tool:tool,ToolUseID:call,Worker:worker,IsError:failed,Detail:"secret fixture"});err!=nil{t.Fatal(err)}
	}
	add(exp,"tool_use","probe","shared","worker",false)
	add(other,"tool_result","probe","shared","worker",true)
	add(exp,"tool_result","probe","shared","other-worker",true)
	add(exp,"tool_result","probe","shared","worker",false)
	add(exp,"tool_use","probe","shared","worker",false)
	add(exp,"tool_result","probe","shared","worker",true)
	add(exp,"tool_result","probe","shared","worker",true)
	add(exp,"tool_use","probe","pending","worker",false)
	add(exp,"tool_use","lookup","","worker",false)
	add(other,"tool_use","unrelated","other","worker",false)
	got,err:=d.ReportToolUsages(t.Context(),exp)
	if err!=nil||len(got)!=2{t.Fatal(got,err)}
	if got[0]!=(ReportToolUsage{Tool:"probe",Total:3,Returned:2,Errors:1,Pending:1})||got[1]!=(ReportToolUsage{Tool:"lookup",Total:1,Pending:1}){t.Fatal(got)}
}
