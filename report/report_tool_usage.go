package report

import (
	"fmt"
	"strings"

	"github.com/Autumn-27/artex/db"
)

func writeReportToolUsage(b *strings.Builder,items []db.ReportToolUsage,unavailable bool,locale string){
	if items==nil&&!unavailable{return}
	tr:=func(zh,tw,en string)string{if locale=="zh-TW"{return tw};return reportText(locale,zh,en)}
	b.WriteString(tr("## 工具使用记录\n\n","## 工具使用記錄\n\n","## Recorded tool usage\n\n"))
	if unavailable{b.WriteString(tr("工具使用记录读取失败。\n\n","工具使用記錄讀取失敗。\n\n","Tool usage records could not be loaded.\n\n"));return}
	b.WriteString(tr("仅统计当前任务记录的调用。返回结果不代表测试成功；未返回结果不代表仍在运行。\n\n","僅統計目前任務記錄的呼叫。回傳結果不代表測試成功；未回傳結果不代表仍在執行。\n\n","Counts reflect calls recorded for this task only. A returned result does not establish test success; a missing result does not establish that a call is still running.\n\n"))
	if len(items)==0{b.WriteString(tr("暂无工具调用记录。\n\n","尚無工具呼叫記錄。\n\n","No tool calls recorded.\n\n"));return}
	for _,s:=range items{
		fmt.Fprintf(b,tr("- %s：调用 %d；已返回 %d；错误 %d；无结果 %d\n","- %s：呼叫 %d；已回傳 %d；錯誤 %d；無結果 %d\n","- %s: calls %d; returned %d; errors %d; no result %d\n"),db.NormalizeCoverageField(s.Tool),s.Total,s.Returned,s.Errors,s.Pending)
	}
	b.WriteString("\n")
}
