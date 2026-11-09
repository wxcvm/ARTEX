package server

import (
    "archive/zip"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)

const (
    offlineBundleMaxBytes int64 = 1024 << 20
    offlineBundleMaxManifest int64 = 64 << 10
)

type offlineBundleTool struct {
    ID string `json:"id"`
    Version string `json:"version"`
    Source string `json:"source"`
    ArchiveSHA256 string `json:"release_zip_sha256"`
    BinarySHA256 string `json:"binary_sha256"`
    SizeBytes int64 `json:"size_bytes"`
}
type offlineBundleManifest struct {
    SchemaVersion int `json:"schema_version"`
    Target string `json:"target"`
    CreatedAt string `json:"created_at"`
    Tools []offlineBundleTool `json:"tools"`
}

// The bundle contains only ARTEX-managed current executables, not arbitrary
// tools discovered through PATH or user-controlled directories.
func exportOfflineToolBundle(root string, dst io.Writer) (int, error) {
    externalToolInstallMu.Lock()
    defer externalToolInstallMu.Unlock()

    catalog, err := loadExternalToolManifest()
    if err != nil { return 0, err }
    writer := zip.NewWriter(dst)
    manifest := offlineBundleManifest{
        SchemaVersion: 1, Target: "windows/amd64",
        CreatedAt: time.Now().UTC().Format(time.RFC3339),
        Tools: make([]offlineBundleTool, 0),
    }
    var total int64
    for _, spec := range catalog.Tools {
        receipt, err := managedToolSlotReceipt(root, spec, "current")
        if err != nil { return 0, fmt.Errorf("%s: %w", spec.ID, err) }
        if receipt == nil { continue }
        exePath := filepath.Join(root, spec.ID, "current", spec.Binary+".exe")
        info, err := os.Lstat(exePath)
        if err != nil { return 0, err }
        if !info.Mode().IsRegular() || info.Size() < 1 || info.Size() > maxExternalToolBinary ||
            total+info.Size() > offlineBundleMaxBytes {
            return 0, fmt.Errorf("invalid or oversized managed binary: %s", spec.ID)
        }
        input, err := os.Open(exePath)
        if err != nil { return 0, err }
        header := &zip.FileHeader{Name: "bin/" + spec.ID + "/" + spec.Binary + ".exe", Method: zip.Deflate}
        header.SetMode(0644)
        out, err := writer.CreateHeader(header)
        if err != nil { input.Close(); return 0, err }
        hash := sha256.New()
        size, copyErr := io.Copy(io.MultiWriter(out, hash), io.LimitReader(input, maxExternalToolBinary+1))
        closeErr := input.Close()
        if copyErr != nil { return 0, copyErr }
        if closeErr != nil { return 0, closeErr }
        if size != info.Size() {
            return 0, fmt.Errorf("managed tool changed during export: %s", spec.ID)
        }
        total += size
        manifest.Tools = append(manifest.Tools, offlineBundleTool{
            ID: spec.ID, Version: receipt.Version, Source: spec.Source,
            ArchiveSHA256: receipt.SHA256, BinarySHA256: hex.EncodeToString(hash.Sum(nil)),
            SizeBytes: size,
        })
    }
    if len(manifest.Tools) == 0 { return 0, errors.New("no ARTEX-managed tools to export") }
    manifestFile, err := writer.Create("manifest.json")
    if err != nil { return 0, err }
    if err := json.NewEncoder(manifestFile).Encode(manifest); err != nil { return 0, err }
    if err := writer.Close(); err != nil { return 0, err }
    return len(manifest.Tools), nil
}

func offlineBundleFileMap(files []*zip.File) (map[string]*zip.File, error) {
    result := make(map[string]*zip.File, len(files))
    for _, entry := range files {
        if entry.FileInfo().IsDir() || !entry.FileInfo().Mode().IsRegular() {
            return nil, fmt.Errorf("unsupported archive member %q", entry.Name)
        }
        if _, exists := result[entry.Name]; exists {
            return nil, fmt.Errorf("duplicate archive member %q", entry.Name)
        }
        result[entry.Name] = entry
    }
    return result, nil
}

type stagedOfflineTool struct {
    spec externalToolSpec
    ready string
}

// Validate the ENTIRE ZIP before changing any current tool. Every filename
// comes from the embedded catalog, never from user-controlled archive paths.
func stageOfflineToolBundle(root, archivePath string) (string, []stagedOfflineTool, error) {
    zr, err := zip.OpenReader(archivePath)
    if err != nil { return "", nil, err }
    defer zr.Close()
    catalog, err := loadExternalToolManifest()
    if err != nil { return "", nil, err }
    if len(zr.File) < 2 || len(zr.File) > len(catalog.Tools)+1 {
        return "", nil, errors.New("unexpected number of offline bundle members")
    }
    files, err := offlineBundleFileMap(zr.File)
    if err != nil { return "", nil, err }
    manifestFile, ok := files["manifest.json"]
    if !ok || manifestFile.UncompressedSize64 > uint64(offlineBundleMaxManifest) {
        return "", nil, errors.New("offline bundle manifest missing or too large")
    }
    reader, err := manifestFile.Open()
    if err != nil { return "", nil, err }
    var manifest offlineBundleManifest
    decodeErr := json.NewDecoder(io.LimitReader(reader, offlineBundleMaxManifest)).Decode(&manifest)
    reader.Close()
    if decodeErr != nil { return "", nil, fmt.Errorf("invalid offline manifest: %w", decodeErr) }
    if manifest.SchemaVersion != 1 || manifest.Target != "windows/amd64" ||
        len(manifest.Tools) < 1 || len(manifest.Tools) > len(catalog.Tools) ||
        len(manifest.Tools)+1 != len(files) {
        return "", nil, errors.New("unsupported offline bundle manifest")
    }
    specs := make(map[string]externalToolSpec, len(catalog.Tools))
    for _, spec := range catalog.Tools { specs[spec.ID] = spec }
    stage, err := os.MkdirTemp(root, ".offline-")
    if err != nil { return "", nil, err }
    cleanup := func(err error) (string, []stagedOfflineTool, error) {
        os.RemoveAll(stage)
        return "", nil, err
    }
    seen := make(map[string]bool)
    var total int64
    staged := make([]stagedOfflineTool, 0, len(manifest.Tools))
    for _, entry := range manifest.Tools {
        spec, ok := specs[entry.ID]
        if !ok || seen[entry.ID] || entry.Source != spec.Source ||
            !externalToolReleaseTag.MatchString(entry.Version) ||
            !externalToolDigest.MatchString(entry.ArchiveSHA256) ||
            !externalToolDigest.MatchString(entry.BinarySHA256) ||
            entry.SizeBytes < 1 || entry.SizeBytes > maxExternalToolBinary ||
            total+entry.SizeBytes > offlineBundleMaxBytes {
            return cleanup(fmt.Errorf("invalid offline tool entry %q", entry.ID))
        }
        seen[entry.ID] = true
        member := files["bin/"+spec.ID+"/"+spec.Binary+".exe"]
        if member == nil || member.UncompressedSize64 != uint64(entry.SizeBytes) {
            return cleanup(fmt.Errorf("missing or invalid executable for %s", spec.ID))
        }
        ready := filepath.Join(stage, spec.ID, "ready")
        if err := os.MkdirAll(ready, 0755); err != nil { return cleanup(err) }
        src, err := member.Open()
        if err != nil { return cleanup(err) }
        output, err := os.OpenFile(filepath.Join(ready, spec.Binary+".exe"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0755)
        if err != nil { src.Close(); return cleanup(err) }
        hash := sha256.New()
        n, copyErr := io.Copy(io.MultiWriter(output, hash), io.LimitReader(src, entry.SizeBytes+1))
        src.Close()
        closeErr := output.Close()
        if copyErr != nil { return cleanup(copyErr) }
        if closeErr != nil { return cleanup(closeErr) }
        if n != entry.SizeBytes || !strings.EqualFold(hex.EncodeToString(hash.Sum(nil)), entry.BinarySHA256) {
            return cleanup(fmt.Errorf("binary SHA-256 mismatch for %s", spec.ID))
        }
        total += n
        receipt := managedToolReceipt{
            Version: entry.Version, SHA256: strings.ToLower(entry.ArchiveSHA256),
            InstalledAt: time.Now().UTC().Format(time.RFC3339),
        }
        bytes, err := json.MarshalIndent(receipt, "", "  ")
        if err != nil { return cleanup(err) }
        if err := os.WriteFile(filepath.Join(ready, "receipt.json"), bytes, 0644); err != nil {
            return cleanup(err)
        }
        staged = append(staged, stagedOfflineTool{spec: spec, ready: ready})
    }
    return stage, staged, nil
}

// Each tool switches independently, preserving its old current version for
// rollback. The surrounding import verifies all bundle contents first.
func activateOfflineTool(root string, item stagedOfflineTool) error {
    idDir := filepath.Join(root, item.spec.ID)
    if err := os.MkdirAll(idDir, 0755); err != nil { return err }
    current := managedToolSlot(root, item.spec, "current")
    previous := managedToolSlot(root, item.spec, ".previous")
    backup := filepath.Join(idDir, ".offline-replaced-previous")
    for _, slot := range []string{"current", ".previous"} {
        if _, err := managedToolSlotReceipt(root, item.spec, slot); err != nil { return err }
    }
    if _, err := os.Lstat(backup); err == nil {
        return errors.New("incomplete earlier offline activation; inspect .offline-replaced-previous")
    } else if !errors.Is(err, os.ErrNotExist) { return err }
    _, currentErr := os.Lstat(current)
    if currentErr != nil && !errors.Is(currentErr, os.ErrNotExist) { return currentErr }
    hadCurrent := currentErr == nil
    if !hadCurrent {
        // A missing current version does not invalidate an older saved backup.
        if err := os.Rename(item.ready, current); err != nil { return err }
        return nil
    }
    hadPrev := false
    if _, err := os.Lstat(previous); err == nil {
        if err := os.Rename(previous, backup); err != nil { return err }
        hadPrev = true
    } else if !errors.Is(err, os.ErrNotExist) { return err }
    if err := os.Rename(current, previous); err != nil {
        if hadPrev { _ = os.Rename(backup, previous) }
        return err
    }
    if err := os.Rename(item.ready, current); err != nil {
        _ = os.Rename(previous, current)
        if hadPrev { _ = os.Rename(backup, previous) }
        return fmt.Errorf("failed to activate %s (is the EXE running?): %w", item.spec.ID, err)
    }
    if hadPrev {
        // The oldest managed executable is replaced by the formerly current
        // one. Any unrelated user files are kept in a separate hidden folder.
        _ = os.Remove(filepath.Join(backup, item.spec.Binary+".exe"))
        _ = os.Remove(filepath.Join(backup, "receipt.json"))
        if err := os.Remove(backup); err != nil {
            preserved := filepath.Join(idDir, ".offline-preserved-"+time.Now().UTC().Format("20060102150405.000000000"))
            _ = os.Rename(backup, preserved)
        }
    }
    return nil
}

func importOfflineToolBundle(root, archivePath string) ([]string, error) {
    externalToolInstallMu.Lock()
    defer externalToolInstallMu.Unlock()
    stage, staged, err := stageOfflineToolBundle(root, archivePath)
    if err != nil { return nil, err }
    defer os.RemoveAll(stage)
    installed := make([]string, 0, len(staged))
    for _, item := range staged {
        if err := activateOfflineTool(root, item); err != nil {
            return installed, fmt.Errorf("%s: %w", item.spec.ID, err)
        }
        installed = append(installed, item.spec.ID)
    }
    return installed, nil
}

func (s *Server) exportOfflineTools(w http.ResponseWriter, r *http.Request) {
    root, err := managedExternalToolRoot()
    if err != nil { writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()}); return }
    if err := os.MkdirAll(root, 0755); err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()}); return
    }
    out, err := os.CreateTemp(root, ".offline-export-*.zip")
    if err != nil { writeJSON(w, 500, map[string]any{"error": err.Error()}); return }
    defer func() { out.Close(); os.Remove(out.Name()) }()
    if _, err := exportOfflineToolBundle(root, out); err != nil {
        writeJSON(w, http.StatusConflict, map[string]any{"error": err.Error()}); return
    }
    if _, err := out.Seek(0, io.SeekStart); err != nil {
        writeJSON(w, 500, map[string]any{"error": err.Error()}); return
    }
    w.Header().Set("Content-Type", "application/zip")
    w.Header().Set("Content-Disposition", "attachment; filename=ARTEX-Tools-Offline-windows-amd64.zip")
    http.ServeContent(w, r, "ARTEX-Tools-Offline-windows-amd64.zip", time.Now(), out)
}

func (s *Server) importOfflineTools(w http.ResponseWriter, r *http.Request) {
    root, err := managedExternalToolRoot()
    if err != nil { writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()}); return }
    if err := os.MkdirAll(root, 0755); err != nil {
        writeJSON(w, 500, map[string]any{"error": err.Error()}); return
    }
    r.Body = http.MaxBytesReader(w, r.Body, offlineBundleMaxBytes)
    out, err := os.CreateTemp(root, ".offline-upload-*.zip")
    if err != nil { writeJSON(w, 500, map[string]any{"error": err.Error()}); return }
    defer func() { out.Close(); os.Remove(out.Name()) }()
    n, err := io.Copy(out, r.Body)
    if err != nil {
        var tooLarge *http.MaxBytesError
        if errors.As(err, &tooLarge) {
            writeJSON(w, http.StatusRequestEntityTooLarge, map[string]any{"error": "offline bundle exceeds 1 GiB"})
        } else {
            writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
        }
        return
    }
    if n == 0 { writeJSON(w, 400, map[string]any{"error": "offline bundle is empty"}); return }
    if err := out.Close(); err != nil { writeJSON(w, 500, map[string]any{"error": err.Error()}); return }
    installed, err := importOfflineToolBundle(root, out.Name())
    if err != nil {
        writeJSON(w, http.StatusConflict, map[string]any{"error": err.Error(), "installed": installed})
        return
    }
    writeJSON(w, http.StatusOK, map[string]any{"ok": true, "installed": installed})
}
