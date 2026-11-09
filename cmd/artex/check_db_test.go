package main

import (
	"context"
	"net/url"
	"os"
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
