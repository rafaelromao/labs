package singleflight

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCallerCancellationDoesNotCancelSharedWork(t *testing.T) {
	started, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	finish := func() { once.Do(func() { close(release) }) }
	defer finish()
	var calls atomic.Int32
	g := New(1, context.Background(), func(ctx context.Context, key string) (int, error) {
		calls.Add(1)
		close(started)
		<-release
		return 7, nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan error, 1)
	go func() { _, err := g.Do(ctx, "a"); result <- err }()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("load not started")
	}
	if _, err := g.Do(context.Background(), "b"); !errors.Is(err, ErrCapacity) {
		t.Fatal(err)
	}
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("caller stuck")
	}
	// The shared flight remains despite its first caller leaving: it still holds the only slot.
	if _, err := g.Do(context.Background(), "b"); !errors.Is(err, ErrCapacity) {
		t.Fatal("shared work removed prematurely", err)
	}
	finish()
	if calls.Load() != 1 {
		t.Fatal("duplicate load")
	}
}

func TestFailureCleanupAllowsRetry(t *testing.T) {
	var calls atomic.Int32
	g := New(1, context.Background(), func(context.Context, string) (int, error) {
		if calls.Add(1) == 1 {
			return 0, errors.New("failed")
		}
		return 9, nil
	})
	if _, err := g.Do(context.Background(), "a"); err == nil {
		t.Fatal("missing error")
	}
	if value, err := g.Do(context.Background(), "a"); err != nil || value != 9 {
		t.Fatal(value, err)
	}
}
