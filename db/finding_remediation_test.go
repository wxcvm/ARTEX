package db

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func remediationVersion(n int64) *int64 { return &n }

func TestRemediationDatabaseWorkflow(t *testing.T) {
	d, err := Open(testDSN(t))
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	ctx := t.Context()
	title := fmt.Sprintf("Remediation QA %d", time.Now().UnixNano())
	finding, err := d.AddFinding(0, 0, "qa", title, SeverityHigh, "original summary", "original material", "qa", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Exec("DELETE FROM findings WHERE id=$1", finding)
	id := fmt.Sprint(finding)
	read := func() FindingRemediationDetail {
		t.Helper()
		current, err := d.GetFindingRemediation(ctx, id, 0)
		if err != nil {
			t.Fatal(err)
		}
		return current
	}
	if current := read(); current.Item.Stage != remediationReview || current.Item.Version != 0 || len(current.Events) != 0 {
		t.Fatalf("implicit record: %+v", current)
	}
	items, total, err := d.ListFindingRemediations(ctx, RemediationFilter{Stage: remediationReview, Severity: SeverityHigh, Query: strings.ToLower(title)})
	if err != nil || total != 1 || len(items) != 1 || items[0].FindingID != id {
		t.Fatalf("query: %+v total=%d (%v)", items, total, err)
	}
	change := func(action string, version int64) RemediationTransition {
		return RemediationTransition{Action: action, ExpectedVersion: remediationVersion(version), Actor: "ARTEX", Note: "Human review"}
	}
	apply := func(input RemediationTransition, stage string) FindingRemediationDetail {
		t.Helper()
		detail, err := d.TransitionFindingRemediation(ctx, id, input)
		if err != nil || detail.Item.Stage != stage || detail.Item.Version != *input.ExpectedVersion+1 || len(detail.Events) == 0 || detail.Events[0].Action != input.Action {
			t.Fatalf("%s: %+v (%v)", input.Action, detail, err)
		}
		return detail
	}

	// Two first writes with the same version may not create duplicate events.
	start, results := make(chan struct{}), make(chan error, 2)
	for range 2 {
		go func() {
			<-start
			_, err := d.TransitionFindingRemediation(ctx, id, change("review", 0))
			results <- err
		}()
	}
	close(start)
	a, b := <-results, <-results
	if !((a == nil && errors.Is(b, ErrRemediationConflict)) || (b == nil && errors.Is(a, ErrRemediationConflict))) {
		t.Fatalf("expected one first-write conflict: %v, %v", a, b)
	}

	plan := change("update_plan", 1)
	plan.Plan, plan.Owner, plan.DueDate = "Upgrade the dependency", "维护人", "2027-01-25"
	current := apply(plan, remediationFix)
	if current.Item.Owner != plan.Owner || current.Events[0].Plan != plan.Plan || current.Item.DueDate != plan.DueDate {
		t.Fatalf("plan and history diverged: %+v", current)
	}
	submit := change("submit", 2)
	submit.Reference, submit.ExpectedEvidenceVersion = "QA-001", remediationVersion(0)
	apply(submit, remediationAcceptance)
	advanceEvidence := func(version int64) {
		t.Helper()
		if _, err := d.Exec("UPDATE findings SET evidence_version=$2 WHERE id=$1", finding, version); err != nil {
			t.Fatal(err)
		}
	}
	advanceEvidence(1)
	accept := change("accept", 3)
	for _, evidence := range []int64{0, 1} {
		accept.ExpectedEvidenceVersion = remediationVersion(evidence)
		if _, err := d.TransitionFindingRemediation(ctx, id, accept); !errors.Is(err, ErrRemediationConflict) {
			t.Fatalf("stale evidence %d accepted, got %v", evidence, err)
		}
	}
	if current := read(); !current.Item.NeedsReview || len(current.Events) != 3 {
		t.Fatalf("conflict must not write: %+v", current)
	}
	apply(change("return", 3), remediationFix)
	submit.ExpectedVersion, submit.ExpectedEvidenceVersion = remediationVersion(4), remediationVersion(1)
	apply(submit, remediationAcceptance)
	accept.ExpectedVersion, accept.ExpectedEvidenceVersion = remediationVersion(5), remediationVersion(1)
	current = apply(accept, remediationAccepted)
	if current.Item.NeedsReview || !sameEvidence(current.Item.AcceptedEvidenceVersion, 1) {
		t.Fatalf("fresh acceptance: %+v", current.Item)
	}
	advanceEvidence(2)
	if current := read(); !current.Item.NeedsReview || current.Item.Stage != remediationAccepted {
		t.Fatalf("later evidence cannot silently inherit acceptance: %+v", current)
	}
	apply(change("reopen", 6), remediationFix)

	for _, disposition := range []string{"false_positive", "ignored", "duplicate", "risk_accepted"} {
		if _, err := d.SetFindingStatus(finding, disposition); err != nil {
			t.Fatal(err)
		}
		items, total, err := d.ListFindingRemediations(ctx, RemediationFilter{Query: id})
		if err != nil || total != 0 || len(items) != 0 {
			t.Fatalf("excluded disposition %s: %+v total=%d (%v)", disposition, items, total, err)
		}
		if current := read(); !current.Item.ReadOnly || len(current.Events) != 7 {
			t.Fatalf("excluded detail: %+v", current)
		}
		plan.ExpectedVersion = remediationVersion(7)
		if _, err := d.TransitionFindingRemediation(ctx, id, plan); !errors.Is(err, ErrRemediationReadOnly) {
			t.Fatalf("disposition %s still writable: %v", disposition, err)
		}
	}
	var summary, evidence string
	if err := d.QueryRow("SELECT summary,evidence FROM findings WHERE id=$1", finding).Scan(&summary, &evidence); err != nil ||
		summary != "original summary" || evidence != "original material" {
		t.Fatalf("remediation changed finding data (%v)", err)
	}
}
