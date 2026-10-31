// Package windowmetrics maintains an exact bounded processing-watermark window.
package windowmetrics

import ()

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

type Window struct{}

func New(width int64, capacity int) *Window {
	panic("TODO: implement")
}

// Add advances the watermark even if admission subsequently fails at capacity.
func (w *Window) Add(event Event) error {
	panic("TODO: implement")
}

func (w *Window) Query(now int64, k int) (Snapshot, error) {
	panic("TODO: implement")
}
