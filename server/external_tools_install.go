package server

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	externalToolAPIBase = "https://api.github.com"
	externalToolDownloadBase = "https://github.com"
	maxExternalToolArchive = 256 << 20
	maxExternalToolBinary = 350 << 20
)

var (
	externalToolInstallMu sync.Mutex
	externalToolReleaseTag = regexp.MustCompile(`^v?[0-9]+\.[0-9]+\.[0-9]+(?:[-+][a-zA-Z0-9.]+)?$`)
	externalToolDigest = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
)

type externalToolRelease struct {
	ID string `json:"id"`
	Version string `json:"version"`
	Asset string `json:"asset"`
	SizeBytes int64 `json:"size_bytes"`
	SHA256 string `json:"sha256"`
	Source string `json:"source"`
	downloadURL string
}

type managedToolReceipt struct {
	Version string `json:"version"`
	SHA256 string `json:"sha256"`
	InstalledAt string `json:"installed_at"`
}

func externalToolSpecByID(id string) (externalToolSpec, bool) {
	manifest, err := loadExternalToolManifest()
	if err != nil { return externalToolSpec{}, false }
	for _, spec := range manifest.Tools {
		if spec.ID == id { return spec, true }
	}
	return externalToolSpec{}, false
}

func externalToolRepo(spec externalToolSpec) (string, error) {
	repo := strings.TrimPrefix(spec.Source, "https://github.com/")
	parts := strings.Split(repo, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" ||
		strings.ContainsAny(repo, "?#\\ \t\n") {
		return "", fmt.Errorf("unsupported official source for %s", spec.ID)
	}
	return repo, nil
}

// The GitHub release API reports the SHA-256 of each uploaded asset. This
// verifies bytes against that API record, but is NOT a publisher signature.
func fetchExternalToolRelease(ctx context.Context, client *http.Client, apiBase, downloadBase string, spec externalToolSpec) (externalToolRelease, error) {
	repo, err := externalToolRepo(spec)
	if err != nil { return externalToolRelease{}, err }
	url := strings.TrimRight(apiBase, "/") + "/repos/" + repo + "/releases/latest"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil { return externalToolRelease{}, err }
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "ARTEX-tool-manager")
	resp, err := client.Do(req)
	if err != nil { return externalToolRelease{}, fmt.Errorf("GitHub release query: %w", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return externalToolRelease{}, fmt.Errorf("GitHub release query: HTTP %d", resp.StatusCode)
	}
	var release struct {
		Tag string `json:"tag_name"`
		Assets []struct {
			Name string `json:"name"`
			Size int64 `json:"size"`
			Digest string `json:"digest"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 2<<20)).Decode(&release); err != nil {
		return externalToolRelease{}, fmt.Errorf("GitHub release payload: %w", err)
	}
	if !externalToolReleaseTag.MatchString(release.Tag) {
		return externalToolRelease{}, fmt.Errorf("invalid official release version %q", release.Tag)
	}
	filename := spec.Binary + "_" + strings.TrimPrefix(release.Tag, "v") + "_windows_amd64.zip"
	for _, asset := range release.Assets {
		if asset.Name != filename { continue }
		digest := strings.TrimPrefix(asset.Digest, "sha256:")
		if !strings.HasPrefix(asset.Digest, "sha256:") || !externalToolDigest.MatchString(digest) {
			return externalToolRelease{}, errors.New("official release asset does not provide a valid SHA-256 digest")
		}
		if asset.Size <= 0 || asset.Size > maxExternalToolArchive {
			return externalToolRelease{}, fmt.Errorf("official archive size unsupported: %d bytes", asset.Size)
		}
		return externalToolRelease{
			ID: spec.ID, Version: release.Tag, Asset: filename,
			SizeBytes: asset.Size, SHA256: strings.ToLower(digest),
			Source: spec.Source + "/releases/tag/" + release.Tag,
			downloadURL: strings.TrimRight(downloadBase, "/") + "/" + repo +
				"/releases/download/" + release.Tag + "/" + filename,
		}, nil
	}
	return externalToolRelease{}, fmt.Errorf("official release %s has no %s asset", release.Tag, filename)
}

func managedExternalToolRoot() (string, error) {
	if runtime.GOOS != "windows" || runtime.GOARCH != "amd64" {
		return "", errors.New("当前仅支持 Windows x64 工具安装；Linux 和其他架构仍可手动配置工具")
	}
	base := os.Getenv("LOCALAPPDATA")
	if base == "" { return "", errors.New("LOCALAPPDATA is not configured") }
	return filepath.Join(base, "ARTEX", "tools"), nil
}

func externalToolManagedReceipt(root string, spec externalToolSpec) *managedToolReceipt {
	data, err := os.ReadFile(filepath.Join(root, spec.ID, "current", "receipt.json"))
	if err != nil { return nil }
	var receipt managedToolReceipt
	if json.Unmarshal(data, &receipt) != nil { return nil }
	if !externalToolReleaseTag.MatchString(receipt.Version) { return nil }
	return &receipt
}

func installVerifiedExternalTool(ctx context.Context, client *http.Client, spec externalToolSpec, asset externalToolRelease, root string) (string, error) {
	externalToolInstallMu.Lock()
	defer externalToolInstallMu.Unlock()

	idDir := filepath.Join(root, spec.ID)
	if err := os.MkdirAll(idDir, 0755); err != nil { return "", err }
	stage, err := os.MkdirTemp(idDir, ".install-")
	if err != nil { return "", err }
	defer os.RemoveAll(stage)
	archiveFile := filepath.Join(stage, "package.zip")
	out, err := os.Create(archiveFile)
	if err != nil { return "", err }

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.downloadURL, nil)
	if err != nil { out.Close(); return "", err }
	req.Header.Set("User-Agent", "ARTEX-tool-manager")
	resp, err := client.Do(req)
	if err != nil { out.Close(); return "", fmt.Errorf("download: %w", err) }
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close(); out.Close()
		return "", fmt.Errorf("official download: HTTP %d", resp.StatusCode)
	}
	h := sha256.New()
	n, copyErr := io.Copy(io.MultiWriter(out, h), io.LimitReader(resp.Body, maxExternalToolArchive+1))
	resp.Body.Close()
	closeErr := out.Close()
	if copyErr != nil { return "", fmt.Errorf("download: %w", copyErr) }
	if closeErr != nil { return "", closeErr }
	if n != asset.SizeBytes || n > maxExternalToolArchive {
		return "", fmt.Errorf("archive size mismatch: expected %d, got %d", asset.SizeBytes, n)
	}
	if checksum := hex.EncodeToString(h.Sum(nil)); checksum != asset.SHA256 {
		return "", fmt.Errorf("SHA-256 mismatch: expected %s, got %s", asset.SHA256, checksum)
	}

	zr, err := zip.OpenReader(archiveFile)
	if err != nil { return "", fmt.Errorf("invalid official archive: %w", err) }
	defer zr.Close()
	binaryName := spec.Binary + ".exe"
	var binary *zip.File
	for _, entry := range zr.File {
		if path.Base(entry.Name) != binaryName { continue }
		if binary != nil || entry.FileInfo().Mode()&os.ModeSymlink != 0 || entry.FileInfo().IsDir() ||
			entry.UncompressedSize64 > maxExternalToolBinary {
			return "", errors.New("archive contains an unsafe or duplicate executable entry")
		}
		binary = entry
	}
	if binary == nil { return "", fmt.Errorf("archive is missing %s", binaryName) }
	installDir := filepath.Join(stage, "ready")
	if err := os.MkdirAll(installDir, 0755); err != nil { return "", err }
	rc, err := binary.Open()
	if err != nil { return "", err }
	dst, err := os.OpenFile(filepath.Join(installDir, binaryName), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0755)
	if err != nil { rc.Close(); return "", err }
	_, copyErr = io.Copy(dst, io.LimitReader(rc, maxExternalToolBinary+1))
	rc.Close()
	closeErr = dst.Close()
	if copyErr != nil { return "", copyErr }
	if closeErr != nil { return "", closeErr }
	if stat, err := os.Stat(filepath.Join(installDir, binaryName)); err != nil || stat.Size() <= 0 || stat.Size() > maxExternalToolBinary {
		return "", errors.New("extracted executable has invalid size")
	}
	receipt := managedToolReceipt{Version: asset.Version, SHA256: asset.SHA256, InstalledAt: time.Now().UTC().Format(time.RFC3339)}
	contents, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil { return "", err }
	if err := os.WriteFile(filepath.Join(installDir, "receipt.json"), contents, 0644); err != nil { return "", err }

	current := filepath.Join(idDir, "current")
	previous := filepath.Join(idDir, ".previous")
	if err := os.RemoveAll(previous); err != nil { return "", err }
	hadCurrent := false
	if _, err := os.Stat(current); err == nil {
		if err := os.Rename(current, previous); err != nil {
			return "", fmt.Errorf("cannot replace running tool; close it and retry: %w", err)
		}
		hadCurrent = true
	} else if !os.IsNotExist(err) { return "", err }
	if err := os.Rename(installDir, current); err != nil {
		if hadCurrent {
			_ = os.Rename(previous, current) // restore previous working version
		}
		return "", fmt.Errorf("install failed: %w", err)
	}
	return filepath.Join(current, binaryName), nil
}

func (s *Server) externalToolLatestRelease(w http.ResponseWriter, r *http.Request) {
	spec, found := externalToolSpecByID(r.PathValue("id"))
	if !found { writeJSON(w, 404, map[string]any{"error": "unknown external tool"}); return }
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	asset, err := fetchExternalToolRelease(ctx, &http.Client{Timeout: 15*time.Second}, externalToolAPIBase, externalToolDownloadBase, spec)
	if err != nil { writeJSON(w, 502, map[string]any{"error": err.Error()}); return }
	writeJSON(w, 200, asset)
}

func (s *Server) externalToolInstall(w http.ResponseWriter, r *http.Request) {
	spec, found := externalToolSpecByID(r.PathValue("id"))
	if !found { writeJSON(w, 404, map[string]any{"error": "unknown external tool"}); return }
	root, err := managedExternalToolRoot()
	if err != nil { writeJSON(w, 400, map[string]any{"error": err.Error()}); return }
	var request struct {
		Version string `json:"version"`
		SHA256 string `json:"sha256"`
	}
	if json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&request) != nil ||
		request.Version == "" || !externalToolDigest.MatchString(request.SHA256) {
		writeJSON(w, 400, map[string]any{"error": "version and SHA-256 from the release preview are required"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()
	client := &http.Client{Timeout: 4*time.Minute}
	asset, err := fetchExternalToolRelease(ctx, client, externalToolAPIBase, externalToolDownloadBase, spec)
	if err != nil { writeJSON(w, 502, map[string]any{"error": err.Error()}); return }
	if asset.Version != request.Version || asset.SHA256 != strings.ToLower(request.SHA256) {
		writeJSON(w, 409, map[string]any{"error": "official release changed; check the latest version again"})
		return
	}
	path, err := installVerifiedExternalTool(ctx, client, spec, asset, root)
	if err != nil { writeJSON(w, 502, map[string]any{"error": err.Error()}); return }
	writeJSON(w, 200, map[string]any{"ok": true, "version": asset.Version, "sha256": asset.SHA256, "path": path})
}
