package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBoundedRetryAndProbe(t *testing.T) {
	now := time.Unix(0, 0)
	b := New(2, time.Second, func() time.Time { return now })
	calls := 0
	task := func(context.Context) error { calls++; return errors.New("failed") }
	eligible := func(error) bool { return true }
	delay := func(int) time.Duration { return 0 }
	if err := b.Execute(context.Background(), 5, task, eligible, delay); !errors.Is(err, ErrOpen) || calls != 2 {
		t.Fatal(err, calls)
	}
	now = now.Add(time.Second)
	p, err := b.acquire()
	if err != nil || !p.probe {
		t.Fatal(p, err)
	}
	if _, err := b.acquire(); !errors.Is(err, ErrOpen) {
		t.Fatal(err)
	}
	b.finish(p, nil)
	if err := b.Execute(context.Background(), 1, func(context.Context) error { return nil }, eligible, delay); err != nil {
		t.Fatal(err)
	}
}
func TestLateSuccessCannotCloseNewOpenCircuit(t *testing.T) {
	b := New(1, time.Second, time.Now)
	old, _ := b.acquire()
	failed, _ := b.acquire()
	b.finish(failed, errors.New("failed"))
	b.finish(old, nil)
	if _, err := b.acquire(); !errors.Is(err, ErrOpen) {
		t.Fatal(err)
	}
}
