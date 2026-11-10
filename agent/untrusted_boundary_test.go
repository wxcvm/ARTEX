package agent

import "testing"

func TestUntrustedBoundaryEnabled(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  bool
	}{
		{"", false}, {"0", false}, {"false", false}, {"yes", false},
		{"1", true}, {"true", true}, {"on", true}, {" TRUE ", true},
	} {
		t.Run(tc.value, func(t *testing.T) {
			t.Setenv("ARTEX_UNTRUSTED_BOUNDARY", tc.value)
			if got := untrustedBoundaryEnabled(); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}
