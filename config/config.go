// Package config loads runtime configuration from a JSON file, with environment
// variables taking precedence. Currently it carries the PostgreSQL connection.
package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Database is the PostgreSQL connection config. Either set DSN directly, or set
// the component fields and a DSN is assembled from them.
type Database struct {
	DSN      string `json:"dsn"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

// Config is the on-disk config file shape.
type Config struct {
	Database Database `json:"database"`
	SkillDir string   `json:"skill_dir"`
	// KnowledgeDir points at a directory holding knowledge.db (the deployment
	// knowledge base). Empty = the kb_search / kb_get tools are not registered.
	KnowledgeDir string `json:"knowledge_dir"`
}

// BaseDir is the directory that anchors all runtime artifacts (config.json and
// the data/ store). It is the directory the running binary lives in, so a
// distributed executable keeps its files next to itself on any OS (Windows,
// Linux, …) regardless of the working directory it is launched from.
//
// When launched via `go run`, the binary is throwaway: it sits either in a temp
// build dir (cache miss → fresh link) OR straight inside the Go build cache
// (cache hit → run from $GOCACHE/.../...-d). We detect both and fall back to the
// current working directory so dev artifacts (data/, transcripts) and config
// resolve against the project dir, not the throwaway binary's location.
func BaseDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	dir := filepath.Dir(exe)
	if isGoRunDir(dir) {
		return "." // throwaway `go run` binary → use CWD
	}
	return dir
}

// isGoRunDir reports whether dir is where `go run` parked its executable: under
// the system temp dir (cache miss), or anywhere inside a Go build cache
// (…/go-build/…, the cache-hit case — NOT under os.TempDir(), which is why the
// old temp-only check failed intermittently). In both cases the binary is
// throwaway, so config/data must resolve against the CWD.
func isGoRunDir(dir string) bool {
	if tmp := os.TempDir(); tmp != "" {
		if rel, err := filepath.Rel(tmp, dir); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return true
		}
	}
	for _, seg := range strings.Split(filepath.ToSlash(dir), "/") {
		if seg == "go-build" {
			return true
		}
	}
	return false
}

// Path returns the config file path. Resolution order:
//  1. env ARTEX_CONFIG (explicit override)
//  2. ./config.json in the current working directory (running from the project
//     dir — robust no matter where `go run` placed the temp/cached binary)
//  3. config.json next to the executable (a distributed binary keeps it beside)
//
// The first existing file wins. If none exist, the CWD path is returned so the
// "not found" message points at the project dir the user most likely expected.
func Path() string {
	if v := strings.TrimSpace(os.Getenv("ARTEX_CONFIG")); v != "" {
		return v
	}
	var candidates []string
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(cwd, "config.json"))
	}
	candidates = append(candidates, filepath.Join(BaseDir(), "config.json"))
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return candidates[0]
}

// readConfig never returns partially decoded settings. Keep diagnostics free
// of JSON values: configuration can contain database credentials.
func readConfig(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	return decodeConfig(b, path)
}

// decodeConfig is shared by runtime loading and read-only diagnostics.
func decodeConfig(b []byte, path string) (Config, error) {
	var c *Config
	if err := json.Unmarshal(b, &c); err != nil || c == nil {
		return Config{}, fmt.Errorf("invalid JSON configuration in %s: expected a configuration object with valid field types", path)
	}
	return *c, nil
}

// Load reads optional settings. Missing, unreadable or invalid configuration
// yields a zero Config; callers requiring database configuration use readConfig.
func Load() Config {
	c, _ := readConfig(Path())
	return c
}

// SkillDir returns the skill root directory with precedence:
//
//	env ARTEX_SKILL_DIR  >  config file (skill_dir)  >  BaseDir()/skills
//
// The directory is created if it does not exist.
// KnowledgeDir returns the optional knowledge base directory with the same
// precedence as SkillDir:
//
//	env ARTEX_KB_DIR  >  config file (knowledge_dir)  >  "" (feature off)
//
// Unlike SkillDir it never creates the directory: an absent knowledge base must
// leave no empty state behind, and the tools simply stay unregistered.
func KnowledgeDir() string {
	if v := strings.TrimSpace(os.Getenv("ARTEX_KB_DIR")); v != "" {
		return v
	}
	return strings.TrimSpace(Load().KnowledgeDir)
}

func SkillDir() string {
	var d string
	if v := strings.TrimSpace(os.Getenv("ARTEX_SKILL_DIR")); v != "" {
		d = v
	} else if v := strings.TrimSpace(Load().SkillDir); v != "" {
		d = v
	} else {
		d = filepath.Join(BaseDir(), "skills")
	}
	_ = os.MkdirAll(d, 0o755)
	return d
}

// PostgresDSN resolves the connection string with precedence:
//
//	env ARTEX_PG_DSN  >  config file (database.dsn, or assembled from fields)
//
// There is NO built-in fallback: when neither source supplies a database config,
// it returns an error naming the config path it inspected, so startup fails loudly
// instead of silently connecting to a wrong default. source describes where the
// DSN came from (for startup logging).
func PostgresDSN() (dsn, source string, err error) {
	if v := strings.TrimSpace(os.Getenv("ARTEX_PG_DSN")); v != "" {
		return v, "environment variable ARTEX_PG_DSN", nil
	}
	path := Path()
	cfg, err := readConfig(path)
	if err != nil && !os.IsNotExist(err) {
		return "", "", fmt.Errorf("read database config: %w", err)
	}
	db := cfg.Database
	if d := strings.TrimSpace(db.DSN); d != "" {
		return d, "config file " + path + " (database.dsn)", nil
	}
	if db.Host != "" || db.DBName != "" || db.User != "" {
		return db.buildDSN(), "config file " + path + " (database fields)", nil
	}
	return "", "", fmt.Errorf("no database config found: the environment variable ARTEX_PG_DSN is not set, and config file %s provides no database (dsn or host/user/dbname). Create the config file or set the environment variable and retry", path)
}

func (d Database) buildDSN() string {
	host := d.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := d.Port
	if port == 0 {
		port = 5432
	}
	ssl := d.SSLMode
	if ssl == "" {
		ssl = "disable"
	}
	u := url.URL{
		Scheme: "postgres",
		Host:   host + ":" + strconv.Itoa(port),
		Path:   "/" + d.DBName,
	}
	if d.User != "" {
		if d.Password != "" {
			u.User = url.UserPassword(d.User, d.Password)
		} else {
			u.User = url.User(d.User)
		}
	}
	u.RawQuery = url.Values{"sslmode": {ssl}}.Encode()
	return u.String()
}

// Redact returns a diagnostic-only PostgreSQL address. Invalid or keyword-style
// DSNs must not be echoed, and URI query parameters can contain credentials too.
func Redact(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil || u.Opaque != "" ||
		(!strings.EqualFold(u.Scheme, "postgres") && !strings.EqualFold(u.Scheme, "postgresql")) {
		return "[redacted PostgreSQL DSN]"
	}
	u.User = nil
	u.RawQuery = ""
	u.ForceQuery = false
	u.Fragment = ""
	u.RawFragment = ""
	return u.String()
}
