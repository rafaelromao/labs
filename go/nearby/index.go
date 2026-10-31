// Package nearby provides spherical search with conservative latitude-band candidates.
package nearby

import (
	"math"
)

const EarthRadiusMeters = 6371008.8

type Point struct{ Latitude, Longitude float64 }

type Hit struct {
	ID             string
	DistanceMeters float64
}

type Index struct{}

func New() *Index {
	panic("TODO: implement")
}

// valid is provided for Distance; reuse it for your own coordinate checks.
func valid(p Point) bool {
	return !math.IsNaN(p.Latitude) && !math.IsNaN(p.Longitude) && p.Latitude >= -90 && p.Latitude <= 90 && p.Longitude >= -180 && p.Longitude <= 180
}

// Put atomically inserts or moves an object, using latitude then longitude degrees.
func (i *Index) Put(id string, p Point) error {
	panic("TODO: implement")
}

func (i *Index) Delete(id string) bool {
	panic("TODO: implement")
}

// Distance is provided (geofence and trajectory use it): haversine distance in
// metres on a sphere of radius EarthRadiusMeters. It panics on invalid coordinates.
func Distance(a, b Point) float64 {
	if !valid(a) || !valid(b) {
		panic("nearby: invalid coordinates")
	}
	radians := math.Pi / 180
	dlat := (b.Latitude - a.Latitude) * radians
	dlon := (b.Longitude - a.Longitude) * radians
	h := math.Pow(math.Sin(dlat/2), 2) + math.Cos(a.Latitude*radians)*math.Cos(b.Latitude*radians)*math.Pow(math.Sin(dlon/2), 2)
	return 2 * EarthRadiusMeters * math.Asin(math.Sqrt(math.Max(0, math.Min(1, h))))
}

// Radius includes points at the boundary under floating-point spherical distance.
func (i *Index) Radius(center Point, meters float64) ([]Hit, error) {
	panic("TODO: implement")
}

func (i *Index) Nearest(center Point, k int) ([]Hit, error) {
	panic("TODO: implement")
}
