package ratelimiter

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBurstRefillAndCardinality(t *testing.T) {
	now := time.Unix(100, 0)
	l, err := New(2, 2, 1, time.Minute, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if ok, err := l.Allow("a"); !ok || err != nil {
			t.Fatal(ok, err)
		}
	}
	if ok, _ := l.Allow("a"); ok {
		t.Fatal("burst exceeded")
	}
	now = now.Add(500 * time.Millisecond)
	if ok, _ := l.Allow("a"); !ok {
		t.Fatal("missing refill")
	}
	if _, err := l.Allow("b"); !errors.Is(err, ErrClients) {
		t.Fatal(err)
	}
	now = now.Add(time.Minute)
	if ok, err := l.Allow("b"); !ok || err != nil {
		t.Fatal(ok, err)
	}
}

func TestConcurrentAdmissionDoesNotExceedBurst(t *testing.T) {
	now := time.Unix(100, 0)
	l, _ := New(1, 5, 1, time.Minute, func() time.Time { return now })
	var allowed atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, err := l.Allow("a")
			if err == nil && ok {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()
	if allowed.Load() != 5 {
		t.Fatalf("allowed %d", allowed.Load())
	}
}
