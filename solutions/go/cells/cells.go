// Package cells demonstrates pinned H3 and S2 APIs, NOT radius query completeness.
package cells

import (
	"errors"
	"github.com/golang/geo/s2"
	"github.com/uber/h3-go/v4"
	"math"
)

func IDs(latitude, longitude float64, resolution, level int) (h3.Cell, s2.CellID, error) {
	if math.IsNaN(latitude) || math.IsNaN(longitude) || latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 || resolution < 0 || resolution > 15 || level < 0 || level > 30 {
		return 0, 0, errors.New("invalid coordinate or level")
	}
	h, err := h3.LatLngToCell(h3.LatLng{Lat: latitude, Lng: longitude}, resolution)
	if err != nil {
		return 0, 0, err
	}
	s := s2.CellIDFromLatLng(s2.LatLngFromDegrees(latitude, longitude)).Parent(level)
	return h, s, nil
}
