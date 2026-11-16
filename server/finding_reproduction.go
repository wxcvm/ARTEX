package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/Autumn-27/artex/db"
)

func (s *Server) registerFindingReproduction(mux *http.ServeMux) {
	path := "/api/exploration/findings/{id}/reproduction"
	mux.HandleFunc("GET "+path, s.getFindingReproduction)
	mux.HandleFunc("PUT "+path, s.putFindingReproduction)
}

func (s *Server) getFindingReproduction(w http.ResponseWriter, r *http.Request) {
	id, ok := s.findingTrafficAccess(w, r, false)
	if !ok {
		return
	}
	steps, err := s.m.pg.ListFindingReproduction(r.Context(), id)
	if err != nil {
		evidenceError(w, err)
		return
	}
	writeJSON(w, 200, steps)
}

func (s *Server) putFindingReproduction(w http.ResponseWriter, r *http.Request) {
	id, ok := s.findingTrafficAccess(w, r, true)
	if !ok {
		return
	}
	var body struct {
		Steps []db.FindingReproductionStep `json:"steps"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&body); err != nil || body.Steps == nil {
		writeErr(w, 400, "steps array is required")
		return
	}
	for i := range body.Steps {
		step := &body.Steps[i]
		step.Title = strings.TrimSpace(step.Title)
		step.Observation = strings.TrimSpace(step.Observation)
		step.BindingID = strings.TrimSpace(step.BindingID)
		if step.Title == "" || utf8.RuneCountInString(step.Title) > 4000 || utf8.RuneCountInString(step.Observation) > 4000 {
			writeErr(w, 400, "step title is required; title and observation must be at most 4000 characters")
			return
		}
	}
	if err := s.m.pg.ReplaceFindingReproduction(r.Context(), id, body.Steps); err != nil {
		evidenceError(w, err)
		return
	}
	s.getFindingReproduction(w, r)
}
