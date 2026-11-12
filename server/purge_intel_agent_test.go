package server

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
)

func TestTaskPurgeIntelCompanyBoundary(t *testing.T) {
	task := &db.Task{CompanyIDs:[]int64{4,7}}
	for _, tc := range []struct{ id int64; allowed bool }{{4,true},{7,true},{8,false},{0,false},{-1,false}} {
		if got:=taskOwnsIntelCompany(task,tc.id);got!=tc.allowed { t.Fatalf("company %d allowed=%v",tc.id,got) }
	}
	if taskOwnsIntelCompany(nil,4) {t.Fatal("missing task allowed")}
	if _,err:=readTaskPurgeIntel(context.Background(),nil,0,purgeIntelAgentRequest{CompanyID:4});err==nil {t.Fatal("taskless read allowed")}
}

func TestTaskPurgeIntelReadDatabase(t *testing.T) {
	dsn,_,err:=db.DSN();if err!=nil{t.Skip("postgres not configured")}
	d,err:=db.Open(dsn);if err!=nil{t.Skipf("postgres unavailable: %v",err)};defer d.Close()
	cs:=d.Companies()
	id,_,_,_,_,err:=cs.CreateCompanyWithScope(fmt.Sprintf("purge-agent-%d",time.Now().UnixNano()),"",[]db.ScopeInput{
		{Kind:"domain",Value:"example.test"},{Kind:"cidr",Value:"192.0.2.0/24"}},"test")
	if err!=nil{t.Fatal(err)}
	defer d.Exec("DELETE FROM companies WHERE id=$1",id)
	task,err:=d.CreateTaskWithOptions("purge intel read test","review saved evidence",db.TaskCreateOptions{CompanyIDs:[]int64{id}})
	if err!=nil{t.Fatal(err)}
	defer d.Exec("DELETE FROM tasks WHERE id=$1",task.ID)
	b:=fixturePurgeIntel()
	payload,err:=json.Marshal(b);if err!=nil{t.Fatal(err)}
	var sid int64
	if err=d.QueryRow("INSERT INTO purge_intel_snapshots(company_id,digest,payload) VALUES($1,$2,$3::jsonb) RETURNING id",id,"fixture",string(payload)).Scan(&sid);err!=nil{t.Fatal(err)}
	ctx:=context.Background()
	list,err:=readTaskPurgeIntel(ctx,d,task.ID,purgeIntelAgentRequest{CompanyID:id})
	if err!=nil{t.Fatal(err)}
	items:=list.(map[string]any)["items"].([]purgeIntelAgentSnapshot)
	if len(items)!=1||items[0].ID!=sid {t.Fatalf("bad snapshot list: %+v",items)}
	v,err:=readTaskPurgeIntel(ctx,d,task.ID,purgeIntelAgentRequest{CompanyID:id,SnapshotID:sid})
	if err!=nil{t.Fatal(err)}
	m:=v.(map[string]any)
	if m["historical_only"]!=true||m["port_observations"]!=1||m["snapshot_id"]!=sid {t.Fatalf("bad analysis: %+v",m)}
	for _,r:=range []purgeIntelAgentRequest{{CompanyID:id+1,SnapshotID:sid},{CompanyID:id,SnapshotID:sid+1}} {
		if _,err:=readTaskPurgeIntel(ctx,d,task.ID,r);err==nil{t.Fatal("unassociated company or missing snapshot accepted")}
	}
	b.IPs[0].IP="198.51.100.1"
	payload,_=json.Marshal(b)
	if _,err=d.Exec("UPDATE purge_intel_snapshots SET payload=$1::jsonb WHERE id=$2",string(payload),sid);err!=nil{t.Fatal(err)}
	if _,err:=readTaskPurgeIntel(ctx,d,task.ID,purgeIntelAgentRequest{CompanyID:id,SnapshotID:sid});err==nil{t.Fatal("out-of-scope snapshot accepted")}
}
