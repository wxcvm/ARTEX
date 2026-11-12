package agent

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"
)

const modelAuditCapacity = 256

// ModelOutboundEvent describes one completed HTTP RoundTrip, not a TCP connection
// or a full streamed completion. Never retain headers, bodies, paths, queries,
// credentials, or raw error strings.
type ModelOutboundEvent struct {
	Sequence        uint64    `json:"sequence"`
	At              time.Time `json:"at"`
	Scheme          string    `json:"scheme"`
	Host            string    `json:"host"`
	Port            string    `json:"port"`
	ProxyConfigured bool      `json:"proxy_configured"`
	Status          int       `json:"status"`
	Result          string    `json:"result"`
	ElapsedMS       int64     `json:"elapsed_ms"`
}

type ModelOutboundSnapshot struct {
	Capacity int                  `json:"capacity"`
	Total    uint64               `json:"total"`
	Events   []ModelOutboundEvent `json:"events"`
}

type modelAudit struct {
	mu    sync.Mutex
	ring  []ModelOutboundEvent
	total uint64
}

func newModelAudit(capacity int) *modelAudit {
	return &modelAudit{ring: make([]ModelOutboundEvent, capacity)}
}

var outboundAudit = newModelAudit(modelAuditCapacity)

func (a *modelAudit) add(event ModelOutboundEvent) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.total++
	event.Sequence = a.total
	a.ring[(a.total-1)%uint64(len(a.ring))] = event
}

func (a *modelAudit) snapshot() ModelOutboundSnapshot {
	a.mu.Lock()
	defer a.mu.Unlock()
	count := min(a.total, uint64(len(a.ring)))
	events := make([]ModelOutboundEvent, 0, count)
	for i := uint64(0); i < count; i++ {
		events = append(events, a.ring[(a.total-1-i)%uint64(len(a.ring))])
	}
	return ModelOutboundSnapshot{Capacity: len(a.ring), Total: a.total, Events: events}
}

// ModelOutboundAudit returns a copy of the process-local audit snapshot.
func ModelOutboundAudit() ModelOutboundSnapshot { return outboundAudit.snapshot() }

type modelAuditTransport struct {
	base            http.RoundTripper
	audit           *modelAudit
	proxyConfigured bool
}

func (t modelAuditTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	event := ModelOutboundEvent{At: start.UTC(), ProxyConfigured: t.proxyConfigured}
	if req.URL != nil {
		event.Scheme = strings.ToLower(req.URL.Scheme)
		event.Host = req.URL.Hostname()
		event.Port = req.URL.Port()
		if event.Port == "" {
			switch event.Scheme {
			case "http":
				event.Port = "80"
			case "https":
				event.Port = "443"
			}
		}
	}
	resp, err := t.base.RoundTrip(req)
	event.ElapsedMS = time.Since(start).Milliseconds()
	if resp != nil {
		event.Status = resp.StatusCode
	}
	switch {
	case errors.Is(err, context.Canceled):
		event.Result = "canceled"
	case errors.Is(err, context.DeadlineExceeded):
		event.Result = "timeout"
	case err != nil || resp == nil:
		event.Result = "transport_error"
	default:
		event.Result = "http_response"
	}
	t.audit.add(event)
	return resp, err
}
