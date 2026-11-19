package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrRemediationNotFound = errors.New("finding remediation not found")
	ErrRemediationReadOnly = errors.New("finding remediation is read-only")
	ErrRemediationConflict = errors.New("finding remediation has changed or does not allow this action")
	ErrRemediationInvalid  = errors.New("invalid finding remediation input")
)

const (
	remediationReview     = "pending_review"
	remediationFix        = "pending_fix"
	remediationAcceptance = "pending_acceptance"
	remediationAccepted   = "accepted"
	remediationEventPage  = 50
)

type RemediationFilter struct {
	Stage, Severity, Query string
	Page, Limit            int
}

type RemediationTransition struct {
	Action, Note, Reference, Plan, Owner, DueDate, Actor string
	ExpectedVersion, ExpectedEvidenceVersion             *int64
}

type FindingRemediation struct {
	FindingID                string `json:"finding_id"`
	Title                    string `json:"title"`
	Severity                 string `json:"severity"`
	FindingStatus            string `json:"finding_status"`
	Stage                    string `json:"stage"`
	Version                  int64  `json:"version"`
	EvidenceVersion          int64  `json:"evidence_version"`
	Owner                    string `json:"owner,omitempty"`
	DueDate                  string `json:"due_date,omitempty"`
	Plan                     string `json:"plan,omitempty"`
	Reference                string `json:"reference,omitempty"`
	NeedsReview              bool   `json:"needs_review"`
	ReadOnly                 bool   `json:"read_only"`
	SubmittedEvidenceVersion *int64 `json:"submitted_evidence_version,omitempty"`
	AcceptedEvidenceVersion  *int64 `json:"accepted_evidence_version,omitempty"`
}

type FindingRemediationEvent struct {
	ID              int64     `json:"id"`
	Version         int64     `json:"version"`
	Action          string    `json:"action"`
	FromStage       string    `json:"from_stage"`
	ToStage         string    `json:"to_stage"`
	EvidenceVersion int64     `json:"evidence_version"`
	Actor           string    `json:"actor"`
	Note            string    `json:"note"`
	Reference       string    `json:"reference,omitempty"`
	Plan            string    `json:"plan,omitempty"`
	Owner           string    `json:"owner,omitempty"`
	DueDate         string    `json:"due_date,omitempty"`
	Time            time.Time `json:"time"`
}

type FindingRemediationDetail struct {
	Item       FindingRemediation        `json:"item"`
	Events     []FindingRemediationEvent `json:"events"`
	NextBefore int64                     `json:"next_before,omitempty"`
}

const remediationTitle = `COALESCE(NULLIF(f.name,''),NULLIF(f.vulnclass,''),'Finding #'||f.id::text)`
const remediationStage = `COALESCE(r.stage,'pending_review')`
const remediationExcluded = `f.status IN ('false_positive','ignored','duplicate','risk_accepted')`
const remediationReadOnly = `(` + remediationExcluded + ` OR t.deleted_at IS NOT NULL OR t.archived_at IS NOT NULL OR COALESCE(a.state<>'archive_failed',false))`
const remediationFrom = ` FROM findings f
 LEFT JOIN finding_remediations r ON r.finding_id=f.id
 LEFT JOIN tasks t ON t.id=f.task_id
 LEFT JOIN task_archives a ON a.task_id=f.task_id`
const remediationColumns = `f.id::text,` + remediationTitle + `,f.severity,f.status,` + remediationStage +
	`,COALESCE(r.version,0),f.evidence_version,COALESCE(r.owner,''),COALESCE(to_char(r.due_date,'YYYY-MM-DD'),''),
 COALESCE(r.plan,''),COALESCE(r.reference,''),r.submitted_evidence_version,r.accepted_evidence_version,` + remediationReadOnly

func knownRemediationStage(stage string) bool {
	switch stage {
	case remediationReview, remediationFix, remediationAcceptance, remediationAccepted:
		return true
	default:
		return false
	}
}

func (item *FindingRemediation) setNeedsReview() {
	item.NeedsReview = (item.Stage == remediationAcceptance && !sameEvidence(item.SubmittedEvidenceVersion, item.EvidenceVersion)) ||
		(item.Stage == remediationAccepted && !sameEvidence(item.AcceptedEvidenceVersion, item.EvidenceVersion))
}

func sameEvidence(version *int64, current int64) bool { return version != nil && *version == current }

func scanFindingRemediation(row interface{ Scan(...any) error }) (item FindingRemediation, err error) {
	err = row.Scan(&item.FindingID, &item.Title, &item.Severity, &item.FindingStatus, &item.Stage,
		&item.Version, &item.EvidenceVersion, &item.Owner, &item.DueDate, &item.Plan, &item.Reference,
		&item.SubmittedEvidenceVersion, &item.AcceptedEvidenceVersion, &item.ReadOnly)
	item.setNeedsReview()
	return item, err
}

func (d *DB) ListFindingRemediations(ctx context.Context, filter RemediationFilter) ([]FindingRemediation, int64, error) {
	filter.Query = strings.TrimSpace(filter.Query)
	if (filter.Stage != "" && !knownRemediationStage(filter.Stage)) || !validRemediationText(filter.Query, 200) {
		return nil, 0, ErrRemediationInvalid
	}
	switch filter.Severity {
	case "", "critical", "high", "medium", "low", "info":
	default:
		return nil, 0, ErrRemediationInvalid
	}
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 25
	}
	if filter.Page < 1 || filter.Page > 1_000_000 || filter.Limit < 1 || filter.Limit > 100 {
		return nil, 0, ErrRemediationInvalid
	}
	const where = ` WHERE t.deleted_at IS NULL AND t.archived_at IS NULL AND NOT(` + remediationExcluded + `)
 AND ($1::text='' OR ` + remediationStage + `=$1)
 AND ($2::text='' OR f.severity=$2)
 AND ($3::text='' OR strpos(lower(` + remediationTitle + `),lower($3))>0
     OR strpos(lower(COALESCE(r.owner,'')),lower($3))>0 OR f.id::text=$3)`
	tx, err := d.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	args := []any{filter.Stage, filter.Severity, filter.Query}
	var total int64
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*)`+remediationFrom+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := tx.QueryContext(ctx, `SELECT `+remediationColumns+remediationFrom+where+
		` ORDER BY f.id DESC LIMIT $4 OFFSET $5`, append(args, filter.Limit, int64(filter.Page-1)*int64(filter.Limit))...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]FindingRemediation, 0)
	for rows.Next() {
		item, err := scanFindingRemediation(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, tx.Commit()
}

func parseRemediationFindingID(id string) (int64, error) {
	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil || n <= 0 || strconv.FormatInt(n, 10) != id {
		return 0, ErrRemediationInvalid
	}
	return n, nil
}

func readFindingRemediation(ctx context.Context, tx *sql.Tx, id int64) (FindingRemediation, error) {
	item, err := scanFindingRemediation(tx.QueryRowContext(ctx, `SELECT `+remediationColumns+remediationFrom+` WHERE f.id=$1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return item, ErrRemediationNotFound
	}
	return item, err
}

func readRemediationEvents(ctx context.Context, tx *sql.Tx, id, before int64) ([]FindingRemediationEvent, int64, error) {
	rows, err := tx.QueryContext(ctx, `SELECT version,action,from_stage,to_stage,evidence_version,
 actor,note,reference,plan,owner,COALESCE(to_char(due_date,'YYYY-MM-DD'),''),created_at
 FROM finding_remediation_events WHERE finding_id=$1 AND ($2::bigint=0 OR version<$2)
 ORDER BY version DESC LIMIT $3`, id, before, remediationEventPage+1)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]FindingRemediationEvent, 0)
	for rows.Next() {
		var item FindingRemediationEvent
		if err := rows.Scan(&item.Version, &item.Action, &item.FromStage, &item.ToStage, &item.EvidenceVersion,
			&item.Actor, &item.Note, &item.Reference, &item.Plan, &item.Owner, &item.DueDate, &item.Time); err != nil {
			return nil, 0, err
		}
		item.ID = item.Version
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	if len(items) > remediationEventPage {
		return items[:remediationEventPage], items[remediationEventPage-1].Version, nil
	}
	return items, 0, nil
}

func (d *DB) GetFindingRemediation(ctx context.Context, findingID string, before int64) (FindingRemediationDetail, error) {
	id, err := parseRemediationFindingID(findingID)
	if err != nil || before < 0 {
		return FindingRemediationDetail{}, ErrRemediationInvalid
	}
	tx, err := d.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return FindingRemediationDetail{}, err
	}
	defer tx.Rollback()
	item, err := readFindingRemediation(ctx, tx, id)
	if err != nil {
		return FindingRemediationDetail{}, err
	}
	events, next, err := readRemediationEvents(ctx, tx, id, before)
	if err != nil {
		return FindingRemediationDetail{}, err
	}
	return FindingRemediationDetail{Item: item, Events: events, NextBefore: next}, tx.Commit()
}

func validRemediationText(value string, limit int) bool {
	return utf8.ValidString(value) && !strings.ContainsRune(value, 0) && utf8.RuneCountInString(value) <= limit
}

func (change RemediationTransition) validated() (RemediationTransition, error) {
	change.Note = strings.TrimSpace(change.Note)
	change.Reference = strings.TrimSpace(change.Reference)
	change.Plan = strings.TrimSpace(change.Plan)
	change.Owner = strings.TrimSpace(change.Owner)
	change.DueDate = strings.TrimSpace(change.DueDate)
	change.Actor = strings.TrimSpace(change.Actor)
	if change.ExpectedVersion == nil || *change.ExpectedVersion < 0 || change.Actor == "" ||
		!validRemediationText(change.Actor, 120) || !validRemediationText(change.Note, 2000) ||
		!validRemediationText(change.Reference, 200) || !validRemediationText(change.Plan, 3000) ||
		!validRemediationText(change.Owner, 120) {
		return change, ErrRemediationInvalid
	}
	if change.Action != "update_plan" && (change.Note == "" || change.Plan != "" || change.Owner != "" || change.DueDate != "") {
		return change, ErrRemediationInvalid
	}
	if change.Action != "submit" && change.Reference != "" {
		return change, ErrRemediationInvalid
	}
	if change.Action == "submit" || change.Action == "accept" {
		if change.ExpectedEvidenceVersion == nil || *change.ExpectedEvidenceVersion < 0 {
			return change, ErrRemediationInvalid
		}
	} else if change.ExpectedEvidenceVersion != nil {
		return change, ErrRemediationInvalid
	}
	switch change.Action {
	case "review", "accept", "return", "reopen":
	case "update_plan":
		if change.Plan == "" {
			return change, ErrRemediationInvalid
		}
		if change.DueDate != "" {
			date, err := time.Parse(time.DateOnly, change.DueDate)
			if err != nil || date.Year() == 0 || date.Format(time.DateOnly) != change.DueDate {
				return change, ErrRemediationInvalid
			}
		}
	case "submit":
		if change.Reference == "" {
			return change, ErrRemediationInvalid
		}
	default:
		return change, ErrRemediationInvalid
	}
	return change, nil
}

func applyRemediationChange(item FindingRemediation, change RemediationTransition) (FindingRemediation, error) {
	if item.ReadOnly {
		return item, ErrRemediationReadOnly
	}
	if item.Version != *change.ExpectedVersion ||
		(change.ExpectedEvidenceVersion != nil && item.EvidenceVersion != *change.ExpectedEvidenceVersion) {
		return item, ErrRemediationConflict
	}
	next := item
	switch change.Action {
	case "review":
		if item.Stage != remediationReview {
			return item, ErrRemediationConflict
		}
		next.Stage = remediationFix
	case "update_plan":
		if item.Stage != remediationFix {
			return item, ErrRemediationConflict
		}
		next.Plan, next.Owner, next.DueDate = change.Plan, change.Owner, change.DueDate
	case "submit":
		if item.Stage != remediationFix {
			return item, ErrRemediationConflict
		}
		next.Stage = remediationAcceptance
		next.Reference = change.Reference
		evidence := item.EvidenceVersion
		next.SubmittedEvidenceVersion = &evidence
	case "accept":
		if item.Stage != remediationAcceptance || !sameEvidence(item.SubmittedEvidenceVersion, item.EvidenceVersion) {
			return item, fmt.Errorf("%w: acceptance requires current submitted evidence", ErrRemediationConflict)
		}
		next.Stage = remediationAccepted
		evidence := item.EvidenceVersion
		next.AcceptedEvidenceVersion = &evidence
	case "return", "reopen":
		if (change.Action == "return" && item.Stage != remediationAcceptance) ||
			(change.Action == "reopen" && item.Stage != remediationAccepted) {
			return item, ErrRemediationConflict
		}
		next.Stage = remediationFix
		next.Reference, next.SubmittedEvidenceVersion, next.AcceptedEvidenceVersion = "", nil, nil
	default:
		return item, ErrRemediationInvalid
	}
	next.Version++
	next.setNeedsReview()
	return next, nil
}

func (d *DB) TransitionFindingRemediation(ctx context.Context, findingID string, change RemediationTransition) (FindingRemediationDetail, error) {
	id, err := parseRemediationFindingID(findingID)
	if err != nil {
		return FindingRemediationDetail{}, err
	}
	change, err = change.validated()
	if err != nil {
		return FindingRemediationDetail{}, err
	}
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return FindingRemediationDetail{}, err
	}
	defer tx.Rollback()
	if err := coordinateWithSchemaMigration(tx); err != nil {
		return FindingRemediationDetail{}, err
	}
	// The task lock orders this work with the archive queue; the finding lock
	// orders it with evidence and disposition edits and serializes first writes.
	if err := LockFindingEvidenceTx(tx, id, nil); err != nil {
		switch {
		case errors.Is(err, ErrFindingNotFound):
			err = ErrRemediationNotFound
		case errors.Is(err, ErrTaskArchiveState), errors.Is(err, sql.ErrNoRows):
			err = ErrRemediationReadOnly
		}
		return FindingRemediationDetail{}, err
	}
	current, err := readFindingRemediation(ctx, tx, id)
	if err != nil {
		return FindingRemediationDetail{}, err
	}
	next, err := applyRemediationChange(current, change)
	if err != nil {
		return FindingRemediationDetail{}, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO finding_remediations
 (finding_id,stage,version,owner,due_date,plan,reference,submitted_evidence_version,accepted_evidence_version)
 VALUES($1,$2,$3,$4,NULLIF($5,'')::date,$6,$7,$8,$9)
 ON CONFLICT (finding_id) DO UPDATE SET stage=EXCLUDED.stage,version=EXCLUDED.version,owner=EXCLUDED.owner,
 due_date=EXCLUDED.due_date,plan=EXCLUDED.plan,reference=EXCLUDED.reference,
 submitted_evidence_version=EXCLUDED.submitted_evidence_version,accepted_evidence_version=EXCLUDED.accepted_evidence_version`,
		id, next.Stage, next.Version, next.Owner, next.DueDate, next.Plan, next.Reference,
		next.SubmittedEvidenceVersion, next.AcceptedEvidenceVersion); err != nil {
		return FindingRemediationDetail{}, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO finding_remediation_events
 (finding_id,version,action,from_stage,to_stage,evidence_version,actor,note,reference,plan,owner,due_date)
 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NULLIF($12,'')::date)`,
		id, next.Version, change.Action, current.Stage, next.Stage, next.EvidenceVersion,
		change.Actor, change.Note, change.Reference, change.Plan, change.Owner, change.DueDate); err != nil {
		return FindingRemediationDetail{}, err
	}
	events, before, err := readRemediationEvents(ctx, tx, id, 0)
	if err != nil {
		return FindingRemediationDetail{}, err
	}
	return FindingRemediationDetail{Item: next, Events: events, NextBefore: before}, tx.Commit()
}
