package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/config"
	"github.com/Autumn-27/artex/db"
)

// PrivacyOverview is an aggregate snapshot. Unlike the detail endpoints it
// never contains addresses, profile IDs, paths, or individual request events.
type PrivacyOverview struct {
	CheckedAt time.Time `json:"checked_at"`
	Listener ListenerOverview `json:"listener"`
	ModelKeys ModelKeysOverview `json:"model_keys"`
	Config ConfigOverview `json:"config"`
	Outbound OutboundOverview `json:"outbound"`
	ExecutionNetwork ExecutionNetworkOverview `json:"execution_network"`
}

type ListenerOverview struct {
	Scope string `json:"scope"`
	BindingSource string `json:"binding_source"`
	PublicReachability string `json:"public_reachability"`
}

type ModelKeysOverview struct {
	State string `json:"state"`
	KeyState string `json:"key_state"`
	Total int `json:"total"`
	Empty int `json:"empty"`
	PlaintextOnly int `json:"plaintext_only"`
	EncryptedOnly int `json:"encrypted_only"`
	Mixed int `json:"mixed"`
	Verified int `json:"verified"`
	Failed int `json:"failed"`
	Unchecked int `json:"unchecked"`
}

type ConfigOverview struct {
	FileState string `json:"file_state"`
	FilePermissions string `json:"file_permissions"`
	DatabaseSource string `json:"database_source"`
	ProfilesState string `json:"profiles_state"`
	ProfilesChecked int `json:"profiles_checked"`
	Findings int `json:"findings"`
	IncompleteFindings int `json:"incomplete_findings"`
}

type OutboundOverview struct {
	Total uint64 `json:"total"`
	Retained int `json:"retained"`
	Capacity int `json:"capacity"`
	HTTPResponses int `json:"http_responses"`
	HTTPNon2xx int `json:"http_non_2xx"`
	TransportFailures int `json:"transport_failures"`
	Canceled int `json:"canceled"`
	Timeouts int `json:"timeouts"`
	ProxyConfigured int `json:"proxy_configured"`
}

// ExecutionNetworkOverview reports Mode B configuration plus a short-lived,
// aggregate egress comparison. It never exposes the observed IP or claims
// visibility into Docker routes, individual hops, VPS2 or a working kill switch.
type ExecutionNetworkOverview struct {
	Mode string `json:"mode"`
	Configuration string `json:"configuration"`
	Verification string `json:"verification"`
	Scope string `json:"scope"`
	FinalExit string `json:"final_exit"`
	ActualEgress string `json:"actual_egress"`
	ObservedAt *time.Time `json:"observed_at"`
	FailClosed string `json:"fail_closed"`
	VerifyBeforeTask bool `json:"verify_before_task"`
	CoveredProcesses []string `json:"covered_processes"`
}

func executionNetworkOverview(getenv func(string) string) ExecutionNetworkOverview {
	result := ExecutionNetworkOverview{
		Mode: "off", Configuration: "not_configured", Verification: "not_tested",
		Scope: "none", FinalExit: "not_tested", ActualEgress: "not_tested",
		FailClosed: "not_tested", VerifyBeforeTask: false, CoveredProcesses: []string{},
	}
	socket := getenv("ARTEX_RUNNER_SOCKET")
	proxy := getenv("ARTEX_RUNNER_PROXY_ADDR")
	if socket == "" && proxy == "" {
		return result
	}
	result.Mode = "other"
	result.Configuration = "partial"
	if socket != "" {
		result.Scope = "runner_builtins_intended"
	}
	if socket == "/run/artex-runner/runner.sock" && proxy == "http://mode-b-egress:7890" {
		result.Mode = "B"
		result.Configuration = "configured_only"
		result.CoveredProcesses = []string{"linux_docker_tool_executor_proxy_aware"}
	}
	return result
}

func summarizeModelKeys(s db.ModelKeyStatus) ModelKeysOverview {
	return ModelKeysOverview{State:"complete", KeyState:s.KeyState, Total:s.Total, Empty:s.Empty,
		PlaintextOnly:s.PlaintextOnly, EncryptedOnly:s.EncryptedOnly, Mixed:s.Mixed,
		Verified:s.Verified, Failed:s.Failed, Unchecked:s.Unchecked}
}

func summarizeConfig(s config.SensitiveConfigStatus) ConfigOverview {
	out := ConfigOverview{FileState:s.FileState, FilePermissions:s.FilePermissions, DatabaseSource:s.DatabaseSource,
		ProfilesState:s.ProfilesState, ProfilesChecked:s.ProfilesChecked, Findings:len(s.Findings)}
	for _, f := range s.Findings {
		if f.Code == "url_not_inspected" || f.Code == "query_not_inspected" {
			out.IncompleteFindings++
		}
	}
	return out
}

func summarizeOutbound(s agent.ModelOutboundSnapshot) OutboundOverview {
	out := OutboundOverview{Total:s.Total, Retained:len(s.Events), Capacity:s.Capacity}
	for _, e := range s.Events {
		if e.ProxyConfigured { out.ProxyConfigured++ }
		switch e.Result {
		case "http_response":
			out.HTTPResponses++
			if e.Status < 200 || e.Status >= 300 { out.HTTPNon2xx++ }
		case "canceled": out.Canceled++
		case "timeout": out.Timeouts++
		default: out.TransportFailures++
		}
	}
	return out
}

func (s *Server) privacyOverviewSnapshot(ctx context.Context) PrivacyOverview {
	e := s.exposureSnapshot()
	out := PrivacyOverview{
		Listener: ListenerOverview{Scope:e.Scope, BindingSource:e.BindingSource, PublicReachability:e.PublicReachability},
		ModelKeys: ModelKeysOverview{State:"unavailable", KeyState:"not_inspected"},
		Config: summarizeConfig(s.sensitiveConfigSnapshot(ctx)),
		Outbound: summarizeOutbound(agent.ModelOutboundAudit()),
		ExecutionNetwork: applyModeBProbeSnapshot(executionNetworkOverview(os.Getenv)),
	}
	if s.m != nil && s.m.pg != nil {
		if status, err := s.m.pg.ModelKeyDiagnostics(ctx); err == nil {
			out.ModelKeys = summarizeModelKeys(status)
		}
	}
	out.CheckedAt = time.Now().UTC()
	return out
}

func (s *Server) privacyOverview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, s.privacyOverviewSnapshot(r.Context()))
}

// privacyBridgeOverview uses the same aggregate snapshot as the administrator
// overview, but only the separately scoped opt-in credential can read it.
func (s *Server) privacyBridgeOverview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, struct {
		Service string `json:"service"`
		SchemaVersion int `json:"schema_version"`
		Overview PrivacyOverview `json:"overview"`
	}{Service: "artex", SchemaVersion: 1, Overview: s.privacyOverviewSnapshot(r.Context())})
}
