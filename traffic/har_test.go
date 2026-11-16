package traffic

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

const harFixture = `{"log":{"entries":[
{"startedDateTime":"2026-11-16T18:01:02.123+08:00","request":{"method":"GET","url":"https://example.invalid/record","httpVersion":"HTTP/2","headers":[{"name":"Accept","value":"text/plain"},{"name":"Accept","value":"application/json"}],"bodySize":0},"response":{"status":200,"httpVersion":"HTTP/2","headers":[{"name":"Set-Cookie","value":"a=1"},{"name":"Set-Cookie","value":"b=2"}],"content":{"size":3,"text":"abc","mimeType":"text/plain"}}},
{"startedDateTime":"2026-11-16T10:01:03Z","request":{"method":"POST","url":"https://example.invalid/binary","headers":[],"bodySize":3,"postData":{"text":"abc"}},"response":{"status":201,"headers":[],"content":{"size":4,"text":"AAEC/w==","encoding":"base64","mimeType":"application/octet-stream"}}},
{"startedDateTime":"2026-11-16T10:01:04Z","request":{"method":"GET","url":"https://example.invalid/missing","headers":[],"bodySize":0},"response":{"status":200,"headers":[],"content":{"size":5}}},
{"startedDateTime":"2026-11-16T10:01:05Z","request":{"method":"GET","url":"https://example.invalid/unknown","headers":[],"bodySize":0},"response":{"status":0,"headers":[],"content":{"size":0}}},
{"startedDateTime":"2026-11-16T10:01:06Z","request":{"method":"POST","url":"https://example.invalid/form","headers":[],"bodySize":3,"postData":{"mimeType":"application/x-www-form-urlencoded"}},"response":{"status":200,"headers":[],"content":{"size":0}}},
{"startedDateTime":"2026-11-16T10:01:07","request":{"method":"GET","url":"https://example.invalid/time","headers":[],"bodySize":0},"response":{"status":200,"headers":[],"content":{"size":0}}},
null
]}}`

func TestImportHARPreservesRecordingAndDeduplicates(t *testing.T) {
	tr,_:=openTraffic(t)
	result,err:=tr.ImportHAR(context.Background(),[]byte(harFixture))
	if err!=nil || result.Entries!=7 || result.Imported!=2 || result.Skipped!=5 || result.Duplicates!=0 || result.Failed!=0 {
		t.Fatalf("import: %+v %v",result,err)
	}
	for _,reason:=range []string{"response_body_missing","unobserved_status","request_body_missing","unknown_capture_time","invalid_entry"} {
		if result.Reasons[reason]!=1 { t.Fatalf("missing reason %s: %+v",reason,result) }
	}
	rows,total,err:=tr.Page(PageQuery{Host:"example.invalid",RespMin:-1,RespMax:-1},0,10)
	if err!=nil || total!=2 || len(rows)!=2 { t.Fatalf("page: %d %+v %v",total,rows,err) }
	var ids []string
	for _,row:=range rows {
		if !strings.HasPrefix(row.ID,"har-") { t.Fatalf("missing provenance: %s",row.ID) }
		ids=append(ids,row.ID)
	}
	check:=func(store *Traffic) {
		t.Helper()
		if err:=store.ReadEvidence(context.Background(),ids,func(e EvidenceExchange) error {
			request,err:=io.ReadAll(e.Request); if err!=nil { return err }
			response,err:=io.ReadAll(e.Response); if err!=nil { return err }
			switch e.URL {
			case "https://example.invalid/record":
				want,_:=time.Parse(time.RFC3339Nano,"2026-11-16T18:01:02.123+08:00")
				if e.TS!=want.Unix() || e.Status!=200 || string(response)!="abc" ||
					!strings.Contains(e.ReqHead,"Accept: text/plain\nAccept: application/json") ||
					!strings.Contains(e.RespHead,"Set-Cookie: a=1\nSet-Cookie: b=2") { t.Fatalf("record changed: %+v",e) }
			case "https://example.invalid/binary":
				if e.Status!=201 || string(request)!="abc" || !bytes.Equal(response,[]byte{0,1,2,255}) { t.Fatalf("binary changed: %v",response) }
			default: t.Fatalf("unexpected URL: %s",e.URL)
			}
			return nil
		}); err!=nil { t.Fatal(err) }
	}
	check(tr)
	retry,err:=tr.ImportHAR(context.Background(),[]byte(harFixture))
	if err!=nil || retry.Imported!=0 || retry.Duplicates!=2 || retry.Skipped!=5 { t.Fatalf("retry: %+v %v",retry,err) }
	dir:=t.TempDir()
	if n,err:=tr.ExportHosts([]string{"example.invalid"},dir); err!=nil || n!=2 { t.Fatalf("archive export: %d %v",n,err) }
	restored,_:=openTraffic(t)
	if n,err:=restored.ImportArchive(dir); err!=nil || n!=2 { t.Fatalf("archive restore: %d %v",n,err) }
	check(restored)
}

func TestImportHARRejectsMalformedContainer(t *testing.T) {
	tr,_:=openTraffic(t)
	for _,data:=range []string{"{}","{\"log\":{\"entries\":null}}","{\"log\":{\"entries\":{}}}","not json"} {
		if _,err:=tr.ImportHAR(context.Background(),[]byte(data)); err==nil { t.Fatalf("accepted malformed HAR: %s",data) }
	}
	if n,_:=tr.Count(); n!=0 { t.Fatalf("malformed HAR wrote %d records",n) }
}
