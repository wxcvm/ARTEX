package server

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Autumn-27/artex/knowledge"
)

// ARTEX_KNOWLEDGE_FILE enables a read-only import preview: no DB writes,
// tool registration or promotion of external claims into task evidence.
func (s *Server) knowledgeCatalog(w http.ResponseWriter, r *http.Request) (*knowledge.Catalog, bool) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Context().Err() != nil {
		writeErr(w, http.StatusRequestTimeout, "Knowledge request cancelled")
		return nil, false
	}
	file := strings.TrimSpace(os.Getenv("ARTEX_KNOWLEDGE_FILE"))
	if file == "" {
		writeErr(w, http.StatusServiceUnavailable, "Knowledge preview is not configured")
		return nil, false
	}
	c, err := knowledge.Load(file)
	if err != nil {
		// A broken catalog must not be shown as an empty library. Do not leak paths.
		writeErr(w, http.StatusServiceUnavailable, "Knowledge catalog could not be read or validated")
		return nil, false
	}
	if r.Context().Err() != nil {
		writeErr(w, http.StatusRequestTimeout, "Knowledge request cancelled")
		return nil, false
	}
	return c, true
}
func (s *Server) knowledgeStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if strings.TrimSpace(os.Getenv("ARTEX_KNOWLEDGE_FILE")) == "" {
		writeJSON(w, http.StatusOK, map[string]any{"enabled": false, "mode": "read_only_preview"})
		return
	}
	if c, ok := s.knowledgeCatalog(w, r); ok {
		writeJSON(w, http.StatusOK, map[string]any{"enabled": true, "mode": "read_only_preview", "source": c.Source, "summary": c.Summary})
	}
}
func (s *Server) knowledgeList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	offset, limit := 0, 20
	for key, target := range map[string]*int{"offset": &offset, "limit": &limit} {
		if q.Has(key) {
			n, err := strconv.Atoi(q.Get(key))
			if err != nil || n < 0 || (key == "limit" && (n < 1 || n > 100)) {
				writeErr(w, http.StatusBadRequest, "Invalid knowledge pagination")
				return
			}
			*target = n
		}
	}
	if len(q.Get("q")) > 512 || (q.Get("conflicts") != "" && q.Get("conflicts") != "true" && q.Get("conflicts") != "false") {
		writeErr(w, http.StatusBadRequest, "Invalid knowledge search parameters")
		return
	}
	if c, ok := s.knowledgeCatalog(w, r); ok {
		items, total := c.Search(q.Get("q"), q.Get("conflicts") == "true", offset, limit)
		writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": total, "offset": offset, "limit": limit, "source": c.Source})
	}
}
func (s *Server) knowledgeEntry(w http.ResponseWriter, r *http.Request) {
	if c, ok := s.knowledgeCatalog(w, r); ok {
		for _, e := range c.Entries {
			if e.ID == r.PathValue("id") {
				writeJSON(w, http.StatusOK, map[string]any{"entry": e, "source": c.Source})
				return
			}
		}
		writeErr(w, http.StatusNotFound, "Knowledge entry not found")
	}
}
