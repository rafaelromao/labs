// Package spatial uses bounded planar integer coordinates, NOT longitude/latitude.
package spatial

import (
	"errors"
	"sort"
	"sync"
)

const Limit int64 = 1000000

type Point struct{ X, Y int64 }
type Polygon struct {
	shell []Point
	holes [][]Point
}

func valid(p Point) bool        { return p.X >= -Limit && p.X <= Limit && p.Y >= -Limit && p.Y <= Limit }
func cross(a, b, c Point) int64 { return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X) }
func on(a, b, p Point) bool {
	return cross(a, b, p) == 0 && p.X >= min(a.X, b.X) && p.X <= max(a.X, b.X) && p.Y >= min(a.Y, b.Y) && p.Y <= max(a.Y, b.Y)
}
func sign(v int64) int {
	if v < 0 {
		return -1
	}
	if v > 0 {
		return 1
	}
	return 0
}
func intersects(a, b, c, d Point) bool {
	return on(a, b, c) || on(a, b, d) || on(c, d, a) || on(c, d, b) || (sign(cross(a, b, c))*sign(cross(a, b, d)) < 0 && sign(cross(c, d, a))*sign(cross(c, d, b)) < 0)
}

// relation returns -1 outside, 0 boundary, 1 inside using exact orientation.
func relation(ring []Point, p Point) int {
	inside := false
	for i, a := range ring {
		b := ring[(i+1)%len(ring)]
		if on(a, b, p) {
			return 0
		}
		if (a.Y > p.Y) != (b.Y > p.Y) {
			c := cross(a, b, p)
			if (c > 0) == (b.Y > a.Y) {
				inside = !inside
			}
		}
	}
	if inside {
		return 1
	}
	return -1
}
func ringValid(r []Point) bool {
	if len(r) < 3 || len(r) > 4096 {
		return false
	}
	area := int64(0)
	for i, a := range r {
		b := r[(i+1)%len(r)]
		if !valid(a) || a == b {
			return false
		}
		area += a.X*b.Y - a.Y*b.X
		for j := i + 1; j < len(r); j++ {
			if j == i+1 || (i == 0 && j == len(r)-1) {
				continue
			}
			if intersects(a, b, r[j], r[(j+1)%len(r)]) {
				return false
			}
		}
		c := r[(i+2)%len(r)]
		if cross(a, b, c) == 0 {
			return false
		}
	}
	return area != 0
}
func ringsIntersect(a, b []Point) bool {
	for i, p := range a {
		for j, q := range b {
			if intersects(p, a[(i+1)%len(a)], q, b[(j+1)%len(b)]) {
				return true
			}
		}
	}
	return false
}
func NewPolygon(shell []Point, holes [][]Point) (*Polygon, error) {
	if !ringValid(shell) {
		return nil, errors.New("spatial: invalid shell")
	}
	p := &Polygon{shell: append([]Point(nil), shell...)}
	for _, h := range holes {
		if !ringValid(h) || relation(shell, h[0]) != 1 || ringsIntersect(shell, h) {
			return nil, errors.New("spatial: invalid hole")
		}
		for _, other := range p.holes {
			if ringsIntersect(other, h) || relation(other, h[0]) >= 0 || relation(h, other[0]) >= 0 {
				return nil, errors.New("spatial: overlapping holes")
			}
		}
		p.holes = append(p.holes, append([]Point(nil), h...))
	}
	return p, nil
}
func (p *Polygon) Covers(point Point) bool {
	if !valid(point) {
		panic("spatial: coordinate domain")
	}
	r := relation(p.shell, point)
	if r < 0 {
		return false
	}
	if r == 0 {
		return true
	}
	for _, h := range p.holes {
		r = relation(h, point)
		if r == 0 {
			return true
		}
		if r > 0 {
			return false
		}
	}
	return true
}

type membership struct {
	sequence int64
	point    Point
	inside   bool
}
type Transition struct {
	ID                 string
	Inside             bool
	Revision, Sequence int64
}
type Fence struct {
	mu       sync.Mutex
	shape    *Polygon
	revision int64
	capacity int
	objects  map[string]membership
}

func NewFence(shape *Polygon, revision int64, capacity int) *Fence {
	if shape == nil || revision < 0 || capacity <= 0 {
		panic("invalid fence")
	}
	return &Fence{shape: shape, revision: revision, capacity: capacity, objects: make(map[string]membership)}
}
func (f *Fence) Observe(id string, sequence int64, p Point) (*Transition, error) {
	if sequence < 0 || !valid(p) {
		return nil, errors.New("invalid observation")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	old, exists := f.objects[id]
	if exists && sequence <= old.sequence {
		return nil, nil
	}
	if !exists && len(f.objects) == f.capacity {
		return nil, errors.New("capacity")
	}
	inside := f.shape.Covers(p)
	f.objects[id] = membership{sequence, p, inside}
	if inside == old.inside {
		return nil, nil
	}
	return &Transition{id, inside, f.revision, sequence}, nil
}

// Replace reclassifies retained observations atomically; caller publishes results outside the lock.
func (f *Fence) Replace(shape *Polygon, revision int64) ([]Transition, error) {
	if shape == nil {
		return nil, errors.New("nil shape")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if revision <= f.revision {
		return nil, errors.New("stale revision")
	}
	changes := []Transition{}
	for id, m := range f.objects {
		inside := shape.Covers(m.point)
		if inside != m.inside {
			changes = append(changes, Transition{id, inside, revision, m.sequence})
		}
		m.inside = inside
		f.objects[id] = m
	}
	f.shape = shape
	f.revision = revision
	sort.Slice(changes, func(i, j int) bool { return changes[i].ID < changes[j].ID })
	return changes, nil
}
