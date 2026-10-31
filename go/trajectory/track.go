// Package trajectory retains bounded, ordered geographic observations.
package trajectory

import (
	"github.com/rafaelromao/labs/go/nearby"
)

type Observation struct {
	Millis   int64
	Position nearby.Point
}

type Segment struct {
	FromMillis, ToMillis         int64
	Meters, SpeedMetersPerSecond float64
	Implausible                  bool
}

type Track struct{}

func New(capacity int) *Track {
	panic("TODO: implement")
}

// Put accepts out-of-order observations, replacing duplicate timestamps. Retain
// only the newest capacity timestamps; an old insertion may be immediately pruned.
func (t *Track) Put(o Observation) error {
	panic("TODO: implement")
}

func (t *Track) Segments(maxSpeed float64) ([]Segment, error) {
	panic("TODO: implement")
}
