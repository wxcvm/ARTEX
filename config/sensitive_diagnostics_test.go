package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestInspectConfigURL(t *testing.T) {
	tests := []struct {
		url, field string
		schemes []string
		want []string
	}{
		{"https://user:synthetic-password@example.test/v1?api_key=synthetic-key", "base_url", []string{"http", "https"}, []string{"url_userinfo", "sensitive_query"}},
		{"http://example.test/v1", "base_url", []string{"http", "https"}, []string{"http_non_loopback"}},
		{"http://127.0.0.1:11434/v1", "base_url", []string{"http", "https"}, nil},
		{"http://[::1]:11434/v1", "base_url", []string{"http", "https"}, nil},
		{"socks5://user:synthetic-proxy@localhost:1080", "proxy", []string{"http", "https", "socks5"}, []string{"url_userinfo"}},
		{"postgres://u:p@localhost/db?sslmode=require", "database.dsn", []string{"postgres", "postgresql"}, []string{"url_userinfo"}},
		{"host=localhost password=synthetic-password", "database.dsn", []string{"postgres", "postgresql"}, []string{"url_not_inspected"}},
		{"https://example.test/v1?api_key=%ZZ", "base_url", []string{"http", "https"}, []string{"query_not_inspected"}},
		{"", "base_url", []string{"http", "https"}, nil},
	}
	for _, tc := range tests {
		got := InspectConfigURL(tc.url, "model_profile", tc.field, "42", tc.schemes...)
		if len(got) != len(tc.want) { t.Fatalf("unexpected finding count for %s: %v", tc.field, got) }
		for i, code := range tc.want {
			if got[i].Code != code || got[i].ProfileID != "42" { t.Fatalf("unexpected finding: %+v", got[i]) }
		}
		data, err := json.Marshal(got)
		if err != nil || strings.Contains(string(data), "synthetic-") || strings.Contains(string(data), "example.test") { t.Fatal("diagnostic exposed a value") }
	}
}

func TestSensitiveConfigDiagnostics(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("ARTEX_CONFIG", path)
	t.Setenv("ARTEX_PG_DSN", "postgres://env:synthetic-env@localhost/db")
	t.Setenv("ARTEX_LLM_BASE_URL", "https://example.test/v1?token=synthetic-token")
	t.Setenv("ARTEX_LLM_PROXY", "socks5://u:synthetic-proxy@localhost:1080")
	raw := []byte(`{"database":{"password":"synthetic-file","dsn":"postgres://u:synthetic-dsn@localhost/db"}}`)
	if err := os.WriteFile(path, raw, 0600); err != nil { t.Fatal(err) }
	s := SensitiveConfigDiagnostics()
	if s.FileState != "readable" || s.DatabaseSource != "environment" || len(s.Findings) != 5 {
		t.Fatalf("unexpected diagnostics: %+v", s)
	}
	if runtime.GOOS == "windows" {
		if s.FilePermissions != "windows_acl_not_inspected" { t.Fatal("Windows permissions incorrectly asserted") }
	} else if s.FilePermissions != "no_group_or_other_read" { t.Fatalf("unexpected permissions: %s", s.FilePermissions) }
	b, err := json.Marshal(s)
	if err != nil || strings.Contains(string(b), "synthetic-") || strings.Contains(string(b), "example.test") || strings.Contains(string(b), path) {
		t.Fatal("diagnostic exposed sensitive configuration")
	}
	after, err := os.ReadFile(path)
	if err != nil || string(after) != string(raw) { t.Fatal("configuration was modified") }
	t.Setenv("ARTEX_PG_DSN", "")
	if got := SensitiveConfigDiagnostics(); got.DatabaseSource != "file_dsn" { t.Fatalf("wrong precedence: %+v", got) }
	if err := os.WriteFile(path, []byte("{"), 0600); err != nil { t.Fatal(err) }
	if got := SensitiveConfigDiagnostics(); got.FileState != "invalid" { t.Fatalf("invalid file incorrectly classified: %+v", got) }
	if err := os.Remove(path); err != nil { t.Fatal(err) }
	if got := SensitiveConfigDiagnostics(); got.FileState != "missing" { t.Fatalf("missing file incorrectly classified: %+v", got) }
}
