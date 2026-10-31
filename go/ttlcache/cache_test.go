package ttlcache

import (
	"fmt"
	"math"
	"sync"
	"testing"
	"time"
)

type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock() *fakeClock { return &fakeClock{now: time.Unix(100, 0)} }

func (f *fakeClock) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.now
}

func (f *fakeClock) Advance(d time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.now = f.now.Add(d)
}

func mustCache(t *testing.T, capacity int, ttl time.Duration, clock *fakeClock) *Cache[string, int] {
	t.Helper()
	cache, err := NewWithClock[string, int](capacity, ttl, clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	return cache
}

func TestExpiresAtInclusiveDeadline(t *testing.T) {
	clock := newFakeClock()
	cache := mustCache(t, 2, time.Minute, clock)
	cache.Put("k", 7)
	clock.Advance(time.Minute - time.Nanosecond)
	if value, ok := cache.Get("k"); !ok || value != 7 {
		t.Fatalf("before deadline: got (%d, %v)", value, ok)
	}
	clock.Advance(time.Nanosecond)
	if _, ok := cache.Get("k"); ok {
		t.Fatal("entry should expire at its deadline")
	}
}

func TestGetPromotesAndEvictsLeastRecentlyUsed(t *testing.T) {
	clock := newFakeClock()
	cache := mustCache(t, 2, time.Hour, clock)
	cache.Put("a", 1)
	cache.Put("b", 2)
	if _, ok := cache.Get("a"); !ok {
		t.Fatal("expected a hit")
	}
	cache.Put("c", 3)
	if _, ok := cache.Get("b"); ok {
		t.Fatal("least-recently-used key b should be evicted")
	}
	if got := cache.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
}

func TestUpdateRefreshesValueDeadlineAndRecency(t *testing.T) {
	clock := newFakeClock()
	cache := mustCache(t, 2, time.Minute, clock)
	cache.Put("a", 1)
	cache.Put("b", 2)
	clock.Advance(30 * time.Second)
	cache.Put("a", 9)
	clock.Advance(30 * time.Second)
	if got, ok := cache.Get("a"); !ok || got != 9 {
		t.Fatalf("updated a: got (%d, %v), want (9, true)", got, ok)
	}
	if _, ok := cache.Get("b"); ok {
		t.Fatal("b should expire at its original deadline")
	}
}

func TestPutSweepsExpiredEntriesBeforeEviction(t *testing.T) {
	clock := newFakeClock()
	cache := mustCache(t, 2, time.Minute, clock)
	cache.Put("a", 1)
	cache.Put("b", 2)
	clock.Advance(time.Minute)
	cache.Put("c", 3)
	cache.Put("d", 4)
	if got := cache.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
	if _, ok := cache.Get("c"); !ok {
		t.Fatal("live entry c was unexpectedly evicted")
	}
}

func TestInvalidConfiguration(t *testing.T) {
	if _, err := New[string, int](0, time.Second); err == nil {
		t.Fatal("expected invalid capacity error")
	}
	if _, err := New[string, int](1, 0); err == nil {
		t.Fatal("expected invalid TTL error")
	}
}

func TestInvalidKeysAreRejectedBeforeMutation(t *testing.T) {
	for _, key := range []any{math.NaN(), struct{ F float64 }{math.NaN()}, []int{1}} {
		t.Run(fmt.Sprintf("%T", key), func(t *testing.T) {
			cache, err := New[any, int](1, time.Hour)
			if err != nil {
				t.Fatal(err)
			}
			cache.Put("stable", 1)
			for _, operation := range []func(){
				func() { cache.Put(key, 2) },
				func() { cache.Get(key) },
				func() { cache.Delete(key) },
			} {
				func() {
					defer func() {
						if recover() == nil {
							t.Error("expected invalid key to panic")
						}
					}()
					operation()
				}()
			}
			if value, ok := cache.Get("stable"); !ok || value != 1 || cache.Len() != 1 {
				t.Fatal("invalid key mutated the cache")
			}
		})
	}
}

func TestNilAndDeleteOfExpiredResident(t *testing.T) {
	clock := newFakeClock()
	cache, err := NewWithClock[any, *int](2, time.Minute, clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	cache.Put(nil, nil)
	if value, ok := cache.Get(nil); !ok || value != nil {
		t.Fatal("nil key and value should be distinguishable from a miss")
	}
	clock.Advance(time.Minute)
	if !cache.Delete(nil) || cache.Delete(nil) {
		t.Fatal("Delete reports resident removal, including expired entries")
	}
}

func TestConcurrentOperationsRemainBounded(t *testing.T) {
	clock := newFakeClock()
	cache := mustCache(t, 64, time.Hour, clock)
	const workers, iterations = 8, 1000
	var wg sync.WaitGroup
	start := make(chan struct{})
	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			<-start
			for i := 0; i < iterations; i++ {
				key := fmt.Sprintf("key-%d", (worker*iterations+i)%128)
				cache.Put(key, i)
				cache.Get(key)
				cache.Len()
			}
		}(worker)
	}
	close(start)
	wg.Wait()
	if got := cache.Len(); got > 64 {
		t.Fatalf("Len() = %d, capacity is 64", got)
	}
}
