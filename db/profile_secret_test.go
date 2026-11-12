package db

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestProfileSecretEnvelope(t *testing.T) {
	const secret = "test-model-credential-not-real"
	t.Setenv(profileSecretEnv, "")
	plain, encrypted, err := encodeProfileSecret(secret)
	if err != nil || plain != secret || encrypted != "" {
		t.Fatal("legacy mode did not preserve the plaintext representation")
	}
	t.Setenv(profileSecretEnv, strings.Repeat("01", 32))
	plain, encrypted, err = encodeProfileSecret(secret)
	if err != nil || plain != "" || !strings.HasPrefix(encrypted, profileSecretPrefix) || strings.Contains(encrypted, secret) {
		t.Fatal("expected only an encrypted representation")
	}
	_, second, err := encodeProfileSecret(secret)
	if err != nil || second == encrypted {
		t.Fatal("encryption did not use a fresh nonce")
	}
	got, err := decodeProfileSecret("", encrypted)
	if err != nil || got != secret {
		t.Fatal("encrypted key did not round-trip")
	}
	data, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(encrypted, profileSecretPrefix))
	if err != nil {
		t.Fatal(err)
	}
	data[len(data)-1] ^= 1
	corrupt := profileSecretPrefix + base64.RawStdEncoding.EncodeToString(data)
	if got, err := decodeProfileSecret(secret, corrupt); got != "" || !errors.Is(err, errProfileSecretUnavailable) {
		t.Fatal("damaged ciphertext was accepted or fell back to plaintext")
	}
	for _, value := range []string{"", strings.Repeat("02", 32)} {
		t.Setenv(profileSecretEnv, value)
		if got, err := decodeProfileSecret("", encrypted); got != "" || !errors.Is(err, errProfileSecretUnavailable) {
			t.Fatal("missing/wrong encryption key was accepted")
		}
		if got, err := decodeProfileSecret(secret, ""); got != secret || err != nil {
			t.Fatal("legacy plaintext stopped working")
		}
	}
	t.Setenv(profileSecretEnv, "invalid")
	if _, _, err := encodeProfileSecret(secret); !errors.Is(err, errProfileSecretConfig) {
		t.Fatal("invalid configuration silently downgraded")
	}
}

func TestEncryptedModelProfileStorage(t *testing.T) {
	t.Setenv(profileSecretEnv, "")
	d, err := Open(testDSN(t))
	if err != nil {
		t.Skipf("postgres unavailable (%v)", err)
	}
	defer d.Close()
	const secret = "test-model-secret-abcdef"
	name := fmt.Sprintf("test-encrypted-profile-%d", time.Now().UnixNano())
	p := &LLMProfile{Name: name, Format: "openai", Model: "test", APIKey: secret}
	id, err := d.SaveProfile(p)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Exec("DELETE FROM llm_profiles WHERE id=$1", id)
	p.ID = id
	readStorage := func() (string, string) {
		t.Helper()
		var plain, encrypted string
		if err := d.QueryRow("SELECT COALESCE(api_key,''),COALESCE(api_key_ciphertext,'') FROM llm_profiles WHERE id=$1", id).Scan(&plain, &encrypted); err != nil {
			t.Fatal(err)
		}
		return plain, encrypted
	}
	if plain, encrypted := readStorage(); plain != secret || encrypted != "" {
		t.Fatal("legacy storage changed")
	}
	t.Setenv(profileSecretEnv, strings.Repeat("11", 32))
	if got, err := d.ProfileByID(id); err != nil || got == nil || got.APIKey != secret {
		t.Fatal("old profile could not be loaded after enabling encryption")
	}
	// A blank edit intentionally retains the existing storage representation.
	p.APIKey = ""
	p.Model = "edited"
	if _, err := d.SaveProfile(p); err != nil {
		t.Fatal(err)
	}
	if plain, encrypted := readStorage(); plain != secret || encrypted != "" {
		t.Fatal("blank update changed legacy key")
	}
	// Explicitly re-saving the model key migrates only this row.
	p.APIKey = secret
	if _, err := d.SaveProfile(p); err != nil {
		t.Fatal(err)
	}
	plain, encrypted := readStorage()
	if plain != "" || !strings.HasPrefix(encrypted, profileSecretPrefix) || strings.Contains(encrypted, secret) {
		t.Fatal("database retained the model key in plaintext")
	}
	p.APIKey = ""
	p.Model = "edited-again"
	if _, err := d.SaveProfile(p); err != nil {
		t.Fatal(err)
	}
	if a, b := readStorage(); a != "" || b != encrypted {
		t.Fatal("blank update replaced encrypted key")
	}
	if got, err := d.ProfileByID(id); err != nil || got == nil || got.APIKey != secret || got.Model != p.Model {
		t.Fatal("encrypted profile read failed")
	}
	chain, err := d.PoolProfiles()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, candidate := range chain {
		if candidate.ID == id {
			found = candidate.APIKey == secret
		}
	}
	if !found {
		t.Fatal("encrypted profile missing from model pool")
	}
	list, err := d.ListProfiles()
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(list)
	if err != nil || strings.Contains(string(b), secret) || strings.Contains(string(b), encrypted) {
		t.Fatal("profile list exposed a secret")
	}

	// An encrypted record must not fall back to plaintext if its key is lost.
	t.Setenv(profileSecretEnv, "")
	if got, err := d.ProfileByID(id); got != nil || !errors.Is(err, errProfileSecretUnavailable) {
		t.Fatal("encrypted profile read without key unexpectedly succeeded")
	}
	p.APIKey = "replacement"
	if _, err := d.SaveProfile(p); !errors.Is(err, errProfileSecretDowngrade) {
		t.Fatal("missing encryption key allowed a plaintext replacement")
	}
	if a, b := readStorage(); a != "" || b != encrypted {
		t.Fatal("rejected replacement changed stored key")
	}
	t.Setenv(profileSecretEnv, strings.Repeat("22", 32))
	if got, err := d.ProfileByID(id); got != nil || !errors.Is(err, errProfileSecretUnavailable) {
		t.Fatal("wrong key was accepted")
	}
	t.Setenv(profileSecretEnv, strings.Repeat("11", 32))
	if got, err := d.ProfileByID(id); err != nil || got == nil || got.APIKey != secret {
		t.Fatal("restoring the original encryption key did not restore the profile")
	}

	// New encrypted profiles also work; keyless Ollama keeps both columns empty.
	newID, err := d.SaveProfile(&LLMProfile{Name: name + "-new", Format: "openai", Model: "test", APIKey: "tiny"})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Exec("DELETE FROM llm_profiles WHERE id=$1", newID)
	var freshPlain, freshCipher, hint string
	if err := d.QueryRow("SELECT COALESCE(api_key,''),COALESCE(api_key_ciphertext,''),COALESCE(api_key_hint,'') FROM llm_profiles WHERE id=$1", newID).Scan(&freshPlain, &freshCipher, &hint); err != nil {
		t.Fatal(err)
	}
	if freshPlain != "" || freshCipher == "" || hint != "••••" {
		t.Fatal("new encrypted profile or short-key masking failed")
	}
	localID, err := d.SaveProfile(&LLMProfile{Name: name + "-local", Format: "ollama", Model: "local"})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Exec("DELETE FROM llm_profiles WHERE id=$1", localID)
	if got, err := d.ProfileByID(localID); err != nil || got == nil || got.APIKey != "" {
		t.Fatal("keyless local profile changed")
	}
}
