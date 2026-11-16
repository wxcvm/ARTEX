package report

import (
	"fmt"
	"strings"

	"github.com/Autumn-27/artex/db"
)

// Summaries count only included findings. Unknown severity/status is explicit,
// not silently treated as informational or confirmed.
func writeFindingSummary(b *strings.Builder,fs []findingView,locale string){
	tr:=func(zh,en string)string{return reportText(locale,zh,en)}
	severity:=map[string]int{}
	status:=map[string]int{}
	for _,f:=range fs{
		severity[severityKey(f.Severity)]++
		key:=f.Status
		if !db.ValidFindingStatus(key){key=""}
		status[key]++
	}
	b.WriteString(tr("### 风险等级分布\n\n","### Severity distribution\n\n"))
	for _,key:=range []string{"critical","high","medium","low","info","unknown"}{
		fmt.Fprintf(b,"- %s: %d\n",reportSeverity(locale,key),severity[key])
	}
	b.WriteString(tr("\n### 处置状态分布\n\n","\n### Finding status distribution\n\n"))
	for _,key:=range []string{db.FindingConfirmed,db.FindingPending,db.FindingInProgress,db.FindingFixed,db.FindingResolved,db.FindingFalsePositive,db.FindingIgnored,db.FindingDuplicate,db.FindingRiskAccepted,""}{
		fmt.Fprintf(b,"- %s: %d\n",taskFindingStatus(locale,key),status[key])
	}
	b.WriteString("\n")
}

// The ledger contains agent-reported observations; it does not prove every
// test was executed. Keep recorded outcomes and their evidence/reason distinct.
func writeCoverageEntries(b *strings.Builder,c *db.CoverageLedger,locale string){
	if len(c.Entries)==0{return}
	tr:=func(zh,en string)string{return reportText(locale,zh,en)}
	b.WriteString(tr("\n### 已记录检查项目\n\n","\n### Recorded checks\n\n"))
	b.WriteString(tr("以下为已保存的检查记录，不代表独立验证。\n\n","The following are saved check records, not independent verification.\n\n"))
	for _,e:=range c.Entries{
		fmt.Fprintf(b,"- %s / %s: %s\n",e.Surface,e.RiskArea,coverageOutcomeLabel(locale,e.Outcome))
		if e.Evidence!=""{fmt.Fprintf(b,"  - %s: %s\n",tr("证据","Evidence"),e.Evidence)}
		if e.Reason!=""{fmt.Fprintf(b,"  - %s: %s\n",tr("原因","Reason"),e.Reason)}
	}
	b.WriteString("\n")
}

func coverageOutcomeLabel(locale,key string)string{
	switch key{
	case db.OutcomeReported:return reportText(locale,"已报告发现","Finding reported")
	case db.OutcomeNoIssueFound:return reportText(locale,"记录为未发现问题","Reported no issue found")
	case db.OutcomeRuledOut:return reportText(locale,"已排除","Ruled out")
	case db.OutcomeNotApplicable:return reportText(locale,"不适用","Not applicable")
	case db.OutcomeNeedsFollowUp:return reportText(locale,"待验证","Needs follow-up")
	default:return reportText(locale,"未知","Unknown")
	}
}
