// Package lfucache implements bounded LFU with LRU ties using frequency buckets.
package lfucache

import (
	"container/list"
	"reflect"
	"sync"
)

type entry[K comparable, V any] struct {
	key       K
	value     V
	frequency uint64
	position  *list.Element
}

// Cache must not be copied after use. Values remain caller-owned aliases.
type Cache[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	minimum  uint64
	entries  map[K]*entry[K, V]
	buckets  map[uint64]*list.List // oldest access at front
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <= 0 {
		panic("lfucache: capacity must be positive")
	}
	return &Cache[K, V]{capacity: capacity, entries: make(map[K]*entry[K, V]), buckets: make(map[uint64]*list.List)}
}
func validateKey[K comparable](key K) {
	v := reflect.ValueOf(key)
	if (v.IsValid() && !v.Comparable()) || key != key {
		panic("lfucache: invalid key")
	}
}
func (c *Cache[K, V]) Get(key K) (V, bool) {
	validateKey(key)
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.touch(e)
	return e.value, true
}
func (c *Cache[K, V]) Put(key K, value V) {
	validateKey(key)
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.entries[key]; ok {
		e.value = value
		c.touch(e)
		return
	}
	if len(c.entries) == c.capacity {
		bucket := c.buckets[c.minimum]
		victim := bucket.Front().Value.(*entry[K, V])
		bucket.Remove(victim.position)
		delete(c.entries, victim.key)
		if bucket.Len() == 0 {
			delete(c.buckets, c.minimum)
		}
	}
	e := &entry[K, V]{key: key, value: value, frequency: 1}
	c.append(e)
	c.entries[key] = e
	c.minimum = 1
}
func (c *Cache[K, V]) append(e *entry[K, V]) {
	bucket := c.buckets[e.frequency]
	if bucket == nil {
		bucket = list.New()
		c.buckets[e.frequency] = bucket
	}
	e.position = bucket.PushBack(e)
}
func (c *Cache[K, V]) touch(e *entry[K, V]) {
	old := e.frequency
	bucket := c.buckets[old]
	bucket.Remove(e.position)
	if bucket.Len() == 0 {
		delete(c.buckets, old)
	}
	if e.frequency < ^uint64(0) {
		e.frequency++
	}
	c.append(e)
	if c.minimum == old && c.buckets[old] == nil {
		c.minimum = e.frequency
	}
}
