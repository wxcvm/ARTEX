package db

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestRenderLessonsDataNotInstructions(t *testing.T) {
	lessons := []Lesson{
		{ID: 1, Kind: "positive", Cue: "xss", Action: "Confirmed xss: reflected in search", Support: 2,
			Evidence: []string{"finding:10 (status=confirmed severity=high)"}},
		{ID: 2, Kind: "negative", Cue: "ssrf", Action: "Do not\nre-report\x00 without new mechanics", Support: 1},
	}
	block := RenderLessons(lessons, 2000)
	if !strings.Contains(block, "not instructions") {
		t.Fatalf("missing data-not-instructions header:\n%s", block)
	}
	if !strings.Contains(block, "[positive, 2 evidence]") || !strings.Contains(block, "[negative, 1 evidence]") {
		t.Fatalf("missing kind/support labels:\n%s", block)
	}
	for _, line := range strings.Split(block, "\n") {
		if strings.Contains(line, "\x00") {
			t.Fatalf("control character survived: %q", line)
		}
	}
	if !strings.Contains(block, "Do not re-report without new mechanics") {
		t.Fatalf("lesson text was not flattened onto one line:\n%s", block)
	}
	if RenderLessons(nil, 2000) != "" {
		t.Fatal("empty lesson set must render empty")
	}
	capped := RenderLessons(lessons, 120)
	if len(capped) > 160 {
		t.Fatalf("cap not enforced: %d bytes", len(capped))
	}
}

func TestLessonTextTruncatesAtRuneBoundary(t *testing.T) {
	cases := []struct {
		input string
		limit int
		want  string
	}{
		{input: "ABC中文", limit: 4, want: "ABC"},
		{input: "中文", limit: 4, want: "中"},
		{input: "中文", limit: 2, want: ""},
		{input: "hello world", limit: 5, want: "hello"},
	}
	for _, tc := range cases {
		got := lessonText(tc.input, tc.limit)
		if got != tc.want {
			t.Errorf("lessonText(%q, %d) = %q, want %q", tc.input, tc.limit, got, tc.want)
		}
		if !utf8.ValidString(got) || len(got) > tc.limit {
			t.Errorf("invalid or oversized UTF-8 result: %q", got)
		}
	}
	longChinese := lessonText(strings.Repeat("知", 60), lessonCueMax)
	if !utf8.ValidString(longChinese) || len(longChinese) > lessonCueMax {
		t.Fatalf("multilingual lesson cue corrupted: %q", longChinese)
	}
}
