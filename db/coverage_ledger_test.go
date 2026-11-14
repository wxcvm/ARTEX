package db

import "testing"

func TestCoverageLedgerGaps(t *testing.T) {
	entries := []CoverageEntry{
		{Surface: "/a", RiskArea: "input", Outcome: OutcomeNoIssueFound},
		{Surface: "/b", RiskArea: "auth", Outcome: OutcomeNeedsFollowUp},
		{Surface: "/c", RiskArea: "config", Outcome: OutcomeRuledOut},
	}
	l := BuildCoverageLedger(entries, []string{"missing", "INPUT"}, CoverageCompletenessFor("timeout", false, false, false))
	if len(l.Gaps) != 4 || len(l.Unexamined) != 1 || l.Unexamined[0] != "missing" || l.Completeness.Complete {
		t.Fatalf("unexpected ledger: %+v", l)
	}
	e, err := normalizeCoverageEntry(CoverageEntry{Surface: " a\nb ", RiskArea: " INPUT ", Outcome: OutcomeReported})
	if err != nil || e.Surface != "a b" || e.RiskArea != "input" {
		t.Fatalf("normalization: %+v, %v", e, err)
	}
}

func TestCoverageLedgerNormalCompletionDoesNotCloseGaps(t *testing.T) {
	l := BuildCoverageLedger(nil, []string{"unexamined"}, CoverageCompletenessFor("done", false, false, false))
	if len(l.Gaps) != 1 || len(l.Unexamined) != 1 {
		t.Fatalf("normal completion must not mark unknown classes clean: %+v", l)
	}
}
