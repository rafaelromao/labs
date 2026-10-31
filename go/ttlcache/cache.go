package ttlcache

import (
	"errors"
	"time"
)

var ErrInvalidConfig = errors.New("ttlcache: capacity and ttl must be positive")

// Cache is a bounded, linearizable LRU cache with a fixed time-to-live.
// Keys must be dynamically comparable and equal to themselves; invalid keys panic
// before mutation. Values are returned as-is; mutable objects remain caller-owned.
type Cache[K comparable, V any] struct{}

func New[K comparable, V any](capacity int, ttl time.Duration) (*Cache[K, V], error) {
	panic("TODO: implement")
}

// NewWithClock injects time for deterministic expiry tests. The supplied clock
// must be safe for concurrent calls and should be monotonic for elapsed-time use.
func NewWithClock[K comparable, V any](capacity int, ttl time.Duration, now func() time.Time) (*Cache[K, V], error) {
	panic("TODO: implement")
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	panic("TODO: implement")
}

func (c *Cache[K, V]) Put(key K, value V) {
	panic("TODO: implement")
}

func (c *Cache[K, V]) Delete(key K) bool {
	panic("TODO: implement")
}

func (c *Cache[K, V]) Len() int {
	panic("TODO: implement")
}
