package traffic

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestCompareRecordedStoredBytesAndJSON(t *testing.T) {
	tr,_:=openTraffic(t)
	const fixture=`{"log":{"entries":[
{"startedDateTime":"2026-11-16T10:00:00Z","request":{"method":"GET","url":"https://example.invalid/record","headers":[{"name":"Accept","value":"a"},{"name":"Accept","value":"b"}],"bodySize":0},"response":{"status":200,"headers":[{"name":"X-Mode","value":"private-one"},{"name":"X-Old","value":"a"}],"content":{"text":"{\"same\":{\"a\":1,\"b\":2},\"changed\":1,\"removed\":true}"}}},
{"startedDateTime":"2026-11-16T10:00:01Z","request":{"method":"GET","url":"https://example.invalid/record","headers":[{"name":"accept","value":"b"},{"name":"ACCEPT","value":"a"}],"bodySize":0},"response":{"status":201,"headers":[{"name":"x-mode","value":"private-two"},{"name":"X-New","value":"b"}],"content":{"text":"{\"same\":{\"b\":2,\"a\":1},\"changed\":2,\"added\":true}"}}}
]}}`
	result,err:=tr.ImportHAR(context.Background(),[]byte(fixture))
	if err!=nil || result.Imported!=2 { t.Fatalf("import: %+v %v",result,err) }
	rows,total,err:=tr.Page(PageQuery{Host:"example.invalid",RespMin:-1,RespMax:-1},0,10)
	if err!=nil || total!=2 || len(rows)!=2 { t.Fatalf("rows: %+v %v",rows,err) }
	a,b:=rows[0],rows[1]
	if a.Status!=200 { a,b=b,a }
	got,err:=tr.CompareRecorded(context.Background(),a.ID,b.ID)
	if err!=nil { t.Fatal(err) }
	if !got.SameURL || got.StatusRelation!="different" || got.RequestBody.Relation!="equal" || got.ResponseBody.Relation!="different" { t.Fatalf("comparison: %+v",got) }
	if got.RequestHeaders.ChangedCount!=1 || got.RequestHeaders.Changed[0]!="accept" ||
		got.ResponseHeaders.AddedCount!=1 || got.ResponseHeaders.RemovedCount!=1 || got.ResponseHeaders.ChangedCount!=1 {
		t.Fatalf("headers: %+v %+v",got.RequestHeaders,got.ResponseHeaders)
	}
	j:=got.ResponseBody.JSONFields
	if !j.Available || j.AddedCount!=1 || j.Added[0]!="added" || j.RemovedCount!=1 || j.Removed[0]!="removed" || j.ChangedCount!=1 || j.Changed[0]!="changed" {
		t.Fatalf("JSON: %+v",j)
	}
	raw,_:=json.Marshal(got)
	if strings.Contains(string(raw),"private-one") || strings.Contains(string(raw),"private-two") { t.Fatal("header value escaped comparison") }
	if _,err:=tr.CompareRecorded(context.Background(),a.ID,a.ID); err==nil { t.Fatal("accepted same ID") }
	if _,err:=tr.CompareRecorded(context.Background(),a.ID,"missing"); err==nil { t.Fatal("accepted missing ID") }
}

func TestRecordedBodyComparisonUsesFullStoredBytes(t *testing.T) {
	ctx:=context.Background()
	a:=strings.Repeat("a",compareJSONBudget+10)
	b:=a[:len(a)-1]+"b"
	x,err:=inspectRecordedBody(ctx,strings.NewReader(a),int64(len(a)),"")
	if err!=nil { t.Fatal(err) }
	y,err:=inspectRecordedBody(ctx,strings.NewReader(b),int64(len(b)),"")
	if err!=nil { t.Fatal(err) }
	if c:=compareRecordedBodies(x,y); c.Relation!="different" || c.JSONFields.Available { t.Fatalf("large body: %+v",c) }
	z,err:=inspectRecordedBody(ctx,strings.NewReader(a),int64(len(a)+1),"")
	if err!=nil { t.Fatal(err) }
	if c:=compareRecordedBodies(x,z); c.Relation!="unknown" { t.Fatalf("length mismatch: %+v",c) }
	bad,err:=inspectRecordedBody(ctx,strings.NewReader(a),int64(len(a)),"wrong")
	if err!=nil { t.Fatal(err) }
	if c:=compareRecordedBodies(x,bad); c.Relation!="unknown" { t.Fatalf("bad blob hash: %+v",c) }
}
