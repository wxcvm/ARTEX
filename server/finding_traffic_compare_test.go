package server

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/traffic"
)

func TestCompareFindingTrafficAfterHotCleanup(t *testing.T) {
	s,f,request := trafficEvidenceServer(t)
	seedServerEvidenceFlow(t,s,"compare-a",[]byte("before"))
	seedServerEvidenceFlow(t,s,"compare-b",[]byte("after"))
	base := fmt.Sprintf("/api/exploration/findings/%d/traffic",f.FindingID)
	w := request("POST",base,`{"traffic_refs":[{"traffic_id":"compare-a","role":"baseline"},{"traffic_id":"compare-b","role":"verification"}]}`)
	var list db.FindingTraffic
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(),&list) != nil || len(list.Bindings)!=2 { t.Fatalf("bind: %d %s",w.Code,w.Body) }
	if _,err:=s.m.traffic.DeleteHost("evidence.local");err!=nil { t.Fatal(err) }
	path := fmt.Sprintf("%s/compare?baseline=%d&candidate=%d&version=%d",base,list.Bindings[0].ID,list.Bindings[1].ID,list.Version)
	w = request("GET",path,"")
	var result struct { traffic.RecordedComparison; EvidenceVersion int64 `json:"evidence_version"` }
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(),&result)!=nil || result.EvidenceVersion!=list.Version || result.ResponseBody.Relation!="different" || result.RequestBody.Relation!="equal" {
		t.Fatalf("snapshot compare: %d %s",w.Code,w.Body)
	}
	if w:=request("GET",fmt.Sprintf("%s/compare?baseline=%d&candidate=%d&version=%d",base,list.Bindings[0].ID,list.Bindings[1].ID,list.Version+1),"");w.Code!=409 { t.Fatalf("stale version: %d %s",w.Code,w.Body) }
	if w:=request("GET",fmt.Sprintf("%s/compare?baseline=%d&candidate=9223372036854775807",base,list.Bindings[0].ID),"");w.Code!=404 { t.Fatalf("unbound snapshot: %d %s",w.Code,w.Body) }
}
