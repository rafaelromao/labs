// Package dag executes finite dependency graphs in bounded parallel waves.
package dag

import (
	"context"
	"errors"
	"sort"

	"github.com/rafaelromao/labs/solutions/go/fanout"
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
	if concurrency <= 0 || maxTasks <= 0 || len(tasks) > maxTasks {
		return nil, errors.New("dag: invalid limits")
	}
	remaining := make(map[string]int)
	children := make(map[string][]string)
	copied := make(map[string]Task[V])
	for id, t := range tasks {
		if t.Run == nil {
			return nil, errors.New("dag: nil task")
		}
		seen := make(map[string]bool)
		for _, dep := range t.DependsOn {
			if _, ok := tasks[dep]; !ok || seen[dep] {
				return nil, errors.New("dag: missing or duplicate dependency")
			}
			seen[dep] = true
			children[dep] = append(children[dep], id)
		}
		remaining[id] = len(t.DependsOn)
		copied[id] = Task[V]{append([]string(nil), t.DependsOn...), t.Run}
	}
	// Validate using a separate indegree copy, before invoking any callback.
	degree := make(map[string]int)
	ready := []string{}
	for id, n := range remaining {
		degree[id] = n
		if n == 0 {
			ready = append(ready, id)
		}
	}
	visited := 0
	for len(ready) > 0 {
		id := ready[len(ready)-1]
		ready = ready[:len(ready)-1]
		visited++
		for _, child := range children[id] {
			degree[child]--
			if degree[child] == 0 {
				ready = append(ready, child)
			}
		}
	}
	if visited != len(tasks) {
		return nil, errors.New("dag: cycle")
	}
	result := make(map[string]Outcome[V])
	ready = nil
	for id, n := range remaining {
		if n == 0 {
			ready = append(ready, id)
		}
	}
	for len(ready) > 0 {
		sort.Strings(ready)
		calls := make(map[string]func(context.Context) (V, error))
		for _, id := range ready {
			blocked := false
			for _, dep := range copied[id].DependsOn {
				if prior := result[dep]; prior.Err != nil || prior.Skipped {
					blocked = true
				}
			}
			if blocked {
				result[id] = Outcome[V]{Skipped: true}
			} else {
				calls[id] = copied[id].Run
			}
		}
		outcomes, err := fanout.Run(ctx, calls, concurrency, maxTasks)
		if err != nil {
			return nil, err
		}
		for id, out := range outcomes {
			result[id] = Outcome[V]{Value: out.Value, Err: out.Err}
		}
		next := []string{}
		for _, id := range ready {
			for _, child := range children[id] {
				remaining[child]--
				if remaining[child] == 0 {
					next = append(next, child)
				}
			}
		}
		ready = next
	}
	return result, nil
}
