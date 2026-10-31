// Package windowmetrics maintains an exact bounded processing-watermark window.
package windowmetrics

import (
	"errors"
	"sort"
	"sync"
)

type Event struct {
	Time  int64
	Value int64
}
type Snapshot struct {
	Count      int
	Maximum    int64
	HasMaximum bool
	Top        []int64
}
type Window struct {
	mu          sync.Mutex
	width       int64
	capacity    int
	watermark   int64
	initialized bool
	events      []Event
}

func New(width int64, capacity int) *Window {
	if width <= 0 || capacity <= 0 {
		panic("windowmetrics: positive width and capacity required")
	}
	return &Window{width: width, capacity: capacity}
}

// Advance retains events in (now-width, now]. Time must be nonnegative and cannot regress.
func (w *Window) advance(now int64) error {
	if now < 0 || (w.initialized && now < w.watermark) {
		return errors.New("windowmetrics: invalid or late time")
	}
	w.watermark = now
	w.initialized = true
	i := 0
	for i < len(w.events) && now-w.events[i].Time >= w.width {
		i++
	}
	copy(w.events, w.events[i:])
	w.events = w.events[:len(w.events)-i]
	return nil
}

// Add advances the watermark even if admission subsequently fails at capacity.
func (w *Window) Add(event Event) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.advance(event.Time); err != nil {
		return err
	}
	if len(w.events) == w.capacity {
		return errors.New("windowmetrics: capacity reached")
	}
	w.events = append(w.events, event)
	return nil
}
func (w *Window) Query(now int64, k int) (Snapshot, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if k < 0 {
		return Snapshot{}, errors.New("windowmetrics: negative k")
	}
	if err := w.advance(now); err != nil {
		return Snapshot{}, err
	}
	values := make([]int64, len(w.events))
	for i, e := range w.events {
		values[i] = e.Value
	}
	sort.Slice(values, func(i, j int) bool { return values[i] > values[j] })
	s := Snapshot{Count: len(values), HasMaximum: len(values) > 0}
	if len(values) > 0 {
		s.Maximum = values[0]
	}
	if k > len(values) {
		k = len(values)
	}
	s.Top = values[:k]
	return s, nil
}
