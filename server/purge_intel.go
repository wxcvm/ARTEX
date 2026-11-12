package server

import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "net/netip"
 "regexp"
 "strconv"
 "strings"
 "time"

 "github.com/Autumn-27/artex/db"
)

const purgeIntelMaxBody = 8 << 20

type purgeIntelEvidence struct {
 Source string `json:"source"`
 CollectedAt *string `json:"collected_at"`
 ObservedAt *string `json:"observed_at,omitempty"`
}
type purgeIntelDNS struct {
 Type string `json:"type"`
 Value string `json:"value"`
 TTL *int64 `json:"ttl,omitempty"`
}
type purgeIntelName struct {
 Name string `json:"name"`
 Evidence []purgeIntelEvidence `json:"evidence"`
 DNS []purgeIntelDNS `json:"dns"`
 DNSCollectedAt *string `json:"dns_collected_at,omitempty"`
 DNSErrors []struct { Type string `json:"type"`; Code string `json:"code"` } `json:"dns_errors,omitempty"`
}
type purgeIntelPort struct {
 Port int `json:"port"`
 Proto string `json:"proto"`
 Evidence []purgeIntelEvidence `json:"evidence"`
}
type purgeIntelIP struct {
 IP string `json:"ip"`
 PassivePorts []purgeIntelPort `json:"passive_ports"`
}
type purgeIntelEdge struct {
 From string `json:"from"`
 To string `json:"to"`
 Relation string `json:"relation"`
 CollectedAt *string `json:"collected_at"`
}
type purgeIntelBundle struct {
 OK bool `json:"ok"`
 Schema string `json:"schema"`
 Engagement string `json:"engagement"`
 Target string `json:"target"`
 GeneratedAt string `json:"generated_at"`
 Mode string `json:"mode"`
 Names []purgeIntelName `json:"names"`
 IPs []purgeIntelIP `json:"ips"`
 Edges []purgeIntelEdge `json:"edges"`
 Excluded map[string]int `json:"excluded"`
 Truncated bool `json:"truncated"`
 Note string `json:"note"`
}

var purgeIntelLabel = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)
func validPurgeIntelName(s string) bool {
 if len(s)>253 || !strings.Contains(s,".") {return false}
 for _, label := range strings.Split(s,".") { if !purgeIntelLabel.MatchString(label) {return false} }
 return true
}
func purgeIntelNameAllowed(name string, scope []db.ScopeRule) bool {
 if !validPurgeIntelName(name) {return false}
 for _, rule := range scope {
  if rule.Kind!="domain" {continue}
  d:=strings.TrimSuffix(strings.ToLower(rule.Domain),".")
  if name==d || strings.HasSuffix(name,"."+d) {return true}
 }
 return false
}
func purgeIntelIPAllowed(ip string, scope []db.ScopeRule) bool {
 addr,err:=netip.ParseAddr(ip)
 if err!=nil || addr.Zone()!="" {return false}
 for _,rule:=range scope {
  if rule.Kind!="ip" && rule.Kind!="cidr" {continue}
  prefix,err:=netip.ParsePrefix(rule.Net)
  if err==nil && prefix.Contains(addr) {return true}
 }
 return false
}
func validIntelTime(s *string) bool {
 if s==nil {return true}
 if len(*s)>64 {return false}
 // Provider observations can use a timestamp without a UTC offset.
 for _,layout:=range []string{time.RFC3339Nano,"2006-01-02T15:04:05.999999999","2006-01-02"} {
  if _,err:=time.Parse(layout,*s);err==nil{return true}
 }
 return false
}
func validIntelEvidence(e []purgeIntelEvidence, port bool) bool {
 if len(e)>4 {return false}
 for _,item:=range e {
  allowed:=item.Source=="internetdb" || item.Source=="shodan"
  if !port {allowed=item.Source=="crt.sh" || item.Source=="subfinder" || item.Source=="dns" || item.Source=="registered-scope"}
  if !allowed || !validIntelTime(item.CollectedAt) || !validIntelTime(item.ObservedAt) {return false}
 }
 return true
}
func validatePurgeIntel(b *purgeIntelBundle, scope []db.ScopeRule) error {
 if !b.OK || b.Schema!="purge.asset-intel.v1" || b.Mode!="offline-passive" || b.Engagement=="" || len(b.Engagement)>256 ||
 !purgeIntelNameAllowed(b.Target,scope) || !validIntelTime(&b.GeneratedAt) || len(b.Names)>500 || len(b.IPs)>500 || len(b.Edges)>16000 {
  return errors.New("invalid inventory envelope or target outside company scope")
 }
 names:=map[string]bool{}
 ips:=map[string]bool{}
 for _,n:=range b.Names {
  if names[n.Name] || !purgeIntelNameAllowed(n.Name,scope) || (n.Name!=b.Target && !strings.HasSuffix(n.Name,"."+b.Target)) ||
   !validIntelEvidence(n.Evidence,false) || len(n.DNS)>32 || len(n.DNSErrors)>3 || !validIntelTime(n.DNSCollectedAt) {
   return errors.New("invalid or out-of-scope hostname evidence")
  }
  names[n.Name]=true
  for _,d:=range n.DNSErrors {
   if d.Type!="A" && d.Type!="AAAA" && d.Type!="CNAME" {return errors.New("invalid DNS error type")}
   switch d.Code {case "ENODATA","ENOTFOUND","ETIMEOUT","ESERVFAIL","EREFUSED","DNS_ERROR": default:return errors.New("invalid DNS error code")}
  }
 }
 for _,ip:=range b.IPs {
  if ips[ip.IP] || !purgeIntelIPAllowed(ip.IP,scope) || len(ip.PassivePorts)>64 {return errors.New("invalid or out-of-scope IP evidence")}
  ips[ip.IP]=true
  seen:=map[string]bool{}
  for _,p:=range ip.PassivePorts {
   key:=fmt.Sprintf("%d/%s",p.Port,p.Proto)
   if seen[key] || p.Port<1 || p.Port>65535 || (p.Proto!="tcp" && p.Proto!="udp" && p.Proto!="unknown") || !validIntelEvidence(p.Evidence,true) {return errors.New("invalid historical port")}
   seen[key]=true
  }
 }
 dnsEdges:=map[string]bool{}
 for _,n:=range b.Names {
  for _,record:=range n.DNS {
   if record.TTL!=nil && (*record.TTL<0 || *record.TTL>4294967295) {return errors.New("invalid DNS TTL")}
   switch record.Type {
   case "A","AAAA":
    a,err:=netip.ParseAddr(record.Value)
    if err!=nil || !ips[record.Value] || (record.Type=="A")!=a.Is4() {return errors.New("invalid or out-of-scope DNS address")}
   case "CNAME":
    if !purgeIntelNameAllowed(record.Value,scope) {return errors.New("out-of-scope DNS alias")}
   default:return errors.New("invalid DNS record")
   }
   dnsEdges[n.Name+"|"+record.Type+"|"+record.Value]=true
  }
 }
 seenEdges:=map[string]bool{}
 for _,edge:=range b.Edges {
  key:=edge.From+"|"+edge.Relation+"|"+edge.To
  if seenEdges[key] || !names[edge.From] || !dnsEdges[key] || !validIntelTime(edge.CollectedAt) ||
    (edge.Relation=="CNAME" && !names[edge.To]) || (edge.Relation!="CNAME" && !ips[edge.To]) {return errors.New("invalid inventory relationship")}
  seenEdges[key]=true
 }
 if len(b.Excluded)>8 {return errors.New("invalid excluded counts")}
 for k,v:=range b.Excluded {
  switch k {case "dns_outside_scope","names_over_limit","ips_over_limit","dns_over_limit","ports_over_limit":default:return errors.New("invalid excluded count")}
  if v<0 {return errors.New("invalid excluded count")}
 }
 // Free text from the source is not executable instructions or trusted evidence.
 b.Note="Historical passive evidence; not verified assets, live ports or findings."
 return nil
}

func decodePurgeIntel(w http.ResponseWriter,r *http.Request)(purgeIntelBundle,error){
 var b purgeIntelBundle
 r.Body=http.MaxBytesReader(w,r.Body,purgeIntelMaxBody)
 d:=json.NewDecoder(r.Body);d.DisallowUnknownFields()
 if err:=d.Decode(&b);err!=nil{return b,err}
 if err:=d.Decode(new(any));err!=io.EOF{return b,errors.New("trailing JSON")}
 return b,nil
}

func (s *Server) purgeIntelCompany(w http.ResponseWriter,r *http.Request)(int64,[]db.ScopeRule,bool){
 id,err:=strconv.ParseInt(r.PathValue("id"),10,64)
 if err!=nil || id<1 {writeErr(w,400,"invalid company id");return 0,nil,false}
 if s.m==nil || s.m.pg==nil {writeErr(w,503,"database unavailable");return 0,nil,false}
 cs:=s.m.pg.Companies()
 company,err:=cs.GetCompany(id)
 if err!=nil {writeErr(w,503,"company lookup unavailable");return 0,nil,false}
 if company==nil {writeErr(w,404,"company not found");return 0,nil,false}
 scope,err:=cs.GetScope(id)
 if err!=nil {writeErr(w,503,"company scope unavailable");return 0,nil,false}
 return id,scope,true
}

func (s *Server) importPurgeIntel(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Cache-Control","no-store")
 id,scope,ok:=s.purgeIntelCompany(w,r);if !ok{return}
 b,err:=decodePurgeIntel(w,r)
 if err!=nil {writeErr(w,400,"invalid inventory JSON");return}
 if err=validatePurgeIntel(&b,scope);err!=nil {writeErr(w,400,err.Error());return}
 payload,err:=json.Marshal(b);if err!=nil {writeErr(w,400,"invalid inventory");return}
 // Ignore only export time when deduplicating the same evidence.
 b.GeneratedAt=""
 key,err:=json.Marshal(b);if err!=nil {writeErr(w,400,"invalid inventory");return}
 sum:=sha256.Sum256(key);digest:=hex.EncodeToString(sum[:])
 var snapshotID int64
 err=s.m.pg.QueryRowContext(r.Context(),`INSERT INTO purge_intel_snapshots(company_id,digest,payload) VALUES($1,$2,$3::jsonb)
 ON CONFLICT(company_id,digest) DO UPDATE SET received_at=now() RETURNING id`,id,digest,string(payload)).Scan(&snapshotID)
 if err!=nil {writeErr(w,503,"inventory save failed");return}
 writeJSON(w,200,map[string]any{"id":snapshotID,"company_id":id,"digest":digest,"names":len(b.Names),"ips":len(b.IPs),"historical_only":true})
}

func (s *Server) listPurgeIntel(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Cache-Control","no-store")
 id,_,ok:=s.purgeIntelCompany(w,r);if !ok{return}
 rows,err:=s.m.pg.QueryContext(r.Context(),`SELECT id,digest,received_at::text FROM purge_intel_snapshots WHERE company_id=$1 ORDER BY received_at DESC,id DESC LIMIT 50`,id)
 if err!=nil {writeErr(w,503,"inventory read failed");return};defer rows.Close()
 type item struct{ID int64 `json:"id"`;Digest string `json:"digest"`;ReceivedAt string `json:"received_at"`}
 items:=[]item{}
 for rows.Next(){var x item;if rows.Scan(&x.ID,&x.Digest,&x.ReceivedAt)!=nil{writeErr(w,503,"inventory read failed");return};items=append(items,x)}
 if rows.Err()!=nil{writeErr(w,503,"inventory read failed");return}
 writeJSON(w,200,map[string]any{"items":items,"historical_only":true})
}

func (s *Server) getPurgeIntel(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Cache-Control","no-store")
 id,scope,ok:=s.purgeIntelCompany(w,r);if !ok{return}
 sid,err:=strconv.ParseInt(r.PathValue("snapshot"),10,64);if err!=nil||sid<1{writeErr(w,400,"invalid snapshot id");return}
 var payload []byte
 err=s.m.pg.QueryRowContext(r.Context(),`SELECT payload FROM purge_intel_snapshots WHERE company_id=$1 AND id=$2`,id,sid).Scan(&payload)
 if err!=nil{writeErr(w,404,"inventory unavailable");return}
 var b purgeIntelBundle
 if json.Unmarshal(payload,&b)!=nil || validatePurgeIntel(&b,scope)!=nil {writeErr(w,409,"inventory no longer matches company scope");return}
 writeJSON(w,200,b)
}
