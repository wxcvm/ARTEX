package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// Managed operations use the same lock as downloads/updates, so a rollback or
// removal cannot run while another request is replacing the active executable.

func managedToolSlot(root string, spec externalToolSpec, slot string) string {
	return filepath.Join(root, spec.ID, slot)
}

// Only ARTEX-generated slots with a receipt and the expected executable can
// be operated on. Hand-placed binaries elsewhere are never affected.
func managedToolSlotReceipt(root string, spec externalToolSpec, slot string) (*managedToolReceipt, error) {
	dir := managedToolSlot(root, spec, slot)
	info, err := os.Lstat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a managed directory", slot)
	}
	binary := filepath.Join(dir, externalToolBinary(spec.Binary))
	info, err = os.Lstat(binary)
	if err != nil || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s does not contain a regular %s executable", slot, spec.Binary)
	}
	receiptBytes, err := os.ReadFile(filepath.Join(dir, "receipt.json"))
	if err != nil {
		return nil, fmt.Errorf("%s has no ARTEX installation receipt: %w", slot, err)
	}
	var receipt managedToolReceipt
	if err := json.Unmarshal(receiptBytes, &receipt); err != nil ||
		!externalToolReleaseTag.MatchString(receipt.Version) ||
		!externalToolDigest.MatchString(receipt.SHA256) {
		return nil, fmt.Errorf("%s has an invalid ARTEX installation receipt", slot)
	}
	return &receipt, nil
}

// Previous becomes current and the formerly current version becomes the new
// previous version. Repeating the operation reverses it again.
func rollbackManagedExternalTool(root string, spec externalToolSpec) (string, error) {
	externalToolInstallMu.Lock()
	defer externalToolInstallMu.Unlock()

	current := managedToolSlot(root, spec, "current")
	previous := managedToolSlot(root, spec, ".previous")
	staging := managedToolSlot(root, spec, ".rollback-swap")
	currentReceipt, err := managedToolSlotReceipt(root, spec, "current")
	if err != nil { return "", err }
	previousReceipt, err := managedToolSlotReceipt(root, spec, ".previous")
	if err != nil { return "", err }
	if currentReceipt == nil || previousReceipt == nil {
		return "", errors.New("no managed current and previous versions available for rollback")
	}
	if _, err := os.Lstat(staging); err == nil {
		return "", errors.New("incomplete previous rollback; inspect .rollback-swap before retrying")
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if err := os.Rename(current, staging); err != nil {
		return "", fmt.Errorf("cannot close/move current tool: %w", err)
	}
	if err := os.Rename(previous, current); err != nil {
		_ = os.Rename(staging, current)
		return "", fmt.Errorf("cannot activate previous tool: %w", err)
	}
	if err := os.Rename(staging, previous); err != nil {
		// Don't delete a working version or its receipt, even if cleanup failed.
		_ = os.Rename(current, previous)
		_ = os.Rename(staging, current)
		return "", fmt.Errorf("cannot finish version swap: %w", err)
	}
	return previousReceipt.Version, nil
}

// Remove only ARTEX-managed executables and their receipts. Unknown files in
// those directories are kept; never delete user-provided tools or PATH entries.
func uninstallManagedExternalTool(root string, spec externalToolSpec) error {
	externalToolInstallMu.Lock()
	defer externalToolInstallMu.Unlock()

	current, err := managedToolSlotReceipt(root, spec, "current")
	if err != nil { return err }
	previous, err := managedToolSlotReceipt(root, spec, ".previous")
	if err != nil { return err }
	if current == nil && previous == nil {
		return errors.New("no ARTEX-managed installation found")
	}
	for _, slot := range []string{"current", ".previous"} {
		receipt := current
		if slot == ".previous" { receipt = previous }
		if receipt == nil { continue }
		dir := managedToolSlot(root, spec, slot)
		if err := os.Remove(filepath.Join(dir, externalToolBinary(spec.Binary))); err != nil {
			return fmt.Errorf("cannot remove %s executable (may be running): %w", slot, err)
		}
		if err := os.Remove(filepath.Join(dir, "receipt.json")); err != nil {
			return fmt.Errorf("cannot remove %s receipt: %w", slot, err)
		}
		_ = os.Remove(dir) // Keep the directory if the user placed other files there.
	}
	_ = os.Remove(filepath.Join(root, spec.ID))
	return nil
}

func (s *Server) externalToolRollback(w http.ResponseWriter, r *http.Request) {
	spec, found := externalToolSpecByID(r.PathValue("id"))
	if !found { writeJSON(w, http.StatusNotFound, map[string]any{"error":"unknown external tool"}); return }
	root, err := managedExternalToolRoot()
	if err != nil { writeJSON(w, http.StatusBadRequest, map[string]any{"error":err.Error()}); return }
	version, err := rollbackManagedExternalTool(root, spec)
	if err != nil { writeJSON(w, http.StatusConflict, map[string]any{"error":err.Error()}); return }
	writeJSON(w, http.StatusOK, map[string]any{"ok":true, "version":version})
}

func (s *Server) externalToolUninstall(w http.ResponseWriter, r *http.Request) {
	spec, found := externalToolSpecByID(r.PathValue("id"))
	if !found { writeJSON(w, http.StatusNotFound, map[string]any{"error":"unknown external tool"}); return }
	root, err := managedExternalToolRoot()
	if err != nil { writeJSON(w, http.StatusBadRequest, map[string]any{"error":err.Error()}); return }
	if err := uninstallManagedExternalTool(root, spec); err != nil {
		writeJSON(w, http.StatusConflict, map[string]any{"error":err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok":true})
}
