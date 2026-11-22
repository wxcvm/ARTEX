package server

import "github.com/Autumn-27/artex/agent"

func (e *Engine) honeypotFor(taskID string) *agent.HoneypotRegistry {
	if v, ok := e.honeypots.Load(taskID); ok {
		return v.(*agent.HoneypotRegistry)
	}
	v, _ := e.honeypots.LoadOrStore(taskID, agent.NewHoneypotRegistry())
	return v.(*agent.HoneypotRegistry)
}
