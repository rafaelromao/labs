package workerpool

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestSaturationDrainAndClose(t *testing.T) {
	p, err := New(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	started, release := make(chan struct{}), make(chan struct{})
	defer p.Close()
	var once sync.Once
	finish := func() { once.Do(func() { close(release) }) }
	defer finish()
	first, err := p.Submit(func() error { close(started); <-release; return nil })
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("worker did not start")
	}
	second, err := p.Submit(func() error { return errors.New("expected") })
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.Submit(func() error { return nil }); !errors.Is(err, ErrFull) {
		t.Fatalf("got %v", err)
	}
	p.Close()
	p.Close()
	if _, err := p.Submit(func() error { return nil }); !errors.Is(err, ErrClosed) {
		t.Fatalf("got %v", err)
	}
	finish()
	select {
	case err := <-first:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("first task stuck")
	}
	select {
	case err := <-second:
		if err == nil || err.Error() != "expected" {
			t.Fatalf("got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("queued task lost")
	}
	select {
	case <-p.Done():
	case <-time.After(time.Second):
		t.Fatal("pool did not drain")
	}
}

func TestFailureDoesNotKillWorker(t *testing.T) {
	p, _ := New(1, 1)
	defer p.Close()
	result, err := p.Submit(func() error { panic("boom") })
	if err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-result:
		if err == nil {
			t.Fatal("missing panic result")
		}
	case <-time.After(time.Second):
		t.Fatal("missing completion")
	}
	next, err := p.Submit(func() error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	p.Close()
	select {
	case err := <-next:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("worker lost")
	}
	select {
	case <-p.Done():
	case <-time.After(time.Second):
		t.Fatal("drain did not finish")
	}
}
