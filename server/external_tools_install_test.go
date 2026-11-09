package server

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testExternalZip(t *testing.T, binary string, contents string) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	file, err := writer.Create("release/" + binary + ".exe")
	if err != nil { t.Fatal(err) }
	if _, err = file.Write([]byte(contents)); err != nil { t.Fatal(err) }
	if err = writer.Close(); err != nil { t.Fatal(err) }
	return buffer.Bytes()
}

func TestExternalToolOfficialReleaseAndInstall(t *testing.T) {
	spec := externalToolSpec{ID: "nuclei", Binary: "nuclei", Source: "https://github.com/projectdiscovery/nuclei"}
	version := "v9.8.7"
	archive := testExternalZip(t, "nuclei", "first working release")
	sum := sha256.Sum256(archive)
	digest := hex.EncodeToString(sum[:])
	var requestCount int
	mux := http.NewServeMux()
	mux.HandleFunc("GET /repos/projectdiscovery/nuclei/releases/latest", func(w http.ResponseWriter, _ *http.Request) {
		payload := map[string]any{
			"tag_name": version,
			"assets": []map[string]any{{
				"name": "nuclei_" + strings.TrimPrefix(version, "v") + "_windows_amd64.zip",
				"size": len(archive),
				"digest": "sha256:" + digest,
			}},
		}
		_ = json.NewEncoder(w).Encode(payload)
	})
	mux.HandleFunc("GET /projectdiscovery/nuclei/releases/download/", func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		_, _ = w.Write(archive)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	client := srv.Client()
	root := t.TempDir()

	release, err := fetchExternalToolRelease(context.Background(), client, srv.URL, srv.URL, spec)
	if err != nil { t.Fatal(err) }
	if release.Version != version || release.SHA256 != digest || release.SizeBytes != int64(len(archive)) {
		t.Fatalf("unexpected release: %+v", release)
	}
	target, err := installVerifiedExternalTool(context.Background(), client, spec, release, root)
	if err != nil { t.Fatal(err) }
	if requestCount != 1 { t.Fatalf("unexpected download count: %d", requestCount) }
	contents, err := os.ReadFile(target)
	if err != nil || string(contents) != "first working release" {
		t.Fatalf("installed executable: %q, %v", contents, err)
	}
	receipt := externalToolManagedReceipt(root, spec)
	if receipt == nil || receipt.Version != "v9.8.7" || receipt.SHA256 != digest {
		t.Fatalf("missing install receipt: %+v", receipt)
	}

	// Updating replaces the active version but keeps one prior working version.
	version = "v9.8.8"
	archive = testExternalZip(t, "nuclei", "newer release")
	sum = sha256.Sum256(archive)
	digest = hex.EncodeToString(sum[:])
	updated, err := fetchExternalToolRelease(context.Background(), client, srv.URL, srv.URL, spec)
	if err != nil { t.Fatal(err) }
	_, err = installVerifiedExternalTool(context.Background(), client, spec, updated, root)
	if err != nil { t.Fatal(err) }
	previous, err := os.ReadFile(filepath.Join(root, spec.ID, ".previous", "nuclei.exe"))
	if err != nil || string(previous) != "first working release" {
		t.Fatalf("previous binary not preserved: %q, %v", previous, err)
	}
	current, err := os.ReadFile(filepath.Join(root, spec.ID, "current", "nuclei.exe"))
	if err != nil || string(current) != "newer release" {
		t.Fatalf("updated executable: %q, %v", current, err)
	}
}

func TestExternalToolBadDownloadPreservesInstall(t *testing.T) {
	spec := externalToolSpec{ID: "ffuf", Binary: "ffuf"}
	zipContents := testExternalZip(t, "ffuf", "expected official release")
	hash := sha256.Sum256(zipContents)
	asset := externalToolRelease{
		Version: "v2.3.0", SHA256: hex.EncodeToString(hash[:]),
		SizeBytes: int64(len(zipContents)),
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(bytes.Repeat([]byte("X"), len(zipContents)))
	}))
	defer srv.Close()
	asset.downloadURL = srv.URL
	root := t.TempDir()
	current := filepath.Join(root, "ffuf", "current")
	if err := os.MkdirAll(current, 0755); err != nil { t.Fatal(err) }
	if err := os.WriteFile(filepath.Join(current, "ffuf.exe"), []byte("existing ffuf"), 0755); err != nil { t.Fatal(err) }
	if _, err := installVerifiedExternalTool(context.Background(), srv.Client(), spec, asset, root); err == nil || !strings.Contains(err.Error(), "SHA-256 mismatch") {
		t.Fatalf("expected digest mismatch; got %v", err)
	}
	actual, err := os.ReadFile(filepath.Join(current, "ffuf.exe"))
	if err != nil || string(actual) != "existing ffuf" {
		t.Fatalf("failed install replaced existing tool: %q, %v", actual, err)
	}
}

func TestExternalToolMissingDigestRejected(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"tag_name":"v1.2.3","assets":[{"name":"ffuf_1.2.3_windows_amd64.zip","size":100}]}`)
	}))
	defer srv.Close()
	spec := externalToolSpec{ID: "ffuf", Binary: "ffuf", Source: "https://github.com/ffuf/ffuf"}
	_, err := fetchExternalToolRelease(context.Background(), srv.Client(), srv.URL, srv.URL, spec)
	if err == nil || !strings.Contains(err.Error(), "SHA-256") {
		t.Fatalf("expected missing digest rejection: %v", err)
	}
}

func TestExternalToolInvalidSourceRejected(t *testing.T) {
	for _, source := range []string{"https://example.org/team/ffuf", "https://github.com/ffuf/ffuf/extra"} {
		if _, err := externalToolRepo(externalToolSpec{ID: "ffuf", Source: source}); err == nil {
			t.Fatalf("accepted non-catalog source %q", source)
		}
	}
}
