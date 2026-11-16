package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReleaseUpdateSelection(t *testing.T){
	var rel githubRelease
	rel.Tag="v0.5.2"
	rel.PublishedAt="2026-11-16T00:00:00Z"
	for _,name:=range []string{"artex-0.5.2-windows-amd64.zip","artex-0.5.2-linux-amd64.zip","ARTEX-Windows-Setup-amd64.exe","ARTEX-Windows-Desktop-amd64.zip","SHA256SUMS","unrelated.exe"}{
		rel.Assets=append(rel.Assets,struct{Name string `json:"name"`;Size int64 `json:"size"`}{name,100})
	}
	for _,tc:=range []struct{current string;comparable,update bool}{
		{"v0.5.1",true,true},{"0.5.2",true,false},{"v0.5.3",true,false},{"v0.5.2-beta.1",true,true},{"dev",false,false},{"v0.5",false,false},
	}{
		got,err:=describeRelease(tc.current,"windows","amd64",rel)
		if err!=nil||got.Comparable!=tc.comparable||got.HasUpdate!=tc.update||len(got.Assets)!=4||got.Assets[0].Name!="ARTEX-Windows-Setup-amd64.exe"{t.Fatal(tc,got,err)}
		for _,asset:=range got.Assets{if !strings.HasPrefix(asset.URL,releasePage+"/download/v0.5.2/"){t.Fatal(asset)}}
	}
	got,err:=describeRelease("v0.5.1","linux","amd64",rel)
	if err!=nil||len(got.Assets)!=2||got.Assets[0].Name!="artex-0.5.2-linux-amd64.zip"{t.Fatal(got,err)}
	rel.Prerelease=true
	if _,err:=describeRelease("v0.5.1","windows","amd64",rel);err==nil{t.Fatal("accepted prerelease")}
}
func TestFetchReleaseProtocol(t *testing.T){
	mock:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		if r.Header.Get("Accept")!="application/vnd.github+json"{t.Error("missing GitHub accept header")}
		_,_=w.Write([]byte(`{"tag_name":"v0.5.2","assets":[{"name":"ARTEX-Windows-Setup-amd64.exe","size":100}]}`))
	}))
	defer mock.Close()
	rel,err:=fetchStableRelease(context.Background(),mock.Client(),mock.URL)
	if err!=nil||rel.Tag!="v0.5.2"||len(rel.Assets)!=1{t.Fatal(rel,err)}
}
