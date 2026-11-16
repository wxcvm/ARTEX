package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

const releaseAPI = "https://api.github.com/repos/wxcvm/ARTEX/releases/latest"
const releasePage = "https://github.com/wxcvm/ARTEX/releases"

type releaseUpdateAsset struct {
	Name string `json:"name"`
	Size int64 `json:"size"`
	URL string `json:"url"`
}
type releaseUpdate struct {
	Current string `json:"current"`
	Latest string `json:"latest"`
	Comparable bool `json:"comparable"`
	HasUpdate bool `json:"has_update"`
	ReleaseURL string `json:"release_url"`
	PublishedAt string `json:"published_at"`
	Assets []releaseUpdateAsset `json:"assets"`
}
type githubRelease struct {
	Tag string `json:"tag_name"`
	PublishedAt string `json:"published_at"`
	Draft bool `json:"draft"`
	Prerelease bool `json:"prerelease"`
	Assets []struct{
		Name string `json:"name"`
		Size int64 `json:"size"`
	} `json:"assets"`
}

func releaseVersion(raw string)string{
	raw=strings.TrimSpace(raw)
	if !strings.HasPrefix(raw,"v"){raw="v"+raw}
	if !semver.IsValid(raw){return ""}
	// Go semver also accepts shortened v1 / v1.2; ARTEX publishes full tags.
	core:=strings.SplitN(strings.SplitN(raw,"-",2)[0],"+",2)[0]
	if strings.Count(core,".")!=2{return ""}
	return raw
}
func describeRelease(current,goos,goarch string,rel githubRelease)(releaseUpdate,error){
	out:=releaseUpdate{Current:current,Latest:rel.Tag,Assets:[]releaseUpdateAsset{},PublishedAt:rel.PublishedAt}
	latest:=releaseVersion(rel.Tag)
	if latest==""||rel.Draft||rel.Prerelease||semver.Prerelease(latest)!=""{return out,errors.New("invalid stable release")}
	out.ReleaseURL=releasePage+"/tag/"+url.PathEscape(rel.Tag)
	if cur:=releaseVersion(current);cur!=""{
		out.Comparable=true
		out.HasUpdate=semver.Compare(cur,latest)<0
	}
	version:=strings.TrimPrefix(latest,"v")
	names:=[]string{"artex-"+version+"-"+goos+"-"+goarch+".zip"}
	if goos=="windows"&&goarch=="amd64"{
		names=append([]string{"ARTEX-Windows-Setup-amd64.exe","ARTEX-Windows-Desktop-amd64.zip"},names...)
	}
	names=append(names,"SHA256SUMS")
	for _,name:=range names{
		for _,asset:=range rel.Assets{
			if asset.Name==name&&asset.Size>0{
				out.Assets=append(out.Assets,releaseUpdateAsset{Name:name,Size:asset.Size,URL:releasePage+"/download/"+url.PathEscape(rel.Tag)+"/"+url.PathEscape(name)})
				break
			}
		}
	}
	return out,nil
}

func fetchStableRelease(ctx context.Context,client *http.Client,endpoint string)(githubRelease,error){
	var rel githubRelease
	req,err:=http.NewRequestWithContext(ctx,http.MethodGet,endpoint,nil)
	if err!=nil{return rel,err}
	req.Header.Set("Accept","application/vnd.github+json")
	req.Header.Set("User-Agent","ARTEX-release-check")
	resp,err:=client.Do(req)
	if err!=nil{return rel,err}
	defer resp.Body.Close()
	if resp.StatusCode!=http.StatusOK{return rel,fmt.Errorf("release service returned HTTP %d",resp.StatusCode)}
	if err=json.NewDecoder(io.LimitReader(resp.Body,4<<20)).Decode(&rel);err!=nil{return rel,err}
	return rel,nil
}

// Read-only, explicit update check: no installation, process exit or writes.
func(s *Server)checkReleaseUpdate(w http.ResponseWriter,r *http.Request){
	transport:=http.DefaultTransport.(*http.Transport).Clone()
	defer transport.CloseIdleConnections()
	if s.m!=nil {
		if raw:=s.m.GlobalProxy();raw!=""{
			proxy,err:=url.Parse(raw)
			if err!=nil{writeErr(w,400,"invalid configured proxy");return}
			transport.Proxy=http.ProxyURL(proxy)
		}
	}
	client:=&http.Client{Transport:transport,Timeout:15*time.Second}
	rel,err:=fetchStableRelease(r.Context(),client,releaseAPI)
	if err!=nil{writeErr(w,502,"Could not check the official release: "+err.Error());return}
	result,err:=describeRelease(BuildVersion,runtime.GOOS,runtime.GOARCH,rel)
	if err!=nil{writeErr(w,502,err.Error());return}
	w.Header().Set("Cache-Control","no-store")
	writeJSON(w,200,result)
}
