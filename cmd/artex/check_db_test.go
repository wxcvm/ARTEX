package main

import (
	"context"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCheckDatabaseConnection(t *testing.T) {
	raw := os.Getenv("ARTEX_PG_DSN")
	if raw == "" {
		t.Skip("PostgreSQL service is not configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := checkDatabaseConnection(ctx); err != nil {
		t.Fatalf("check PostgreSQL: %v", err)
	}
}

func TestCheckDatabaseConnectionWrongPassword(t *testing.T) {
	raw := os.Getenv("ARTEX_PG_DSN")
	if raw == "" {
		t.Skip("PostgreSQL service is not configured")
	}
	u, err := url.Parse(raw)
	if err != nil || u.User == nil {
		t.Skip("cannot replace password in configured test DSN")
	}
	u.User = url.UserPassword(u.User.Username(), "artex-intentionally-incorrect-password")
	t.Setenv("ARTEX_PG_DSN", u.String())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := checkDatabaseConnection(ctx); err == nil {
		t.Fatal("incorrect PostgreSQL password was accepted")
	}
}

// The explicit create-db operation should not touch an existing database.
func TestCreateConfiguredDatabaseExisting(t *testing.T) {
    if os.Getenv("ARTEX_PG_DSN") == "" {
        t.Skip("PostgreSQL service is not configured")
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    created, err := createConfiguredDatabase(ctx)
    if err != nil {
        t.Fatalf("check existing database: %v", err)
    }
    if created {
        t.Fatal("existing configured database was unexpectedly created")
    }
}

func TestMaintenanceDSN(t *testing.T) {
    target, name, err := maintenanceDSN("postgres://user:pass@localhost:5432/a%20b?sslmode=disable")
    if err != nil || name != "a b" {
        t.Fatalf("name=%q err=%v", name, err)
    }
    if !strings.Contains(target, "/postgres?") {
        t.Fatalf("unexpected maintenance DSN: %s", target)
    }
    for _, input := range []string{
        "postgres://user:pass@localhost:5432/",
        "postgres://user:pass@localhost:5432/a/b",
        "https://user:pass@localhost:5432/artex",
    } {
        if _, _, err := maintenanceDSN(input); err == nil {
            t.Fatalf("invalid DSN accepted: %s", input)
        }
    }
}
