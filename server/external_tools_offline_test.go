package server

import (
    "archive/zip"
    "bytes"
    "encoding/json"
    "io"
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func offlineToolFixture(t *testing.T, root string, spec externalToolSpec, slot, version, body string) {
    t.Helper()
    platformBinary := writeManagedToolFixture(t, root, spec, slot, version, body)
    windowsBinary := filepath.Join(filepath.Dir(platformBinary), spec.Binary+".exe")
    if windowsBinary != platformBinary {
        if err := os.WriteFile(windowsBinary, []byte(body), 0755); err != nil { t.Fatal(err) }
    }
}

func offlineExportFixture(t *testing.T, root string) []byte {
    t.Helper()
    var output bytes.Buffer
    count, err := exportOfflineToolBundle(root, &output)
    if err != nil || count < 1 { t.Fatalf("offline export: %d tools: %v", count, err) }
    return output.Bytes()
}

func saveOfflineArchive(t *testing.T, data []byte) string {
    t.Helper()
    file := filepath.Join(t.TempDir(), "ARTEX-Tools-Offline.zip")
    if err := os.WriteFile(file, data, 0644); err != nil { t.Fatal(err) }
    return file
}

func TestExternalToolOfflineRoundTrip(t *testing.T) {
    source := t.TempDir()
    nuclei := externalToolSpec{ID: "nuclei", Binary: "nuclei"}
    ffuf := externalToolSpec{ID: "ffuf", Binary: "ffuf"}
    offlineToolFixture(t, source, nuclei, "current", "v3.11.1", "nuclei exported")
    offlineToolFixture(t, source, ffuf, "current", "v2.3.0", "ffuf exported")
    archive := offlineExportFixture(t, source)
    zr, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
    if err != nil { t.Fatal(err) }
    if len(zr.File) != 3 { t.Fatalf("expected manifest plus 2 tools, got %d", len(zr.File)) }
    var manifest offlineBundleManifest
    for _, f := range zr.File {
        if f.Name != "manifest.json" { continue }
        reader, err := f.Open()
        if err != nil { t.Fatal(err) }
        if err := json.NewDecoder(reader).Decode(&manifest); err != nil { t.Fatal(err) }
        reader.Close()
    }
    if manifest.SchemaVersion != 1 || manifest.Target != "windows/amd64" || len(manifest.Tools) != 2 {
        t.Fatalf("unexpected manifest: %+v", manifest)
    }
    target := t.TempDir()
    offlineToolFixture(t, target, nuclei, "current", "v3.10.0", "nuclei existing")
    offlineToolFixture(t, target, nuclei, ".previous", "v3.9.0", "nuclei older")
    olderNote := filepath.Join(target, nuclei.ID, ".previous", "my-notes.txt")
    if err := os.WriteFile(olderNote, []byte("personal"), 0644); err != nil { t.Fatal(err) }
    custom := filepath.Join(target, "user-tools", "nuclei.exe")
    if err := os.MkdirAll(filepath.Dir(custom), 0755); err != nil { t.Fatal(err) }
    if err := os.WriteFile(custom, []byte("user binary"), 0755); err != nil { t.Fatal(err) }

    installed, err := importOfflineToolBundle(target, saveOfflineArchive(t, archive))
    if err != nil || len(installed) != 2 { t.Fatalf("offline import: %v (%v)", installed, err) }
    if got := readToolContents(t, filepath.Join(target, nuclei.ID, "current", "nuclei.exe")); got != "nuclei exported" {
        t.Fatalf("new nuclei: %q", got)
    }
    if got := readToolContents(t, filepath.Join(target, nuclei.ID, ".previous", externalToolBinary(nuclei.Binary))); got != "nuclei existing" {
        t.Fatalf("previous nuclei: %q", got)
    }
    if got := readToolContents(t, filepath.Join(target, ffuf.ID, "current", "ffuf.exe")); got != "ffuf exported" {
        t.Fatalf("new ffuf: %q", got)
    }
    if got := readToolContents(t, custom); got != "user binary" { t.Fatalf("custom changed: %q", got) }
    preserved, err := filepath.Glob(filepath.Join(target, nuclei.ID, ".offline-preserved-*", "my-notes.txt"))
    if err != nil || len(preserved) != 1 || readToolContents(t, preserved[0]) != "personal" {
        t.Fatalf("old personal notes not preserved: %v, %v", preserved, err)
    }
}

func TestExternalToolOfflineRejectTamperingBeforeAnyInstall(t *testing.T) {
    source := t.TempDir()
    for _, spec := range []externalToolSpec{{ID: "nuclei", Binary: "nuclei"}, {ID: "ffuf", Binary: "ffuf"}} {
        offlineToolFixture(t, source, spec, "current", "v1.2.3", "trusted "+spec.ID)
    }
    original := offlineExportFixture(t, source)
    reader, err := zip.NewReader(bytes.NewReader(original), int64(len(original)))
    if err != nil { t.Fatal(err) }
    var modified bytes.Buffer
    writer := zip.NewWriter(&modified)
    for _, file := range reader.File {
        src, err := file.Open()
        if err != nil { t.Fatal(err) }
        data, err := io.ReadAll(src)
        src.Close()
        if err != nil { t.Fatal(err) }
        if strings.HasSuffix(file.Name, "/ffuf.exe") { data = []byte("altered ffuf") }
        dest, err := writer.Create(file.Name)
        if err != nil { t.Fatal(err) }
        if _, err := dest.Write(data); err != nil { t.Fatal(err) }
    }
    if err := writer.Close(); err != nil { t.Fatal(err) }
    target := t.TempDir()
    spec := externalToolSpec{ID: "nuclei", Binary: "nuclei"}
    offlineToolFixture(t, target, spec, "current", "v0.9.0", "keep old")
    installed, err := importOfflineToolBundle(target, saveOfflineArchive(t, modified.Bytes()))
    if err == nil || !strings.Contains(err.Error(), "SHA-256 mismatch") || len(installed) != 0 {
        t.Fatalf("altered executable accepted: %v, %v", installed, err)
    }
    current := filepath.Join(target, "nuclei", "current", externalToolBinary("nuclei"))
    if got := readToolContents(t, current); got != "keep old" { t.Fatalf("previous installation replaced: %q", got) }
    if _, err := os.Stat(filepath.Join(target, "ffuf", "current", "ffuf.exe")); !os.IsNotExist(err) {
        t.Fatalf("other tool installed before validation: %v", err)
    }
}

func TestExternalToolOfflineRejectUnknownMember(t *testing.T) {
    source := t.TempDir()
    offlineToolFixture(t, source, externalToolSpec{ID: "nuclei", Binary: "nuclei"}, "current", "v1.2.3", "nuclei")
    valid := offlineExportFixture(t, source)
    reader, err := zip.NewReader(bytes.NewReader(valid), int64(len(valid)))
    if err != nil { t.Fatal(err) }
    var modified bytes.Buffer
    writer := zip.NewWriter(&modified)
    for _, file := range reader.File {
        src, err := file.Open()
        if err != nil { t.Fatal(err) }
        out, err := writer.Create(file.Name)
        if err != nil { t.Fatal(err) }
        if _, err := io.Copy(out, src); err != nil { t.Fatal(err) }
        src.Close()
    }
    bogus, err := writer.Create("../unwanted.txt")
    if err != nil { t.Fatal(err) }
    if _, err := bogus.Write([]byte("unwanted")); err != nil { t.Fatal(err) }
    if err := writer.Close(); err != nil { t.Fatal(err) }
    target := t.TempDir()
    if _, err := importOfflineToolBundle(target, saveOfflineArchive(t, modified.Bytes())); err == nil {
        t.Fatal("bundle with extra member accepted")
    }
    if _, err := os.Stat(filepath.Join(target, "nuclei", "current")); !os.IsNotExist(err) {
        t.Fatalf("unexpected install after rejected bundle: %v", err)
    }
}

func TestExternalToolOfflineEmptyExport(t *testing.T) {
    var out bytes.Buffer
    if _, err := exportOfflineToolBundle(t.TempDir(), &out); err == nil {
        t.Fatal("empty managed tool set should not produce bundle")
    }
}

func TestExternalToolOfflineKeepPreviousIfCurrentMissing(t *testing.T) {
    source := t.TempDir()
    spec := externalToolSpec{ID: "dnsx", Binary: "dnsx"}
    offlineToolFixture(t, source, spec, "current", "v1.2.3", "new dnsx")
    target := t.TempDir()
    offlineToolFixture(t, target, spec, ".previous", "v1.1.0", "old dnsx")
    _, err := importOfflineToolBundle(target, saveOfflineArchive(t, offlineExportFixture(t, source)))
    if err != nil { t.Fatal(err) }
    if got := readToolContents(t, filepath.Join(target, spec.ID, ".previous", externalToolBinary(spec.Binary))); got != "old dnsx" {
        t.Fatalf("previous deleted: %q", got)
    }
}
