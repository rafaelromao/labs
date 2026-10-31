// Package temporalstore stores timestamped values with predecessor lookup.
package temporalstore

import (
	"reflect"
	"sort"
	"sync"
)

type version[V any] struct {
	timestamp int64
	value     V
}

// Store is an in-memory, concurrency-safe time-indexed store.
// Timestamps are caller-defined integer units. Values are not copied.
type Store[K comparable, V any] struct {
	mu     sync.RWMutex
	values map[K][]version[V]
}

func New[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{values: make(map[K][]version[V])}
}

// Put replaces the value at an existing timestamp, even for out-of-order writes.
func (s *Store[K, V]) Put(key K, timestamp int64, value V) {
	validateKey(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	items := s.values[key]
	i := sort.Search(len(items), func(i int) bool { return items[i].timestamp >= timestamp })
	if i < len(items) && items[i].timestamp == timestamp {
		items[i].value = value
	} else {
		items = append(items, version[V]{})
		copy(items[i+1:], items[i:])
		items[i] = version[V]{timestamp, value}
	}
	s.values[key] = items
}

// Get returns the value at the greatest timestamp not greater than at.
func (s *Store[K, V]) Get(key K, at int64) (V, bool) {
	validateKey(key)
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := s.values[key]
	i := sort.Search(len(items), func(i int) bool { return items[i].timestamp > at })
	if i == 0 {
		var zero V
		return zero, false
	}
	return items[i-1].value, true
}

func validateKey[K comparable](key K) {
	value := reflect.ValueOf(key)
	if (value.IsValid() && !value.Comparable()) || key != key {
		panic("temporalstore: key must be comparable and equal to itself")
	}
}
