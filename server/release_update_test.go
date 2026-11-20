package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReleaseUpdateSelection(t *testing.T) {
	var rel githubRelease
	rel.Tag = "v0.5.2"
	rel.PublishedAt = "2026-11-16T00:00:00Z"
	for _, name := range []string{"artex-0.5.2-windows-amd64.zip", "artex-0.5.2-linux-amd64.zip", "ARTEX-Windows-Setup-amd64.exe", "ARTEX-Windows-Desktop-amd64.zip", "SHA256SUMS", "unrelated.exe"} {
		rel.Assets = append(rel.Assets, struct {
			Name string `json:"name"`
			Size int64  `json:"size"`
		}{name, 100})
	}
	for _, tc := range []struct {
		current    string
		comparable bool
		update     bool
	}{
		{"v0.5.1", true, true}, {"0.5.2", true, false}, {"v0.5.3", true, false}, {"v0.5.2-beta.1", true, true}, {"dev", false, false}, {"v0.5", false, false},
	} {
		got, err := describeRelease(tc.current, "windows", "amd64", rel)
		if err != nil || got.Comparable != tc.comparable || got.HasUpdate != tc.update || len(got.Assets) != 4 || got.Assets[0].Name != "ARTEX-Windows-Setup-amd64.exe" {
			t.Fatal(tc, got, err)
		}
		for _, asset := range got.Assets {
			if !strings.HasPrefix(asset.URL, releasePage+"/download/v0.5.2/") {
				t.Fatal(asset)
			}
		}
	}
	got, err := describeRelease("v0.5.1", "linux", "amd64", rel)
	if err != nil || len(got.Assets) != 2 || got.Assets[0].Name != "artex-0.5.2-linux-amd64.zip" {
		t.Fatal(got, err)
	}
	rel.Prerelease = true
	if _, err := describeRelease("v0.5.1", "windows", "amd64", rel); err == nil {
		t.Fatal("accepted prerelease")
	}
}
func TestFetchReleaseProtocol(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/vnd.github+json" {
			t.Error("missing GitHub accept header")
		}
		_, _ = w.Write([]byte(`{"tag_name":"v0.5.2","assets":[{"name":"ARTEX-Windows-Setup-amd64.exe","size":100}]}`))
	}))
	defer mock.Close()
	rel, err := fetchStableRelease(context.Background(), mock.Client(), mock.URL)
	if err != nil || rel.Tag != "v0.5.2" || len(rel.Assets) != 1 {
		t.Fatal(rel, err)
	}
}

// The mirror chain must answer when the first endpoint refuses (403/timeout)
// and must report source=mirror so the UI can label the answer.
func TestFetchReleaseFallbackToMirror(t *testing.T) {
	official := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer official.Close()
	payload := `{"tag_name":"v0.6.0","published_at":"2026-11-18T02:09:04Z","assets":[` +
		`{"name":"ARTEX-Windows-Setup-amd64.exe","size":50824851},` +
		`{"name":"ARTEX-Windows-Desktop-amd64.zip","size":70664366},` +
		`{"name":"artex-0.6.0-windows-amd64.zip","size":16731871},` +
		`{"name":"SHA256SUMS","size":668}]}`
	mirror := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/repos/wxcvm/ARTEX/releases/latest") {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(payload))
	}))
	defer mirror.Close()
	endpoints := []string{official.URL, mirror.URL}
	result, source, errs := fetchStableReleaseWithFallback(context.Background(), mockClientFor(official, mirror), endpoints, "v0.5.2", "windows", "amd64")
	if len(errs) != 1 || source != "mirror" || result.Latest != "v0.6.0" || !result.HasUpdate {
		t.Fatal(result, source, errs)
	}
	if result.Source != "" { // source is set by the handler, not the fetch helper
		t.Fatal("unexpected pre-set source")
	}
	if len(result.Assets) != 4 {
		t.Fatal("asset selection changed", result.Assets)
	}
}

// A mirror that serves a broken payload must not end the chain: the next
// endpoint still gets probed and errors accumulate.
func TestFetchReleaseFallbackSkipsBrokenMirror(t *testing.T) {
	broken := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not json at all"))
	}))
	defer broken.Close()
	good := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"tag_name":"v0.6.0","assets":[{"name":"artex-0.6.0-linux-amd64.zip","size":10}]}`))
	}))
	defer good.Close()
	result, source, errs := fetchStableReleaseWithFallback(context.Background(), mockClientFor(broken, good), []string{broken.URL, good.URL}, "v0.5.0", "linux", "amd64")
	if source != "mirror" || result.Latest != "v0.6.0" {
		t.Fatal(result, source, errs)
	}
	if len(errs) != 1 {
		t.Fatal("expected the broken endpoint to be reported", errs)
	}
}

// When every endpoint fails the caller gets all errors, so the API response
// can explain each failed hop instead of a bare 502.
func TestFetchReleaseAllEndpointsFail(t *testing.T) {
	dead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer dead.Close()
	result, source, errs := fetchStableReleaseWithFallback(context.Background(), dead.Client(), []string{dead.URL, dead.URL}, "v0.5.0", "linux", "amd64")
	if source != "" || result.Latest != "" || len(errs) != 2 {
		t.Fatal(result, source, errs)
	}
	joined := joinErrors(errs)
	if !strings.Contains(joined, "403") || !strings.Contains(joined, ";") {
		t.Fatal(joined)
	}
}

// endpointShort keeps multi-hop mirror URLs readable in error text.
func TestEndpointShort(t *testing.T) {
	if got := endpointShort("https://gh-proxy.com/https://api.github.com"); got != "gh-proxy.com" {
		t.Fatal(got)
	}
	if got := endpointShort("https://api.github.com"); got != "api.github.com" {
		t.Fatal(got)
	}
}

// ARTEX_RELEASE_API_MIRRORS replaces the whole chain (override mode); invalid
// entries are dropped and a fully invalid value falls back to the default.
func TestReleaseAPIEndpointsOverride(t *testing.T) {
	t.Setenv("ARTEX_RELEASE_API_MIRRORS", " https://mirror.example/ , ,https://mirror2.example ")
	got := releaseAPIEndpoints()
	if len(got) != 2 || got[0] != "https://mirror.example" || got[1] != "https://mirror2.example" {
		t.Fatal(got)
	}
	t.Setenv("ARTEX_RELEASE_API_MIRRORS", " , ")
	if got := releaseAPIEndpoints(); len(got) == 0 || got[0] != "https://api.github.com" {
		t.Fatal(got)
	}
}

// The HTTP handler: first hop 403, mirror answers, response carries
// source=mirror and the note, and no 502 is returned.
func TestCheckReleaseUpdateHandlerMirror(t *testing.T) {
	official := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer official.Close()
	mirror := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"tag_name":"v0.6.0","assets":[{"name":"SHA256SUMS","size":668}]}`))
	}))
	defer mirror.Close()
	s := &Server{}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/update/check", nil)
	req = req.WithContext(withReleaseEndpoints(req.Context(), official.URL, mirror.URL))
	s.checkReleaseUpdate(rr, req)
	if rr.Code != 200 {
		t.Fatal(rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"source":"mirror"`) {
		t.Fatal(rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "镜像") {
		t.Fatal("missing mirror note", rr.Body.String())
	}
}

// All hops down → 502 whose body lists every endpoint that failed.
func TestCheckReleaseUpdateHandlerAllFail(t *testing.T) {
	dead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer dead.Close()
	s := &Server{}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/update/check", nil)
	req = req.WithContext(withReleaseEndpoints(req.Context(), dead.URL))
	s.checkReleaseUpdate(rr, req)
	if rr.Code != 502 {
		t.Fatal(rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "release check failed on all endpoints") {
		t.Fatal(rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "127.0.0.1") {
		t.Fatal("joined errors should name the failing hosts", rr.Body.String())
	}
}

func mockClientFor(_ ...*httptest.Server) *http.Client {
	return &http.Client{}
}

func withReleaseEndpoints(ctx context.Context, endpoints ...string) context.Context {
	return context.WithValue(ctx, releaseEndpointsKey{}, endpoints)
}
