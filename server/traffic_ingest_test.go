package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/traffic"
)

// ingestServer wires a Server with a live traffic store, like
// trafficEvidenceServer but without the finding plumbing.
func ingestServer(t *testing.T) (*Server, func(method, path, body string) *httptest.ResponseRecorder) {
	t.Helper()
	m, err := NewManager(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { m.Close() })
	m.traffic, err = traffic.Open(filepath.Join(m.dir, "traffic"), ":0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s := New(ctx, m, t.TempDir(), t.TempDir(), t.TempDir())
	s.archiveWG.Wait()
	<-s.side.done
	token, err := signJWT(s.jwtKey)
	if err != nil {
		t.Fatal(err)
	}
	return s, func(method, path, body string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, path, strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, r)
		return w
	}
}

func TestTrafficIngestFlow(t *testing.T) {
	s, request := ingestServer(t)

	body := fmt.Sprintf(`{"exchanges":[
		{"method":"POST","url":"https://api.example.com/v1/login","status":200,
		 "requestHeaders":{"content-type":"application/json"},"requestBody":"{\"u\":\"u1\"}",
		 "responseHeaders":{"content-type":"application/json"},"responseBody":"{\"token\":\"t\"}"},
		{"method":"GET","url":"https://api.example.com/v1/profile","status":200,
		 "requestHeaders":{},"requestBody":"",
		 "responseHeaders":{},"responseBody":"{\"name\":\"n\"}"}
	]}`)
	w := request("POST", "/api/traffic/ingest", body)
	if w.Code != 200 {
		t.Fatalf("ingest status = %d, body = %s", w.Code, w.Body)
	}
	if !strings.Contains(w.Body.String(), `"ingested":2`) || !strings.Contains(w.Body.String(), `"failed":0`) {
		t.Fatalf("unexpected counts: %s", w.Body)
	}

	// Both entries are visible on the normal traffic API, like proxy traffic.
	w = request("GET", "/api/traffic", "")
	if w.Code != 200 || !strings.Contains(w.Body.String(), "api.example.com") {
		t.Fatalf("traffic page after ingest: %d %s", w.Code, w.Body)
	}

	// One bad URL among good ones: the batch reports per-entry failure.
	mixed := `{"exchanges":[
		{"method":"GET","url":"https://good.example.com/ok","status":200,"responseBody":"ok"},
		{"method":"GET","url":"::::","status":200,"responseBody":"bad"}
	]}`
	w = request("POST", "/api/traffic/ingest", mixed)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"ingested":1`) || !strings.Contains(w.Body.String(), `"failed":1`) {
		t.Fatalf("mixed batch: %d %s", w.Code, w.Body)
	}

	// Guards: empty batch, oversized batch, oversized body truncation.
	if w = request("POST", "/api/traffic/ingest", `{"exchanges":[]}`); w.Code != 400 {
		t.Fatalf("empty batch status = %d", w.Code)
	}
	big := `{"exchanges":[` + strings.Repeat(`{"method":"GET","url":"https://x.example.com/"},`, maxIngestBatch) + `{}]}`
	if w = request("POST", "/api/traffic/ingest", big); w.Code != 400 {
		t.Fatalf("oversized batch status = %d", w.Code)
	}
	trunc := `{"exchanges":[{"method":"GET","url":"https://big.example.com/","status":200,"responseBody":"` + strings.Repeat("x", maxIngestBodySize+1024) + `"}]}`
	w = request("POST", "/api/traffic/ingest", trunc)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"ingested":1`) {
		t.Fatalf("truncating batch: %d %s", w.Code, w.Body)
	}
	if n, _ := s.m.Traffic().Count(); n != 4 { // 2 + 1 + 1
		t.Fatalf("final count = %d, want 4", n)
	}
}

func TestTrafficIngestRequiresAuth(t *testing.T) {
	s, _ := ingestServer(t)
	r := httptest.NewRequest("POST", "/api/traffic/ingest", strings.NewReader(`{"exchanges":[]}`))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, r)
	if w.Code != 401 {
		t.Fatalf("unauthenticated ingest status = %d, want 401", w.Code)
	}
}

func TestTrafficIngestCaptureGroup(t *testing.T) {
    s,request:=ingestServer(t)
    id:="00000000-0000-4000-8000-000000000099"
    body:=`{"exchanges":[{"captureId":"`+id+`","method":"GET","url":"https://browser.example.test/record","status":200,"responseBody":"captured"}]}`
    w:=request("POST","/api/traffic/ingest",body)
    if w.Code!=200||!strings.Contains(w.Body.String(),`"ingested":1`){t.Fatalf("ingest: %d %s",w.Code,w.Body)}
    rows,total,err:=s.m.Traffic().Page(traffic.PageQuery{Host:"browser.example.test",RespMin:-1,RespMax:-1},0,10)
    if err!=nil||total!=1||len(rows)!=1{t.Fatalf("rows: %+v %d %v",rows,total,err)}
    if err:=s.m.Traffic().ReadEvidence(context.Background(),[]string{rows[0].ID},func(e traffic.EvidenceExchange)error{
        if e.CaptureID!=id{t.Fatalf("capture ID lost on HTTP ingest: %q",e.CaptureID)}
        return nil
    });err!=nil{t.Fatal(err)}
}


func TestTrafficListCaptureGroupFilter(t *testing.T) {
    _,request:=ingestServer(t)
    body:=`{"exchanges":[
       {"captureId":"batch-alpha","method":"GET","url":"https://example.invalid/a","status":200},
       {"captureId":"batch-beta","method":"GET","url":"https://example.invalid/b","status":200},
       {"method":"GET","url":"https://example.invalid/legacy","status":200}
    ]}`
    if w:=request("POST","/api/traffic/ingest",body);w.Code!=200{
        t.Fatalf("ingest: %d %s",w.Code,w.Body)
    }
    for _,tc:=range []struct{query string; want int}{
        {"capture_id=batch-alpha",1},{"capture_id=batch-beta",1},
        {"capture_id=unknown",0},{"",3},
    }{
        url:="/api/traffic"
        if tc.query!=""{url+="?"+tc.query}
        w:=request("GET",url,"")
        if w.Code!=200||!strings.Contains(w.Body.String(),fmt.Sprintf(`"total":%d`,tc.want)){
            t.Fatalf("filter %q: %d %s",tc.query,w.Code,w.Body)
        }
    }
}

func TestTrafficCaptureGroupsAPI(t *testing.T) {
	_, request := ingestServer(t)
	payload := `{"exchanges":[
		{"captureId":"group-one","method":"GET","url":"https://one.example.invalid/a","status":200},
		{"captureId":"group-one","method":"GET","url":"https://two.example.invalid/b","status":404},
		{"captureId":"group-two","method":"GET","url":"https://two.example.invalid/c","status":0},
		{"method":"GET","url":"https://legacy.example.invalid/d","status":200}
	]}`
	if w := request("POST", "/api/traffic/ingest", payload); w.Code != 200 {
		t.Fatalf("ingest: %d %s", w.Code, w.Body)
	}
	w := request("GET", "/api/traffic/groups?size=10", "")
	if w.Code != 200 {
		t.Fatalf("groups API: %d %s", w.Code, w.Body)
	}
	var result struct {
		Groups  []traffic.CaptureGroupSummary `json:"groups"`
		HasMore bool                          `json:"has_more"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.HasMore || len(result.Groups) != 2 {
		t.Fatalf("groups response: %+v", result)
	}
	byID := map[string]traffic.CaptureGroupSummary{}
	for _, g := range result.Groups {
		byID[g.CaptureID] = g
	}
	if g := byID["group-one"]; g.Count != 2 || g.Hosts != 2 || g.ErrorResponses != 1 || g.UnknownStatuses != 0 {
		t.Fatalf("group one: %+v", g)
	}
	if g := byID["group-two"]; g.Count != 1 || g.UnknownStatuses != 1 {
		t.Fatalf("group two: %+v", g)
	}
	page := request("GET", "/api/traffic/groups?size=1&page=0", "")
	if page.Code != 200 || !strings.Contains(page.Body.String(), `"has_more":true`) {
		t.Fatalf("first group page: %d %s", page.Code, page.Body)
	}
	page = request("GET", "/api/traffic/groups?size=1&page=1", "")
	if page.Code != 200 || !strings.Contains(page.Body.String(), `"has_more":false`) {
		t.Fatalf("second group page: %d %s", page.Code, page.Body)
	}
}

func TestTrafficIngestReceiptDeduplicatesLostReply(t *testing.T) {
	s,request:=ingestServer(t)
	body:=`{"exchanges":[
		{"captureId":"group-retry","exchangeId":"request-one","method":"GET","url":"https://example.invalid/a","status":0},
		{"captureId":"group-retry","exchangeId":"request-two","method":"GET","url":"https://example.invalid/b","status":404}
	]}`
	check:=func(wantIngested,wantDuplicates int){
		t.Helper()
		w:=request("POST","/api/traffic/ingest",body)
		if w.Code!=200{t.Fatalf("ingest status=%d %s",w.Code,w.Body)}
		var result struct{Ingested int `json:"ingested"`;Duplicates int `json:"duplicates"`;Failed int `json:"failed"`}
		if err:=json.Unmarshal(w.Body.Bytes(),&result);err!=nil{t.Fatal(err)}
		if result.Ingested!=wantIngested||result.Duplicates!=wantDuplicates||result.Failed!=0{t.Fatalf("receipt: %+v",result)}
	}
	check(2,0)
	check(0,2)
	if n,err:=s.m.Traffic().Count();err!=nil||n!=2{t.Fatalf("retry duplicated evidence: %d %v",n,err)}
	groups,_,err:=s.m.Traffic().CaptureGroups(0,10)
	if err!=nil||len(groups)!=1||groups[0].Count!=2||groups[0].UnknownStatuses!=1||groups[0].ErrorResponses!=1{
		t.Fatalf("capture group after retry: %+v %v",groups,err)
	}
}
