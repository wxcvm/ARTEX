package server

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/traffic"
)

func TestHARImportAPI(t *testing.T) {
	s,request:=ingestServer(t)
	const body=`{"log":{"entries":[{"startedDateTime":"2026-11-16T10:00:00Z","request":{"method":"GET","url":"https://example.invalid/har","headers":[],"bodySize":0},"response":{"status":204,"headers":[],"content":{"size":0}}}]}}`
	for i:=0;i<2;i++ {
		w:=request("POST","/api/traffic/har/import",body)
		var result traffic.HARImportResult
		if w.Code!=200 || json.Unmarshal(w.Body.Bytes(),&result)!=nil || result.Imported!=1-i || result.Duplicates!=i { t.Fatalf("import %d: %d %s",i,w.Code,w.Body) }
	}
	if w:=request("POST","/api/traffic/har/import","{}");w.Code!=400 { t.Fatalf("invalid HAR: %d %s",w.Code,w.Body) }
	r:=httptest.NewRequest("POST","/api/traffic/har/import",strings.NewReader(body))
	w:=httptest.NewRecorder()
	s.Handler().ServeHTTP(w,r)
	if w.Code!=401 { t.Fatalf("unauthenticated import: %d",w.Code) }
	if n,err:=s.m.Traffic().Count();err!=nil || n!=1 { t.Fatalf("traffic count: %d %v",n,err) }
}
