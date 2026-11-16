package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrReportSnapshotNotFound=errors.New("report snapshot not found")

type TaskReportSnapshot struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Format string `json:"format"`
	Scope string `json:"scope"`
	Locale string `json:"locale"`
	TemplateName string `json:"template_name"`
	Size int64 `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
const reportSnapshotCols=`id,name,format,scope,locale,template_name,octet_length(payload),created_at`
func scanReportSnapshot(row interface{Scan(...any)error})(TaskReportSnapshot,error){
	var item TaskReportSnapshot
	err:=row.Scan(&item.ID,&item.Name,&item.Format,&item.Scope,&item.Locale,&item.TemplateName,&item.Size,&item.CreatedAt)
	if errors.Is(err,sql.ErrNoRows){err=ErrReportSnapshotNotFound}
	return item,err
}
func(d *DB) SaveTaskReportSnapshot(ctx context.Context,taskID int64,item TaskReportSnapshot,payload []byte)(TaskReportSnapshot,error){
	var out TaskReportSnapshot
	err:=d.WithEvidenceTx(ctx,func(tx *sql.Tx)error{
		if err:=LockTaskEvidenceTx(tx,taskID);err!=nil{return err}
		var err error
		out,err=scanReportSnapshot(tx.QueryRowContext(ctx,`INSERT INTO task_report_snapshots(id,task_id,name,format,scope,locale,template_name,payload)
VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING `+reportSnapshotCols,uuid.NewString(),taskID,item.Name,item.Format,item.Scope,item.Locale,item.TemplateName,payload))
		return err
	})
	return out,err
}
func(d *DB) ListTaskReportSnapshots(ctx context.Context,taskID int64,offset int)([]TaskReportSnapshot,error){
	rows,err:=d.QueryContext(ctx,"SELECT "+reportSnapshotCols+" FROM task_report_snapshots WHERE task_id=$1 ORDER BY created_at DESC,id LIMIT 50 OFFSET $2",taskID,offset)
	if err!=nil{return nil,err}
	defer rows.Close()
	out:=[]TaskReportSnapshot{}
	for rows.Next(){item,err:=scanReportSnapshot(rows);if err!=nil{return nil,err};out=append(out,item)}
	return out,rows.Err()
}
func(d *DB) ReadTaskReportSnapshot(ctx context.Context,taskID int64,id string)(string,[]byte,error){
	var format string
	var payload []byte
	err:=d.QueryRowContext(ctx,"SELECT format,payload FROM task_report_snapshots WHERE task_id=$1 AND id=$2",taskID,id).Scan(&format,&payload)
	if errors.Is(err,sql.ErrNoRows){err=ErrReportSnapshotNotFound}
	return format,payload,err
}
func(d *DB) DeleteTaskReportSnapshot(ctx context.Context,taskID int64,id string)error{
	return d.WithEvidenceTx(ctx,func(tx *sql.Tx)error{
		if err:=LockTaskEvidenceTx(tx,taskID);err!=nil{return err}
		res,err:=tx.ExecContext(ctx,"DELETE FROM task_report_snapshots WHERE task_id=$1 AND id=$2",taskID,id)
		if err!=nil{return err}
		n,err:=res.RowsAffected()
		if err==nil&&n==0{return ErrReportSnapshotNotFound}
		return err
	})
}
func restoreTaskReportSnapshotsTx(tx *sql.Tx,snapshot *TaskArchiveSnapshot)error{
	raw:=snapshot.Tables["task_report_snapshots"]
	rows,err:=decodeArchiveRows(raw)
	if err!=nil{return err}
	for _,row:=range rows{
		tid,ok:=jsonInt64(row["task_id"])
		if !ok||tid!=snapshot.TaskID{return errors.New("archived report belongs to another task")}
	}
	return insertArchiveRows(tx,"task_report_snapshots",raw)
}
