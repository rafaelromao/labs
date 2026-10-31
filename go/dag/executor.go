// Package dag executes finite dependency graphs in bounded parallel waves.
package dag

import (
	"context"
)

type Task[V any] struct {
	DependsOn []string
	Run       func(context.Context) (V, error)
}

type Outcome[V any] struct {
	Value   V
	Err     error
	Skipped bool
}

// Execute validates everything before effects, then processes topological waves.
// Task inputs must not change during the call; results do not feed task arguments.
func Execute[V any](ctx context.Context, tasks map[string]Task[V], concurrency, maxTasks int) (map[string]Outcome[V], error) {
	panic("TODO: implement")
}
