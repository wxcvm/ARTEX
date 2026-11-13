package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

// Exercise the actual PostgreSQL role/database permissions against the CI
// PostgreSQL service. No external installation or desktop UI is required.
func TestProvisionDedicatedDatabase(t *testing.T) {
	adminDSN := os.Getenv("ARTEX_PG_DSN")
	if adminDSN == "" {
		t.Skip("PostgreSQL service is not configured")
	}
	u, err := url.Parse(adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("artex_bootstrap_%d", time.Now().UnixNano())
	role := name + "_role"
	pass := "test'password\\with-quotes"
	u.Path = "/" + name
	targetDSN := u.String()
	maintenance, _, err := maintenanceDSN(targetDSN)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	admin, err := sql.Open("pgx", maintenance)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), "DROP DATABASE IF EXISTS "+pgx.Identifier{name}.Sanitize())
		_, _ = admin.ExecContext(context.Background(), "DROP ROLE IF EXISTS "+pgx.Identifier{role}.Sanitize())
	})
	if err := provisionDedicatedDatabase(ctx, targetDSN, role, pass); err != nil {
		t.Fatalf("initial provision: %v", err)
	}
	if err := provisionDedicatedDatabase(ctx, targetDSN, role, pass); err != nil {
		t.Fatalf("idempotent provision: %v", err)
	}
	if err := provisionDedicatedDatabase(ctx, targetDSN, role, "wrong password"); err == nil {
		t.Fatal("existing role password was unexpectedly replaced")
	}
	var super, createdb, createrole bool
	if err := admin.QueryRowContext(ctx, "SELECT rolsuper, rolcreatedb, rolcreaterole FROM pg_roles WHERE rolname=$1", role).
		Scan(&super, &createdb, &createrole); err != nil {
		t.Fatal(err)
	}
	if super || createdb || createrole {
		t.Fatal("dedicated app account has administrative privileges")
	}
	u.User = url.UserPassword(role, pass)
	app, err := sql.Open("pgx", u.String())
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()
	if _, err := app.ExecContext(ctx, "CREATE TABLE artex_bootstrap_smoke (id INT PRIMARY KEY)"); err != nil {
		t.Fatalf("dedicated account cannot create ARTEX tables: %v", err)
	}
	if _, err := app.ExecContext(ctx, "INSERT INTO artex_bootstrap_smoke VALUES (1)"); err != nil {
		t.Fatalf("dedicated account cannot write ARTEX tables: %v", err)
	}
}

func TestProvisionDoesNotAdoptExistingDatabase(t *testing.T) {
	adminDSN := os.Getenv("ARTEX_PG_DSN")
	if adminDSN == "" {
		t.Skip("PostgreSQL service is not configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	role := fmt.Sprintf("artex_uncreated_%d", time.Now().UnixNano())
	err := provisionDedicatedDatabase(ctx, adminDSN, role, "some-password")
	if err == nil || !strings.Contains(err.Error(), "existing database ownership was not changed") {
		t.Fatalf("expected existing database ownership refusal, got %v", err)
	}
	maintenance, _, err := maintenanceDSN(adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := sql.Open("pgx", maintenance)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	var exists bool
	if err := conn.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname=$1)", role).Scan(&exists); err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("role was created despite an existing database belonging to someone else")
	}
}

func TestProvisionResumesAuthenticatedRoleWithoutDatabase(t *testing.T) {
	dsn := os.Getenv("ARTEX_PG_DSN")
	if dsn == "" {
		t.Skip("PostgreSQL service is not configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("artex_resume_%d", time.Now().UnixNano())
	role := name + "_role"
	pass := "resume-only-test-password"
	u.Path = "/" + name
	target := u.String()
	maintenance, _, err := maintenanceDSN(target)
	if err != nil {
		t.Fatal(err)
	}
	admin, err := sql.Open("pgx", maintenance)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), "DROP DATABASE IF EXISTS "+pgx.Identifier{name}.Sanitize())
		_, _ = admin.ExecContext(context.Background(), "DROP ROLE IF EXISTS "+pgx.Identifier{role}.Sanitize())
	})
	var statement string
	if err := admin.QueryRowContext(ctx, "SELECT format('CREATE ROLE %I LOGIN PASSWORD %L', $1::text, $2::text)", role, pass).Scan(&statement); err != nil {
		t.Fatal(err)
	}
	if _, err := admin.ExecContext(ctx, statement); err != nil {
		t.Fatal(err)
	}
	if err := provisionDedicatedDatabase(ctx, target, role, "incorrect"); err == nil {
		t.Fatal("Unrelated role adopted without valid credentials")
	}
	var exists bool
	if err := admin.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", name).Scan(&exists); err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("Failed authentication created a database")
	}
	if err := provisionDedicatedDatabase(ctx, target, role, pass); err != nil {
		t.Fatalf("Resume after CREATE ROLE: %v", err)
	}
}
