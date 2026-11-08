package db

import (
	"context"
	"testing"
)

func TestInstalledSkillDoesNotRemainInGapBacklog(t *testing.T) {
	d, err := Open(testDSN(t))
	if err != nil {
		t.Skipf("postgres unavailable: %v", err)
	}
	defer d.Close()

	const unresolved = "zz-gap-unresolved-case"
	const installed = "zz-gap-installed-case"
	clean := func() {
		_, _ = d.Exec("DELETE FROM skill_usage WHERE skill IN ($1,$2)", unresolved, installed)
	}
	clean()
	defer clean()

	for _, row := range []*SkillUsage{
		{Skill: unresolved, AgentKey: "worker", Found: false},
		{Skill: installed, AgentKey: "worker", Found: false},
		{Skill: installed, AgentKey: "worker", Found: true},
	} {
		if err := d.InsertSkillUsage(row); err != nil {
			t.Fatal(err)
		}
	}

	hasMiss := func(skills []SkillStat, name string) bool {
		for _, s := range skills {
			if s.Skill == name {
				return true
			}
		}
		return false
	}

	before, err := d.MissingSkillStats(200)
	if err != nil {
		t.Fatal(err)
	}
	if !hasMiss(before, installed) || !hasMiss(before, unresolved) {
		t.Fatal("missing history not recorded before filtering")
	}

	after, err := d.MissingSkillStats(200, installed)
	if err != nil {
		t.Fatal(err)
	}
	if hasMiss(after, installed) || !hasMiss(after, unresolved) {
		t.Fatalf("current backlog filtered incorrectly: %+v", after)
	}

	allGaps, err := d.SkillGaps(context.Background(), 500)
	if err != nil {
		t.Fatal(err)
	}
	resolvedGaps, err := d.SkillGaps(context.Background(), 500, installed)
	if err != nil {
		t.Fatal(err)
	}
	hasGap := func(rows []SkillGap, name string) bool {
		for _, row := range rows {
			if row.Skill == name {
				return true
			}
		}
		return false
	}
	if !hasGap(allGaps, installed) || !hasGap(allGaps, unresolved) {
		t.Fatal("original skill gap history missing")
	}
	if hasGap(resolvedGaps, installed) || !hasGap(resolvedGaps, unresolved) {
		t.Fatal("installed name still appears in scoped gap report")
	}

	// No ledger mutation: removing a Skill from the library can expose its
	// historical gaps again.
	again, err := d.MissingSkillStats(200)
	if err != nil {
		t.Fatal(err)
	}
	if !hasMiss(again, installed) {
		t.Fatal("historical gap was destroyed by filtering")
	}
}
