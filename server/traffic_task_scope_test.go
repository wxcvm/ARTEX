package server

import (
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
)

func TestTaskTrafficBounds(t *testing.T) {
	created := time.Unix(100, 0)
	first := time.Unix(200, 0)
	finished := time.Unix(400, 0)
	now := time.Unix(500, 0)
	for _, tc := range []struct {
		name string
		task db.Task
		min, max int64
	}{
		{"never started", db.Task{CreatedAt: created}, 100, 500},
		{"running", db.Task{CreatedAt: created, FirstRunAt: &first}, 200, 500},
		{"finished", db.Task{CreatedAt: created, FirstRunAt: &first, CompletedAt: &finished}, 200, 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			min, max := taskTrafficBounds(&tc.task, now)
			if min != tc.min || max != tc.max {
				t.Fatalf("got (%d,%d), want (%d,%d)", min, max, tc.min, tc.max)
			}
		})
	}
}
