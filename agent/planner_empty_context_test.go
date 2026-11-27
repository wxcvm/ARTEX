package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestPlannerRejectsNilExplorationStore(t *testing.T) {
	p := &Planner{}
	var activities []db.Activity
	_, _, err := p.Plan(context.Background(), 42, nil, nil, "", nil, func(a db.Activity) {
		activities = append(activities, a)
	})
	if !errors.Is(err, ErrPlannerEmptyContext) {
		t.Fatalf("expected ErrPlannerEmptyContext, got %v", err)
	}
	if len(activities) != 1 || !activities[0].IsError {
		t.Fatalf("expected one error activity, got %+v", activities)
	}
}

func TestPlannerRejectsEmptyGoalAndEmptyGoalNodes(t *testing.T) {
	d := testDB(t)
	defer d.Close()
	task, err := d.CreateTask("empty goal context", "", nil, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer d.DeleteTask(task.ID)
	ts := d.Exploration(task.ExplorationID)

	p := &Planner{}
	var activities []db.Activity
	_, _, err = p.Plan(context.Background(), task.ID, nil, ts, "", nil, func(a db.Activity) {
		activities = append(activities, a)
	})
	if !errors.Is(err, ErrPlannerEmptyContext) {
		t.Fatalf("expected ErrPlannerEmptyContext, got %v", err)
	}
	if len(activities) != 1 || !activities[0].IsError {
		t.Fatalf("expected one error activity, got %+v", activities)
	}
}

func TestFirstNonEmptyGoalNodeSkipsBlankGoals(t *testing.T) {
	d := testDB(t)
	defer d.Close()
	task, err := d.CreateTask("goal node context", "", nil, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer d.DeleteTask(task.ID)
	ts := d.Exploration(task.ExplorationID)

	if _, err := ts.AddNode(db.KindGoal, map[string]any{"text": "   "}, 0, "open", "test", nil); err != nil {
		t.Fatal(err)
	}
	if got, err := firstNonEmptyGoalNode(ts); err != nil || got != "" {
		t.Fatalf("blank goal must be ignored: got=%q err=%v", got, err)
	}
	if _, err := ts.AddNode(db.KindGoal, map[string]any{"text": "  Verify the supplied task objective  "}, 0, "open", "test", nil); err != nil {
		t.Fatal(err)
	}
	if got, err := firstNonEmptyGoalNode(ts); err != nil || got != "Verify the supplied task objective" {
		t.Fatalf("expected trimmed non-empty goal: got=%q err=%v", got, err)
	}
	if got, err := hasNonEmptyGoalNodes(ts); err != nil || !got {
		t.Fatalf("non-empty goal must count: got=%v err=%v", got, err)
	}
}
