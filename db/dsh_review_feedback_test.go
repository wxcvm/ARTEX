package db

import (
    "strings"
    "testing"
)

const validDSHReview = `{
 "schema":"artex.dsh.review-feedback.v1","source":"dsh-purge","engagement":"lab-1",
 "exported_at":"2026-11-10T08:00:00Z",
 "reviews":[{"finding_id":51,"task_id":7,"evidence_version":8,
 "status":"needs_evidence","note":"Check updated evidence","updated_at":"2026-11-10T08:00:00Z"}]
}`

func TestParseDSHReviewFeedback(t *testing.T) {
    got, err := ParseDSHReviewFeedback([]byte(validDSHReview))
    if err != nil { t.Fatal(err) }
    if got.Engagement != "lab-1" || len(got.Reviews) != 1 ||
        got.Reviews[0].FindingID != 51 || *got.Reviews[0].EvidenceVersion != 8 ||
        got.Reviews[0].Status != "needs_evidence" {
        t.Fatalf("unexpected record: %+v", got)
    }
    for _, tc := range []struct{name, raw string}{
        {"wrong schema", strings.Replace(validDSHReview, "review-feedback.v1", "bad.v1", 1)},
        {"wrong source", strings.Replace(validDSHReview, "dsh-purge", "other", 1)},
        {"zero source ID", strings.Replace(validDSHReview, `"finding_id":51`, `"finding_id":0`, 1)},
        {"invalid task ID", strings.Replace(validDSHReview, `"task_id":7`, `"task_id":-1`, 1)},
        {"invalid evidence version", strings.Replace(validDSHReview, `"evidence_version":8`, `"evidence_version":-1`, 1)},
        {"verdict injection", strings.Replace(validDSHReview, `"needs_evidence"`, `"confirmed"`, 1)},
        {"duplicate record", strings.Replace(validDSHReview, `"updated_at":"2026-11-10T08:00:00Z"}]`, `"updated_at":"2026-11-10T08:00:00Z"},{"finding_id":51,"status":"reviewed"}]`, 1)},
        {"empty records", strings.Replace(validDSHReview, `"reviews":[{"finding_id":51,"task_id":7,"evidence_version":8,
 "status":"needs_evidence","note":"Check updated evidence","updated_at":"2026-11-10T08:00:00Z"}]`, `"reviews":[]`, 1)},
    } {
        t.Run(tc.name, func(t *testing.T) {
            if _, err := ParseDSHReviewFeedback([]byte(tc.raw)); err == nil {
                t.Fatal("invalid DSH review feedback accepted")
            }
        })
    }
}
