package spatial

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
)

func square(a, b int64) []Point { return []Point{{a, a}, {b, a}, {b, b}, {a, b}} }
func TestPolygonSharedFixtureAndRevision(t *testing.T) {
	p, err := NewPolygon(square(0, 10), [][]Point{square(3, 7)})
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open("../../../fixtures/polygon.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		var x, y int64
		var want bool
		if _, err := fmt.Sscan(line, &x, &y, &want); err != nil {
			t.Fatal(err)
		}
		if p.Covers(Point{x, y}) != want {
			t.Fatal(line)
		}
	}
	if err := scan.Err(); err != nil {
		t.Fatal(err)
	}
	f := NewFence(p, 1, 2)
	event, err := f.Observe("a", 1, Point{1, 1})
	if err != nil || event == nil || !event.Inside {
		t.Fatal(event, err)
	}
	smaller, _ := NewPolygon(square(2, 8), nil)
	changes, err := f.Replace(smaller, 2)
	if err != nil || len(changes) != 1 || changes[0].Inside {
		t.Fatal(changes, err)
	}
	if event, err := f.Observe("a", 1, Point{5, 5}); err != nil || event != nil {
		t.Fatal("stale observation", event, err)
	}
}
func TestInvalidGeometry(t *testing.T) {
	for _, shell := range [][]Point{{{0, 0}, {5, 5}, {0, 5}, {5, 0}}, {{0, 0}, {1, 0}, {2, 0}}, {{0, 0}, {0, 0}, {1, 1}}} {
		if _, err := NewPolygon(shell, nil); err == nil {
			t.Fatal(shell)
		}
	}
	if _, err := NewPolygon(square(0, 10), [][]Point{square(0, 2)}); err == nil {
		t.Fatal("touching hole")
	}
	if _, err := NewPolygon(square(0, 10), [][]Point{square(2, 6), square(4, 8)}); err == nil {
		t.Fatal("overlap")
	}
}
func TestQuadTreeAgainstMap(t *testing.T) {
	q := NewQuadTree()
	model := map[string]Point{}
	r := rand.New(rand.NewSource(17))
	for n := 0; n < 1200; n++ {
		id := fmt.Sprint(r.Intn(100))
		p := Point{int64(r.Intn(101) - 50), int64(r.Intn(101) - 50)}
		if n%4 == 0 {
			q.Delete(id)
			delete(model, id)
		} else {
			q.Put(id, p)
			model[id] = p
		}
		box := Box{Point{-20, -20}, Point{20, 20}}
		want := map[string]Point{}
		for id, p := range model {
			if box.Contains(p) {
				want[id] = p
			}
		}
		if !reflect.DeepEqual(q.Range(box), want) {
			t.Fatal(n)
		}
	}
	for i := 0; i < 100; i++ {
		q.Put(fmt.Sprint(i), Point{})
	}
	if len(q.Range(Box{Point{}, Point{}})) != 100 {
		t.Fatal("coincident points")
	}
}
