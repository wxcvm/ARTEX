package server

import (
	"time"

	"github.com/Autumn-27/artex/db"
)

// taskTrafficBounds uses a task's actual first run to avoid borrowing captures
// from earlier runs against the same host. Completed tasks have a fixed upper
// bound; live/paused tasks extend to the current time.
func taskTrafficBounds(task *db.Task, now time.Time) (int64, int64) {
	start := task.CreatedAt
	if task.FirstRunAt != nil {
		start = *task.FirstRunAt
	}
	end := now
	if task.CompletedAt != nil {
		end = *task.CompletedAt
	}
	return start.Unix(), end.Unix()
}
