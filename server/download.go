package server

import (
	"fmt"
	"net/url"
	"strings"
)

// downloadDisposition keeps a printable ASCII fallback for old clients and
// carries the original UTF-8 name in filename* (RFC 6266 / RFC 8187). Keep the
// findings export and workspace download headers consistent.
func downloadDisposition(name string) string {
	name = strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, name)
	if strings.TrimSpace(name) == "" {
		name = "download"
	}
	fallback := strings.Map(func(r rune) rune {
		if r > 0x7e || r == '"' || r == '\\' || r == '/' {
			return '_'
		}
		return r
	}, name)
	// QueryEscape escapes all parameter separators. A space must be %20, not
	// form encoding's '+', because filename* is not a URL query string.
	encoded := strings.ReplaceAll(url.QueryEscape(name), "+", "%20")
	return fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, fallback, encoded)
}
