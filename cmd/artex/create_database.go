package main

import (
    "context"
    "database/sql"
    "fmt"
    "net/url"
    "strings"

    "github.com/Autumn-27/artex/db"
    "github.com/jackc/pgx/v5"
)

// createConfiguredDatabase is deliberately separate from db.Open: the
// operator-initiated setup command creates only a missing database, and never
// runs ARTEX migrations or alters roles/passwords.
func createConfiguredDatabase(ctx context.Context) (bool, error) {
    dsn, _, err := db.DSN()
    if err != nil {
        return false, err
    }
    target, name, err := maintenanceDSN(dsn)
    if err != nil {
        return false, err
    }
    admin, err := sql.Open("pgx", target)
    if err != nil {
        return false, err
    }
    defer admin.Close()
    if err := admin.PingContext(ctx); err != nil {
        return false, fmt.Errorf("connect to PostgreSQL maintenance database: %w", err)
    }
    var exists bool
    err = admin.QueryRowContext(ctx,
        "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", name).Scan(&exists)
    if err != nil {
        return false, err
    }
    if exists {
        return false, nil
    }
    if _, err := admin.ExecContext(ctx, "CREATE DATABASE "+pgx.Identifier{name}.Sanitize()); err != nil {
        return false, fmt.Errorf("create database %q: %w", name, err)
    }
    return true, nil
}

func maintenanceDSN(dsn string) (string, string, error) {
    parsed, err := url.Parse(dsn)
    if err != nil {
        return "", "", err
    }
    if (parsed.Scheme != "postgres" && parsed.Scheme != "postgresql") || parsed.Host == "" {
        return "", "", fmt.Errorf("database DSN must be a postgres:// or postgresql:// URL")
    }
    name := strings.TrimPrefix(parsed.Path, "/")
    if name == "" || strings.Contains(name, "/") {
        return "", "", fmt.Errorf("database DSN must specify one target database name")
    }
    parsed.Path = "/postgres"
    parsed.RawPath = ""
    return parsed.String(), name, nil
}
