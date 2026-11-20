package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

const releaseAPI = "https://api.github.com/repos/wxcvm/ARTEX/releases/latest"
const releasePage = "https://github.com/wxcvm/ARTEX/releases"

// releaseAPIFallbacks are tried in order AFTER the official endpoint (and after
// the configured global proxy) when the official query fails. Users behind
// restricted networks hit api.github.com with 403 rate-limit or TLS reset and
// previously got a bare 502 with no recovery. These read-only mirrors serve the
// same JSON for the same path, so the shape check in fetchStableRelease still
// applies. The first mirror is the release-check default; the full list (and
// the primary endpoint) can be overridden with ARTEX_RELEASE_API_MIRRORS as a
// comma-separated list of API base URLs.
var releaseAPIFallbacks = []string{
	"https://gh-proxy.com/https://api.github.com",
	"https://ghproxy.net/https://api.github.com",
	"https://mirror.ghproxy.com/https://api.github.com",
}

type releaseUpdateAsset struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}
type releaseUpdate struct {
	Current     string               `json:"current"`
	Latest      string               `json:"latest"`
	Comparable  bool                 `json:"comparable"`
	HasUpdate   bool                 `json:"has_update"`
	ReleaseURL  string               `json:"release_url"`
	PublishedAt string               `json:"published_at"`
	Assets      []releaseUpdateAsset `json:"assets"`
	Source      string               `json:"source"`
	Note        string               `json:"note,omitempty"`
}
type githubRelease struct {
	Tag         string `json:"tag_name"`
	PublishedAt string `json:"published_at"`
	Draft       bool   `json:"draft"`
	Prerelease  bool   `json:"prerelease"`
	Assets      []struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
	} `json:"assets"`
}

func releaseVersion(raw string) string {
	raw = strings.TrimSpace(raw)
	if !strings.HasPrefix(raw, "v") {
		raw = "v" + raw
	}
	if !semver.IsValid(raw) {
		return ""
	}
	// Go semver also accepts shortened v1 / v1.2; ARTEX publishes full tags.
	core := strings.SplitN(strings.SplitN(raw, "-", 2)[0], "+", 2)[0]
	if strings.Count(core, ".") != 2 {
		return ""
	}
	return raw
}

func describeRelease(current, goos, goarch string, rel githubRelease) (releaseUpdate, error) {
	out := releaseUpdate{Current: current, Latest: rel.Tag, Assets: []releaseUpdateAsset{}, PublishedAt: rel.PublishedAt}
	latest := releaseVersion(rel.Tag)
	if latest == "" || rel.Draft || rel.Prerelease || semver.Prerelease(latest) != "" {
		return out, errors.New("invalid stable release")
	}
	out.ReleaseURL = releasePage + "/tag/" + url.PathEscape(rel.Tag)
	if cur := releaseVersion(current); cur != "" {
		out.Comparable = true
		out.HasUpdate = semver.Compare(cur, latest) < 0
	}
	version := strings.TrimPrefix(latest, "v")
	names := []string{"artex-" + version + "-" + goos + "-" + goarch + ".zip"}
	if goos == "windows" && goarch == "amd64" {
		names = append([]string{"ARTEX-Windows-Setup-amd64.exe", "ARTEX-Windows-Desktop-amd64.zip"}, names...)
	}
	names = append(names, "SHA256SUMS")
	for _, name := range names {
		for _, asset := range rel.Assets {
			if asset.Name == name && asset.Size > 0 {
				out.Assets = append(out.Assets, releaseUpdateAsset{Name: name, Size: asset.Size, URL: releasePage + "/download/" + url.PathEscape(rel.Tag) + "/" + url.PathEscape(name)})
				break
			}
		}
	}
	return out, nil
}

type releaseEndpointsKey struct{}

// releaseAPIEndpoints returns the endpoints to probe, in order: the override
// list from ARTEX_RELEASE_API_MIRRORS when set, otherwise the official
// api.github.com followed by the default mirror chain. Tests can inject local
// servers via the releaseEndpointsKey context value.
func releaseAPIEndpoints() []string {
	if raw := strings.TrimSpace(os.Getenv("ARTEX_RELEASE_API_MIRRORS")); raw != "" {
		var out []string
		for _, part := range strings.Split(raw, ",") {
			if base := strings.TrimRight(strings.TrimSpace(part), "/"); base != "" {
				out = append(out, base)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	return append([]string{"https://api.github.com"}, releaseAPIFallbacks...)
}

func releaseAPIEndpointsFor(ctx context.Context) []string {
	if v, ok := ctx.Value(releaseEndpointsKey{}).([]string); ok && len(v) > 0 {
		return v
	}
	return releaseAPIEndpoints()
}

// fetchStableRelease probes endpoints in order and returns the first
// well-formed release payload. Every failure is collected so the caller can
// surface WHY the check failed instead of a bare 502.
func fetchStableRelease(ctx context.Context, client *http.Client, endpoint string) (githubRelease, error) {
	var rel githubRelease
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return rel, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "ARTEX-release-check")
	resp, err := client.Do(req)
	if err != nil {
		return rel, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return rel, fmt.Errorf("release service returned HTTP %d", resp.StatusCode)
	}
	if err = json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&rel); err != nil {
		return rel, err
	}
	return rel, nil
}

// fetchStableReleaseWithFallback probes every endpoint. It returns the first
// payload that decodes AND passes describeRelease's stability gate; that gate
// also rejects mirror-decoy/empty payloads. All per-endpoint errors are joined
// so callers can show the exact failing step ("官方源 403，镜像 1 超时…").
func fetchStableReleaseWithFallback(ctx context.Context, client *http.Client, endpoints []string, current, goos, goarch string) (releaseUpdate, string, []error) {
	var errs []error
	for i, endpoint := range endpoints {
		if err := ctx.Err(); err != nil {
			return releaseUpdate{}, "", append(errs, err)
		}
		rel, err := fetchStableRelease(ctx, client, strings.TrimRight(endpoint, "/")+"/repos/wxcvm/ARTEX/releases/latest")
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", endpointShort(endpoint), err))
			continue
		}
		out, err := describeRelease(current, goos, goarch, rel)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", endpointShort(endpoint), err))
			continue
		}
		source := "official"
		if i > 0 {
			source = "mirror"
		}
		// errs carries the failures of earlier hops so the handler can tell
		// "answered by mirror after official 403" apart from "answered directly".
		return out, source, errs
	}
	return releaseUpdate{}, "", errs
}

func endpointShort(endpoint string) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		return endpoint
	}
	return u.Host
}

// Read-only, explicit update check: no installation, process exit or writes.
// Proxy → official GitHub API → community mirrors. The frontend receives the
// endpoint that answered, and a joined error listing every failed hop when all
// of them fail, so "check update" no longer dead-ends on restricted networks.
func (s *Server) checkReleaseUpdate(w http.ResponseWriter, r *http.Request) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	defer transport.CloseIdleConnections()
	if s.m != nil {
		if raw := s.m.GlobalProxy(); raw != "" {
			proxy, err := url.Parse(raw)
			if err != nil {
				writeErr(w, 400, "invalid configured proxy")
				return
			}
			transport.Proxy = http.ProxyURL(proxy)
		}
	}
	client := &http.Client{Transport: transport, Timeout: 15 * time.Second}
	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()
	endpoints := releaseAPIEndpointsFor(r.Context())
	result, source, errs := fetchStableReleaseWithFallback(ctx, client, endpoints, BuildVersion, runtime.GOOS, runtime.GOARCH)
	if len(errs) == len(endpoints) {
		writeErr(w, 502, "release check failed on all endpoints: "+joinErrors(errs))
		return
	}
	result.Source = source
	if source == "mirror" {
		note := "官方 GitHub API 不可达，结果来自社区镜像；版本信息与官方发布一致，下载链接仍指向 github.com。"
		if len(errs) > 0 {
			note += "已跳过：" + joinErrors(errs)
		}
		result.Note = note
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, 200, result)
}

func joinErrors(errs []error) string {
	msgs := make([]string, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			msgs = append(msgs, err.Error())
		}
	}
	return strings.Join(msgs, "; ")
}
