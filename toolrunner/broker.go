package toolrunner

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ServeBroker owns the controller-facing Unix socket but executes no tools.
// The actual tool executor is on a separate Docker container with no mount
// of the broker's socket, credentials, or control-plane data directory.
func ServeBroker(ctx context.Context, socket, target string) error {
	u, err := url.Parse(target)
	if err != nil || u == nil || u.Scheme != "http" || u.Hostname() == "" || u.Port() == "" || u.User != nil || u.Path != "" {
		return errors.New("runner broker requires an http://host:port executor target")
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok := r.Method == http.MethodPost && (r.URL.Path == "/call" || r.URL.Path == "/close")
		ok = ok || (r.Method == http.MethodGet && r.URL.Path == "/health")
		if !ok {
			http.Error(w, "unsupported runner broker operation", http.StatusNotFound)
			return
		}
		proxy.ServeHTTP(w, r)
	})
	return serveOnUnix(ctx, socket, handler)
}

// ServeExecutor exposes only SDK tool operations on the private execution
// network. It does not mount the Unix control socket or controller secrets.
func ServeExecutor(ctx context.Context, listen string) error {
	if listen == "" {
		return fmt.Errorf("runner executor TCP listen address required")
	}
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	defer ln.Close()
	service := NewService()
	defer service.Close()
	return serveHTTP(ctx, ln, service.Handler())
}
