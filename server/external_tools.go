package server

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

//go:embed external-tools.json
var externalToolManifestFS embed.FS

type externalToolSpec struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	Binary                 string   `json:"binary"`
	Category               string   `json:"category"`
	Description            string   `json:"description"`
	Source                 string   `json:"source"`
	VersionArgs            []string `json:"version_args"`
	DockerReferenceVersion string   `json:"docker_reference_version,omitempty"`
}

type externalToolManifest struct {
	SchemaVersion int                `json:"schema_version"`
	Tools         []externalToolSpec `json:"tools"`
}

type externalToolStatus struct {
	externalToolSpec
	Status  string `json:"status"` // missing | installed | unverified
	Path    string `json:"path,omitempty"`
	Origin  string `json:"origin,omitempty"` // custom | managed | bundled | path
	Version string `json:"version,omitempty"`
	Detail  string `json:"detail,omitempty"`
	ManagedVersion string `json:"managed_version,omitempty"`
	ManagedPath string `json:"managed_path,omitempty"`
}

type externalToolRoot struct {
	dir    string
	origin string
}

func loadExternalToolManifest() (externalToolManifest, error) {
	data, err := externalToolManifestFS.ReadFile("external-tools.json")
	if err != nil {
		return externalToolManifest{}, err
	}
	var manifest externalToolManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return externalToolManifest{}, err
	}
	if manifest.SchemaVersion != 1 || len(manifest.Tools) == 0 {
		return externalToolManifest{}, fmt.Errorf("unsupported or empty external tool manifest")
	}
	return manifest, nil
}

func externalToolBinary(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

// Only the known executable names from the embedded catalog are looked up.
// User-provided directories are opt-in via ARTEX_TOOL_DIRS (OS path-list syntax).
func externalToolRoots() []externalToolRoot {
	roots := []externalToolRoot{}
	for _, dir := range filepath.SplitList(os.Getenv("ARTEX_TOOL_DIRS")) {
		if strings.TrimSpace(dir) != "" {
			roots = append(roots, externalToolRoot{dir: dir, origin: "custom"})
		}
	}
	if local := os.Getenv("LOCALAPPDATA"); local != "" {
		roots = append(roots, externalToolRoot{dir: filepath.Join(local, "ARTEX", "tools"), origin: "managed"})
	}
	if exe, err := os.Executable(); err == nil {
		base := filepath.Dir(exe)
		roots = append(roots,
			externalToolRoot{dir: filepath.Join(base, "data", "tools"), origin: "bundled"},
			externalToolRoot{dir: filepath.Join(base, "tools"), origin: "bundled"})
	}
	return roots
}

func externalToolExecutable(spec externalToolSpec, roots []externalToolRoot) (string, string) {
	binary := externalToolBinary(spec.Binary)
	for _, root := range roots {
		candidates := []string{filepath.Join(root.dir, binary), filepath.Join(root.dir, spec.ID, binary), filepath.Join(root.dir, spec.ID, "current", binary)}
		// Versioned Tool Manager layout: tools/<id>/<version>/<binary>.
		versioned, _ := filepath.Glob(filepath.Join(root.dir, spec.ID, "*", binary))
		sort.Sort(sort.Reverse(sort.StringSlice(versioned)))
		candidates = append(candidates, versioned...)
		for _, candidate := range candidates {
			fi, err := os.Stat(candidate)
			if err == nil && fi.Mode().IsRegular() && (runtime.GOOS == "windows" || fi.Mode()&0111 != 0) {
				path, err := filepath.Abs(candidate)
				if err == nil {
					return path, root.origin
				}
			}
		}
	}
	if path, err := exec.LookPath(binary); err == nil {
		return path, "path"
	}
	return "", ""
}

// Version checks intentionally pass only the catalog's fixed --version-style
// arguments. No targets are supplied; this is inventory, not a scan.
func probeExternalTool(ctx context.Context, spec externalToolSpec, roots []externalToolRoot) externalToolStatus {
	status := externalToolStatus{externalToolSpec: spec, Status: "missing"}
	if local := os.Getenv("LOCALAPPDATA"); local != "" {
		root := filepath.Join(local, "ARTEX", "tools")
		if receipt := externalToolManagedReceipt(root, spec); receipt != nil {
			status.ManagedVersion = receipt.Version
			status.ManagedPath = filepath.Join(root, spec.ID, "current", externalToolBinary(spec.Binary))
		}
	}
	path, origin := externalToolExecutable(spec, roots)
	if path == "" {
		return status
	}
	status.Path, status.Origin = path, origin
	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(timeout, path, spec.VersionArgs...)
	var output cappedToolOutput
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()
	raw := strings.TrimSpace(output.String())
	if err != nil {
		status.Status = "unverified"
		if timeout.Err() != nil {
			status.Detail = "版本检查超时"
		} else {
			status.Detail = "版本检查失败：" + err.Error()
		}
		if len(raw) > 0 {
			status.Version = firstToolLine(raw)
		}
		return status
	}
	status.Status = "installed"
	status.Version = firstToolLine(raw)
	if status.Version == "" {
		status.Detail = "命令运行成功，但没有返回版本"
	}
	return status
}

type cappedToolOutput struct {
	bytes.Buffer
}

func (w *cappedToolOutput) Write(data []byte) (int, error) {
	original := len(data)
	if remaining := 4096 - w.Len(); remaining > 0 {
		if len(data) > remaining {
			data = data[:remaining]
		}
		_, _ = w.Buffer.Write(data)
	}
	return original, nil
}

func firstToolLine(raw string) string {
	line := strings.SplitN(strings.ReplaceAll(raw, "\r", ""), "\n", 2)[0]
	if len(line) > 180 {
		line = line[:180]
	}
	return line
}

func externalToolInventory(ctx context.Context, tools []externalToolSpec, roots []externalToolRoot) []externalToolStatus {
	items := make([]externalToolStatus, len(tools))
	var wg sync.WaitGroup
	for i := range tools {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			items[i] = probeExternalTool(ctx, tools[i], roots)
		}(i)
	}
	wg.Wait()
	return items
}

func (s *Server) listExternalTools(w http.ResponseWriter, r *http.Request) {
	manifest, err := loadExternalToolManifest()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"tools": externalToolInventory(r.Context(), manifest.Tools, externalToolRoots()),
		"host":  runtime.GOOS + "/" + runtime.GOARCH,
	})
}
