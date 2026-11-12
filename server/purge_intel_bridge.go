package server

import (
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

func (s *Server) importPurgeIntelBridge(w http.ResponseWriter,r *http.Request) {
 id,ok:=purgeIntelBridgeCompany()
 if !ok || !validPurgeIntelBridgeCredential(r.Header.Get("X-Artex-Purge-Intel-Token")) {
  writeErr(w,http.StatusUnauthorized,"Unauthorized");return
 }
 if r.Header.Get("X-Artex-Purge-Intel-Company")!=strconv.FormatInt(id,10) {
  writeErr(w,http.StatusConflict,"receiving company does not match configured company");return
 }
 // The receiving company is configured by the administrator, never by the sender.
 r.SetPathValue("id",strconv.FormatInt(id,10))
 s.importPurgeIntel(w,r)
}
