package agent

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/kb"
	actool "github.com/Autumn-27/norma/tool"
)

func TestKnowledgeWorkflowUsesOnlyAvailableTools(t *testing.T) {
	tools := []actool.CoreTool{ragTool("kb_search", "Search"), ragTool("kb_get", "Read"), ragTool("Skill", "Load"), ragTool("ingest_scan_output", "Import"), ragTool("record_coverage", "Record")}
	for _, name := range []string{"kb_search", "kb_get", "Skill", "ingest_scan_output", "record_coverage"} {
		if !strings.Contains(knowledgeWorkflowGuidance(tools), name) {
			t.Fatalf("guidance missing %s", name)
		}
	}
	if got := knowledgeWorkflowGuidance(tools[:1]); got != "" {
		t.Fatalf("incomplete knowledge binding produced guidance: %q", got)
	}
	if got := knowledgeWorkflowGuidance(nil); got != "" {
		t.Fatalf("empty catalog produced guidance: %q", got)
	}
}

func TestKnowledgeToolsAvailableToMainPlannerWorker(t *testing.T) {
	old := Knowledge()
	t.Cleanup(func() { SetKnowledge(old) })
	k, err := kb.Open(filepath.Join("..", "knowledge"))
	if err != nil || k == nil {
		t.Fatalf("open shipped knowledge: %v", err)
	}
	defer k.Close()
	SetKnowledge(k)
	ts := NewToolSet(nil, "")
	for name, tools := range map[string][]actool.CoreTool{"mainagent": ts.MainAgentTools(), "planner": ts.PlannerTools(), "worker": ts.WorkerTools()} {
		var search, get actool.CoreTool
		for _, tool := range tools {
			switch tool.Name() {
			case "kb_search": search = tool
			case "kb_get": get = tool
			}
		}
		if search == nil || get == nil {
			t.Fatalf("%s missing knowledge tools", name)
		}
		res, err := search.Call(t.Context(), json.RawMessage(`{"q":"1panel","limit":1}`), nil)
		if err != nil || res.IsError {
			t.Fatalf("%s search: %v, %s", name, err, res.Flatten())
		}
		var hit struct { Hits []struct { ID int64 `json:"id"` } `json:"hits"` }
		if err := json.Unmarshal([]byte(res.Flatten()), &hit); err != nil || len(hit.Hits) == 0 {
			t.Fatalf("%s search result: %v, %s", name, err, res.Flatten())
		}
		raw, _ := json.Marshal(map[string]int64{"id":hit.Hits[0].ID})
		res, err = get.Call(t.Context(), raw, nil)
		if err != nil || res.IsError {
			t.Fatalf("%s get: %v, %s", name, err, res.Flatten())
		}
	}
	SetKnowledge(nil)
	if strings.Contains(knowledgeWorkflowGuidance(ts.MainAgentTools()), "kb_search") {
		t.Fatal("disabled knowledge advertised")
	}
}
