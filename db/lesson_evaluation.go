package db

import (
	"context"
)

// LessonEvaluationGroup reports observed outcomes across completed tasks in a
// single exposure group. Totals are descriptive, not causal attribution.
type LessonEvaluationGroup struct {
	Exposure       string `json:"exposure"` // inject | shadow | mixed | none
	Tasks          int64  `json:"tasks"`
	Findings       int64  `json:"findings"`
	Confirmed      int64  `json:"confirmed"`
	FalsePositive  int64  `json:"false_positive"`
	Duplicate      int64  `json:"duplicate"`
	ToolCalls      int64  `json:"tool_calls"`
	SkillLoads     int64  `json:"skill_loads"`
	SkillMisses    int64  `json:"skill_misses"`
	InputTokens    int64  `json:"input_tokens"`
	OutputTokens   int64  `json:"output_tokens"`
	WallSeconds    int64  `json:"wall_seconds"`
}

type LessonEvaluation struct {
	TaskID     int64                   `json:"task_id"`
	CategoryID int64                   `json:"category_id"`
	CompanyID  int64                   `json:"company_id"`
	Groups     []LessonEvaluationGroup `json:"groups"`
}

// TaskLessonEvaluation compares completed, unarchived tasks with the same
// company and category as taskID (including matching NULL scopes).
// "none" means no recorded lesson exposure, not a randomized control.
// A task that was both shadowed and injected has its own "mixed" group.
func (d *DB) TaskLessonEvaluation(ctx context.Context, taskID int64) (LessonEvaluation, error) {
	categoryID, companyID, err := taskLessonScope(ctx, d, taskID)
	if err != nil {
		return LessonEvaluation{}, err
	}
	result := LessonEvaluation{
		TaskID: taskID, CategoryID: categoryID, CompanyID: companyID,
		Groups: []LessonEvaluationGroup{},
	}
	rows, err := d.QueryContext(ctx, `
		WITH eligible AS (
			SELECT id, exploration_id, completed_at,
			       COALESCE(first_run_at, created_at) AS started_at
			FROM tasks
			WHERE status='done' AND completed_at IS NOT NULL
			  AND deleted_at IS NULL AND archived_at IS NULL
			  AND company_id IS NOT DISTINCT FROM NULLIF($1::bigint,0)
			  AND category_id IS NOT DISTINCT FROM NULLIF($2::bigint,0)
		), metrics AS (
			SELECT
				CASE
					WHEN lu.injected > 0 AND lu.shadowed > 0 THEN 'mixed'
					WHEN lu.injected > 0 THEN 'inject'
					WHEN lu.shadowed > 0 THEN 'shadow'
					ELSE 'none'
				END AS exposure,
				f.findings, f.confirmed, f.false_positive, f.duplicate,
				tu.calls AS tool_calls,
				su.loads AS skill_loads, su.misses AS skill_misses,
				a.input_tokens, a.output_tokens,
				GREATEST(0, EXTRACT(EPOCH FROM (t.completed_at - t.started_at))::bigint) AS wall_seconds
			FROM eligible t
			CROSS JOIN LATERAL (
				SELECT COUNT(*) AS findings,
				       COUNT(*) FILTER (WHERE status IN ('confirmed','resolved','fixed')) AS confirmed,
				       COUNT(*) FILTER (WHERE status='false_positive') AS false_positive,
				       COUNT(*) FILTER (WHERE status='duplicate') AS duplicate
				FROM findings WHERE task_id=t.id
			) f
			CROSS JOIN LATERAL (
				SELECT COUNT(*) AS calls FROM tool_usage WHERE task_id=t.id
			) tu
			CROSS JOIN LATERAL (
				SELECT COUNT(*) FILTER (WHERE found) AS loads,
				       COUNT(*) FILTER (WHERE NOT found) AS misses
				FROM skill_usage WHERE task_id=t.id
			) su
			CROSS JOIN LATERAL (
				SELECT COALESCE(SUM(input_tokens),0)::bigint AS input_tokens,
				       COALESCE(SUM(output_tokens),0)::bigint AS output_tokens
				FROM activity WHERE exploration_id=t.exploration_id AND kind='result'
			) a
			CROSS JOIN LATERAL (
				SELECT COUNT(*) FILTER (WHERE mode='inject') AS injected,
				       COUNT(*) FILTER (WHERE mode='shadow') AS shadowed
				FROM rsi_lesson_usage WHERE task_id=t.id
			) lu
		)
		SELECT exposure, COUNT(*)::bigint,
		       COALESCE(SUM(findings),0)::bigint,
		       COALESCE(SUM(confirmed),0)::bigint,
		       COALESCE(SUM(false_positive),0)::bigint,
		       COALESCE(SUM(duplicate),0)::bigint,
		       COALESCE(SUM(tool_calls),0)::bigint,
		       COALESCE(SUM(skill_loads),0)::bigint,
		       COALESCE(SUM(skill_misses),0)::bigint,
		       COALESCE(SUM(input_tokens),0)::bigint,
		       COALESCE(SUM(output_tokens),0)::bigint,
		       COALESCE(SUM(wall_seconds),0)::bigint
		FROM metrics GROUP BY exposure
		ORDER BY CASE exposure
			WHEN 'inject' THEN 0 WHEN 'shadow' THEN 1
			WHEN 'mixed' THEN 2 ELSE 3 END`, companyID, categoryID)
	if err != nil {
		return LessonEvaluation{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var g LessonEvaluationGroup
		if err := rows.Scan(
			&g.Exposure, &g.Tasks, &g.Findings, &g.Confirmed,
			&g.FalsePositive, &g.Duplicate, &g.ToolCalls,
			&g.SkillLoads, &g.SkillMisses, &g.InputTokens,
			&g.OutputTokens, &g.WallSeconds,
		); err != nil {
			return LessonEvaluation{}, err
		}
		result.Groups = append(result.Groups, g)
	}
	return result, rows.Err()
}
