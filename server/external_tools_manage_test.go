package server

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func writeManagedToolFixture(t *testing.T, root string, spec externalToolSpec, slot string, version string, value string) string {
	t.Helper()
	dir := managedToolSlot(root, spec, slot)
	if err := os.MkdirAll(dir, 0755); err != nil { t.Fatal(err) }
	exe := filepath.Join(dir, externalToolBinary(spec.Binary))
	if err := os.WriteFile(exe, []byte(value), 0755); err != nil { t.Fatal(err) }
	data, err := json.Marshal(managedToolReceipt{Version: version, SHA256: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"})
	if err != nil { t.Fatal(err) }
	if err := os.WriteFile(filepath.Join(dir, "receipt.json"), data, 0644); err != nil { t.Fatal(err) }
	return exe
}

func readToolContents(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil { t.Fatal(err) }
	return string(data)
}

func TestExternalToolRollbackSwapsVersions(t *testing.T) {
	root := t.TempDir()
	spec := externalToolSpec{ID: "nuclei", Binary: "nuclei"}
	current := writeManagedToolFixture(t, root, spec, "current", "v3.11.1", "new release")
	previous := writeManagedToolFixture(t, root, spec, ".previous", "v3.10.0", "old release")

	version, err := rollbackManagedExternalTool(root, spec)
	if err != nil || version != "v3.10.0" { t.Fatalf("rollback version=%q, err=%v", version, err) }
	if got := readToolContents(t, current); got != "old release" { t.Fatalf("current=%q", got) }
	if got := readToolContents(t, previous); got != "new release" { t.Fatalf("previous=%q", got) }
	receipt, err := managedToolSlotReceipt(root, spec, "current")
	if err != nil || receipt == nil || receipt.Version != "v3.10.0" {
		t.Fatalf("rollback receipt=%+v, err=%v", receipt, err)
	}
	version, err = rollbackManagedExternalTool(root, spec)
	if err != nil || version != "v3.11.1" { t.Fatalf("second rollback version=%q, err=%v", version, err) }
	if got := readToolContents(t, current); got != "new release" { t.Fatalf("second current=%q", got) }
}

func TestExternalToolRollbackKeepsUnmanagedFiles(t *testing.T) {
	root := t.TempDir()
	spec := externalToolSpec{ID: "ffuf", Binary: "ffuf"}
	_, err := rollbackManagedExternalTool(root, spec)
	if err == nil { t.Fatal("rollback should need two managed versions") }
	current := writeManagedToolFixture(t, root, spec, "current", "v2.3.0", "current release")
	_, err = rollbackManagedExternalTool(root, spec)
	if err == nil { t.Fatal("rollback should not work without a prior version") }
	if got := readToolContents(t, current); got != "current release" { t.Fatalf("current corrupted: %q", got) }
}

func TestExternalToolUninstallPreservesCustomFiles(t *testing.T) {
	root := t.TempDir()
	spec := externalToolSpec{ID: "httpx", Binary: "httpx"}
	current := writeManagedToolFixture(t, root, spec, "current", "v1.12.0", "current")
	previous := writeManagedToolFixture(t, root, spec, ".previous", "v1.11.0", "previous")
	extra := filepath.Join(filepath.Dir(current), "my-settings.txt")
	if err := os.WriteFile(extra, []byte("keep"), 0644); err != nil { t.Fatal(err) }
	custom := filepath.Join(root, "user-supplied", "httpx.exe")
	if err := os.MkdirAll(filepath.Dir(custom), 0755); err != nil { t.Fatal(err) }
	if err := os.WriteFile(custom, []byte("outside of manager"), 0755); err != nil { t.Fatal(err) }

	if err := uninstallManagedExternalTool(root, spec); err != nil { t.Fatal(err) }
	for _, file := range []string{current, previous} {
		if _, err := os.Stat(file); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("managed executable remains %q: %v", file, err)
		}
	}
	if got := readToolContents(t, extra); got != "keep" { t.Fatalf("user settings removed: %q", got) }
	if got := readToolContents(t, custom); got != "outside of manager" { t.Fatalf("custom binary removed: %q", got) }
	if err := uninstallManagedExternalTool(root, spec); err == nil {
		t.Fatal("second uninstall should report that nothing managed remains")
	}
}

func TestExternalToolUninstallWithoutPriorVersion(t *testing.T) {
	root := t.TempDir()
	spec := externalToolSpec{ID: "dnsx", Binary: "dnsx"}
	current := writeManagedToolFixture(t, root, spec, "current", "v1.3.1", "current")
	if err := uninstallManagedExternalTool(root, spec); err != nil { t.Fatal(err) }
	if _, err := os.Stat(current); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("executable still present: %v", err)
	}
}

func TestExternalToolInvalidReceiptDoesNotDelete(t *testing.T) {
	root := t.TempDir()
	spec := externalToolSpec{ID: "katana", Binary: "katana"}
	current := writeManagedToolFixture(t, root, spec, "current", "v1.8.0", "user file")
	if err := os.WriteFile(filepath.Join(filepath.Dir(current), "receipt.json"), []byte("{}"), 0644); err != nil { t.Fatal(err) }
	if err := uninstallManagedExternalTool(root, spec); err == nil { t.Fatal("invalid receipt should be rejected") }
	if got := readToolContents(t, current); got != "user file" { t.Fatalf("unverified file deleted: %q", got) }
}

func TestExternalToolPreviousBackupIsNotActive(t *testing.T) {
	root := t.TempDir()
	spec := externalToolSpec{ID: "artex-inventory-backup-test", Binary: "artex-inventory-backup-test"}
	writeManagedToolFixture(t, root, spec, ".previous", "v1.0.0", "saved backup")
	path, origin := externalToolExecutable(spec, []externalToolRoot{{dir: root, origin: "managed"}})
	if path != "" || origin != "" {
		t.Fatalf("backup presented as active executable: %q (%q)", path, origin)
	}
}
