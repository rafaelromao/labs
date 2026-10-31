package geofence

import (
	"github.com/rafaelromao/labs/solutions/go/nearby"
	"testing"
)

func TestTransitionsHysteresisAndStale(t *testing.T) {
	e, err := New(nearby.Point{}, 1000, 500, 1)
	if err != nil {
		t.Fatal(err)
	}
	r, err := e.Observe("a", 1, nearby.Point{})
	if err != nil || !r.Changed || !r.Inside {
		t.Fatal(r, err)
	}
	r, _ = e.Observe("a", 2, nearby.Point{Longitude: 0.01})
	if !r.Inside || r.Changed {
		t.Fatal(r)
	}
	r, _ = e.Observe("a", 1, nearby.Point{Longitude: 1})
	if !r.Stale || !r.Inside {
		t.Fatal(r)
	}
	r, _ = e.Observe("a", 3, nearby.Point{Longitude: 0.02})
	if !r.Changed || r.Inside {
		t.Fatal(r)
	}
	if _, err := e.Observe("b", 1, nearby.Point{}); err == nil {
		t.Fatal("capacity")
	}
	e.Forget("a")
	if _, err := e.Observe("b", 1, nearby.Point{}); err != nil {
		t.Fatal(err)
	}
}
