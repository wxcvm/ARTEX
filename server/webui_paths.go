package server

import (
	"path"
	"strings"
)

// Next's segment prefetch URLs flatten nested segment files with dots, while
// static export writes those same segments into directories. Exact-file
// serving is attempted first; this mapping is only the segment-file fallback.
func nextSegmentPath(requestPath string) string {
	name := path.Base(requestPath)
	if !strings.HasPrefix(name, "__next.") || !strings.HasSuffix(name, ".txt") {
		return ""
	}
	parts := strings.Split(strings.TrimSuffix(strings.TrimPrefix(name, "__next."), ".txt"), ".")
	if len(parts) < 2 {
		return ""
	}
	for _, part := range parts {
		if part == "" || part == "." || part == ".." || strings.ContainsAny(part, "/\\") {
			return ""
		}
	}
	return path.Join(path.Dir(requestPath), "__next."+parts[0], strings.Join(parts[1:], "/")+".txt")
}
