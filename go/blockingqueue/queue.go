// Package blockingqueue implements a bounded FIFO with close and cancellation.
package blockingqueue

import (
	"context"
	"errors"
)

var ErrClosed = errors.New("blockingqueue: closed")

var ErrCapacity = errors.New("blockingqueue: capacity must be positive")

// Queue drains accepted items after close. It must not be copied after use.
type Queue[T any] struct{}

func New[T any](capacity int) (*Queue[T], error) {
	panic("TODO: implement")
}

// Put blocks while full. If cancellation races with progress, either may win.
func (q *Queue[T]) Put(ctx context.Context, value T) error {
	panic("TODO: implement")
}

// Take blocks while empty and open; after close it drains then returns ErrClosed.
func (q *Queue[T]) Take(ctx context.Context) (T, error) {
	panic("TODO: implement")
}

// Close is idempotent and wakes all waiting producers and consumers.
func (q *Queue[T]) Close() {
	panic("TODO: implement")
}
