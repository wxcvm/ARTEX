package server

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestClaimNextReportsFrontierReadFailure(t *testing.T) {
	// A closed sql.DB fails before a network connection is attempted. This
	// exercises the real DB query error path without a PostgreSQL dependency.
	conn, err := sql.Open("pgx", "postgres://localhost/artex_test")
	if err != nil { t.Fatal(err) }
	if err := conn.Close(); err != nil { t.Fatal(err) }
	store := (&db.DB{DB: conn}).Exploration(1)
	node, err := NewEngine(nil).claimNext(&Task{ID:"fixture", Store:store}, "worker")
	if err == nil || !strings.Contains(err.Error(),"list open intents") || node != nil {
		t.Fatalf("closed-DB frontier returned node=%+v error=%v; want a reported query failure",node,err)
	}
}


func TestIntentClaimTokensAreUnique(t *testing.T) {
	first, err := newIntentClaimToken()
	if err != nil {
		t.Fatal(err)
	}
	second, err := newIntentClaimToken()
	if err != nil {
		t.Fatal(err)
	}
	if first == "" || second == "" || first == second {
		t.Fatalf("claim tokens must be non-empty and unique: first=%q second=%q", first, second)
	}
}
