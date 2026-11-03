package llmpool

import "context"

// queuePersistenceLocked coalesces changes under the same lock that updates m.
// A blocked database cannot block health reads, failover, or manual recovery,
// nor create one goroutine per failed request. The worker exits when idle.
func (r *Registry) queuePersistenceLocked(id int64) {
	if r.persist == nil && r.forget == nil {
		return
	}
	if r.pending == nil {
		r.pending = make(map[int64]struct{})
	}
	r.pending[id] = struct{}{}
	if r.flushDone == nil {
		r.flushDone = make(chan struct{})
		go r.persistLoop()
	}
}

func (r *Registry) persistLoop() {
	for {
		r.mu.Lock()
		if len(r.pending) == 0 {
			close(r.flushDone)
			r.flushDone = nil
			r.mu.Unlock()
			return
		}
		// Drain a batch so a continuously changing profile cannot starve others.
		ids := make([]int64, 0, len(r.pending))
		for id := range r.pending {
			ids = append(ids, id)
		}
		r.mu.Unlock()
		for _, id := range ids {
			r.mu.Lock()
			delete(r.pending, id)
			st, exists := r.m[id]
			var snap State
			if exists {
				snap = *st
			}
			r.mu.Unlock()
			// Changes while I/O is in flight re-queue this id. In particular,
			// Pass/Reset will delete AFTER an old Trip finishes saving, never
			// before it. Callbacks must not launch their own unordered writes.
			if exists && r.persist != nil {
				r.persist(id, snap)
			} else if !exists && r.forget != nil {
				r.forget(id)
			}
		}
	}
}

// FlushPersistence waits until queued callbacks have returned, or ctx expires.
// Persistence remains best-effort: callbacks report their own I/O errors. This
// is not a transaction or a guarantee against abrupt process termination.
// Do not call it from a persistence callback (that would wait on itself).
func (r *Registry) FlushPersistence(ctx context.Context) error {
	for {
		r.mu.Lock()
		done := r.flushDone
		r.mu.Unlock()
		if done == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
		}
	}
}
