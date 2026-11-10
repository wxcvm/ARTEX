package agent

import (
	"os"
	"strings"
)

// untrustedBoundaryEnabled is opt-in; the default agent behavior is unchanged.
func untrustedBoundaryEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("ARTEX_UNTRUSTED_BOUNDARY"))) {
	case "1", "true", "on":
		return true
	}
	return false
}
