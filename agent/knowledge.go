package agent

import (
	"context"
	"encoding/json"

	"github.com/Autumn-27/artex/kb"
	actool "github.com/Autumn-27/norma/tool"
)

// knowledgeBase is the process-wide read-only knowledge base. It is optional:
// when nil the kb_* tools are not registered at all, so a deployment without a
// knowledge/ directory behaves exactly as before.
var knowledgeBase *kb.KB

// SetKnowledge installs the knowledge base behind kb_search / kb_get. Call it
// once at startup, before any task is created.
func SetKnowledge(k *kb.KB) { knowledgeBase = k }

// Knowledge returns the installed handle (nil when the feature is off).
func Knowledge() *kb.KB { return knowledgeBase }

// kbTools returns the two knowledge tools, or nil when no base is installed.
func kbTools() []actool.CoreTool {
	if knowledgeBase == nil {
		return nil
	}
	return []actool.CoreTool{kbSearchTool(), kbGetTool()}
}

// withKB appends the knowledge tools to a role's base list. A composite literal
// cannot spread a slice, so roles that want them go through this wrapper.
func withKB(tools []actool.CoreTool) []actool.CoreTool {
	return append(tools, kbTools()...)
}

func kbSearchTool() actool.CoreTool {
	return readTool("kb_search",
		"Search the deployment knowledge base before hand-crafting a payload: past engagements with POCs, playbooks and negative results. q takes a product name / component / CVE / technique keyword (Chinese or English). Returns compact hits (id, title, tier, snippet). tier=ammo is verified and directly reusable; tier=negative means that direction was already proven dead, so do not retry it. Read the full entry with kb_get(id). Read-only, no side effects.",
		obj(map[string]any{
			"q":     str("keywords: product / component / CVE / technique, e.g. '1panel 验证码绕过' or 'phpwind 写文件'"),
			"limit": intp("max hits to return, default 8, max 30"),
		}, "q"),
		func(_ context.Context, in json.RawMessage) (actool.Result, error) {
			var a struct {
				Q     string `json:"q"`
				Limit int    `json:"limit"`
			}
			_ = json.Unmarshal(in, &a)
			hits, degraded, err := knowledgeBase.Search(a.Q, a.Limit)
			if err != nil {
				return actool.Errorf("kb_search: " + err.Error()), nil
			}
			if len(hits) == 0 {
				return actool.Text("no matching knowledge entries"), nil
			}
			return jsonResult(map[string]any{"hits": hits, "count": len(hits), "degraded": degraded})
		})
}

func kbGetTool() actool.CoreTool {
	return readTool("kb_get",
		"Read one knowledge entry in full by the id returned from kb_search. The body is the recorded PoC / playbook (targets are placeholders such as <target>). Read-only.",
		obj(map[string]any{
			"id": idp("entry id taken from kb_search"),
		}, "id"),
		func(_ context.Context, in json.RawMessage) (actool.Result, error) {
			var a struct {
				ID int64 `json:"id"`
			}
			_ = json.Unmarshal(in, &a)
			if a.ID <= 0 {
				return actool.Errorf("id is required (take it from kb_search)"), nil
			}
			e, err := knowledgeBase.Get(a.ID)
			if err != nil {
				return actool.Errorf("kb_get: " + err.Error()), nil
			}
			if e == nil {
				return actool.Errorf("no knowledge entry with that id"), nil
			}
			return jsonResult(e)
		})
}
