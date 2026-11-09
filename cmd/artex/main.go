// Command artex runs the ARTEX backend: the dual SQLite graph stores,
// the event-driven exploration engine, and the JSON HTTP API consumed by the
// shadcn/ui frontend.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/Autumn-27/artex/agent"
	"github.com/Autumn-27/artex/config"
	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/server"
)

// version is the build version, injected at release time via
// -ldflags "-X main.version=<tag>". Defaults to "dev" for local builds.
var version = "dev"

const banner = `
    _    ____ _____ _______  __
   / \  |  _ \_   _| ____\ \/ /
  / _ \ | |_) || | |  _|  \  /
 / ___ \|  _ < | | | |___ /  \
/_/   \_\_| \_\|_| |_____/_/\_\
`

// printBanner writes the startup banner + version/runtime info to stdout.
func printBanner(addr string) {
	fmt.Print(banner)
	fmt.Println("  AI autonomous penetration testing system")
	fmt.Printf("  version %s  ·  %s/%s  ·  %s  ·  listening on %s\n\n",
		version, runtime.GOOS, runtime.GOARCH, runtime.Version(), addr)
}

// main maps run's result onto the process exit code. The body lives in a
// function that can *return* rather than os.Exit, so its own defers still run.
func main() {
	os.Exit(run())
}

func run() int {
	var (
		addr    = flag.String("addr", ":8787", "HTTP listen address")
		dataDir = flag.String("data", filepath.Join(config.BaseDir(), "data"), "data directory for SQLite stores (default: data/ next to the executable)")
		proxy   = flag.String("proxy", "127.0.0.1:8788", "traffic recording proxy address (empty to disable)")
		checkDB = flag.Bool("check-db", false, "check PostgreSQL connectivity without starting the server or applying schema changes")
		createDB = flag.Bool("create-db", false, "explicitly create the configured PostgreSQL database if it is missing; do not apply schema")
		setupAppDB = flag.Bool("setup-app-db", false, "create a dedicated login role and database owned by it (new setups only)")
	)
	flag.Parse()
	if (*checkDB && *createDB) || (*checkDB && *setupAppDB) || (*createDB && *setupAppDB) {
		fmt.Fprintln(os.Stderr, "use only one database action at a time")
		return 2
	}
	if *setupAppDB {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		dsn, _, err := db.DSN()
		if err == nil {
			err = provisionDedicatedDatabase(ctx, dsn,
				os.Getenv("ARTEX_APP_DB_USER"), os.Getenv("ARTEX_APP_DB_PASSWORD"))
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "PostgreSQL app role setup failed: %v\n", err)
			return 1
		}
		fmt.Println("PostgreSQL dedicated role and database ready")
		return 0
	}
	if *createDB {
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()
		created, err := createConfiguredDatabase(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "PostgreSQL database creation failed: %v\n", err)
			return 1
		}
		if created {
			fmt.Println("PostgreSQL database created")
		} else {
			fmt.Println("PostgreSQL database already exists (unchanged)")
		}
		return 0
	}
	if *checkDB {
		ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
		defer cancel()
		if err := checkDatabaseConnection(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "PostgreSQL connection failed: %v\n", err)
			return 1
		}
		fmt.Println("PostgreSQL connection OK")
		return 0
	}

	// hand the build version to the server package so GET /api/health can report it
	// to the frontend top bar.
	server.BuildVersion = version

	printBanner(*addr)

	// capture backend logs into the in-memory sink (still to stderr) so the /logs
	// page can show a live log stream. Do this first, to catch startup logs too.
	server.StartLogCapture()

	// surface which config file the binary reads (absolute, so `go run`'s relative
	// "config.json" — resolved against the CWD — is unambiguous).
	cfgPath := config.Path()
	if abs, e := filepath.Abs(cfgPath); e == nil {
		cfgPath = abs
	}
	if _, e := os.Stat(cfgPath); e == nil {
		log.Printf("[config] config file: %s", cfgPath)
	} else {
		log.Printf("[config] config file: %s (not found — will only try the environment variable ARTEX_PG_DSN)", cfgPath)
	}

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, shutdown := shutdownContext(sigCtx)
	defer shutdown(agent.AbortShutdown)

	mgr, err := server.NewManager(*dataDir, *proxy)
	if err != nil {
		log.Fatalf("open stores: %v", err)
	}
	defer mgr.Close()

	skillDir := config.SkillDir()
	if abs, err := filepath.Abs(skillDir); err == nil {
		skillDir = abs
	}
	log.Printf("[config] skill directory: %s", skillDir)
	jwtKeyDir := config.BaseDir()
	if override := os.Getenv("ARTEX_JWT_KEY_DIR"); override != "" {
		jwtKeyDir = override
	}
	srv := server.New(ctx, mgr, skillDir, *dataDir, jwtKeyDir)
	httpSrv := &http.Server{
		Addr:              *addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("ARTEX %s backend listening on %s (data=%s, workers=%d)", version, *addr, *dataDir, mgr.Workers())
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve: %v", err)
		}
	}()

	// Wait for a stop signal.
	<-ctx.Done()

	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdownCtx)
	return 0
}

// shutdownContext deliberately does not derive from signalCtx. If it did, the
// parent's plain context.Canceled could win the race before AbortShutdown was
// attached to the child, losing the diagnostic cause in every running Agent.
func shutdownContext(signalCtx context.Context) (context.Context, context.CancelCauseFunc) {
	ctx, shutdown := context.WithCancelCause(context.Background())
	go func() {
		select {
		case <-signalCtx.Done():
			shutdown(agent.AbortShutdown)
		case <-ctx.Done():
		}
	}()
	return ctx, shutdown
}

// checkDatabaseConnection only authenticates and pings. Unlike db.Open, it does
// not create databases, apply migrations, or modify any existing records.
func checkDatabaseConnection(ctx context.Context) error {
	dsn, _, err := db.DSN()
	if err != nil {
		return err
	}
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.PingContext(ctx)
}
