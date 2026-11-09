package agent

import (
    "encoding/json"
    "strings"
    "testing"
)

func TestPlannerOverviewCompactKeepsCriticalContext(t *testing.T) {
    data := map[string]any{
        "task": map[string]any{"description": "fixture", "goal": "review evidence"},
        "running_intents_total": 1,
        "running_intents": []map[string]any{{"id": int64(7), "state": "running", "summary": "older task still running"}},
        "hints": []map[string]any{{"id": int64(8), "state": "active", "text": "human instruction", "traffic_refs": []string{"ref-5"}}},
        "stagnation_warning": map[string]any{"instruction": "change approach"},
        "task_credentials": []map[string]any{{"username": "test-account", "alias": "fixture"}},
        "task_credentials_note": "reuse existing test account",
        "cold_digests": []map[string]any{{
            "id": int64(42), "member_count": 4,
            "body": strings.Repeat("historic result ", 180),
            "unresolved": []map[string]any{{"intent_id": int64(9), "state": "blocked", "reason": "dependency unfinished"}},
        }},
        "cold_digests_more": []int64{43, 44},
    }
    before, err := json.Marshal(data)
    if err != nil { t.Fatal(err) }
    view := plannerOverviewText(data)
    for _, needle := range []string{"older task still running", "human instruction", "ref-5",
        "change approach", "test-account", "reuse existing test account", "historic result",
        "dependency unfinished", "42,43,44", "expand_digest(42)", "graph_overview"} {
        if !strings.Contains(view, needle) { t.Errorf("compact overview lost %q", needle) }
    }
    after, err := json.Marshal(data)
    if err != nil { t.Fatal(err) }
    if string(before) != string(after) { t.Fatal("read-only compact converter mutated overview data") }
}

func TestPlannerOverviewUsesCompactOnlyWhenSmaller(t *testing.T) {
    tiny := map[string]any{"facts": 1}
    got := renderGraphOverview(tiny)
    if !strings.Contains(got, "\"" + "facts" + "\":1") || strings.Contains(got, "read-only compact") {
        t.Fatalf("tiny snapshot should use original JSON: %s", got)
    }

    rows := make([]map[string]any, 0, 80)
    for i := 0; i < 80; i++ {
        rows = append(rows, map[string]any{"id": i+1, "summary": strings.Repeat("repeated observation ", 30)})
    }
    large := map[string]any{
        "task": map[string]any{"description": "test", "goal": "review"},
        "recent_facts": rows,
        "frontier_open": 0,
        "running_intents_total": 0,
    }
    raw, err := json.Marshal(large)
    if err != nil { t.Fatal(err) }
    got = renderGraphOverview(large)
    if !strings.Contains(got, "read-only compact graph_overview") || len(got) >= len(raw) {
        t.Fatalf("large overview not compact: compact=%d raw=%d", len(got), len(raw))
    }
    if !strings.Contains(got, "按原始工具查阅") || !strings.Contains(got, "graph_overview") {
        t.Fatal("older rows must have an explicit recovery path")
    }
}
