package server

import (
	"net/http"
	"runtime"

	"github.com/Autumn-27/artex/config"
)

type StoragePermissionItem struct {
	Kind string `json:"kind"`
	State string `json:"state"`
	BroadRead bool `json:"broad_read"`
	OtherRead bool `json:"other_read"`
	InheritanceEnabled bool `json:"inheritance_enabled"`
}

type StoragePermissions struct {
	Platform string `json:"platform"`
	Items []StoragePermissionItem `json:"items"`
}

// Fixed application-owned locations only. No caller-supplied path, recursive
// walk, file contents, permission changes, or raw error strings.
func (s *Server) storagePermissionsSnapshot() StoragePermissions {
	return StoragePermissions{Platform: runtime.GOOS, Items: []StoragePermissionItem{
		inspectStoragePermission("config_file", config.Path()),
		inspectStoragePermission("jwt_key", s.privacyJWTPath),
		inspectStoragePermission("data_root", s.privacyDataRoot),
	}}
}

func (s *Server) storagePermissions(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, s.storagePermissionsSnapshot())
}
