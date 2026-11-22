package agent

import (
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
)

func TestHoneypotRegistryRequiresExplicitUnexpiredSignal(t *testing.T) {
	r := NewHoneypotRegistry()
	if r.Record(db.HoneypotAssessment{Host: "weak.example.test", Score: 65, Level: "suspected", Signals: []string{"cloud ASN"}}) {
		t.Fatal("weak heuristic must not quarantine a target")
	}
	if r.Record(db.HoneypotAssessment{Host: "expired.example.test", Score: 95, Level: "suspected", Signals: []string{"nmap.service.product:cowrie"}, ExpiresAt: time.Now().Add(-time.Minute)}) {
		t.Fatal("expired assessment must not quarantine a target")
	}
	if r.AllQuarantined(nil) {
		t.Fatal("empty target set must never be treated as quarantined")
	}
	if r.AllQuarantined([]string{"weak.example.test"}) {
		t.Fatal("unmatched target must remain eligible")
	}
}

func TestHoneypotRegistryNormalizesAliasesAndRequiresAllTargets(t *testing.T) {
	r := NewHoneypotRegistry()
	now := time.Now()
	r.Record(db.HoneypotAssessment{
		Host: "https://Cowrie.Example.Test:443/path",
		Score: 95, Level: "suspected",
		Signals: []string{"nmap.service.product:cowrie"},
		ObservedAt: now, ExpiresAt: now.Add(time.Hour),
	})
	if !r.IsQuarantined("cowrie.example.test.") {
		t.Fatal("URL and hostname forms should resolve to the same key")
	}
	if r.AllQuarantined([]string{"cowrie.example.test", "198.51.100.77"}) {
		t.Fatal("a mixed intent must not be stopped when only one target is quarantined")
	}
	r.Record(db.HoneypotAssessment{
		Host: "198.51.100.77", Score: 95, Level: "suspected",
		Signals: []string{"httpx.title:opencanary"},
		ObservedAt: now, ExpiresAt: now.Add(time.Hour),
	})
	if !r.AllQuarantined([]string{"cowrie.example.test", "198.51.100.77"}) {
		t.Fatal("all explicit target identities should be quarantined")
	}
	items := r.List()
	if len(items) != 2 || items[0].Host >= items[1].Host {
		t.Fatalf("registry list should be stable and sorted: %+v", items)
	}
}

func TestAssetTargetHostsKeepsURLDomainAndIPAliases(t *testing.T) {
	got := AssetTargetHosts([]*db.Asset{
		{Type: "service", URL: "https://app.example.test:8443/login", IP: "198.51.100.77"},

		{Type: "ip", IP: "198.51.100.77"},
		{Type: "root_domain", Domain: "example.test"},
	})
	want := []string{"198.51.100.77", "app.example.test", "example.test"}
	if len(got) != len(want) {
		t.Fatalf("hosts=%v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("hosts=%v want %v", got, want)
		}
	}
}
