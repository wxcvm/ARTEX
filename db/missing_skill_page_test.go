package db

import (
    "fmt"
    "testing"
    "time"
)

func TestMissingSkillStatsPageCountsBeyondFirstPage(t *testing.T) {
    d, err := Open(testDSN(t))
    if err != nil { t.Skipf("postgres unavailable (%v)", err) }
    defer d.Close()

    prefix := fmt.Sprintf("zz-miss-page-%d-", time.Now().UnixNano())
    installed := prefix + "installed"
    defer func() { _, _ = d.Exec("DELETE FROM skill_usage WHERE skill LIKE $1", prefix+"%") }()
    _, baselineCount, baselineCalls, err := d.MissingSkillStatsPage(20, 1, installed)
    if err != nil { t.Fatal(err) }

    // 26 distinct skill names, two misses each. Installed names must not
    // count as unresolved, including when their historical misses remain.
    for i := 0; i < 26; i++ {
        name := fmt.Sprintf("%s%02d", prefix, i)
        if i == 25 { name = installed }
        for k := 0; k < 2; k++ {
            if err := d.InsertSkillUsage(&SkillUsage{
                Skill: name, AgentKey: "worker", Found: false,
            }); err != nil { t.Fatal(err) }
        }
    }

    first, total, calls, err := d.MissingSkillStatsPage(20, 1, installed)
    if err != nil { t.Fatal(err) }
    if total != baselineCount+25 || calls != baselineCalls+50 || len(first) != 20 {
        t.Fatalf("page totals: items=%d total=%d calls=%d baseline=%d/%d", len(first), total, calls, baselineCount, baselineCalls)
    }

    legacy, err := d.MissingSkillStats(20, installed)
    if err != nil { t.Fatal(err) }
    if len(legacy) != len(first) { t.Fatalf("legacy count=%d page count=%d", len(legacy), len(first)) }
    for i := range first {
        if legacy[i].Skill != first[i].Skill { t.Fatalf("legacy order mismatch at index %d", i) }
    }

    seen := make(map[string]int, 25)
    for page := 1; page <= (total-1)/20+1; page++ {
        items, gotTotal, gotCalls, err := d.MissingSkillStatsPage(20, page, installed)
        if err != nil { t.Fatal(err) }
        if gotTotal != total || gotCalls != calls { t.Fatalf("totals change by page: %d/%d", gotTotal, gotCalls) }
        for _, item := range items {
            if item.Skill == installed { t.Fatal("already installed skill must be excluded") }
            seen[item.Skill] = item.Calls
        }
    }
    for i := 0; i < 25; i++ {
        name := fmt.Sprintf("%s%02d", prefix, i)
        if got := seen[name]; got != 2 { t.Fatalf("missing %s across pages, calls=%d", name, got) }
    }

    empty, beyondTotal, beyondCalls, err := d.MissingSkillStatsPage(20, total+2, installed)
    if err != nil || len(empty) != 0 || beyondTotal != total || beyondCalls != calls {
        t.Fatalf("out-of-range page: items=%d total=%d calls=%d err=%v", len(empty), beyondTotal, beyondCalls, err)
    }
}
