package db

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFindingEvidenceLockNamespace(t *testing.T) {
	for _, reserved := range []int64{7337741001, 7337741002, 7337741003} {
		if findingEvidenceLockKey == reserved {
			t.Fatal("evidence lock collides with migration/test/company lock")
		}
	}
}

func TestFindingTrafficLegacyArchiveDefaults(t *testing.T) {
	d, err := Open(testDSN(t))
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	if err = d.EnsureLLMRecordsTable(); err != nil {
		t.Fatal(err)
	}
	if err = d.EnsureLLMUsageTable(); err != nil {
		t.Fatal(err)
	}
	for _, version := range []int{1, 2} {
		t.Run(fmt.Sprint(version), func(t *testing.T) {
			task, err := d.CreateTask("legacy archive evidence defaults", "fixture", nil, 0, 0)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { d.Exec(`DELETE FROM task_archives WHERE task_id=$1`, task.ID); d.DeleteTask(task.ID) }()
			f, err := d.Exploration(task.ExplorationID).RecordFinding(t.Context(), RecordFindingInput{TaskID: task.ID, ExplorationID: task.ExplorationID, Summary: "legacy", Severity: "low"})
			if err != nil {
				t.Fatal(err)
			}
			if err = d.SetPaused(task.ID, true); err != nil {
				t.Fatal(err)
			}
			archive, err := d.QueueTaskArchive(task.ID)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = d.ClaimTaskArchiveJob(t.Context()); err != nil {
				t.Fatal(err)
			}
			snapshot, err := d.SnapshotTaskArchive(task.ID)
			if err != nil {
				t.Fatal(err)
			}
			if err = d.CompleteTaskArchive(archive.ID, snapshot, "/tmp/legacy-test.tar.zst", "fixture", 1, 1); err != nil {
				t.Fatal(err)
			}
			snapshot.FormatVersion = version
			delete(snapshot.Tables, "finding_traffic_bindings")
			delete(snapshot.Tables, "traffic_evidence_snapshots")
			var rows []map[string]any
			if err = json.Unmarshal(snapshot.Tables["findings"], &rows); err != nil {
				t.Fatal(err)
			}
			for _, row := range rows {
				delete(row, "evidence_version")
				delete(row, "report_evidence_version")
			}
			snapshot.Tables["findings"], _ = json.Marshal(rows)
			if _, err = d.QueueTaskArchiveRestore(archive.ID); err != nil {
				t.Fatal(err)
			}
			if _, err = d.ClaimTaskArchiveJob(t.Context()); err != nil {
				t.Fatal(err)
			}
			if _, err = d.RestoreTaskArchive(archive.ID, snapshot, 0); err != nil {
				t.Fatal(err)
			}
			got, err := d.GetFindingTraffic(t.Context(), f.FindingID)
			if err != nil {
				t.Fatal(err)
			}
			if got.Version != 0 || got.ReportVersion != 0 || len(got.Bindings) != 0 {
				t.Fatal(got)
			}
		})
	}
}


func TestHARPageSnapshotLegacyHashAndArchive(t *testing.T) {
	// Existing v3 archives lack these optional fields. Their original hash
	// must not change when decoded under the new schema.
	var legacy TrafficEvidenceSnapshot
	if err := json.Unmarshal([]byte(`{"id":"","source_traffic_id":"har-old","captured_at":1,"url":"https://fixture.local/","method":"GET","status":200,"content_type":"","req_hash":"a","resp_hash":"b","req_len":0,"resp_len":2,"created_at":"0001-01-01T00:00:00Z"}`),&legacy); err != nil { t.Fatal(err) }
	oldID := TrafficSnapshotID(legacy)
	legacy.ID = oldID
	raw,err := json.Marshal(legacy)
	if err != nil { t.Fatal(err) }
	var fields map[string]any
	if err := json.Unmarshal(raw,&fields); err != nil { t.Fatal(err) }
	if _,ok := fields["har_page_state"]; ok { t.Fatal("legacy state serialized") }
	if _,ok := fields["har_page_ordinal"]; ok { t.Fatal("legacy ordinal serialized") }
	current := legacy
	current.HARPageState,current.HARPageOrdinal = "unique",2
	current.ID = TrafficSnapshotID(current)
	if current.ID == oldID { t.Fatal("new provenance missing from hash") }
	rows,err := json.Marshal([]TrafficEvidenceSnapshot{legacy,current})
	if err != nil { t.Fatal(err) }
	archive := &TaskArchiveSnapshot{FormatVersion:3,Tables:map[string]json.RawMessage{"traffic_evidence_snapshots":rows}}
	got,err := ArchiveEvidenceSnapshots(archive)
	if err != nil || len(got)!=2 { t.Fatalf("archive: %+v %v",got,err) }
	for _,v := range got {
		if v.ID != TrafficSnapshotID(v) { t.Fatal("archive hash mismatch") }
	}
	if got[0].HARPageState!="" || got[1].HARPageState!="unique" || got[1].HARPageOrdinal!=2 { t.Fatal(got) }
}

func TestCaptureGroupSnapshotHashOptional(t *testing.T) {
    var legacy TrafficEvidenceSnapshot
    if err:=json.Unmarshal([]byte(`{"id":"","source_traffic_id":"capture-old","captured_at":1,"url":"https://fixture.local/","method":"GET","status":200,"content_type":"","req_hash":"a","resp_hash":"b","req_len":0,"resp_len":2,"created_at":"0001-01-01T00:00:00Z"}`),&legacy);err!=nil{t.Fatal(err)}
    legacyID:=TrafficSnapshotID(legacy)
    raw,err:=json.Marshal(legacy)
    if err!=nil{t.Fatal(err)}
    var fields map[string]any
    if err:=json.Unmarshal(raw,&fields);err!=nil{t.Fatal(err)}
    if _,ok:=fields["capture_id"];ok{t.Fatal("legacy hash serialized absent capture")}
    newer:=legacy
    newer.CaptureID="capture-new"
    if TrafficSnapshotID(newer)==legacyID{t.Fatal("capture group not covered by snapshot hash")}
    newer.ID=TrafficSnapshotID(newer)
    rows,err:=json.Marshal([]TrafficEvidenceSnapshot{newer})
    if err!=nil{t.Fatal(err)}
    archive:=&TaskArchiveSnapshot{FormatVersion:3,Tables:map[string]json.RawMessage{"traffic_evidence_snapshots":rows}}
    got,err:=ArchiveEvidenceSnapshots(archive)
    if err!=nil||len(got)!=1||got[0].CaptureID!="capture-new"{t.Fatalf("archive: %+v %v",got,err)}
}
