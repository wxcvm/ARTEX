package db

import (
	"fmt"
	"testing"
	"time"
)

func TestHoneypotAssessmentRoundTripAndExpiry(t *testing.T) {
	d, assets, _ := testSetup(t)
	defer d.Close()

	host := fmt.Sprintf("hp-%d.example.test", time.Now().UnixNano())
	taskID := time.Now().UnixNano()
	defer d.Exec("DELETE FROM target_honeypot_assessments WHERE task_id=$1 AND host=$2", taskID, host)

	if err := assets.RecordHoneypotAssessment(taskID, HoneypotAssessment{
		Host:    "https://" + host + ":443/",
		Score:   95,
		Level:   "suspected",
		Signals: []string{"nmap.service.product:cowrie"},
	}); err != nil {
		t.Fatal(err)
	}
	got, err := assets.HoneypotAssessments(taskID, []string{host, "https://" + host + "/"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Host != host || got[0].Score != 95 ||
		got[0].Level != "suspected" || len(got[0].Signals) != 1 ||
		got[0].Signals[0] != "nmap.service.product:cowrie" ||
		!got[0].ExpiresAt.After(time.Now()) {
		t.Fatalf("assessment round trip: %+v", got)
	}

	if _, err := d.Exec("UPDATE target_honeypot_assessments SET expires_at=now()-interval '1 second' WHERE task_id=$1 AND host=$2", taskID, host); err != nil {
		t.Fatal(err)
	}
	got, err = assets.HoneypotAssessments(taskID, []string{host})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("expired assessment must not suppress a target: %+v", got)
	}
}

func TestRecordHoneypotAssessmentRejectsWeakOrIncompleteSignals(t *testing.T) {
	d, assets, _ := testSetup(t)
	defer d.Close()
	for _, assessment := range []HoneypotAssessment{
		{Host: "weak.example.test", Score: 70, Level: "suspected", Signals: []string{"cloud ASN"}},
		{Host: "empty.example.test", Score: 95, Level: "suspected"},
		{Host: "bad-level.example.test", Score: 95, Level: "confirmed", Signals: []string{"nmap.service.product:cowrie"}},
	} {
		if err := assets.RecordHoneypotAssessment(time.Now().UnixNano(), assessment); err == nil {
			t.Fatalf("expected invalid assessment to be rejected: %+v", assessment)
		}
	}
}
