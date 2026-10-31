// Package routing searches immutable snapshots of directed weighted graphs.
package routing

import (
	"container/heap"
	"errors"
	"math"
)

type Edge struct {
	To   string
	Cost float64
}
type Route struct {
	Nodes []string
	Cost  float64
	Found bool
}
type item struct {
	node           string
	cost, priority float64
}
type queue []item

func (q queue) Len() int { return len(q) }
func (q queue) Less(i, j int) bool {
	if q[i].priority == q[j].priority {
		return q[i].node < q[j].node
	}
	return q[i].priority < q[j].priority
}
func (q queue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q *queue) Push(v any)   { *q = append(*q, v.(item)) }
func (q *queue) Pop() any     { old := *q; n := len(old); v := old[n-1]; *q = old[:n-1]; return v }

// Find uses Dijkstra when heuristic is nil, otherwise A*. For optimality the
// heuristic must be finite, nonnegative, admissible, and zero at goal. Inputs
// must not be mutated during the call. All destination nodes must be declared.
func Find(graph map[string][]Edge, start, goal string, heuristic func(string) float64) (Route, error) {
	bad := func() (Route, error) {
		return Route{}, errors.New("routing: invalid graph, heuristic, or numeric range")
	}
	if _, ok := graph[start]; !ok {
		return bad()
	}
	if _, ok := graph[goal]; !ok {
		return bad()
	}
	h := make(map[string]float64, len(graph))
	for node, edges := range graph {
		value := 0.0
		if heuristic != nil {
			value = heuristic(node)
		}
		if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
			return bad()
		}
		h[node] = value
		for _, e := range edges {
			if _, ok := graph[e.To]; !ok || e.Cost < 0 || math.IsNaN(e.Cost) || math.IsInf(e.Cost, 0) {
				return bad()
			}
		}
	}
	if h[goal] != 0 {
		return bad()
	}
	distance := map[string]float64{start: 0}
	previous := make(map[string]string)
	q := &queue{{node: start, cost: 0, priority: h[start]}}
	heap.Init(q)
	for q.Len() > 0 {
		current := heap.Pop(q).(item)
		if current.cost != distance[current.node] {
			continue
		}
		if current.node == goal {
			nodes := []string{goal}
			for node := goal; node != start; {
				node = previous[node]
				nodes = append(nodes, node)
			}
			for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
				nodes[i], nodes[j] = nodes[j], nodes[i]
			}
			return Route{Nodes: nodes, Cost: current.cost, Found: true}, nil
		}
		for _, edge := range graph[current.node] {
			candidate := current.cost + edge.Cost
			if math.IsInf(candidate, 0) || math.IsInf(candidate+h[edge.To], 0) {
				return bad()
			}
			old, ok := distance[edge.To]
			if !ok || candidate < old {
				distance[edge.To] = candidate
				previous[edge.To] = current.node
				heap.Push(q, item{edge.To, candidate, candidate + h[edge.To]})
			}
		}
	}
	return Route{}, nil
}
