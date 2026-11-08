package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestTaskLessonEvaluationCohorts(t *testing.T) {
	m, err := NewManager(t.TempDir(), "")
	if err != nil {
		t.Skipf("postgres unavailable (%v) - skipping", err)
	}
	defer m.Close()

	// A unique category and company keep the comparison independent of
	// unrelated, completed tasks already present in a developer database.
	seed, err := m.CreateTask("lesson evaluation seed", "compare outcomes", nil, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	var created []int64
	seedID, _ := strconv.ParseInt(seed.ID, 10, 64)
	created = append(created, seedID)

	var categoryID, companyID, otherCompanyID int64
	key := "lesson-evaluation-" + seed.ID
	if err := m.pg.QueryRow("INSERT INTO task_categories(name,nkey) VALUES($1,$1) RETURNING id",
		key+"-category").Scan(&categoryID); err != nil {
		t.Fatal(err)
	}
	if err := m.pg.QueryRow("INSERT INTO companies(name,nkey) VALUES($1,$1) RETURNING id",
		key+"-company").Scan(&companyID); err != nil {
		t.Fatal(err)
	}
	if err := m.pg.QueryRow("INSERT INTO companies(name,nkey) VALUES($1,$1) RETURNING id",
		key+"-other").Scan(&otherCompanyID); err != nil {
		t.Fatal(err)
	}
	var lessonID int64
	defer func() {
		for _, id := range created {
			_ = m.pg.DeleteTask(id)
		}
		if lessonID != 0 {
			_, _ = m.pg.Exec("DELETE FROM rsi_lessons WHERE id=$1", lessonID)
		}
		_, _ = m.pg.Exec("DELETE FROM task_categories WHERE id=$1", categoryID)
		_, _ = m.pg.Exec("DELETE FROM companies WHERE id IN ($1,$2)", companyID, otherCompanyID)
	}()

	addTask := func(name string, compID int64, done, archived bool) (int64, int64) {
		t.Helper()
		var id, expID int64
		if name == "inject" {
			id = seedID
			expID = seed.ExpID
		} else {
			task, err := m.CreateTask("lesson evaluation "+name, "compare", nil, 0, 0)
			if err != nil {
				t.Fatal(err)
			}
			id, _ = strconv.ParseInt(task.ID, 10, 64)
			expID = task.ExpID
			created = append(created, id)
		}
		if _, err := m.pg.Exec(`
UPDATE tasks SET category_id=$1, company_id=$2,
 status=CASE WHEN $3::boolean THEN 'done' ELSE 'created' END,
 completed_at=CASE WHEN $3::boolean THEN now() ELSE NULL END,
 first_run_at=now() - interval '90 seconds',
 archived_at=CASE WHEN $4::boolean THEN now() ELSE NULL END
WHERE id=$5`, categoryID, compID, done, archived, id); err != nil {
			t.Fatal(err)
		}
		return id, expID
	}

	injected, expID := addTask("inject", companyID, true, false)
	shadow, _ := addTask("shadow", companyID, true, false)
	none, _ := addTask("none", companyID, true, false)
	mixed, _ := addTask("mixed", companyID, true, false)
	foreign, _ := addTask("foreign", otherCompanyID, true, false)
	running, _ := addTask("running", companyID, false, false)
	archived, _ := addTask("archived", companyID, true, true)

	if err := m.pg.QueryRow(`
INSERT INTO rsi_lessons(dedupe_key,kind,cue,action)
VALUES($1,'positive','xss','recorded prior finding') RETURNING id`,
		key).Scan(&lessonID); err != nil {
		t.Fatal(err)
	}
	if _, err := m.pg.Exec(`
INSERT INTO rsi_lesson_usage(lesson_id,task_id,mode)
VALUES ($1,$2,'inject'),($1,$3,'shadow'),($1,$4,'shadow'),($1,$4,'inject')`,
		lessonID, injected, shadow, mixed); err != nil {
		t.Fatal(err)
	}

	for _, item := range []struct {
		id     int64
		status string
	}{
		{injected, "confirmed"}, {injected, "fixed"}, {injected, "false_positive"},
		{shadow, "confirmed"}, {shadow, "duplicate"},
		{mixed, "confirmed"},
		{foreign, "confirmed"}, {running, "confirmed"}, {archived, "confirmed"},
	} {
		if _, err := m.pg.Exec(`
INSERT INTO findings(task_id,vulnclass,name,severity,status)
VALUES ($1,'test','fixture','medium',$2)`, item.id, item.status); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := m.pg.Exec(`
INSERT INTO tool_usage(tool_key,task_id) VALUES('test-tool',$1),('test-tool',$1)`,
		injected); err != nil {
		t.Fatal(err)
	}
	if _, err := m.pg.Exec(`
INSERT INTO skill_usage(skill,task_id,found)
VALUES('api-recon',$1,true),('missing',$1,false)`, injected); err != nil {
		t.Fatal(err)
	}
	if _, err := m.pg.Exec(`
INSERT INTO activity(exploration_id,kind,input_tokens,output_tokens)
VALUES ($1,'result',10,5), ($1,'tool_result',100,100)`, expID); err != nil {
		t.Fatal(err)
	}

	s := New(context.Background(), m, t.TempDir(), t.TempDir(), t.TempDir())
	token, err := signJWT(s.jwtKey)
	if err != nil {
		t.Fatal(err)
	}
	get := func(path string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		s.Handler().ServeHTTP(rec, req)
		return rec
	}

	rec := get("/api/tasks/" + strconv.FormatInt(injected, 10) + "/lessons/evaluation")
	if rec.Code != http.StatusOK {
		t.Fatalf("evaluation status=%d body=%s", rec.Code, rec.Body.String())
	}
	var result db.LessonEvaluation
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.TaskID != injected || result.CategoryID != categoryID || result.CompanyID != companyID {
		t.Fatalf("scope: %+v", result)
	}
	if len(result.Groups) != 4 {
		t.Fatalf("groups=%+v; want all four exposure kinds", result.Groups)
	}
	groups := map[string]db.LessonEvaluationGroup{}
	for _, group := range result.Groups {
		groups[group.Exposure] = group
		if group.Tasks != 1 {
			t.Fatalf("unexpected scope or status contamination in %+v", group)
		}
	}
	inj := groups["inject"]
	if inj.Findings != 3 || inj.Confirmed != 2 || inj.FalsePositive != 1 ||
		inj.ToolCalls != 2 || inj.SkillLoads != 1 || inj.SkillMisses != 1 ||
		inj.InputTokens != 10 || inj.OutputTokens != 5 {
		t.Fatalf("injected task metrics incorrect: %+v", inj)
	}
	if groups["shadow"].Confirmed != 1 || groups["shadow"].Duplicate != 1 {
		t.Fatalf("shadow task metrics incorrect: %+v", groups["shadow"])
	}
	if groups["mixed"].Confirmed != 1 || groups["none"].Findings != 0 {
		t.Fatalf("mixed and unrecorded groups: %+v", groups)
	}
	if rec := get("/api/tasks/not-a-number/lessons/evaluation"); rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid id returned %d", rec.Code)
	}
	if rec := get("/api/tasks/999999999/lessons/evaluation"); rec.Code != http.StatusNotFound {
		t.Fatalf("unknown task returned %d", rec.Code)
	}
}
