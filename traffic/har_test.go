package traffic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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


func TestHARPageAttributionRoundTrip(t *testing.T) {
	tr,_ := openTraffic(t)
	// Only the page ordinal and association state are persisted, never the
	// original page identifier, title or unrelated page fields.
	var doc struct { Log struct { Entries []json.RawMessage `json:"entries"` } `json:"log"` }
	if err := json.Unmarshal([]byte(harFixture), &doc); err != nil { t.Fatal(err) }
	var first map[string]any
	if err := json.Unmarshal(doc.Log.Entries[0], &first); err != nil { t.Fatal(err) }
	refs := []any{"PAGE_PRIVATE","DUP_PRIVATE","not-listed",nil,42,"PAGE_PRIVATE"}
	entries := make([]any,0,len(refs))
	for _, ref := range refs {
		e := make(map[string]any,len(first)+1)
		for k,v := range first { e[k] = v }
		e["pageref"] = ref
		entries = append(entries,e)
	}
	raw,err := json.Marshal(map[string]any{"log":map[string]any{
		"pages":[]any{map[string]any{"id":"PAGE_PRIVATE","title":"SECRET_TITLE"},map[string]any{"id":"DUP_PRIVATE"},map[string]any{"id":"DUP_PRIVATE"}},
		"entries":entries,
	}})
	if err != nil { t.Fatal(err) }
	result,err := tr.ImportHAR(context.Background(),raw)
	if err != nil || result.Imported != len(refs) { t.Fatalf("import: %+v %v",result,err) }
	rows,total,err := tr.Page(PageQuery{Host:"example.invalid",RespMin:-1,RespMax:-1},0,10)
	if err != nil || total != len(refs) || len(rows) != len(refs) { t.Fatalf("page: %v %v",rows,err) }
	want := []HARPageAttribution{{State:"unique",Ordinal:1},{State:"ambiguous"},{State:"unknown"},{State:"missing"},{State:"unknown"},{State:"unique",Ordinal:1}}
	check := func(store *Traffic) {
		t.Helper()
		for _,row := range rows {
			i := strings.LastIndex(row.ID,"-")
			var index int
			if _,err := fmt.Sscanf(row.ID[i+1:],"%d",&index); err != nil { t.Fatal(err) }
			err := store.ReadEvidence(context.Background(),[]string{row.ID},func(e EvidenceExchange) error {
				if e.HARPage == nil || *e.HARPage != want[index] { t.Fatalf("page %d: %+v",index,e.HARPage) }
				return nil
			})
			if err != nil { t.Fatal(err) }
		}
	}
	check(tr)
	if retry,err := tr.ImportHAR(context.Background(),raw); err != nil || retry.Duplicates != len(refs) { t.Fatalf("retry: %+v %v",retry,err) }
	dir := t.TempDir()
	if n,err := tr.ExportHosts([]string{"example.invalid"},dir); err != nil || n != int64(len(refs)) { t.Fatalf("export: %d %v",n,err) }
	manifest,err := os.ReadFile(filepath.Join(dir,"traffic.json"))
	if err != nil { t.Fatal(err) }
	for _,s := range []string{"PAGE_PRIVATE","DUP_PRIVATE","SECRET_TITLE","not-listed"} {
		if strings.Contains(string(manifest),s) { t.Fatalf("raw page metadata leaked: %s",s) }
	}
	restored,_ := openTraffic(t)
	if n,err := restored.ImportArchive(dir); err != nil || n != int64(len(refs)) { t.Fatalf("restore: %d %v",n,err) }
	check(restored)
	if n,err := restored.DeleteHost("example.invalid"); err != nil || n != int64(len(refs)) { t.Fatalf("delete: %d %v",n,err) }
	var remaining int
	if err := restored.db.QueryRow("SELECT count(*) FROM har_page_attribution").Scan(&remaining); err != nil || remaining != 0 { t.Fatalf("orphaned page records: %d %v",remaining,err) }
}
