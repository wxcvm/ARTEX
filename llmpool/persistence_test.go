package llmpool

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func flushRegistry(t *testing.T, r *Registry) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.FlushPersistence(ctx); err != nil {
		t.Fatalf("flush circuit-breaker persistence: %v", err)
	}
}

func waitPersistenceSignal(t *testing.T, signal <-chan struct{}) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(5 * time.Second):
		t.Fatal("persistence callback did not start")
	}
}

// Simulate a slow old write. Recovery must remain non-blocking, and the final
// durable state must be a deletion, not a resurrected open breaker.
func TestPersistenceRecoveryCannotBeOvertakenByOldTrip(t *testing.T) {
	for _, action := range []string{"pass", "reset"} {
		t.Run(action, func(t *testing.T) {
			started, release := make(chan struct{}), make(chan struct{})
			var once sync.Once
			t.Cleanup(func() { once.Do(func() { close(release) }) })
			var disk *State
			var operations []string
			r := NewRegistry(func(_ int64, st State) {
				close(started)
				<-release
				disk = &st
				operations = append(operations, "save")
			}, func(int64) {
				disk = nil
				operations = append(operations, "delete")
			})
			r.Trip(1, "status 402", true)
			waitPersistenceSignal(t, started)
			if action == "pass" {
				r.Pass(1)
			} else {
				r.Reset(1)
			}
			if r.IsOpen(1) {
				t.Fatal("in-memory recovery waited for database I/O")
			}
			once.Do(func() { close(release) })
			flushRegistry(t, r)
			if disk != nil || len(operations) != 2 || operations[0] != "save" || operations[1] != "delete" {
				t.Fatalf("old state resurrected: disk=%+v operations=%v", disk, operations)
			}
		})
	}
}

// A burst during slow I/O needs one latest-state write, not an unbounded queue
// of obsolete snapshots. A deletion followed by another failure must persist
// the NEW failure rather than clear it when the old write finally returns.
func TestPersistenceCoalescesLatestState(t *testing.T) {
	started, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	t.Cleanup(func() { once.Do(func() { close(release) }) })
	disk := make(map[int64]State)
	writes := 0
	var r *Registry
	r = NewRegistry(func(id int64, st State) {
		// Registry locks must not be held during a callback.
		_ = r.Snapshot()
		writes++
		if writes == 1 {
			close(started)
			<-release
		}
		disk[id] = st
	}, func(id int64) { delete(disk, id) })
	r.Trip(1, "old failure", true)
	waitPersistenceSignal(t, started)
	r.Reset(1)
	for range 1000 {
		r.Trip(1, "new failure", false)
	}
	r.Trip(2, "second profile", true)
	latest := r.Snapshot()
	r.mu.Lock()
	pending := len(r.pending)
	r.mu.Unlock()
	if pending != 2 {
		t.Fatalf("pending=%d, want one entry per changed profile", pending)
	}
	once.Do(func() { close(release) })
	flushRegistry(t, r)
	if writes != 3 || disk[1] != latest[1] || disk[2] != latest[2] {
		t.Fatalf("latest state not mirrored: writes=%d disk=%+v latest=%+v", writes, disk, latest)
	}
	// A later burst starts a fresh worker after the first exits.
	r.Reset(1)
	r.Reset(2)
	flushRegistry(t, r)
	if len(disk) != 0 {
		t.Fatalf("subsequent recovery was lost: %+v", disk)
	}
}

func TestPersistenceFlushCanBeCancelled(t *testing.T) {
	started, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	t.Cleanup(func() { once.Do(func() { close(release) }) })
	r := NewRegistry(func(int64, State) {
		close(started)
		<-release
	}, nil)
	r.Trip(1, "failed", true)
	waitPersistenceSignal(t, started)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := r.FlushPersistence(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("flush error=%v, want context.Canceled", err)
	}
	once.Do(func() { close(release) })
	flushRegistry(t, r)
}

func TestPersistenceWithoutCallbacksDoesNotStartWorker(t *testing.T) {
	r := NewRegistry(nil, nil)
	r.Trip(1, "failed", true)
	r.Pass(1)
	r.Reset(2)
	if r.pending != nil || r.flushDone != nil {
		t.Fatal("registry without storage started a persistence worker")
	}
	flushRegistry(t, r)
}
