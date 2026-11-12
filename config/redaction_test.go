package config

import (
	"strings"
	"testing"
)

func TestRedactLogText(t *testing.T) {
	cases := []struct{ input, want string }{
		{"proxy socks5://user:proxy-pass@127.0.0.1:1080 failed", "proxy socks5://***@127.0.0.1:1080 failed"},
		{"postgres://dbuser:db-pass@localhost/db?password=other-pass&sslmode=require", "postgres://***@localhost/db?password=***&sslmode=require"},
		{"Authorization: Bearer secret-token.123", "Authorization: Bearer ***"},
		{"Proxy-Authorization: Basic c2VjcmV0OmFh==", "Proxy-Authorization: Basic ***"},
		{"request Bearer abc.def", "request Bearer ***"},
		{"https://example.test/v1?api_key=private&model=test", "https://example.test/v1?api_key=***&model=test"},
		{"password='a secret with spaces' host=localhost", "password='***' host=localhost"},
		{"{\"api_key\":\"secret with \\\"quotes\\\"\",\"model\":\"test\"}", "{\"api_key\":\"***\",\"model\":\"test\"}"},
		{"OPENAI_API_KEY=private-key", "OPENAI_API_KEY=***"},
		{"connection refused; model unavailable", "connection refused; model unavailable"},
	}
	for _, tc := range cases {
		got := RedactLogText(tc.input)
		if got != tc.want {
			t.Errorf("redaction mismatch: got %q want %q", got, tc.want)
		}
		if RedactLogText(got) != got {
			t.Errorf("redaction is not idempotent: %q", got)
		}
	}
}

func TestRedactDatabaseDSN(t *testing.T) {
	for _, input := range []string{
		"postgres://user:secret@localhost/db?password=another-secret&sslmode=require",
		"postgres://user:secret%zz@localhost/db",
		"host=localhost password='secret with spaces' user=user",
		"postgres://localhost/db?password=secret",
	} {
		got := Redact(input)
		if strings.Contains(got, "secret") || got == input {
			t.Errorf("DSN not redacted: %q", got)
		}
	}
	if got := Redact("postgres://user:secret@localhost:5432/db?sslmode=require"); got != "postgres://localhost:5432/db" {
		t.Errorf("safe connection context lost: %q", got)
	}
}
