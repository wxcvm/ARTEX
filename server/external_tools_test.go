package server

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExternalToolManifest(t *testing.T) {
	manifest, err := loadExternalToolManifest()
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Tools) != 7 {
		t.Fatalf("want 7 recon entries, got %d", len(manifest.Tools))
	}
	seen := map[string]bool{}
	for _, tool := range manifest.Tools {
		if seen[tool.ID] || tool.ID == "" || tool.Binary == "" || len(tool.VersionArgs) == 0 {
			t.Fatalf("invalid tool entry: %+v", tool)
		}
		if !strings.HasPrefix(tool.Source, "https://github.com/") {
			t.Fatalf("tool %s has no official GitHub source: %s", tool.ID, tool.Source)
		}
		seen[tool.ID] = true
	}
}

func TestExternalToolMissing(t *testing.T) {
	spec := externalToolSpec{ID: "artex-unlikely-probe-name", Binary: "artex-unlikely-probe-name", VersionArgs: []string{"-version"}}
	result := probeExternalTool(context.Background(), spec, []externalToolRoot{{dir: t.TempDir(), origin: "custom"}})
	if result.Status != "missing" || result.Path != "" {
		t.Fatalf("missing tool incorrectly detected: %+v", result)
	}
}

func TestExternalToolVersionHelper(t *testing.T) {
	if os.Getenv("ARTEX_PROBE_HELPER") != "1" {
		return
	}
	fmt.Println("example-recon version v9.8.7")
	os.Exit(0)
}

func TestExternalToolVersionedDirectory(t *testing.T) {
	root := t.TempDir()
	binary := externalToolBinary("example-recon")
	versionDir := filepath.Join(root, "example-recon", "v9.8.7")
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		t.Fatal(err)
	}
	// Use the Go test binary itself so the probe tests real process execution
	// on Windows and Linux without requiring any third-party tool downloads.
	testBinary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	src, err := os.Open(testBinary)
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()
	target := filepath.Join(versionDir, binary)
	dst, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		t.Fatal(err)
	}
	if err := dst.Close(); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ARTEX_PROBE_HELPER", "1")
	spec := externalToolSpec{ID: "example-recon", Binary: "example-recon", VersionArgs: []string{"-test.run=TestExternalToolVersionHelper"}}
	status := probeExternalTool(context.Background(), spec, []externalToolRoot{{dir: root, origin: "managed"}})
	if status.Status != "installed" || status.Origin != "managed" ||
		!strings.Contains(status.Version, "v9.8.7") || status.Path != target {
		t.Fatalf("unexpected status: %+v", status)
	}
}

func TestExternalToolPreferExplicitDirectory(t *testing.T) {
	base := t.TempDir()
	custom := filepath.Join(base, "custom")
	managed := filepath.Join(base, "managed")
	if err := os.MkdirAll(custom, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(managed, 0755); err != nil {
		t.Fatal(err)
	}
	bin := externalToolBinary("example-recon")
	for _, dir := range []string{custom, managed} {
		if err := os.WriteFile(filepath.Join(dir, bin), []byte("not a real executable"), 0755); err != nil {
			t.Fatal(err)
		}
	}
	path, origin := externalToolExecutable(externalToolSpec{ID: "example-recon", Binary: "example-recon"}, []externalToolRoot{
		{dir: custom, origin: "custom"},
		{dir: managed, origin: "managed"},
	})
	if origin != "custom" || path != filepath.Join(custom, bin) {
		t.Fatalf("explicit directory lost priority: %s (%s)", path, origin)
	}
}

func TestFirstToolLine(t *testing.T) {
	if got := firstToolLine("nuclei 3.11.1\r\nmore"); got != "nuclei 3.11.1" {
		t.Fatalf("version extraction: %q", got)
	}
}
