package llmrec

import "context"

// TaskOwner is the immutable fencing identity carried by task execution contexts.
// Tool wrappers use it to bind side effects to the same durable run lease as LLM calls.
type TaskOwner struct {
	TaskID int64
	RunID  string
	Epoch  int64
}

type taskOwnerContextKey struct{}

// WithTaskOwner binds the active durable task owner to ctx. A zero-value owner
// removes no existing value; callers should only use it after acquiring a lease.
func WithTaskOwner(ctx context.Context, owner TaskOwner) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, taskOwnerContextKey{}, owner)
}

// TaskOwnerFrom returns the task's durable run and fencing epoch, when present.
func TaskOwnerFrom(ctx context.Context) (TaskOwner, bool) {
	if ctx == nil {
		return TaskOwner{}, false
	}
	owner, ok := ctx.Value(taskOwnerContextKey{}).(TaskOwner)
	if !ok || owner.TaskID <= 0 || owner.RunID == "" || owner.Epoch <= 0 {
		return TaskOwner{}, false
	}
	return owner, true
}

type toolUseIDContextKey struct{}

// WithToolUseID makes the provider's stable tool-use identifier available to
// tool implementations that support native idempotency keys.
func WithToolUseID(ctx context.Context, toolUseID string) context.Context {
	if ctx == nil || toolUseID == "" {
		return ctx
	}
	return context.WithValue(ctx, toolUseIDContextKey{}, toolUseID)
}

// ToolUseIDFrom returns the stable tool-use identifier, when one is present.
func ToolUseIDFrom(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	value, _ := ctx.Value(toolUseIDContextKey{}).(string)
	return value
}
