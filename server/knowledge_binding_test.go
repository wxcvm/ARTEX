package server

import (
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestKnowledgeBindingMigrationPreservesCustomizedTools(t *testing.T) {
	dsn, _, err := db.DSN()
	if err != nil {
		t.Skipf("no database config: %v", err)
	}
	pg, err := db.Open(dsn)
	if err != nil {
		t.Skipf("postgres unavailable: %v", err)
	}
	defer pg.Close()
	tx, err := pg.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	for _, q := range []string{
		`CREATE TEMP TABLE settings (key text PRIMARY KEY, value text)`,
		`CREATE TEMP TABLE tools (key text PRIMARY KEY, system bool, enabled bool, agents jsonb)`,
		`INSERT INTO tools VALUES ('kb_search',true,true,'["planner","worker"]'), ('kb_get',true,false,'["planner","worker"]'), ('other',true,true,'["planner","worker"]')`,
	} {
		if _, err := tx.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	if err := migrateKnowledgeMainBindingTx(tx); err != nil {
		t.Fatal(err)
	}
	assertBound := func(key string, want bool) {
		t.Helper()
		var bound bool
		if err := tx.QueryRow(`SELECT agents ? 'mainagent' FROM tools WHERE key=$1`, key).Scan(&bound); err != nil || bound != want {
			t.Fatalf("%s: bound=%v, want=%v, err=%v", key, bound, want, err)
		}
	}
	assertBound("kb_search", true)
	assertBound("kb_get", false)
	assertBound("other", false)
	if _, err := tx.Exec(`UPDATE tools SET agents='["planner","worker"]' WHERE key='kb_search'`); err != nil {
		t.Fatal(err)
	}
	if err := migrateKnowledgeMainBindingTx(tx); err != nil {
		t.Fatal(err)
	}
	assertBound("kb_search", false)
}
