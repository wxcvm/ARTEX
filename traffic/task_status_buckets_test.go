package traffic

import (
	"reflect"
	"testing"
)

func TestStatusBucketsMatchFullFilteredSet(t *testing.T) {
	tr, _ := openTraffic(t)
	for _, row := range []struct {
		id, host, capture string
		ts int64
		status int
	}{
		{"a-unobserved", "a.example", "", 100, 0},
		{"a-ok", "a.example", "group-one", 150, 200},
		{"a-client-error", "a.example", "group-one", 170, 404},
		{"a-server-error", "a.example", "", 210, 503},
		{"b-ok", "b.example", "group-one", 165, 201},
	} {
		_, err := tr.DB().Exec(`INSERT INTO exchanges(id,ts,host,method,url_template,url,status,content_type,req_len,resp_len,path)
VALUES(?,?,?,?,?,?,?,?,?,?,?)`, row.id, row.ts, row.host, "GET", "/", "http://"+row.host+"/", row.status, "text/plain", 0, 0, "")
		if err != nil { t.Fatal(err) }
		if row.capture != "" {
			if _, err := tr.DB().Exec(`INSERT INTO capture_attribution(exchange_id,capture_id) VALUES (?,?)`, row.id, row.capture); err != nil {
				t.Fatal(err)
			}
		}
	}
	tests := []struct {
		name string
		filter PageQuery
		want map[int]int
	}{
		{"global includes unobserved", PageQuery{RespMin:-1,RespMax:-1}, map[int]int{0:1,2:2,4:1,5:1}},
		{"task and time", PageQuery{TaskHosts:[]string{"a.example"},TsMin:120,TsMax:200,RespMin:-1,RespMax:-1}, map[int]int{2:1,4:1}},
		{"capture group intersects task scope", PageQuery{TaskHosts:[]string{"a.example"},CaptureID:"group-one",TsMin:120,TsMax:200,RespMin:-1,RespMax:-1}, map[int]int{2:1,4:1}},
		{"status class", PageQuery{Status:"4xx",RespMin:-1,RespMax:-1}, map[int]int{4:1}},
		{"empty task never becomes global", PageQuery{TaskHosts:[]string{},RespMin:-1,RespMax:-1}, map[int]int{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buckets, err := tr.StatusBuckets(tc.filter)
			if err != nil { t.Fatal(err) }
			if !reflect.DeepEqual(buckets, tc.want) {
				t.Fatalf("StatusBuckets = %v, want %v", buckets, tc.want)
			}
		})
	}
	// Query one visible record but still aggregate *all* matching rows.
	filter := PageQuery{TaskHosts: []string{"a.example"}, TsMin:120, TsMax:200, RespMin:-1, RespMax:-1}
	page, total, err := tr.Page(filter, 0, 1)
	if err != nil || len(page) != 1 || total != 2 {
		t.Fatalf("Page: len=%d total=%d err=%v", len(page), total, err)
	}
	buckets, err := tr.StatusBuckets(filter)
	if err != nil || buckets[2] != 1 || buckets[4] != 1 {
		t.Fatalf("full filtered distribution = %v, err=%v", buckets, err)
	}
}
