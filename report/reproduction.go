package report

import (
	"fmt"
	"strings"

	"github.com/Autumn-27/artex/db"
)

func findingReproductionMarkdown(f *db.DBFinding, locale, heading string) string {
	if len(f.ReproductionSteps) == 0 {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "\n%s %s\n\n", heading, reportText(locale, "人工复现记录", "Recorded reproduction steps"))
	for i, step := range f.ReproductionSteps {
		fmt.Fprintf(&b, "%d. %s\n", i+1, strings.Join(strings.Fields(step.Title), " "))
		if step.Observation != "" {
			fmt.Fprintf(&b, "   %s: %s\n", reportText(locale, "观察结果", "Observation"), strings.ReplaceAll(step.Observation, "\n", "\n   "))
		}
		if step.BindingID != "" {
			fmt.Fprintf(&b, "   %s #%s\n", reportText(locale, "关联流量证据", "Linked traffic evidence"), step.BindingID)
		}
	}
	b.WriteString("\n")
	return b.String()
}
