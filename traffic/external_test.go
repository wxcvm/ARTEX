package traffic

import (
	"strings"
	"testing"
)

func TestRecordExternalStoresLikeProxyTraffic(t *testing.T) {
	tr, _ := openTraffic(t)

	big := strings.Repeat("A", 2<<20) // spills to the blob store
	err := tr.RecordExternal(ExternalExchange{
		Method:          "POST",
		URL:             "https://api.example.com/v1/login",
		Status:          200,
		RequestHeaders:  map[string]string{"content-type": "application/json", "authorization": "Bearer xyz"},
		RequestBody:     `{"user":"u1","password":"p1"}`,
		ResponseHeaders: map[string]string{"content-type": "application/json"},
		ResponseBody:    `{"token":"tok"}`,
	})
	if err != nil {
		t.Fatalf("RecordExternal small: %v", err)
	}
	if err := tr.RecordExternal(ExternalExchange{
		Method:          "GET",
		URL:             "https://api.example.com/v1/items?page=1",
		Status:          404,
		ResponseHeaders: map[string]string{"content-type": "text/html"},
		ResponseBody:    big,
	}); err != nil {
		t.Fatalf("RecordExternal spilled: %v", err)
	}

	if n, _ := tr.Count(); n != 2 {
		t.Fatalf("count = %d, want 2", n)
	}
	rows, total, err := tr.Page(PageQuery{Query: "login", RespMin: -1, RespMax: -1}, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(rows) != 1 {
		t.Fatalf("filtered page = %d rows (total %d), want 1", len(rows), total)
	}
	req, resp, err := tr.Get(rows[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(req, "POST /v1/login") || !strings.Contains(req, "Bearer xyz") || !strings.Contains(req, `"user":"u1"`) {
		t.Fatalf("request side incomplete: %q", req)
	}
	if !strings.Contains(resp, `{"token":"tok"}`) {
		t.Fatalf("response side incomplete: %q", resp)
	}

	// The spilled entry must be FTS-searchable by its body content too.
	rows, total, err = tr.Page(PageQuery{Body: strings.Repeat("A", 64), RespMin: -1, RespMax: -1}, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(rows) != 1 {
		t.Fatalf("spilled-body FTS search = %d rows (total %d), want 1", len(rows), total)
	}
}

func TestRecordExternalRejectsBadInput(t *testing.T) {
	tr, _ := openTraffic(t)
	cases := []ExternalExchange{
		{Method: "GET", URL: "not a url"},
		{Method: "GET", URL: ""},
		{Method: "", URL: "https://api.example.com/"},
	}
	for i, c := range cases {
		if err := tr.RecordExternal(c); err == nil {
			t.Fatalf("case %d: expected an error for %+v", i, c)
		}
	}
}

func TestRecordExternalPreservesUnobservedStatus(t *testing.T) {
	tr, _ := openTraffic(t)
	for _, e := range []ExternalExchange{
		{Method: "GET", URL: "https://example.invalid/unobserved", ResponseBody: "partial capture"},
		{Method: "GET", URL: "https://example.invalid/observed", Status: 200, ResponseBody: "ok"},
	} {
		if err := tr.RecordExternal(e); err != nil {
			t.Fatal(err)
		}
	}
	rows, total, err := tr.Page(PageQuery{Status: "0", RespMin: -1, RespMax: -1}, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(rows) != 1 || rows[0].Status != 0 {
		t.Fatalf("unknown-status filter: rows=%+v total=%d", rows, total)
	}
	_, resp, err := tr.Get(rows[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp, "HTTP status unavailable") || !strings.Contains(resp, "partial capture") || strings.Contains(resp, "HTTP 200") {
		t.Fatalf("unobserved status was misrepresented: %q", resp)
	}
	rows, total, err = tr.Page(PageQuery{Status: "2xx", RespMin: -1, RespMax: -1}, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(rows) != 1 || rows[0].Status != 200 {
		t.Fatalf("2xx must exclude unknown status: rows=%+v total=%d", rows, total)
	}
}
