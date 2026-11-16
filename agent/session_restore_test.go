package agent

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"net/http"
	"net/http/httptest"
	"io"
	"testing"

	"github.com/Autumn-27/norma/llm"
	"github.com/Autumn-27/norma/transcript"
)

func TestChatSessionRestorePreservesHistoryAndRejectsCorruption(t *testing.T) {
	var calls atomic.Int32
	var request atomic.Value
	srv:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		calls.Add(1)
		body,err:=io.ReadAll(r.Body)
		if err!=nil{t.Error(err)}
		request.Store(string(body))
		w.Header().Set("Content-Type","application/json")
		_,_=io.WriteString(w,`{"id":"msg_fixture","type":"message","role":"assistant","model":"fixture","content":[{"type":"text","text":"fixture response"}],"stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":10,"output_tokens":2}}`)
	}))
	defer srv.Close()
	p,err:=(Config{Format:llm.FormatAnthropic,BaseURL:srv.URL,APIKey:"fixture",Model:"fixture"}).NewProvider()
	if err!=nil{t.Fatal(err)}
	root:=t.TempDir()
	store:=transcript.NewStore(filepath.Join(root,"transcripts"))
	chat:=NewChatAgent(p,"fixture",root,store,200000)
	chat.SetNonStreaming(func()bool{return true})
	run:=func(id,msg string)(string,error){return chat.Chat(context.Background(),"assistant",id,msg,4,0,false,nil)}
	if text,err:=run("fresh","first user message");err!=nil||!strings.Contains(text,"fixture response"){t.Fatal(text,err)}
	if text,err:=run("fresh","second user message");err!=nil||!strings.Contains(text,"fixture response"){t.Fatal(text,err)}
	if got,_:=request.Load().(string);!strings.Contains(got,"first user message")||!strings.Contains(got,"second user message"){t.Fatal("prior conversation not restored")}
	before:=calls.Load()
	// Start with a valid transcript then simulate a broken append after a crash.
	writer:=store.NewWriter("damaged","")
	writer.RecordMessage(llm.UserText("preserve this history"),llm.Usage{})
	path:=store.MainPath("damaged")
	f,err:=os.OpenFile(path,os.O_APPEND|os.O_WRONLY,0600);if err!=nil{t.Fatal(err)}
	if _,err=f.WriteString("{broken JSON\n");err!=nil{f.Close();t.Fatal(err)};if err=f.Close();err!=nil{t.Fatal(err)}
	original,err:=os.ReadFile(path);if err!=nil{t.Fatal(err)}
	if _,err=run("damaged","must not become an empty-history request");err==nil||!strings.Contains(err.Error(),"restore session"){t.Fatalf("expected restore error, got %v",err)}
	if calls.Load()!=before{t.Fatal("provider called after failed restore")}
	after,err:=os.ReadFile(path);if err!=nil||string(after)!=string(original){t.Fatal("damaged transcript changed",err)}
}
