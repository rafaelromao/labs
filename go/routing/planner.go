// Package routing searches immutable snapshots of directed weighted graphs.
package routing

import ()

type Edge struct {
	To   string
	Cost float64
}

type Route struct {
	Nodes []string
	Cost  float64
	Found bool
}

// Find uses Dijkstra when heuristic is nil, otherwise A*. For optimality the
// heuristic must be finite, nonnegative, admissible, and zero at goal. Inputs
// must not be mutated during the call. All destination nodes must be declared.
func Find(graph map[string][]Edge, start, goal string, heuristic func(string) float64) (Route, error) {
	panic("TODO: implement")
}
