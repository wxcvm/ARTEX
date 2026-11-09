package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"
)

// provisionDedicatedDatabase creates a login-only role and a database owned by
// that role on a fresh PostgreSQL installation. It does not alter existing role
// passwords, database owners, or application records.
func provisionDedicatedDatabase(ctx context.Context, adminDSN, appUser, appPassword string) error {
	appUser = strings.TrimSpace(appUser)
	if appUser == "" || appPassword == "" {
		return fmt.Errorf("dedicated role name and password must be supplied")
	}
	maintenance, dbName, err := maintenanceDSN(adminDSN)
	if err != nil {
		return err
	}
	if dbName == "postgres" {
		return fmt.Errorf("the postgres maintenance database cannot be provisioned for ARTEX")
	}
	admin, err := sql.Open("pgx", maintenance)
	if err != nil {
		return err
	}
	defer admin.Close()
	if err := admin.PingContext(ctx); err != nil {
		return fmt.Errorf("connect to PostgreSQL maintenance database: %w", err)
	}

	var owner string
	err = admin.QueryRowContext(ctx, "SELECT pg_get_userbyid(datdba) FROM pg_database WHERE datname=$1", dbName).Scan(&owner)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("check target database: %w", err)
	}
	dbExists := err == nil
	if dbExists && owner != appUser {
		return fmt.Errorf("database %q already belongs to role %q; existing database ownership was not changed", dbName, owner)
	}

	var roleExists bool
	if err := admin.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname=$1)", appUser).Scan(&roleExists); err != nil {
		return fmt.Errorf("check dedicated role: %w", err)
	}
	if roleExists && !dbExists {
		return fmt.Errorf("role %q already exists; existing role password was not changed", appUser)
	}

	if !roleExists {
		var statement string
		// PostgreSQL format(%I, %L) quotes role identifiers and passwords
		// correctly, including quote characters in operator-supplied passwords.
		if err := admin.QueryRowContext(ctx,
			"SELECT format('CREATE ROLE %I LOGIN PASSWORD %L', $1::text, $2::text)",
			appUser, appPassword).Scan(&statement); err != nil {
			return fmt.Errorf("prepare role: %w", err)
		}
		if _, err := admin.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("create dedicated role %q: %w", appUser, err)
		}
	}

	if !dbExists {
		statement := "CREATE DATABASE " + pgx.Identifier{dbName}.Sanitize() +
			" OWNER " + pgx.Identifier{appUser}.Sanitize()
		if _, err := admin.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("create database %q owned by %q: %w", dbName, appUser, err)
		}
	}

	// Confirm the credentials actually work and the role can create the tables
	// ARTEX seeds on initial startup, without applying the actual schema here.
	parsed, err := url.Parse(adminDSN)
	if err != nil {
		return err
	}
	parsed.User = url.UserPassword(appUser, appPassword)
	app, err := sql.Open("pgx", parsed.String())
	if err != nil {
		return err
	}
	defer app.Close()
	if err := app.PingContext(ctx); err != nil {
		return fmt.Errorf("verify dedicated role login: %w", err)
	}
	var canCreate bool
	if err := app.QueryRowContext(ctx,
		"SELECT has_schema_privilege(current_user, 'public', 'CREATE')").Scan(&canCreate); err != nil {
		return fmt.Errorf("verify schema privileges: %w", err)
	}
	if !canCreate {
		return fmt.Errorf("dedicated role %q cannot CREATE tables in public schema", appUser)
	}
	return nil
}
