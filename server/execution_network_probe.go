package server

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Autumn-27/artex/toolrunner"
)

type modeBProbeSnapshot struct {
	Verification string
	FinalExit    string
	ActualEgress string
	ObservedAt   time.Time
}

var modeBProbeState struct {
	sync.RWMutex
	snapshot modeBProbeSnapshot
}

func recordModeBProbe(verification, finalExit, actualEgress string) time.Time {
	now := time.Now().UTC()
	modeBProbeState.Lock()
	modeBProbeState.snapshot = modeBProbeSnapshot{
		Verification: verification,
		FinalExit: finalExit,
		ActualEgress: actualEgress,
		ObservedAt: now,
	}
	modeBProbeState.Unlock()
	return now
}

func applyModeBProbeSnapshot(out ExecutionNetworkOverview) ExecutionNetworkOverview {
	if out.Mode != "B" {
		return out
	}
	modeBProbeState.RLock()
	snapshot := modeBProbeState.snapshot
	modeBProbeState.RUnlock()
	if snapshot.ObservedAt.IsZero() {
		return out
	}
	out.Verification = snapshot.Verification
	out.FinalExit = snapshot.FinalExit
	out.ActualEgress = snapshot.ActualEgress
	observedAt := snapshot.ObservedAt
	out.ObservedAt = &observedAt
	return out
}

func validExpectedModeBExit(raw string) (string, bool) {
	ip := net.ParseIP(strings.TrimSpace(raw))
	if ip == nil || ip.To4() == nil || !ip.IsGlobalUnicast() || ip.IsPrivate() ||
		ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsMulticast() || ip.IsUnspecified() {
		return "", false
	}
	return ip.To4().String(), true
}

// verifyExecutionNetwork is an explicit operator action. The request is sent
// through the isolated executor's fixed private gateway; the controller does
// not make an outbound request and receives no raw observed IP.
func (s *Server) verifyExecutionNetwork(w http.ResponseWriter, r *http.Request) {
	socket := os.Getenv("ARTEX_RUNNER_SOCKET")
	verifyExecutionNetworkWith(w, r, os.Getenv, func(ctx context.Context, expected string) (toolrunner.EgressProbeResponse, error) {
		return toolrunner.NewClient(socket).ProbeEgress(ctx, expected)
	})
}

func verifyExecutionNetworkWith(
	w http.ResponseWriter,
	r *http.Request,
	getenv func(string) string,
	probe func(context.Context, string) (toolrunner.EgressProbeResponse, error),
) {
	var input struct {
		ExpectedExitIP string `json:"expected_exit_ip"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeErr(w, http.StatusBadRequest, "expected_exit_ip must be a public IPv4 address")
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeErr(w, http.StatusBadRequest, "invalid request body")
		return
	}
	expected, ok := validExpectedModeBExit(input.ExpectedExitIP)
	if !ok {
		writeErr(w, http.StatusBadRequest, "expected_exit_ip must be a public IPv4 address")
		return
	}
	if executionNetworkOverview(getenv).Mode != "B" {
		writeErr(w, http.StatusConflict, "ARTEX Mode B isolated runner is not configured")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	result, err := probe(ctx, expected)
	if err != nil {
		recordModeBProbe("probe_failed", "unavailable", "unavailable")
		writeErr(w, http.StatusBadGateway, "isolated executor egress probe failed; no direct fallback was attempted")
		return
	}
	if !result.OK || result.Status != "exit_ip_match" {
		observedAt := recordModeBProbe("exit_ip_mismatch", "mismatch", "mismatch")
		writeJSON(w, http.StatusOK, map[string]any{
			"ok": false, "status": "exit_ip_mismatch", "observed_at": observedAt,
		})
		return
	}
	observedAt := recordModeBProbe("exit_ip_match", "matched_expected", "matches_expected")
	writeJSON(w, http.StatusOK, map[string]any{
		"ok": true, "status": "exit_ip_match", "observed_at": observedAt,
	})
}
