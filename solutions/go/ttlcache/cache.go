package ttlcache

import (
	"container/list"
	"errors"
	"reflect"
	"sync"
	"time"
)

var ErrInvalidConfig = errors.New("ttlcache: capacity and ttl must be positive")

type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
}

// Cache is a bounded, linearizable LRU cache with a fixed time-to-live.
// Keys must be dynamically comparable and equal to themselves; invalid keys panic
// before mutation. Values are returned as-is; mutable objects remain caller-owned.
type Cache[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	ttl      time.Duration
	now      func() time.Time
	byKey    map[K]*list.Element
	lru      *list.List // front is most-recently used
}

func New[K comparable, V any](capacity int, ttl time.Duration) (*Cache[K, V], error) {
	return NewWithClock[K, V](capacity, ttl, time.Now)
}

// NewWithClock injects time for deterministic expiry tests. The supplied clock
// must be safe for concurrent calls and should be monotonic for elapsed-time use.
func NewWithClock[K comparable, V any](capacity int, ttl time.Duration, now func() time.Time) (*Cache[K, V], error) {
	if capacity <= 0 || ttl <= 0 || now == nil {
		return nil, ErrInvalidConfig
	}
	return &Cache[K, V]{
		capacity: capacity,
		ttl:      ttl,
		now:      now,
		byKey:    make(map[K]*list.Element, capacity),
		lru:      list.New(),
	}, nil
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	validateKey(key)
	c.mu.Lock()
	defer c.mu.Unlock()

	element, ok := c.byKey[key]
	if !ok {
		var zero V
		return zero, false
	}
	item := element.Value.(*entry[K, V])
	if !c.now().Before(item.expiresAt) {
		c.remove(element)
		var zero V
		return zero, false
	}
	c.lru.MoveToFront(element)
	return item.value, true
}

func (c *Cache[K, V]) Put(key K, value V) {
	validateKey(key)
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	c.removeExpired(now)
	if element, ok := c.byKey[key]; ok {
		item := element.Value.(*entry[K, V])
		item.value = value
		item.expiresAt = now.Add(c.ttl)
		c.lru.MoveToFront(element)
		return
	}
	if len(c.byKey) == c.capacity {
		c.remove(c.lru.Back())
	}
	item := &entry[K, V]{key: key, value: value, expiresAt: now.Add(c.ttl)}
	c.byKey[key] = c.lru.PushFront(item)
}

func (c *Cache[K, V]) Delete(key K) bool {
	validateKey(key)
	c.mu.Lock()
	defer c.mu.Unlock()
	element, ok := c.byKey[key]
	if !ok {
		return false
	}
	c.remove(element)
	return true
}

func validateKey[K comparable](key K) {
	value := reflect.ValueOf(key)
	if (value.IsValid() && !value.Comparable()) || key != key {
		panic("ttlcache: key must be comparable and equal to itself")
	}
}

func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.removeExpired(c.now())
	return len(c.byKey)
}

func (c *Cache[K, V]) removeExpired(now time.Time) {
	for element := c.lru.Back(); element != nil; {
		previous := element.Prev()
		item := element.Value.(*entry[K, V])
		if !now.Before(item.expiresAt) {
			c.remove(element)
		}
		element = previous
	}
}

func (c *Cache[K, V]) remove(element *list.Element) {
	if element == nil {
		return
	}
	item := element.Value.(*entry[K, V])
	delete(c.byKey, item.key)
	c.lru.Remove(element)
}
