package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/llmrec"
	actool "github.com/Autumn-27/norma/tool"
)

type journaledTool struct {
	actool.CoreTool
	pg       *db.DB
	fallback agent.RunInfo
}

func journalTaskTools(tools []actool.CoreTool, pg *db.DB, fallback agent.RunInfo) []actool.CoreTool {
	if len(tools) == 0 || pg == nil {
		return tools
	}
	out := make([]actool.CoreTool, len(tools))
	for i, tool := range tools {
		if _, alreadyWrapped := tool.(*journaledTool); alreadyWrapped {
			out[i] = tool
			continue
		}
		out[i] = &journaledTool{CoreTool: tool, pg: pg, fallback: fallback}
	}
	return out
}

func (t *journaledTool) Call(ctx context.Context, input json.RawMessage, tc *actool.ToolContext) (actool.Result, error) {
	ri := agent.RunInfoFrom(ctx)
	if ri.TaskID == 0 {
		ri = t.fallback
	}
	if ri.TaskID <= 0 || t.pg == nil || t.CoreTool.IsReadOnly(input) {
		return t.CoreTool.Call(ctx, input, tc)
	}
	owner, ok := llmrec.TaskOwnerFrom(ctx)
	if !ok || owner.TaskID != ri.TaskID {
		return actool.Result{}, fmt.Errorf("durable tool execution owner is unavailable for task %d", ri.TaskID)
	}
	if tc == nil || tc.ToolUseID == "" {
		return actool.Result{}, fmt.Errorf("durable tool execution requires the provider tool_use_id")
	}
	if ri.SessionID == "" {
		return actool.Result{}, fmt.Errorf("durable tool execution requires a task session id")
	}
	if ri.AgentKey == "" {
		ri.AgentKey = t.fallback.AgentKey
	}
	sum := sha256.Sum256(append(append([]byte(t.Name()), 0), input...))
	req := db.ToolExecutionRequest{
		TaskID: owner.TaskID, RunID: owner.RunID, OwnerEpoch: owner.Epoch,
		IntentID: ri.IntentID, SessionID: ri.SessionID, ToolUseID: tc.ToolUseID,
		ToolKey: t.Name(), AgentKey: ri.AgentKey, InputHash: hex.EncodeToString(sum[:]),
	}
	cached, replay, err := t.pg.BeginToolExecution(req)
	if err != nil {
		return actool.Result{}, fmt.Errorf("durable tool execution blocked: %w", err)
	}
	if replay {
		var result actool.Result
		if err := json.Unmarshal(cached, &result); err != nil {
			return actool.Result{}, fmt.Errorf("decode stored tool result for %s: %w", t.Name(), err)
		}
		return result, nil
	}

	ctx = llmrec.WithToolUseID(ctx, tc.ToolUseID)
	result, callErr := t.CoreTool.Call(ctx, input, tc)
	if callErr != nil {
		_ = t.pg.MarkToolExecutionUnknown(req, "tool returned an error; side-effect outcome is not confirmed")
		return result, callErr
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		_ = t.pg.MarkToolExecutionUnknown(req, "tool result could not be serialized after execution")
		return actool.Result{}, fmt.Errorf("serialize tool result for %s: %w", t.Name(), err)
	}
	if err := t.pg.CompleteToolExecution(req, encoded); err != nil {
		_ = t.pg.MarkToolExecutionUnknown(req, "tool completed but its result could not be durably committed")
		return actool.Result{}, fmt.Errorf("persist completed tool result for %s: %w", t.Name(), err)
	}
	return result, nil
}
