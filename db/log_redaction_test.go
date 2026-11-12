package db

import (
	"testing"
)

func TestHistoricLogRedaction(t *testing.T) {
	d, err := Open(testDSN(t))
	if err != nil {
		t.Skipf("postgres unavailable: %v", err)
	}
	defer d.Close()
	const raw = "provider error: api_key=synthetic-test-key"
	const want = "provider error: api_key=***"
	id, err := d.InsertLog("error", "test-privacy", raw)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Exec("DELETE FROM server_logs WHERE id=$1", id)
	rows, err := d.ListLogsBefore(id+1, 1)
	if err != nil || len(rows) != 1 || rows[0].ID != id || rows[0].Text != want {
		t.Fatal("historic log pagination failed to mask credentials")
	}
	rows, err = d.RecentLogs(100)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, l := range rows {
		if l.ID == id {
			found = true
			if l.Text != want {
				t.Fatal("recent logs exposed a credential")
			}
		}
	}
	if !found {
		t.Fatal("fixture missing from recent logs")
	}
	var stored string
	if err := d.QueryRow("SELECT text FROM server_logs WHERE id=$1", id).Scan(&stored); err != nil || stored != raw {
		t.Fatal("historical source record was changed")
	}
}
