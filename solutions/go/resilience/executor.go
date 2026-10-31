// Package resilience combines bounded eligible retries with a circuit breaker.
package resilience

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrOpen = errors.New("resilience: circuit open")

type Breaker struct {
	mu                  sync.Mutex
	threshold, failures int
	cooldown            time.Duration
	now                 func() time.Time
	state               int
	until               time.Time
	generation          uint64
}
type permit struct {
	generation uint64
	probe      bool
}

// States: 0 closed, 1 open, 2 single half-open probe in flight.
func New(threshold int, cooldown time.Duration, now func() time.Time) *Breaker {
	if threshold <= 0 || cooldown <= 0 || now == nil {
		panic("resilience: invalid configuration")
	}
	return &Breaker{threshold: threshold, cooldown: cooldown, now: now}
}
func (b *Breaker) acquire() (permit, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == 1 && !b.now().Before(b.until) {
		b.state = 2
		return permit{b.generation, true}, nil
	}
	if b.state != 0 {
		return permit{}, ErrOpen
	}
	return permit{generation: b.generation}, nil
}
func (b *Breaker) finish(p permit, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if p.generation != b.generation {
		return
	}
	if p.probe {
		if err == nil {
			b.state = 0
			b.failures = 0
			b.generation++
		} else {
			b.open()
		}
		return
	}
	if b.state != 0 {
		return
	}
	if err == nil {
		b.failures = 0
	} else {
		b.failures++
		if b.failures >= b.threshold {
			b.open()
		}
	}
}
func (b *Breaker) open() { b.state = 1; b.until = b.now().Add(b.cooldown); b.generation++ }
func invoke(ctx context.Context, task func(context.Context) error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("resilience: task panic: %v", p)
		}
	}()
	return task(ctx)
}

// Execute counts each attempt against the breaker. backoff supplies a nonnegative
// delay and may include caller-owned jitter. task must honor ctx for hard IO bounds.
func (b *Breaker) Execute(ctx context.Context, maxAttempts int, task func(context.Context) error, retryable func(error) bool, backoff func(int) time.Duration) error {
	if maxAttempts <= 0 || task == nil || retryable == nil || backoff == nil {
		return errors.New("resilience: invalid execution configuration")
	}
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		p, err := b.acquire()
		if err != nil {
			return err
		}
		err = invoke(ctx, task)
		b.finish(p, err)
		if err == nil || attempt == maxAttempts || !retryable(err) {
			return err
		}
		delay := backoff(attempt)
		if delay < 0 {
			return errors.New("resilience: negative delay")
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return nil
}
