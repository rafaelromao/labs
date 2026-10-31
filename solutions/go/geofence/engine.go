// Package geofence tracks circular spherical fences with spatial hysteresis.
package geofence

import (
	"errors"
	"github.com/rafaelromao/labs/solutions/go/nearby"
	"math"
	"sync"
)

type state struct {
	sequence int64
	inside   bool
}
type Result struct {
	Inside, Changed, Stale bool
	Sequence               int64
}
type Engine struct {
	mu             sync.Mutex
	center         nearby.Point
	radius, margin float64
	limit          int
	objects        map[string]state
}

func New(center nearby.Point, radius, margin float64, limit int) (*Engine, error) {
	if math.IsNaN(center.Latitude) || math.IsNaN(center.Longitude) || center.Latitude < -90 || center.Latitude > 90 || center.Longitude < -180 || center.Longitude > 180 || radius < 0 || margin < 0 || math.IsNaN(radius) || math.IsNaN(margin) || math.IsInf(radius+margin, 0) || limit <= 0 {
		return nil, errors.New("geofence: invalid configuration")
	}
	return &Engine{center: center, radius: radius, margin: margin, limit: limit, objects: make(map[string]state)}, nil
}

// Observe requires increasing per-object sequence numbers. Entry is at distance
// <= radius; an inside object exits only at distance > radius+margin.
func (e *Engine) Observe(id string, sequence int64, position nearby.Point) (Result, error) {
	if sequence < 0 || math.IsNaN(position.Latitude) || math.IsNaN(position.Longitude) || position.Latitude < -90 || position.Latitude > 90 || position.Longitude < -180 || position.Longitude > 180 {
		return Result{}, errors.New("geofence: invalid observation")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	old, exists := e.objects[id]
	if exists && sequence <= old.sequence {
		return Result{Inside: old.inside, Stale: true, Sequence: old.sequence}, nil
	}
	if !exists && len(e.objects) == e.limit {
		return Result{}, errors.New("geofence: object capacity")
	}
	d := nearby.Distance(e.center, position)
	inside := d <= e.radius
	if old.inside {
		inside = d <= e.radius+e.margin
	}
	e.objects[id] = state{sequence, inside}
	return Result{Inside: inside, Changed: inside != old.inside, Sequence: sequence}, nil
}

// Forget releases memory and sequence history; a new observation establishes new state.
func (e *Engine) Forget(id string) { e.mu.Lock(); defer e.mu.Unlock(); delete(e.objects, id) }
