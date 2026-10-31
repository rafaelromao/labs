// Package resilience combines bounded eligible retries with a circuit breaker.
package resilience

import (
	"context"
	"errors"
	"time"
)

var ErrOpen = errors.New("resilience: circuit open")

type Breaker struct{}

// States: 0 closed, 1 open, 2 single half-open probe in flight.
func New(threshold int, cooldown time.Duration, now func() time.Time) *Breaker {
	panic("TODO: implement")
}

// Execute counts each attempt against the breaker. backoff supplies a nonnegative
// delay and may include caller-owned jitter. task must honor ctx for hard IO bounds.
func (b *Breaker) Execute(ctx context.Context, maxAttempts int, task func(context.Context) error, retryable func(error) bool, backoff func(int) time.Duration) error {
	panic("TODO: implement")
}
