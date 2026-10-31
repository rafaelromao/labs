// Package nearby provides spherical search with conservative latitude-band candidates.
package nearby

import (
	"errors"
	"math"
	"sort"
	"sync"
)

const EarthRadiusMeters = 6371008.8

type Point struct{ Latitude, Longitude float64 }
type Hit struct {
	ID             string
	DistanceMeters float64
}
type Index struct {
	mu     sync.RWMutex
	points map[string]Point
	bands  map[int]map[string]struct{}
}

func New() *Index {
	return &Index{points: make(map[string]Point), bands: make(map[int]map[string]struct{})}
}
func band(p Point) int { return int(math.Floor(p.Latitude)) }
func (i *Index) removeBand(id string, p Point) {
	b := band(p)
	delete(i.bands[b], id)
	if len(i.bands[b]) == 0 {
		delete(i.bands, b)
	}
}
func valid(p Point) bool {
	return !math.IsNaN(p.Latitude) && !math.IsNaN(p.Longitude) && p.Latitude >= -90 && p.Latitude <= 90 && p.Longitude >= -180 && p.Longitude <= 180
}

// Put atomically inserts or moves an object, using latitude then longitude degrees.
func (i *Index) Put(id string, p Point) error {
	if !valid(p) {
		return errors.New("nearby: invalid coordinates")
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	if old, ok := i.points[id]; ok {
		i.removeBand(id, old)
	}
	i.points[id] = p
	b := band(p)
	if i.bands[b] == nil {
		i.bands[b] = make(map[string]struct{})
	}
	i.bands[b][id] = struct{}{}
	return nil
}
func (i *Index) Delete(id string) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	_, ok := i.points[id]
	if ok {
		i.removeBand(id, i.points[id])
	}
	delete(i.points, id)
	return ok
}
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
	if !valid(center) || math.IsNaN(meters) || math.IsInf(meters, 0) || meters < 0 {
		return nil, errors.New("nearby: invalid query")
	}
	i.mu.RLock()
	defer i.mu.RUnlock()
	hits := make([]Hit, 0)
	// Angular latitude difference cannot exceed great-circle angular distance.
	// Expand the cover outward for floating arithmetic; refinement is unchanged.
	delta := meters/EarthRadiusMeters*180/math.Pi + 1e-9
	low, high := center.Latitude-delta, center.Latitude+delta
	for b, ids := range i.bands {
		if float64(b) > high || float64(b+1) < low {
			continue
		}
		for id := range ids {
			p := i.points[id]
			d := Distance(center, p)
			if d <= meters {
				hits = append(hits, Hit{id, d})
			}
		}
	}
	sort.Slice(hits, func(a, b int) bool {
		if hits[a].DistanceMeters == hits[b].DistanceMeters {
			return hits[a].ID < hits[b].ID
		}
		return hits[a].DistanceMeters < hits[b].DistanceMeters
	})
	return hits, nil
}

func (i *Index) Nearest(center Point, k int) ([]Hit, error) {
	if k < 0 {
		return nil, errors.New("nearby: negative k")
	}
	hits, err := i.Radius(center, math.Pi*EarthRadiusMeters)
	if err != nil {
		return nil, err
	}
	if k < len(hits) {
		hits = hits[:k]
	}
	return hits, nil
}
