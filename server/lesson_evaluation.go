package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// taskLessonEvaluation is an on-demand, read-only observational comparison.
// It never changes lesson state, task configuration, or agent prompts.
func (s *Server) taskLessonEvaluation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeErr(w, 400, "invalid task id")
		return
	}
	report, err := s.m.PG().TaskLessonEvaluation(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		writeErr(w, 404, "task not found")
		return
	}
	if err != nil {
		log.Printf("[rsi] lesson evaluation failed for task %d: %v", id, err)
		writeErr(w, 500, "lesson evaluation failed")
		return
	}
	writeJSON(w, 200, report)
}
