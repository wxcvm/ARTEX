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

func TestCompareRecordedFormAndXML(t *testing.T) {
	ctx := context.Background()
	read := func(raw, ct string) recordedBody {
		t.Helper()
		b, err := inspectRecordedBody(ctx, strings.NewReader(raw), int64(len(raw)), "", ct)
		if err != nil { t.Fatal(err) }
		return b
	}
	a := read("tag=a&tag=b&old=1&same=x+y", "application/x-www-form-urlencoded")
	b := read("same=x%20y&tag=b&tag=a&new=2", "application/x-www-form-urlencoded; charset=UTF-8")
	f := compareRecordedBodies(a,b).FormFields
	if !f.Available || f.AddedCount != 1 || f.Added[0] != "new" || f.RemovedCount != 1 || f.Removed[0] != "old" || f.ChangedCount != 1 || f.Changed[0] != "tag" {
		t.Fatalf("form changes: %+v", f)
	}
	if compareRecordedBodies(a,read("tag=%ZZ","application/x-www-form-urlencoded")).FormFields.Available {
		t.Fatal("invalid form treated as parsed")
	}
	x := read(`<r xmlns="urn:x"><a>one</a><a>two</a><old/></r>`,"application/xml")
	y := read(`<r xmlns="urn:x"><a>two</a><a>one</a><new/></r>`,"application/example+xml")
	d := compareRecordedBodies(x,y).XMLChildren
	if !d.Available || d.AddedCount != 1 || d.Added[0] != "{urn:x}new" || d.RemovedCount != 1 || d.ChangedCount != 1 || d.Changed[0] != "{urn:x}a" {
		t.Fatalf("XML changes: %+v", d)
	}
	for _, raw := range []string{`<r><a></r>`, `<!DOCTYPE r SYSTEM "https://example.invalid/external"><r/>`, `<r/><r/>`} {
		if compareRecordedBodies(x,read(raw,"application/xml")).XMLChildren.Available {
			t.Fatalf("invalid XML compared: %s", raw)
		}
	}
	if compareRecordedBodies(x,read("<other><a/></other>","application/xml")).XMLChildren.Available {
		t.Fatal("different roots treated as comparable")
	}
}
