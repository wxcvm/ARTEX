package server

import (
    "io"
    "net/http"

    "github.com/Autumn-27/artex/db"
)

// Import only explicit DSH human-review feedback. It never changes ARTEX's
// own vulnerability disposition or attempts to merge source databases.
func (s *Server) importDSHReviewFeedback(w http.ResponseWriter, r *http.Request) {
    if s.m.pg == nil { writeErr(w, 503, "database unavailable"); return }
    const maxBytes = 1 << 20
    body, err := io.ReadAll(io.LimitReader(r.Body, maxBytes+1))
    if err != nil { writeErr(w, 400, err.Error()); return }
    if len(body) > maxBytes { writeErr(w, 413, "review feedback exceeds 1 MiB"); return }
    input, err := db.ParseDSHReviewFeedback(body)
    if err != nil { writeErr(w, 400, err.Error()); return }
    outcome, err := s.m.pg.ImportDSHReviewFeedback(input)
    if err != nil { writeErr(w, 500, err.Error()); return }
    writeJSON(w, 200, outcome)
}

func (s *Server) listDSHReviewFeedback(w http.ResponseWriter, r *http.Request) {
    if s.m.pg == nil { writeErr(w, 503, "database unavailable"); return }
    rows, err := s.m.pg.ListDSHReviewFeedback(50)
    if err != nil { writeErr(w, 500, err.Error()); return }
    writeJSON(w, 200, map[string]any{"items": rows, "total_shown": len(rows)})
}
