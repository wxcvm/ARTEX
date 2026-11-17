package server

import (
 "context"
 "fmt"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestLogSSEReplayAndReconnectCursor(t *testing.T){
 sink:=&logSinkT{cap:20,subs:map[chan LogLine]struct{}{}}
 sink.add(LogLine{Text:"startup1"});sink.add(LogLine{Text:"startup2"})
 sink.prependRestored([]LogLine{{Text:"old1"},{Text:"old2"}})
 request:=func(since,last string)string{
  t.Helper()
  ctx,cancel:=context.WithCancel(context.Background());cancel()
  r:=httptest.NewRequest("GET","/api/logs/stream?since="+since,nil).WithContext(ctx)
  if last!=""{r.Header.Set("Last-Event-ID",last)}
  w:=httptest.NewRecorder();streamLogsFromSink(w,r,sink)
  return w.Body.String()
 }
 all:=request("0","")
 prev:=-3
 for _,line:=range strings.Split(all,"\n"){
  if !strings.HasPrefix(line,"id: "){continue}
  var id int
  if _,err:=fmt.Sscanf(line,"id: %d",&id);err!=nil{t.Fatal(err)}
  if id<=prev{t.Fatalf("out-of-order SSE id=%d after %d",id,prev)}
  prev=id
 }
 if prev!=2||!strings.Contains(all,"id: -2\n"){t.Fatalf("missing history: %q",all)}
 resumed:=request("0","1")
 if !strings.Contains(resumed,"id: 2\n")||strings.Contains(resumed,"id: 1\n")||strings.Contains(resumed,"old1"){t.Fatalf("reconnect cursor failed: %q",resumed)}
 negative:=request("0","-1")
 if !strings.Contains(negative,"id: 1\n")||strings.Contains(negative,"id: -2\n"){t.Fatalf("negative cursor failed: %q",negative)}
}
