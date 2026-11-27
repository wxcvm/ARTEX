package db

import "testing"

func TestToolUsageLedger(t *testing.T) {
	d, err := Open(testDSN(t))
	if err != nil {
		t.Skipf("postgres unavailable (%v) - skipping", err)
	}
	defer d.Close()

	const (
		toolA = "zz_test_tool_usage_a"
		toolB = "zz_test_tool_usage_b"
	)
	cleanup := func() {
		_, _ = d.Exec(`DELETE FROM tool_usage WHERE tool_key IN ($1,$2)`, toolA, toolB)
	}
	cleanup()
	defer cleanup()

	rows := []*ToolUsage{
		{ToolKey: toolA, AgentKey: "worker", TaskID: 991, ExplorationID: 5, IntentID: 7},
		{ToolKey: toolA, AgentKey: "planner", TaskID: 992, ExplorationID: 6},
		{ToolKey: toolA, AgentKey: "chatbot", SessionID: "conv-1"},
		{ToolKey: toolB, AgentKey: "worker", TaskID: 991},
	}
	for _, row := range rows {
		if err := d.InsertToolUsage(row); err != nil {
			t.Fatalf("insert %s: %v", row.ToolKey, err)
		}
	}

	counts, err := d.ToolUsageCounts()
	if err != nil {
		t.Fatal(err)
	}
	if got := counts[toolA]; got != 3 {
		t.Errorf("%s calls: want 3, got %d", toolA, got)
	}
	if got := counts[toolB]; got != 1 {
		t.Errorf("%s calls: want 1, got %d", toolB, got)
	}
}

func TestToolUsageDeduplicatesStableToolUseID(t *testing.T) {
	d, taskID := budgetJournalTestTask(t)
	row := &ToolUsage{
		ToolKey: "insert_assets", AgentKey: "worker", TaskID: taskID,
		ExplorationID: 1, IntentID: 7, SessionID: "session-tool-usage",
		ToolUseID: "tool-use-stable",
	}
	if err := d.InsertToolUsage(row); err != nil {
		t.Fatal(err)
	}
	if err := d.InsertToolUsage(row); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := d.QueryRow(`SELECT COUNT(*) FROM tool_usage
WHERE task_id=$1 AND session_id=$2 AND tool_use_id=$3`,
		taskID, row.SessionID, row.ToolUseID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("duplicate replay created %d usage rows; want 1", count)
	}
}

