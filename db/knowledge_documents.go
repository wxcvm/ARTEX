package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrKnowledgeDocumentNotFound = errors.New("knowledge document not found")

type KnowledgeDocument struct {
	ID string `json:"id"`
	TaskID string `json:"task_id,omitempty"`
	Title string `json:"title"`
	Category string `json:"category"`
	Content string `json:"content,omitempty"`
	SourceFilename string `json:"source_filename,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type KnowledgeDocumentPage struct {
	Items []KnowledgeDocument `json:"items"`
	HasMore bool `json:"has_more"`
}

const knowledgeDocumentCols = `id, COALESCE(task_id::text,''), title, category, content, source_filename, created_at, updated_at`
const knowledgeScope = `(task_id IS NOT DISTINCT FROM NULLIF($1::bigint,0) OR ($2::boolean AND task_id IS NULL))`

func scanKnowledgeDocument(row interface{ Scan(...any) error }) (KnowledgeDocument,error) {
	var doc KnowledgeDocument
	err := row.Scan(&doc.ID,&doc.TaskID,&doc.Title,&doc.Category,&doc.Content,&doc.SourceFilename,&doc.CreatedAt,&doc.UpdatedAt)
	if errors.Is(err,sql.ErrNoRows) { err = ErrKnowledgeDocumentNotFound }
	return doc,err
}

// ListKnowledgeDocuments is literal, case-insensitive search. A global view
// never reads task-owned data; a task can explicitly include global documents.
func (d *DB) ListKnowledgeDocuments(ctx context.Context, taskID int64, includeGlobal bool, query, category string, offset int) (KnowledgeDocumentPage,error) {
	out := KnowledgeDocumentPage{Items: []KnowledgeDocument{}}
	rows,err := d.QueryContext(ctx, `SELECT id, COALESCE(task_id::text,''), title, category,
CASE WHEN $3='' THEN left(content,240) ELSE substring(content from greatest(strpos(lower(content),lower($3))-80,1) for 240) END,
source_filename,created_at,updated_at
FROM knowledge_documents WHERE `+knowledgeScope+`
AND ($3='' OR strpos(lower(title),lower($3))>0 OR strpos(lower(category),lower($3))>0 OR strpos(lower(content),lower($3))>0)
AND ($4='' OR category=$4)
ORDER BY CASE WHEN task_id IS NOT NULL THEN 0 ELSE 1 END, updated_at DESC, id
LIMIT 51 OFFSET $5`, taskID,includeGlobal,query,category,offset)
	if err != nil { return out,err }
	defer rows.Close()
	for rows.Next() {
		doc,err := scanKnowledgeDocument(rows)
		if err != nil { return out,err }
		if len(out.Items)==50 { out.HasMore=true; break }
		out.Items=append(out.Items,doc)
	}
	return out,rows.Err()
}

func (d *DB) GetKnowledgeDocument(ctx context.Context, taskID int64, includeGlobal bool, id string) (KnowledgeDocument,error) {
	return scanKnowledgeDocument(d.QueryRowContext(ctx,`SELECT `+knowledgeDocumentCols+` FROM knowledge_documents WHERE `+knowledgeScope+` AND id=$3`,taskID,includeGlobal,id))
}

func (d *DB) SaveKnowledgeDocument(ctx context.Context, taskID int64, id string, doc KnowledgeDocument) (KnowledgeDocument,error) {
	var out KnowledgeDocument
	err := d.WithEvidenceTx(ctx,func(tx *sql.Tx) error {
		if err:=LockTaskEvidenceTx(tx,taskID); err!=nil {return err}
		var err error
		if id=="" {
			out,err=scanKnowledgeDocument(tx.QueryRowContext(ctx,`INSERT INTO knowledge_documents(id,task_id,title,category,content,source_filename)
VALUES($1,NULLIF($2::bigint,0),$3,$4,$5,$6) RETURNING `+knowledgeDocumentCols,uuid.NewString(),taskID,doc.Title,doc.Category,doc.Content,doc.SourceFilename))
		} else {
			out,err=scanKnowledgeDocument(tx.QueryRowContext(ctx,`UPDATE knowledge_documents SET title=$3,category=$4,content=$5,updated_at=now()
WHERE id=$1 AND task_id IS NOT DISTINCT FROM NULLIF($2::bigint,0) RETURNING `+knowledgeDocumentCols,id,taskID,doc.Title,doc.Category,doc.Content))
		}
		return err
	})
	return out,err
}

func (d *DB) DeleteKnowledgeDocument(ctx context.Context,taskID int64,id string) error {
	return d.WithEvidenceTx(ctx,func(tx *sql.Tx) error {
		if err:=LockTaskEvidenceTx(tx,taskID); err!=nil {return err}
		res,err:=tx.ExecContext(ctx,`DELETE FROM knowledge_documents WHERE id=$1 AND task_id IS NOT DISTINCT FROM NULLIF($2::bigint,0)`,id,taskID)
		if err!=nil{return err}
		n,err:=res.RowsAffected()
		if err!=nil{return err}
		if n==0{return ErrKnowledgeDocumentNotFound}
		return nil
	})
}

func restoreKnowledgeDocumentsTx(tx *sql.Tx,snapshot *TaskArchiveSnapshot) error {
	raw:=snapshot.Tables["knowledge_documents"]
	rows,err:=decodeArchiveRows(raw)
	if err!=nil{return err}
	for _,row:=range rows {
		tid,ok:=jsonInt64(row["task_id"])
		if !ok||tid!=snapshot.TaskID{return errors.New("archived knowledge document belongs to another scope")}
	}
	return insertArchiveRows(tx,"knowledge_documents",raw)
}
