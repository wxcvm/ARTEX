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


func TestCyberStrikeSupplementManifest(t *testing.T) {
	var first, supplement knowledgeReferencePack
	if err := json.Unmarshal(cyberStrikeReferencePackJSON, &first); err != nil { t.Fatal(err) }
	if err := json.Unmarshal(cyberStrikeSupplementPackJSON, &supplement); err != nil { t.Fatal(err) }
	if supplement.ID != "cyberstrike-csf-supplement-v1" || len(supplement.Entries) != 9 ||
		supplement.Revision != first.Revision || supplement.UpstreamLicense != first.UpstreamLicense ||
		supplement.Framework != first.Framework || referencePackMarker(first) == referencePackMarker(supplement) {
		t.Fatalf("unexpected supplemental manifest: %+v", supplement)
	}
	codes := map[string]bool{}
	ids := map[string]bool{}
	for _, pack := range []knowledgeReferencePack{first, supplement} {
		for _, entry := range pack.Entries {
			id := referenceDocumentID(pack, entry)
			if codes[entry.Code] || ids[id] || len(entry.SourceBlob) != 40 ||
				!strings.HasPrefix(entry.SourcePath, ".cyberstrike/skill/NIST/CSF_v2.0/") ||
				!strings.HasPrefix(entry.Title, entry.Code+" ") || strings.TrimSpace(entry.Body) == "" || entry.Keywords == "" {
				t.Fatalf("invalid or duplicate reference: %+v", entry)
			}
			codes[entry.Code], ids[id] = true, true
		}
	}
}

func TestKnowledgeReferenceSupplementPreservesEarlierPack(t *testing.T) {
	d, err := Open(testDSN(t))
	if err != nil { t.Skipf("postgres unavailable (%v)", err) }
	defer d.Close()
	ctx := context.Background()
	var first, supplement knowledgeReferencePack
	if err := json.Unmarshal(cyberStrikeReferencePackJSON, &first); err != nil { t.Fatal(err) }
	if err := json.Unmarshal(cyberStrikeSupplementPackJSON, &supplement); err != nil { t.Fatal(err) }
	prefix := "test-supplement-" + uuid.NewString()
	first.ID, supplement.ID = prefix+"-first", prefix+"-second"
	first.Entries, supplement.Entries = first.Entries[:2], supplement.Entries[:1]
	a, b, c := referenceDocumentID(first, first.Entries[0]), referenceDocumentID(first, first.Entries[1]), referenceDocumentID(supplement, supplement.Entries[0])
	defer d.Exec(`DELETE FROM knowledge_documents WHERE id IN ($1,$2,$3)`,a,b,c)
	defer d.Exec(`DELETE FROM settings WHERE key IN ($1,$2)`,referencePackMarker(first),referencePackMarker(supplement))
	if err := d.seedKnowledgeReferencePack(first); err != nil { t.Fatal(err) }
	doc, err := d.GetKnowledgeDocument(ctx,0,false,a)
	if err != nil { t.Fatal(err) }
	doc.Content = "retained user edit"
	if _,err := d.SaveKnowledgeDocument(ctx,0,a,doc); err != nil { t.Fatal(err) }
	if err := d.DeleteKnowledgeDocument(ctx,0,b); err != nil { t.Fatal(err) }
	if err := d.seedKnowledgeReferencePack(first); err != nil { t.Fatal(err) }
	if err := d.seedKnowledgeReferencePack(supplement); err != nil { t.Fatal(err) }
	doc,err = d.GetKnowledgeDocument(ctx,0,false,a)
	if err != nil || doc.Content != "retained user edit" { t.Fatalf("earlier edit: %+v %v",doc,err) }
	if _,err := d.GetKnowledgeDocument(ctx,0,false,b); err != ErrKnowledgeDocumentNotFound { t.Fatalf("earlier deletion: %v",err) }
	doc,err = d.GetKnowledgeDocument(ctx,0,false,c)
	if err != nil || doc.TaskID != "" || doc.SourceFilename != supplement.Entries[0].SourcePath { t.Fatalf("supplement: %+v %v",doc,err) }
}


func TestCyberStrikeThirdPackManifest(t *testing.T) {
	var first, second, third knowledgeReferencePack
	for i, item := range []struct{ raw []byte; pack *knowledgeReferencePack }{
		{cyberStrikeReferencePackJSON,&first}, {cyberStrikeSupplementPackJSON,&second}, {cyberStrikeThirdPackJSON,&third},
	} {
		if err := json.Unmarshal(item.raw,item.pack); err != nil { t.Fatal(err) }
		if len(item.pack.Entries)!=9 { t.Fatalf("pack %d: unexpected entry count",i) }
	}
	if third.ID!="cyberstrike-csf-third-v1" || third.Revision!=first.Revision || third.Framework!=first.Framework || third.UpstreamLicense!=first.UpstreamLicense {
		t.Fatalf("third pack metadata: %+v",third)
	}
	codes, ids, markers := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _,pack := range []knowledgeReferencePack{first,second,third} {
		marker := referencePackMarker(pack)
		if markers[marker] { t.Fatal("duplicate installation marker") }
		markers[marker]=true
		for _,e := range pack.Entries {
			id:=referenceDocumentID(pack,e)
			if codes[e.Code] || ids[id] || len(e.SourceBlob)!=40 || !strings.HasPrefix(e.SourcePath,".cyberstrike/skill/NIST/CSF_v2.0/") || !strings.HasPrefix(e.Title,e.Code+" ") || strings.TrimSpace(e.Body)=="" || e.Keywords=="" {
				t.Fatalf("invalid or duplicate reference: %+v",e)
			}
			codes[e.Code],ids[id]=true,true
		}
	}
	if len(codes)!=27 { t.Fatal("unexpected total reference count") }
}


func TestCyberStrikeFourthPackManifestAndCrossPackUniqueness(t *testing.T) {
	raws := [][]byte{cyberStrikeReferencePackJSON, cyberStrikeSupplementPackJSON, cyberStrikeThirdPackJSON, cyberStrikeFourthPackJSON}
	var baseline knowledgeReferencePack
	codes, ids, markers := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for i, raw := range raws {
		var pack knowledgeReferencePack
		if err := json.Unmarshal(raw, &pack); err != nil { t.Fatal(err) }
		if i == 0 { baseline = pack }
		if i == 3 && pack.ID != "cyberstrike-csf-fourth-v1" { t.Fatalf("fourth pack: %s", pack.ID) }
		if pack.Revision != baseline.Revision || pack.Framework != baseline.Framework || pack.UpstreamLicense != baseline.UpstreamLicense { t.Fatalf("mismatched metadata %s",pack.ID) }
		if i == 3 && len(pack.Entries) != 18 { t.Fatalf("fourth pack count %d",len(pack.Entries)) }
		marker := referencePackMarker(pack)
		if markers[marker] { t.Fatalf("duplicate installation marker %s",marker) }
		markers[marker] = true
		for _, e := range pack.Entries {
			id := referenceDocumentID(pack,e)
			if codes[e.Code] || ids[id] || len(e.SourceBlob) != 40 || e.Keywords == "" || strings.TrimSpace(e.Body) == "" || !strings.HasPrefix(e.Title,e.Code+" ") || !strings.HasPrefix(e.SourcePath,".cyberstrike/skill/NIST/CSF_v2.0/") { t.Fatalf("invalid/duplicate %v",e) }
			for _, provenance := range []string{pack.Revision, e.SourceBlob, pack.Authority, "Reference only"} { if !strings.Contains(referenceDocumentContent(pack,e),provenance) { t.Fatalf("lost %s for %s",provenance,e.Code) } }
			codes[e.Code], ids[id] = true,true
		}
	}
	if len(codes) != 45 { t.Fatalf("want 45 distinct codes, got %d",len(codes)) }
}

func TestCyberStrikeFourthPackIncrementalPreservesEdits(t *testing.T) {
	d,err:=Open(testDSN(t)); if err!=nil{t.Skipf("postgres unavailable (%v)",err)}
	defer d.Close()
	ctx:=context.Background()
	var first,fourth knowledgeReferencePack
	if err:=json.Unmarshal(cyberStrikeReferencePackJSON,&first);err!=nil{t.Fatal(err)}
	if err:=json.Unmarshal(cyberStrikeFourthPackJSON,&fourth);err!=nil{t.Fatal(err)}
	prefix:="test-fourth-"+uuid.NewString()
	first.ID,fourth.ID=prefix+"-first",prefix+"-fourth"
	first.Entries,fourth.Entries=first.Entries[:2],fourth.Entries[:2]
	a,b:=referenceDocumentID(first,first.Entries[0]),referenceDocumentID(first,first.Entries[1])
	c,e:=referenceDocumentID(fourth,fourth.Entries[0]),referenceDocumentID(fourth,fourth.Entries[1])
	defer d.Exec(`DELETE FROM knowledge_documents WHERE id IN ($1,$2,$3,$4)`,a,b,c,e)
	defer d.Exec(`DELETE FROM settings WHERE key IN ($1,$2)`,referencePackMarker(first),referencePackMarker(fourth))
	if err:=d.seedKnowledgeReferencePack(first);err!=nil{t.Fatal(err)}
	doc,err:=d.GetKnowledgeDocument(ctx,0,false,a);if err!=nil{t.Fatal(err)}
	doc.Content="user-edited-first-pack"
	if _,err:=d.SaveKnowledgeDocument(ctx,0,a,doc);err!=nil{t.Fatal(err)}
	if err:=d.DeleteKnowledgeDocument(ctx,0,b);err!=nil{t.Fatal(err)}
	if err:=d.seedKnowledgeReferencePack(fourth);err!=nil{t.Fatal(err)}
	if err:=d.seedKnowledgeReferencePack(first);err!=nil{t.Fatal(err)}
	if err:=d.seedKnowledgeReferencePack(fourth);err!=nil{t.Fatal(err)}
	doc,err=d.GetKnowledgeDocument(ctx,0,false,a);if err!=nil||doc.Content!="user-edited-first-pack"{t.Fatalf("overwrite older edit: %+v %v",doc,err)}
	if _,err:=d.GetKnowledgeDocument(ctx,0,false,b);err!=ErrKnowledgeDocumentNotFound{t.Fatalf("reinstalled deleted reference: %v",err)}
	for _,id:=range []string{c,e}{doc,err=d.GetKnowledgeDocument(ctx,0,false,id);if err!=nil||doc.TaskID!=""||!strings.Contains(doc.Content,fourth.Revision){t.Fatalf("new pack %s: %+v %v",id,doc,err)}}
}
