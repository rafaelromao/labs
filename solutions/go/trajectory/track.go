// Package trajectory retains bounded, ordered geographic observations.
package trajectory

import (
	"errors"
	"github.com/rafaelromao/labs/solutions/go/nearby"
	"math"
	"sort"
	"sync"
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
type Track struct {
	mu       sync.Mutex
	capacity int
	samples  []Observation
}

func New(capacity int) *Track {
	if capacity <= 0 {
		panic("trajectory: positive capacity required")
	}
	return &Track{capacity: capacity}
}

// Put accepts out-of-order observations, replacing duplicate timestamps. Retain
// only the newest capacity timestamps; an old insertion may be immediately pruned.
func (t *Track) Put(o Observation) error {
	p := o.Position
	if o.Millis < 0 || math.IsNaN(p.Latitude) || math.IsNaN(p.Longitude) || p.Latitude < -90 || p.Latitude > 90 || p.Longitude < -180 || p.Longitude > 180 {
		return errors.New("trajectory: invalid observation")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	i := sort.Search(len(t.samples), func(i int) bool { return t.samples[i].Millis >= o.Millis })
	if i < len(t.samples) && t.samples[i].Millis == o.Millis {
		t.samples[i] = o
	} else {
		t.samples = append(t.samples, Observation{})
		copy(t.samples[i+1:], t.samples[i:])
		t.samples[i] = o
	}
	if len(t.samples) > t.capacity {
		copy(t.samples, t.samples[1:])
		t.samples = t.samples[:t.capacity]
	}
	return nil
}
func (t *Track) Segments(maxSpeed float64) ([]Segment, error) {
	if maxSpeed <= 0 || math.IsNaN(maxSpeed) || math.IsInf(maxSpeed, 0) {
		return nil, errors.New("trajectory: invalid speed threshold")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	segments := make([]Segment, 0, max(0, len(t.samples)-1))
	for i := 1; i < len(t.samples); i++ {
		a, b := t.samples[i-1], t.samples[i]
		meters := nearby.Distance(a.Position, b.Position)
		speed := meters / (float64(b.Millis-a.Millis) / 1000)
		segments = append(segments, Segment{a.Millis, b.Millis, meters, speed, speed > maxSpeed})
	}
	return segments, nil
}
