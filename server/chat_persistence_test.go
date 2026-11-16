package server

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestUserMessageWriteFailureDoesNotStartTurn(t *testing.T){
	s,f,req:=trafficEvidenceServer(t)
	pg:=s.m.pg
	finding,err:=pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	taskID:=strconv.FormatInt(*finding.TaskID,10)
	task:=s.m.ResolveTask(taskID)
	if task==nil{t.Fatal("missing task")}
	// A targeted NOT VALID constraint affects only new user rows for this
	// fixture, without deleting existing records or contacting a model.
	taskConstraint:=fmt.Sprintf("chat_write_fixture_%d",task.ExpID)
	if _,err:=pg.Exec(fmt.Sprintf("ALTER TABLE activity ADD CONSTRAINT %s CHECK (exploration_id<>%d OR kind<>'user') NOT VALID",taskConstraint,task.ExpID));err!=nil{t.Fatal(err)}
	defer pg.Exec("ALTER TABLE activity DROP CONSTRAINT "+taskConstraint)
	for i:=0;i<2;i++{
		w:=req("POST","/api/chat?task="+taskID,`{"message":"fixture message"}`)
		if w.Code!=503||!strings.Contains(w.Body.String(),"no turn was started"){t.Fatal(w.Code,w.Body)}
		if w=req("GET","/api/tasks/"+taskID+"/chat/status","");w.Code!=200||!strings.Contains(w.Body.String(),`"running":false`){t.Fatal(w.Code,w.Body)}
	}
	profile,err:=pg.SaveProfile(&db.LLMProfile{Name:"write-failure-fixture-"+taskID,Format:"anthropic",Model:"fixture",APIKey:"sk-test"})
	if err!=nil{t.Fatal(err)}
	defer pg.DeleteProfile(profile)
	conv,err:=pg.CreateConversation("auto","write failure fixture",&profile)
	if err!=nil{t.Fatal(err)}
	defer pg.DeleteConversation(conv.ID)
	convConstraint:=fmt.Sprintf("conv_write_fixture_%d",conv.ID)
	if _,err:=pg.Exec(fmt.Sprintf("ALTER TABLE conversation_activities ADD CONSTRAINT %s CHECK (conversation_id<>%d OR kind<>'user') NOT VALID",convConstraint,conv.ID));err!=nil{t.Fatal(err)}
	defer pg.Exec("ALTER TABLE conversation_activities DROP CONSTRAINT "+convConstraint)
	path:=fmt.Sprintf("/api/conversations/%d/messages",conv.ID)
	for i:=0;i<2;i++{
		w:=req("POST",path,`{"message":"fixture message"}`)
		if w.Code!=503||!strings.Contains(w.Body.String(),"no turn was started"){t.Fatal(w.Code,w.Body)}
		s.chatMu.Lock()
		busy:=s.chatBusy[s.convBusyKey(conv.ID)]
		s.chatMu.Unlock()
		if busy{t.Fatal("failed write left conversation busy")}
	}
}
