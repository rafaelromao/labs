package algorithms

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestIndexedQueueModel(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	q := NewQueue()
	model := map[string]int64{}
	for step := 0; step < 3000; step++ {
		key := fmt.Sprint(r.Intn(40))
		if r.Intn(3) == 0 {
			_, exists := model[key]
			if q.Remove(key) != exists {
				t.Fatal("remove")
			}
			delete(model, key)
		} else {
			v := r.Int63n(100) - 50
			q.Put(key, v)
			model[key] = v
		}
		if !q.Valid() || q.Len() != len(model) {
			t.Fatal("heap invariant", step)
		}
		var want Priority
		found := false
		for k, v := range model {
			if !found || v < want.Value || (v == want.Value && k < want.Key) {
				want = Priority{k, v}
				found = true
			}
		}
		got, ok := q.Peek()
		if ok != found || (ok && got != want) {
			t.Fatal(step, got, want)
		}
	}
	var last Priority
	for q.Len() > 0 {
		v, _ := q.Pop()
		if v.Value < last.Value && last.Key != "" {
			t.Fatal("order")
		}
		last = v
		if !q.Valid() {
			t.Fatal("pop")
		}
	}
}
func TestRangeModelsAndSharedFixture(t *testing.T) {
	f := NewFenwick(3)
	s := NewSegment(3)
	a := make([]int64, 3)
	file, err := os.Open("../../../fixtures/structures.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		var i, l, h int
		var d, want int64
		if _, err := fmt.Sscan(line, &i, &d, &l, &h, &want); err != nil {
			t.Fatal(err)
		}
		a[i] += d
		f.Add(i, d)
		s.Set(i, a[i])
		if f.Sum(l, h) != want || s.Sum(l, h) != want {
			t.Fatal(line)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	r := rand.New(rand.NewSource(7))
	a = make([]int64, 37)
	f = NewFenwick(len(a))
	s = NewSegment(len(a))
	for step := 0; step < 2000; step++ {
		i := r.Intn(len(a))
		d := int64(r.Intn(41) - 20)
		a[i] += d
		f.Add(i, d)
		s.Set(i, a[i])
		l := r.Intn(len(a) + 1)
		h := l + r.Intn(len(a)-l+1)
		var want int64
		for _, v := range a[l:h] {
			want += v
		}
		if f.Sum(l, h) != want || s.Sum(l, h) != want {
			t.Fatal(step, l, h)
		}
	}
}
func TestTrieUnionBloomSearch(t *testing.T) {
	var trie Trie
	for _, s := range []string{"", "car", "cart", "cat"} {
		trie.Insert(s)
	}
	if !trie.Delete("car") || trie.Has("car") || !trie.Has("cart") {
		t.Fatal("prefix deletion")
	}
	if !reflect.DeepEqual(trie.Prefix("ca"), []string{"cart", "cat"}) {
		t.Fatal(trie.Prefix("ca"))
	}
	d := NewDisjoint(10)
	labels := make([]int, 10)
	for i := range labels {
		labels[i] = i
	}
	r := rand.New(rand.NewSource(13))
	for step := 0; step < 100; step++ {
		a, b := r.Intn(10), r.Intn(10)
		same := labels[a] == labels[b]
		if d.Union(a, b) == same {
			t.Fatal("union")
		}
		old, new := labels[b], labels[a]
		for i := range labels {
			if labels[i] == old {
				labels[i] = new
			}
		}
		for i := range labels {
			count := 0
			for j := range labels {
				if (labels[i] == labels[j]) != (d.Find(i) == d.Find(j)) {
					t.Fatal("components")
				}
				if labels[j] == labels[i] {
					count++
				}
			}
			if d.Size(i) != count {
				t.Fatal("size")
			}
		}
	}
	bloom := NewBloom(128, 4)
	for i := 0; i < 100; i++ {
		bloom.Add(fmt.Sprint(i))
	}
	for i := 0; i < 100; i++ {
		if !bloom.Has(fmt.Sprint(i)) {
			t.Fatal("false negative")
		}
	}
	a := []int64{2, 2, 7}
	if Lower(a, 2) != 0 || Upper(a, 2) != 2 || Upper(a, 1) != 0 {
		t.Fatal("boundaries")
	}
}
func TestStreamsGraphsAndIntervals(t *testing.T) {
	r := rand.New(rand.NewSource(9))
	for n := 1; n < 60; n++ {
		a := make([]int64, n)
		for i := range a {
			a[i] = int64(r.Intn(100) - 50)
		}
		width := 1 + r.Intn(n)
		got := Maxima(a, width)
		for i, v := range got {
			want := a[i]
			for _, x := range a[i : i+width] {
				if x > want {
					want = x
				}
			}
			if v != want {
				t.Fatal("max")
			}
		}
		sorted := append([]int64{}, a...)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i] > sorted[j] })
		k := r.Intn(n + 1)
		if !reflect.DeepEqual(TopK(a, k), sorted[:k]) {
			t.Fatal("topK")
		}
	}
	if !reflect.DeepEqual(MergeSorted([][]int64{{1, 4}, {}, {2, 3, 5}}), []int64{1, 2, 3, 4, 5}) {
		t.Fatal("merge")
	}
	if !reflect.DeepEqual(MergeIntervals([]Interval{{5, 8}, {0, 5}, {9, 10}}), []Interval{{0, 8}, {9, 10}}) {
		t.Fatal("intervals")
	}
	g := [][]int{{1, 2}, {3}, {3}, {}}
	if len(Traverse(g, 0, false)) != 4 || len(Traverse(g, 0, true)) != 4 {
		t.Fatal("traversal")
	}
	order, ok := Topological(g)
	if !ok || order[0] != 0 || order[3] != 3 {
		t.Fatal(order)
	}
	if _, ok := Topological([][]int{{1}, {0}}); ok {
		t.Fatal("cycle")
	}
}
