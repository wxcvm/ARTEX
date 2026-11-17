package toolrunner

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// installPublicCA puts only a validated public certificate chain in the
// runner's own temporary directory. The controller's certificate file, data
// root and CA signing key are never mounted in the runner.
func installPublicCA(env []string, chain string) ([]string, error) {
	rest := []byte(chain)
	parsed := 0
	for len(rest) > 0 {
		block, tail := pem.Decode(rest)
		if block == nil || block.Type != "CERTIFICATE" {
			if len(strings.TrimSpace(string(rest))) > 0 {
				return nil, errors.New("unexpected PEM content")
			}
			break
		}
		if _, err := x509.ParseCertificate(block.Bytes); err != nil {
			return nil, err
		}
		parsed++
		rest = tail
	}
	if parsed == 0 {
		return nil, errors.New("missing public certificate")
	}
	digest := sha256.Sum256([]byte(chain))
	path := filepath.Join(os.TempDir(), "artex-runner-ca-"+hex.EncodeToString(digest[:])+".pem")
	if err := os.WriteFile(path, []byte(chain), 0600); err != nil {
		return nil, err
	}
	out := append([]string(nil), env...)
	for i, kv := range out {
		key, _, ok := strings.Cut(kv, "=")
		if !ok {
			continue
		}
		switch key {
		case "SSL_CERT_FILE", "CURL_CA_BUNDLE", "REQUESTS_CA_BUNDLE", "NODE_EXTRA_CA_CERTS":
			out[i] = key + "=" + path
		}
	}
	return out, nil
}
