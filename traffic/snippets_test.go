package traffic

import (
	"strings"
	"testing"
)

func TestBodyMatchSnippetsInlineAndUnicode(t *testing.T) {
	tr, _ := openTraffic(t)
	if !tr.fts { t.Skip("requires FTS5") }
	tr.record(newFlow("api.example.com", "POST", "/login", nil, []byte("The login failed: Invalid Credentials near 你好世界!")))
	page, count, err := tr.Page(PageQuery{Body:"invalid credentials",RespMin:-1,RespMax:-1},0,10)
	if err!=nil || count!=1 || len(page)!=1 { t.Fatalf("page=%v count=%d err=%v",page,count,err) }
	snips:=tr.BodyMatchSnippets(page, "invalid credentials")
	sn, ok := snips[page[0].ID]
	if !ok || sn.Match!="Invalid Credentials" || !strings.Contains(sn.Before,"login failed") {
		t.Fatalf("bad snippet: %+v, present=%v",sn,ok)
	}
	zh:=tr.BodyMatchSnippets(page,"你好世界")
	if zh[page[0].ID].Match!="你好世界" { t.Fatalf("unicode: %+v",zh) }
	if got:=tr.BodyMatchSnippets(page,"ab");len(got)!=0 { t.Fatalf("short term unexpectedly matched: %v",got) }
	if got:=tr.BodyMatchSnippets(page,"absent");len(got)!=0 { t.Fatalf("missing term: %v",got) }
}

func TestBodyMatchSnippetsSpilledBodyBeyondPreview(t *testing.T) {
	tr, _ := openTraffic(t)
	if !tr.fts { t.Skip("requires FTS5") }
	// Keyword is past the 8 KiB SQLite preview in the blob, yet within the
	// indexed range. This must not display a fake "no context" result.
	payload:=[]byte(strings.Repeat("x",maxInlineBody+16)+"UniqueSearchNeedle"+strings.Repeat("y",80))
	tr.record(newFlow("large.example.com", "GET", "/large", nil,payload))
	page,count,err:=tr.Page(PageQuery{Body:"UniqueSearchNeedle",RespMin:-1,RespMax:-1},0,5)
	if err!=nil||count!=1||len(page)!=1 { t.Fatalf("page count=%d len=%d err=%v",count,len(page),err) }
	got:=tr.BodyMatchSnippets(page,"uniquesearchneedle")
	if got[page[0].ID].Match!="UniqueSearchNeedle" { t.Fatalf("spilled body snippet: %+v",got) }
}

func TestMatchTextSnippetPlainTextBounds(t *testing.T) {
	s,ok:=matchTextSnippet(strings.Repeat("A",90)+"<script>alert(1)</script>"+strings.Repeat("B",90),"<script>")
	if !ok || s.Match!="<script>" || !strings.HasPrefix(s.Before,"…") || !strings.HasSuffix(s.After,"…") {
		t.Fatalf("snippet: %+v found=%v",s,ok)
	}
}
