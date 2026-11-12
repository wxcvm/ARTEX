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
 var count int
 if err:=d.QueryRow("SELECT count(*) FROM purge_intel_snapshots WHERE company_id=$1",id).Scan(&count);err!=nil||count!=1{t.Fatalf("snapshot count %d: %v",count,err)}
 if err:=d.QueryRow("SELECT count(*) FROM assets WHERE company_id=$1",id).Scan(&count);err!=nil||count!=0{t.Fatalf("verified assets changed: %d %v",count,err)}
 b.IPs[0].IP="198.51.100.1"
 if w=request("POST",path,b);w.Code!=400{t.Fatal("out-of-scope import accepted")}
}
