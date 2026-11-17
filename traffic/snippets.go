package traffic

import (
	"database/sql"
	"strings"
	"unicode"
	"unicode/utf8"
)

// MatchSnippet is plain text, not HTML. The UI can mark the matching portion
// without interpreting captured response content as executable markup.
type MatchSnippet struct {
	Before string `json:"before"`
	Match  string `json:"match"`
	After  string `json:"after"`
}

// BodyMatchSnippets uses the original stored text, not FTS5 snippet(): ARTEX's
// trigram FTS is contentless. Only the current visible page is inspected.
func (t *Traffic) BodyMatchSnippets(page []ExchangeMeta, keyword string) map[string]MatchSnippet {
	out := make(map[string]MatchSnippet)
	keyword = strings.TrimSpace(keyword)
	if utf8.RuneCountInString(keyword) < minTrigram {
		return out
	}
	for _, meta := range page {
		req, resp, err := t.Get(meta.ID)
		if err != nil {
			continue // a deleted/legacy unreadable capture must not fail the list
		}
		if snippet, ok := matchTextSnippet(resp, keyword); ok {
			out[meta.ID] = snippet
			continue
		}
		if snippet, ok := matchTextSnippet(req, keyword); ok {
			out[meta.ID] = snippet
			continue
		}
		if snippet, ok := matchTextSnippet(meta.URL, keyword); ok {
			out[meta.ID] = snippet // URL is also part of the original FTS text
			continue
		}

		// Spilled bodies expose an 8 KiB inline preview; a match later in the
		// first maxIndexBody bytes is still indexed and deserves a real snippet.
		var reqHash, respHash sql.NullString
		if err := t.db.QueryRow("SELECT req_blob,resp_blob FROM exchange_bodies WHERE id=?", meta.ID).Scan(&reqHash, &respHash); err != nil {
			continue
		}
		for _, hash := range []sql.NullString{respHash, reqHash} {
			if !hash.Valid || hash.String == "" {
				continue
			}
			body, _, err := t.BlobRange(hash.String, 0, maxIndexBody)
			if err != nil {
				continue
			}
			if snippet, ok := matchTextSnippet(string(body), keyword); ok {
				out[meta.ID] = snippet
				break
			}
		}
	}
	return out
}

func matchTextSnippet(text, keyword string) (MatchSnippet, bool) {
	haystack := []rune(text)
	needle := []rune(keyword)
	if len(needle) == 0 || len(needle) > len(haystack) {
		return MatchSnippet{}, false
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		ok := true
		for j, r := range needle {
			if unicode.ToLower(haystack[i+j]) != unicode.ToLower(r) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		start := i - 64
		if start < 0 { start = 0 }
		end := i + len(needle) + 64
		if end > len(haystack) { end = len(haystack) }
		before := strings.Join(strings.Fields(string(haystack[start:i])), " ")
		after := strings.Join(strings.Fields(string(haystack[i+len(needle):end])), " ")
		if start > 0 { before = "…" + before }
		if end < len(haystack) { after += "…" }
		if before != "" { before += " " }
		if after != "" { after = " " + after }
		return MatchSnippet{Before: before, Match: string(haystack[i:i+len(needle)]), After: after}, true
	}
	return MatchSnippet{}, false
}
