// Package singleflight bounds and shares in-flight loads per string key.
package singleflight

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrCapacity = errors.New("singleflight: distinct-key capacity reached")

type flight[V any] struct {
	done  chan struct{}
	value V
	err   error
}

type Group[V any] struct {
	mu       sync.Mutex
	capacity int
	owner    context.Context
	load     func(context.Context, string) (V, error)
	flights  map[string]*flight[V]
	closed   bool
	wg       sync.WaitGroup
}

// New uses owner, not a caller's context, for shared work. The owner must arrange
// cancellation and completion of cooperative loads before releasing dependencies.
func New[V any](capacity int, owner context.Context, load func(context.Context, string) (V, error)) *Group[V] {
	if capacity <= 0 || owner == nil || load == nil {
		panic("singleflight: invalid configuration")
	}
	return &Group[V]{capacity: capacity, owner: owner, load: load, flights: make(map[string]*flight[V])}
}

// Do shares only in-flight work; completed results are not cached. Caller
// cancellation stops this wait without cancelling the shared load.
func (g *Group[V]) Do(ctx context.Context, key string) (V, error) {
	var zero V
	if err := ctx.Err(); err != nil {
		return zero, err
	}
	g.mu.Lock()
	if g.closed {
		g.mu.Unlock()
		return zero, errors.New("singleflight: admission closed")
	}
	if err := g.owner.Err(); err != nil {
		g.mu.Unlock()
		return zero, err
	}
	f := g.flights[key]
	if f == nil {
		if len(g.flights) == g.capacity {
			g.mu.Unlock()
			return zero, ErrCapacity
		}
		f = &flight[V]{done: make(chan struct{})}
		g.flights[key] = f
		g.wg.Add(1)
		go g.run(key, f)
	}
	g.mu.Unlock()
	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	case <-f.done:
		return f.value, f.err
	}
}

func (g *Group[V]) run(key string, f *flight[V]) {
	defer g.wg.Done()
	func() {
		defer func() {
			if p := recover(); p != nil {
				f.err = fmt.Errorf("singleflight: load panicked: %v", p)
			}
		}()
		f.value, f.err = g.load(g.owner, key)
	}()
	g.mu.Lock()
	delete(g.flights, key)
	close(f.done)
	g.mu.Unlock()
}

// CloseAndWait closes admission before joining, preventing Add/Wait races.
// The external owner may cancel its context first; loaders must cooperate.
func (g *Group[V]) CloseAndWait() { g.mu.Lock(); g.closed = true; g.mu.Unlock(); g.wg.Wait() }
