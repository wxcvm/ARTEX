package agent

import "testing"

func TestGraphOverviewKeepsOlderRunningIntentAndActiveHints(t *testing.T) {
    d := testDB(t)
    defer d.Close()

    task, err := d.CreateTask("overview state fixture", "verify active work", nil, 0, 0)
    if err != nil { t.Fatal(err) }
    defer func() { _ = d.DeleteTask(task.ID) }()
    store := d.Exploration(task.ExplorationID)

    runningID, err := store.AddIntent(map[string]any{"summary": "older work still running"}, 1, nil, "planner")
    if err != nil { t.Fatal(err) }
    if err := store.SetIntentState(runningID, "running"); err != nil { t.Fatal(err) }

    // The running intent is older than the legacy 300-intent overview window.
    if _, err := d.Exec(`INSERT INTO exploration_nodes(exploration_id, kind, payload, state)
        SELECT $1, 'intent', '{}', 'done' FROM generate_series(1, 305)`, task.ExplorationID); err != nil {
        t.Fatal(err)
    }
    if _, err := d.Exec(`INSERT INTO exploration_nodes(exploration_id,kind,payload,state) VALUES
        ($1,'hint','{"text":"active human instruction"}','active'),
        ($1,'hint','{"text":"consumed human instruction"}','consumed')`, task.ExplorationID); err != nil {
        t.Fatal(err)
    }

    overview := NewToolSet(store, "planner").graphOverviewData()
    if got := overview["running_intents_total"]; got != 1 {
        t.Fatalf("running_intents_total=%v want 1", got)
    }
    running, ok := overview["running_intents"].([]map[string]any)
    if !ok || len(running) != 1 || running[0]["id"] != runningID {
        t.Fatalf("missing older running intent: %#v", overview["running_intents"])
    }
    hints, ok := overview["hints"].([]map[string]any)
    if !ok || len(hints) != 1 || hints[0]["text"] != "active human instruction" {
        t.Fatalf("consumed hint leaked or active hint missing: %#v", overview["hints"])
    }
    if overview["done_intents_total"] != 305 {
        t.Fatalf("unexpected terminal intent count: %v", overview["done_intents_total"])
    }
}

