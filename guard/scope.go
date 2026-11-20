// Task-scope enforcement at the tool boundary.
//
// guard.go records that "the RoE authorization-scope mechanism was removed; a
// replacement may be added later." This file is that replacement. Two documents
// in this repository treat the surface as non-negotiable, so leaving it removed
// makes their claims untrue:
//
//   - RSI.md, Loop C contract item 3, lists "scope and authorization checks"
//     among the immutable surfaces that harness optimization may never search.
//   - AUTONOMY_ROADMAP.md §5 states that discovering a new target outside the
//     constraints does not grant authorization for it.
//
// Only the enforcement half went away: agent/goals.go still requires the planner
// to register the task's stated asset scope through add_task_scope ("this task's
// authorization boundary, and also the denominator of asset test coverage"), and
// agent/constraints.go still renders those constraints into the system prompt. A
// prompt-level boundary is a request to the model, not a guarantee about what the
// worker then does with a shell. This file makes it a guarantee.
//
// Enforcement runs in the direction that helps autonomy rather than hurting it.
// Every turn a worker spends on a host the task never declared is budget spent
// off-target, and it can produce findings the report is not entitled to claim. A
// precise block that names the boundary lets the agent pivot inside the scope
// instead of retrying a host it was never given.
//
// Deliberate limits:
//
//   - No registered scope => nothing is blocked. Tasks that declare no scope keep
//     today's behavior exactly; ScopeOff is the default mode.
//   - Loopback, RFC 6761/8375 special-use names, and a fixed infrastructure
//     allowlist (tool downloads, module proxies, CVE databases) are never treated
//     as test targets. Blocking them would break normal operation while
//     protecting nobody.
//   - Scope kinds with no host shape (company, icp, keyword) cannot be resolved
//     to hosts without the company database, so they are recorded in the reason
//     text and never used to allow or block on their own. A scope that resolves
//     to nothing host-shaped is not enforceable and allows everything.
//   - Target extraction is a lexical heuristic over the tool input, biased toward
//     missing a target rather than inventing one.
package guard

import (
	"fmt"
	"net"
	"net/netip"
	"regexp"
	"strings"
	"time"

	"github.com/Autumn-27/norma/hook"
)

// ScopeMode selects how a resolved scope is applied to a tool call.
type ScopeMode int

const (
	// ScopeOff consults no scope and changes nothing. The zero value, so a Guard
	// that was never given a source behaves exactly as it did before.
	ScopeOff ScopeMode = iota
	// ScopeWarn audits an out-of-scope target or a scope lookup failure and
	// lets the call proceed.
	ScopeWarn
	// ScopeEnforce blocks gated calls when a target is out of scope or the
	// configured scope source cannot be read.
	ScopeEnforce
)

// String renders the mode as the token ParseScopeMode accepts.
func (m ScopeMode) String() string {
	switch m {
	case ScopeWarn:
		return "warn"
	case ScopeEnforce:
		return "enforce"
	default:
		return "off"
	}
}

// ParseScopeMode accepts "off", "warn", or "enforce" (case- and trim-insensitive,
// with "on" as an alias for enforce). ok is false for an unrecognized value so a
// caller can report the bad configuration instead of silently picking a mode.
func ParseScopeMode(s string) (ScopeMode, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "off", "none", "disabled":
		return ScopeOff, true
	case "warn", "warning", "audit":
		return ScopeWarn, true
	case "enforce", "on", "block":
		return ScopeEnforce, true
	}
	return ScopeOff, false
}

// ScopeEntry is one registered authorization boundary, normalized to the shape
// enforcement needs.
type ScopeEntry struct {
	// Kind is one of the db.TaskScope kinds: cidr, ip, domain, root_domain,
	// subdomain, company, icp, keyword.
	Kind string
	// Value is the boundary itself: a CIDR, an IP, a hostname, a root domain, or a
	// free-text value for the kinds that have no host shape.
	Value string
}

// ScopeSource resolves the current task's registered scope. The planner can add
// scope during a run through add_task_scope, so Guard reads it per call instead
// of snapshotting it at construction (see scopeCacheTTL).
type ScopeSource interface {
	TaskScope() ([]ScopeEntry, error)
}

// ScopeSourceFunc adapts a function to ScopeSource.
type ScopeSourceFunc func() ([]ScopeEntry, error)

// TaskScope calls f.
func (f ScopeSourceFunc) TaskScope() ([]ScopeEntry, error) { return f() }

// scopeCacheTTL bounds how stale a cached scope may be. add_task_scope can widen
// the boundary mid-task; a short TTL means the widening takes effect on the next
// call without one query per tool call.
const scopeCacheTTL = 2 * time.Second

// maxTargetsPerCall bounds the work one tool input can create.
const maxTargetsPerCall = 128

// scopeRule is one compiled boundary. Exactly one of host/prefix/addr is set.
type scopeRule struct {
	kind   string
	host   string
	prefix netip.Prefix
	addr   netip.Addr
	raw    string
}

// scopeSet is a compiled scope: the host and network rules that can decide a
// call, plus the free-text values that can only be shown to a human.
type scopeSet struct {
	rules    []scopeRule
	keywords []string
}

// enforceable reports whether the scope contains at least one boundary that can
// decide a host. A keyword-only or empty scope cannot, and must not block.
func (s scopeSet) enforceable() bool { return len(s.rules) > 0 }

// describe renders the scope for a block message, capped so a /8-heavy task does
// not paste a hundred entries into the model's context.
func (s scopeSet) describe() string {
	if !s.enforceable() {
		if len(s.keywords) == 0 {
			return "no registered scope"
		}
		return "no enforceable host or network boundary (registered: " + strings.Join(capStrings(s.keywords, 6), ", ") + ")"
	}
	raw := make([]string, 0, len(s.rules))
	for _, r := range s.rules {
		raw = append(raw, r.raw)
	}
	return strings.Join(capStrings(raw, 8), ", ")
}

// capStrings returns at most n items, appending an ellipsis marker when cut.
func capStrings(in []string, n int) []string {
	if len(in) <= n {
		return in
	}
	out := make([]string, 0, n+1)
	out = append(out, in[:n]...)
	return append(out, "…")
}

// compileScope normalizes entries into a scopeSet. Unparseable host or network
// entries are dropped and not treated as a wildcard: a boundary that cannot be
// read must not silently allow the whole internet.
func compileScope(entries []ScopeEntry) scopeSet {
	var set scopeSet
	for _, e := range entries {
		kind := strings.ToLower(strings.TrimSpace(e.Kind))
		value := strings.TrimSpace(e.Value)
		if value == "" {
			continue
		}
		switch kind {
		case "cidr":
			if p, err := netip.ParsePrefix(stripZone(value)); err == nil {
				set.rules = append(set.rules, scopeRule{kind: kind, prefix: p.Masked(), raw: value})
				continue
			}
			if a, err := netip.ParseAddr(stripZone(value)); err == nil {
				set.rules = append(set.rules, scopeRule{kind: "ip", addr: a.Unmap(), raw: value})
			}
		case "ip":
			if a, err := netip.ParseAddr(stripZone(value)); err == nil {
				set.rules = append(set.rules, scopeRule{kind: kind, addr: a.Unmap(), raw: value})
			}
		case "subdomain":
			// A subdomain entry is an exact host, same as "domain".
			if h := normalizeHost(value); h != "" {
				set.rules = append(set.rules, scopeRule{kind: "domain", host: h, raw: value})
			}
		case "domain", "root_domain":
			if h := normalizeHost(value); h != "" {
				set.rules = append(set.rules, scopeRule{kind: kind, host: h, raw: value})
			}
		default:
			// company|icp|keyword: no host shape. Kept for the reason text only.
			set.keywords = append(set.keywords, value)
		}
	}
	return set
}

// allowed reports whether host falls inside the scope. An empty host means the
// call named no target and cannot be a violation.
func (s scopeSet) allowed(host string) bool {
	host = normalizeHost(host)
	if host == "" {
		return true
	}
	if isLocalOrInfrastructure(host) {
		return true
	}
	if a, err := netip.ParseAddr(stripZone(host)); err == nil {
		a = a.Unmap()
		for _, r := range s.rules {
			if r.prefix.IsValid() && r.prefix.Contains(a) {
				return true
			}
			if r.addr.IsValid() && r.addr == a {
				return true
			}
		}
		return false
	}
	for _, r := range s.rules {
		if r.host == "" {
			continue
		}
		if host == r.host {
			return true
		}
		// Label-boundary suffix only: example.com must not admit evil-example.com.
		if r.kind == "root_domain" && strings.HasSuffix(host, "."+r.host) {
			return true
		}
	}
	return false
}

// stripZone drops an IPv6 zone suffix ("fe80::1%4") that netip cannot parse in
// the plain-address form used here.
func stripZone(v string) string {
	if i := strings.LastIndexByte(v, '%'); i > 0 {
		return v[:i]
	}
	return v
}

// stripHostPort drops a trailing :port from host:port and [ipv6]:port, returning
// the bare host. Bare hosts and bare IPv6 (whose colons are ambiguous) come back
// unchanged, matching db.stripHostPort's intent without importing db.
func stripHostPort(v string) string {
	v = strings.TrimSpace(v)
	if host, _, err := net.SplitHostPort(v); err == nil {
		return host
	}
	return v
}

// normalizeHost canonicalizes a host for comparison: no quotes, no userinfo, no
// port, no trailing root dot, no IPv6 brackets, lower case.
func normalizeHost(v string) string {
	v = strings.TrimSpace(v)
	v = strings.Trim(v, "\"'`")
	if v == "" {
		return ""
	}
	if i := strings.LastIndexByte(v, '@'); i >= 0 {
		v = v[i+1:]
	}
	v = stripHostPort(v)
	v = strings.Trim(v, "[]")
	v = strings.TrimSuffix(v, ".")
	v = strings.ToLower(v)
	if v == "" {
		return ""
	}
	// Reject anything that survived normalization but cannot be a host.
	if strings.ContainsAny(v, " \t\r\n/\\,;") {
		return ""
	}
	return v
}

// localNames are special-use names (RFC 6761/8375 plus the near-universal
// developer convention) that never name a third-party target.
var localNames = map[string]bool{
	"localhost":                true,
	"localhost.localdomain":    true,
	"ip6-localhost":            true,
	"metadata.google.internal": true,
}

var localSuffixes = []string{".localhost", ".local", ".internal", ".test", ".invalid", ".example"}

// infrastructureHosts are where the agent gets its tools, modules, and
// vulnerability data. They are not test targets and must not be gated: a build
// that cannot reach the Go proxy or a release download is broken, not safe.
var infrastructureHosts = []string{
	"github.com", "raw.githubusercontent.com", "codeload.github.com",
	"objects.githubusercontent.com", "api.github.com", "gist.github.com",
	"proxy.golang.org", "sum.golang.org", "storage.googleapis.com",
	"registry.npmjs.org", "npmjs.com", "pypi.org", "files.pythonhosted.org",
	"nvd.nist.gov", "cve.org", "osv.dev", "exploit-db.com",
	"cisa.gov", "attack.mitre.org", "mitre.org",
}

// isLocalOrInfrastructure reports whether a host is out of the target space by
// nature: loopback, unspecified, special-use, or a vendor the agent fetches from.
func isLocalOrInfrastructure(host string) bool {
	if localNames[host] {
		return true
	}
	if a, err := netip.ParseAddr(stripZone(host)); err == nil {
		a = a.Unmap()
		if a.IsLoopback() || a.IsUnspecified() {
			return true
		}
	}
	for _, suffix := range localSuffixes {
		if strings.HasSuffix(host, suffix) {
			return true
		}
	}
	for _, infra := range infrastructureHosts {
		if host == infra || strings.HasSuffix(host, "."+infra) {
			return true
		}
	}
	return false
}

var (
	reURLTarget  = regexp.MustCompile(`(?i)\b[a-z][a-z0-9+.\-]{0,15}://[^\s"'<>()\\]+`)
	reIPv4Target = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	reHostTarget = regexp.MustCompile(`(?i)\b[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?)+\b`)
)

// nonHostSuffixes are final labels that make a dotted token a file or a module
// path rather than a hostname. Without this list a command naming jwt.txt or
// best64.rule would read as two out-of-scope hosts.
var nonHostSuffixes = map[string]bool{
	"php": true, "php3": true, "php4": true, "php5": true, "phtml": true,
	"jsp": true, "jspx": true, "jspf": true, "asp": true, "aspx": true,
	"ashx": true, "asmx": true, "cfm": true, "cgi": true, "pl": true,
	"html": true, "htm": true, "xhtml": true, "shtml": true,
	"js": true, "mjs": true, "cjs": true, "ts": true, "tsx": true, "jsx": true,
	"vue": true, "css": true, "scss": true, "less": true,
	"json": true, "jsonl": true, "yaml": true, "yml": true, "toml": true,
	"ini": true, "cfg": true, "conf": true, "properties": true, "env": true,
	"txt": true, "md": true, "markdown": true, "rst": true, "log": true,
	"out": true, "err": true, "sh": true, "bash": true, "zsh": true,
	"ps1": true, "psm1": true, "bat": true, "cmd": true, "vbs": true,
	"py": true, "pyc": true, "pyw": true, "go": true, "rb": true, "java": true,
	"kt": true, "cs": true, "c": true, "cc": true, "cpp": true, "h": true,
	"hpp": true, "rs": true, "lua": true, "sql": true,
	"db": true, "sqlite": true, "sqlite3": true, "mdb": true,
	"lock": true, "sum": true, "mod": true,
	"pem": true, "crt": true, "cer": true, "key": true, "p12": true,
	"pfx": true, "csr": true,
	"zip": true, "tar": true, "gz": true, "tgz": true, "bz2": true,
	"xz": true, "7z": true, "rar": true, "iso": true, "img": true,
	"bin": true, "exe": true, "dll": true, "so": true, "dylib": true,
	"apk": true, "dex": true, "jar": true, "war": true, "ear": true,
	"class": true, "pyz": true, "whl": true, "egg": true, "rpm": true, "deb": true,
	"pcap": true, "pcapng": true, "cap": true, "rule": true, "hc22000": true,
	"potfile": true, "wordlist": true, "dic": true, "lst": true,
	"csv": true, "tsv": true, "xls": true, "xlsx": true, "doc": true,
	"docx": true, "pdf": true, "ppt": true, "pptx": true,
	"png": true, "jpg": true, "jpeg": true, "gif": true, "svg": true,
	"webp": true, "ico": true, "mp4": true, "mp3": true, "wav": true,
	"ttf": true, "woff": true, "woff2": true,
}

// TargetsFromInput returns the host-shaped targets a tool call mentions, without
// duplicates. Order is by extraction class — URL authorities first, then IPv4
// literals, then bare hostnames — and first-seen within each class, so the output
// is deterministic for a given input. File names, module paths, and version
// strings are dropped; the heuristic errs toward returning nothing rather than
// returning a target that is not there.
func TargetsFromInput(raw []byte) []string {
	if len(raw) == 0 {
		return nil
	}
	text := string(raw)
	seen := make(map[string]struct{}, 8)
	out := make([]string, 0, 8)

	add := func(candidate string) {
		if len(out) >= maxTargetsPerCall {
			return
		}
		host := normalizeHost(candidate)
		if !isHostShaped(host) {
			return
		}
		if _, dup := seen[host]; dup {
			return
		}
		seen[host] = struct{}{}
		out = append(out, host)
	}

	for _, m := range reURLTarget.FindAllString(text, -1) {
		i := strings.Index(m, "://")
		if i < 0 {
			continue
		}
		authority := m[i+3:]
		if j := strings.IndexAny(authority, "/?#"); j >= 0 {
			authority = authority[:j]
		}
		add(authority)
	}
	for _, m := range reIPv4Target.FindAllString(text, -1) {
		add(m)
	}
	for _, m := range reHostTarget.FindAllString(text, -1) {
		add(m)
	}
	return out
}

// isHostShaped reports whether a normalized token can be a hostname or address.
func isHostShaped(host string) bool {
	if host == "" {
		return false
	}
	if _, err := netip.ParseAddr(stripZone(host)); err == nil {
		return true
	}
	if !strings.Contains(host, ".") {
		return false
	}
	labels := strings.Split(host, ".")
	last := labels[len(labels)-1]
	if len(last) < 2 || len(last) > 24 {
		return false
	}
	if nonHostSuffixes[last] {
		return false
	}
	allDigits := true
	for _, c := range last {
		if c < '0' || c > '9' {
			allDigits = false
			break
		}
	}
	return !allDigits
}

// UseScope installs the scope source and mode, dropping any cached resolution.
// Call it before the guard is attached to a session; it is safe to call again to
// switch modes at runtime.
func (g *Guard) UseScope(mode ScopeMode, src ScopeSource) {
	g.scopeMu.Lock()
	defer g.scopeMu.Unlock()
	g.scopeMode, g.scopeSource = mode, src
	g.scopeCache = scopeSnapshot{}
}

// ScopeMode reports the installed mode (ScopeOff when none was installed).
func (g *Guard) ScopeMode() ScopeMode {
	g.scopeMu.Lock()
	defer g.scopeMu.Unlock()
	return g.scopeMode
}

type scopeSnapshot struct {
	set     scopeSet
	at      time.Time
	present bool
}

// resolvedScope returns the task's compiled scope, refreshing the cache when it
// is older than scopeCacheTTL. ok is false when the source could not be read.
func (g *Guard) resolvedScope(src ScopeSource) (scopeSet, bool) {
	g.scopeMu.Lock()
	if g.scopeCache.present && time.Since(g.scopeCache.at) < scopeCacheTTL {
		set := g.scopeCache.set
		g.scopeMu.Unlock()
		return set, true
	}
	g.scopeMu.Unlock()

	entries, err := src.TaskScope()
	if err != nil {
		// Leave the decision to checkScope and let the next call retry the lookup.
		return scopeSet{}, false
	}
	set := compileScope(entries)
	g.scopeMu.Lock()
	g.scopeCache = scopeSnapshot{set: set, at: time.Now(), present: true}
	g.scopeMu.Unlock()
	return set, true
}

// scopeGatedTool reports whether a tool's input names a target the call will
// actually reach. Gating is deliberately restricted to that surface: a file
// write, a report edit, or a knowledge entry may quote an out-of-scope URL
// without touching it, and blocking those would refuse legitimate work while
// protecting nothing. The agent's shell and its HTTP/browser tools are where
// target access happens, so they are what the gate reads.
//
// Tool names arriving from third-party MCP servers are outside this list in v1.
// That is a real gap, not an oversight: an unknown MCP tool that reaches the
// network is not gated. It is named here so it is not mistaken for coverage.
func scopeGatedTool(tool string) bool {
	switch strings.ToLower(strings.TrimSpace(tool)) {
	case "bash", "shell_open", "shell_send",
		"webfetch", "web_fetch", "http", "http_request", "request",
		"browser", "browser_navigate", "navigate", "playwright":
		return true
	}
	return false
}

// checkScope evaluates the installed scope for gated tool calls.
// It returns the block result and true when the call must be refused; warn mode
// audits lookup failures and violations without blocking.
func (g *Guard) checkScope(tool, command string, input []byte) (hook.Result, bool) {
	g.scopeMu.Lock()
	mode, src := g.scopeMode, g.scopeSource
	g.scopeMu.Unlock()
	if mode == ScopeOff || src == nil || !scopeGatedTool(tool) {
		return hook.Result{}, false
	}
	set, ok := g.resolvedScope(src)
	if !ok {
		msg := "task scope could not be verified because the scope source is unavailable"
		if mode == ScopeWarn {
			g.record(tool, "warn", msg, command)
			return hook.Result{}, false
		}
		return g.block(systemBlockMessage(msg)), true
	}
	if !set.enforceable() {
		return hook.Result{}, false
	}
	var offending []string
	for _, target := range TargetsFromInput(input) {
		if !set.allowed(target) {
			offending = append(offending, target)
		}
	}
	if len(offending) == 0 {
		return hook.Result{}, false
	}
	msg := scopeBlockMessage(offending, set.describe())
	if mode == ScopeWarn {
		g.record(tool, "warn", msg, command)
		return hook.Result{}, false
	}
	return g.block(msg), true
}

// scopeBlockMessage reuses the platform-control framing so the worker registers
// the refusal as a policy decision rather than a target-side defense it should
// try to defeat.
func scopeBlockMessage(targets []string, scope string) string {
	verb := "are"
	if len(targets) == 1 {
		verb = "is"
	}
	return systemBlockMessage(fmt.Sprintf(
		"%s %s outside this task's declared test scope (%s). Do not retry or obfuscate it. "+
			"If the user's task goal explicitly names it, widen the scope with add_task_scope; "+
			"otherwise record it as an out-of-scope fact and keep working inside the scope",
		strings.Join(targets, ", "), verb, scope))
}
