// Package fanout collects all dependency outcomes with bounded concurrency.
package fanout

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
)

type Outcome[V any] struct {
	Value V
	Err   error
}

// Run owns and joins its workers. It returns after all accepted work completes,
// including after cancellation; dependencies must cooperate with ctx.
func Run[V any](ctx context.Context, tasks map[string]func(context.Context) (V, error), concurrency, maxTasks int) (map[string]Outcome[V], error) {
	if concurrency <= 0 || maxTasks <= 0 || len(tasks) > maxTasks {
		return nil, errors.New("fanout: invalid limits or excessive input")
	}
	keys := make([]string, 0, len(tasks))
	copied := make(map[string]func(context.Context) (V, error), len(tasks))
	for key, task := range tasks {
		if task == nil {
			return nil, errors.New("fanout: nil task")
		}
		keys = append(keys, key)
		copied[key] = task
	}
	sort.Strings(keys)
	jobs := make(chan string, len(keys))
	for _, key := range keys {
		jobs <- key
	}
	close(jobs)
	result := make(map[string]Outcome[V], len(keys))
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < min(concurrency, len(keys)); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for key := range jobs {
				out := call(ctx, copied[key])
				mu.Lock()
				result[key] = out
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return result, nil
}
func call[V any](ctx context.Context, task func(context.Context) (V, error)) (out Outcome[V]) {
	defer func() {
		if p := recover(); p != nil {
			out.Err = fmt.Errorf("fanout: panic: %v", p)
		}
	}()
	if out.Err = ctx.Err(); out.Err != nil {
		return out
	}
	out.Value, out.Err = task(ctx)
	return out
}
