package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type FindingScreenshot struct {
	ID        string    `json:"id"`
	BindingID string    `json:"binding_id,omitempty"`
	Caption   string    `json:"caption"`
	MIME      string    `json:"mime"`
	CreatedAt time.Time `json:"created_at"`
}

const screenshotCols = `id, COALESCE(binding_id::text,''), caption, mime, created_at`

func (d *DB) ListFindingScreenshots(ctx context.Context, findingID int64) ([]FindingScreenshot, error) {
	rows, err := d.QueryContext(ctx, `SELECT `+screenshotCols+` FROM finding_screenshots WHERE finding_id=$1 ORDER BY created_at,id`, findingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []FindingScreenshot{}
	for rows.Next() {
		var s FindingScreenshot
		if err := rows.Scan(&s.ID, &s.BindingID, &s.Caption, &s.MIME, &s.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (d *DB) CreateFindingScreenshot(ctx context.Context, findingID int64, bindingID *int64, caption, mime string, image []byte) (out FindingScreenshot, err error) {
	err = d.WithEvidenceTx(ctx, func(tx *sql.Tx) error {
		if err := LockFindingEvidenceTx(tx, findingID, nil); err != nil {
			return err
		}
		if bindingID != nil {
			var note, method, url string
			err := tx.QueryRow(`SELECT b.note, s.method, s.url FROM finding_traffic_bindings b JOIN traffic_evidence_snapshots s ON s.id=b.snapshot_id WHERE b.id=$1 AND b.finding_id=$2`, *bindingID, findingID).Scan(&note, &method, &url)
			if errors.Is(err, sql.ErrNoRows) {
				return ErrEvidenceNotFound
			}
			if err != nil {
				return err
			}
			if caption == "" {
				caption = note
				if caption == "" {
					caption = method + " " + url
				}
			}
		}
		return tx.QueryRow(`INSERT INTO finding_screenshots(id,finding_id,binding_id,caption,mime,image) VALUES($1,$2,$3,$4,$5,$6) RETURNING `+screenshotCols,
			uuid.NewString(), findingID, bindingID, caption, mime, image).Scan(&out.ID, &out.BindingID, &out.Caption, &out.MIME, &out.CreatedAt)
	})
	return
}

func (d *DB) ReadFindingScreenshot(ctx context.Context, findingID int64, id string) (mime string, data []byte, err error) {
	err = d.QueryRowContext(ctx, `SELECT mime,image FROM finding_screenshots WHERE finding_id=$1 AND id=$2`, findingID, id).Scan(&mime, &data)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrEvidenceNotFound
	}
	return
}

func (d *DB) ChangeFindingScreenshot(ctx context.Context, findingID int64, id string, caption *string) error {
	return d.WithEvidenceTx(ctx, func(tx *sql.Tx) error {
		if err := LockFindingEvidenceTx(tx, findingID, nil); err != nil {
			return err
		}
		var res sql.Result
		var err error
		if caption == nil {
			res, err = tx.Exec(`DELETE FROM finding_screenshots WHERE finding_id=$1 AND id=$2`, findingID, id)
		} else {
			res, err = tx.Exec(`UPDATE finding_screenshots SET caption=$3 WHERE finding_id=$1 AND id=$2`, findingID, id, *caption)
		}
		if err != nil {
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if n == 0 {
			return ErrEvidenceNotFound
		}
		return nil
	})
}

// The optional table is absent in earlier archives. Verify ownership before
// restoring so an archive cannot attach a screenshot to another task's finding.
func restoreFindingScreenshotsTx(tx *sql.Tx, snapshot *TaskArchiveSnapshot) error {
	rows, err := decodeArchiveRows(snapshot.Tables["finding_screenshots"])
	if err != nil {
		return err
	}
	for _, row := range rows {
		fid, ok := jsonInt64(row["finding_id"])
		if !ok || fid <= 0 {
			return errors.New("invalid archived screenshot finding id")
		}
		var owned bool
		if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM findings WHERE id=$1 AND task_id=$2)`, fid, snapshot.TaskID).Scan(&owned); err != nil {
			return err
		}
		if !owned {
			return errors.New("archived screenshot belongs to another task")
		}
		if row["binding_id"] != nil {
			bid, ok := jsonInt64(row["binding_id"])
			if !ok || bid <= 0 {
				return errors.New("invalid archived screenshot binding")
			}
			if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM finding_traffic_bindings WHERE id=$1 AND finding_id=$2)`, bid, fid).Scan(&owned); err != nil {
				return err
			}
			if !owned {
				return errors.New("archived screenshot binding belongs to another finding")
			}
		}
	}
	return insertArchiveRows(tx, "finding_screenshots", snapshot.Tables["finding_screenshots"])
}

type FindingScreenshotData struct {
	FindingScreenshot
	Image []byte `json:"-"`
}

// LoadFindingScreenshots reads only findings selected for the report.
func (d *DB) LoadFindingScreenshots(ctx context.Context, findings []*DBFinding, includeImages bool) error {
	if len(findings) == 0 {
		return nil
	}
	return d.WithEvidenceTx(ctx, func(tx *sql.Tx) error {
		for _, f := range findings {
			query := `SELECT ` + screenshotCols + `, NULL::bytea FROM finding_screenshots WHERE finding_id=$1 ORDER BY created_at,id`
			if includeImages {
				query = `SELECT ` + screenshotCols + `, image FROM finding_screenshots WHERE finding_id=$1 ORDER BY created_at,id`
			}
			rows, err := tx.QueryContext(ctx, query, f.ID)
			if err != nil {
				return err
			}
			var shots []FindingScreenshotData
			for rows.Next() {
				var shot FindingScreenshotData
				if err := rows.Scan(&shot.ID, &shot.BindingID, &shot.Caption, &shot.MIME, &shot.CreatedAt, &shot.Image); err != nil {
					rows.Close()
					return err
				}
				shots = append(shots, shot)
			}
			err = rows.Err()
			rows.Close()
			if err != nil {
				return err
			}
			f.Screenshots = shots
		}
		return nil
	})
}
