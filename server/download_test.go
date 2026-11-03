package server

import (
	"mime"
	"strings"
	"testing"
)

func TestDownloadDisposition(t *testing.T) {
	for _, tt := range []struct {
		name     string
		filename string
		fallback string
	}{
		{"findings-20261102.csv", "findings-20261102.csv", "findings-20261102.csv"},
		{"漏洞报告.md", "漏洞报告.md", "____.md"},
		{"report final 1.md", "report final 1.md", "report final 1.md"},
		{`report "final"; 50%+1.md`, `report "final"; 50%+1.md`, "report _final_; 50%+1.md"},
		{"l'été.txt", "l'été.txt", "l'_t_.txt"},
		{"report\r\nX-Header: injected\x00.txt", "reportX-Header: injected.txt", "reportX-Header: injected.txt"},
		{"", "download", "download"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			header := downloadDisposition(tt.name)
			for _, r := range header {
				if r < 0x20 || r > 0x7e {
					t.Fatalf("header contains a non-printable-ASCII character: %q", header)
				}
			}
			kind, params, err := mime.ParseMediaType(header)
			if err != nil || kind != "attachment" || params["filename"] != tt.filename {
				t.Fatalf("parse %q = %q, %v, %v; want filename %q", header, kind, params, err, tt.filename)
			}
			legacy := strings.SplitN(header, "; filename*=", 2)[0]
			_, params, err = mime.ParseMediaType(legacy)
			if err != nil || params["filename"] != tt.fallback {
				t.Fatalf("legacy parse = %v, %v; want %q", params, err, tt.fallback)
			}
		})
	}
}
