package trajectory

import (
	"errors"
	"math"
	"sort"
)

// XY are local projected metres, finite and within +/-1,000,000. NOT lat/lon.
type XY struct{ X, Y float64 }
type RoadSegment struct {
	ID   string
	A, B XY
}
type Match struct {
	ID        string
	Position  XY
	Distance  float64
	Ambiguous bool
}

func planar(p XY) bool {
	return !math.IsNaN(p.X) && !math.IsNaN(p.Y) && math.Abs(p.X) <= 1e6 && math.Abs(p.Y) <= 1e6
}
func Project(p, a, b XY) (XY, error) {
	if !planar(p) || !planar(a) || !planar(b) {
		return XY{}, errors.New("invalid projected coordinate")
	}
	dx, dy := b.X-a.X, b.Y-a.Y
	length := dx*dx + dy*dy
	if length == 0 {
		return a, nil
	}
	t := math.Max(0, math.Min(1, ((p.X-a.X)*dx+(p.Y-a.Y)*dy)/length))
	return XY{a.X + t*dx, a.Y + t*dy}, nil
}

// MatchNearest is a pointwise approximation, not a connected-route map matcher.
func MatchNearest(p XY, segments []RoadSegment, tolerance float64) (Match, error) {
	if len(segments) == 0 || tolerance < 0 || math.IsNaN(tolerance) || math.IsInf(tolerance, 0) {
		return Match{}, errors.New("invalid match input")
	}
	matches := []Match{}
	seen := map[string]bool{}
	for _, s := range segments {
		if seen[s.ID] {
			return Match{}, errors.New("duplicate ID")
		}
		seen[s.ID] = true
		q, err := Project(p, s.A, s.B)
		if err != nil {
			return Match{}, err
		}
		matches = append(matches, Match{ID: s.ID, Position: q, Distance: math.Hypot(p.X-q.X, p.Y-q.Y)})
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Distance == matches[j].Distance {
			return matches[i].ID < matches[j].ID
		}
		return matches[i].Distance < matches[j].Distance
	})
	out := matches[0]
	out.Ambiguous = len(matches) > 1 && matches[1].Distance-out.Distance <= tolerance
	return out, nil
}
