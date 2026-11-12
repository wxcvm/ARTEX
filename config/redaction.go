package config

import (
	"regexp"
)

// These patterns are intended for diagnostic text, not raw request/response
// captures or evidence. Keep the surrounding error context and mask only
// recognizable credential representations.
var (
	logURLUserinfo = regexp.MustCompile("(?i)([a-z][a-z0-9+.-]*://)[^\\s/?#\"'<>]+@")
	logAuthValue = regexp.MustCompile("(?i)(\\b(?:proxy-)?authorization\\b[\"']?[ \\t]*[:=][ \\t]*[\"']?(?:bearer|basic)[ \\t]+)[A-Za-z0-9._~+/-]+=*")
	logBearer = regexp.MustCompile("(?i)(\\bbearer[ \\t]+)[A-Za-z0-9._~+/-]+=*")
	logQuerySecret = regexp.MustCompile("(?i)([?&](?:api[_-]?key|access[_-]?token|refresh[_-]?token|token|password|secret|key)=)[^\\s&#\"'<>]+")
	logNamedSecret = regexp.MustCompile("(?i)(\\b(?:api[_-]?key|password|client[_-]?secret|access[_-]?token|refresh[_-]?token|proxy[_-]?password|pgpassword|ARTEX_LLM_SECRET_KEY|ARTEX_PRIVACY_BRIDGE_TOKEN|OPENAI_API_KEY|ANTHROPIC_API_KEY)\\b[\"']?[ \\t]*[:=][ \\t]*)(?:\"(?:\\\\.|[^\"\\\\])*\"|'(?:\\\\.|[^'\\\\])*'|[^\\s,;}\\]&]+)")
)

// RedactLogText masks common credential representations before diagnostics are
// displayed or persisted. It cannot recognize arbitrary unlabelled secrets.
func RedactLogText(text string) string {
	text = logURLUserinfo.ReplaceAllString(text, "1***@")
	text = logAuthValue.ReplaceAllString(text, "1***")
	text = logBearer.ReplaceAllString(text, "1***")
	text = logQuerySecret.ReplaceAllString(text, "1***")
	return logNamedSecret.ReplaceAllStringFunc(text, func(match string) string {
		indices := logNamedSecret.FindStringSubmatchIndex(match)
		if len(indices) < 4 {
			return "***"
		}
		prefix := match[:indices[3]]
		value := match[indices[3]:]
		if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
			return prefix + string(value[0]) + "***" + string(value[0])
		}
		return prefix + "***"
	})
}
