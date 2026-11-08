package server

import (
	"context"
	"fmt"
	"strconv"
	"testing"
)

// A generic lesson has NULL scope columns. In PostgreSQL, ordering a
// (nullable_column = task_scope) comparison DESC puts NULL before TRUE by
// default. This test catches that regression with a real database and LIMIT.
func TestTaskLessonsPreferMatchingScopes(t *testing.T) {
	m, err := NewManager(t.TempDir(), "")
	if err != nil {
		t.Skipf("postgres unavailable (%v) - skipping", err)
	}
	defer m.Close()

	task, err := m.CreateTask("scoped lesson ranking", "ranking check", nil, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	taskID, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = m.pg.DeleteTask(taskID) }()

	prefix := "lesson-rank-" + task.ID
	var categoryID, companyID, otherCategoryID, otherCompanyID int64
	for _, entry := range []struct {
		table string
		label string
		id    *int64
	}{
		{"task_categories", "category", &categoryID},
		{"companies", "company", &companyID},
		{"task_categories", "other-category", &otherCategoryID},
		{"companies", "other-company", &otherCompanyID},
	} {
		key := prefix + "-" + entry.label
		if err := m.pg.QueryRow(fmt.Sprintf(
			"INSERT INTO %s(name,nkey) VALUES($1,$2) RETURNING id", entry.table,
		), key, key).Scan(entry.id); err != nil {
			t.Fatal(err)
		}
	}
	defer func() {
		// The task deletion is a separate defer; foreign keys allow the
		// cleanup here once the test has finished.
		_, _ = m.pg.Exec("DELETE FROM rsi_lessons WHERE dedupe_key LIKE $1", prefix+"-%")
		_, _ = m.pg.Exec("UPDATE tasks SET category_id=NULL, company_id=NULL WHERE id=$1", taskID)
		_, _ = m.pg.Exec("DELETE FROM task_categories WHERE id IN ($1,$2)", categoryID, otherCategoryID)
		_, _ = m.pg.Exec("DELETE FROM companies WHERE id IN ($1,$2)", companyID, otherCompanyID)
	}()
	if _, err := m.pg.Exec("UPDATE tasks SET category_id=$1,company_id=$2 WHERE id=$3",
		categoryID, companyID, taskID); err != nil {
		t.Fatal(err)
	}

	add := func(name, kind string, category, company int64) int64 {
		t.Helper()
		var id int64
		err := m.pg.QueryRow(`
INSERT INTO rsi_lessons(dedupe_key,kind,cue,action,category_id,company_id)
VALUES ($1,$2,$3,$4,NULLIF($5,0),NULLIF($6,0)) RETURNING id`,
			prefix+"-"+name, kind, name, "sample "+name, category, company,
		).Scan(&id)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := m.pg.Exec(`
INSERT INTO rsi_lesson_evidence(lesson_id,source_kind,source_id)
VALUES ($1,'fixture',$2)`, id, id); err != nil {
			t.Fatal(err)
		}
		return id
	}

	both := add("both", "caution", categoryID, companyID)
	company := add("company-only", "negative", 0, companyID)
	category := add("category-only", "positive", categoryID, 0)
	add("generic", "positive", 0, 0)
	add("wrong-company", "positive", 0, otherCompanyID)
	add("wrong-category", "positive", otherCategoryID, 0)

	scoped, err := m.pg.TaskLessons(context.Background(), taskID, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(scoped) != 3 {
		t.Fatalf("scope ranking returned %d lessons, want 3", len(scoped))
	}
	for i, want := range []int64{both, company, category} {
		if scoped[i].ID != want {
			t.Fatalf("rank %d: lesson %d (%s), want %d; got %+v",
				i, scoped[i].ID, scoped[i].Cue, want, scoped)
		}
	}
	topTwo, err := m.pg.TaskLessons(context.Background(), taskID, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(topTwo) != 2 || topTwo[0].ID != both || topTwo[1].ID != company {
		t.Fatalf("LIMIT displaced scoped lessons: %+v", topTwo)
	}

	// A task with no company or category may only see generic lessons.
	genericTask, err := m.CreateTask("generic ranking", "no company", nil, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	genericTaskID, err := strconv.ParseInt(genericTask.ID, 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = m.pg.DeleteTask(genericTaskID) }()
	generic, err := m.pg.TaskLessons(context.Background(), genericTaskID, 30)
	if err != nil {
		t.Fatal(err)
	}
	for _, lesson := range generic {
		if lesson.CategoryID != 0 || lesson.CompanyID != 0 {
			t.Fatalf("unscoped task received a scoped lesson: %+v", lesson)
		}
	}
}
