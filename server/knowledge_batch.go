package server

import (
	"context"
	"net/http"

	"github.com/Autumn-27/artex/knowledgeindex"
)

type knowledgeBatchStatus struct {
	State string `json:"state"`
	Model string `json:"model,omitempty"`
	Current string `json:"current,omitempty"`
	Succeeded int `json:"succeeded"`
	Failed int `json:"failed"`
	Reused int `json:"reused"`
	Embedded int `json:"embedded"`
	Errors []string `json:"errors,omitempty"`
}

type knowledgeBatchJob struct {
	status knowledgeBatchStatus
	cancel context.CancelFunc
}

func (s *Server) knowledgeBatchSnapshot(taskID int64) knowledgeBatchStatus {
	s.knowledgeBatchMu.Lock()
	defer s.knowledgeBatchMu.Unlock()
	if job:=s.knowledgeBatches[taskID];job!=nil {
		out:=job.status
		out.Errors=append([]string(nil),out.Errors...)
		return out
	}
	return knowledgeBatchStatus{State:"idle"}
}

func (s *Server) getKnowledgeBatch(w http.ResponseWriter,r *http.Request){
	id,ok:=s.knowledgeDocumentTask(w,r);if !ok{return}
	writeJSON(w,200,s.knowledgeBatchSnapshot(id))
}

func (s *Server) startKnowledgeBatch(w http.ResponseWriter,r *http.Request){
	id,ok:=s.knowledgeDocumentTask(w,r);if !ok{return}
	c,err:=knowledgeindex.Load(s.m.pg)
	if err!=nil{writeErr(w,500,"failed to read embedding settings");return}
	if !c.Enabled{writeErr(w,400,"enable local embedding first");return}
	s.knowledgeBatchMu.Lock()
	if existing:=s.knowledgeBatches[id];existing!=nil&&(existing.status.State=="running"||existing.status.State=="stopping"){
		s.knowledgeBatchMu.Unlock()
		writeJSON(w,200,s.knowledgeBatchSnapshot(id));return
	}
	parent:=s.ctx
	if parent==nil{parent=context.Background()}
	ctx,cancel:=context.WithCancel(parent)
	job:=&knowledgeBatchJob{status:knowledgeBatchStatus{State:"running",Model:c.Model},cancel:cancel}
	if s.knowledgeBatches==nil{s.knowledgeBatches=make(map[int64]*knowledgeBatchJob)}
	s.knowledgeBatches[id]=job
	s.knowledgeBatchMu.Unlock()
	go s.runKnowledgeBatch(ctx,id,c,job)
	writeJSON(w,202,s.knowledgeBatchSnapshot(id))
}

func (s *Server) stopKnowledgeBatch(w http.ResponseWriter,r *http.Request){
	id,ok:=s.knowledgeDocumentTask(w,r);if !ok{return}
	s.knowledgeBatchMu.Lock()
	if job:=s.knowledgeBatches[id];job!=nil&&job.status.State=="running"{
		job.status.State="stopping"
		job.cancel()
	}
	s.knowledgeBatchMu.Unlock()
	writeJSON(w,200,s.knowledgeBatchSnapshot(id))
}

func (s *Server) runKnowledgeBatch(ctx context.Context,id int64,c knowledgeindex.Config,job *knowledgeBatchJob){
	defer job.cancel()
	update:=func(fn func(*knowledgeBatchStatus)){
		s.knowledgeBatchMu.Lock();defer s.knowledgeBatchMu.Unlock()
		fn(&job.status)
	}
	finish:=func(state,msg string){
		update(func(st *knowledgeBatchStatus){
			if ctx.Err()!=nil{state="cancelled";msg=""}
			st.State=state;st.Current=""
			if msg!=""{st.Errors=append(st.Errors,msg)}
		})
	}
	after:=""
	for {
		if ctx.Err()!=nil{finish("cancelled","");return}
		items,err:=s.m.pg.ListPendingKnowledgeIndexes(ctx,id,c.Key(),after)
		if err!=nil{finish("error","could not list pending documents");return}
		for _,doc:=range items{
			if ctx.Err()!=nil{finish("cancelled","");return}
			// A job uses the selected model consistently. Stop if settings change
			// instead of continuing to index the library for an obsolete model.
			current,err:=knowledgeindex.Load(s.m.pg)
			if err!=nil||!current.Enabled||current.Key()!=c.Key(){
				finish("error","embedding settings changed or became unavailable; restart indexing");return
			}
			update(func(st *knowledgeBatchStatus){st.Current=doc.Title})
			result,err:=knowledgeindex.BuildIncremental(ctx,s.m.pg,id,doc.ID,c)
			if err!=nil&&ctx.Err()!=nil{finish("cancelled","");return}
			update(func(st *knowledgeBatchStatus){
				if err!=nil{
					st.Failed++
					if len(st.Errors)>=20{st.Errors=st.Errors[1:]}
					st.Errors=append(st.Errors,doc.Title+": "+err.Error())
				}else{
					st.Succeeded++;st.Reused+=result.Reused;st.Embedded+=result.Embedded
				}
			})
			after=doc.ID
		}
		if len(items)<50{finish("finished","");return}
	}
}
