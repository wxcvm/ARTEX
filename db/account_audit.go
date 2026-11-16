package db

import (
	"context"
	"time"
)

type AccountAuditEvent struct {
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Action string `json:"action"`
	Result string `json:"result"`
}

// RecordAccountAudit stores only server-selected action and result labels.
func (d *DB) RecordAccountAudit(ctx context.Context, action, result string) error {
	_, err := d.ExecContext(ctx, "INSERT INTO account_audit(action,result) VALUES ($1,$2)", action, result)
	return err
}

// ListAccountAudit returns newest-first rows, using an exclusive stable ID cursor.
func (d *DB) ListAccountAudit(ctx context.Context, action, result string, before int64, limit int) ([]AccountAuditEvent, error) {
	rows, err := d.QueryContext(ctx, `SELECT id,created_at,action,result FROM account_audit
		WHERE ($1='' OR action=$1) AND ($2='' OR result=$2) AND ($3=0 OR id<$3)
		ORDER BY id DESC LIMIT $4`, action, result, before, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	out := []AccountAuditEvent{}
	for rows.Next() {
		var e AccountAuditEvent
		if err := rows.Scan(&e.ID,&e.CreatedAt,&e.Action,&e.Result); err != nil { return nil, err }
		out = append(out,e)
	}
	return out, rows.Err()
}
