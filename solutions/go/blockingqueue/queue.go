// Package blockingqueue implements a bounded FIFO with close and cancellation.
package blockingqueue

import (
	"context"
	"errors"
	"sync"
)

var ErrClosed = errors.New("blockingqueue: closed")
var ErrCapacity = errors.New("blockingqueue: capacity must be positive")

// Queue drains accepted items after close. It must not be copied after use.
type Queue[T any] struct {
	mu       sync.Mutex
	items    []T
	capacity int
	closed   bool
	changed  chan struct{}
}

func New[T any](capacity int) (*Queue[T], error) {
	if capacity <= 0 {
		return nil, ErrCapacity
	}
	return &Queue[T]{capacity: capacity, changed: make(chan struct{})}, nil
}

// Put blocks while full. If cancellation races with progress, either may win.
func (q *Queue[T]) Put(ctx context.Context, value T) error {
	for {
		q.mu.Lock()
		if q.closed {
			q.mu.Unlock()
			return ErrClosed
		}
		if err := ctx.Err(); err != nil {
			q.mu.Unlock()
			return err
		}
		if len(q.items) < q.capacity {
			q.items = append(q.items, value)
			q.signal()
			q.mu.Unlock()
			return nil
		}
		changed := q.changed
		q.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-changed:
		}
	}
}

// Take blocks while empty and open; after close it drains then returns ErrClosed.
func (q *Queue[T]) Take(ctx context.Context) (T, error) {
	for {
		q.mu.Lock()
		if err := ctx.Err(); err != nil {
			q.mu.Unlock()
			var zero T
			return zero, err
		}
		if len(q.items) > 0 {
			value := q.items[0]
			var zero T
			q.items[0] = zero
			q.items = q.items[1:]
			q.signal()
			q.mu.Unlock()
			return value, nil
		}
		if q.closed {
			q.mu.Unlock()
			var zero T
			return zero, ErrClosed
		}
		changed := q.changed
		q.mu.Unlock()
		select {
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		case <-changed:
		}
	}
}

// Close is idempotent and wakes all waiting producers and consumers.
func (q *Queue[T]) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if !q.closed {
		q.closed = true
		q.signal()
	}
}

func (q *Queue[T]) signal() { close(q.changed); q.changed = make(chan struct{}) }
