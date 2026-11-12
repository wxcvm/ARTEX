package server

import (
	"net/http"

	"github.com/Autumn-27/artex/agent"
)

// modelOutboundAudit is an administrator-only snapshot of completed model HTTP
// attempts. The scoped Purge bridge credential is not authorized here.
func (s *Server) modelOutboundAudit(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, agent.ModelOutboundAudit())
}
