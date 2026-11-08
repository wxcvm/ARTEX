package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
)

func TestFindingsSARIFEmpty(t *testing.T) {
	out, err := FindingsSARIF(nil)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Version string `json:"version"`
		Runs []struct {
			Tool struct {
				Driver struct {
					Name string `json:"name"`
					Rules []json.RawMessage `json:"rules"`
				} `json:"driver"`
			} `json:"tool"`
			Results []json.RawMessage `json:"results"`
		} `json:"runs"`
	}
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Version != "2.1.0" || len(doc.Runs) != 1 ||
		doc.Runs[0].Tool.Driver.Name != "ARTEX" ||
		doc.Runs[0].Results == nil || doc.Runs[0].Tool.Driver.Rules == nil {
		t.Fatalf("invalid empty SARIF envelope: %s", out)
	}
}

func TestFindingsSARIFStableTriageAndSafeEvidence(t *testing.T) {
	task := int64(9)
	now := time.Date(2026, 11, 8, 0, 0, 0, 0, time.UTC)
	findings := []*db.DBFinding{
		{ID: 11, TaskID: &task, Name: "Fixed", VulnClass: "Access Control", Severity: "low", Status: db.FindingFixed, CreatedAt: now.Add(time.Hour), AssetIDs: []int64{4, 1}},
		{ID: 12, Name: "Pending", VulnClass: "Access Control", Severity: "medium", Status: db.FindingPending, CreatedAt: now, Evidence: "PRIVATE-EVIDENCE-DO-NOT-EXPORT"},
		{ID: 13, Name: "Confirmed", VulnClass: "Configuration", Severity: "critical", Status: db.FindingConfirmed, Summary: "Observed problem", CreatedAt: now},
	}
	before := append([]*db.DBFinding(nil), findings...)
	out, err := FindingsSARIF(findings)
	if err != nil {
		t.Fatal(err)
	}
	again, err := FindingsSARIF(findings)
	if err != nil || !bytes.Equal(out, again) {
		t.Fatal("SARIF output must be deterministic")
	}
	for i := range before {
		if findings[i] != before[i] {
			t.Fatal("export mutated caller order")
		}
	}
	if strings.Contains(string(out), "PRIVATE-EVIDENCE-DO-NOT-EXPORT") ||
		strings.Contains(string(out), "physicalLocation") {
		t.Fatalf("export leaked raw evidence or invented source location: %s", out)
	}
	var doc struct {
		Runs []struct {
			Tool struct {
				Driver struct {
					Rules []struct{ ID string `json:"id"` } `json:"rules"`
				} `json:"driver"`
			} `json:"tool"`
			Results []struct {
				RuleID string `json:"ruleId"`
				Level string `json:"level"`
				Kind string `json:"kind"`
				PartialFingerprints map[string]string `json:"partialFingerprints"`
				Properties struct {
					FindingID int64 `json:"findingId"`
					TaskID *int64 `json:"taskId"`
					AssetIDs []int64 `json:"assetIds"`
				} `json:"properties"`
			} `json:"results"`
		} `json:"runs"`
	}
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	r := doc.Runs[0]
	if len(r.Tool.Driver.Rules) != 2 || len(r.Results) != 3 {
		t.Fatalf("rules/results mismatch: %s", out)
	}
	if r.Results[0].Properties.FindingID != 13 || r.Results[0].Kind != "fail" || r.Results[0].Level != "error" {
		t.Fatalf("confirmed critical result not first: %+v", r.Results[0])
	}
	if r.Results[1].Properties.FindingID != 12 || r.Results[1].Kind != "review" || r.Results[1].Level != "warning" {
		t.Fatalf("unreviewed finding misclassified: %+v", r.Results[1])
	}
	last := r.Results[2]
	if last.Kind != "pass" || last.PartialFingerprints["artexFindingId"] != "11" || last.Properties.TaskID == nil ||
		*last.Properties.TaskID != 9 || len(last.Properties.AssetIDs) != 2 ||
		last.Properties.AssetIDs[0] != 1 || last.Properties.AssetIDs[1] != 4 {
		t.Fatalf("resolved metadata missing: %+v", last)
	}
	if r.Results[1].RuleID != last.RuleID || r.Results[0].RuleID == last.RuleID {
		t.Fatal("rule IDs must be stable within a class and different across classes")
	}
}
