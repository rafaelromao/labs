package routing

import ()

type Road struct {
	ID, From, To string
	Cost         float64
}

type Turn struct{ Incoming, Outgoing string }

// FindWithTurns searches incoming-edge states, not junction-only states.
// Returned Nodes are road IDs. Zero heuristic preserves Dijkstra optimality.
func FindWithTurns(roads []Road, blocked []Turn, start, goal string) (Route, error) {
	panic("TODO: implement")
}
