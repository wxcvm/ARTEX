package server

import (
 "crypto/rand"
 "crypto/sha256"
 "crypto/subtle"
 "encoding/hex"
 "encoding/json"
 "net/http"
 "os"
 "strconv"
 "strings"
)

func purgeIntelBridgeCompany() (int64,bool) {
 id,err:=strconv.ParseInt(strings.TrimSpace(os.Getenv("ARTEX_PURGE_INTEL_COMPANY_ID")),10,64)
 return id,err==nil&&id>0
}

func validPurgeIntelBridgeCredential(provided string) bool {
 expected:=strings.TrimSpace(os.Getenv("ARTEX_PURGE_INTEL_TOKEN"))
 if len(expected)<32 || !validSetupToken(provided,expected) {return false}
 if _,ok:=purgeIntelBridgeCompany();!ok{return false}
 // Never allow a previously issued read-only token to become a write token.
 for _,key:=range []string{"ARTEX_PRIVACY_BRIDGE_TOKEN","ARTEX_PRIVACY_OVERVIEW_TOKEN","ARTEX_SETUP_TOKEN"} {
  if expected==strings.TrimSpace(os.Getenv(key)){return false}
 }
 return true
}

const purgeIntelPairingKey = "purge.intel_pairing.v1"

type purgeIntelPairing struct {
 CompanyID int64 `json:"company_id"`
 Digest string `json:"digest"`
}

func purgeIntelEnvOverride() bool {
 return strings.TrimSpace(os.Getenv("ARTEX_PURGE_INTEL_TOKEN")) != "" ||
  strings.TrimSpace(os.Getenv("ARTEX_PURGE_INTEL_COMPANY_ID")) != ""
}

func purgeIntelTokenDigest(token string) string {
 sum := sha256.Sum256([]byte(token))
 return hex.EncodeToString(sum[:])
}

func (s *Server) savedPurgeIntelPairing() (purgeIntelPairing, error) {
 var pair purgeIntelPairing
 if s.m == nil || s.m.pg == nil { return pair, nil }
 raw, _, err := s.m.pg.GetSetting(purgeIntelPairingKey)
 if err != nil || raw == "" { return pair, err }
 err = json.Unmarshal([]byte(raw), &pair)
 return pair, err
}

func (s *Server) authorizePurgeIntelPairing(provided string) (int64, bool) {
 if purgeIntelEnvOverride() {
  id, ok := purgeIntelBridgeCompany()
  return id, ok && validPurgeIntelBridgeCredential(provided)
 }
 pair, err := s.savedPurgeIntelPairing()
 if err != nil || pair.CompanyID < 1 || len(pair.Digest) != 64 || len(provided) < 32 { return 0, false }
 for _, key := range []string{"ARTEX_PRIVACY_BRIDGE_TOKEN", "ARTEX_PRIVACY_OVERVIEW_TOKEN", "ARTEX_SETUP_TOKEN"} {
  if provided == strings.TrimSpace(os.Getenv(key)) { return 0, false }
 }
 return pair.CompanyID, subtle.ConstantTimeCompare([]byte(pair.Digest), []byte(purgeIntelTokenDigest(provided))) == 1
}

// Admin-only management routes use the existing JWT, never the bridge token.
func (s *Server) managePurgeIntelPairing(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Cache-Control", "no-store")
 if purgeIntelEnvOverride() {
  if r.Method != http.MethodGet { writeErr(w, 409, "pairing is managed by host environment"); return }
  id, ok := purgeIntelBridgeCompany()
  configured := ok && validPurgeIntelBridgeCredential(strings.TrimSpace(os.Getenv("ARTEX_PURGE_INTEL_TOKEN")))
  writeJSON(w, 200, map[string]any{"source":"environment", "configured":configured, "company_id":id})
  return
 }
 pg := s.pg(w)
 if pg == nil { return }
 if r.Method == http.MethodGet {
  pair, err := s.savedPurgeIntelPairing()
  if err != nil { writeErr(w, 503, "pairing configuration unavailable"); return }
  writeJSON(w, 200, map[string]any{"source":"saved", "configured":pair.CompanyID > 0 && len(pair.Digest)==64, "company_id":pair.CompanyID})
  return
 }
 if r.Method == http.MethodDelete {
  if pg.SetSetting(purgeIntelPairingKey, "") != nil { writeErr(w, 503, "cannot revoke pairing"); return }
  writeJSON(w, 200, map[string]any{"source":"saved", "configured":false, "company_id":0})
  return
 }
 var req struct { CompanyID int64 `json:"company_id"` }
 if err := decode(r, &req); err != nil || req.CompanyID < 1 { writeErr(w, 400, "invalid company id"); return }
 r.SetPathValue("id", strconv.FormatInt(req.CompanyID, 10))
 id, _, ok := s.purgeIntelCompany(w, r)
 if !ok { return }
 var secret [32]byte
 if _, err := rand.Read(secret[:]); err != nil { writeErr(w, 500, "cannot generate pairing token"); return }
 token := hex.EncodeToString(secret[:])
 pair := purgeIntelPairing{CompanyID:id, Digest:purgeIntelTokenDigest(token)}
 raw, err := json.Marshal(pair)
 if err != nil || pg.SetSetting(purgeIntelPairingKey, string(raw)) != nil { writeErr(w, 503, "cannot save pairing"); return }
 // The random token is returned only once, on explicit creation/rotation.
 writeJSON(w, 200, map[string]any{"source":"saved", "configured":true, "company_id":id, "token":token})
}

func (s *Server) importPurgeIntelBridge(w http.ResponseWriter,r *http.Request) {
 id,ok:=s.authorizePurgeIntelPairing(r.Header.Get("X-Artex-Purge-Intel-Token"))
 if !ok {
  writeErr(w,http.StatusUnauthorized,"Unauthorized");return
 }
 if r.Header.Get("X-Artex-Purge-Intel-Company")!=strconv.FormatInt(id,10) {
  writeErr(w,http.StatusConflict,"receiving company does not match configured company");return
 }
 // The receiving company is configured by the administrator, never by the sender.
 r.SetPathValue("id",strconv.FormatInt(id,10))
 s.importPurgeIntel(w,r)
}

// Read-only pair check; no snapshot creation and no company scope contents returned.
func (s *Server) purgeIntelBridgeStatus(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Cache-Control", "no-store")
 id, ok := s.authorizePurgeIntelPairing(r.Header.Get("X-Artex-Purge-Intel-Token"))
 if !ok {
  writeErr(w, http.StatusUnauthorized, "Unauthorized")
  return
 }
 if r.Header.Get("X-Artex-Purge-Intel-Company") != strconv.FormatInt(id, 10) {
  writeErr(w, http.StatusConflict, "receiving company does not match configured company")
  return
 }
 r.SetPathValue("id", strconv.FormatInt(id, 10))
 _, scope, ok := s.purgeIntelCompany(w, r)
 if !ok { return }
 writeJSON(w, http.StatusOK, map[string]any{
  "service": "artex", "schema_version": 1, "company_id": id, "scope_count": len(scope), "paired": true,
 })
}

// Company-bound, read-only overview. Task metadata is deliberately limited
// to identifiers/status; no cross-company targets, prompts or evidence bodies.
func (s *Server) purgeIntelBridgeOverview(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Cache-Control", "no-store")
 id, ok := s.authorizePurgeIntelPairing(r.Header.Get("X-Artex-Purge-Intel-Token"))
 if !ok { writeErr(w, 401, "Unauthorized"); return }
 if r.Header.Get("X-Artex-Purge-Intel-Company") != strconv.FormatInt(id, 10) {
  writeErr(w, 409, "receiving company does not match configured company"); return
 }
 r.SetPathValue("id", strconv.FormatInt(id, 10))
 _, _, ok = s.purgeIntelCompany(w, r)
 if !ok { return }
 pg := s.m.pg
 type taskItem struct {
  ID int64 `json:"id"`
  Status string `json:"status"`
  Paused bool `json:"paused"`
  Queued bool `json:"queued"`
 }
 type snapshotItem struct {
  ID int64 `json:"id"`
  ReceivedAt string `json:"received_at"`
 }
 tasks := []taskItem{}
 snapshots := []snapshotItem{}
 var taskCount, assetCount, findingCount, snapshotCount int64
 const taskWhere = `t.deleted_at IS NULL AND EXISTS (SELECT 1 FROM task_scope ts WHERE ts.task_id=t.id AND ts.kind='company' AND ts.company_id=$1)`
 if err := pg.QueryRowContext(r.Context(), `SELECT count(*) FROM tasks t WHERE `+taskWhere,id).Scan(&taskCount); err != nil { writeErr(w,503,"task overview unavailable"); return }
 rows, err := pg.QueryContext(r.Context(), `SELECT t.id,t.status,t.paused,t.queued FROM tasks t WHERE `+taskWhere+` ORDER BY t.id DESC LIMIT 50`,id)
 if err != nil { writeErr(w,503,"task overview unavailable"); return }
 for rows.Next() { var item taskItem; if rows.Scan(&item.ID,&item.Status,&item.Paused,&item.Queued)!=nil { rows.Close(); writeErr(w,503,"task overview unavailable"); return }; tasks=append(tasks,item) }
 err=rows.Err(); rows.Close()
 if err!=nil { writeErr(w,503,"task overview unavailable"); return }
 if err=pg.QueryRowContext(r.Context(), `SELECT count(*) FROM assets WHERE company_id=$1`,id).Scan(&assetCount); err!=nil { writeErr(w,503,"asset overview unavailable"); return }
 // Existing company asset attribution, not all findings of a multi-company task.
 if err=pg.QueryRowContext(r.Context(), `SELECT count(*) FROM findings f WHERE EXISTS (SELECT 1 FROM assets a WHERE a.company_id=$1 AND f.asset_ids @> jsonb_build_array(a.id))`,id).Scan(&findingCount); err!=nil { writeErr(w,503,"finding overview unavailable"); return }
 if err=pg.QueryRowContext(r.Context(), `SELECT count(*) FROM purge_intel_snapshots WHERE company_id=$1`,id).Scan(&snapshotCount); err!=nil { writeErr(w,503,"snapshot overview unavailable"); return }
 rows,err=pg.QueryContext(r.Context(),`SELECT id,received_at::text FROM purge_intel_snapshots WHERE company_id=$1 ORDER BY received_at DESC,id DESC LIMIT 20`,id)
 if err!=nil { writeErr(w,503,"snapshot overview unavailable"); return }
 defer rows.Close()
 for rows.Next() { var item snapshotItem; if rows.Scan(&item.ID,&item.ReceivedAt)!=nil { writeErr(w,503,"snapshot overview unavailable"); return }; snapshots=append(snapshots,item) }
 if rows.Err()!=nil { writeErr(w,503,"snapshot overview unavailable"); return }
 writeJSON(w,200,map[string]any{"service":"artex","schema_version":1,"company_id":id,
  "task_count":taskCount,"asset_count":assetCount,"linked_finding_count":findingCount,"snapshot_count":snapshotCount,
  "tasks":tasks,"snapshots":snapshots})
}


// Paired company is resolved from the credential, never a query parameter.
func (s *Server) purgeIntelReadCompany(w http.ResponseWriter, r *http.Request) (int64, bool) {
 w.Header().Set("Cache-Control", "no-store")
 id, ok := s.authorizePurgeIntelPairing(r.Header.Get("X-Artex-Purge-Intel-Token"))
 if !ok { writeErr(w, 401, "Unauthorized"); return 0,false }
 if r.Header.Get("X-Artex-Purge-Intel-Company") != strconv.FormatInt(id,10) {
  writeErr(w,409,"receiving company does not match configured company"); return 0,false
 }
 r.SetPathValue("id",strconv.FormatInt(id,10))
 return id,true
}

func (s *Server) purgeIntelBridgeAssets(w http.ResponseWriter, r *http.Request) {
 id,ok:=s.purgeIntelReadCompany(w,r); if !ok{return}
 if _,_,ok=s.purgeIntelCompany(w,r);!ok{return}
 var before int64
 if raw:=r.URL.Query().Get("before");raw!="" {
  n,err:=strconv.ParseInt(raw,10,64);if err!=nil||n<=0{writeErr(w,400,"invalid asset cursor");return};before=n
 }
 type item struct {
  ID int64 `json:"id"`
  Type string `json:"type"`
  Domain string `json:"domain"`
  IP string `json:"ip"`
  Port int `json:"port"`
  AppName string `json:"app_name"`
  ServiceType string `json:"service_type"`
  LastSeen string `json:"last_seen"`
 }
 items:=[]item{}
 // Do not expose stored auth/params/extra or URLs containing credentials.
 rows,err:=s.m.pg.QueryContext(r.Context(),`SELECT id,type,COALESCE(domain,''),COALESCE(ip,''),COALESCE(port,0),COALESCE(app_name,''),COALESCE(service_type,''),last_seen::text
 FROM assets WHERE company_id=$1 AND ($2::bigint=0 OR id<$2) ORDER BY id DESC LIMIT 51`,id,before)
 if err!=nil{writeErr(w,503,"asset list unavailable");return};defer rows.Close()
 for rows.Next(){var a item;if rows.Scan(&a.ID,&a.Type,&a.Domain,&a.IP,&a.Port,&a.AppName,&a.ServiceType,&a.LastSeen)!=nil{writeErr(w,503,"asset list unavailable");return};items=append(items,a)}
 if rows.Err()!=nil{writeErr(w,503,"asset list unavailable");return}
 var next int64
 if len(items)>50{items=items[:50];next=items[49].ID}
 writeJSON(w,200,map[string]any{"service":"artex","schema_version":1,"company_id":id,"items":items,"next_before":next})
}

func (s *Server) purgeIntelBridgeSnapshot(w http.ResponseWriter, r *http.Request) {
 id,ok:=s.purgeIntelReadCompany(w,r);if !ok{return}
 sid,err:=strconv.ParseInt(r.URL.Query().Get("snapshot"),10,64)
 if err!=nil||sid<=0{writeErr(w,400,"invalid snapshot id");return}
 r.SetPathValue("snapshot",strconv.FormatInt(sid,10))
 // Reuse existing snapshot lookup and current company scope validation.
 b,ok:=s.loadPurgeIntelSnapshot(w,r);if !ok{return}
 a:=analyzePurgeIntel(b)
 writeJSON(w,200,map[string]any{"service":"artex","schema_version":1,"company_id":id,"snapshot_id":sid,"analysis":a})
}
