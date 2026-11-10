package server

import (
	"net"
	"net/http"
	"strings"
)

// ExposureStatus describes the configured HTTP listener, not an internet scan.
// Firewalls, NAT, Docker port publication and reverse proxies are outside what
// this process can establish by inspecting its listener address.
type ExposureStatus struct {
	ListenAddress      string `json:"listen_address"`
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

func (s *Server) exposureStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, classifyListener(s.listenAddress))
}
