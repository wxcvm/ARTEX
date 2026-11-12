package db

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestModelKeyDiagnosticsCounts(t *testing.T) {
	const secret = "test-credential-never-returned"
	key := strings.Repeat("11", 32)
	t.Setenv(profileSecretEnv, key)
	_, encrypted, err := encodeProfileSecret(secret)
	if err != nil {
		t.Fatal(err)
	}

	snapshot := func() ModelKeyStatus {
		s := newModelKeyStatus()
		s.observe(1, false, "")
		s.observe(2, true, "")
		s.observe(3, false, encrypted)
		s.observe(4, true, encrypted)
		s.observe(5, false, profileSecretPrefix+"broken")
		return s
	}
	s := snapshot()
	if s.KeyState != "configured" || s.Total != 5 || s.Empty != 1 || s.PlaintextOnly != 1 ||
		s.EncryptedOnly != 2 || s.Mixed != 1 || s.Verified != 2 || s.Failed != 1 || s.Unchecked != 0 {
		t.Fatalf("unexpected diagnostic counts: %+v", s)
	}
	if s.Profiles[0].Storage != "empty" || s.Profiles[1].Integrity != "not_applicable" ||
		s.Profiles[2].Integrity != "verified" || s.Profiles[3].Storage != "mixed" ||
		s.Profiles[4].Integrity != "failed" || s.Profiles[4].ID != "5" {
		t.Fatalf("unexpected profile states: %+v", s.Profiles)
	}
	b, err := json.Marshal(s)
	if err != nil || strings.Contains(string(b), secret) || strings.Contains(string(b), encrypted) ||
		strings.Contains(string(b), key) {
		t.Fatal("diagnostics must not expose key material")
	}

	t.Setenv(profileSecretEnv, "")
	s = snapshot()
	if s.KeyState != "missing" || s.Unchecked != 3 || s.Verified != 0 || s.Failed != 0 {
		t.Fatalf("missing key incorrectly classified: %+v", s)
	}
	t.Setenv(profileSecretEnv, "invalid")
	s = snapshot()
	if s.KeyState != "invalid" || s.Unchecked != 3 {
		t.Fatalf("invalid key incorrectly classified: %+v", s)
	}
	t.Setenv(profileSecretEnv, strings.Repeat("22", 32))
	s = snapshot()
	if s.KeyState != "configured" || s.Failed != 3 || s.Verified != 0 {
		t.Fatalf("wrong key incorrectly classified: %+v", s)
	}
}

func TestModelKeyDiagnosticsPostgres(t *testing.T) {
	t.Setenv(profileSecretEnv, strings.Repeat("33", 32))
	d, err := Open(testDSN(t))
	if err != nil {
		t.Skipf("postgres unavailable: %v", err)
	}
	defer d.Close()
	const secret = "fixture-secret-not-real"
	name := fmt.Sprintf("diagnostic-%d", time.Now().UnixNano())
	id, err := d.SaveProfile(&LLMProfile{Name: name, Format: "openai", Model: "test", APIKey: secret})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Exec("DELETE FROM llm_profiles WHERE id=$1", id)
	readStorage := func() (plain, encrypted string) {
		t.Helper()
		if err := d.QueryRow("SELECT COALESCE(api_key,''), COALESCE(api_key_ciphertext,'') FROM llm_profiles WHERE id=$1", id).Scan(&plain, &encrypted); err != nil {
			t.Fatal(err)
		}
		return
	}
	_, original := readStorage()
	check := func(storage, integrity string) {
		t.Helper()
		s, err := d.ModelKeyDiagnostics(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		if s.Total != s.Empty+s.PlaintextOnly+s.EncryptedOnly+s.Mixed ||
			s.EncryptedOnly+s.Mixed != s.Verified+s.Failed+s.Unchecked {
			t.Fatalf("inconsistent diagnostic totals: %+v", s)
		}
		for _, p := range s.Profiles {
			if p.ID == fmt.Sprint(id) {
				if p.Storage != storage || p.Integrity != integrity {
					t.Fatalf("unexpected profile diagnostics: %+v", p)
				}
				return
			}
		}
		t.Fatal("fixture profile missing")
	}
	check("encrypted_only", "verified")
	if plain, encrypted := readStorage(); plain != "" || encrypted != original {
		t.Fatal("diagnostics changed the encrypted row")
	}
	if _, err := d.Exec("UPDATE llm_profiles SET api_key=$1 WHERE id=$2", secret, id); err != nil {
		t.Fatal(err)
	}
	check("mixed", "verified")
	t.Setenv(profileSecretEnv, "")
	check("mixed", "unchecked")
	if plain, encrypted := readStorage(); plain != secret || encrypted != original {
		t.Fatal("diagnostics changed the mixed row")
	}
}
