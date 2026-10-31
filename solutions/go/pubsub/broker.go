// Package pubsub implements bounded, best-effort local broadcast.
package pubsub

import (
	"errors"
	"sync"
)

type Broker[T any] struct {
	mu              sync.Mutex
	closed          bool
	limit, capacity int
	subscribers     map[*Subscription[T]]struct{}
}
type Subscription[T any] struct {
	broker   *Broker[T]
	messages chan T
}
type Delivery struct{ Delivered, Dropped int }

func New[T any](limit, capacity int) *Broker[T] {
	if limit <= 0 || capacity <= 0 {
		panic("pubsub: positive limits required")
	}
	return &Broker[T]{limit: limit, capacity: capacity, subscribers: make(map[*Subscription[T]]struct{})}
}
func (b *Broker[T]) Subscribe() (*Subscription[T], error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return nil, errors.New("pubsub: closed")
	}
	if len(b.subscribers) == b.limit {
		return nil, errors.New("pubsub: subscriber limit")
	}
	s := &Subscription[T]{broker: b, messages: make(chan T, b.capacity)}
	b.subscribers[s] = struct{}{}
	return s, nil
}
func (s *Subscription[T]) Messages() <-chan T { return s.messages }
func (s *Subscription[T]) Unsubscribe() {
	b := s.broker
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.subscribers[s]; ok {
		delete(b.subscribers, s)
		close(s.messages)
	}
}

// Publish attempts each subscriber once under the broker lock, dropping the new
// item for each full queue. Payloads are shared; callers must not mutate them.
func (b *Broker[T]) Publish(value T) (Delivery, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	result := Delivery{}
	if b.closed {
		return result, errors.New("pubsub: closed")
	}
	for s := range b.subscribers {
		select {
		case s.messages <- value:
			result.Delivered++
		default:
			result.Dropped++
		}
	}
	return result, nil
}
func (b *Broker[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.closed {
		b.closed = true
		for s := range b.subscribers {
			close(s.messages)
			delete(b.subscribers, s)
		}
	}
}
