package report

import (
	"encoding/json"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestFindingsSARIFSourceProvenance(t *testing.T) {
	const a = "0123456789abcdef0123456789abcdef"
	const b = "fedcba9876543210fedcba9876543210"
	items := []*db.DBFinding{{ID: 42, Name: "Review note"}}
	for _, id := range []string{"", a, b, a} {
		payload, err := FindingsSARIFWithSource(items, nil, id)
		if err != nil {
			t.Fatal(err)
		}
		var log sarifLog
		if err := json.Unmarshal(payload, &log); err != nil {
			t.Fatal(err)
		}
		if len(log.Runs) != 1 || len(log.Runs[0].Results) != 1 || log.Runs[0].Results[0].Properties.FindingID != 42 {
			t.Fatal("existing finding identity was lost")
		}
		props := log.Runs[0].Properties
		if id == "" && props != nil || id != "" && (props == nil || props.ARTEXSourceID != id) {
			t.Fatalf("unexpected report provenance: %+v", props)
		}
	}
}
