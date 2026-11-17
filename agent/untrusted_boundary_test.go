package agent

import "testing"

func TestUntrustedBoundaryEnabled(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value string
		want  bool
	}{
		{"default", "", true},
		{"enabled-numeric", "1", true},
		{"enabled-word", "true", true},
		{"enabled-on", "on", true},
		{"enabled-case-and-space", " TRUE ", true},
		{"disabled-numeric", "0", false},
		{"disabled-word", "false", false},
		{"disabled-off", "off", false},
		{"disabled-no", "no", false},
		{"disabled-case-and-space", " OFF ", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ARTEX_UNTRUSTED_BOUNDARY", tc.value)
			if got := untrustedBoundaryEnabled(); got != tc.want {
				t.Fatalf("ARTEX_UNTRUSTED_BOUNDARY=%q: got %v, want %v", tc.value, got, tc.want)
			}
		})
	}
}
