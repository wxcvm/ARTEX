package report

import (
	"encoding/json"
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
