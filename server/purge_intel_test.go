package server

import (
 "encoding/json"
 "fmt"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "time"

 "github.com/Autumn-27/artex/db"
)

func fixturePurgeIntel() purgeIntelBundle {
 return purgeIntelBundle{OK:true,Schema:"purge.asset-intel.v1",Engagement:"lab",Target:"example.test",
  GeneratedAt:"2026-11-12T00:00:00Z",Mode:"offline-passive",
  Names:[]purgeIntelName{{Name:"a.example.test",Evidence:[]purgeIntelEvidence{{Source:"crt.sh"}},
   DNS:[]purgeIntelDNS{{Type:"A",Value:"192.0.2.10"}}}},
  IPs:[]purgeIntelIP{{IP:"192.0.2.10",PassivePorts:[]purgeIntelPort{{Port:443,Proto:"unknown",Evidence:[]purgeIntelEvidence{{Source:"internetdb"}}}}}},
  Edges:[]purgeIntelEdge{{From:"a.example.test",To:"192.0.2.10",Relation:"A"}},
 }
}
func fixturePurgeScope() []db.ScopeRule {
 return []db.ScopeRule{{Kind:"domain",Domain:"example.test"},{Kind:"cidr",Net:"192.0.2.0/24"}}
}
func TestPurgeIntelValidation(t *testing.T) {
 b:=fixturePurgeIntel()
 if err:=validatePurgeIntel(&b,fixturePurgeScope());err!=nil{t.Fatal(err)}
 for _,change:=range []func(*purgeIntelBundle){
  func(b *purgeIntelBundle){b.Names[0].Name="evil.test"},
  func(b *purgeIntelBundle){b.IPs[0].IP="198.51.100.1"},
  func(b *purgeIntelBundle){b.Names[0].DNS[0].Value="198.51.100.1"},
  func(b *purgeIntelBundle){b.Edges[0].To="198.51.100.1"},
  func(b *purgeIntelBundle){b.IPs[0].PassivePorts[0].Port=65536},
  func(b *purgeIntelBundle){b.IPs[0].PassivePorts[0].Proto="open"},
  func(b *purgeIntelBundle){b.Names[0].Evidence[0].Source="untrusted"},
  func(b *purgeIntelBundle){b.Mode="active"},
 } {
  x:=fixturePurgeIntel();change(&x)
  if validatePurgeIntel(&x,fixturePurgeScope())==nil{t.Fatal("accepted invalid evidence")}
 }
 if validatePurgeIntel(&b,[]db.ScopeRule{{Kind:"domain",Domain:"example.test"}})==nil{t.Fatal("domain authorization must not authorize IP")}
}
func TestPurgeIntelAuthAndJSON(t *testing.T) {
 s:=&Server{jwtKey:[]byte("01234567890123456789012345678901")}
 t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN","listener-test")
 t.Setenv("ARTEX_PRIVACY_OVERVIEW_TOKEN","overview-test")
 mux:=http.NewServeMux()
 mux.HandleFunc("POST /api/companies/{id}/purge-intel",s.importPurgeIntel)
 h:=s.requireAuth(mux)
 for _,header:=range []string{"","X-Artex-Privacy-Token","X-Artex-Privacy-Overview-Token"}{
  r:=httptest.NewRequest("POST","/api/companies/1/purge-intel",nil)
  if header!=""{r.Header.Set(header,map[string]string{"X-Artex-Privacy-Token":"listener-test","X-Artex-Privacy-Overview-Token":"overview-test"}[header])}
  w:=httptest.NewRecorder();h.ServeHTTP(w,r)
  if w.Code!=401{t.Fatalf("unexpected auth status %d",w.Code)}
 }
 for _,raw:=range []string{`{"unknown":1}`,`{} {}`} {
  if _,err:=decodePurgeIntel(httptest.NewRecorder(),httptest.NewRequest("POST","/",strings.NewReader(raw)));err==nil{t.Fatal("accepted malformed JSON")}
 }
}
func TestPurgeIntelDatabaseRoundTrip(t *testing.T) {
 dsn,_,err:=db.DSN();if err!=nil{t.Skip("postgres not configured")}
 d,err:=db.Open(dsn);if err!=nil{t.Skipf("postgres unavailable: %v",err)};defer d.Close()
 cs:=d.Companies()
 id,_,_,_,_,err:=cs.CreateCompanyWithScope(fmt.Sprintf("purge-intel-%d",time.Now().UnixNano()),"",[]db.ScopeInput{
  {Kind:"domain",Value:"example.test"},{Kind:"cidr",Value:"192.0.2.0/24"}},"test")
 if err!=nil{t.Fatal(err)}
 defer d.Exec("DELETE FROM companies WHERE id=$1",id)
 // Creating scope can attribute existing assets to the new company. Compare
 // the actual asset rows rather than assuming an initially empty inventory.
 assetSnapshot:=func() string {
  t.Helper()
  var snapshot string
  if err:=d.QueryRow("SELECT COALESCE(jsonb_agg(to_jsonb(a) ORDER BY a.id), '[]'::jsonb)::text FROM assets a WHERE a.company_id=$1",id).Scan(&snapshot);err!=nil{t.Fatal(err)}
  return snapshot
 }
 beforeAssets:=assetSnapshot()
 assertAssetsUnchanged:=func(){
  t.Helper()
  if after:=assetSnapshot();after!=beforeAssets{t.Fatal("verified assets changed during historical intel import")}
 }
 s:=&Server{m:&Manager{pg:d},jwtKey:[]byte("01234567890123456789012345678901")}
 mux:=http.NewServeMux()
 mux.HandleFunc("POST /api/companies/{id}/purge-intel",s.importPurgeIntel)
 mux.HandleFunc("GET /api/companies/{id}/purge-intel",s.listPurgeIntel)
 mux.HandleFunc("GET /api/companies/{id}/purge-intel/{snapshot}",s.getPurgeIntel)
 mux.HandleFunc("GET /api/companies/{id}/purge-intel/{snapshot}/analysis",s.getPurgeIntelAnalysis)
 mux.HandleFunc("POST /api/purge/intel-bridge",s.importPurgeIntelBridge)
 mux.HandleFunc("GET /api/purge/intel-bridge/status",s.purgeIntelBridgeStatus)
 mux.HandleFunc("GET /api/purge/intel-bridge/overview",s.purgeIntelBridgeOverview)
 mux.HandleFunc("GET /api/purge/intel-bridge/assets",s.purgeIntelBridgeAssets)
 mux.HandleFunc("GET /api/purge/intel-bridge/snapshot",s.purgeIntelBridgeSnapshot)
 mux.HandleFunc("GET /api/purge/intel-bridge/task",s.purgeIntelBridgeTask)
 mux.HandleFunc("GET /api/purge/intel-bridge/findings",s.purgeIntelBridgeFindings)
 mux.HandleFunc("GET /api/purge/intel-bridge/report",s.purgeIntelBridgeReport)
 mux.HandleFunc("GET /api/purge/intel-pairing",s.managePurgeIntelPairing)
 mux.HandleFunc("POST /api/purge/intel-pairing",s.managePurgeIntelPairing)
 mux.HandleFunc("DELETE /api/purge/intel-pairing",s.managePurgeIntelPairing)
 h:=s.requireAuth(mux)
 token,err:=signJWT(s.jwtKey);if err!=nil{t.Fatal(err)}
 path:=fmt.Sprintf("/api/companies/%d/purge-intel",id)
 request:=func(method,path string,b purgeIntelBundle)*httptest.ResponseRecorder{
  payload,_:=json.Marshal(b)
  r:=httptest.NewRequest(method,path,strings.NewReader(string(payload)))
  r.Header.Set("Authorization","Bearer "+token)
  w:=httptest.NewRecorder();h.ServeHTTP(w,r);return w
 }
 b:=fixturePurgeIntel()
 w:=request("POST",path,b);if w.Code!=200{t.Fatalf("import: %d %s",w.Code,w.Body.String())}
 var first struct{ID int64 `json:"id"`};if json.Unmarshal(w.Body.Bytes(),&first)!=nil||first.ID==0{t.Fatal("missing snapshot ID")}
 b.GeneratedAt="2026-11-13T00:00:00Z"
 w=request("POST",path,b);var second struct{ID int64 `json:"id"`};_ = json.Unmarshal(w.Body.Bytes(),&second)
 if w.Code!=200||second.ID!=first.ID{t.Fatal("duplicate snapshot not merged")}
 w=request("GET",fmt.Sprintf("%s/%d",path,first.ID),b);if w.Code!=200{t.Fatalf("read: %d %s",w.Code,w.Body.String())}
 w=request("GET",path,b);if w.Code!=200{t.Fatalf("list: %d",w.Code)}
 w=request("GET",fmt.Sprintf("%s/%d/analysis",path,first.ID),b)
 if w.Code!=200{t.Fatalf("analysis: %d %s",w.Code,w.Body.String())}
 var analysis purgeIntelAnalysis
 if json.Unmarshal(w.Body.Bytes(),&analysis)!=nil||analysis.PortObservations!=1||analysis.UnknownTransportPorts!=1||!analysis.HistoricalOnly{t.Fatalf("invalid analysis: %+v",analysis)}

 var count int
 if err:=d.QueryRow("SELECT count(*) FROM purge_intel_snapshots WHERE company_id=$1",id).Scan(&count);err!=nil||count!=1{t.Fatalf("snapshot count %d: %v",count,err)}
 assertAssetsUnchanged()
 // The direct bridge reuses the validated import path and the same snapshot.
 const secret="fixture-only-long-intelligence-token-123"
 t.Setenv("ARTEX_PURGE_INTEL_TOKEN",secret)
 t.Setenv("ARTEX_PURGE_INTEL_COMPANY_ID",fmt.Sprint(id))
 bridge:=func(company string)*httptest.ResponseRecorder{
  payload,_:=json.Marshal(b)
  r:=httptest.NewRequest("POST","/api/purge/intel-bridge",strings.NewReader(string(payload)))
  r.Header.Set("X-Artex-Purge-Intel-Token",secret)
  r.Header.Set("X-Artex-Purge-Intel-Company",company)
  w:=httptest.NewRecorder();h.ServeHTTP(w,r);return w
 }
 statusReq:=httptest.NewRequest("GET","/api/purge/intel-bridge/status",nil)
 statusReq.Header.Set("X-Artex-Purge-Intel-Token",secret)
 statusReq.Header.Set("X-Artex-Purge-Intel-Company",fmt.Sprint(id))
 statusReply:=httptest.NewRecorder();h.ServeHTTP(statusReply,statusReq)
 if statusReply.Code!=200{t.Fatalf("pair status: %d %s",statusReply.Code,statusReply.Body.String())}
 var pairing struct{ CompanyID int64 `json:"company_id"`; Paired bool `json:"paired"` }
 if json.Unmarshal(statusReply.Body.Bytes(),&pairing)!=nil||pairing.CompanyID!=id||!pairing.Paired{t.Fatal("invalid pair status")}
 if err:=d.QueryRow("SELECT count(*) FROM purge_intel_snapshots WHERE company_id=$1",id).Scan(&count);err!=nil||count!=1{t.Fatal("status modified snapshots")}
 if w=bridge(fmt.Sprint(id+1));w.Code!=409{t.Fatalf("mismatched company: %d",w.Code)}
 if w=bridge(fmt.Sprint(id));w.Code!=200{t.Fatalf("bridge import: %d %s",w.Code,w.Body.String())}
 t.Setenv("ARTEX_PURGE_INTEL_TOKEN","")
 t.Setenv("ARTEX_PURGE_INTEL_COMPANY_ID","")
 admin:=func(method,body string)*httptest.ResponseRecorder{
  r:=httptest.NewRequest(method,"/api/purge/intel-pairing",strings.NewReader(body))
  r.Header.Set("Authorization","Bearer "+token)
  w:=httptest.NewRecorder();h.ServeHTTP(w,r);return w
 }
 // Preserve test DB state; this integration fixture uses its own database.
 t.Cleanup(func(){_ = d.SetSetting(purgeIntelPairingKey,"")})
 created:=admin("POST",fmt.Sprintf(`{"company_id":%d}`,id))
 var pairResponse struct{ Token string `json:"token"`; CompanyID int64 `json:"company_id"` }
 if created.Code!=200||json.Unmarshal(created.Body.Bytes(),&pairResponse)!=nil||len(pairResponse.Token)!=64||pairResponse.CompanyID!=id{t.Fatalf("create pairing failed: %d",created.Code)}
 stored,_,err:=d.GetSetting(purgeIntelPairingKey);if err!=nil||strings.Contains(stored,pairResponse.Token){t.Fatal("plaintext pairing stored")}
 if got:=admin("GET","");got.Code!=200||strings.Contains(got.Body.String(),pairResponse.Token)||strings.Contains(got.Body.String(),"digest"){t.Fatal("status exposed secret")}
 if gotID,ok:=s.authorizePurgeIntelPairing(pairResponse.Token);!ok||gotID!=id{t.Fatal("saved pairing not accepted")}
 overviewReq:=httptest.NewRequest("GET","/api/purge/intel-bridge/overview",nil)
 overviewReq.Header.Set("X-Artex-Purge-Intel-Token",pairResponse.Token)
 overviewReq.Header.Set("X-Artex-Purge-Intel-Company",fmt.Sprint(id))
 overviewReply:=httptest.NewRecorder();h.ServeHTTP(overviewReply,overviewReq)
 var overview struct{ CompanyID int64 `json:"company_id"`; SnapshotCount int64 `json:"snapshot_count"`; Snapshots []struct{ID int64 `json:"id"`} `json:"snapshots"` }
 if overviewReply.Code!=200||json.Unmarshal(overviewReply.Body.Bytes(),&overview)!=nil||overview.CompanyID!=id||overview.SnapshotCount!=1||len(overview.Snapshots)!=1||overview.Snapshots[0].ID!=first.ID{t.Fatalf("overview failed: %d %s",overviewReply.Code,overviewReply.Body.String())}
 bridgeRead:=func(path,company string)*httptest.ResponseRecorder{
  r:=httptest.NewRequest("GET",path,nil);r.Header.Set("X-Artex-Purge-Intel-Token",pairResponse.Token);r.Header.Set("X-Artex-Purge-Intel-Company",company);w:=httptest.NewRecorder();h.ServeHTTP(w,r);return w
 }
 assetsReply:=bridgeRead("/api/purge/intel-bridge/assets",fmt.Sprint(id))
 var assetPage struct{CompanyID int64 `json:"company_id"`;Items []map[string]any `json:"items"`}
 if assetsReply.Code!=200||json.Unmarshal(assetsReply.Body.Bytes(),&assetPage)!=nil||assetPage.CompanyID!=id||assetPage.Items==nil{t.Fatalf("asset page failed: %d",assetsReply.Code)}
 for _,a:=range assetPage.Items{for _,key:=range []string{"auth","params","extra","url"}{if _,exists:=a[key];exists{t.Fatal("private asset field exposed")}}}
 if got:=bridgeRead("/api/purge/intel-bridge/assets?before=-1",fmt.Sprint(id));got.Code!=400{t.Fatal("bad cursor accepted")}
 if got:=bridgeRead("/api/purge/intel-bridge/assets",fmt.Sprint(id+1));got.Code!=409{t.Fatal("other company accepted")}
 snap:=bridgeRead(fmt.Sprintf("/api/purge/intel-bridge/snapshot?snapshot=%d",first.ID),fmt.Sprint(id))
 var snapBody struct{CompanyID int64 `json:"company_id"`;SnapshotID int64 `json:"snapshot_id"`;Analysis purgeIntelAnalysis `json:"analysis"`}
 if snap.Code!=200||json.Unmarshal(snap.Body.Bytes(),&snapBody)!=nil||snapBody.CompanyID!=id||snapBody.SnapshotID!=first.ID||!snapBody.Analysis.HistoricalOnly||snapBody.Analysis.PortObservations!=1{t.Fatalf("snapshot read failed: %d",snap.Code)}
 if got:=bridgeRead("/api/purge/intel-bridge/snapshot?snapshot=0",fmt.Sprint(id));got.Code!=400{t.Fatal("invalid snapshot accepted")}
 assertAssetsUnchanged()
 // Isolated task/finding fixture with explicit company attribution.
 task,err:=d.CreateTaskWithOptions("bridge detail","goal",db.TaskCreateOptions{Name:"paired task",CompanyIDs:[]int64{id}})
 if err!=nil{t.Fatal(err)}
 defer d.DeleteTaskCascade(task.ID,true,true)
 aid,err:=d.Assets().UpsertRootDomain(db.UpsertRootDomainReq{Domain:fmt.Sprintf("bridge-%d.example.test",time.Now().UnixNano()),TaskID:task.ID})
 if err!=nil{t.Fatal(err)}
 if _,err=d.Exec("UPDATE assets SET company_id=$1 WHERE id=$2",id,aid);err!=nil{t.Fatal(err)}
 fid,err:=d.AddFinding(task.ID,0,"test","paired finding","low","summary","evidence","worker",[]int64{aid})
 if err!=nil{t.Fatal(err)}
 if _,err=d.Exec("UPDATE findings SET report=$1 WHERE id=$2","fixture report",fid);err!=nil{t.Fatal(err)}
 // Task creation deliberately mutates the verified asset database. Use its
 // own baseline to assert that the bridge detail reads are read-only.
 beforeDetailAssets:=assetSnapshot()
 td:=bridgeRead(fmt.Sprintf("/api/purge/intel-bridge/task?task=%d",task.ID),fmt.Sprint(id))
 if td.Code!=200||!strings.Contains(td.Body.String(),"paired task"){t.Fatalf("task detail failed %d %s",td.Code,td.Body.String())}
 fl:=bridgeRead(fmt.Sprintf("/api/purge/intel-bridge/findings?task=%d",task.ID),fmt.Sprint(id))
 if fl.Code!=200||!strings.Contains(fl.Body.String(),"paired finding"){t.Fatalf("findings failed %d %s",fl.Code,fl.Body.String())}
 fr:=bridgeRead(fmt.Sprintf("/api/purge/intel-bridge/report?finding=%d",fid),fmt.Sprint(id))
 if fr.Code!=200||!strings.Contains(fr.Body.String(),"fixture report"){t.Fatalf("report failed %d %s",fr.Code,fr.Body.String())}
 if got:=bridgeRead(fmt.Sprintf("/api/purge/intel-bridge/report?finding=%d",fid),fmt.Sprint(id+1));got.Code!=409{t.Fatal("other company accepted")}
 if after:=assetSnapshot();after!=beforeDetailAssets{t.Fatal("bridge detail reads changed verified assets")}
 if _,err=d.Exec("UPDATE assets SET company_id=NULL WHERE id=$1",aid);err!=nil{t.Fatal(err)}
 if got:=bridgeRead(fmt.Sprintf("/api/purge/intel-bridge/report?finding=%d",fid),fmt.Sprint(id));got.Code!=404{t.Fatal("unattributed report accepted")}
 if _,err=d.DeleteTaskCascade(task.ID,true,true);err!=nil{t.Fatal(err)}
 // Fixture creation, reassignment and cascade cleanup are intentional asset
 // writes. Compare the rejected import against the post-cleanup inventory,
 // rather than the snapshot before those writes (which also tracks timestamps).
 beforeRejectedAssets:=assetSnapshot()
 if got:=admin("DELETE","");got.Code!=200{t.Fatalf("revoke failed: %d",got.Code)}
 if _,ok:=s.authorizePurgeIntelPairing(pairResponse.Token);ok{t.Fatal("revoked pairing accepted")}
 if err:=d.QueryRow("SELECT count(*) FROM purge_intel_snapshots WHERE company_id=$1",id).Scan(&count);err!=nil||count!=1{t.Fatalf("bridge duplicated snapshot: %d %v",count,err)}
 b.IPs[0].IP="198.51.100.1"
 if w=request("POST",path,b);w.Code!=400{t.Fatal("out-of-scope import accepted")}
 if after:=assetSnapshot();after!=beforeRejectedAssets{t.Fatal("rejected import or pairing revoke changed verified assets")}
}

func TestPurgeIntelBridgeCredentialIsolation(t *testing.T) {
 const secret="fixture-only-long-intelligence-token-123"
 t.Setenv("ARTEX_PURGE_INTEL_TOKEN",secret)
 t.Setenv("ARTEX_PURGE_INTEL_COMPANY_ID","12")
 t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN","listener-test")
 t.Setenv("ARTEX_PRIVACY_OVERVIEW_TOKEN","overview-test")
 t.Setenv("ARTEX_SETUP_TOKEN","setup-test")
 s:=&Server{jwtKey:[]byte("01234567890123456789012345678901")}
 mux:=http.NewServeMux()
 mux.HandleFunc("POST /api/purge/intel-bridge",s.importPurgeIntelBridge)
 mux.HandleFunc("GET /api/purge/intel-bridge/status",s.purgeIntelBridgeStatus)
 mux.HandleFunc("GET /api/purge/intel-bridge/overview",s.purgeIntelBridgeOverview)
 mux.HandleFunc("GET /api/purge/intel-bridge/assets",s.purgeIntelBridgeAssets)
 mux.HandleFunc("GET /api/purge/intel-bridge/snapshot",s.purgeIntelBridgeSnapshot)
 mux.HandleFunc("GET /api/purge/intel-bridge/task",s.purgeIntelBridgeTask)
 mux.HandleFunc("GET /api/purge/intel-bridge/findings",s.purgeIntelBridgeFindings)
 mux.HandleFunc("GET /api/purge/intel-bridge/report",s.purgeIntelBridgeReport)
 mux.HandleFunc("GET /api/purge/intel-pairing",s.managePurgeIntelPairing)
 mux.HandleFunc("POST /api/purge/intel-pairing",s.managePurgeIntelPairing)
 mux.HandleFunc("DELETE /api/purge/intel-pairing",s.managePurgeIntelPairing)
 mux.HandleFunc("POST /api/companies/{id}/purge-intel",s.importPurgeIntel)
 h:=s.requireAuth(mux)
 req:=func(method,path,token,company string)int{
  r:=httptest.NewRequest(method,path,nil)
  r.Header.Set("X-Artex-Purge-Intel-Token",token)
  r.Header.Set("X-Artex-Purge-Intel-Company",company)
  w:=httptest.NewRecorder();h.ServeHTTP(w,r);return w.Code
 }
 for _,tc:=range []struct{method,path,token,company string}{
  {"POST","/api/purge/intel-bridge","","12"},
  {"POST","/api/purge/intel-bridge","listener-test","12"},
  {"POST","/api/purge/intel-bridge","overview-test","12"},
  {"GET","/api/purge/intel-bridge",secret,"12"},
  {"GET","/api/purge/intel-bridge/status","listener-test","12"},
  {"GET","/api/purge/intel-bridge/status","","12"},
  {"GET","/api/purge/intel-bridge/overview","","12"},
  {"GET","/api/purge/intel-bridge/assets","","12"},
  {"GET","/api/purge/intel-bridge/task","","12"},
  {"GET","/api/purge/intel-bridge/findings","listener-test","12"},
  {"GET","/api/purge/intel-bridge/report","overview-test","12"},
  {"POST","/api/purge/intel-bridge/report",secret,"12"},
  {"GET","/api/purge/intel-bridge/snapshot","listener-test","12"},
  {"POST","/api/purge/intel-bridge/assets",secret,"12"},
  {"GET","/api/purge/intel-bridge/overview","listener-test","12"},
  {"POST","/api/purge/intel-bridge/overview",secret,"12"},
  {"POST","/api/purge/intel-bridge/status",secret,"12"},
  {"GET","/api/purge/intel-pairing",secret,"12"},
  {"POST","/api/purge/intel-pairing",secret,"12"},
  {"DELETE","/api/purge/intel-pairing",secret,"12"},
  {"POST","/api/companies/12/purge-intel",secret,"12"},
 } {if got:=req(tc.method,tc.path,tc.token,tc.company);got!=401{t.Fatalf("unexpected authorization %d",got)}}
 if got:=req("GET","/api/purge/intel-bridge/overview",secret,"13");got!=409{t.Fatalf("overview company mismatch %d",got)}
 if got:=req("GET","/api/purge/intel-bridge/status",secret,"13");got!=409{t.Fatalf("status company mismatch %d",got)}
 if got:=req("GET","/api/purge/intel-bridge/status",secret,"12");got!=503{t.Fatalf("status expected DB unavailable, got %d",got)}
 if got:=req("POST","/api/purge/intel-bridge",secret,"13");got!=409{t.Fatalf("company mismatch %d",got)}
 if got:=req("POST","/api/purge/intel-bridge",secret,"12");got!=503{t.Fatalf("expected DB unavailable, got %d",got)}
 t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN",secret)
 if validPurgeIntelBridgeCredential(secret){t.Fatal("reused read-only token accepted")}
 t.Setenv("ARTEX_PRIVACY_BRIDGE_TOKEN","listener-test")
 t.Setenv("ARTEX_PURGE_INTEL_COMPANY_ID","")
 if validPurgeIntelBridgeCredential(secret){t.Fatal("unbound token accepted")}
}

func TestPurgeIntelAnalysisEvidenceGaps(t *testing.T) {
 b:=fixturePurgeIntel()
 timestamp:="2026-11-12T00:00:00Z"
 b.Names[0].DNSCollectedAt=&timestamp
 b.Names[0].Evidence=append(b.Names[0].Evidence,purgeIntelEvidence{Source:"dns"})
 b.Names=append(b.Names,purgeIntelName{Name:"b.example.test",Evidence:[]purgeIntelEvidence{{Source:"registered-scope"}}})
 b.IPs=append(b.IPs,purgeIntelIP{IP:"192.0.2.11"})
 b.Truncated=true
 b.Excluded=map[string]int{"dns_outside_scope":3}
 a:=analyzePurgeIntel(b)
 if !a.HistoricalOnly||!a.Truncated||a.ExcludedDNS!=3||a.Names!=2||a.IPs!=2||a.PortObservations!=1||a.UnknownTransportPorts!=1||a.PortsWithoutObservationTime!=1{t.Fatalf("invalid counts: %+v",a)}
 if len(a.NamesWithMultipleSources)!=1||a.NamesWithMultipleSources[0]!="a.example.test"||len(a.NamesWithoutDNS)!=1||a.NamesWithoutDNS[0]!="b.example.test" {t.Fatalf("invalid name gaps: %+v",a)}
 if len(a.NamesWithoutRetainedAddress)!=1||len(a.IPsWithoutDomainLink)!=1||len(a.IPsWithoutPortEvidence)!=1||a.SourceCounts["registered-scope"]!=0{t.Fatalf("invalid evidence gaps: %+v",a)}
 b.Names[1].DNS=[]purgeIntelDNS{{Type:"A",Value:"192.0.2.10"}}
 b.Edges=append(b.Edges,purgeIntelEdge{From:"b.example.test",To:"192.0.2.10",Relation:"A"})
 if got:=analyzePurgeIntel(b);got.SharedAddresses!=1{t.Fatalf("shared address count=%d",got.SharedAddresses)}
}
