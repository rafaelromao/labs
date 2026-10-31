package routing

import (
	"math"
	"reflect"
	"testing"
)

func TestDijkstraAndAdmissibleAStar(t *testing.T) {
	g := map[string][]Edge{"a": {{"b", 4}, {"c", 1}}, "c": {{"b", 1}, {"z", 10}}, "b": {{"z", 1}}, "z": {}, "isolated": {}}
	for _, h := range []func(string) float64{nil, func(node string) float64 { return map[string]float64{"a": 3, "c": 2, "b": 1}[node] }} {
		r, err := Find(g, "a", "z", h)
		if err != nil || !r.Found || r.Cost != 3 || !reflect.DeepEqual(r.Nodes, []string{"a", "c", "b", "z"}) {
			t.Fatal(r, err)
		}
	}
	if r, err := Find(g, "z", "isolated", nil); err != nil || r.Found {
		t.Fatal(r, err)
	}
	if r, err := Find(g, "a", "a", nil); err != nil || !r.Found || r.Cost != 0 {
		t.Fatal(r, err)
	}
}

func TestInvalidEdges(t *testing.T) {
	for _, cost := range []float64{-1, math.NaN(), math.Inf(1)} {
		if _, err := Find(map[string][]Edge{"a": {{"b", cost}}, "b": {}}, "a", "b", nil); err == nil {
			t.Fatal("accepted invalid edge")
		}
	}
}
