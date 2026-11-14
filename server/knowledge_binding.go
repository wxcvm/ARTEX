package server

import (
	"database/sql"

	"github.com/Autumn-27/artex/db"
)

// Migrate only the old, enabled default binding. Preserve disabled tools and
// custom agent selections. The marker is written atomically with the change,
// so later manual edits are never overwritten on restart.
func migrateKnowledgeMainBinding(pg *db.DB) error {
	tx, err := pg.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := migrateKnowledgeMainBindingTx(tx); err != nil {
		return err
	}
	return tx.Commit()
}

func migrateKnowledgeMainBindingTx(tx *sql.Tx) error {
	result, err := tx.Exec(`INSERT INTO settings(key,value) VALUES ('migration.knowledge_main_binding.v1','true') ON CONFLICT (key) DO NOTHING`)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 0 {
		if _, err := tx.Exec(`UPDATE tools SET agents = agents || '["mainagent"]'::jsonb
WHERE key IN ('kb_search','kb_get') AND system AND enabled
AND agents = '["planner","worker"]'::jsonb`); err != nil {
			return err
		}
	}
	return nil
}
