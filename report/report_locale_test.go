package report

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
)

func TestLocalizedTaskMarkdown(t *testing.T) {
	payload := json.RawMessage(`{"name":"Original finding","severity":"high","summary":"Original summary","evidence":{"poc":"Original evidence"}}`)
	in := Input{Title: "Original title", Goal: "Original goal", GeneratedAt: time.Unix(0, 0), AssetCounts: map[string]int{"root_domain": 1}, Findings: []*db.Node{{Payload: payload}}}
	for _, lang := range []string{"zh-CN", "en"} {
		md := MarkdownLocalized(in, lang)
		for _, value := range []string{"Original title", "Original goal", "Original finding", "Original summary", "Original evidence"} {
			if !strings.Contains(md, value) {
				t.Fatalf("%s lost original content %q", lang, value)
			}
		}
		if lang == "zh-CN" && (!strings.Contains(md, "## 发现结果") || !strings.Contains(md, "高危") || !strings.Contains(md, "根域名")) {
			t.Fatal("missing Chinese headings or labels")
		}
		if lang == "en" && !strings.Contains(md, "## Findings") {
			t.Fatal("missing English heading")
		}
	}
	if ReportLocale("") != "zh-CN" || ReportLocale("en") != "en" {
		t.Fatal("unexpected default locale")
	}
}

func TestLocalizedFindingExportsPreserveEvidence(t *testing.T) {
	f := &db.DBFinding{ID: 17, Name: "Original finding", Severity: "high", Status: "confirmed", Evidence: "Header: original", Report: "## Original report\nOriginal text", Summary: "Original summary"}
	for _, lang := range []string{"zh-CN", "en"} {
		for _, md := range []string{FindingsMarkdownLocalized([]*db.DBFinding{f}, time.Unix(0, 0), lang), SingleFindingMarkdownLocalized(f, time.Unix(0, 0), lang)} {
			for _, value := range []string{"Original finding", "Original summary", "Header: original", "## Original report\nOriginal text"} {
				if !strings.Contains(md, value) {
					t.Fatalf("%s lost original content %q", lang, value)
				}
			}
			if lang == "zh-CN" && !strings.Contains(md, "详细报告") {
				t.Fatal("missing Chinese report label")
			}
		}
	}
}

func TestTaskReportFindingScope(t *testing.T) {
	in := Input{
		Findings: []*db.Node{
			{ID: 1, Payload: json.RawMessage(`{"name":"confirmed sample","severity":"high"}`)},
			{ID: 2, Payload: json.RawMessage(`{"name":"pending sample","severity":"low"}`)},
			{ID: 3, Payload: json.RawMessage(`{"name":"false positive sample","severity":"medium"}`)},
			{ID: 4, Payload: json.RawMessage(`{"name":"unknown sample","severity":"low"}`)},
		},
		FindingMeta: map[int64]db.FindingMeta{
			1: {Status: db.FindingConfirmed},
			2: {Status: db.FindingPending},
			3: {Status: db.FindingFalsePositive},
		},
	}
	all := MarkdownLocalized(in, "en")
	for _, want := range []string{"Confirmed findings: **1**", "Findings included: **4**", "confirmed sample", "pending sample", "false positive sample", "unknown sample", "Status: Unknown"} {
		if !strings.Contains(all, want) {
			t.Fatalf("all report missing %q", want)
		}
	}
	in.ConfirmedOnly = true
	for _, locale := range []string{"en", "zh-CN"} {
		md := MarkdownLocalized(in, locale)
		if !strings.Contains(md, "confirmed sample") {
			t.Fatal("confirmed finding missing")
		}
		for _, excluded := range []string{"pending sample", "false positive sample", "unknown sample"} {
			if strings.Contains(md, excluded) {
				t.Fatalf("confirmed report includes %q", excluded)
			}
		}
	}
}

func TestTaskReportLinkedEvidence(t *testing.T) {
	in := Input{
		Findings: []*db.Node{
			{ID: 1, Payload: json.RawMessage(`{"name":"included sample","severity":"high"}`)},
			{ID: 2, Payload: json.RawMessage(`{"name":"excluded sample","severity":"low"}`)},
		},
		FindingMeta: map[int64]db.FindingMeta{
			1: {Status: db.FindingConfirmed},
			2: {Status: db.FindingPending},
		},
		ConfirmedOnly: true,
		FindingDetails: map[int64]*db.DBFinding{
			1: {ID: 11, EvidenceVersion: 2, ReportEvidenceVersion: 1, Report: "previous report", TrafficBindings: []db.FindingTrafficBinding{
				{ID: 9, Role: "baseline", Note: "first saved note", Snapshot: db.TrafficEvidenceSnapshot{Method: "GET", URL: "http://example.test/", Status: 200}},
				{ID: 8, Role: "supporting", Note: "second saved note"},
			}},
			2: {ID: 12, TrafficBindings: []db.FindingTrafficBinding{{ID: 10, Note: "excluded evidence"}}},
		},
	}
	md := MarkdownLocalized(in, "en")
	for _, want := range []string{"#### Related Traffic Evidence", "Evidence #9", "first saved note", "second saved note", "The evidence has changed", "not inferred reproduction steps"} {
		if !strings.Contains(md, want) {
			t.Fatalf("missing %q", want)
		}
	}
	if strings.Contains(md, "excluded evidence") || strings.Index(md, "first saved note") > strings.Index(md, "second saved note") {
		t.Fatal("wrong evidence scope or order")
	}
}

func TestTraditionalReportPreservesOriginalEvidence(t *testing.T) {
	if ReportLocale("zh-TW") != "zh-TW" { t.Fatal("locale lost") }
	f := &db.DBFinding{ID: 1, Name: "原始发现", Severity: "high", Status: "confirmed", Evidence: "原始证据", Report: "原始详细报告"}
	md := SingleFindingMarkdownLocalized(f, time.Unix(0,0), "zh-TW")
	for _, want := range []string{"## 詳細報告", "## 證據", "原始证据", "原始详细报告", "原始发现"} {
		if !strings.Contains(md,want) { t.Fatalf("missing %q",want) }
	}
	in:=Input{Title:"原始任务",Goal:"原始目标",AssetCounts:map[string]int{"app":1}}
	task:=MarkdownLocalized(in,"zh-TW")
	for _,want:=range []string{"安全測試報告","任務目標","應用","原始任务","原始目标"} {if !strings.Contains(task,want){t.Fatalf("missing %q",want)}}
	data,err:=DOCXWithLayout(task,nil,"zh-TW","professional")
	if err!=nil{t.Fatal(err)}
	z,err:=zip.NewReader(bytes.NewReader(data),int64(len(data)))
	if err!=nil{t.Fatal(err)}
	for _,file:=range z.File{if file.Name=="word/document.xml"{r,err:=file.Open();if err!=nil{t.Fatal(err)};raw,err:=io.ReadAll(r);r.Close();if err!=nil{t.Fatal(err)};for _,want:=range []string{"目錄","頁碼","安全測試報告","原始任务"}{if !strings.Contains(string(raw),want){t.Fatalf("docx missing %q",want)}};return}}
	t.Fatal("missing document.xml")
}
