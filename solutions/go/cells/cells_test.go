package cells

import (
	"github.com/golang/geo/s2"
	"github.com/uber/h3-go/v4"
	"testing"
)

func TestCellRoundTripsAtPolesAndAntimeridian(t *testing.T) {
	for _, p := range [][2]float64{{0, 0}, {89.9, 179.9}, {-89.9, -179.9}, {0, 180}} {
		h, s, err := IDs(p[0], p[1], 7, 12)
		if err != nil {
			t.Fatal(err)
		}
		center, err := h.LatLng()
		if err != nil {
			t.Fatal(err)
		}
		round, err := h3.LatLngToCell(center, 7)
		if err != nil || round != h {
			t.Fatal("H3 round trip", err)
		}
		if s.Level() != 12 || !s2.CellFromCellID(s).ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(p[0], p[1]))) {
			t.Fatal("S2 containment")
		}
	}
}
