package agent

import (
	"strings"
	"testing"
)

func TestScanIngestDetectsExplicitNmapHoneypotFingerprint(t *testing.T) {
	xml := `<nmaprun>
<host><status state="up"/><address addr="198.51.100.77" addrtype="ipv4"/>
<ports><port protocol="tcp" portid="22"><state state="open"/>
<service name="ssh" product="Cowrie SSH Honeypot" version="2.5"/>
</port></ports></host></nmaprun>`
	s := parseScan(t, "nmap-xml", xml)
	if len(s.Honeypots) != 1 {
		t.Fatalf("assessments=%+v, want one", s.Honeypots)
	}
	a := s.Honeypots[0]
	if a.Host != "198.51.100.77" || a.Score < 90 || a.Level != "suspected" ||
		!strings.Contains(strings.Join(a.Signals, " "), "nmap.service.product:cowrie") {
		t.Fatalf("unexpected assessment: %+v", a)
	}
	foundFact := false
	for _, fact := range s.Facts {
		if fact.Source == factSourceHoneypot && strings.Contains(fact.Summary, a.Host) {
			foundFact = true
		}
	}
	if !foundFact {
		t.Fatal("explicit honeypot evidence must be retained as a separate graph fact")
	}
}

func TestScanIngestHttpxStoresHostnameAndIPAliases(t *testing.T) {
	line := `{"url":"https://decoy.example.test","ip":"198.51.100.77","status_code":200,"title":"OpenCanary service","tech":["OpenCanary"],"webserver":"OpenCanary"}`
	s := parseScan(t, "httpx-jsonl", line)
	if len(s.Honeypots) != 2 {
		t.Fatalf("assessments=%+v, want hostname and IP aliases", s.Honeypots)
	}
	hosts := map[string]bool{}
	for _, a := range s.Honeypots {
		hosts[a.Host] = true
		if a.Score < 90 || a.Level != "suspected" {
			t.Fatalf("assessment is not high-confidence: %+v", a)
		}
	}
	if !hosts["decoy.example.test"] || !hosts["198.51.100.77"] {
		t.Fatalf("hostname/IP aliases missing: %+v", s.Honeypots)
	}
}

func TestScanIngestDoesNotFlagOrdinaryCloudAndWebFingerprints(t *testing.T) {
	line := `{"url":"https://normal.example.test","ip":"198.51.100.77","status_code":200,"title":"Welcome to nginx","tech":["Nginx","Cloudflare"],"webserver":"nginx"}`
	s := parseScan(t, "httpx-jsonl", line)
	if len(s.Honeypots) != 0 {
		t.Fatalf("ordinary hosting fingerprints must not quarantine a target: %+v", s.Honeypots)
	}
}
