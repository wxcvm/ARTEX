package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/knowledgeindex"
)

func (s *Server) getKnowledgeEmbedding(w http.ResponseWriter,r *http.Request){
	c,err:=knowledgeindex.Load(s.m.pg)
	if err!=nil{writeErr(w,500,"failed to read embedding settings");return}
	writeJSON(w,200,c)
}

func (s *Server) putKnowledgeEmbedding(w http.ResponseWriter,r *http.Request){
	var c knowledgeindex.Config
	if err:=json.NewDecoder(http.MaxBytesReader(w,r.Body,64<<10)).Decode(&c);err!=nil{writeErr(w,400,"invalid embedding settings");return}
	c,err:=c.Normalize()
	if err!=nil{writeErr(w,400,err.Error());return}
	c,err=knowledgeindex.Save(s.m.pg,c)
	if err!=nil{writeErr(w,500,"failed to save embedding settings");return}
	writeJSON(w,200,c)
}

func (s *Server) testKnowledgeEmbedding(w http.ResponseWriter,r *http.Request){
	var c knowledgeindex.Config
	if err:=json.NewDecoder(http.MaxBytesReader(w,r.Body,64<<10)).Decode(&c);err!=nil{writeErr(w,400,"invalid embedding settings");return}
	c.Enabled=true
	v,err:=knowledgeindex.Embed(r.Context(),c,[]string{"ARTEX local embedding connection test"})
	if err!=nil{writeJSON(w,200,map[string]any{"ok":false,"error":err.Error()});return}
	writeJSON(w,200,map[string]any{"ok":true,"dimension":len(v[0])})
}

func (s *Server) getKnowledgeIndex(w http.ResponseWriter,r *http.Request){
	taskID,ok:=s.knowledgeDocumentTask(w,r);if !ok{return}
	docID:=r.PathValue("document")
	if _,err:=s.m.pg.GetKnowledgeDocument(r.Context(),taskID,false,docID);err!=nil{knowledgeDocumentError(w,err);return}
	c,err:=knowledgeindex.Load(s.m.pg)
	if err!=nil{writeErr(w,500,"failed to read embedding settings");return}
	count:=0
	if c.Enabled{
		count,err=s.m.pg.KnowledgeVectorCount(r.Context(),taskID,docID,c.Key())
		if err!=nil{knowledgeDocumentError(w,err);return}
	}
	writeJSON(w,200,map[string]any{"enabled":c.Enabled,"chunks":count,"model":c.Model})
}

func (s *Server) buildKnowledgeIndex(w http.ResponseWriter,r *http.Request){
	taskID,ok:=s.knowledgeDocumentTask(w,r);if !ok{return}
	c,err:=knowledgeindex.Load(s.m.pg)
	if err!=nil{writeErr(w,500,"failed to read embedding settings");return}
	if !c.Enabled{writeErr(w,400,"enable local embedding first");return}
	var result knowledgeindex.BuildResult
	if r.URL.Query().Get("incremental")=="true"{
		result,err=knowledgeindex.BuildIncremental(r.Context(),s.m.pg,taskID,r.PathValue("document"),c)
	}else{
		result.Chunks,err=knowledgeindex.Build(r.Context(),s.m.pg,taskID,r.PathValue("document"),c)
		result.Embedded=result.Chunks
	}
	if err!=nil{
		switch {
		case errors.Is(err,db.ErrKnowledgeIndexStale):writeErr(w,409,err.Error())
		case errors.Is(err,db.ErrKnowledgeDocumentNotFound),errors.Is(err,db.ErrTaskArchiveState):knowledgeDocumentError(w,err)
		default:writeErr(w,502,"local index build failed: "+err.Error())
		}
		return
	}
	writeJSON(w,200,map[string]any{"enabled":true,"chunks":result.Chunks,"model":c.Model,"reused":result.Reused,"embedded":result.Embedded,"skipped":result.Skipped})
}

func (s *Server) listPendingKnowledgeIndexes(w http.ResponseWriter,r *http.Request){
	taskID,ok:=s.knowledgeDocumentTask(w,r);if !ok{return}
	c,err:=knowledgeindex.Load(s.m.pg)
	if err!=nil{writeErr(w,500,"failed to read embedding settings");return}
	if !c.Enabled{writeErr(w,400,"enable local embedding first");return}
	items,err:=s.m.pg.ListPendingKnowledgeIndexes(r.Context(),taskID,c.Key(),r.URL.Query().Get("after"))
	if err!=nil{knowledgeDocumentError(w,err);return}
	next:=""
	if len(items)==50{next=items[len(items)-1].ID}
	writeJSON(w,200,map[string]any{"items":items,"next":next,"model":c.Model})
}
