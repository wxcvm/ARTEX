package server

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCompareTrafficAPI(t *testing.T) {
	s,request:=ingestServer(t)
	if w:=request("GET","/api/traffic/compare?baseline=x&candidate=x",""); w.Code!=400 { t.Fatalf("same ID: %d %s",w.Code,w.Body) }
	if w:=request("GET","/api/traffic/compare?baseline=x&candidate=y",""); w.Code!=422 { t.Fatalf("missing record: %d %s",w.Code,w.Body) }
	r:=httptest.NewRequest("GET","/api/traffic/compare?baseline=x&candidate=y",strings.NewReader(""))
	w:=httptest.NewRecorder()
	s.Handler().ServeHTTP(w,r)
	if w.Code!=401 { t.Fatalf("unauthenticated: %d",w.Code) }
}
