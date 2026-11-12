package server

import (
	"net/http"
	"strconv"

	"github.com/Autumn-27/artex/config"
)

// sensitiveConfigDiagnostics reports metadata only. It does not resolve
// endpoints, open connections to providers, or change stored configuration.
func (s *Server) sensitiveConfigDiagnostics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	status := config.SensitiveConfigDiagnostics()
	if s.m == nil || s.m.pg == nil {
		status.ProfilesState = "unavailable"
		writeJSON(w, http.StatusOK, status)
		return
	}
	rows, err := s.m.pg.QueryContext(r.Context(), `SELECT id, COALESCE(base_url,''), COALESCE(proxy,'') FROM llm_profiles ORDER BY id`)
	if err != nil {
		status.ProfilesState = "unavailable"
		writeJSON(w, http.StatusOK, status)
		return
	}
	defer rows.Close()
	status.ProfilesState = "complete"
	for rows.Next() {
		var id int64
		var baseURL, proxy string
		if err := rows.Scan(&id, &baseURL, &proxy); err != nil {
			status.ProfilesState = "partial"
			break
		}
		profileID := strconv.FormatInt(id, 10)
		status.ProfilesChecked++
		status.Findings = append(status.Findings, config.InspectConfigURL(baseURL, "model_profile", "base_url", profileID, "http", "https")...)
		status.Findings = append(status.Findings, config.InspectConfigURL(proxy, "model_profile", "proxy", profileID, "http", "https", "socks5")...)
	}
	if rows.Err() != nil {
		status.ProfilesState = "partial"
	}
	writeJSON(w, http.StatusOK, status)
}
