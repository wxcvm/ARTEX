package db

import "testing"

func TestFilteredFindingStatsMatchesList(t *testing.T) {
	d, err := Open(testDSN(t))
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	const marker = "__filtered_stats_regression__"
	cleanup := func() { _, _ = d.Exec(`DELETE FROM findings WHERE vulnclass=$1`, marker) }
	cleanup()
	defer cleanup()
	for _, tc := range []struct{ severity, status, name string }{
		{"critical", FindingPending, "needle first"},
		{"high", FindingResolved, "other"},
		{"low", FindingPending, "needle second"},
	} {
		id, err := d.AddFinding(0, 0, marker, tc.name, tc.severity, "summary", "evidence", "test", nil)
		if err != nil { t.Fatal(err) }
		if _, err := d.SetFindingStatus(id, tc.status); err != nil { t.Fatal(err) }
	}
	for _, tc := range []struct {
		name string
		filter FindingFilter
		total, pending, critical, high, low int
	}{
		{"class", FindingFilter{VulnClass: marker}, 3, 2, 1, 1, 1},
		{"status and query", FindingFilter{VulnClass: marker, Status: FindingPending, Query: "needle"}, 2, 2, 1, 0, 1},
		{"severity", FindingFilter{VulnClass: marker, Severity: "high"}, 1, 0, 0, 1, 0},
		{"empty", FindingFilter{VulnClass: marker, Query: "no-such-result"}, 0, 0, 0, 0, 0},
		{"unlinked", FindingFilter{VulnClass: marker, AssetScope: FindingUnassignedAsset}, 3, 2, 1, 1, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			st, err := d.FilteredFindingStats(tc.filter)
			if err != nil { t.Fatal(err) }
			_, total, err := d.ListFindingsPage(tc.filter, 1, 1)
			if err != nil { t.Fatal(err) }
			if st.Total != tc.total || total != tc.total || st.Pending != tc.pending || st.Critical != tc.critical || st.High != tc.high || st.Low != tc.low {
				t.Fatalf("stats=%+v list=%d want=%d", st, total, tc.total)
			}
			found := false
			for _, option := range st.VulnClasses { if option == marker { found = true } }
			if !found { t.Fatal("global filter options lost") }
		})
	}
}
