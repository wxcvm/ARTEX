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

func TestImportDSHReviewFeedbackPreservesNewerAndEqualTimestamps(t *testing.T) {
    d, err := Open(testDSN(t))
    if err != nil {
        t.Skipf("postgres unavailable (%v) — skipping", err)
    }
    defer d.Close()
    fid, err := d.AddFinding(0, 0, "review-import-test", "Review import test", "info", "summary", "evidence", "tester", nil)
    if err != nil {
        t.Fatal(err)
    }
    defer d.DeleteFinding(fid)
    input := DSHReviewFeedback{
        Schema: DSHReviewFeedbackSchema,
        Source: "dsh-purge",
        Engagement: "review-import-test",
        Reviews: []DSHReviewEntry{{
            FindingID: fid, Status: "needs_evidence", Note: "newer review",
            UpdatedAt: "2026-11-10T09:00:00Z",
        }},
    }
    check := func(wantImported, wantSkipped int) {
        t.Helper()
        got, err := d.ImportDSHReviewFeedback(input)
        if err != nil {
            t.Fatal(err)
        }
        if got.Imported != wantImported || got.SkippedNotNewer != wantSkipped || got.Missing != 0 || got.TaskMismatch != 0 {
            t.Fatalf("import result: got %+v, want imported=%d skipped=%d", got, wantImported, wantSkipped)
        }
    }
    check(1, 0)
    check(0, 1) // Identical re-import is a no-op.
    input.Reviews[0].Status = "pending"
    input.Reviews[0].Note = "same timestamp, different content"
    check(0, 1) // Equal timestamp must not overwrite the saved review.
    input.Reviews[0].UpdatedAt = "2026-11-10T08:00:00Z"
    check(0, 1)
    var status, note string
    if err := d.QueryRow("SELECT review_status,note FROM finding_dsh_reviews WHERE finding_id=$1 AND dsh_engagement=$2", fid, input.Engagement).Scan(&status, &note); err != nil {
        t.Fatal(err)
    }
    if status != "needs_evidence" || note != "newer review" {
        t.Fatalf("saved newer review changed: %q %q", status, note)
    }
    input.Reviews[0].UpdatedAt = "2026-11-10T10:00:00Z"
    input.Reviews[0].Status = "reviewed"
    input.Reviews[0].Note = "new review"
    check(1, 0)
    if err := d.QueryRow("SELECT review_status,note FROM finding_dsh_reviews WHERE finding_id=$1 AND dsh_engagement=$2", fid, input.Engagement).Scan(&status, &note); err != nil {
        t.Fatal(err)
    }
    if status != "reviewed" || note != "new review" {
        t.Fatalf("new review not saved: %q %q", status, note)
    }
}
