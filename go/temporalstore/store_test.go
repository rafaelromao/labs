package temporalstore

import (
	"sync"
	"testing"
)

func TestPredecessorAndOutOfOrderReplacement(t *testing.T) {
	s := New[string, string]()
	s.Put("a", 20, "twenty")
	s.Put("a", 10, "ten")
	s.Put("a", 20, "updated")
	for _, tt := range []struct {
		at    int64
		value string
		ok    bool
	}{{9, "", false}, {10, "ten", true}, {19, "ten", true}, {20, "updated", true}, {21, "updated", true}} {
		value, ok := s.Get("a", tt.at)
		if value != tt.value || ok != tt.ok {
			t.Fatalf("at %d: got %q, %v", tt.at, value, ok)
		}
	}
	if _, ok := s.Get("missing", 100); ok {
		t.Fatal("unexpected hit")
	}
}

func TestConcurrentWrites(t *testing.T) {
	s := New[string, int]()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) { defer wg.Done(); s.Put("a", int64(i), i); s.Get("a", int64(i)) }(i)
	}
	wg.Wait()
	if value, ok := s.Get("a", 99); !ok || value != 99 {
		t.Fatalf("got %d, %v", value, ok)
	}
}
