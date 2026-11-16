package server

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
)

func TestBroadcasterOverflowDisconnectsOnlySlowSubscriber(t *testing.T){
	b:=NewBroadcaster()
	slow,stopSlow:=b.Subscribe("task")
	fast,stopFast:=b.Subscribe("task")
	defer stopSlow()
	defer stopFast()
	other,stopOther:=b.Subscribe("other")
	defer stopOther()
	for i:=int64(1);i<=300;i++{
		b.Publish("task",db.Activity{ID:i})
		if got:=<-fast;got.ID!=i{t.Fatal("fast subscriber lost activity",got.ID,i)}
	}
	// The slow stream must terminate after its ordered buffered prefix, not
	// silently deliver a later ID that would prevent database replay.
	var count int64
	for got:=range slow{
		count++
		if got.ID!=count{t.Fatal("wrong buffered prefix",got.ID,count)}
	}
	if count!=256{t.Fatal("unexpected buffered prefix length",count)}
	stopSlow() // overflow already closed it; unsubscribe remains idempotent
	b.Publish("other",db.Activity{ID:999})
	if got:=<-other;got.ID!=999{t.Fatal("other task affected")}
	next,stopNext:=b.Subscribe("task")
	defer stopNext()
	b.Publish("task",db.Activity{ID:301})
	if got:=<-next;got.ID!=301{t.Fatal("reconnect failed")}
	if got:=<-fast;got.ID!=301{t.Fatal("fast subscriber disconnected")}
}

func TestConcurrentActivityPublishMatchesDatabaseOrder(t *testing.T){
	s,f,_:=trafficEvidenceServer(t)
	finding,err:=s.m.pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	task:=s.m.ResolveTask(strconv.FormatInt(*finding.TaskID,10))
	if task==nil{t.Fatal("missing task")}
	e:=NewEngine(s.m)
	_,cursor,err:=task.Store.ActivityList(nil,0,1000)
	if err!=nil{t.Fatal(err)}
	ch,stop:=e.Broadcaster().Subscribe(task.ID)
	defer stop()
	const n=32
	var wg sync.WaitGroup
	for i:=0;i<n;i++{
		wg.Add(1)
		go func(i int){
			defer wg.Done()
			a:=e.emitActivity(task,db.Activity{Worker:"fixture",Kind:"text",Summary:fmt.Sprintf("ordered-%d",i)})
			if a.ID<=0{t.Error("persist failed")}
		}(i)
	}
	wg.Wait()
	items,_,err:=task.Store.ActivityList(nil,cursor,1000)
	if err!=nil||len(items)!=n{t.Fatal("database events",len(items),err)}
	last:=cursor
	for _,item:=range items{
		select{
		case a:=<-ch:
			if a.ID!=item.ID||a.ID<=last{t.Fatal("live order differs from DB",a.ID,item.ID,last)}
			last=a.ID
		default:t.Fatal("missing live activity",item.ID)
		}
	}
}

func TestActivityStreamReconnectReplaysAfterLastEventID(t *testing.T){
	s,f,_:=trafficEvidenceServer(t)
	finding,err:=s.m.pg.GetFinding(f.FindingID)
	if err!=nil||finding==nil||finding.TaskID==nil{t.Fatal("missing task",err)}
	task:=s.m.ResolveTask(strconv.FormatInt(*finding.TaskID,10))
	if task==nil||s.engine==nil{t.Fatal("missing task or engine")}
	ids:=[]int64{}
	for i:=0;i<3;i++{
		a:=s.engine.emitActivity(task,db.Activity{Worker:"fixture",Kind:"text",Summary:fmt.Sprintf("reconnect-%d",i)})
		if a.ID<=0{t.Fatal("persist failed")}
		ids=append(ids,a.ID)
	}
	server:=httptest.NewServer(s.Handler())
	defer server.Close()
	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	req,err:=http.NewRequestWithContext(ctx,"GET",server.URL+"/api/exploration/activity/stream?task="+task.ID+"&since=0",nil)
	if err!=nil{t.Fatal(err)}
	token,err:=signJWT(s.jwtKey)
	if err!=nil{t.Fatal(err)}
	req.Header.Set("Authorization","Bearer "+token)
	req.Header.Set("Last-Event-ID",strconv.FormatInt(ids[0],10))
	resp,err:=server.Client().Do(req)
	if err!=nil{t.Fatal(err)}
	defer resp.Body.Close()
	if resp.StatusCode!=200||!strings.Contains(resp.Header.Get("Content-Type"),"text/event-stream"){t.Fatal(resp.Status,resp.Header)}
	scanner:=bufio.NewScanner(resp.Body)
	got:=[]int64{}
	for len(got)<2&&scanner.Scan(){
		if line:=scanner.Text();strings.HasPrefix(line,"id: "){
			id,err:=strconv.ParseInt(strings.TrimPrefix(line,"id: "),10,64)
			if err!=nil{t.Fatal(err)}
			got=append(got,id)
		}
	}
	if err:=scanner.Err();err!=nil{t.Fatal(err)}
	if len(got)!=2||got[0]!=ids[1]||got[1]!=ids[2]{t.Fatal("incorrect replay",got,ids)}
	cancel()
}
