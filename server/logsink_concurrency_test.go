package server

import (
	"sync"
	"testing"
)

func TestLogSinkConcurrentUnsubscribe(t *testing.T) {
	s := &logSinkT{cap: 1024, subs: map[chan LogLine]struct{}{}}
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				_, cancel := s.subscribe()
				cancel()
				cancel()
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				s.add(LogLine{Text: "event"})
			}
		}()
	}
	wg.Wait()
}

func TestLogSinkSlowSubscriberIsolatedAndReplayRetained(t *testing.T) {
	s := &logSinkT{cap: 600, subs: map[chan LogLine]struct{}{}}
	slow, cancelSlow := s.subscribe()
	defer cancelSlow()
	fast, cancelFast := s.subscribe()
	defer cancelFast()
	for i := 1; i <= 300; i++ {
		s.add(LogLine{Text: "event"})
		if got := <-fast; got.Seq != int64(i) {
			t.Fatalf("fast subscriber sequence = %d, want %d", got.Seq, i)
		}
	}
	var n int64
	for l := range slow {
		n++
		if l.Seq != n {
			t.Fatalf("slow prefix sequence = %d, want %d", l.Seq, n)
		}
	}
	if n != 256 {
		t.Fatalf("slow subscriber received %d buffered events, want 256", n)
	}
	replay, cursor := s.recent(n, 600)
	if cursor != 300 || len(replay) != 44 || replay[0].Seq != 257 || replay[43].Seq != 300 {
		t.Fatalf("unexpected replay: %d events, cursor %d", len(replay), cursor)
	}
	cancelSlow()
}

func TestLogSinkConcurrentPublicationPreservesSequence(t *testing.T) {
	s := &logSinkT{cap: 256, subs: map[chan LogLine]struct{}{}, dbCh: make(chan dbWriteReq, 256)}
	ch, cancel := s.subscribe()
	defer cancel()
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 25; j++ {
				s.add(LogLine{Text: "event"})
			}
		}()
	}
	wg.Wait()
	for want := int64(1); want <= 200; want++ {
		if got := <-ch; got.Seq != want {
			t.Fatalf("subscriber sequence = %d, want %d", got.Seq, want)
		}
		if got := <-s.dbCh; got.seq != want {
			t.Fatalf("database queue sequence = %d, want %d", got.seq, want)
		}
	}
}
