package server

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

// A real PostgreSQL constraint failure must restore the previously running
// engine barrier without cancelling its execution context or changing the row.
func TestPauseCommitFailureCompensatesWithPostgres(t *testing.T) {
	s, _, _ := trafficEvidenceServer(t)
	task, err := s.m.CreateTask("pause rollback", "local fixture", nil, 3600, 0)
	if err != nil {
		t.Fatal(err)
	}
	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.m.pg.Exec(`DELETE FROM tasks WHERE id=$1`, id) })
	constraint := fmt.Sprintf("test_pause_failure_%d", id)
	_, err = s.m.pg.Exec(fmt.Sprintf(`ALTER TABLE tasks ADD CONSTRAINT %s CHECK (id <> %d OR NOT paused)`, constraint, id))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := s.m.pg.Exec(`ALTER TABLE tasks DROP CONSTRAINT ` + constraint); err != nil {
			t.Error(err)
		}
	})
	ctx := s.engine.execContextFor(context.Background(), task.ID)
	_, err = s.applyTaskControl(task, "pause")
	if err == nil {
		t.Fatal("expected PostgreSQL constraint failure")
	}
	var paused bool
	if err := s.m.pg.QueryRow(`SELECT paused FROM tasks WHERE id=$1`, id).Scan(&paused); err != nil {
		t.Fatal(err)
	}
	if paused || task.lifecycleSnapshot().Paused || s.engine.IsPaused(task.ID) || ctx.Err() != nil {
		t.Fatalf("pause commit failed but row/barrier/context diverged: row=%v ctx=%v", paused, ctx.Err())
	}
}

// Pausing a task must freeze its wall-clock timeout budget: the admission that
// ends the pause credits the paused interval back onto deadline_at, so a task
// paused overnight (or while the machine is off) resumes with its remaining
// time instead of immediately settling into the timeout terminal state (#1).
func TestPauseShiftsTaskDeadlineOnResume(t *testing.T) {
	s, _, _ := trafficEvidenceServer(t)
	m := s.m
	pg := s.m.pg

	task, err := m.CreateTask("pause deadline test", "verify the deadline is frozen while paused", nil, 3600, 0)
	if err != nil {
		t.Fatal(err)
	}
	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { pg.Exec(`DELETE FROM tasks WHERE id=$1`, id) })

	// A running task with a stamped deadline (first_run_at + 1h budget).
	if _, err := m.StampTaskFirstRun(task.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := pg.Exec(`UPDATE tasks SET status='running' WHERE id=$1`, id); err != nil {
		t.Fatal(err)
	}
	var before time.Time
	if err := pg.QueryRow(`SELECT deadline_at FROM tasks WHERE id=$1`, id).Scan(&before); err != nil {
		t.Fatal(err)
	}

	// Pause stamps paused_at alongside the user pause.
	if err := m.ApplyTaskPause(task.ID); err != nil {
		t.Fatal(err)
	}
	var pausedAt sql.NullTime
	if err := pg.QueryRow(`SELECT paused_at FROM tasks WHERE id=$1`, id).Scan(&pausedAt); err != nil {
		t.Fatal(err)
	}
	if !pausedAt.Valid || pausedAt.Time.After(time.Now()) {
		t.Fatalf("pause should stamp paused_at at the current time, got %v", pausedAt.Time)
	}

	// Set a known persisted pause interval instead of depending on wall sleeps.
	const paused = 10 * time.Minute
	if _, err := pg.Exec(`UPDATE tasks SET paused_at=now()-interval '10 minutes' WHERE id=$1`, id); err != nil {
		t.Fatal(err)
	}
	if got := s.engine.taskDeadline(task); got != before.Unix() {
		t.Fatalf("initial engine deadline=%d want=%d", got, before.Unix())
	}

	// Admission out of the paused state (a paused task keeps its lifecycle status).
	if err := m.ApplyTaskAdmission(task.ID, "running", "running", false, "resume", false); err != nil {
		t.Fatal(err)
	}

	var after time.Time
	var cleared sql.NullTime
	if err := pg.QueryRow(`SELECT deadline_at, paused_at FROM tasks WHERE id=$1`, id).Scan(&after, &cleared); err != nil {
		t.Fatal(err)
	}
	if shifted := after.Sub(before); shifted < paused-100*time.Millisecond {
		t.Fatalf("resume should credit the paused interval onto the deadline, got shift=%v (paused=%v)", shifted, paused)
	}
	if cleared.Valid {
		t.Fatalf("admission should clear paused_at, got %v", cleared.Time)
	}
	// The live handle must mirror the shifted deadline (the engine clock context
	// is built from it).
	if got := task.lifecycleSnapshot().DeadlineAt; got != after.Unix() {
		t.Fatalf("live handle deadline = %d, want the shifted %d", got, after.Unix())
	}
	if got := s.engine.taskDeadline(task); got != after.Unix() {
		t.Fatalf("engine retained stale deadline=%d want=%d", got, after.Unix())
	}
}
