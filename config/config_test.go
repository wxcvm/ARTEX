package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPostgresDSNPrecedence(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	os.WriteFile(cfgPath, []byte(`{"database":{"host":"10.1.2.3","port":6000,"user":"u","password":"p","dbname":"d","sslmode":"require"}}`), 0o644)
	t.Setenv("ARTEX_CONFIG", cfgPath)

	// no env DSN → built from config file fields
	t.Setenv("ARTEX_PG_DSN", "")
	got, _, err := PostgresDSN()
	want := "postgres://u:p@10.1.2.3:6000/d?sslmode=require"
	if err != nil || got != want {
		t.Fatalf("from file: got %q err %v want %q", got, err, want)
	}

	// env wins over config file
	t.Setenv("ARTEX_PG_DSN", "postgres://envwins/x")
	if got, _, err := PostgresDSN(); err != nil || got != "postgres://envwins/x" {
		t.Fatalf("env should win, got %q err %v", got, err)
	}

	// no env, no file → error (no built-in default)
	t.Setenv("ARTEX_PG_DSN", "")
	t.Setenv("ARTEX_CONFIG", filepath.Join(dir, "nope.json"))
	if got, _, err := PostgresDSN(); err == nil {
		t.Fatalf("missing config should error, got %q", got)
	}

	// full dsn in config file is used verbatim
	os.WriteFile(cfgPath, []byte(`{"database":{"dsn":"postgres://full/dsn"}}`), 0o644)
	t.Setenv("ARTEX_CONFIG", cfgPath)
	if got, _, err := PostgresDSN(); err != nil || got != "postgres://full/dsn" {
		t.Fatalf("file dsn verbatim, got %q err %v", got, err)
	}
}

func TestPostgresDSNRejectsInvalidConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("ARTEX_CONFIG", path)
	t.Setenv("ARTEX_PG_DSN", "")
	for _, raw := range []string{
		`{"database":{"host":"localhost",`,
		`{"database":{"host":"localhost","user":"example","port":"secret-fixture"},"skill_dir":"partial"}`,
		`{"database":{"dsn":"postgres://example/fixture"},"knowledge_dir":42}`,
		`null`,
		`[]`,
	} {
		if err := os.WriteFile(path, []byte(raw), 0o600); err != nil {
			t.Fatal(err)
		}
		dsn, source, err := PostgresDSN()
		if err == nil || !strings.Contains(err.Error(), "invalid JSON configuration") || dsn != "" || source != "" {
			t.Fatalf("expected explicit configuration error, got source=%q err=%v", source, err)
		}
		if strings.Contains(err.Error(), "secret-fixture") {
			t.Fatal("configuration error exposed a value")
		}
		if got := Load(); got != (Config{}) {
			t.Fatal("invalid configuration returned partially decoded settings")
		}
	}
	t.Setenv("ARTEX_PG_DSN", "postgres://override/fixture")
	if got, _, err := PostgresDSN(); err != nil || got != "postgres://override/fixture" {
		t.Fatalf("environment override must remain authoritative: %v", err)
	}
}

func TestPostgresDSNReportsUnreadableConfig(t *testing.T) {
	t.Setenv("ARTEX_CONFIG", t.TempDir())
	t.Setenv("ARTEX_PG_DSN", "")
	if _, _, err := PostgresDSN(); err == nil || !strings.Contains(err.Error(), "read database config") {
		t.Fatalf("expected configuration read error, got %v", err)
	}
}
