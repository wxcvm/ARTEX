package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"
)

type accountAuditResponse struct {
	http.ResponseWriter
	status int
}

func (w *accountAuditResponse) WriteHeader(status int) {
	if w.status != 0 { return }
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *accountAuditResponse) Write(p []byte) (int,error) {
	if w.status == 0 { w.WriteHeader(http.StatusOK) }
	return w.ResponseWriter.Write(p)
}

// Only wrap the three JSON account handlers, not streaming or general API routes.
// The audit is best-effort and does not change authentication or its response.
func (s *Server) auditAccountAction(action string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter,r *http.Request) {
		aw := &accountAuditResponse{ResponseWriter:w}
		next(aw,r)
		result := "error"
		switch {
		case aw.status >= 200 && aw.status < 300: result = "success"
		case aw.status >= 400 && aw.status < 500: result = "failure"
		}
		if s.m == nil || s.m.pg == nil { return }
		// A disconnected client should not erase the outcome of an account operation.
		ctx,cancel := context.WithTimeout(context.Background(),2*time.Second)
		defer cancel()
		if err := s.m.pg.RecordAccountAudit(ctx,action,result); err != nil {
			log.Print("[account-audit] could not persist account event")
		}
	}
}

// Protected by the existing general JWT middleware, not an /api/auth/ exemption.
func (s *Server) getAccountAudit(w http.ResponseWriter,r *http.Request) {
	pg := s.pg(w)
	if pg == nil { return }
	q := r.URL.Query()
	action,result := q.Get("action"),q.Get("result")
	if action != "" && action != "login" && action != "initialize_password" && action != "change_password" {
		writeErr(w,400,"Invalid account audit action"); return
	}
	if result != "" && result != "success" && result != "failure" && result != "error" {
		writeErr(w,400,"Invalid account audit result"); return
	}
	before,_ := strconv.ParseInt(q.Get("before"),10,64)
	if before < 0 { writeErr(w,400,"Invalid account audit cursor"); return }
	limit,_ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 || limit > 200 { limit = 50 }
	rows,err := pg.ListAccountAudit(r.Context(),action,result,before,limit+1)
	if err != nil { writeErr(w,500,"Could not load account audit"); return }
	more := len(rows)>limit
	if more { rows=rows[:limit] }
	w.Header().Set("Cache-Control","no-store")
	writeJSON(w,200,map[string]any{"items":rows,"has_more":more})
}
