package server

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/traffic"
)

// Compare the immutable snapshots from one finding version. The original hot
// traffic store may already have been cleaned or disabled.
func (s *Server) compareFindingTraffic(w http.ResponseWriter, r *http.Request) {
	id, ok := s.findingTrafficAccess(w, r, false)
	if !ok { return }
	q := r.URL.Query()
	a, errA := strconv.ParseInt(q.Get("baseline"),10,64)
	b, errB := strconv.ParseInt(q.Get("candidate"),10,64)
	if errA != nil || errB != nil || a <= 0 || b <= 0 || a == b {
		writeErr(w,400,"two distinct binding IDs are required")
		return
	}
	list, err := s.m.pg.GetFindingTraffic(r.Context(),id)
	if err != nil { evidenceError(w,err); return }
	if raw := q.Get("version"); raw != "" {
		version, err := strconv.ParseInt(raw,10,64)
		if err != nil { writeErr(w,400,"invalid evidence version"); return }
		if version != list.Version { evidenceError(w,db.ErrEvidenceConflict); return }
	}
	var left,right *db.FindingTrafficBinding
	for i := range list.Bindings {
		item := &list.Bindings[i]
		if item.ID == a { left = item }
		if item.ID == b { right = item }
	}
	if left == nil || right == nil { evidenceError(w,db.ErrEvidenceNotFound); return }
	result, err := s.compareSavedEvidence(r.Context(),*left,*right)
	if err != nil { evidenceError(w,err); return }
	writeJSON(w,200,struct {
		traffic.RecordedComparison
		EvidenceVersion int64 `json:"evidence_version"`
	}{result,list.Version})
}

func (s *Server) compareSavedEvidence(ctx context.Context, a,b db.FindingTrafficBinding) (traffic.RecordedComparison,error) {
	store := s.evidenceStore()
	var files []*os.File
	defer func(){ for _, f := range files { _ = f.Close() } }()
	read := func(binding db.FindingTrafficBinding) (traffic.EvidenceExchange,error) {
		v := binding.Snapshot
		req,_,err := store.OpenBody(v,"request")
		if err != nil { return traffic.EvidenceExchange{},err }
		files = append(files,req)
		resp,_,err := store.OpenBody(v,"response")
		if err != nil { return traffic.EvidenceExchange{},err }
		files = append(files,resp)
		return traffic.EvidenceExchange{ID:strconv.FormatInt(binding.ID,10),URL:v.URL,Method:v.Method,ContentType:v.ContentType,
			TS:v.CapturedAt,Status:v.Status,ReqHead:v.ReqHead,RespHead:v.RespHead,
			ReqLen:v.ReqLen,RespLen:v.RespLen,ReqHash:v.ReqHash,RespHash:v.RespHash,Request:req,Response:resp},nil
	}
	left,err := read(a)
	if err != nil { return traffic.RecordedComparison{},err }
	right,err := read(b)
	if err != nil { return traffic.RecordedComparison{},err }
	return traffic.CompareEvidencePair(ctx,left,right)
}
