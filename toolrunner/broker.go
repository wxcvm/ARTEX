package toolrunner

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"time"
)

// ServeBroker forwards control-plane Unix RPC to a distinct execution Unix
// socket. The executor never sees the controller-facing socket; neither
// process exposes a TCP tool-management endpoint.
func ServeBroker(ctx context.Context, socket, executorSocket string) error {
	if executorSocket == "" || !filepath.IsAbs(executorSocket) || executorSocket == socket {
		return errors.New("broker requires a distinct absolute executor Unix socket")
	}
	dialer := &net.Dialer{Timeout: 3 * time.Second}
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "executor"})
	proxy.Transport = &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, "unix", executorSocket)
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok := r.Method == http.MethodPost && (r.URL.Path == "/call" || r.URL.Path == "/egress" || r.URL.Path == "/close")
		ok = ok || (r.Method == http.MethodGet && r.URL.Path == "/health")
		if !ok {
			http.Error(w, "unsupported runner broker operation", http.StatusNotFound)
			return
		}
		proxy.ServeHTTP(w, r)
	})
	return serveOnUnix(ctx, socket, handler)
}
