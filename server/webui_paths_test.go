package server

import "testing"

func TestOnboardingStaticSegmentPaths(t *testing.T) {
	cases := map[string]string{
		"welcome/__next.!KG1haW4p.welcome.__PAGE__.txt":          "welcome/__next.!KG1haW4p/welcome/__PAGE__.txt",
		"system/knowledge/__next.!KG1haW4p.system.knowledge.txt": "system/knowledge/__next.!KG1haW4p/system/knowledge.txt",
		"welcome/__next._full.txt":                               "",
		"welcome/index.html":                                     "",
		"welcome/__next.a..txt":                                  "",
		"welcome/__next.a.b.js":                                  "",
	}
	for input, want := range cases {
		if got := nextSegmentPath(input); got != want {
			t.Errorf("%s: got %q want %q", input, got, want)
		}
	}
}
