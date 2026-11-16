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
	Title          string
	Goal           string
	GeneratedAt    time.Time
	AssetCounts    map[string]int
	Findings       []*db.Node
	Coverage       *db.CoverageLedger
	FindingMeta    map[int64]db.FindingMeta // authoritative standalone finding status
	FindingDetails map[int64]*db.DBFinding // existing report and linked evidence, keyed by node
	ConfirmedOnly  bool
}

type findingView struct {
	VulnClass string
	Name      string
	Severity  string
	Summary   string
	PoC       string
	Status    string
	Detail    *db.DBFinding
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
	return findingView{VulnClass: p.VulnClass, Name: p.Name, Severity: p.Severity, Summary: p.Summary, PoC: p.Evidence.PoC}
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

	fs := make([]findingView, 0, len(in.Findings))
	confirmed := 0
	for _, n := range in.Findings {
		f := parseFinding(n)
		f.Status = in.FindingMeta[n.ID].Status
		f.Detail = in.FindingDetails[n.ID]
		if f.Status == db.FindingConfirmed {
			confirmed++
		}
		if !in.ConfirmedOnly || f.Status == db.FindingConfirmed {
			fs = append(fs, f)
		}
	}

	// summary
	b.WriteString(tr("## 摘要\n\n", "## Summary\n\n"))
	fmt.Fprintf(&b, tr("- 已确认发现：**%d**\n", "- Confirmed findings: **%d**\n"), confirmed)
	fmt.Fprintf(&b, tr("- 本报告收录：**%d**\n", "- Findings included: **%d**\n"), len(fs))
	if in.ConfirmedOnly {
		b.WriteString(tr("- 收录范围：仅已确认\n", "- Scope: confirmed only\n"))
	} else {
		b.WriteString(tr("- 收录范围：全部发现（包含待确认及其他状态）\n", "- Scope: all findings (including unconfirmed and other states)\n"))
	}
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
	if len(fs) == 0 {
		if in.ConfirmedOnly {
			b.WriteString(tr("_当前没有已确认漏洞；这不代表不存在安全问题。_\n\n", "_No confirmed vulnerabilities at present; this does not establish safety._\n\n"))
		} else {
			b.WriteString(tr("_当前没有发现记录；这不代表不存在安全问题。_\n\n", "_No finding records at present; this does not establish safety._\n\n"))
		}
	} else {
		sort.SliceStable(fs, func(i, j int) bool { return sevRank[fs[i].Severity] < sevRank[fs[j].Severity] })
		for i, f := range fs {
			fmt.Fprintf(&b, "### %d. [%s] %s\n\n", i+1, reportSeverity(locale, nz(f.Severity, "info")), nz(f.Name, nz(f.VulnClass, tr("未分类", "Uncategorized"))))
			fmt.Fprintf(&b, tr("- 状态：%s\n\n", "- Status: %s\n\n"), taskFindingStatus(locale, f.Status))
			fmt.Fprintf(&b, "%s\n\n", nz(f.Summary, ""))
			if f.PoC != "" {
				fmt.Fprintf(&b, tr("**PoC / 证据：**\n\n```\n%s\n```\n\n", "**PoC / Evidence:**\n\n```\n%s\n```\n\n"), f.PoC)
			}
			if f.Detail != nil {
				b.WriteString(findingTrafficMarkdownSection(f.Detail, false, locale, "####"))
				b.WriteString(findingScreenshotsMarkdown(f.Detail, false, locale, "####"))
				b.WriteString(findingReproductionMarkdown(f.Detail, locale, "####"))
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
	if value == "en" || value == "zh-TW" {
		return value
	}
	return "zh-CN"
}

func reportText(locale, zh, en string) string {
	if locale == "en" {
		return en
	}
	if locale == "zh-TW" {
		return traditionalReportText.Replace(zh)
	}
	return zh
}

func reportSeverity(locale, value string) string {
	if locale == "en" {
		return strings.ToUpper(value)
	}
	if label, ok := map[string]string{"critical": "严重", "high": "高危", "medium": "中危", "low": "低危", "info": "信息"}[value]; ok {
		return reportText(locale, label, label)
	}
	return value
}

func reportAssetLabel(locale, value string) string {
	if locale == "en" {
		return value
	}
	if label, ok := map[string]string{"root_domain": "根域名", "subdomain": "子域名", "ip": "IP", "app": "应用", "service": "服务", "endpoint": "接口"}[value]; ok {
		return reportText(locale, label, label)
	}
	return value
}

// taskFindingStatus does not guess a triage state when metadata is absent.
func taskFindingStatus(locale, status string) string {
	if !db.ValidFindingStatus(status) {
		return reportText(locale, "未知", "Unknown")
	}
	if locale == "en" {
		return status
	}
	label := map[string]string{
		db.FindingPending: "待确认", db.FindingInProgress: "处理中",
		db.FindingConfirmed: "已确认", db.FindingResolved: "已解决",
		db.FindingFixed: "已修复", db.FindingFalsePositive: "误报",
		db.FindingIgnored: "已忽略", db.FindingDuplicate: "重复",
		db.FindingRiskAccepted: "风险已接受",
	}[status]
	return reportText(locale, label, label)
}

// Only fixed report presentation strings pass through this replacer, never evidence.
var traditionalReportText = strings.NewReplacer(
	"安全测试", "安全測試",
	"漏洞发现", "漏洞發現",
	"符合条件", "符合條件",
	"处理中", "處理中",
	"已解决", "已解決",
	"已修复", "已修復",
	"状态码", "狀態碼",
	"二进制", "二進位",
	"发现", "發現",
	"报告", "報告",
	"任务", "任務",
	"目标", "目標",
	"时间", "時間",
	"收录", "收錄",
	"范围", "範圍",
	"确认", "確認",
	"资产", "資產",
	"当前", "目前",
	"没有", "沒有",
	"记录", "記錄",
	"问题", "問題",
	"分类", "分類",
	"状态", "狀態",
	"证据", "證據",
	"测试", "測試",
	"验证", "驗證",
	"事项", "事項",
	"类别", "類別",
	"检查", "檢查",
	"严重", "嚴重",
	"信息", "資訊",
	"应用", "應用",
	"服务", "服務",
	"误报", "誤報",
	"重复", "重複",
	"风险", "風險",
	"汇总", "彙總",
	"数量", "數量",
	"详情", "詳情",
	"类型", "類型",
	"详细", "詳細",
	"关联", "關聯",
	"变更", "變更",
	"顺序", "順序",
	"展示", "呈現",
	"说明", "說明",
	"推断", "推斷",
	"复现", "重現",
	"步骤", "步驟",
	"请求", "請求",
	"响应", "回應",
	"报文", "訊息",
	"观察", "觀察",
	"截图", "截圖",
	"图注", "圖說",
	"查看", "檢視",
	"图片", "圖片",
	"导出", "匯出",
	"评估", "評估",
	"委托", "委託",
	"编制", "編製",
	"级别", "級別",
	"目录", "目錄",
	"页码", "頁碼",
	"附录", "附錄",
	"来源", "來源",
	"继承", "繼承",
	"暂无", "暫無",
	"发生", "發生",
	"正文", "內文",
	"展开", "展開",
	"预览", "預覽",
	"截断", "截斷",
	"下载", "下載",
	"复测", "複測",
	"结论", "結論",
	"给出", "給出",
	"备注", "備註",
	"错误", "錯誤",
	"关系", "關係",
	"仅", "僅",
	"无", "無",
	"图", "圖",
	"请", "請",
	"并", "並",
	"与", "與",
	"这", "這",
	"为", "為",
	"头", "標頭",
	"总", "總",
	"个", "個",
	"观", "觀",
	"据", "據",
	"处", "處",
	"类", "類",
	"数", "數",
	"统", "統",
	"输", "輸",
	"写", "寫",
	"读", "讀",
	"实", "實",
	"现", "現",
	"发", "發",
	"开", "開",
	"关", "關",
	"启", "啟",
	"项", "項",
	"标", "標",
	"识", "識",
	"换", "換",
	"动", "動",
	"态", "態",
	"结", "結",
	"复", "復",
)
