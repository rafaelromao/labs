package blockingqueue

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestFIFOAndCloseDrains(t *testing.T) {
	q, _ := New[int](2)
	ctx := context.Background()
	if err := q.Put(ctx, 1); err != nil {
		t.Fatal(err)
	}
	if err := q.Put(ctx, 2); err != nil {
		t.Fatal(err)
	}
	q.Close()
	q.Close()
	for _, want := range []int{1, 2} {
		got, err := q.Take(ctx)
		if err != nil || got != want {
			t.Fatalf("got %d, %v; want %d", got, err, want)
		}
	}
	if _, err := q.Take(ctx); !errors.Is(err, ErrClosed) {
		t.Fatal(err)
	}
	if err := q.Put(ctx, 3); !errors.Is(err, ErrClosed) {
		t.Fatal(err)
	}
}

func TestBlockedPutWakesOnClose(t *testing.T) {
	q, _ := New[int](1)
	_ = q.Put(context.Background(), 1)
	result := make(chan error, 1)
	ready := make(chan struct{})
	ctx := &observedContext{Context: context.Background(), ready: ready}
	go func() { result <- q.Put(ctx, 2) }()
	select {
	case <-ready:
	case <-time.After(time.Second):
		t.Fatal("producer did not check predicate")
	}
	q.Close()
	select {
	case err := <-result:
		if !errors.Is(err, ErrClosed) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("producer did not finish")
	}
}

// Err is checked under the queue lock before the full predicate and change
// channel capture. Closing after this signal exercises the check-to-wait gap.
type observedContext struct {
	context.Context
	once  sync.Once
	ready chan struct{}
}

func (c *observedContext) Err() error { c.once.Do(func() { close(c.ready) }); return c.Context.Err() }
func TestBlockedTakeCancelledAfterPredicateCheck(t *testing.T) {
	q, _ := New[int](1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ready := make(chan struct{})
	observed := &observedContext{Context: ctx, ready: ready}
	result := make(chan error, 1)
	go func() { _, err := q.Take(observed); result <- err }()
	select {
	case <-ready:
	case <-time.After(time.Second):
		t.Fatal("predicate")
	}
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("blocked cancellation")
	}
}

func TestCancellationAndInvalidCapacity(t *testing.T) {
	if _, err := New[int](0); !errors.Is(err, ErrCapacity) {
		t.Fatal(err)
	}
	q, _ := New[int](1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := q.Take(ctx); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if err := q.Put(ctx, 1); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
