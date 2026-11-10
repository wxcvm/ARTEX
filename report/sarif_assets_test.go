package report

import (
    "encoding/json"
    "strings"
    "testing"
    "time"

    "github.com/Autumn-27/artex/db"
)

func TestFindingsSARIFWithAssetsKeepsIdentityWithoutURLSecrets(t *testing.T) {
    port := 8443
    findings := []*db.DBFinding{{
        ID: 51, Name: "Service finding", VulnClass: "Configuration",
        Severity: "medium", Status: db.FindingPending,
        AssetIDs: []int64{8, 2, 99}, CreatedAt: time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC),
    }}
    assets := map[int64]*db.Asset{
        8: {ID: 8, Type: "service", URL: "https://user:secret@api.example.test:8443/private?token=DO-NOT-EXPORT"},
        2: {ID: 2, Type: "ip", IP: "192.0.2.8", Domain: "HOST.EXAMPLE.TEST.", Port: &port},
    }
    out, err := FindingsSARIFWithAssets(findings, assets)
    if err != nil { t.Fatal(err) }
    if strings.Contains(string(out), "DO-NOT-EXPORT") || strings.Contains(string(out), "secret") ||
       strings.Contains(string(out), "/private") {
        t.Fatalf("SARIF leaked raw URL: %s", out)
    }
    var doc struct {
        Runs []struct {
            Results []struct {
                Properties struct {
                    AssetIDs []int64 `json:"assetIds"`
                    AssetRefs []sarifAssetRef `json:"assetRefs"`
                } `json:"properties"`
            } `json:"results"`
        } `json:"runs"`
    }
    if err := json.Unmarshal(out, &doc); err != nil { t.Fatal(err) }
    p := doc.Runs[0].Results[0].Properties
    if len(p.AssetIDs) != 3 || len(p.AssetRefs) != 2 {
        t.Fatalf("must retain unresolved source IDs without guessing: %+v", p)
    }
    if p.AssetRefs[0].ID != 2 || p.AssetRefs[0].IP != "192.0.2.8" ||
       p.AssetRefs[0].Domain != "host.example.test" || *p.AssetRefs[0].Port != 8443 {
        t.Fatalf("ip asset identity mismatch: %+v", p.AssetRefs[0])
    }
    if p.AssetRefs[1].ID != 8 || p.AssetRefs[1].Domain != "api.example.test" ||
       p.AssetRefs[1].IP != "" || *p.AssetRefs[1].Port != 8443 {
        t.Fatalf("service URL host/port mismatch: %+v", p.AssetRefs[1])
    }
    legacy, err := FindingsSARIF(findings)
    if err != nil { t.Fatal(err) }
    if strings.Contains(string(legacy), "assetRefs") {
        t.Fatal("legacy exporter should not invent resolved asset identities")
    }
}
