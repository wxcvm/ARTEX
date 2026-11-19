package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Autumn-27/artex/db"
)

const maxRemediationEventBody = 64 << 10

// The remediation surface has a deliberately narrow dependency: persist and
// read annotations on existing findings. It has no task, model or tool hooks.
type remediationStore interface {
	ListFindingRemediations(context.Context, db.RemediationFilter) ([]db.FindingRemediation, int64, error)
	GetFindingRemediation(context.Context, string, int64) (db.FindingRemediationDetail, error)
	TransitionFindingRemediation(context.Context, string, db.RemediationTransition) (db.FindingRemediationDetail, error)
}

type remediationHTTP struct {
	store  func(http.ResponseWriter) remediationStore
	jwtKey []byte
}

func (s *Server) registerRemediation(mux *http.ServeMux) {
	remediationHTTP{
		store: func(w http.ResponseWriter) remediationStore {
			pg := s.pg(w)
			if pg == nil {
				return nil
			}
			return pg
		},
		jwtKey: s.jwtKey,
	}.register(mux)
}

// The caller mounts these routes behind requireAuth, like the finding routes.
// Neither the privacy nor pairing bridge has an exemption for this surface.
func (h remediationHTTP) register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/remediation", h.list)
	mux.HandleFunc("GET /api/exploration/findings/{id}/remediation", h.detail)
	mux.HandleFunc("POST /api/exploration/findings/{id}/remediation/events", h.event)
}

type remediationPage struct {
	Items    []db.FindingRemediation `json:"items"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
}

func (h remediationHTTP) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	q := r.URL.Query()
	filter := db.RemediationFilter{
		Stage:    strings.TrimSpace(q.Get("stage")),
		Severity: strings.TrimSpace(q.Get("severity")),
		Query:    strings.TrimSpace(q.Get("q")),
		Page:     1,
		Limit:    25,
	}
	switch filter.Stage {
	case "", "pending_review", "pending_fix", "pending_acceptance", "accepted":
	default:
		writeErr(w, http.StatusBadRequest, "invalid remediation stage")
		return
	}
	if filter.Severity != "" && filter.Severity != "info" && !db.ValidSeverity(filter.Severity) {
		writeErr(w, http.StatusBadRequest, "invalid finding severity")
		return
	}
	if !validRemediationText(filter.Query, 200, false) {
		writeErr(w, http.StatusBadRequest, "query must be at most 200 characters of plain text")
		return
	}
	for _, param := range []struct {
		name string
		max  int
		dst  *int
	}{{"page", 1_000_000, &filter.Page}, {"limit", 100, &filter.Limit}} {
		if raw := q.Get(param.name); raw != "" {
			n, err := strconv.Atoi(raw)
			if err != nil || n < 1 || n > param.max {
				writeErr(w, http.StatusBadRequest, "invalid "+param.name)
				return
			}
			*param.dst = n
		}
	}
	store := h.store(w)
	if store == nil {
		return
	}
	items, total, err := store.ListFindingRemediations(r.Context(), filter)
	if err != nil {
		writeRemediationError(w, err)
		return
	}
	if items == nil {
		items = []db.FindingRemediation{}
	}
	writeJSON(w, http.StatusOK, remediationPage{Items: items, Total: total, Page: filter.Page, PageSize: filter.Limit})
}

func remediationFindingID(w http.ResponseWriter, r *http.Request) (string, bool) {
	// An inherited task view can refer to a historical snapshot. This workbench
	// only tracks the canonical finding, not the task's inherited view.
	if r.URL.Query().Get("context_task") != "" {
		writeErr(w, http.StatusBadRequest, "remediation only supports canonical findings")
		return "", false
	}
	id, ok := pathInt(r, "id")
	if !ok || id <= 0 {
		writeErr(w, http.StatusBadRequest, "invalid finding id")
		return "", false
	}
	return strconv.FormatInt(id, 10), true
}

func (h remediationHTTP) detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	id, ok := remediationFindingID(w, r)
	if !ok {
		return
	}
	var before int64
	if raw := r.URL.Query().Get("before"); raw != "" {
		var err error
		before, err = strconv.ParseInt(raw, 10, 64)
		if err != nil || before <= 0 {
			writeErr(w, http.StatusBadRequest, "invalid history cursor")
			return
		}
	}
	store := h.store(w)
	if store == nil {
		return
	}
	detail, err := store.GetFindingRemediation(r.Context(), id, before)
	if err != nil {
		writeRemediationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

func (h remediationHTTP) event(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	id, ok := remediationFindingID(w, r)
	if !ok {
		return
	}
	token := parseVerifiedJWT(extractToken(r), h.jwtKey)
	if token == nil {
		writeErr(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	actor, err := token.Claims.GetSubject()
	if err != nil || strings.TrimSpace(actor) == "" || !validRemediationText(actor, 120, false) {
		writeErr(w, http.StatusUnauthorized, "missing signed account identity")
		return
	}
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" {
		writeErr(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxRemediationEventBody)
	change, err := readRemediationEvent(r.Body)
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeErr(w, http.StatusRequestEntityTooLarge, "remediation event is too large")
			return
		}
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	change.Actor = actor // The actor cannot be supplied in the JSON body.
	store := h.store(w)
	if store == nil {
		return
	}
	detail, err := store.TransitionFindingRemediation(r.Context(), id, change)
	if err != nil {
		writeRemediationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

type remediationEventRequest struct {
	Action                  string  `json:"action"`
	ExpectedVersion         *int64  `json:"expected_version"`
	ExpectedEvidenceVersion *int64  `json:"expected_evidence_version"`
	Note                    *string `json:"note"`
	Reference               *string `json:"reference"`
	Plan                    *string `json:"plan"`
	Owner                   *string `json:"owner"`
	DueDate                 *string `json:"due_date"`
}

func readRemediationEvent(body io.Reader) (db.RemediationTransition, error) {
	var input remediationEventRequest
	var zero db.RemediationTransition
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			return zero, err
		}
		return zero, errors.New("invalid remediation event JSON or unsupported field")
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			return zero, err
		}
		return zero, errors.New("expected one remediation event")
	}
	if input.ExpectedVersion == nil || *input.ExpectedVersion < 0 {
		return zero, errors.New("expected_version is required and must not be negative")
	}
	switch input.Action {
	case "review", "update_plan", "submit", "accept", "return", "reopen":
	default:
		return zero, errors.New("invalid remediation action")
	}
	fields := []struct {
		name      string
		value     *string
		max       int
		multiline bool
		dst       *string
	}{
		{"note", input.Note, 2000, true, &zero.Note},
		{"reference", input.Reference, 200, false, &zero.Reference},
		{"plan", input.Plan, 3000, true, &zero.Plan},
		{"owner", input.Owner, 120, false, &zero.Owner},
		{"due_date", input.DueDate, 10, false, &zero.DueDate},
	}
	for _, field := range fields {
		if field.value == nil {
			continue
		}
		*field.dst = strings.TrimSpace(*field.value)
		if !validRemediationText(*field.dst, field.max, field.multiline) {
			return db.RemediationTransition{}, fmt.Errorf("%s must be at most %d characters of plain text", field.name, field.max)
		}
	}
	if input.Action != "update_plan" && (input.Plan != nil || input.Owner != nil || input.DueDate != nil) {
		return db.RemediationTransition{}, errors.New("plan, owner and due_date are only for update_plan")
	}
	if input.Action != "submit" && input.Reference != nil {
		return db.RemediationTransition{}, errors.New("reference is only for submit")
	}
	if input.Action == "submit" || input.Action == "accept" {
		if input.ExpectedEvidenceVersion == nil || *input.ExpectedEvidenceVersion < 0 {
			return db.RemediationTransition{}, errors.New("expected_evidence_version is required and must not be negative")
		}
	} else if input.ExpectedEvidenceVersion != nil {
		return db.RemediationTransition{}, errors.New("expected_evidence_version is only for submit or accept")
	}
	if input.Action != "update_plan" && zero.Note == "" {
		return db.RemediationTransition{}, errors.New("a note is required for this action")
	}
	if input.Action == "submit" && zero.Reference == "" {
		return db.RemediationTransition{}, errors.New("a verification reference is required for submit")
	}
	if input.Action == "update_plan" && zero.Plan == "" {
		return db.RemediationTransition{}, errors.New("a remediation plan is required")
	}
	if zero.DueDate != "" {
		d, err := time.Parse(time.DateOnly, zero.DueDate)
		if err != nil || d.Year() < 1 || d.Format(time.DateOnly) != zero.DueDate {
			return db.RemediationTransition{}, errors.New("due_date must be a valid YYYY-MM-DD date")
		}
	}
	zero.Action = input.Action
	zero.ExpectedVersion = input.ExpectedVersion
	zero.ExpectedEvidenceVersion = input.ExpectedEvidenceVersion
	return zero, nil
}

// Notes/plans may contain line breaks; references and names are single-line
// labels. UI renders them as plain text and the server never dereferences them.
func validRemediationText(text string, maxChars int, multiline bool) bool {
	if !utf8.ValidString(text) || utf8.RuneCountInString(text) > maxChars {
		return false
	}
	for _, r := range text {
		if unicode.IsControl(r) && !(multiline && (r == '\n' || r == '\r' || r == '\t')) {
			return false
		}
	}
	return true
}

func writeRemediationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, db.ErrRemediationNotFound):
		writeErr(w, http.StatusNotFound, "finding not available for remediation")
	case errors.Is(err, db.ErrRemediationReadOnly):
		writeErr(w, http.StatusConflict, "finding is read-only")
	case errors.Is(err, db.ErrRemediationConflict):
		writeErr(w, http.StatusConflict, "remediation state or evidence has changed; refresh and review before resubmitting")
	case errors.Is(err, db.ErrRemediationInvalid):
		writeErr(w, http.StatusBadRequest, "invalid remediation data")
	default:
		writeErr(w, http.StatusServiceUnavailable, "remediation data is temporarily unavailable")
	}
}
