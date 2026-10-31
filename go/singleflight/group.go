// Package singleflight bounds and shares in-flight loads per string key.
package singleflight

import (
	"context"
	"errors"
)

var ErrCapacity = errors.New("singleflight: distinct-key capacity reached")

type Group[V any] struct{}

// New uses owner, not a caller's context, for shared work. The owner must arrange
// cancellation and completion of cooperative loads before releasing dependencies.
func New[V any](capacity int, owner context.Context, load func(context.Context, string) (V, error)) *Group[V] {
	panic("TODO: implement")
}

// Do shares only in-flight work; completed results are not cached. Caller
// cancellation stops this wait without cancelling the shared load.
func (g *Group[V]) Do(ctx context.Context, key string) (V, error) {
	panic("TODO: implement")
}

// CloseAndWait closes admission before joining, preventing Add/Wait races.
// The external owner may cancel its context first; loaders must cooperate.
func (g *Group[V]) CloseAndWait() {
	panic("TODO: implement")
}
