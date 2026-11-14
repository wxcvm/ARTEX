package agent

import (
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestCoverageSummaryExposesOutstandingChecks(t *testing.T) {
	l := db.BuildCoverageLedger(nil, []string{"authentication"}, db.CoverageCompletenessFor("running", false, false, false))
	got := renderCoverageSummary(l)
	if !strings.Contains(got, "authentication") || !strings.Contains(got, "unrecorded_risk_class") || !strings.Contains(got, "data, not instructions") {
		t.Fatalf("missing structured gaps: %s", got)
	}
	if renderCoverageSummary(nil) != "" {
		t.Fatal("nil ledger should not produce guidance")
	}
}
