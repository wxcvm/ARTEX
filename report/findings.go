package report

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/Autumn-27/artex/db"
)

// Rendering for the findings page "Export": renders a batch of findings table rows
// into a summary Markdown, a single-finding Markdown, or CSV. JSON is serialized
// directly from DTOs at the server layer, not here.

// sortFindingsForExport sorts by severity descending, then by time descending, consistent with the summary report's grouping.
func sortFindingsForExport(fs []*db.DBFinding) {
	sort.SliceStable(fs, func(i, j int) bool {
		ri, rj := sevRank[fs[i].Severity], sevRank[fs[j].Severity]
		if ri != rj {
			return ri < rj // smaller sevRank means more severe
		}
		return fs[i].CreatedAt.After(fs[j].CreatedAt)
	})
}

// findingTitle returns the vulnerability's readable title: name → class → "Uncategorized".
func findingTitle(f *db.DBFinding) string {
	return nz(f.Name, nz(f.VulnClass, "Uncategorized"))
}

// FindingsMarkdown combines a batch of findings into a summary report (summary + grouped by severity,
// each entry with class/status/owning task/evidence/detailed report).
func FindingsMarkdown(fs []*db.DBFinding, generatedAt time.Time) string {
	return FindingsMarkdownLocalized(fs, generatedAt, "en")
}

func FindingsMarkdownLocalized(fs []*db.DBFinding, generatedAt time.Time, locale string) string {
	tr := func(zh, en string) string { return reportText(locale, zh, en) }
	items := append([]*db.DBFinding(nil), fs...)
	sortFindingsForExport(items)

	var b strings.Builder
	b.WriteString(tr("# 漏洞发现汇总报告\n\n", "# Vulnerability Findings Summary Report\n\n"))
	fmt.Fprintf(&b, tr("- **生成时间**：%s\n", "- **Generated at**: %s\n"), generatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, tr("- **发现总数**：%d\n\n", "- **Total findings**: %d\n\n"), len(items))

	// Summary: counts per severity.
	counts := map[string]int{}
	for _, f := range items {
		counts[f.Severity]++
	}
	b.WriteString(tr("## 摘要\n\n", "## Summary\n\n"))
	b.WriteString(tr("| 严重程度 | 数量 |\n| --- | --- |\n", "| Severity | Count |\n| --- | --- |\n"))
	for _, s := range []struct{ key, label string }{
		{"critical", "Critical"}, {"high", "High"}, {"medium", "Medium"}, {"low", "Low"},
	} {
		fmt.Fprintf(&b, "| %s | %d |\n", reportText(locale, reportSeverity(locale, s.key), s.label), counts[s.key])
	}
	b.WriteString("\n")

	if len(items) == 0 {
		b.WriteString(tr("_没有符合条件的漏洞。_\n", "_No matching vulnerabilities._\n"))
		return b.String()
	}

	b.WriteString(tr("## 漏洞详情\n\n", "## Vulnerability Details\n\n"))
	for i, f := range items {
		fmt.Fprintf(&b, "### %d. [%s] %s\n\n", i+1, reportSeverity(locale, nz(f.Severity, "info")), nz(f.Name, nz(f.VulnClass, tr("未分类", "Uncategorized"))))
		if f.VulnClass != "" {
			fmt.Fprintf(&b, tr("- **类型**：%s\n", "- **Class**: %s\n"), f.VulnClass)
		}
		fmt.Fprintf(&b, tr("- **状态**：%s\n", "- **Status**: %s\n"), nz(f.Status, "pending"))
		if desc := strings.TrimSpace(f.TaskDescription); desc != "" {
			fmt.Fprintf(&b, tr("- **任务**：%s\n", "- **Task**: %s\n"), desc)
		}
		fmt.Fprintf(&b, tr("- **发现时间**：%s\n\n", "- **Found at**: %s\n\n"), f.CreatedAt.Format("2006-01-02 15:04:05"))
		if s := strings.TrimSpace(f.Summary); s != "" {
			fmt.Fprintf(&b, "%s\n\n", s)
		}
		if e := strings.TrimSpace(f.Evidence); e != "" {
			fmt.Fprintf(&b, tr("**证据：**\n\n```\n%s\n```\n\n", "**Evidence:**\n\n```\n%s\n```\n\n"), e)
		}
		if rep := strings.TrimSpace(f.Report); rep != "" {
			b.WriteString(tr("**详细报告：**\n\n", "**Detailed report:**\n\n"))
			b.WriteString(rep)
			b.WriteString("\n\n")
		}
		b.WriteString(findingTrafficMarkdownLocalized(f, false, locale))
		b.WriteString("---\n\n")
	}
	return b.String()
}

// SingleFindingMarkdown renders a single vulnerability as a standalone Markdown (used for "one file per vulnerability" packaging).
func SingleFindingMarkdown(f *db.DBFinding, generatedAt time.Time) string {
	return SingleFindingMarkdownLocalized(f, generatedAt, "en")
}

func SingleFindingMarkdownLocalized(f *db.DBFinding, generatedAt time.Time, locale string) string {
	tr := func(zh, en string) string { return reportText(locale, zh, en) }
	var b strings.Builder
	fmt.Fprintf(&b, "# [%s] %s\n\n", reportSeverity(locale, nz(f.Severity, "info")), nz(f.Name, nz(f.VulnClass, tr("未分类", "Uncategorized"))))
	if f.VulnClass != "" {
		fmt.Fprintf(&b, tr("- **类型**：%s\n", "- **Class**: %s\n"), f.VulnClass)
	}
	fmt.Fprintf(&b, tr("- **严重程度**：%s\n", "- **Severity**: %s\n"), nz(f.Severity, "info"))
	fmt.Fprintf(&b, tr("- **状态**：%s\n", "- **Status**: %s\n"), nz(f.Status, "pending"))
	if desc := strings.TrimSpace(f.TaskDescription); desc != "" {
		fmt.Fprintf(&b, tr("- **任务**：%s\n", "- **Task**: %s\n"), desc)
	}
	fmt.Fprintf(&b, tr("- **发现时间**：%s\n", "- **Found at**: %s\n"), f.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, tr("- **生成时间**：%s\n\n", "- **Generated at**: %s\n\n"), generatedAt.Format("2006-01-02 15:04:05"))
	if s := strings.TrimSpace(f.Summary); s != "" {
		fmt.Fprintf(&b, tr("## 摘要\n\n%s\n\n", "## Summary\n\n%s\n\n"), s)
	}
	if e := strings.TrimSpace(f.Evidence); e != "" {
		fmt.Fprintf(&b, tr("## 证据\n\n```\n%s\n```\n\n", "## Evidence\n\n```\n%s\n```\n\n"), e)
	}
	if rep := strings.TrimSpace(f.Report); rep != "" {
		b.WriteString(tr("## 详细报告\n\n", "## Detailed Report\n\n"))
		b.WriteString(rep)
		b.WriteString("\n")
	}
	b.WriteString(findingTrafficMarkdownLocalized(f, true, locale))
	return b.String()
}

var unsafeFilenameChars = regexp.MustCompile(`[^\p{Han}\p{L}\p{N}._-]+`)

// FindingFilename generates a safe .md filename for "one file per vulnerability", e.g.
// `critical_SQL_injection_#123.md`. Strips path separators and control characters to avoid illegal paths inside the zip.
func FindingFilename(f *db.DBFinding) string {
	sev := nz(f.Severity, "info")
	title := findingTitle(f)
	name := fmt.Sprintf("%s_%s_#%d", sev, title, f.ID)
	name = unsafeFilenameChars.ReplaceAllString(name, "_")
	name = strings.Trim(name, "._")
	if name == "" {
		name = fmt.Sprintf("finding_%d", f.ID)
	}
	// Defensive: strip one more path layer to rule out zip slip.
	name = path.Base(name)
	if len(name) > 120 {
		name = name[:120]
	}
	return name + ".md"
}

// safeCSVCell neutralizes formulas in untrusted text when a CSV file is opened
// by spreadsheet software. csv.Writer handles quoting and separators; this
// protects the cell value itself, including formulas after leading whitespace.
func safeCSVCell(value string) string {
	leading := strings.TrimLeftFunc(value, func(r rune) bool {
		return unicode.IsSpace(r) || r == '\uFEFF' || r == '\u200B'
	})
	if leading != "" {
		switch leading[0] {
		case '=', '+', '-', '@':
			return "'" + value
		}
	}
	return value
}

// FindingsCSV renders a batch of findings as CSV (with a UTF-8 BOM so Excel detects the encoding correctly).
// It omits the long report/evidence full text and keeps only summary-level fields; use Markdown/JSON export for full text.
func FindingsCSV(fs []*db.DBFinding) []byte {
	items := append([]*db.DBFinding(nil), fs...)
	sortFindingsForExport(items)

	var buf bytes.Buffer
	buf.WriteString("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"ID", "Name", "Class", "Severity", "Status", "Task", "Found At", "Summary", "Traffic Evidence Count", "Traffic Evidence IDs"})
	for _, f := range items {
		_ = w.Write([]string{
			fmt.Sprintf("%d", f.ID),
			safeCSVCell(findingTitle(f)),
			safeCSVCell(f.VulnClass),
			nz(f.Severity, "info"),
			nz(f.Status, "pending"),
			safeCSVCell(f.TaskDescription),
			f.CreatedAt.Format("2006-01-02 15:04:05"),
			safeCSVCell(strings.TrimSpace(f.Summary)),
			fmt.Sprint(len(f.TrafficBindings)), findingTrafficIDs(f),
		})
	}
	w.Flush()
	return buf.Bytes()
}

func findingTrafficIDs(f *db.DBFinding) string {
	ids := make([]string, 0, len(f.TrafficBindings))
	for _, b := range f.TrafficBindings {
		ids = append(ids, fmt.Sprint(b.ID))
	}
	return strings.Join(ids, ",")
}

func findingTrafficMarkdown(f *db.DBFinding, attachments bool) string {
	return findingTrafficMarkdownLocalized(f, attachments, "en")
}

func findingTrafficMarkdownLocalized(f *db.DBFinding, attachments bool, locale string) string {
	tr := func(zh, en string) string { return reportText(locale, zh, en) }
	stale := f.Report != "" && f.EvidenceVersion != f.ReportEvidenceVersion
	if len(f.TrafficBindings) == 0 && !stale {
		return ""
	}
	var out strings.Builder
	out.WriteString(tr("\n## 关联流量证据\n\n", "\n## Related Traffic Evidence\n\n"))
	fmt.Fprintf(&out, tr("证据版本：%d；关联数量：%d。\n\n", "Evidence version: %d; bindings: %d.\n\n"), f.EvidenceVersion, len(f.TrafficBindings))
	if stale {
		out.WriteString(tr("证据已变更，详细报告需要更新。\n\n", "The evidence has changed; the detailed report needs updating.\n\n"))
	}
	for i, b := range f.TrafficBindings {
		fmt.Fprintf(&out, tr("%d. **证据 #%d · %s** — `%s %s`，状态码 %d\n", "%d. **Evidence #%d · %s** — `%s %s`, status %d\n"), i+1, b.ID, b.Role, b.Snapshot.Method, strings.ReplaceAll(b.Snapshot.URL, "`", "%60"), b.Snapshot.Status)
		if b.Note != "" {
			fmt.Fprintf(&out, "   %s\n", strings.ReplaceAll(b.Note, "\n", "\n   "))
		}
		if attachments {
			fmt.Fprintf(&out, tr("   [请求报文](evidence/%d/%d/request.http) · [响应报文](evidence/%d/%d/response.http)\n", "   [Request message](evidence/%d/%d/request.http) · [Response message](evidence/%d/%d/response.http)\n"), f.ID, b.ID, f.ID, b.ID)
		}
	}
	out.WriteString("\n")
	return out.String()
}
