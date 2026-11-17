package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Autumn-27/artex/traffic"
)

// These fixed public IP-echo endpoints are contacted only by a user-requested
// proxy test. A failure at one service is not a failed proxy if another works.
var proxyEgressEndpoints = []string{
	"https://ifconfig.me/ip",
	"https://icanhazip.com",
	"https://api.ip.sb/ip",
	"https://api.ipify.org",
}

// probeProxyEgress accepts an HTTP client so tests can use a local echo server.
// The transport is configured by testGlobalProxy to use the submitted proxy.
func probeProxyEgress(ctx context.Context, client *http.Client, endpoints []string) (string, time.Duration, error) {
	lastErr := errors.New("no IP-echo endpoint is reachable")
	for _, endpoint := range endpoints {
		if err := ctx.Err(); err != nil {
			return "", 0, err
		}
		callCtx, cancel := context.WithTimeout(ctx, 6*time.Second)
		req, err := http.NewRequestWithContext(callCtx, http.MethodGet, endpoint, nil)
		if err != nil {
			cancel()
			lastErr = err
			continue
		}
		start := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			cancel()
			lastErr = err
			continue
		}
		payload, readErr := io.ReadAll(io.LimitReader(resp.Body, 64))
		_ = resp.Body.Close()
		cancel()
		if readErr != nil {
			lastErr = readErr
			continue
		}
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("IP-echo endpoint returned HTTP %d", resp.StatusCode)
			continue
		}
		ip := strings.TrimSpace(string(payload))
		if net.ParseIP(ip) == nil {
			lastErr = errors.New("IP-echo endpoint returned a non-IP response")
			continue
		}
		return ip, time.Since(start), nil
	}
	return "", 0, lastErr
}

// testGlobalProxy probes the current form's unsaved proxy URL (or the saved URL
// when blank) without changing settings. This makes it possible to diagnose
// egress connectivity before applying a new global proxy to agents.
func (s *Server) testGlobalProxy(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Proxy string `json:"global_proxy"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil && !errors.Is(err, io.EOF) {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	raw := strings.TrimSpace(input.Proxy)
	if raw == "" {
		raw = s.m.GlobalProxy()
	}
	if raw == "" {
		writeJSON(w, http.StatusOK, map[string]any{"ok": false, "error": "No global proxy configured"})
		return
	}
	u, err := traffic.ValidateProxyURL(raw)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"ok": false, "error": "Invalid proxy address (expected http, https or socks5 with a host)"})
		return
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = http.ProxyURL(u)
	defer transport.CloseIdleConnections()
	client := &http.Client{Transport: transport}
	ip, elapsed, err := probeProxyEgress(r.Context(), client, proxyEgressEndpoints)
	if err != nil {
		// Raw proxy URLs may contain credentials. Never echo those back in an error.
		writeJSON(w, http.StatusOK, map[string]any{"ok": false, "error": strings.ReplaceAll(err.Error(), raw, "[proxy]")})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "egress_ip": ip, "latency_ms": elapsed.Milliseconds()})
}
