package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestAccountAuditRealLoginAndReadAccess(t *testing.T) {
	m,err:=NewManager(t.TempDir(),"")
	if err!=nil { t.Skipf("postgres unavailable: %v",err) }
	defer m.Close()
	s:=&Server{m:m,jwtKey:[]byte(strings.Repeat("K",32))}
	old,exists,err:=m.pg.GetSetting(authPassKey)
	if err!=nil {t.Fatal(err)}
	defer func(){
		if exists { _=m.pg.SetSetting(authPassKey,old) } else { _,_=m.pg.Exec("DELETE FROM settings WHERE key=$1",authPassKey) }
	}()
	hash,err:=bcrypt.GenerateFromPassword([]byte("synthetic-valid-password"),bcrypt.MinCost)
	if err!=nil { t.Fatal(err) }
	if err=m.pg.SetSetting(authPassKey,string(hash)); err!=nil {t.Fatal(err)}
	var baseline int64
	if err=m.pg.QueryRow("SELECT COALESCE(MAX(id),0) FROM account_audit").Scan(&baseline);err!=nil {t.Fatal(err)}
	defer m.pg.Exec("DELETE FROM account_audit WHERE id>$1",baseline)
	mux:=http.NewServeMux()
	mux.HandleFunc("POST /api/auth/login",s.auditAccountAction("login",s.authLogin))
	mux.HandleFunc("GET /api/logs/account-audit",s.getAccountAudit)
	h:=s.requireAuth(mux)
	call:=func(method,path,body,token string)*httptest.ResponseRecorder{
		r:=httptest.NewRequest(method,path,strings.NewReader(body))
		if token!="" {r.Header.Set("Authorization","Bearer "+token)}
		w:=httptest.NewRecorder();h.ServeHTTP(w,r);return w
	}
	failed:=call("POST","/api/auth/login",`{"username":"synthetic-untrusted-name","password":"synthetic-secret"}`,"")
	if failed.Code!=401 {t.Fatalf("wrong login: %d",failed.Code)}
	good:=call("POST","/api/auth/login",`{"username":"ARTEX","password":"synthetic-valid-password"}`,"")
	if good.Code!=200 {t.Fatalf("valid login: %d %s",good.Code,good.Body.String())}
	var token struct{Token string `json:"token"`}
	if err=json.Unmarshal(good.Body.Bytes(),&token);err!=nil||token.Token=="" {t.Fatal("missing token")}
	if w:=call("GET","/api/logs/account-audit","","");w.Code!=401 {t.Fatalf("anonymous audit read: %d",w.Code)}
	w:=call("GET","/api/logs/account-audit?action=login&limit=1","",token.Token)
	if w.Code!=200 || !strings.Contains(w.Body.String(),`"result":"success"`) || !strings.Contains(w.Body.String(),`"has_more":true`) {t.Fatalf("audit: %d %s",w.Code,w.Body.String())}
	if strings.Contains(w.Body.String(),"synthetic")||strings.Contains(w.Body.String(),token.Token){t.Fatal("secret exposed")}
	rows,err:=m.pg.ListAccountAudit(context.Background(),"login","",0,2)
	if err!=nil||len(rows)!=2||rows[0].Result!="success"||rows[1].Result!="failure" {t.Fatalf("rows: %+v %v",rows,err)}
	older,err:=m.pg.ListAccountAudit(context.Background(),"login","failure",rows[0].ID,1)
	if err!=nil||len(older)!=1||older[0].ID!=rows[1].ID {t.Fatalf("cursor: %+v %v",older,err)}
}
