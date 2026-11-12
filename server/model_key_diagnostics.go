package server

import "net/http"

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
