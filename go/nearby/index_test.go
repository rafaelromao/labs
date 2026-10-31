package nearby

import (
	"math"
	"testing"
)

func TestAntimeridianMovementAndTies(t *testing.T) {
	i := New()
	for id, p := range map[string]Point{"east": {0, 179.9}, "west": {0, -179.9}, "far": {0, 0}} {
		if err := i.Put(id, p); err != nil {
			t.Fatal(err)
		}
	}
	hits, err := i.Radius(Point{0, 180}, 12000)
	if err != nil || len(hits) != 2 {
		t.Fatal(hits, err)
	}
	if err := i.Put("far", Point{0, 180}); err != nil {
		t.Fatal(err)
	}
	hits, err = i.Nearest(Point{0, 180}, 1)
	if err != nil || len(hits) != 1 || hits[0].ID != "far" {
		t.Fatal(hits, err)
	}
	if !i.Delete("far") || i.Delete("far") {
		t.Fatal("delete contract")
	}
	_ = i.Put("a", Point{1, 1})
	_ = i.Put("b", Point{1, 1})
	hits, err = i.Radius(Point{1, 1}, 0)
	if err != nil || len(hits) != 2 || hits[0].ID != "a" || hits[1].ID != "b" {
		t.Fatal(hits, err)
	}
}

func TestInvalidAndAntipodal(t *testing.T) {
	i := New()
	if err := i.Put("a", Point{91, 0}); err == nil {
		t.Fatal("invalid latitude")
	}
	if _, err := i.Radius(Point{}, math.NaN()); err == nil {
		t.Fatal("invalid radius")
	}
	if _, err := i.Nearest(Point{}, -1); err == nil {
		t.Fatal("negative k")
	}
	if d := Distance(Point{0, 0}, Point{0, 180}); math.Abs(d-math.Pi*EarthRadiusMeters) > 0.01 {
		t.Fatal(d)
	}
}
