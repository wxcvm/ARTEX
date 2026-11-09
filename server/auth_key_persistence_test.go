package server

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJWTKeyPersistsAcrossStartup(t *testing.T) {
	keys := filepath.Join(t.TempDir(), "keys")
	data := t.TempDir()
	key1, err := loadOrCreateJWTKey(keys, data)
	if err != nil {
		t.Fatal(err)
	}
	if len(key1) < 32 {
		t.Fatalf("new key length=%d", len(key1))
	}
	token, err := signJWT(key1)
	if err != nil {
		t.Fatal(err)
	}
	key2, err := loadOrCreateJWTKey(keys, data)
	if err != nil {
		t.Fatal(err)
	}
	if string(key1) != string(key2) || !verifyJWT(token, key2) {
		t.Fatal("restart changed JWT key or invalidated existing sessions")
	}
	if _, err := os.Stat(filepath.Join(data, jwtKeyFilename)); !os.IsNotExist(err) {
		t.Fatalf("key unexpectedly stored in browsable data dir: %v", err)
	}
}

func TestJWTKeyMigratesFromBrowsableData(t *testing.T) {
	base := t.TempDir()
	data := filepath.Join(base, "data")
	keys := filepath.Join(base, "keys")
	if err := os.MkdirAll(data, 0o700); err != nil {
		t.Fatal(err)
	}
	legacy := strings.Repeat("A", 32)
	if err := os.WriteFile(filepath.Join(data, jwtKeyFilename), []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}
	actual, err := loadOrCreateJWTKey(keys, data)
	if err != nil {
		t.Fatal(err)
	}
	if string(actual) != legacy {
		t.Fatal("legacy JWT key was not retained")
	}
	if _, err := os.Stat(filepath.Join(data, jwtKeyFilename)); !os.IsNotExist(err) {
		t.Fatalf("legacy key is still in browsable data: %v", err)
	}
}

func TestJWTKeyMigratesFromLegacyExecutableDir(t *testing.T) {
	legacyDir := t.TempDir()
	original := jwtLegacyBaseDir
	jwtLegacyBaseDir = func() string { return legacyDir }
	t.Cleanup(func() { jwtLegacyBaseDir = original })
	legacy := strings.Repeat("B", 32)
	path := filepath.Join(legacyDir, jwtKeyFilename)
	if err := os.WriteFile(path, []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}
	keys := filepath.Join(t.TempDir(), "keys")
	actual, err := loadOrCreateJWTKey(keys, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if string(actual) != legacy {
		t.Fatal("legacy executable key was replaced")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("old key should have been moved to persistent key dir: %v", err)
	}
}

func TestJWTKeyNeverOverwritesExistingKey(t *testing.T) {
	keys := t.TempDir()
	path := filepath.Join(keys, jwtKeyFilename)
	bad := "short-existing-key"
	if err := os.WriteFile(path, []byte(bad), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadOrCreateJWTKey(keys, t.TempDir()); err == nil {
		t.Fatal("corrupted signing key should not be overwritten silently")
	}
	still, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(still) != bad {
		t.Fatal("existing signing key modified")
	}
}

func TestJWTKeyExistingDestinationTakesPrecedence(t *testing.T) {
	keys := t.TempDir()
	data := t.TempDir()
	current := strings.Repeat("C", 32)
	legacy := strings.Repeat("D", 32)
	if err := os.WriteFile(filepath.Join(keys, jwtKeyFilename), []byte(current), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(data, jwtKeyFilename), []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}
	actual, err := loadOrCreateJWTKey(keys, data)
	if err != nil {
		t.Fatal(err)
	}
	if string(actual) != current {
		t.Fatal("existing destination key overwritten by legacy key")
	}
}
