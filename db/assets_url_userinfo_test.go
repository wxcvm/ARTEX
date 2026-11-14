package db

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestAssetURLUserinfoAndArchiveRestore(t *testing.T) {
	d, av2, _ := testSetup(t)
	defer d.Close()

	host := fmt.Sprintf("userinfo-%d.example", time.Now().UnixNano())
	clean := "https://" + host + "/users/@me?q=1"
	first := "https://alice:secret@" + host + "/users/@me?q=1"
	second := "https://bob:other@" + host + "/users/@me?q=1"
	defer func() {
		if _, err := d.Exec(`DELETE FROM assets WHERE url LIKE $1 OR domain=$2 OR root_domain=$2`, "%"+host+"%", host); err != nil {
			t.Errorf("cleanup: %v", err)
		}
	}()

	serviceID, err := av2.UpsertHTTPService(UpsertHTTPServiceReq{URL: first})
	if err != nil {
		t.Fatal(err)
	}
	serviceAgain, err := av2.UpsertHTTPService(UpsertHTTPServiceReq{URL: second})
	if err != nil || serviceAgain != serviceID {
		t.Fatalf("service dedup: got %d, want %d, err=%v", serviceAgain, serviceID, err)
	}
	endpointID, err := av2.UpsertEndpoint(UpsertEndpointReq{URL: first, Method: "GET"})
	if err != nil {
		t.Fatal(err)
	}
	endpointAgain, err := av2.UpsertEndpoint(UpsertEndpointReq{URL: second, Method: "GET"})
	if err != nil || endpointAgain != endpointID {
		t.Fatalf("endpoint dedup: got %d, want %d, err=%v", endpointAgain, endpointID, err)
	}

	for _, id := range []int64{serviceID, endpointID} {
		var stored string
		if err := d.QueryRow(`SELECT url FROM assets WHERE id=$1`, id).Scan(&stored); err != nil || stored != clean {
			t.Fatalf("asset %d URL=%q, want %q, err=%v", id, stored, clean, err)
		}
	}

	raw, _, err := queryArchiveRows(d, `SELECT * FROM assets WHERE id IN ($1,$2) ORDER BY id`, serviceID, endpointID)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := decodeArchiveRows(raw)
	if err != nil || len(rows) != 2 {
		t.Fatalf("archive rows=%d, err=%v", len(rows), err)
	}
	want := map[int64]int64{}
	for i, row := range rows {
		id, ok := jsonInt64(row["id"])
		if !ok {
			t.Fatal("invalid archived asset ID")
		}
		oldID := int64(990001001 + i)
		want[oldID] = id
		row["id"], row["url"] = oldID, first
	}
	archive, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := d.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	mapping, _, err := restoreArchiveAssets(tx, 880001001, archive)
	if err != nil {
		t.Fatal(err)
	}
	for oldID, id := range want {
		if mapping[oldID] != id {
			t.Fatalf("restore ID %d mapped to %d, want %d", oldID, mapping[oldID], id)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
