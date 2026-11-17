package agent

import (
	"os"
	"strings"
)

// untrustedBoundaryEnabled frames tool output as data by default.
// Explicit opt-out preserves the legacy unframed behavior when required.
func untrustedBoundaryEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("ARTEX_UNTRUSTED_BOUNDARY"))) {
	case "0", "false", "off", "no":
		return false
	}
	return true
}
