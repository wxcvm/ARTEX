package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrReportTemplateNotFound = errors.New("report template not found")

// ReportTemplate contains text-only presentation settings, never report data.
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
}
const reportTemplateCols = `id,name,brand,cover_title,subtitle,prepared_for,author,confidentiality,footer,created_at,updated_at`

func scanReportTemplate(row interface{Scan(...any)error})(ReportTemplate,error){
	var t ReportTemplate
	err:=row.Scan(&t.ID,&t.Name,&t.Brand,&t.CoverTitle,&t.Subtitle,&t.PreparedFor,&t.Author,&t.Confidentiality,&t.Footer,&t.CreatedAt,&t.UpdatedAt)
	if errors.Is(err,sql.ErrNoRows){err=ErrReportTemplateNotFound}
	return t,err
}
func(d *DB) ListReportTemplates(ctx context.Context)([]ReportTemplate,error){
	rows,err:=d.QueryContext(ctx,"SELECT "+reportTemplateCols+" FROM report_templates ORDER BY name,id")
	if err!=nil{return nil,err}
	defer rows.Close()
	out:=[]ReportTemplate{}
	for rows.Next(){t,err:=scanReportTemplate(rows);if err!=nil{return nil,err};out=append(out,t)}
	return out,rows.Err()
}
func(d *DB) GetReportTemplate(ctx context.Context,id string)(ReportTemplate,error){
	return scanReportTemplate(d.QueryRowContext(ctx,"SELECT "+reportTemplateCols+" FROM report_templates WHERE id=$1",id))
}
func(d *DB) SaveReportTemplate(ctx context.Context,id string,t ReportTemplate)(ReportTemplate,error){
	if id==""{
		return scanReportTemplate(d.QueryRowContext(ctx,`INSERT INTO report_templates(id,name,brand,cover_title,subtitle,prepared_for,author,confidentiality,footer)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING `+reportTemplateCols,uuid.NewString(),t.Name,t.Brand,t.CoverTitle,t.Subtitle,t.PreparedFor,t.Author,t.Confidentiality,t.Footer))
	}
	return scanReportTemplate(d.QueryRowContext(ctx,`UPDATE report_templates SET name=$2,brand=$3,cover_title=$4,subtitle=$5,prepared_for=$6,author=$7,confidentiality=$8,footer=$9,updated_at=now()
WHERE id=$1 RETURNING `+reportTemplateCols,id,t.Name,t.Brand,t.CoverTitle,t.Subtitle,t.PreparedFor,t.Author,t.Confidentiality,t.Footer))
}
func(d *DB) DeleteReportTemplate(ctx context.Context,id string)error{
	res,err:=d.ExecContext(ctx,"DELETE FROM report_templates WHERE id=$1",id)
	if err!=nil{return err}
	n,err:=res.RowsAffected()
	if err==nil&&n==0{return ErrReportTemplateNotFound}
	return err
}
