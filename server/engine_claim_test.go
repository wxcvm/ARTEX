package server

import (
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestClaimNextReportsFrontierReadFailure(t *testing.T) {
	dsn,_,err:=db.DSN()
	if err!=nil{t.Skip("test PostgreSQL not configured")}
	pg,err:=db.Open(dsn)
	if err!=nil{t.Skipf("test PostgreSQL unavailable: %v",err)}
	store:=pg.Exploration(1)
	if err:=pg.Close();err!=nil{t.Fatal(err)}
	engine:=NewEngine(nil)
	node,err:=engine.claimNext(&Task{ID:"fixture",Store:store},"worker")
	if err==nil||!strings.Contains(err.Error(),"list open intents")||node!=nil{
		t.Fatalf("closed-DB frontier returned node=%+v error=%v; want a reported query error",node,err)
	}
}
