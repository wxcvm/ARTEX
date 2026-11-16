package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sort"
	"time"
)

var ErrKnowledgeIndexStale = errors.New("document changed while indexing; rebuild from the latest version")

type KnowledgeVector struct {
	Offset int
	Snippet string
	Embedding []float32
}

type KnowledgeSemanticHit struct {
	Document KnowledgeDocument
	Offset int
	Score float64
}

// ReplaceKnowledgeVectors commits only if the document version still matches.
// The expensive local model request happens before this short transaction.
func (d *DB) ReplaceKnowledgeVectors(ctx context.Context, taskID int64, docID, key string, version time.Time, chunks []KnowledgeVector) error {
	return d.WithEvidenceTx(ctx,func(tx *sql.Tx) error {
		if err:=LockTaskEvidenceTx(tx,taskID);err!=nil{return err}
		var current time.Time
		err:=tx.QueryRowContext(ctx,`SELECT updated_at FROM knowledge_documents WHERE id=$1 AND task_id IS NOT DISTINCT FROM NULLIF($2::bigint,0) FOR UPDATE`,docID,taskID).Scan(&current)
		if errors.Is(err,sql.ErrNoRows){return ErrKnowledgeDocumentNotFound}
		if err!=nil{return err}
		if !current.Equal(version){return ErrKnowledgeIndexStale}
		if _,err=tx.ExecContext(ctx,"DELETE FROM knowledge_vectors WHERE document_id=$1",docID);err!=nil{return err}
		for i,chunk:=range chunks {
			raw,err:=json.Marshal(chunk.Embedding)
			if err!=nil{return err}
			if _,err=tx.ExecContext(ctx,`INSERT INTO knowledge_vectors(document_id,chunk_number,char_offset,snippet,document_version,index_key,embedding)
VALUES($1,$2,$3,$4,$5,$6,$7::jsonb)`,docID,i,chunk.Offset,chunk.Snippet,version,key,string(raw));err!=nil{return err}
		}
		return nil
	})
}

const knowledgeVectorScope = `(d.task_id IS NULL OR d.task_id=NULLIF($1::bigint,0)) AND v.index_key=$2 AND v.document_version=d.updated_at`

func (d *DB) HasKnowledgeVectors(ctx context.Context,taskID int64,key string)(bool,error){
	var exists bool
	err:=d.QueryRowContext(ctx,`SELECT EXISTS(SELECT 1 FROM knowledge_vectors v JOIN knowledge_documents d ON d.id=v.document_id WHERE `+knowledgeVectorScope+`)`,taskID,key).Scan(&exists)
	return exists,err
}

func (d *DB) KnowledgeVectorCount(ctx context.Context,taskID int64,docID,key string)(int,error){
	var count int
	err:=d.QueryRowContext(ctx,`SELECT count(*) FROM knowledge_vectors v JOIN knowledge_documents d ON d.id=v.document_id
WHERE d.id=$3 AND d.task_id IS NOT DISTINCT FROM NULLIF($1::bigint,0) AND v.index_key=$2 AND v.document_version=d.updated_at`,taskID,key,docID).Scan(&count)
	return count,err
}

// SearchKnowledgeVectors performs an exact normalized dot-product scan of the
// current scope, retaining only the highest scoring chunk for each document.
// No vector extension or separate service is required.
func (d *DB) SearchKnowledgeVectors(ctx context.Context,taskID int64,key string,query []float32,limit int)([]KnowledgeSemanticHit,error){
	rows,err:=d.QueryContext(ctx,`SELECT d.id,COALESCE(d.task_id::text,''),d.title,d.category,d.source_filename,v.snippet,v.char_offset,v.embedding
FROM knowledge_vectors v JOIN knowledge_documents d ON d.id=v.document_id WHERE `+knowledgeVectorScope,taskID,key)
	if err!=nil{return nil,err}
	defer rows.Close()
	best:=map[string]KnowledgeSemanticHit{}
	for rows.Next(){
		var hit KnowledgeSemanticHit
		var raw []byte
		if err:=rows.Scan(&hit.Document.ID,&hit.Document.TaskID,&hit.Document.Title,&hit.Document.Category,&hit.Document.SourceFilename,&hit.Document.Content,&hit.Offset,&raw);err!=nil{return nil,err}
		var vector []float32
		if err:=json.Unmarshal(raw,&vector);err!=nil{return nil,err}
		if len(vector)!=len(query)||len(query)==0{continue}
		for i,x:=range vector{hit.Score+=float64(x)*float64(query[i])}
		old,ok:=best[hit.Document.ID]
		if !ok||hit.Score>old.Score||(hit.Score==old.Score&&hit.Offset<old.Offset){best[hit.Document.ID]=hit}
	}
	if err:=rows.Err();err!=nil{return nil,err}
	out:=make([]KnowledgeSemanticHit,0,len(best))
	for _,hit:=range best{out=append(out,hit)}
	sort.Slice(out,func(i,j int)bool{
		if out[i].Score!=out[j].Score{return out[i].Score>out[j].Score}
		return out[i].Document.ID<out[j].Document.ID
	})
	if limit>0&&len(out)>limit{out=out[:limit]}
	return out,nil
}


type KnowledgeIndexCandidate struct {
	ID string `json:"id"`
	Title string `json:"title"`
}

// ListPendingKnowledgeIndexes pages only the requested scope, excluding
// documents already indexed for their current revision and model.
func (d *DB) ListPendingKnowledgeIndexes(ctx context.Context,taskID int64,key,after string)([]KnowledgeIndexCandidate,error){
	rows,err:=d.QueryContext(ctx,`SELECT d.id,d.title FROM knowledge_documents d
WHERE d.task_id IS NOT DISTINCT FROM NULLIF($1::bigint,0) AND d.id>$3
AND NOT EXISTS(SELECT 1 FROM knowledge_vectors v WHERE v.document_id=d.id AND v.index_key=$2 AND v.document_version=d.updated_at)
ORDER BY d.id LIMIT 50`,taskID,key,after)
	if err!=nil{return nil,err}
	defer rows.Close()
	out:=[]KnowledgeIndexCandidate{}
	for rows.Next(){
		var item KnowledgeIndexCandidate
		if err:=rows.Scan(&item.ID,&item.Title);err!=nil{return nil,err}
		out=append(out,item)
	}
	return out,rows.Err()
}

// ReusableKnowledgeVectors returns previous chunks for the same document and
// model, including old revisions. Exact text equality makes reuse independent
// of chunk position; a forced rebuild bypasses this cache.
func (d *DB) ReusableKnowledgeVectors(ctx context.Context,taskID int64,docID,key string)(map[string][]float32,error){
	rows,err:=d.QueryContext(ctx,`SELECT v.snippet,v.embedding FROM knowledge_vectors v JOIN knowledge_documents d ON d.id=v.document_id
WHERE d.id=$3 AND d.task_id IS NOT DISTINCT FROM NULLIF($1::bigint,0) AND v.index_key=$2`,taskID,key,docID)
	if err!=nil{return nil,err}
	defer rows.Close()
	out:=map[string][]float32{}
	for rows.Next(){
		var snippet string
		var raw []byte
		var vector []float32
		if err:=rows.Scan(&snippet,&raw);err!=nil{return nil,err}
		if err:=json.Unmarshal(raw,&vector);err!=nil{return nil,err}
		if len(vector)>0{out[snippet]=vector}
	}
	return out,rows.Err()
}
