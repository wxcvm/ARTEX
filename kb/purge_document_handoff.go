package kb

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

var ErrInvalidPurgeHandoff = errors.New("invalid Purge handoff filename or content digest")

const purgeHandoffPrefix = "purge-document-handoff-"

var purgeHandoffName = regexp.MustCompile("^" + purgeHandoffPrefix + "([0-9a-f]{16})(?: \\([1-9][0-9]{0,8}\\))?\\.md$")

type PurgeHandoff struct {
	Filename string
	Content  string
}

func CanonicalPurgeHandoffText(text string) string {
	return strings.ReplaceAll(strings.TrimPrefix(text, "\uFEFF"), "\r\n", "\n")
}

// ParsePurgeHandoff returns nil for ordinary uploads. Only a recognized
// handoff name claims the content digest; browser copy suffixes are discarded.
func ParsePurgeHandoff(filename, text string) (*PurgeHandoff, error) {
	name := strings.ToLower(filename)
	if !strings.HasPrefix(name, purgeHandoffPrefix) {
		return nil, nil
	}
	match := purgeHandoffName.FindStringSubmatch(name)
	if len(match) != 2 {
		return nil, ErrInvalidPurgeHandoff
	}
	content := CanonicalPurgeHandoffText(text)
	sum := sha256.Sum256([]byte(content))
	if !strings.HasPrefix(hex.EncodeToString(sum[:]), match[1]) {
		return nil, ErrInvalidPurgeHandoff
	}
	return &PurgeHandoff{Filename: purgeHandoffPrefix + match[1] + ".md", Content: content}, nil
}
