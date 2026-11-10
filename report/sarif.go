package report

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Autumn-27/artex/db"
)

// SARIF is an interoperable view of persisted ARTEX findings. It intentionally
// does not fabricate source-code locations: ARTEX findings may describe hosts,
// services, or HTTP behavior rather than a source file and line number.
type sarifText struct {
	Text string `json:"text"`
}

type sarifRule struct {
	ID               string    `json:"id"`
	ShortDescription sarifText `json:"shortDescription"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri"`
	Rules          []sarifRule `json:"rules"`
}

type sarifResultProperties struct {
	FindingID  int64   `json:"findingId"`
	TaskID     *int64  `json:"taskId,omitempty"`
	AssetIDs   []int64 `json:"assetIds,omitempty"`
	Severity              string `json:"severity"`
	Status                string `json:"status"`
	TrafficCount          int    `json:"trafficCount"`
	EvidenceVersion       int64  `json:"evidenceVersion"`
	ReportEvidenceVersion int64  `json:"reportEvidenceVersion"`
	ReportStale           bool   `json:"reportStale"`
}

type sarifResult struct {
	RuleID              string                `json:"ruleId"`
	Level               string                `json:"level"`
	Kind                string                `json:"kind"`
	Message             sarifText             `json:"message"`
	PartialFingerprints map[string]string     `json:"partialFingerprints"`
	Properties          sarifResultProperties `json:"properties"`
}

type sarifRun struct {
	Tool    struct {
		Driver sarifDriver `json:"driver"`
	} `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifLog struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

func sarifRuleID(class string) string {
	canonical := strings.ToLower(strings.TrimSpace(class))
	if canonical == "" {
		canonical = "uncategorized"
	}
	sum := sha256.Sum256([]byte(canonical))
	return "ARTEX-" + hex.EncodeToString(sum[:8])
}

func sarifLevel(severity string) string {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "critical", "high":
		return "error"
	case "medium":
		return "warning"
	default:
		return "note"
	}
}

// The SARIF kind records the finding's triage outcome; an unreviewed result
// must not be represented as a confirmed vulnerability.
func sarifKind(status string) string {
	switch status {
	case db.FindingConfirmed, db.FindingRiskAccepted:
		return "fail"
	case db.FindingFixed, db.FindingResolved:
		return "pass"
	case db.FindingFalsePositive, db.FindingDuplicate, db.FindingIgnored:
		return "notApplicable"
	default:
		return "review"
	}
}

// FindingsSARIF exports SARIF 2.1.0 without altering the input slice or
// exposing raw HTTP evidence and transcripts in third-party report viewers.
// It contains no physical locations because findings do not have verified
// source path/line data; downstream code-scanning alert support may vary.
func FindingsSARIF(findings []*db.DBFinding) ([]byte, error) {
	items := append([]*db.DBFinding(nil), findings...)
	sortFindingsForExport(items)

	classes := make(map[string]string)
	results := make([]sarifResult, 0, len(items))
	for _, f := range items {
		class := strings.TrimSpace(f.VulnClass)
		if class == "" {
			class = "Uncategorized"
		}
		id := sarifRuleID(class)
		if _, ok := classes[id]; !ok {
			classes[id] = class
		}
		assets := append([]int64(nil), f.AssetIDs...)
		sort.Slice(assets, func(i, j int) bool { return assets[i] < assets[j] })
		message := findingTitle(f)
		if summary := strings.TrimSpace(f.Summary); summary != "" {
			message = fmt.Sprintf("%s: %s", message, summary)
		}
		results = append(results, sarifResult{
			RuleID: id,
			Level:  sarifLevel(f.Severity),
			Kind:   sarifKind(f.Status),
			Message: sarifText{Text: message},
			PartialFingerprints: map[string]string{
				"artexFindingId": strconv.FormatInt(f.ID, 10),
			},
			Properties: sarifResultProperties{
				FindingID: f.ID,
				TaskID: f.TaskID,
				AssetIDs: assets,
				Severity:              f.Severity,
				Status:                f.Status,
				TrafficCount:          f.TrafficCount,
				EvidenceVersion:       f.EvidenceVersion,
				ReportEvidenceVersion: f.ReportEvidenceVersion,
				ReportStale:           f.Report != "" && f.EvidenceVersion != f.ReportEvidenceVersion,
			},
		})
	}

	ruleIDs := make([]string, 0, len(classes))
	for id := range classes {
		ruleIDs = append(ruleIDs, id)
	}
	sort.Strings(ruleIDs)
	rules := make([]sarifRule, 0, len(ruleIDs))
	for _, id := range ruleIDs {
		rules = append(rules, sarifRule{
			ID: id, ShortDescription: sarifText{Text: classes[id]},
		})
	}

	run := sarifRun{Results: results}
	run.Tool.Driver = sarifDriver{
		Name: "ARTEX", InformationURI: "https://github.com/wxcvm/ARTEX",
		Rules: rules,
	}
	return json.MarshalIndent(sarifLog{
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Version: "2.1.0",
		Runs:    []sarifRun{run},
	}, "", "  ")
}
