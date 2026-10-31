// Package reservation provides atomic weighted half-open reservations.
package reservation

import ()

type Interval struct{ Start, End int64 }

type Calendar struct{}

func New(capacity int) *Calendar {
	panic("TODO: implement")
}

func (c *Calendar) Book(id string, v Interval) (bool, error) {
	panic("TODO: implement")
}

func (c *Calendar) BookUnits(id string, v Interval, units int) (bool, error) {
	panic("TODO: implement")
}

// Reschedule retains the old booking on capacity rejection; a missing ID fails.
func (c *Calendar) Reschedule(id string, v Interval, units int) (bool, error) {
	panic("TODO: implement")
}

func (c *Calendar) Cancel(id string) bool {
	panic("TODO: implement")
}

func (c *Calendar) Expire(now int64) int {
	panic("TODO: implement")
}
