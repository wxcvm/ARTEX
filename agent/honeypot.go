package agent

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Autumn-27/artex/db"
	actool "github.com/Autumn-27/norma/tool"
)

var ErrHoneypotTargetQuarantined = errors.New("target quarantined after an explicit honeypot fingerprint")

// HoneypotRegistry shares short-lived target decisions across workers in one
// task. The durable source is PostgreSQL; this in-memory copy avoids a DB query
// on every model/tool call.
type HoneypotRegistry struct {
	mu    sync.RWMutex
	hosts map[string]db.HoneypotAssessment
}

func NewHoneypotRegistry() *HoneypotRegistry {
	return &HoneypotRegistry{hosts: make(map[string]db.HoneypotAssessment)}
}

func normalizeHoneypotTarget(raw string) string {
	raw = strings.TrimSpace(raw)
	if u, err := url.Parse(raw); err == nil && u.Hostname() != "" {
		raw = u.Hostname()
	}
	return strings.TrimSuffix(strings.Trim(strings.ToLower(raw), "[]"), ".")
}

// Record accepts only high-confidence explicit fingerprints.
func (r *HoneypotRegistry) Record(a db.HoneypotAssessment) bool {
	if r == nil || a.Score < 90 || a.Level != "suspected" || len(a.Signals) == 0 {
		return false
	}
	a.Host = normalizeHoneypotTarget(a.Host)
	if a.Host == "" {
		return false
	}
	now := time.Now().UTC()
	if a.ObservedAt.IsZero() {
		a.ObservedAt = now
	}
	if a.ExpiresAt.IsZero() {
		a.ExpiresAt = now.Add(db.HoneypotAssessmentTTL)
	}
	if !a.ExpiresAt.After(now) {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.hosts == nil {
		r.hosts = make(map[string]db.HoneypotAssessment)
	}
	old, ok := r.hosts[a.Host]
	if ok && old.ExpiresAt.After(now) {
		seen := make(map[string]bool, len(old.Signals)+len(a.Signals))
		merged := make([]string, 0, len(old.Signals)+len(a.Signals))
		for _, signal := range append(old.Signals, a.Signals...) {
			if signal != "" && !seen[signal] {
				seen[signal] = true
				merged = append(merged, signal)
			}
		}
		sort.Strings(merged)
		a.Signals = merged
		if old.Score > a.Score {
			a.Score = old.Score
		}
		if old.ExpiresAt.After(a.ExpiresAt) {
			a.ExpiresAt = old.ExpiresAt
		}
		if old.ObservedAt.After(a.ObservedAt) {
			a.ObservedAt = old.ObservedAt
		}
	}
	r.hosts[a.Host] = a
	return true
}

func (r *HoneypotRegistry) RecordAll(items []db.HoneypotAssessment) {
	for _, item := range items {
		r.Record(item)
	}
}

func (r *HoneypotRegistry) Assessment(host string) (db.HoneypotAssessment, bool) {
	if r == nil {
		return db.HoneypotAssessment{}, false
	}
	host = normalizeHoneypotTarget(host)
	if host == "" {
		return db.HoneypotAssessment{}, false
	}
	r.mu.RLock()
	a, ok := r.hosts[host]
	r.mu.RUnlock()
	if !ok || !a.ExpiresAt.After(time.Now()) {
		return db.HoneypotAssessment{}, false
	}
	a.Signals = append([]string(nil), a.Signals...)
	return a, true
}

func (r *HoneypotRegistry) IsQuarantined(host string) bool {
	_, ok := r.Assessment(host)
	return ok
}

// AllQuarantined is deliberately false for an empty target set and true only
// when every target identity attached to this intent has an active assessment.
func (r *HoneypotRegistry) AllQuarantined(hosts []string) bool {
	if r == nil || len(hosts) == 0 {
		return false
	}
	for _, host := range hosts {
		if !r.IsQuarantined(host) {
			return false
		}
	}
	return true
}

func (r *HoneypotRegistry) List() []db.HoneypotAssessment {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	out := make([]db.HoneypotAssessment, 0, len(r.hosts))
	for _, a := range r.hosts {
		if a.ExpiresAt.After(time.Now()) {
			a.Signals = append([]string(nil), a.Signals...)
			out = append(out, a)
		}
	}
	r.mu.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].Host < out[j].Host })
	return out
}

// HoneypotRuntime is attached to one task/worker context. Targets are the
// normalized asset identities owned by the current intent.
type HoneypotRuntime struct {
	Registry *HoneypotRegistry
	Targets  []string
	Stop     func()
}

type honeypotRuntimeKey struct{}

func WithHoneypotRuntime(ctx context.Context, runtime *HoneypotRuntime) context.Context {
	if runtime == nil {
		return ctx
	}
	return context.WithValue(ctx, honeypotRuntimeKey{}, runtime)
}

func HoneypotRuntimeFrom(ctx context.Context) *HoneypotRuntime {
	if ctx == nil {
		return nil
	}
	runtime, _ := ctx.Value(honeypotRuntimeKey{}).(*HoneypotRuntime)
	return runtime
}

// IntentTargetHosts resolves the host identities of the assets anchored to an
// intent. URLs, domains and IPs are retained as aliases so scanner observations
// can be matched whether they identify a host by name or address.
func IntentTargetHosts(as *db.AssetStore, intent *db.Node) []string {
	if as == nil || intent == nil || len(intent.Anchors) == 0 {
		return nil
	}
	assets, err := as.GetByIDs(intent.Anchors)
	if err != nil {
		return nil
	}
	return AssetTargetHosts(assets)
}

func AssetTargetHosts(assets []*db.Asset) []string {
	seen := make(map[string]bool)
	var hosts []string
	add := func(raw string) {
		host := normalizeHoneypotTarget(raw)
		if host != "" && !seen[host] {
			seen[host] = true
			hosts = append(hosts, host)
		}
	}
	for _, asset := range assets {
		if asset == nil {
			continue
		}
		if asset.URL != "" {
			if u, err := url.Parse(asset.URL); err == nil {
				add(u.Hostname())
			}
		}
		add(asset.Domain)
		add(asset.IP)
		add(asset.RootDomain)
		for _, domain := range asset.BoundDomains {
			add(domain)
		}
	}
	sort.Strings(hosts)
	return hosts
}

func (t *ToolSet) quarantinedAssets(ids []int64) ([]db.HoneypotAssessment, bool) {
	if t == nil || t.as == nil || len(ids) == 0 {
		return nil, false
	}
	assets, err := t.as.GetByIDs(ids)
	if err != nil {
		return nil, false
	}
	hosts := AssetTargetHosts(assets)
	if len(hosts) == 0 {
		return nil, false
	}
	assessments, err := t.as.HoneypotAssessments(t.taskID, hosts)
	if err != nil {
		return nil, false
	}
	if t.honeypotRuntime != nil && t.honeypotRuntime.Registry != nil {
		t.honeypotRuntime.Registry.RecordAll(assessments)
		if t.honeypotRuntime.Registry.AllQuarantined(hosts) {
			return assessments, true
		}
		return assessments, false
	}
	byHost := make(map[string]bool, len(assessments))
	for _, a := range assessments {
		byHost[a.Host] = true
	}
	for _, host := range hosts {
		if !byHost[host] {
			return assessments, false
		}
	}
	return assessments, true
}

type honeypotGuardedTool struct {
	actool.CoreTool
	runtime *HoneypotRuntime
}

func (g *honeypotGuardedTool) Call(ctx context.Context, in json.RawMessage, tc *actool.ToolContext) (actool.Result, error) {
	if g.runtime != nil && g.runtime.Registry != nil && g.runtime.Registry.AllQuarantined(g.runtime.Targets) {
		if g.runtime.Stop != nil {
			g.runtime.Stop()
		}
		return actool.Errorf(ErrHoneypotTargetQuarantined.Error() + "; this intent is stopped to avoid repeated tool and LLM work"), nil
	}
	return g.CoreTool.Call(ctx, in, tc)
}

func guardHoneypotTools(tools []actool.CoreTool, runtime *HoneypotRuntime) []actool.CoreTool {
	if runtime == nil || runtime.Registry == nil || len(runtime.Targets) == 0 {
		return tools
	}
	out := make([]actool.CoreTool, len(tools))
	for i, tool := range tools {
		out[i] = &honeypotGuardedTool{CoreTool: tool, runtime: runtime}
	}
	return out
}
