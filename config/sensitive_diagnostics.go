package config

import (
	"encoding/json"
	"errors"
	"net"
	"net/url"
	"os"
	"runtime"
	"strings"
)

// ConfigFinding contains only fixed identifiers, never a value or an address.
type ConfigFinding struct {
	Source string `json:"source"`
	Field string `json:"field"`
	Code string `json:"code"`
	ProfileID string `json:"profile_id,omitempty"`
}

type SensitiveConfigStatus struct {
	FileState string `json:"file_state"`
	FilePermissions string `json:"file_permissions"`
	DatabaseSource string `json:"database_source"`
	ProfilesState string `json:"profiles_state"`
	ProfilesChecked int `json:"profiles_checked"`
	Findings []ConfigFinding `json:"findings"`
}

func sensitiveQueryKey(key string) bool {
	switch strings.ToLower(strings.ReplaceAll(key, "-", "_")) {
	case "key", "api_key", "apikey", "token", "access_token", "refresh_token",
		"password", "passwd", "secret", "client_secret", "authorization", "auth", "sig", "signature":
		return true
	}
	return false
}

// InspectConfigURL inspects a configured URL without resolving DNS, connecting
// to its host, or returning any part of the actual URL.
func InspectConfigURL(raw, source, field, profileID string, schemes ...string) []ConfigFinding {
	out := []ConfigFinding{}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return out
	}
	add := func(code string) {
		out = append(out, ConfigFinding{Source: source, Field: field, ProfileID: profileID, Code: code})
	}
	u, err := url.Parse(raw)
	if err != nil || u.Opaque != "" || u.Scheme == "" {
		add("url_not_inspected")
		return out
	}
	allowed := false
	for _, scheme := range schemes {
		if strings.EqualFold(u.Scheme, scheme) {
			allowed = true
			break
		}
	}
	if !allowed || (u.Hostname() == "" && field != "database.dsn") {
		add("url_not_inspected")
		return out
	}
	if u.User != nil {
		if _, present := u.User.Password(); present || u.User.Username() != "" {
			add("url_userinfo")
		}
	}
	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		add("query_not_inspected")
	}
	for key := range query {
		if sensitiveQueryKey(key) {
			add("sensitive_query")
			break
		}
	}
	if field == "base_url" && strings.EqualFold(u.Scheme, "http") {
		host := u.Hostname()
		ip := net.ParseIP(host)
		if !strings.EqualFold(host, "localhost") && (ip == nil || !ip.IsLoopback()) {
			add("http_non_loopback")
		}
	}
	return out
}

// SensitiveConfigDiagnostics reports local configuration metadata. It inspects
// the file even when an environment override is active, since a dormant file
// can still contain credentials. It never changes configuration.
func SensitiveConfigDiagnostics() SensitiveConfigStatus {
	s := SensitiveConfigStatus{
		FileState: "not_inspected", FilePermissions: "not_inspected",
		DatabaseSource: "missing", ProfilesState: "not_inspected",
		Findings: []ConfigFinding{},
	}
	path := Path()
	data, err := os.ReadFile(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		s.FileState = "missing"
	case err != nil:
		s.FileState = "unreadable"
	default:
		var c Config
		if err := json.Unmarshal(data, &c); err != nil {
			s.FileState = "invalid"
		} else {
			s.FileState = "readable"
			d := c.Database
			if strings.TrimSpace(d.Password) != "" {
				s.Findings = append(s.Findings, ConfigFinding{Source: "config_file", Field: "database.password", Code: "plaintext_password"})
			}
			s.Findings = append(s.Findings, InspectConfigURL(d.DSN, "config_file", "database.dsn", "", "postgres", "postgresql")...)
			if strings.TrimSpace(d.DSN) != "" {
				s.DatabaseSource = "file_dsn"
			} else if d.Host != "" || d.DBName != "" || d.User != "" {
				s.DatabaseSource = "file_fields"
			}
		}
		if runtime.GOOS == "windows" {
			s.FilePermissions = "windows_acl_not_inspected"
		} else if info, err := os.Stat(path); err == nil {
			if info.Mode().Perm()&0044 != 0 {
				s.FilePermissions = "group_or_other_read"
			} else {
				s.FilePermissions = "no_group_or_other_read"
			}
		}
	}
	if dsn := strings.TrimSpace(os.Getenv("ARTEX_PG_DSN")); dsn != "" {
		s.DatabaseSource = "environment"
		s.Findings = append(s.Findings, InspectConfigURL(dsn, "environment", "database.dsn", "", "postgres", "postgresql")...)
	}
	s.Findings = append(s.Findings, InspectConfigURL(os.Getenv("ARTEX_LLM_BASE_URL"), "environment", "base_url", "", "http", "https")...)
	s.Findings = append(s.Findings, InspectConfigURL(os.Getenv("ARTEX_LLM_PROXY"), "environment", "proxy", "", "http", "https", "socks5")...)
	return s
}
