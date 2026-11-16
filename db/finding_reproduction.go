package db

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

type FindingReproductionStep struct {
	Title       string `json:"title"`
	Observation string `json:"observation"`
	BindingID   string `json:"binding_id,omitempty"`
}

type reproductionQueryer interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func listFindingReproduction(ctx context.Context, q reproductionQueryer, findingID int64) ([]FindingReproductionStep, error) {
	rows, err := q.QueryContext(ctx, `SELECT title,observation,COALESCE(binding_id::text,'') FROM finding_reproduction_steps WHERE finding_id=$1 ORDER BY position`, findingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []FindingReproductionStep{}
	for rows.Next() {
		var step FindingReproductionStep
		if err := rows.Scan(&step.Title, &step.Observation, &step.BindingID); err != nil {
			return nil, err
		}
		out = append(out, step)
	}
	return out, rows.Err()
}

func (d *DB) ListFindingReproduction(ctx context.Context, findingID int64) ([]FindingReproductionStep, error) {
	return listFindingReproduction(ctx, d.DB, findingID)
}

func (d *DB) ReplaceFindingReproduction(ctx context.Context, findingID int64, steps []FindingReproductionStep) error {
	return d.WithEvidenceTx(ctx, func(tx *sql.Tx) error {
		if err := LockFindingEvidenceTx(tx, findingID, nil); err != nil {
			return err
		}
		if _, err := tx.Exec(`DELETE FROM finding_reproduction_steps WHERE finding_id=$1`, findingID); err != nil {
			return err
		}
		for i, step := range steps {
			if strings.TrimSpace(step.Title) == "" {
				return errors.New("reproduction step title is required")
			}
			var binding *int64
			if step.BindingID != "" {
				id, err := strconv.ParseInt(step.BindingID, 10, 64)
				if err != nil || id <= 0 {
					return ErrEvidenceNotFound
				}
				var owned bool
				if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM finding_traffic_bindings WHERE id=$1 AND finding_id=$2)`, id, findingID).Scan(&owned); err != nil {
					return err
				}
				if !owned {
					return ErrEvidenceNotFound
				}
				binding = &id
			}
			if _, err := tx.Exec(`INSERT INTO finding_reproduction_steps(finding_id,position,title,observation,binding_id) VALUES($1,$2,$3,$4,$5)`, findingID, i, step.Title, step.Observation, binding); err != nil {
				return err
			}
		}
		return nil
	})
}

func restoreFindingReproductionTx(tx *sql.Tx, snapshot *TaskArchiveSnapshot) error {
	raw := snapshot.Tables["finding_reproduction_steps"]
	rows, err := decodeArchiveRows(raw)
	if err != nil {
		return err
	}
	for _, row := range rows {
		fid, ok := jsonInt64(row["finding_id"])
		if !ok || fid <= 0 {
			return errors.New("invalid archived reproduction finding id")
		}
		var owned bool
		if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM findings WHERE id=$1 AND task_id=$2)`, fid, snapshot.TaskID).Scan(&owned); err != nil {
			return err
		}
		if !owned {
			return errors.New("archived reproduction step belongs to another task")
		}
		if row["binding_id"] != nil {
			bid, ok := jsonInt64(row["binding_id"])
			if !ok || bid <= 0 {
				return errors.New("invalid archived reproduction binding")
			}
			if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM finding_traffic_bindings WHERE id=$1 AND finding_id=$2)`, bid, fid).Scan(&owned); err != nil {
				return err
			}
			if !owned {
				return errors.New("archived reproduction binding belongs to another finding")
			}
		}
	}
	return insertArchiveRows(tx, "finding_reproduction_steps", raw)
}
