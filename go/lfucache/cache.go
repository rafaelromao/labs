// Package lfucache implements bounded LFU with LRU ties using frequency buckets.
package lfucache

import ()

// Cache must not be copied after use. Values remain caller-owned aliases.
type Cache[K comparable, V any] struct{}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	panic("TODO: implement")
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	panic("TODO: implement")
}

func (c *Cache[K, V]) Put(key K, value V) {
	panic("TODO: implement")
}
