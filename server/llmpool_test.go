package server

import (
	"testing"

	"github.com/Autumn-27/artex/llmpool"
)

// The status list drives the "pooling order" strip in the UI, so it has to match the
// order PoolProfiles actually runs: active first, then priority DESC, then id ASC
// (input arrives id-ordered, so equal priorities must keep their relative order).
func TestSortPoolStatusMatchesChainOrder(t *testing.T) {
	in := []LLMPoolMemberStatus{
		{ProfileID: "1", Name: "low", Priority: 0},
		{ProfileID: "2", Name: "active", Priority: 0, Active: true},
		{ProfileID: "3", Name: "high", Priority: 10},
		{ProfileID: "4", Name: "mid-a", Priority: 5},
		{ProfileID: "5", Name: "mid-b", Priority: 5},
	}
	sortPoolStatus(in)

	want := []string{"active", "high", "mid-a", "mid-b", "low"}
	for i, w := range want {
		if in[i].Name != w {
			got := make([]string, len(in))
			for j, m := range in {
				got[j] = m.Name
			}
			t.Fatalf("order = %v, want %v", got, want)
		}
	}
}

// The active profile heads the chain no matter how low its own priority is —
// that's the documented precedence, and the UI must not imply otherwise.
func TestActiveProfileHeadsStatusList(t *testing.T) {
	in := []LLMPoolMemberStatus{
		{ProfileID: "1", Name: "loud", Priority: 999},
		{ProfileID: "2", Name: "active", Priority: -5, Active: true},
	}
	sortPoolStatus(in)
	if in[0].Name != "active" {
		t.Fatalf("head = %q, want the active profile regardless of priority", in[0].Name)
	}
}

func TestPoolCompactionWindowUsesSmallestActualFailoverMember(t *testing.T) {
	tests := []struct {
		name string
		members []int
		configured int
		want int
	}{
		{name:"large primary, smaller fallback",members:[]int{200000,32000},configured:200000,want:32000},
		{name:"primary already smallest",members:[]int{64000,200000},configured:64000,want:64000},
		{name:"zero member metadata ignored",members:[]int{200000,0},configured:200000,want:200000},
		{name:"another fallback smaller",members:[]int{200000,128000,64000},configured:200000,want:64000},
	}
	for _,tt:=range tests{
		t.Run(tt.name,func(t *testing.T){
			members:=make([]*llmpool.Member,0,len(tt.members))
			for _,w:=range tt.members{members=append(members,&llmpool.Member{WindowTokens:w})}
			p:=llmpool.New(members,nil)
			if got:=poolCompactionWindow(p,tt.configured);got!=tt.want{t.Fatalf("window=%d want %d",got,tt.want)}
		})
	}
	if got:=poolCompactionWindow(nil,96000);got!=96000{t.Fatalf("disabled pool changed window: %d",got)}
}
