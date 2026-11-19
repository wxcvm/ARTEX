package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func restoreFindingRemediationTx(tx *sql.Tx, snapshot *TaskArchiveSnapshot, taskID int64) error {
	// These tables are optional in legacy archives. A finding without stored
	// state deliberately restores to implicit pending_review at version zero.
	for _, table := range []string{"finding_remediations", "finding_remediation_events"} {
		raw := snapshot.Tables[table]
		if len(raw) == 0 {
			continue
		}
		var entries []json.RawMessage
		if err := json.Unmarshal(raw, &entries); err != nil {
			return fmt.Errorf("%s: %w", table, err)
		}
		if len(entries) == 0 {
			continue
		}
		// Do not attach an archive's remediation history to findings that
		// belong to another live or archived task.
		result, err := tx.Exec(`INSERT INTO `+table+`
 SELECT archive_row.* FROM json_populate_recordset(NULL::`+table+`,$1::json) archive_row
 JOIN findings f ON f.id=archive_row.finding_id AND f.task_id=$2`, string(raw), taskID)
		if err != nil {
			return fmt.Errorf("%s: %w", table, err)
		}
		if count, err := result.RowsAffected(); err != nil || count != int64(len(entries)) {
			return fmt.Errorf("%s: unable to restore all rows under archived task", table)
		}
	}
	return nil
}
