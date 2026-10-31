// Package geofence tracks circular spherical fences with spatial hysteresis.
package geofence

import (
	"github.com/rafaelromao/labs/go/nearby"
)

type Result struct {
	Inside, Changed, Stale bool
	Sequence               int64
}

type Engine struct{}

func New(center nearby.Point, radius, margin float64, limit int) (*Engine, error) {
	panic("TODO: implement")
}

// Observe requires increasing per-object sequence numbers. Entry is at distance
// <= radius; an inside object exits only at distance > radius+margin.
func (e *Engine) Observe(id string, sequence int64, position nearby.Point) (Result, error) {
	panic("TODO: implement")
}

// Forget releases memory and sequence history; a new observation establishes new state.
func (e *Engine) Forget(id string) {
	panic("TODO: implement")
}
