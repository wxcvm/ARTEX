//go:build !windows

package server

func inspectStoragePermission(kind, _ string) StoragePermissionItem {
	return StoragePermissionItem{Kind: kind, State: "unsupported"}
}
