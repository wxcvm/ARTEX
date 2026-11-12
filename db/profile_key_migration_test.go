package db

import (
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func TestMigrateModelKeysPostgres(t *testing.T) {
	// Never migrate the developer's real model table, even if their test DSN
	// points at an existing database. Restrict every connection to pg_temp.
	cfg, err := pgx.ParseConfig(testDSN(t))
	if err != nil {
		t.Fatal(err)
	}
	cfg.RuntimeParams["search_path"] = "pg_temp"
	sqlDB := stdlib.OpenDB(*cfg)
	defer sqlDB.Close()
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	if err := sqlDB.PingContext(t.Context()); err != nil {
		t.Skipf("postgres unavailable: %v", err)
	}
	d := &DB{sqlDB}
	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := d.ExecContext(t.Context(), query, args...); err != nil {
			t.Fatal(err)
		}
	}
	exec(`CREATE TEMP TABLE llm_profiles (
		id bigint PRIMARY KEY, api_key text, api_key_ciphertext text,
		api_key_hint text, model text NOT NULL,
		CONSTRAINT test_reject_second_migration CHECK (id <> 2 OR api_key IS NOT NULL)
	)`)
	key := strings.Repeat("11", 32)
	t.Setenv(profileSecretEnv, key)
	_, encrypted, err := encodeProfileSecret("already-encrypted-fixture")
	if err != nil {
		t.Fatal(err)
	}
	exec(`INSERT INTO llm_profiles VALUES
		(1, 'legacy-fixture-one', NULL, 'old', 'model-a'),
		(2, 'tiny', NULL, 'old', 'model-b'),
		(3, NULL, $1, 'old', 'model-c'),
		(4, 'mixed-fixture', $1, 'old', 'model-d'),
		(5, NULL, NULL, NULL, 'ollama')`, encrypted)

	type storage struct{ plain, encrypted, hint, model string }
	read := func(id int64) storage {
		t.Helper()
		var s storage
		if err := d.QueryRowContext(t.Context(), `SELECT COALESCE(api_key,''), COALESCE(api_key_ciphertext,''),
			COALESCE(api_key_hint,''), model FROM llm_profiles WHERE id=$1`, id).Scan(&s.plain, &s.encrypted, &s.hint, &s.model); err != nil {
			t.Fatal(err)
		}
		return s
	}
	original := make(map[int64]storage)
	for id := int64(1); id <= 5; id++ {
		original[id] = read(id)
	}
	unchanged := func() {
		t.Helper()
		for id, want := range original {
			if got := read(id); got != want {
				t.Fatalf("fixture row %d changed unexpectedly", id)
			}
		}
	}
	ids := []string{"1", "2"}
	t.Setenv(profileSecretEnv, "")
	if _, err := d.MigrateModelKeys(t.Context(), ids); !errors.Is(err, ErrModelKeyMigrationKey) {
		t.Fatalf("missing key: %v", err)
	}
	unchanged()
	t.Setenv(profileSecretEnv, strings.Repeat("22", 32))
	if _, err := d.MigrateModelKeys(t.Context(), ids); !errors.Is(err, ErrModelKeyMigrationIntegrity) {
		t.Fatalf("wrong key: %v", err)
	}
	unchanged()
	t.Setenv(profileSecretEnv, key)
	for _, stale := range [][]string{{"1"}, {"1", "9"}, {"1", "1"}} {
		if _, err := d.MigrateModelKeys(t.Context(), stale); !errors.Is(err, ErrModelKeyMigrationChanged) {
			t.Fatalf("stale preview: %v", err)
		}
		unchanged()
	}
	// The second write fails after the first succeeds. Both must roll back.
	if _, err := d.MigrateModelKeys(t.Context(), ids); err == nil {
		t.Fatal("expected transaction failure")
	}
	unchanged()
	exec("ALTER TABLE llm_profiles DROP CONSTRAINT test_reject_second_migration")

	result, err := d.MigrateModelKeys(t.Context(), ids)
	if err != nil || result.Migrated != 2 || result.SkippedMixed != 1 {
		t.Fatalf("migration result: %+v, err=%v", result, err)
	}
	for _, id := range []int64{1, 2} {
		got := read(id)
		decoded, err := decodeProfileSecret("", got.encrypted)
		if err != nil || decoded != original[id].plain || got.plain != "" || got.model != original[id].model {
			t.Fatalf("migration failed to preserve key or settings for row %d", id)
		}
		if id == 2 && got.hint != "••••" {
			t.Fatal("short key hint must not reveal the complete key")
		}
	}
	for _, id := range []int64{3, 4, 5} {
		if got := read(id); got != original[id] {
			t.Fatalf("non-migrated row %d changed", id)
		}
	}
	result, err = d.MigrateModelKeys(t.Context(), []string{})
	if err != nil || result.Migrated != 0 || result.SkippedMixed != 1 {
		t.Fatalf("repeat migration: %+v, err=%v", result, err)
	}
}
