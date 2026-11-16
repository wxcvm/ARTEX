package knowledgeindex

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLocalEmbeddingProtocol(t *testing.T){
	server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		if r.URL.Path!="/v1/embeddings"||r.Method!="POST"{t.Errorf("unexpected request %s %s",r.Method,r.URL.Path)}
		var req struct{Model string `json:"model"`;Input []string `json:"input"`}
		if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil||req.Model!="fixture"||len(req.Input)!=2{t.Errorf("request %+v %v",req,err)}
		// The response order need not match the input order.
		_ = json.NewEncoder(w).Encode(map[string]any{"data":[]any{
			map[string]any{"index":1,"embedding":[]float32{0,2}},
			map[string]any{"index":0,"embedding":[]float32{3,4}},
		}})
	}))
	defer server.Close()
	c:=Config{Enabled:true,BaseURL:server.URL,Model:"fixture"}
	v,err:=Embed(context.Background(),c,[]string{"one","two"})
	if err!=nil||len(v)!=2||len(v[0])!=2||v[0][0]!=0.6||v[0][1]!=0.8||v[1][1]!=1{t.Fatal(v,err)}
	for _,base:=range []string{"http://example.com/v1","http://192.168.1.2/v1","http://127.0.0.1:1234/other","http://user@localhost:1234/v1","http://localhost:1234/v1?q=x"}{
		if _,err:=(Config{Enabled:true,BaseURL:base,Model:"fixture"}).Normalize();err==nil{t.Fatal("accepted invalid endpoint",base)}
	}
	normalized,err:=(Config{Enabled:true,BaseURL:"http://localhost:1234/",Model:"fixture"}).Normalize()
	if err!=nil||normalized.BaseURL!="http://127.0.0.1:1234/v1"{t.Fatal(normalized,err)}
}

func TestEmbeddingRejectsBadResponses(t *testing.T){
	for _,raw:=range []string{
		`{"data":[{"index":0,"embedding":[0,0]}]}`,
		`{"data":[{"index":1,"embedding":[1,0]}]}`,
		`{"data":[{"index":0,"embedding":[]}]}`,
		`{"data":[]} `,
	}{
		server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){_,_=w.Write([]byte(raw))}))
		_,err:=Embed(context.Background(),Config{Enabled:true,BaseURL:server.URL,Model:"fixture"},[]string{"x"})
		server.Close()
		if err==nil{t.Fatal("accepted",raw)}
	}
}
func TestUnicodeChunkOffsets(t *testing.T){
	text:=strings.Repeat("中文",900)
	chunks:=Split(text)
	if len(chunks)!=3||chunks[0].Offset!=0||chunks[1].Offset!=700||chunks[2].Offset!=1400{t.Fatal("wrong offsets",len(chunks))}
	runes:=[]rune(text)
	for _,c:=range chunks{
		if c.Snippet!=string(runes[c.Offset:c.Offset+len([]rune(c.Snippet))]){t.Fatal("wrong Unicode slice")}
	}
}
