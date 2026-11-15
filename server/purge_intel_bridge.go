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
