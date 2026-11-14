// Package report renders a deliverable Markdown report from the exploration
// graph: confirmed findings (with severity/evidence) (docs §15 P5).
package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Autumn-27/artex/db"
)

// Input bundles what the report needs.
type Input struct {
	Title       string
	Goal        string
	GeneratedAt time.Time
	AssetCounts map[string]int
	Findings    []*db.Node
	Coverage    *db.CoverageLedger
}

type findingView struct {
	VulnClass string
	Name      string
	Severity  string
	Summary   string
	PoC       string
}

func parseFinding(n *db.Node) findingView {
	var p struct {
		VulnClass string `json:"vulnclass"`
		Name      string `json:"name"`
		Severity  string `json:"severity"`
		Summary   string `json:"summary"`
		Evidence  struct {
			PoC string `json:"poc"`
		} `json:"evidence"`
	}
	_ = json.Unmarshal(n.Payload, &p)
	return findingView{p.VulnClass, p.Name, p.Severity, p.Summary, p.Evidence.PoC}
}

var sevRank = map[string]int{"critical": 0, "high": 1, "medium": 2, "low": 3, "": 4}

// Markdown retains the original English rendering for existing callers.
func Markdown(in Input) string {
	return MarkdownLocalized(in, "en")
}

// MarkdownLocalized translates the report structure, not stored evidence.
func MarkdownLocalized(in Input, locale string) string {
	tr := func(zh, en string) string { return reportText(locale, zh, en) }
	var b strings.Builder
	fmt.Fprintf(&b, tr("# 安全测试报告 — %s\n\n", "# Penetration Test Report — %s\n\n"), nz(in.Title, tr("未命名任务", "Unnamed task")))
	fmt.Fprintf(&b, tr("- **任务目标**：%s\n", "- **Task goal**: %s\n"), nz(in.Goal, tr("（未指定）", "(unspecified)")))
	fmt.Fprintf(&b, tr("- **生成时间**：%s\n\n", "- **Generated at**: %s\n\n"), in.GeneratedAt.Format("2006-01-02 15:04:05"))

	// summary
	b.WriteString(tr("## 摘要\n\n", "## Summary\n\n"))
	fmt.Fprintf(&b, tr("- 已确认发现：**%d**\n", "- Confirmed findings: **%d**\n"), len(in.Findings))
	b.WriteString(tr("- 资产：", "- Assets: "))
	var types []string
	for t := range in.AssetCounts {
		types = append(types, t)
	}
	sort.Strings(types)
	for i, t := range types {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%s %d", reportAssetLabel(locale, t), in.AssetCounts[t])
	}
	b.WriteString("\n\n")

	// findings
	b.WriteString(tr("## 发现结果\n\n", "## Findings\n\n"))
	if len(in.Findings) == 0 {
		b.WriteString(tr("_本次未确认漏洞；这不代表不存在安全问题。_\n\n", "_No vulnerabilities confirmed this time._\n\n"))
	} else {
		fs := make([]findingView, 0, len(in.Findings))
		for _, n := range in.Findings {
			fs = append(fs, parseFinding(n))
		}
		sort.SliceStable(fs, func(i, j int) bool { return sevRank[fs[i].Severity] < sevRank[fs[j].Severity] })
		for i, f := range fs {
			fmt.Fprintf(&b, "### %d. [%s] %s\n\n", i+1, reportSeverity(locale, nz(f.Severity, "info")), nz(f.Name, nz(f.VulnClass, tr("未分类", "Uncategorized"))))
			fmt.Fprintf(&b, "%s\n\n", nz(f.Summary, ""))
			if f.PoC != "" {
				fmt.Fprintf(&b, tr("**PoC / 证据：**\n\n```\n%s\n```\n\n", "**PoC / Evidence:**\n\n```\n%s\n```\n\n"), f.PoC)
			}
		}
	}

	if c := in.Coverage; c != nil {
		b.WriteString(tr("## 测试覆盖与待验证事项\n\n", "## Coverage and open checks\n\n"))
		fmt.Fprintf(&b, tr("- 记录：%d；待验证：%d\n", "- Recorded: %d; open gaps: %d\n"), len(c.Entries), len(c.Gaps))
		b.WriteString(tr("- 未发现问题不等于目标安全；未检查类别仍未验证。\n", "- No finding does not establish safety; unexamined classes remain unverified.\n"))
		if c.Completeness != nil && !c.Completeness.Complete {
			b.WriteString(tr("- 本次任务未正常完成。\n", "- This run did not complete normally.\n"))
		}
		for _, g := range c.Gaps {
			fmt.Fprintf(&b, "- %s\n", g)
		}
		b.WriteString("\n")
	}
	return b.String()
}

func nz(s, d string) string {
	if strings.TrimSpace(s) == "" {
		return d
	}
	return s
}

// ReportLocale normalizes the requested report language; Chinese is the default.
func ReportLocale(value string) string {
	if value == "en" {
		return "en"
	}
	return "zh-CN"
}

func reportText(locale, zh, en string) string {
	if locale == "en" {
		return en
	}
	return zh
}

func reportSeverity(locale, value string) string {
	if locale == "en" {
		return strings.ToUpper(value)
	}
	if label, ok := map[string]string{"critical": "严重", "high": "高危", "medium": "中危", "low": "低危", "info": "信息"}[value]; ok {
		return label
	}
	return value
}

func reportAssetLabel(locale, value string) string {
	if locale == "en" {
		return value
	}
	if label, ok := map[string]string{"root_domain": "根域名", "subdomain": "子域名", "ip": "IP", "app": "应用", "service": "服务", "endpoint": "接口"}[value]; ok {
		return label
	}
	return value
}
