package server

import (
	"errors"
	"fmt"
)

var errTaskSettling = errors.New("task is settling after its deadline; pause is no longer available")

// lockTaskTransition looks up the live handle without holding Manager.mu while
// waiting for its transition lock. Ordering is transitionMu -> taskStateMu ->
// Manager.mu -> lifecycleMu. Missing/deleted rows remain the DB's responsibility.
func (m *Manager) lockTaskTransition(id string) func() {
	m.mu.RLock()
	t := m.tasks[id]
	m.mu.RUnlock()
	if t == nil {
		return func() {}
	}
	t.transitionMu.Lock()
	return t.transitionMu.Unlock
}

// pauseTaskWithCommit shares its admission lock with nextDeadlineAction. The
// execution barrier is installed before persistence, but cancellation happens
// only after a successful commit: a DB failure must not kill a healthy run.
// commit must synchronously persist AND mirror Paused on the live Task.
func (e *Engine) pauseTaskWithCommit(t *Task, cause error, commit func() error) error {
	t.transitionMu.Lock()
	defer t.transitionMu.Unlock()
	e.deleteMu.RLock()
	if e.IsDeleting(t.ID) {
		e.deleteMu.RUnlock()
		return fmt.Errorf("task is being deleted, control unavailable")
	}
	if e.isSettling(t.ID) {
		e.deleteMu.RUnlock()
		return errTaskSettling
	}
	previous := e.IsPaused(t.ID)
	e.paused.Store(t.ID, true)
	e.deleteMu.RUnlock()
	if err := commit(); err != nil {
		// A concurrent delete owns its own barrier; compensation must not clear it.
		e.deleteMu.RLock()
		if !previous && !e.IsDeleting(t.ID) {
			e.paused.Delete(t.ID)
			t.Notify()
		}
		e.deleteMu.RUnlock()
		return err
	}
	e.cancelExec(t.ID, cause)
	return nil
}
