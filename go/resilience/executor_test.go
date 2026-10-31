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
	if err := b.Execute(context.Background(), 1, task, eligible, delay); !errors.Is(err, ErrOpen) || calls != 2 {
		t.Fatal("open circuit must not call the task", err, calls)
	}
	now = now.Add(time.Second)

	// Exactly one probe after the cooldown; concurrent callers are rejected while it runs.
	probing, release := make(chan struct{}), make(chan struct{})
	probe := make(chan error, 1)
	go func() {
		probe <- b.Execute(context.Background(), 1, func(context.Context) error {
			close(probing)
			<-release
			return nil
		}, eligible, delay)
	}()
	<-probing
	if err := b.Execute(context.Background(), 1, func(context.Context) error { return nil }, eligible, delay); !errors.Is(err, ErrOpen) {
		t.Fatal("second caller during probe", err)
	}
	close(release)
	if err := <-probe; err != nil {
		t.Fatal(err)
	}
	if err := b.Execute(context.Background(), 1, func(context.Context) error { return nil }, eligible, delay); err != nil {
		t.Fatal("closed after successful probe", err)
	}
}

func TestLateSuccessCannotCloseNewOpenCircuit(t *testing.T) {
	b := New(1, time.Hour, time.Now)
	eligible := func(error) bool { return false }
	delay := func(int) time.Duration { return 0 }
	started, release := make(chan struct{}), make(chan struct{})
	late := make(chan error, 1)
	go func() {
		late <- b.Execute(context.Background(), 1, func(context.Context) error {
			close(started)
			<-release
			return nil
		}, eligible, delay)
	}()
	<-started
	b.Execute(context.Background(), 1, func(context.Context) error { return errors.New("failed") }, eligible, delay)
	close(release)
	if err := <-late; err != nil {
		t.Fatal(err)
	}
	if err := b.Execute(context.Background(), 1, func(context.Context) error { return nil }, eligible, delay); !errors.Is(err, ErrOpen) {
		t.Fatal("a success admitted before the circuit opened must not close it", err)
	}
}
