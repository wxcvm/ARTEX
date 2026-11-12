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
 s:=&Server{m:&Manager{pg:d},jwtKey:[]byte("01234567890123456789012345678901")}
 mux:=http.NewServeMux()
 mux.HandleFunc("POST /api/companies/{id}/purge-intel",s.importPurgeIntel)
 mux.HandleFunc("GET /api/companies/{id}/purge-intel",s.listPurgeIntel)
 mux.HandleFunc("GET /api/companies/{id}/purge-intel/{snapshot}",s.getPurgeIntel)
 mux.HandleFunc("GET /api/companies/{id}/purge-intel/{snapshot}/analysis",s.getPurgeIntelAnalysis)
 mux.HandleFunc("POST /api/purge/intel-bridge",s.importPurgeIntelBridge)
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
 if err:=d.QueryRow("SELECT count(*) FROM assets WHERE company_id=$1",id).Scan(&count);err!=nil||count!=0{t.Fatalf("verified assets changed: %d %v",count,err)}
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
 if w=bridge(fmt.Sprint(id+1));w.Code!=409{t.Fatalf("mismatched company: %d",w.Code)}
 if w=bridge(fmt.Sprint(id));w.Code!=200{t.Fatalf("bridge import: %d %s",w.Code,w.Body.String())}
 if err:=d.QueryRow("SELECT count(*) FROM purge_intel_snapshots WHERE company_id=$1",id).Scan(&count);err!=nil||count!=1{t.Fatalf("bridge duplicated snapshot: %d %v",count,err)}
 b.IPs[0].IP="198.51.100.1"
 if w=request("POST",path,b);w.Code!=400{t.Fatal("out-of-scope import accepted")}
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
  {"POST","/api/companies/12/purge-intel",secret,"12"},
 } {if got:=req(tc.method,tc.path,tc.token,tc.company);got!=401{t.Fatalf("unexpected authorization %d",got)}}
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
