package workerpool

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestAbortAccountsForQueuedAndRunning(t *testing.T) {
	p, _ := New(1, 1)
	defer p.Abort()
	started := make(chan struct{})
	first, err := p.SubmitContext(func(ctx context.Context) error { close(started); <-ctx.Done(); return ctx.Err() })
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("start")
	}
	second, err := p.Submit(func() error { t.Error("queued task ran"); return nil })
	if err != nil {
		t.Fatal(err)
	}
	p.Abort()
	p.Abort()
	for _, result := range []<-chan error{first, second} {
		select {
		case err := <-result:
			if !errors.Is(err, context.Canceled) && !errors.Is(err, ErrAborted) {
				t.Fatal(err)
			}
		case <-time.After(time.Second):
			t.Fatal("completion")
		}
	}
	select {
	case <-p.Done():
	case <-time.After(time.Second):
		t.Fatal("leaked workers")
	}
}
