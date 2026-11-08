package report

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
)

func TestSafeCSVCell(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"=1+1", "'=1+1"},
		{"+SUM(1,2)", "'+SUM(1,2)"},
		{"-1+2", "'-1+2"},
		{"@SUM(1,2)", "'@SUM(1,2)"},
		{" \t=HYPERLINK(1)", "' \t=HYPERLINK(1)"},
		{"\r\n@SUM(1)", "'\r\n@SUM(1)"},
		{"\uFEFF=1+1", "'\uFEFF=1+1"},
		{"\u200B+1", "'\u200B+1"},
		{"ordinary finding", "ordinary finding"},
		{"123", "123"},
		{"", ""},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			if got := safeCSVCell(tc.input); got != tc.want {
				t.Errorf("safeCSVCell(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestFindingsCSVUntrustedFieldsAreText(t *testing.T) {
	now := time.Date(2026, 11, 8, 0, 0, 0, 0, time.UTC)
	f := &db.DBFinding{
		ID:              42,
		Name:            "=HYPERLINK(\"https://example.test\")",
		VulnClass:       "+SUM(1,2)",
		Severity:        db.SeverityHigh,
		Status:          db.FindingPending,
		TaskDescription: "\t@SUM(1,2)",
		Summary:         "=1+1",
		CreatedAt:       now,
	}
	raw := FindingsCSV([]*db.DBFinding{f})
	if !bytes.HasPrefix(raw, []byte("\xEF\xBB\xBF")) {
		t.Fatal("CSV lost UTF-8 BOM")
	}
	rows, err := csv.NewReader(bytes.NewReader(raw[3:])).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || len(rows[1]) != len(rows[0]) {
		t.Fatalf("unexpected CSV shape: %q", rows)
	}
	got := rows[1]
	for index, want := range map[int]string{
		0: "42",
		1: "'" + f.Name,
		2: "'" + f.VulnClass,
		3: db.SeverityHigh,
		4: db.FindingPending,
		5: "'" + f.TaskDescription,
		7: "'" + f.Summary,
		8: "0",
	} {
		if got[index] != want {
			t.Errorf("CSV col %d = %q, want %q", index, got[index], want)
		}
	}
	if !strings.HasPrefix(f.Name, "=") || f.Summary != "=1+1" {
		t.Fatal("CSV export mutated source finding")
	}
}
