// Package pubsub implements bounded, best-effort local broadcast.
package pubsub

import ()

type Broker[T any] struct{}

type Subscription[T any] struct{}

type Delivery struct{ Delivered, Dropped int }

func New[T any](limit, capacity int) *Broker[T] {
	panic("TODO: implement")
}

func (b *Broker[T]) Subscribe() (*Subscription[T], error) {
	panic("TODO: implement")
}

func (s *Subscription[T]) Messages() <-chan T {
	panic("TODO: implement")
}

func (s *Subscription[T]) Unsubscribe() {
	panic("TODO: implement")
}

// Publish attempts each subscriber once under the broker lock, dropping the new
// item for each full queue. Payloads are shared; callers must not mutate them.
func (b *Broker[T]) Publish(value T) (Delivery, error) {
	panic("TODO: implement")
}

func (b *Broker[T]) Close() {
	panic("TODO: implement")
}
