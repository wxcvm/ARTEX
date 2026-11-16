package db

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestCyberStrikeReferencePackManifest(t *testing.T) {
	var pack knowledgeReferencePack
	if err := json.Unmarshal(cyberStrikeReferencePackJSON, &pack); err != nil {
		t.Fatal(err)
	}
	if pack.ID != "cyberstrike-csf-v1" || len(pack.Entries) != 9 || len(pack.Revision) != 40 || pack.UpstreamLicense != "AGPL-3.0" {
		t.Fatalf("unexpected pack metadata: %+v", pack)
	}
	seen := map[string]bool{}
	for _, e := range pack.Entries {
		id := referenceDocumentID(pack, e)
		if seen[id] || !strings.HasPrefix(e.Title, e.Code+" ") || !strings.HasPrefix(e.SourcePath, ".cyberstrike/skill/NIST/CSF_v2.0/") ||
			len(e.SourceBlob) != 40 || e.Keywords == "" || strings.TrimSpace(e.Body) == "" {
			t.Fatalf("invalid or duplicate entry: %+v", e)
		}
		seen[id] = true
		if _, err := uuid.Parse(id); err != nil { t.Fatal(err) }
		content := referenceDocumentContent(pack, e)
		for _, required := range []string{pack.Revision, e.SourceBlob, e.SourcePath, pack.Authority, "Reference only"} {
			if !strings.Contains(content, required) { t.Fatalf("missing provenance: %s", required) }
		}
	}
}

func TestKnowledgeReferencePackSeedOnceAndSearch(t *testing.T) {
	d, err := Open(testDSN(t))
	if err != nil {
		t.Skipf("postgres unavailable (%v)", err)
	}
	defer d.Close()
	ctx := context.Background()
	var pack knowledgeReferencePack
	if err := json.Unmarshal(cyberStrikeReferencePackJSON, &pack); err != nil { t.Fatal(err) }
	// Use unique identities on the shared test database, not real shipped
	// document IDs or the production marker.
	pack.ID = "test-csf-" + uuid.NewString()
	pack.Category = pack.ID
	pack.Entries = pack.Entries[:2]
	for i := range pack.Entries {
		pack.Entries[i].Title = pack.ID + " " + pack.Entries[i].Title
	}
	a,b := referenceDocumentID(pack, pack.Entries[0]),referenceDocumentID(pack,pack.Entries[1])
	defer d.Exec(`DELETE FROM knowledge_documents WHERE id IN ($1,$2)`,a,b)
	defer d.Exec(`DELETE FROM settings WHERE key=$1`,referencePackMarker(pack))
	if err := d.seedKnowledgeReferencePack(pack); err != nil { t.Fatal(err) }
	hits, err := d.SearchKnowledgeDocuments(ctx,0,pack.ID,10)
	if err != nil || len(hits)!=2 { t.Fatalf("search: %v %v",hits,err) }
	doc,err:=d.GetKnowledgeDocument(ctx,0,false,a)
	if err!=nil || doc.TaskID!="" || doc.SourceFilename!=pack.Entries[0].SourcePath || !strings.Contains(doc.Content,pack.Revision) {
		t.Fatalf("global document: %+v %v",doc,err)
	}
	doc.Content="user edited reference"
	if _,err:=d.SaveKnowledgeDocument(ctx,0,a,doc);err!=nil { t.Fatal(err) }
	if err:=d.DeleteKnowledgeDocument(ctx,0,b);err!=nil { t.Fatal(err) }
	if err:=d.seedKnowledgeReferencePack(pack);err!=nil { t.Fatal(err) }
	doc,err=d.GetKnowledgeDocument(ctx,0,false,a)
	if err!=nil || doc.Content!="user edited reference" { t.Fatalf("edit overwritten: %+v %v",doc,err) }
	if _,err:=d.GetKnowledgeDocument(ctx,0,false,b);err!=ErrKnowledgeDocumentNotFound { t.Fatalf("deleted reference reinstalled: %v",err) }
}
