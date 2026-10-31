package trajectory

import "testing"

func TestProjectionAndAmbiguousParallelRoads(t *testing.T) {
	p, err := Project(XY{5, 4}, XY{}, XY{10, 0})
	if err != nil || p != (XY{5, 0}) {
		t.Fatal(p, err)
	}
	p, _ = Project(XY{-1, 0}, XY{}, XY{10, 0})
	if p != (XY{}) {
		t.Fatal(p)
	}
	m, err := MatchNearest(XY{5, 1}, []RoadSegment{{"a", XY{}, XY{10, 0}}, {"b", XY{0, 2}, XY{10, 2}}}, 0.01)
	if err != nil || m.ID != "a" || !m.Ambiguous {
		t.Fatal(m, err)
	}
}
