package routing

import (
	"errors"
	"math"
)

type Road struct {
	ID, From, To string
	Cost         float64
}
type Turn struct{ Incoming, Outgoing string }

// FindWithTurns searches incoming-edge states, not junction-only states.
// Returned Nodes are road IDs. Zero heuristic preserves Dijkstra optimality.
func FindWithTurns(roads []Road, blocked []Turn, start, goal string) (Route, error) {
	byID := map[string]Road{}
	junctions := map[string]bool{}
	for _, r := range roads {
		if r.ID == "" || math.IsNaN(r.Cost) || math.IsInf(r.Cost, 0) || r.Cost < 0 {
			return Route{}, errors.New("invalid road")
		}
		if _, ok := byID[r.ID]; ok {
			return Route{}, errors.New("duplicate road")
		}
		byID[r.ID] = r
		junctions[r.From] = true
		junctions[r.To] = true
	}
	if !junctions[start] || !junctions[goal] {
		return Route{}, errors.New("unknown endpoint")
	}
	forbidden := map[Turn]bool{}
	for _, t := range blocked {
		a, aok := byID[t.Incoming]
		b, bok := byID[t.Outgoing]
		if !aok || !bok || a.To != b.From {
			return Route{}, errors.New("invalid turn")
		}
		forbidden[t] = true
	}
	g := map[string][]Edge{"start": {}, "goal": {}}
	for _, r := range roads {
		g["road:"+r.ID] = []Edge{}
		if r.From == start {
			g["start"] = append(g["start"], Edge{"road:" + r.ID, r.Cost})
		}
		if r.To == goal {
			g["road:"+r.ID] = append(g["road:"+r.ID], Edge{"goal", 0})
		}
	}
	if start == goal {
		g["start"] = append(g["start"], Edge{"goal", 0})
	}
	for _, a := range roads {
		for _, b := range roads {
			if a.To == b.From && !forbidden[Turn{a.ID, b.ID}] {
				g["road:"+a.ID] = append(g["road:"+a.ID], Edge{"road:" + b.ID, b.Cost})
			}
		}
	}
	r, err := Find(g, "start", "goal", nil)
	if err != nil || !r.Found {
		return r, err
	}
	ids := []string{}
	for _, v := range r.Nodes {
		if len(v) >= 5 && v[:5] == "road:" {
			ids = append(ids, v[5:])
		}
	}
	r.Nodes = ids
	return r, nil
}
