package server

import (
	"errors"
	"net/http"

	"github.com/Autumn-27/artex/db"
)

// This is an administrator-only read-only diagnostic. The scoped Purge bridge
// token remains restricted to the existing listener-status endpoint.
func (s *Server) modelKeyDiagnostics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if s.m == nil || s.m.pg == nil {
		writeErr(w, http.StatusServiceUnavailable, "model key diagnostics: database unavailable")
		return
	}
	status, err := s.m.pg.ModelKeyDiagnostics(r.Context())
	if err != nil {
		// Avoid exposing database errors or connection configuration.
		writeErr(w, http.StatusServiceUnavailable, "model key diagnostics: database read failed")
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (s *Server) migrateModelKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if s.m == nil || s.m.pg == nil {
		writeErr(w, http.StatusServiceUnavailable, "model key migration: database unavailable")
		return
	}
	var body struct {
		ExpectedIDs []string `json:"expected_ids"`
	}
	if err := decode(r, &body); err != nil || body.ExpectedIDs == nil {
		writeErr(w, http.StatusBadRequest, "expected_ids must be provided from the diagnostic preview")
		return
	}
	result, err := s.m.pg.MigrateModelKeys(r.Context(), body.ExpectedIDs)
	switch {
	case errors.Is(err, db.ErrModelKeyMigrationKey):
		writeErr(w, http.StatusConflict, "configure ARTEX_LLM_SECRET_KEY before migration")
	case errors.Is(err, db.ErrModelKeyMigrationIntegrity):
		writeErr(w, http.StatusConflict, "existing ciphertext failed authentication; check the original encryption key")
	case errors.Is(err, db.ErrModelKeyMigrationChanged):
		writeErr(w, http.StatusConflict, "model key storage changed; refresh the diagnostic preview")
	case err != nil:
		writeErr(w, http.StatusServiceUnavailable, "model key migration failed; refresh diagnostics before retrying")
	default:
		writeJSON(w, http.StatusOK, result)
	}
}
