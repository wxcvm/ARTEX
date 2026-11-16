package server

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/report"
)
func(s *Server) registerReportSnapshots(mux *http.ServeMux){
	const base="/api/report/snapshots"
	mux.HandleFunc("GET /api/report/library",s.listReportLibrary)
	mux.HandleFunc("GET "+base,s.listReportSnapshots)
	mux.HandleFunc("POST "+base,s.getReport)
	mux.HandleFunc("GET "+base+"/{snapshot}",s.downloadReportSnapshot)
	mux.HandleFunc("DELETE "+base+"/{snapshot}",s.deleteReportSnapshot)
}
func(s *Server) reportSnapshotTask(w http.ResponseWriter,r *http.Request)(int64,bool){
	raw:=r.URL.Query().Get("task")
	id,err:=strconv.ParseInt(raw,10,64)
	if err!=nil||id<=0{writeErr(w,400,"explicit task is required");return 0,false}
	if s.m.ResolveTask(raw)==nil{writeErr(w,404,"task not found");return 0,false}
	return id,true
}
func(s *Server) saveGeneratedReport(w http.ResponseWriter,r *http.Request,taskID int64,payload []byte,t *db.ReportTemplate){
	name:=strings.TrimSpace(r.URL.Query().Get("name"))
	if name==""{name="Task report"}
	scope:=r.URL.Query().Get("scope")
	if scope==""{scope="all"}
	item:=db.TaskReportSnapshot{Name:name,Format:r.URL.Query().Get("format"),Scope:scope,Locale:report.ReportLocale(r.URL.Query().Get("lang"))}
	if t!=nil{item.TemplateName=t.Name}
	saved,err:=s.m.pg.SaveTaskReportSnapshot(r.Context(),taskID,item,payload)
	if err!=nil{knowledgeDocumentError(w,err);return}
	writeJSON(w,201,saved)
}
func(s *Server) listReportSnapshots(w http.ResponseWriter,r *http.Request){
	id,ok:=s.reportSnapshotTask(w,r);if !ok{return}
	offset:=atoiDefault(r.URL.Query().Get("offset"),0)
	if offset<0{offset=0}
	items,err:=s.m.pg.ListTaskReportSnapshots(r.Context(),id,offset)
	if err!=nil{writeErr(w,500,"could not list saved reports");return}
	writeJSON(w,200,items)
}
func(s *Server) downloadReportSnapshot(w http.ResponseWriter,r *http.Request){
	id,ok:=s.reportSnapshotTask(w,r);if !ok{return}
	format,data,err:=s.m.pg.ReadTaskReportSnapshot(r.Context(),id,r.PathValue("snapshot"))
	if errors.Is(err,db.ErrReportSnapshotNotFound){writeErr(w,404,"report snapshot not found");return}
	if err!=nil{writeErr(w,500,"could not read saved report");return}
	w.Header().Set("Cache-Control","no-store")
	if format=="docx"{w.Header().Set("Content-Type",report.DOCXContentType);w.Header().Set("Content-Disposition",`attachment; filename="task-report.docx"`)}else{
		w.Header().Set("Content-Type","text/markdown; charset=utf-8");w.Header().Set("Content-Disposition",`attachment; filename="task-report.md"`)
	}
	_,_=w.Write(data)
}
func(s *Server) deleteReportSnapshot(w http.ResponseWriter,r *http.Request){
	id,ok:=s.reportSnapshotTask(w,r);if !ok{return}
	err:=s.m.pg.DeleteTaskReportSnapshot(r.Context(),id,r.PathValue("snapshot"))
	if errors.Is(err,db.ErrReportSnapshotNotFound){writeErr(w,404,"report snapshot not found");return}
	if err!=nil{knowledgeDocumentError(w,err);return}
	writeJSON(w,200,map[string]bool{"ok":true})
}

func(s *Server) listReportLibrary(w http.ResponseWriter,r *http.Request){
	pg:=s.pg(w);if pg==nil{return}
	q:=r.URL.Query()
	format:=q.Get("format")
	if format!=""&&format!="md"&&format!="docx"{writeErr(w,400,"invalid report format");return}
	offset:=atoiDefault(q.Get("offset"),0);if offset<0{offset=0}
	items,err:=pg.ListTaskReportLibrary(r.Context(),strings.TrimSpace(q.Get("q")),format,offset)
	if err!=nil{writeErr(w,500,"could not list report library");return}
	writeJSON(w,200,items)
}
