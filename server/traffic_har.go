package server

import (
	"errors"
	"io"
	"net/http"
)

// Import a saved recording through the authenticated traffic API only.
func (s *Server) importTrafficHAR(w http.ResponseWriter,r *http.Request) {
	tr:=s.m.Traffic()
	if tr==nil { writeErr(w,http.StatusServiceUnavailable,"traffic store is not available"); return }
	raw,err:=io.ReadAll(http.MaxBytesReader(w,r.Body,32<<20))
	if err!=nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err,&tooLarge) { writeErr(w,http.StatusRequestEntityTooLarge,"HAR exceeds 32 MiB") } else { writeErr(w,http.StatusBadRequest,"unable to read HAR") }
		return
	}
	result,err:=tr.ImportHAR(r.Context(),raw)
	if err!=nil { writeErr(w,http.StatusBadRequest,err.Error()); return }
	writeJSON(w,http.StatusOK,result)
}
