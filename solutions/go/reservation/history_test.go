package reservation

import (
	"sync"
	"sync/atomic"
	"testing"
)

type call struct {
	kind          string
	start, finish int64
	accepted      bool
}

// Enumerate at most six orders for three operations, respecting real-time precedence.
func legal(history []call) bool {
	var search func([]bool, map[string]bool, int) bool
	search = func(used []bool, state map[string]bool, count int) bool {
		if count == len(history) {
			return true
		}
		for i, c := range history {
			if used[i] {
				continue
			}
			allowed := true
			for j, other := range history {
				if !used[j] && other.finish < c.start {
					allowed = false
				}
			}
			if !allowed {
				continue
			}
			next := map[string]bool{}
			for k, v := range state {
				next[k] = v
			}
			var answer bool
			if c.kind == "cancel" {
				answer = next["a"]
				delete(next, "a")
			} else {
				answer = next[c.kind] || len(next) == 0
				if answer {
					next[c.kind] = true
				}
			}
			if answer != c.accepted {
				continue
			}
			used[i] = true
			if search(used, next, count+1) {
				return true
			}
			used[i] = false
		}
		return false
	}
	return search(make([]bool, len(history)), map[string]bool{}, 0)
}
func TestBoundedConcurrentHistories(t *testing.T) {
	for round := 0; round < 40; round++ {
		c := New(1)
		var clock atomic.Int64
		history := make([]call, 3)
		var wg sync.WaitGroup
		for i, kind := range []string{"a", "b", "cancel"} {
			wg.Add(1)
			go func(i int, kind string) {
				defer wg.Done()
				start := clock.Add(1)
				var accepted bool
				if kind == "cancel" {
					accepted = c.Cancel("a")
				} else {
					var err error
					accepted, err = c.Book(kind, Interval{0, 10})
					if err != nil {
						t.Error(err)
					}
				}
				history[i] = call{kind, start, clock.Add(1), accepted}
			}(i, kind)
		}
		wg.Wait()
		if !legal(history) {
			t.Fatal(history)
		}
	}
	if legal([]call{{"a", 1, 4, true}, {"b", 2, 3, true}}) {
		t.Fatal("oracle failed to reject deliberately faulty overbooking")
	}
}
