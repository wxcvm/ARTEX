package server

import (
	"net"
	"net/http"
	"strings"
)

// ExposureStatus describes the HTTP listener inside the process's current network
// namespace, not an internet scan. Host firewalls, NAT, Docker port publication
// and reverse proxies are outside this assessment.
type ExposureStatus struct {
	ListenAddress      string `json:"listen_address"`
	ConfiguredAddress  string `json:"configured_address"`
	BoundAddress       string `json:"bound_address,omitempty"`
	BindingSource      string `json:"binding_source"`
	Scope              string `json:"scope"`
	PublicReachability string `json:"public_reachability"`
	Detail             string `json:"detail"`
}

func classifyListener(address string) ExposureStatus {
	status := ExposureStatus{
		ListenAddress:      address,
		Scope:              "unknown",
		PublicReachability: "not_tested",
		Detail:             "No public reachability test was performed. Check host firewall, Docker port mappings, NAT and reverse proxies separately.",
	}
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return status
	}
	if host == "" || host == "*" {
		status.Scope = "all_interfaces"
		return status
	}
	if strings.EqualFold(host, "localhost") {
		status.Scope = "loopback"
		return status
	}
	ipText := host
	if zone := strings.IndexByte(ipText, '%'); zone >= 0 {
		ipText = ipText[:zone]
	}
	ip := net.ParseIP(ipText)
	if ip == nil {
		return status // Hostnames may resolve differently on each machine.
	}
	switch {
	case ip.IsUnspecified():
		status.Scope = "all_interfaces"
	case ip.IsLoopback():
		status.Scope = "loopback"
	default:
		status.Scope = "specific_interface"
	}
	return status
}

// SetListenAddress is called once at startup before the HTTP server runs.
func (s *Server) SetListenAddress(address string) {
	s.listenAddress = address
}

// SetBoundListenAddress records net.Listener.Addr().String() once, before serving.
// Unlike the configured address, it includes the OS-selected port and IP.
func (s *Server) SetBoundListenAddress(address string) {
	s.boundListenAddress = address
}

func (s *Server) exposureSnapshot() ExposureStatus {
	address, source := s.listenAddress, "not_inspected"
	if address != "" {
		source = "configured"
	}
	if s.boundListenAddress != "" {
		address, source = s.boundListenAddress, "runtime_listener"
	}
	status := classifyListener(address)
	status.ConfiguredAddress = s.listenAddress
	status.BoundAddress = s.boundListenAddress
	status.BindingSource = source
	return status
}

func (s *Server) exposureStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.exposureSnapshot())
}

// privacyBridgeStatus is guarded by requireAuth's exact-route token check.
// Return only the coarse listener assessment; no raw IPs, paths, keys or JWTs.
func (s *Server) privacyBridgeStatus(w http.ResponseWriter, _ *http.Request) {
	status := s.exposureSnapshot()
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, map[string]any{
		"service":             "artex",
		"scope":               status.Scope,
		"binding_source":      status.BindingSource,
		"public_reachability": status.PublicReachability,
	})
}
