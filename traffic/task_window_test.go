package traffic

import "testing"

func TestPageTaskHostsAndActivityWindow(t *testing.T) {
	tr, _ := openTraffic(t)
	for _, row := range []struct{
		id string
		ts int64
		host string
	}{
		{"a-early", 100, "a.example"},
		{"a-late", 200, "a.example"},
		{"b-middle", 150, "b.example"},
		{"c-middle", 150, "c.example"},
	} {
		if _, err := tr.DB().Exec(`INSERT INTO exchanges(id,ts,host,method,url_template,url,status,content_type,req_len,resp_len,path)
VALUES(?,?,?,?,?,?,?,?,?,?,?)`, row.id,row.ts,row.host,"GET","/","http://"+row.host+"/",200,"text/plain",0,0,""); err!=nil {
			t.Fatal(err)
		}
	}
	if _,err:=tr.DB().Exec(`INSERT INTO capture_attribution(exchange_id,capture_id) VALUES('b-middle','group-one')`);err!=nil { t.Fatal(err) }
	tests := []struct{
		name string
		q PageQuery
		want string
		count int
	}{
		{"combined hosts and window", PageQuery{TaskHosts: []string{"a.example","b.example"},TsMin:110,TsMax:180,RespMin:-1,RespMax:-1},"b-middle",1},
		{"same host later task", PageQuery{TaskHosts: []string{"a.example"},TsMin:190,TsMax:210,RespMin:-1,RespMax:-1},"a-late",1},
		{"no assets means no results", PageQuery{TaskHosts: []string{},RespMin:-1,RespMax:-1},"",0},
		{"global traffic unchanged", PageQuery{RespMin:-1,RespMax:-1},"",4},
		{"existing host filter intersects", PageQuery{Host:"a.",TaskHosts:[]string{"b.example"},RespMin:-1,RespMax:-1},"",0},
		{"capture group intersection", PageQuery{TaskHosts:[]string{"a.example","b.example"},TsMin:110,TsMax:180,CaptureID:"group-one",RespMin:-1,RespMax:-1},"b-middle",1},
	}
	for _,tc:=range tests{
		t.Run(tc.name,func(t *testing.T){
			rows,total,err:=tr.Page(tc.q,0,20)
			if err!=nil||total!=tc.count||len(rows)!=tc.count {
				t.Fatalf("rows=%v total=%d err=%v",rows,total,err)
			}
			if tc.want!="" && rows[0].ID!=tc.want { t.Fatalf("id=%s, want %s",rows[0].ID,tc.want) }
		})
	}
}
