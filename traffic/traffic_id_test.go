package traffic

import (
	"testing"
)

func TestTrafficPageSearchByStableID(t *testing.T) {
	tr, err := Open(t.TempDir(), "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer tr.Close()
	const wanted = "1700000000-0001"
	for _, tc := range []struct{ id, host string }{
		{wanted, "id.example"},
		{"1700000000-0002", "other.example"},
	} {
		if _, err := tr.DB().Exec(`INSERT INTO exchanges(id,ts,host,method,url_template,url,status,content_type,req_len,resp_len,path)
VALUES(?,?,?,?,?,?,?,?,?,?,?)`, tc.id, 1700000000, tc.host, "GET", "/proof", "http://"+tc.host+"/proof", 200, "text/plain", 0, 1, ""); err != nil {
			t.Fatal(err)
		}
	}
	for _, q := range []string{wanted, "0001"} {
		rows, total, err := tr.Page(PageQuery{Query: q, RespMin: -1, RespMax: -1}, 0, 20)
		if err != nil || total != 1 || len(rows) != 1 || rows[0].ID != wanted {
			t.Fatalf("query %q: rows=%v total=%d err=%v", q, rows, total, err)
		}
	}
	rows, total, err := tr.Page(PageQuery{Host: "other.example", Query: wanted, RespMin: -1, RespMax: -1}, 0, 20)
	if err != nil || total != 0 || len(rows) != 0 {
		t.Fatalf("host filter must still apply: rows=%v total=%d err=%v", rows, total, err)
	}
}
