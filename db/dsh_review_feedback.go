package db

import (
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"
    "strings"
)

const DSHReviewFeedbackSchema = "artex.dsh.review-feedback.v1"

type DSHReviewEntry struct {
    FindingID int64 `json:"finding_id"`
    TaskID *int64 `json:"task_id"`
    EvidenceVersion *int64 `json:"evidence_version"`
    Status string `json:"status"`
    Note string `json:"note"`
    UpdatedAt string `json:"updated_at"`
}

type DSHReviewFeedback struct {
    Schema string `json:"schema"`
    Source string `json:"source"`
    Engagement string `json:"engagement"`
    ExportedAt string `json:"exported_at"`
    Reviews []DSHReviewEntry `json:"reviews"`
}

type DSHReviewImportResult struct {
    Imported int `json:"imported"`
    Missing int `json:"missing"`
    TaskMismatch int `json:"task_mismatch"`
}

type DSHReviewRow struct {
    FindingID int64 `json:"finding_id"`
    Engagement string `json:"engagement"`
    SourceTaskID *int64 `json:"source_task_id"`
    EvidenceVersion *int64 `json:"evidence_version"`
    Status string `json:"status"`
    Note string `json:"note"`
    SourceUpdatedAt string `json:"source_updated_at"`
    ARTEXStatus string `json:"artex_status"`
    CurrentEvidenceVersion int64 `json:"current_evidence_version"`
    EvidenceChanged bool `json:"evidence_changed"`
}

func ParseDSHReviewFeedback(data []byte) (DSHReviewFeedback, error) {
    var input DSHReviewFeedback
    if err := json.Unmarshal(data, &input); err != nil { return input, fmt.Errorf("invalid JSON: %w", err) }
    if input.Schema != DSHReviewFeedbackSchema || input.Source != "dsh-purge" {
        return input, errors.New("expected dsh-purge review feedback v1")
    }
    input.Engagement = strings.TrimSpace(input.Engagement)
    if len(input.Engagement) == 0 || len(input.Engagement) > 128 {
        return input, errors.New("invalid DSH engagement")
    }
    if len(input.Reviews) == 0 || len(input.Reviews) > 100 {
        return input, errors.New("expected 1..100 saved review records")
    }
    seen := make(map[int64]bool, len(input.Reviews))
    for _, item := range input.Reviews {
        if item.FindingID <= 0 || seen[item.FindingID] {
            return input, errors.New("finding IDs must be unique positive integers")
        }
        seen[item.FindingID] = true
        if item.TaskID != nil && *item.TaskID <= 0 { return input, errors.New("invalid source task ID") }
        if item.EvidenceVersion != nil && *item.EvidenceVersion < 0 { return input, errors.New("invalid evidence version") }
        switch item.Status {
        case "pending", "in_progress", "needs_evidence", "reviewed":
        default: return input, errors.New("invalid human review status")
        }
        if len(item.Note) > 1000 || len(item.UpdatedAt) > 80 {
            return input, errors.New("review note or timestamp too long")
        }
    }
    return input, nil
}

// Import records human review progress only: no changes to findings.status,
// evidence, task states, scores or traffic. A mismatched task association is
// reported, not rewritten. The operation is atomic for valid source records.
func (d *DB) ImportDSHReviewFeedback(input DSHReviewFeedback) (DSHReviewImportResult, error) {
    var result DSHReviewImportResult
    tx, err := d.Begin()
    if err != nil { return result, err }
    defer tx.Rollback()
    for _, item := range input.Reviews {
        var taskID sql.NullInt64
        err := tx.QueryRow("SELECT task_id FROM findings WHERE id=$1", item.FindingID).Scan(&taskID)
        if errors.Is(err, sql.ErrNoRows) { result.Missing++; continue }
        if err != nil { return result, err }
        if item.TaskID != nil && taskID.Valid && taskID.Int64 != *item.TaskID {
            result.TaskMismatch++
            continue
        }
        _, err = tx.Exec(`INSERT INTO finding_dsh_reviews(
            finding_id, dsh_engagement, source_task_id, inspected_evidence_version,
            review_status, note, source_updated_at)
            VALUES ($1,$2,$3,$4,$5,$6,$7)
            ON CONFLICT (finding_id,dsh_engagement) DO UPDATE SET
                source_task_id=excluded.source_task_id,
                inspected_evidence_version=excluded.inspected_evidence_version,
                review_status=excluded.review_status, note=excluded.note,
                source_updated_at=excluded.source_updated_at, imported_at=now()
            WHERE finding_dsh_reviews.source_updated_at <= excluded.source_updated_at`,
            item.FindingID, input.Engagement, item.TaskID, item.EvidenceVersion,
            item.Status, item.Note, item.UpdatedAt)
        if err != nil { return result, err }
        result.Imported++
    }
    if err := tx.Commit(); err != nil { return result, err }
    return result, nil
}

func (d *DB) ListDSHReviewFeedback(limit int) ([]DSHReviewRow, error) {
    if limit <= 0 || limit > 100 { limit = 50 }
    rows, err := d.Query(`SELECT r.finding_id, r.dsh_engagement, r.source_task_id,
        r.inspected_evidence_version, r.review_status, r.note, r.source_updated_at,
        f.status, f.evidence_version
        FROM finding_dsh_reviews r
        JOIN findings f ON f.id=r.finding_id
        ORDER BY r.imported_at DESC, r.finding_id ASC LIMIT $1`, limit)
    if err != nil { return nil, err }
    defer rows.Close()
    out := []DSHReviewRow{}
    for rows.Next() {
        var row DSHReviewRow
        if err := rows.Scan(&row.FindingID, &row.Engagement, &row.SourceTaskID,
            &row.EvidenceVersion, &row.Status, &row.Note, &row.SourceUpdatedAt,
            &row.ARTEXStatus, &row.CurrentEvidenceVersion); err != nil { return nil, err }
        row.EvidenceChanged = row.EvidenceVersion != nil && *row.EvidenceVersion != row.CurrentEvidenceVersion
        out = append(out, row)
    }
    return out, rows.Err()
}
