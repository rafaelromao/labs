package nearby

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

func TestIndexedMovesMatchExhaustiveOracle(t *testing.T) {
	r := rand.New(rand.NewSource(37))
	idx := New()
	model := map[string]Point{}
	for n := 0; n < 1000; n++ {
		id := fmt.Sprint(r.Intn(100))
		p := Point{Latitude: r.Float64()*180 - 90, Longitude: r.Float64()*360 - 180}
		if n%7 == 0 {
			idx.Delete(id)
			delete(model, id)
		} else {
			if err := idx.Put(id, p); err != nil {
				t.Fatal(err)
			}
			model[id] = p
		}
		center := Point{Latitude: r.Float64()*180 - 90, Longitude: r.Float64()*360 - 180}
		radius := r.Float64() * 2e7
		got, err := idx.Radius(center, radius)
		if err != nil {
			t.Fatal(err)
		}
		want := []Hit{}
		for id, p := range model {
			d := Distance(center, p)
			if d <= radius {
				want = append(want, Hit{id, d})
			}
		}
		sort.Slice(want, func(i, j int) bool {
			if want[i].DistanceMeters == want[j].DistanceMeters {
				return want[i].ID < want[j].ID
			}
			return want[i].DistanceMeters < want[j].DistanceMeters
		})
		if len(got) != len(want) {
			t.Fatal(n, "candidate omission")
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatal(n, i)
			}
		}
	}
}
