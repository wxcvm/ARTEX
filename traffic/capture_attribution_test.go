package traffic

import (
    "context"
    "strings"
    "testing"
)

func TestCaptureAttributionSurvivesArchiveAndDelete(t *testing.T) {
    tr,_:=openTraffic(t)
    captureID:="00000000-0000-4000-8000-000000000001"
    if err:=tr.RecordExternal(ExternalExchange{
        CaptureID:captureID,Method:"GET",URL:"https://example.invalid/capture",Status:200,ResponseBody:"ok",
    });err!=nil{t.Fatal(err)}
    rows,total,err:=tr.Page(PageQuery{Host:"example.invalid",RespMin:-1,RespMax:-1},0,10)
    if err!=nil||total!=1||len(rows)!=1{t.Fatalf("page: %+v %d %v",rows,total,err)}
    if rows[0].CaptureID!=captureID{t.Fatalf("page lost capture group: %q",rows[0].CaptureID)}
    id:=rows[0].ID
    read:=func(store *Traffic){
        t.Helper()
        if err:=store.ReadEvidence(context.Background(),[]string{id},func(e EvidenceExchange)error{
            if e.CaptureID!=captureID{t.Fatalf("lost capture group: %q",e.CaptureID)}
            if e.HARPage!=nil{t.Fatal("capture group misclassified as HAR page")}
            return nil
        });err!=nil{t.Fatal(err)}
    }
    read(tr)
    dir:=t.TempDir()
    if n,err:=tr.ExportHosts([]string{"example.invalid"},dir);err!=nil||n!=1{t.Fatalf("archive export: %d %v",n,err)}
    restored,_:=openTraffic(t)
    if n,err:=restored.ImportArchive(dir);err!=nil||n!=1{t.Fatalf("archive import: %d %v",n,err)}
    read(restored)
    restoredRows,restoredTotal,restoredErr:=restored.Page(PageQuery{Host:"example.invalid",RespMin:-1,RespMax:-1},0,10)
    if restoredErr!=nil||restoredTotal!=1||len(restoredRows)!=1||restoredRows[0].CaptureID!=captureID{t.Fatalf("restored listing: %+v %d %v",restoredRows,restoredTotal,restoredErr)}
    if n,err:=restored.ImportArchive(dir);err!=nil||n!=0{t.Fatalf("archive duplicate: %d %v",n,err)}
    if n,err:=tr.DeleteHost("example.invalid");err!=nil||n!=1{t.Fatalf("delete: %d %v",n,err)}
    var count int
    if err:=tr.db.QueryRow("SELECT COUNT(*) FROM capture_attribution").Scan(&count);err!=nil||count!=0{
        t.Fatalf("orphan capture rows: %d %v",count,err)
    }
}

func TestCaptureAttributionOptionalAndLength(t *testing.T) {
    tr,_:=openTraffic(t)
    if err:=tr.RecordExternal(ExternalExchange{Method:"GET",URL:"https://example.invalid/legacy",Status:200});err!=nil{t.Fatal(err)}
    rows,total,err:=tr.Page(PageQuery{Host:"example.invalid",RespMin:-1,RespMax:-1},0,10)
    if err!=nil||total!=1{t.Fatalf("page: %d %v",total,err)}
    if rows[0].CaptureID!=""{t.Fatalf("legacy row invented group: %q",rows[0].CaptureID)}
    if err:=tr.ReadEvidence(context.Background(),[]string{rows[0].ID},func(e EvidenceExchange)error{
        if e.CaptureID!=""{t.Fatalf("invented attribution: %q",e.CaptureID)}
        return nil
    });err!=nil{t.Fatal(err)}
    if err:=tr.RecordExternal(ExternalExchange{Method:"GET",URL:"https://example.invalid/invalid",CaptureID:strings.Repeat("x",129)});err==nil{t.Fatal("accepted unbounded capture ID")}
}


func TestCaptureGroupPageFilter(t *testing.T) {
    tr,_:=openTraffic(t)
    for _,tc:=range []struct{id,host string}{
        {"batch-a","one.example.invalid"},{"batch-a","two.example.invalid"},
        {"batch-b","one.example.invalid"},{"","one.example.invalid"},
    }{
        if err:=tr.RecordExternal(ExternalExchange{
            CaptureID:tc.id,Method:"GET",URL:"https://"+tc.host+"/",Status:200,
        });err!=nil{t.Fatal(err)}
    }
    q:=PageQuery{CaptureID:"batch-a",RespMin:-1,RespMax:-1}
    rows,total,err:=tr.Page(q,0,20)
    if err!=nil||total!=2||len(rows)!=2{t.Fatalf("batch filter: %+v %d %v",rows,total,err)}
    for _,row:=range rows{if row.CaptureID!="batch-a"{t.Fatalf("wrong filtered group: %+v",row)}}
    q.Host="one.example.invalid"
    rows,total,err=tr.Page(q,0,20)
    if err!=nil||total!=1||len(rows)!=1{t.Fatalf("combined filters: %+v %d %v",rows,total,err)}
    q.Host=""
    q.CaptureID="does-not-exist"
    rows,total,err=tr.Page(q,0,20)
    if err!=nil||total!=0||len(rows)!=0{t.Fatalf("unknown batch: %+v %d %v",rows,total,err)}
    q.CaptureID=""
    rows,total,err=tr.Page(q,0,20)
    if err!=nil||total!=4||len(rows)!=4{t.Fatalf("unfiltered legacy rows: %+v %d %v",rows,total,err)}
    groupCounts:=map[string]int{}
    for _,row:=range rows{groupCounts[row.CaptureID]++}
    if groupCounts["batch-a"]!=2||groupCounts["batch-b"]!=1||groupCounts[""]!=1{t.Fatalf("unfiltered attribution: %+v",groupCounts)}
    searched,err:=tr.Search("one.example.invalid",0,20)
    if err!=nil||len(searched)!=3{t.Fatalf("search attribution: %+v %v",searched,err)}
    seen:=map[string]int{}
    for _,row:=range searched{seen[row.CaptureID]++}
    if seen["batch-a"]!=1||seen["batch-b"]!=1||seen[""]!=1{t.Fatalf("search attribution lost: %+v",seen)}
}
