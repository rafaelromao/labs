// Package temporalstore stores timestamped values with predecessor lookup.
package temporalstore

import ()

// Store is an in-memory, concurrency-safe time-indexed store.
// Timestamps are caller-defined integer units. Values are not copied.
type Store[K comparable, V any] struct{}

func New[K comparable, V any]() *Store[K, V] {
	panic("TODO: implement")
}

// Put replaces the value at an existing timestamp, even for out-of-order writes.
func (s *Store[K, V]) Put(key K, timestamp int64, value V) {
	panic("TODO: implement")
}

// Get returns the value at the greatest timestamp not greater than at.
func (s *Store[K, V]) Get(key K, at int64) (V, bool) {
	panic("TODO: implement")
}
