package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Autumn-27/artex/db"
)

func(s *Server) registerReportTemplates(mux *http.ServeMux){
	const base="/api/report/templates"
	mux.HandleFunc("GET "+base,s.listReportTemplates)
	mux.HandleFunc("POST "+base,s.saveReportTemplate)
	mux.HandleFunc("PUT "+base+"/{template}",s.saveReportTemplate)
	mux.HandleFunc("DELETE "+base+"/{template}",s.deleteReportTemplate)
}
func(s *Server) listReportTemplates(w http.ResponseWriter,r *http.Request){
	pg:=s.pg(w);if pg==nil{return}
	items,err:=pg.ListReportTemplates(r.Context())
	if err!=nil{writeErr(w,500,"could not list report templates");return}
	writeJSON(w,200,items)
}
func validReportTemplate(t *db.ReportTemplate)bool{
	for _,p:=range []*string{&t.Name,&t.Brand,&t.CoverTitle,&t.Subtitle,&t.PreparedFor,&t.Author,&t.Confidentiality,&t.Footer}{
		*p=strings.TrimSpace(*p)
		if len(*p)>4096||strings.ContainsRune(*p,0){return false}
		for _,r:=range *p{if r<32&&r!='\t'&&r!='\n'&&r!='\r'{return false}}
	}
	return t.Name!=""
}
func(s *Server) saveReportTemplate(w http.ResponseWriter,r *http.Request){
	pg:=s.pg(w);if pg==nil{return}
	var t db.ReportTemplate
	if err:=json.NewDecoder(http.MaxBytesReader(w,r.Body,128<<10)).Decode(&t);err!=nil||!validReportTemplate(&t){writeErr(w,400,"invalid report template; name required");return}
	item,err:=pg.SaveReportTemplate(r.Context(),r.PathValue("template"),t)
	if errors.Is(err,db.ErrReportTemplateNotFound){writeErr(w,404,"report template not found");return}
	if err!=nil{writeErr(w,500,"could not save report template");return}
	writeJSON(w,200,item)
}
func(s *Server) deleteReportTemplate(w http.ResponseWriter,r *http.Request){
	pg:=s.pg(w);if pg==nil{return}
	err:=pg.DeleteReportTemplate(r.Context(),r.PathValue("template"))
	if errors.Is(err,db.ErrReportTemplateNotFound){writeErr(w,404,"report template not found");return}
	if err!=nil{writeErr(w,500,"could not delete report template");return}
	writeJSON(w,200,map[string]bool{"ok":true})
}
func(s *Server) resolveReportTemplate(w http.ResponseWriter,r *http.Request)(*db.ReportTemplate,bool){
	id:=r.URL.Query().Get("template")
	if id==""{return nil,true}
	pg:=s.pg(w);if pg==nil{return nil,false}
	t,err:=pg.GetReportTemplate(r.Context(),id)
	if errors.Is(err,db.ErrReportTemplateNotFound){writeErr(w,404,"report template not found");return nil,false}
	if err!=nil{writeErr(w,500,"could not read report template");return nil,false}
	return &t,true
}
