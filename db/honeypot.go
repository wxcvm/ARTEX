package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// HoneypotAssessmentTTL bounds how long an explicit fingerprint suppresses
// repeat work. A short cache avoids treating a mutable service as permanently
// classified while still preventing multi-agent retries within a task.
const HoneypotAssessmentTTL = time.Hour

// HoneypotAssessment is a high-confidence, rule-based observation. It means
// "suspected honeypot", never a definitive claim about the target.
type HoneypotAssessment struct {
	Host       string    `json:"host"`
	Score      int       `json:"score"`
	Level      string    `json:"level"`
	Signals    []string  `json:"signals"`
	ObservedAt time.Time `json:"observed_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

func normalizeHoneypotHost(raw string) string {
	raw = strings.TrimSpace(raw)
	if u, err := url.Parse(raw); err == nil && u.Hostname() != "" {
		raw = u.Hostname()
	}
	return strings.TrimSuffix(strings.Trim(strings.ToLower(raw), "[]"), ".")
}

// RecordHoneypotAssessment stores only explicit, high-confidence signatures.
// Weak heuristics must never automatically suppress a target.
func (s *AssetStore) RecordHoneypotAssessment(taskID int64, a HoneypotAssessment) error {
	if s == nil || s.db == nil {
		return errors.New("asset store is unavailable")
	}
	if taskID <= 0 {
		return errors.New("honeypot assessment task id is required")
	}
	a.Host = normalizeHoneypotHost(a.Host)
	if a.Host == "" {
		return errors.New("honeypot assessment host is required")
	}
	if a.Score < 90 || a.Score > 100 {
		return errors.New("only high-confidence honeypot assessments can be persisted")
	}
	if len(a.Signals) == 0 {
		return errors.New("honeypot assessment requires at least one explicit signal")
	}
	if a.Level == "" {
		a.Level = "suspected"
	}
	if a.Level != "suspected" {
		return errors.New("honeypot assessment level must be suspected")
	}
	now := time.Now().UTC()
	if a.ObservedAt.IsZero() {
		a.ObservedAt = now
	}
	if a.ExpiresAt.IsZero() {
		a.ExpiresAt = now.Add(HoneypotAssessmentTTL)
	}
	if !a.ExpiresAt.After(now) {
		return errors.New("honeypot assessment expiry must be in the future")
	}
	signals, err := json.Marshal(a.Signals)
	if err != nil {
		return fmt.Errorf("marshal honeypot signals: %w", err)
	}
	_, err = s.db.Exec(`
INSERT INTO target_honeypot_assessments(task_id, host, score, level, signals, observed_at, expires_at)
VALUES ($1, $2, $3, $4, $5::jsonb, $6, $7)
ON CONFLICT (task_id, host) DO UPDATE SET
    score = GREATEST(target_honeypot_assessments.score, EXCLUDED.score),
    level = 'suspected',
    signals = CASE
        WHEN EXCLUDED.score >= target_honeypot_assessments.score THEN EXCLUDED.signals
        ELSE target_honeypot_assessments.signals
    END,
    observed_at = GREATEST(target_honeypot_assessments.observed_at, EXCLUDED.observed_at),
    expires_at = GREATEST(target_honeypot_assessments.expires_at, EXCLUDED.expires_at)`,
		taskID, a.Host, a.Score, a.Level, string(signals), a.ObservedAt, a.ExpiresAt)
	return err
}

// HoneypotAssessments returns unexpired assessments for the supplied hosts.
func (s *AssetStore) HoneypotAssessments(taskID int64, hosts []string) ([]HoneypotAssessment, error) {
	if s == nil || s.db == nil {
		return nil, errors.New("asset store is unavailable")
	}
	if taskID <= 0 {
		return nil, errors.New("honeypot assessment task id is required")
	}
	seen := make(map[string]bool, len(hosts))
	normalized := make([]string, 0, len(hosts))
	for _, raw := range hosts {
		host := normalizeHoneypotHost(raw)
		if host != "" && !seen[host] {
			seen[host] = true
			normalized = append(normalized, host)
		}
	}
	if len(normalized) == 0 {
		return nil, nil
	}
	marks := make([]string, len(normalized))
	args := make([]any, len(normalized)+1)
	args[0] = taskID
	for i, host := range normalized {
		marks[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = host
	}
	rows, err := s.db.Query(`
SELECT host, score, level, signals::text, observed_at, expires_at
FROM target_honeypot_assessments
WHERE task_id=$1 AND host IN (`+strings.Join(marks, ",")+`) AND expires_at > now()
ORDER BY score DESC, observed_at DESC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []HoneypotAssessment
	for rows.Next() {
		var a HoneypotAssessment
		var raw []byte
		if err := rows.Scan(&a.Host, &a.Score, &a.Level, &raw, &a.ObservedAt, &a.ExpiresAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(raw, &a.Signals); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// HoneypotAssessmentsForTask returns only assessments matching hosts already
// associated with this task; it never broadens the task's target scope.
func (s *AssetStore) HoneypotAssessmentsForTask(taskID int64) ([]HoneypotAssessment, error) {
	if s == nil || s.db == nil {
		return nil, errors.New("asset store is unavailable")
	}
	if taskID <= 0 {
		return nil, nil
	}
	rows, err := s.db.Query(`
SELECT host, score, level, signals::text, observed_at, expires_at
FROM target_honeypot_assessments
WHERE task_id=$1 AND expires_at > now()
ORDER BY score DESC, observed_at DESC
LIMIT 100`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []HoneypotAssessment
	for rows.Next() {
		var a HoneypotAssessment
		var raw []byte
		if err := rows.Scan(&a.Host, &a.Score, &a.Level, &raw, &a.ObservedAt, &a.ExpiresAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(raw, &a.Signals); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
