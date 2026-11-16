package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrReportTemplateNotFound = errors.New("report template not found")

// ReportSection is literal user-authored report context, separate from findings.
type ReportSection struct {
	Title string `json:"title"`
	Body string `json:"body"`
	Position string `json:"position"`
}

// ReportTemplate contains reusable presentation settings, not stored findings.
type ReportTemplate struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Brand string `json:"brand"`
	CoverTitle string `json:"cover_title"`
	Subtitle string `json:"subtitle"`
	PreparedFor string `json:"prepared_for"`
	Author string `json:"author"`
	Confidentiality string `json:"confidentiality"`
	Footer string `json:"footer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	HasLogo bool `json:"has_logo"`
	Logo []byte `json:"-"`
	Sections []ReportSection `json:"sections"`
}
const reportTemplateCols = `id,name,brand,cover_title,subtitle,prepared_for,author,confidentiality,footer,created_at,updated_at,(logo IS NOT NULL),sections`

func scanReportTemplate(row interface{Scan(...any)error}, withLogo bool)(ReportTemplate,error){
	var t ReportTemplate
	var sections []byte
	dest:=[]any{&t.ID,&t.Name,&t.Brand,&t.CoverTitle,&t.Subtitle,&t.PreparedFor,&t.Author,&t.Confidentiality,&t.Footer,&t.CreatedAt,&t.UpdatedAt,&t.HasLogo,&sections}
	if withLogo{dest=append(dest,&t.Logo)}
	err:=row.Scan(dest...)
	if errors.Is(err,sql.ErrNoRows){err=ErrReportTemplateNotFound}
	if err==nil{err=json.Unmarshal(sections,&t.Sections)}
	return t,err
}
func(d *DB) ListReportTemplates(ctx context.Context)([]ReportTemplate,error){
	rows,err:=d.QueryContext(ctx,"SELECT "+reportTemplateCols+" FROM report_templates ORDER BY name,id")
	if err!=nil{return nil,err}
	defer rows.Close()
	out:=[]ReportTemplate{}
	for rows.Next(){t,err:=scanReportTemplate(rows,false);if err!=nil{return nil,err};out=append(out,t)}
	return out,rows.Err()
}
func(d *DB) GetReportTemplate(ctx context.Context,id string)(ReportTemplate,error){
	return scanReportTemplate(d.QueryRowContext(ctx,"SELECT "+reportTemplateCols+",logo FROM report_templates WHERE id=$1",id),true)
}
func(d *DB) SaveReportTemplate(ctx context.Context,id string,t ReportTemplate)(ReportTemplate,error){
	var sections any
	if t.Sections!=nil{raw,err:=json.Marshal(t.Sections);if err!=nil{return ReportTemplate{},err};sections=string(raw)}
	if id==""{
		return scanReportTemplate(d.QueryRowContext(ctx,`INSERT INTO report_templates(id,name,brand,cover_title,subtitle,prepared_for,author,confidentiality,footer,sections)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,COALESCE($10::jsonb,'[]'::jsonb)) RETURNING `+reportTemplateCols,uuid.NewString(),t.Name,t.Brand,t.CoverTitle,t.Subtitle,t.PreparedFor,t.Author,t.Confidentiality,t.Footer,sections),false)
	}
	return scanReportTemplate(d.QueryRowContext(ctx,`UPDATE report_templates SET name=$2,brand=$3,cover_title=$4,subtitle=$5,prepared_for=$6,author=$7,confidentiality=$8,footer=$9,sections=COALESCE($10::jsonb,sections),updated_at=now()
WHERE id=$1 RETURNING `+reportTemplateCols,id,t.Name,t.Brand,t.CoverTitle,t.Subtitle,t.PreparedFor,t.Author,t.Confidentiality,t.Footer,sections),false)
}
func(d *DB) DeleteReportTemplate(ctx context.Context,id string)error{
	res,err:=d.ExecContext(ctx,"DELETE FROM report_templates WHERE id=$1",id)
	if err!=nil{return err}
	n,err:=res.RowsAffected()
	if err==nil&&n==0{return ErrReportTemplateNotFound}
	return err
}

func(d *DB) SetReportTemplateLogo(ctx context.Context,id string,logo []byte)(ReportTemplate,error){
	return scanReportTemplate(d.QueryRowContext(ctx,`UPDATE report_templates SET logo=$2,updated_at=now() WHERE id=$1 RETURNING `+reportTemplateCols,id,logo),false)
}
