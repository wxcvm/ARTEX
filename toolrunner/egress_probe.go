package toolrunner

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/netip"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const modeBEgressEchoURL = "https://api4.ipify.org"

// EgressProbeRequest is deliberately limited to comparing the executor's
// observed public IPv4 with an operator-supplied expected exit. It accepts no
// arbitrary URL, proxy, command, or destination.
type EgressProbeRequest struct {
	ExpectedExitIP string `json:"expected_exit_ip"`
}

type EgressProbeResponse struct {
	OK     bool   `json:"ok"`
	Status string `json:"status"`
}

var nonPublicIPv4Prefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("169.254.0.0/16"),
	netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.88.99.0/24"),
	netip.MustParsePrefix("192.168.0.0/16"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("224.0.0.0/4"),
	netip.MustParsePrefix("240.0.0.0/4"),
}

func publicIPv4(addr netip.Addr) bool {
	if !addr.Is4() || !addr.IsGlobalUnicast() {
		return false
	}
	for _, prefix := range nonPublicIPv4Prefixes {
		if prefix.Contains(addr) {
			return false
		}
	}
	return true
}

// NormalizePublicIPv4 accepts only globally routable IPv4 literals and returns
// their canonical dotted-decimal representation.
func NormalizePublicIPv4(raw string) (string, bool) {
	addr, err := netip.ParseAddr(strings.TrimSpace(raw))
	if err != nil || !publicIPv4(addr) {
		return "", false
	}
	return addr.String(), true
}

// modeBProxyURL accepts only the private Compose gateway address. The probe
// never inherits proxy variables or falls back to a direct transport.
func modeBProxyURL(raw string) (*url.URL, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u == nil || u.Scheme != "http" ||
		u.Hostname() != "mode-b-egress" || u.Port() != "7890" ||
		u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("Mode B private egress gateway is not configured")
	}
	return u, nil
}

// probeModeBEgress makes one fixed HTTPS request through an explicit proxy.
// endpoint is parameterized only for local unit tests; the HTTP handler always
// supplies modeBEgressEchoURL.
func probeModeBEgress(ctx context.Context, proxyURL *url.URL, endpoint string) (string, error) {
	if proxyURL == nil || endpoint == "" {
		return "", errors.New("Mode B egress probe is not configured")
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyURL(proxyURL),
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		DisableKeepAlives:     true,
	}
	defer transport.CloseIdleConnections()
	client := &http.Client{
		Transport: transport,
		Timeout:   12 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("Mode B egress gateway request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Mode B egress endpoint returned a non-success response")
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 65))
	if err != nil || len(body) == 0 || len(body) > 64 {
		return "", errors.New("Mode B egress endpoint returned an invalid response")
	}
	observed, ok := NormalizePublicIPv4(string(body))
	if !ok {
		return "", errors.New("Mode B egress endpoint did not return a public IPv4")
	}
	return observed, nil
}

func (s *Service) handleEgressProbe(w http.ResponseWriter, r *http.Request) {
	var input EgressProbeRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "expected_exit_ip must be a public IPv4 address", http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	expected, ok := NormalizePublicIPv4(input.ExpectedExitIP)
	if !ok {
		http.Error(w, "expected_exit_ip must be a globally routable public IPv4 address", http.StatusBadRequest)
		return
	}
	proxyURL, err := modeBProxyURL(os.Getenv("ARTEX_RUNNER_PROXY_ADDR"))
	if err != nil {
		http.Error(w, "Mode B private egress gateway is not configured", http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 13*time.Second)
	defer cancel()
	observed, err := probeModeBEgress(ctx, proxyURL, modeBEgressEchoURL)
	if err != nil {
		http.Error(w, "Mode B egress probe failed", http.StatusBadGateway)
		return
	}
	result := EgressProbeResponse{OK: observed == expected, Status: "exit_ip_mismatch"}
	if result.OK {
		result.Status = "exit_ip_match"
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
