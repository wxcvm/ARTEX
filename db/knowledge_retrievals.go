package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type KnowledgeRetrievalHit struct {
	ID string `json:"id"`
	Scope string `json:"scope"`
	Title string `json:"title"`
}

type KnowledgeRetrieval struct {
	ID string `json:"id"`
	AgentKey string `json:"agent_key"`
	Query string `json:"query"`
	Hits []KnowledgeRetrievalHit `json:"hits"`
	Degraded bool `json:"degraded"`
	CreatedAt time.Time `json:"created_at"`
}

// RecordKnowledgeRetrieval stores query and returned reference identities only.
// Task history shares the task archive lock and never becomes global history.
func (d *DB) RecordKnowledgeRetrieval(ctx context.Context,taskID int64,agentKey,query string,hits []KnowledgeRetrievalHit,degraded bool)error{
	if hits==nil{hits=[]KnowledgeRetrievalHit{}}
	raw,err:=json.Marshal(hits)
	if err!=nil{return err}
	return d.WithEvidenceTx(ctx,func(tx *sql.Tx)error{
		if err:=LockTaskEvidenceTx(tx,taskID);err!=nil{return err}
		_,err:=tx.ExecContext(ctx,`INSERT INTO knowledge_retrievals(id,task_id,agent_key,query,hits,degraded)
VALUES($1,NULLIF($2::bigint,0),$3,$4,$5::jsonb,$6)`,uuid.NewString(),taskID,agentKey,query,string(raw),degraded)
		return err
	})
}

func (d *DB) ListKnowledgeRetrievals(ctx context.Context,taskID int64)([]KnowledgeRetrieval,error){
	rows,err:=d.QueryContext(ctx,`SELECT id,agent_key,query,hits,degraded,created_at FROM knowledge_retrievals
WHERE task_id IS NOT DISTINCT FROM NULLIF($1::bigint,0) ORDER BY created_at DESC,id LIMIT 50`,taskID)
	if err!=nil{return nil,err}
	defer rows.Close()
	out:=[]KnowledgeRetrieval{}
	for rows.Next(){
		var item KnowledgeRetrieval
		var raw []byte
		if err:=rows.Scan(&item.ID,&item.AgentKey,&item.Query,&raw,&item.Degraded,&item.CreatedAt);err!=nil{return nil,err}
		if err:=json.Unmarshal(raw,&item.Hits);err!=nil{return nil,err}
		out=append(out,item)
	}
	return out,rows.Err()
}

func restoreKnowledgeRetrievalsTx(tx *sql.Tx,snapshot *TaskArchiveSnapshot)error{
	raw:=snapshot.Tables["knowledge_retrievals"]
	rows,err:=decodeArchiveRows(raw)
	if err!=nil{return err}
	for _,row:=range rows{
		tid,ok:=jsonInt64(row["task_id"])
		if !ok||tid!=snapshot.TaskID{return errors.New("archived retrieval belongs to another scope")}
	}
	return insertArchiveRows(tx,"knowledge_retrievals",raw)
}
