package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/kb"
)

func (s *Server) registerKnowledgeDocuments(mux *http.ServeMux) {
	const base="/api/knowledge/documents"
	mux.HandleFunc("GET "+base,s.listKnowledgeDocuments)
	mux.HandleFunc("GET /api/knowledge/retrievals",s.listKnowledgeRetrievals)
	mux.HandleFunc("POST "+base,s.saveKnowledgeDocument)
	mux.HandleFunc("POST "+base+"/upload",s.uploadKnowledgeDocument)
	mux.HandleFunc("GET "+base+"/{document}",s.getKnowledgeDocument)
	mux.HandleFunc("PUT "+base+"/{document}",s.saveKnowledgeDocument)
	mux.HandleFunc("DELETE "+base+"/{document}",s.deleteKnowledgeDocument)
}

// Omitted task is explicitly global. Invalid, missing or archived task IDs
// never fall through to the global library.
func (s *Server) knowledgeDocumentTask(w http.ResponseWriter,r *http.Request)(int64,bool){
	raw:=r.URL.Query().Get("task")
	if raw==""{return 0,true}
	id,err:=strconv.ParseInt(raw,10,64)
	if err!=nil||id<=0{writeErr(w,400,"invalid task id");return 0,false}
	task,err:=s.m.pg.GetTask(id)
	if err!=nil{writeErr(w,500,"failed to read task");return 0,false}
	if task==nil{writeErr(w,404,"task not found");return 0,false}
	return id,true
}

func knowledgeDocumentError(w http.ResponseWriter,err error){
	switch {
	case errors.Is(err,db.ErrKnowledgeDocumentNotFound),errors.Is(err,sql.ErrNoRows): writeErr(w,404,"document not found")
	case errors.Is(err,db.ErrTaskArchiveState): writeErr(w,409,"task is being archived or is archived")
	default: writeErr(w,500,"knowledge document operation failed")
	}
}

func (s *Server) listKnowledgeDocuments(w http.ResponseWriter,r *http.Request){
	id,ok:=s.knowledgeDocumentTask(w,r);if !ok{return}
	q:=r.URL.Query()
	offset,err:=strconv.Atoi(q.Get("offset"))
	if q.Get("offset")==""{offset,err=0,nil}
	if err!=nil||offset<0{writeErr(w,400,"invalid offset");return}
	page,err:=s.m.pg.ListKnowledgeDocuments(r.Context(),id,q.Get("include_global")=="true",strings.TrimSpace(q.Get("q")),q.Get("category"),offset)
	if err!=nil{knowledgeDocumentError(w,err);return}
	writeJSON(w,200,page)
}
func (s *Server) getKnowledgeDocument(w http.ResponseWriter,r *http.Request){
	id,ok:=s.knowledgeDocumentTask(w,r);if !ok{return}
	doc,err:=s.m.pg.GetKnowledgeDocument(r.Context(),id,r.URL.Query().Get("include_global")=="true",r.PathValue("document"))
	if err!=nil{knowledgeDocumentError(w,err);return}
	writeJSON(w,200,doc)
}
func validKnowledgeDocument(doc *db.KnowledgeDocument)bool{
	doc.Title=strings.TrimSpace(doc.Title)
	doc.Category=strings.TrimSpace(doc.Category)
	return doc.Title!="" && len(doc.Title)<=4096 && len(doc.Category)<=4096 && strings.TrimSpace(doc.Content)!="" && len(doc.Content)<=kb.DocumentTextLimit && !strings.ContainsRune(doc.Title+doc.Category+doc.Content,0)
}
func (s *Server) saveKnowledgeDocument(w http.ResponseWriter,r *http.Request){
	id,ok:=s.knowledgeDocumentTask(w,r);if !ok{return}
	var doc db.KnowledgeDocument
	if err:=json.NewDecoder(http.MaxBytesReader(w,r.Body,8<<20)).Decode(&doc);err!=nil||!validKnowledgeDocument(&doc){writeErr(w,400,"title and nonempty text required; text limit 4 MiB");return}
	doc.SourceFilename="" // only the upload path sets provenance
	out,err:=s.m.pg.SaveKnowledgeDocument(r.Context(),id,r.PathValue("document"),doc)
	if err!=nil{knowledgeDocumentError(w,err);return}
	writeJSON(w,200,out)
}
func (s *Server) deleteKnowledgeDocument(w http.ResponseWriter,r *http.Request){
	id,ok:=s.knowledgeDocumentTask(w,r);if !ok{return}
	if err:=s.m.pg.DeleteKnowledgeDocument(r.Context(),id,r.PathValue("document"));err!=nil{knowledgeDocumentError(w,err);return}
	writeJSON(w,200,map[string]bool{"ok":true})
}
func (s *Server) uploadKnowledgeDocument(w http.ResponseWriter,r *http.Request){
	id,ok:=s.knowledgeDocumentTask(w,r);if !ok{return}
	r.Body=http.MaxBytesReader(w,r.Body,9<<20)
	if err:=r.ParseMultipartForm(8<<20);err!=nil{writeErr(w,400,"invalid or oversized upload");return}
	defer r.MultipartForm.RemoveAll()
	f,h,err:=r.FormFile("file");if err!=nil{writeErr(w,400,"file required");return}
	defer f.Close()
	raw,err:=io.ReadAll(io.LimitReader(f,(8<<20)+1))
	if err!=nil||len(raw)>8<<20{writeErr(w,400,"upload exceeds 8 MiB");return}
	name:=filepath.Base(strings.ReplaceAll(h.Filename,"\\","/"))
	content,err:=kb.ExtractDocumentText(name,raw)
	if err!=nil{writeErr(w,400,err.Error());return}
	doc:=db.KnowledgeDocument{Title:strings.TrimSuffix(name,filepath.Ext(name)),Category:r.FormValue("category"),Content:content,SourceFilename:name}
	if !validKnowledgeDocument(&doc){writeErr(w,400,"invalid document metadata");return}
	out,err:=s.m.pg.SaveKnowledgeDocument(r.Context(),id,"",doc)
	if err!=nil{knowledgeDocumentError(w,err);return}
	writeJSON(w,201,out)
}

func (s *Server) listKnowledgeRetrievals(w http.ResponseWriter,r *http.Request){
	id,ok:=s.knowledgeDocumentTask(w,r)
	if !ok{return}
	items,err:=s.m.pg.ListKnowledgeRetrievals(r.Context(),id)
	if err!=nil{knowledgeDocumentError(w,err);return}
	writeJSON(w,200,items)
}
