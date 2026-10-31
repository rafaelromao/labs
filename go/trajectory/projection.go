package trajectory

import ()

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

func Project(p, a, b XY) (XY, error) {
	panic("TODO: implement")
}

// MatchNearest is a pointwise approximation, not a connected-route map matcher.
func MatchNearest(p XY, segments []RoadSegment, tolerance float64) (Match, error) {
	panic("TODO: implement")
}
