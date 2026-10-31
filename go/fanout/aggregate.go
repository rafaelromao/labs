// Package fanout collects all dependency outcomes with bounded concurrency.
package fanout

import (
	"context"
)

type Outcome[V any] struct {
	Value V
	Err   error
}

// Run owns and joins its workers. It returns after all accepted work completes,
// including after cancellation; dependencies must cooperate with ctx.
func Run[V any](ctx context.Context, tasks map[string]func(context.Context) (V, error), concurrency, maxTasks int) (map[string]Outcome[V], error) {
	panic("TODO: implement")
}
