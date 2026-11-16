package server

import (
	"net/http"
)

// Compare two saved records only; no requests are issued to recorded URLs.
func (s *Server) compareTraffic(w http.ResponseWriter,r *http.Request) {
	tr:=s.m.Traffic()
	if tr==nil { writeErr(w,http.StatusServiceUnavailable,"traffic store is not available"); return }
	q:=r.URL.Query()
	a,b:=q.Get("baseline"),q.Get("candidate")
	if a=="" || b=="" || a==b { writeErr(w,http.StatusBadRequest,"two distinct traffic IDs are required"); return }
	result,err:=tr.CompareRecorded(r.Context(),a,b)
	if err!=nil { writeErr(w,http.StatusUnprocessableEntity,"unable to compare recorded traffic"); return }
	writeJSON(w,http.StatusOK,result)
}
