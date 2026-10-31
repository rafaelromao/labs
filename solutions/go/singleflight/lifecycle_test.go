package singleflight

import (
	"context"
	"testing"
	"time"
)

func TestOwnerCancellationThenJoin(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started := make(chan struct{})
	g := New(1, ctx, func(ctx context.Context, _ string) (int, error) { close(started); <-ctx.Done(); return 0, ctx.Err() })
	result := make(chan error, 1)
	go func() { _, err := g.Do(context.Background(), "a"); result <- err }()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("start")
	}
	cancel()
	joined := make(chan struct{})
	go func() { g.CloseAndWait(); close(joined) }()
	select {
	case <-joined:
	case <-time.After(time.Second):
		t.Fatal("join")
	}
	select {
	case <-result:
	case <-time.After(time.Second):
		t.Fatal("waiter")
	}
	if _, err := g.Do(context.Background(), "b"); err == nil {
		t.Fatal("admission after close")
	}
}
