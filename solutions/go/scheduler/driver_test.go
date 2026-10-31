package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestEarlierDeadlineWakesAndOwnerJoins(t *testing.T) {
	s := New(2)
	origin := time.Now()
	now := func() int64 { return time.Since(origin).Nanoseconds() }
	if err := s.Schedule("later", int64(time.Hour), func() { t.Error("late task ran") }); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx, now) }()
	ran := make(chan struct{})
	if err := s.Schedule("earlier", now(), func() { close(ran) }); err != nil {
		t.Fatal(err)
	}
	select {
	case <-ran:
	case <-time.After(time.Second):
		t.Fatal("missed wakeup")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("owner leak")
	}
}
