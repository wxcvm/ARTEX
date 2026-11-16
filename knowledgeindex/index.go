// Package knowledgeindex provides optional local document embeddings.
package knowledgeindex

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Autumn-27/artex/db"
)

const SettingKey="knowledge_embedding_config"
const chunkSize=800
const chunkOverlap=100
const batchSize=16

type Config struct {
	Enabled bool `json:"enabled"`
	BaseURL string `json:"base_url"`
	Model string `json:"model"`
}

// Normalize only accepts a local machine endpoint. The library never sends
// document text to an external API or follows redirects.
func (c Config) Normalize()(Config,error){
	c.BaseURL=strings.TrimSpace(c.BaseURL)
	c.Model=strings.TrimSpace(c.Model)
	if !c.Enabled&&c.BaseURL==""&&c.Model==""{return c,nil}
	u,err:=url.Parse(c.BaseURL)
	if err!=nil||u==nil||u.User!=nil||u.RawQuery!=""||u.Fragment!=""||(u.Scheme!="http"&&u.Scheme!="https"){return c,errors.New("use a local HTTP(S) embedding base URL")}
	host:=u.Hostname()
	if strings.EqualFold(host,"localhost"){host="127.0.0.1"}
	ip:=net.ParseIP(host)
	if ip==nil||!ip.IsLoopback(){return c,errors.New("embedding endpoint must use localhost or a loopback IP")}
	if u.Port()!=""{u.Host=net.JoinHostPort(host,u.Port())}else if strings.Contains(host,":"){u.Host="["+host+"]"}else{u.Host=host}
	u.Path=strings.TrimRight(u.Path,"/")
	if u.Path==""{u.Path="/v1"}
	if u.Path!="/v1"{return c,errors.New("use the OpenAI-compatible /v1 base URL")}
	u.RawPath=""
	c.BaseURL=u.String()
	if c.Enabled&&c.Model==""{return c,errors.New("embedding model is required")}
	return c,nil
}

func (c Config) Key()string{return "v1:800:100:"+c.BaseURL+"|"+c.Model}

func Load(pg *db.DB)(Config,error){
	raw,ok,err:=pg.GetSetting(SettingKey)
	if err!=nil{return Config{},err}
	if !ok{return Config{},nil}
	var c Config
	if err=json.Unmarshal([]byte(raw),&c);err!=nil{return c,err}
	return c.Normalize()
}

func Save(pg *db.DB,c Config)(Config,error){
	c,err:=c.Normalize()
	if err!=nil{return c,err}
	raw,err:=json.Marshal(c)
	if err!=nil{return c,err}
	return c,pg.SetSetting(SettingKey,string(raw))
}

var localClient=&http.Client{
	Transport:&http.Transport{Proxy:nil, DialContext:(&net.Dialer{Timeout:5*time.Second}).DialContext},
	Timeout:30*time.Second,
	CheckRedirect:func(_ *http.Request,_ []*http.Request)error{return http.ErrUseLastResponse},
}

func Embed(ctx context.Context,c Config,inputs []string)([][]float32,error){
	c,err:=c.Normalize()
	if err!=nil{return nil,err}
	if !c.Enabled||len(inputs)==0{return nil,errors.New("local embedding is disabled or no input")}
	body,err:=json.Marshal(map[string]any{"model":c.Model,"input":inputs})
	if err!=nil{return nil,err}
	req,err:=http.NewRequestWithContext(ctx,http.MethodPost,c.BaseURL+"/embeddings",bytes.NewReader(body))
	if err!=nil{return nil,err}
	req.Header.Set("Content-Type","application/json")
	resp,err:=localClient.Do(req)
	if err!=nil{return nil,fmt.Errorf("local embedding request failed: %w",err)}
	defer resp.Body.Close()
	if resp.StatusCode!=http.StatusOK{return nil,fmt.Errorf("local embedding returned HTTP %d",resp.StatusCode)}
	raw,err:=io.ReadAll(io.LimitReader(resp.Body,(32<<20)+1))
	if err!=nil{return nil,err}
	if len(raw)>32<<20{return nil,errors.New("embedding response too large")}
	var out struct{Data []struct{Index int `json:"index"`;Embedding []float32 `json:"embedding"`} `json:"data"`}
	if err=json.Unmarshal(raw,&out);err!=nil{return nil,err}
	if len(out.Data)!=len(inputs){return nil,errors.New("embedding response count mismatch")}
	vectors:=make([][]float32,len(inputs))
	dimension:=0
	for _,item:=range out.Data{
		if item.Index<0||item.Index>=len(inputs)||vectors[item.Index]!=nil{return nil,errors.New("invalid embedding index")}
		v:=item.Embedding
		if len(v)==0||len(v)>4096{return nil,errors.New("invalid embedding dimensions")}
		if dimension==0{dimension=len(v)}
		if len(v)!=dimension{return nil,errors.New("inconsistent embedding dimensions")}
		var norm float64
		for _,x:=range v{norm+=float64(x)*float64(x)}
		if norm<=0||math.IsNaN(norm)||math.IsInf(norm,0){return nil,errors.New("invalid embedding norm")}
		norm=math.Sqrt(norm)
		for i:=range v{v[i]=float32(float64(v[i])/norm)}
		vectors[item.Index]=v
	}
	return vectors,nil
}

func Split(text string)[]db.KnowledgeVector{
	runes:=[]rune(text)
	out:=[]db.KnowledgeVector{}
	for start:=0;start<len(runes);{
		end:=start+chunkSize
		if end>len(runes){end=len(runes)}
		out=append(out,db.KnowledgeVector{Offset:start,Snippet:string(runes[start:end])})
		if end==len(runes){break}
		start=end-chunkOverlap
	}
	return out
}

// Build indexes a saved document. If edited during embedding, it fails rather
// than publishing vectors for an obsolete revision.
func Build(ctx context.Context,pg *db.DB,taskID int64,docID string,c Config)(int,error){
	result,err:=build(ctx,pg,taskID,docID,c,true)
	return result.Chunks,err
}

// BuildResult distinguishes cached chunks from actual model requests.
type BuildResult struct {
	Chunks int `json:"chunks"`
	Reused int `json:"reused"`
	Embedded int `json:"embedded"`
	Skipped bool `json:"skipped"`
}

func BuildIncremental(ctx context.Context,pg *db.DB,taskID int64,docID string,c Config)(BuildResult,error){
	return build(ctx,pg,taskID,docID,c,false)
}

func build(ctx context.Context,pg *db.DB,taskID int64,docID string,c Config,force bool)(BuildResult,error){
	var result BuildResult
	c,err:=c.Normalize()
	if err!=nil{return result,err}
	if !c.Enabled{return result,errors.New("enable local embedding first")}
	doc,err:=pg.GetKnowledgeDocument(ctx,taskID,false,docID)
	if err!=nil{return result,err}
	chunks:=Split(doc.Content)
	if len(chunks)==0{return result,errors.New("document has no text")}
	result.Chunks=len(chunks)
	dim:=0
	if !force{
		count,err:=pg.KnowledgeVectorCount(ctx,taskID,docID,c.Key())
		if err!=nil{return result,err}
		if count==len(chunks){
			result.Skipped=true
			result.Reused=count
			return result,nil
		}
		cached,err:=pg.ReusableKnowledgeVectors(ctx,taskID,docID,c.Key())
		if err!=nil{return result,err}
		for i:=range chunks{
			v:=cached[chunks[i].Snippet]
			if len(v)==0{continue}
			if dim==0{dim=len(v)}
			if len(v)!=dim{return result,errors.New("cached embedding dimensions differ; force rebuild")}
			chunks[i].Embedding=v
			result.Reused++
		}
	}
	var pending []int
	for i:=range chunks{if chunks[i].Embedding==nil{pending=append(pending,i)}}
	for start:=0;start<len(pending);start+=batchSize{
		end:=start+batchSize
		if end>len(pending){end=len(pending)}
		inputs:=make([]string,end-start)
		for i,j:=range pending[start:end]{inputs[i]=chunks[j].Snippet}
		vectors,err:=Embed(ctx,c,inputs)
		if err!=nil{return result,err}
		for i,v:=range vectors{
			if dim==0{dim=len(v)}
			if len(v)!=dim{return result,errors.New("embedding dimensions changed; force rebuild")}
			chunks[pending[start+i]].Embedding=v
			result.Embedded++
		}
	}
	if err=pg.ReplaceKnowledgeVectors(ctx,taskID,doc.ID,c.Key(),doc.UpdatedAt,chunks);err!=nil{return result,err}
	return result,nil
}

// Search makes no model call unless current-scope indexed documents exist.
func Search(ctx context.Context,pg *db.DB,taskID int64,q string,limit int)([]db.KnowledgeSemanticHit,bool,error){
	c,err:=Load(pg)
	if err!=nil{return nil,false,err}
	if !c.Enabled{return nil,false,nil}
	exists,err:=pg.HasKnowledgeVectors(ctx,taskID,c.Key())
	if err!=nil||!exists{return nil,false,err}
	ctx,cancel:=context.WithTimeout(ctx,8*time.Second)
	defer cancel()
	v,err:=Embed(ctx,c,[]string{q})
	if err!=nil{return nil,true,err}
	hits,err:=pg.SearchKnowledgeVectors(ctx,taskID,c.Key(),v[0],limit)
	return hits,true,err
}
