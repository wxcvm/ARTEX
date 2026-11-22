package kb

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"testing"
)

func TestParsePurgeHandoff(t *testing.T) {
	const content = "# Document\n"
	sum := sha256.Sum256([]byte(content))
	name := fmt.Sprintf("purge-document-handoff-%x", sum[:8])
	for _, filename := range []string{name + ".md", name + " (2).MD"} {
		got, err := ParsePurgeHandoff(filename, "\uFEFF# Document\r\n")
		if err != nil || got == nil || got.Filename != name+".md" || got.Content != content {
			t.Fatalf("parse %q: got %+v %v", filename, got, err)
		}
	}
	if got, err := ParsePurgeHandoff("ordinary.md", content); got != nil || err != nil {
		t.Fatalf("ordinary file handled as a receipt: %+v %v", got, err)
	}
	for _, filename := range []string{name + ".md", "purge-document-handoff-invalid.md"} {
		if _, err := ParsePurgeHandoff(filename, "changed"); !errors.Is(err, ErrInvalidPurgeHandoff) {
			t.Fatalf("unverified receipt %q accepted: %v", filename, err)
		}
	}
}
