package agent

import (
	"strings"
	"testing"
)

func TestScanIngestDetectsExplicitNmapHoneypotFingerprint(t *testing.T) {
	xml := `<nmaprun>
<host><status state="up"/><address addr="198.51.100.77" addrtype="ipv4"/>
<hostnames><hostname name="decoy.example.test" type="user"/></hostnames>
<ports><port protocol="tcp" portid="22"><state state="open"/>
<service name="ssh" product="Cowrie SSH Honeypot" version="2.5"/>
</port></ports></host></nmaprun>`
	s := parseScan(t, "nmap-xml", xml)
	if len(s.Honeypots) != 2 {
		t.Fatalf("assessments=%+v, want IP and scanned-hostname aliases", s.Honeypots)
	}
	hosts := map[string]bool{}
	foundFact := false
	for _, a := range s.Honeypots {
		hosts[a.Host] = true
		if a.Score < 90 || a.Level != "suspected" ||
			!strings.Contains(strings.Join(a.Signals, " "), "nmap.service.product:cowrie") {
			t.Fatalf("unexpected assessment: %+v", a)
		}
	}
	for _, fact := range s.Facts {
		if fact.Source == factSourceHoneypot {
			foundFact = true
		}
	}
	if !hosts["198.51.100.77"] || !hosts["decoy.example.test"] {
		t.Fatalf("IP/domain aliases missing: %+v", s.Honeypots)
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

func TestScanIngestIgnoresHoneypotWordsInUntrustedPageText(t *testing.T) {
	line := `{"url":"https://normal.example.test","ip":"198.51.100.77","status_code":200,"title":"Honeypot detection guide","tech":["Nginx"],"webserver":"nginx","tls":{"subject_dn":"CN=honeypot.example.test"}}`
	s := parseScan(t, "httpx-jsonl", line)
	if len(s.Honeypots) != 0 {
		t.Fatalf("page title or certificate metadata must not trigger quarantine: %+v", s.Honeypots)
	}
}

func TestScanIngestMatchesProductNamesAsWholeTokens(t *testing.T) {
	if got := explicitHoneypotSignals(
		honeypotField{name: "nmap.service.product", value: "Honeywell industrial controller"},
		honeypotField{name: "nmap.service.version", value: "Cowrie 2.5"},
		honeypotField{name: "httpx.title", value: "OpenCanary"},
	); len(got) != 0 {
		t.Fatalf("untrusted or unrelated fields triggered detection: %v", got)
	}
	if got := explicitHoneypotSignals(honeypotField{name: "nmap.service.product", value: "Generic honeypot service"}); len(got) != 0 {
		t.Fatalf("generic honeypot wording is not a product fingerprint: %v", got)
	}
	got := explicitHoneypotSignals(
		honeypotField{name: "nmap.service.product", value: "Cowrie SSH Honeypot"},
		honeypotField{name: "httpx.tech", value: "OpenCanary"},
	)
	if len(got) != 2 {
		t.Fatalf("explicit product identities should match: %v", got)
	}
}
