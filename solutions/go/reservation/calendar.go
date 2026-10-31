// Package reservation provides atomic weighted half-open reservations.
package reservation

import (
	"errors"
	"sort"
	"sync"
)

type Interval struct{ Start, End int64 }
type booking struct {
	Interval
	units int
}
type Calendar struct {
	mu       sync.Mutex
	capacity int
	entries  map[string]booking
}

func New(capacity int) *Calendar {
	if capacity <= 0 {
		panic("positive capacity required")
	}
	return &Calendar{capacity: capacity, entries: make(map[string]booking)}
}
func (c *Calendar) Book(id string, v Interval) (bool, error) { return c.BookUnits(id, v, 1) }
func (c *Calendar) BookUnits(id string, v Interval, units int) (bool, error) {
	return c.change(id, v, units, false)
}

// Reschedule retains the old booking on capacity rejection; a missing ID fails.
func (c *Calendar) Reschedule(id string, v Interval, units int) (bool, error) {
	return c.change(id, v, units, true)
}
func (c *Calendar) change(id string, v Interval, units int, replace bool) (bool, error) {
	if id == "" || v.Start >= v.End || units <= 0 || units > c.capacity {
		return false, errors.New("invalid booking")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	candidate := booking{v, units}
	old, exists := c.entries[id]
	if replace && !exists {
		return false, errors.New("missing booking")
	}
	if exists && !replace {
		if old != candidate {
			return false, errors.New("ID conflict")
		}
		return true, nil
	}
	type endpoint struct {
		time  int64
		units int
		start bool
	}
	events := make([]endpoint, 0, 2*(len(c.entries)+1))
	add := func(b booking) {
		events = append(events, endpoint{b.Start, b.units, true}, endpoint{b.End, b.units, false})
	}
	for key, b := range c.entries {
		if key != id {
			add(b)
		}
	}
	add(candidate)
	sort.Slice(events, func(i, j int) bool {
		if events[i].time == events[j].time {
			return !events[i].start && events[j].start
		}
		return events[i].time < events[j].time
	})
	used := 0
	for _, e := range events {
		if e.start {
			if e.units > c.capacity-used {
				return false, nil
			}
			used += e.units
		} else {
			used -= e.units
		}
	}
	c.entries[id] = candidate
	return true, nil
}
func (c *Calendar) Cancel(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.entries[id]
	delete(c.entries, id)
	return ok
}
func (c *Calendar) Expire(now int64) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := 0
	for id, b := range c.entries {
		if b.End <= now {
			delete(c.entries, id)
			n++
		}
	}
	return n
}
