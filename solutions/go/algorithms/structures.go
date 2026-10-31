// Package algorithms contains sequential teaching structures. Callers own
// synchronization. Range sums require representable signed-64 intermediates.
package algorithms

import (
	"container/heap"
	"crypto/sha256"
	"encoding/binary"
	"sort"
	"strconv"
)

type Priority struct {
	Key   string
	Value int64
}
type IndexedQueue struct {
	items     []Priority
	positions map[string]int
}

func NewQueue() *IndexedQueue    { return &IndexedQueue{positions: make(map[string]int)} }
func (q *IndexedQueue) Len() int { return len(q.items) }
func (q *IndexedQueue) less(a, b int) bool {
	x, y := q.items[a], q.items[b]
	return x.Value < y.Value || (x.Value == y.Value && x.Key < y.Key)
}
func (q *IndexedQueue) swap(a, b int) {
	q.items[a], q.items[b] = q.items[b], q.items[a]
	q.positions[q.items[a].Key] = a
	q.positions[q.items[b].Key] = b
}
func (q *IndexedQueue) fix(i int) {
	if i > 0 && q.less(i, (i-1)/2) {
		for i > 0 && q.less(i, (i-1)/2) {
			p := (i - 1) / 2
			q.swap(i, p)
			i = p
		}
		return
	}
	for {
		child := 2*i + 1
		if child >= q.Len() {
			return
		}
		if child+1 < q.Len() && q.less(child+1, child) {
			child++
		}
		if !q.less(child, i) {
			return
		}
		q.swap(child, i)
		i = child
	}
}

// Put inserts or atomically replaces a priority within this sequential structure.
// region: queue-put
func (q *IndexedQueue) Put(key string, value int64) {
	if i, ok := q.positions[key]; ok {
		q.items[i].Value = value
		q.fix(i)
		return
	}
	q.positions[key] = q.Len()
	q.items = append(q.items, Priority{key, value})
	q.fix(q.Len() - 1)
}

// endregion: queue-put
func (q *IndexedQueue) Peek() (Priority, bool) {
	if q.Len() == 0 {
		return Priority{}, false
	}
	return q.items[0], true
}
func (q *IndexedQueue) Remove(key string) bool {
	i, ok := q.positions[key]
	if !ok {
		return false
	}
	last := q.Len() - 1
	q.swap(i, last)
	q.items = q.items[:last]
	delete(q.positions, key)
	if i < q.Len() {
		q.fix(i)
	}
	return true
}
func (q *IndexedQueue) Pop() (Priority, bool) {
	v, ok := q.Peek()
	if ok {
		q.Remove(v.Key)
	}
	return v, ok
}
func (q *IndexedQueue) Valid() bool {
	if q.Len() != len(q.positions) {
		return false
	}
	for i, v := range q.items {
		if p, ok := q.positions[v.Key]; !ok || p != i {
			return false
		}
		if i > 0 && q.less(i, (i-1)/2) {
			return false
		}
	}
	return true
}

func Lower(a []int64, x int64) int { return sort.Search(len(a), func(i int) bool { return a[i] >= x }) }
func Upper(a []int64, x int64) int { return sort.Search(len(a), func(i int) bool { return a[i] > x }) }

type trieNode struct {
	terminal bool
	children map[byte]*trieNode
}
type Trie struct{ root trieNode }

func ascii(s string) {
	for i := range s {
		if s[i] > 127 {
			panic("algorithms: ASCII required")
		}
	}
}
func (t *Trie) Insert(s string) {
	ascii(s)
	n := &t.root
	for i := range s {
		if n.children == nil {
			n.children = make(map[byte]*trieNode)
		}
		if n.children[s[i]] == nil {
			n.children[s[i]] = &trieNode{}
		}
		n = n.children[s[i]]
	}
	n.terminal = true
}
func (t *Trie) Has(s string) bool {
	ascii(s)
	n := &t.root
	for i := range s {
		n = n.children[s[i]]
		if n == nil {
			return false
		}
	}
	return n.terminal
}
func (t *Trie) Delete(s string) bool {
	ascii(s)
	var remove func(*trieNode, int) bool
	remove = func(n *trieNode, i int) bool {
		if i == len(s) {
			old := n.terminal
			n.terminal = false
			return old
		}
		child := n.children[s[i]]
		if child == nil {
			return false
		}
		ok := remove(child, i+1)
		if !child.terminal && len(child.children) == 0 {
			delete(n.children, s[i])
		}
		return ok
	}
	return remove(&t.root, 0)
}
func (t *Trie) Prefix(s string) []string {
	ascii(s)
	n := &t.root
	for i := range s {
		n = n.children[s[i]]
		if n == nil {
			return []string{}
		}
	}
	out := []string{}
	var walk func(*trieNode, string)
	walk = func(n *trieNode, p string) {
		if n.terminal {
			out = append(out, p)
		}
		for b := 0; b < 128; b++ {
			if c := n.children[byte(b)]; c != nil {
				walk(c, p+string(byte(b)))
			}
		}
	}
	walk(n, s)
	return out
}

type Disjoint struct {
	parent, size []int
	Count        int
}

func NewDisjoint(n int) *Disjoint {
	if n < 0 {
		panic("negative size")
	}
	d := &Disjoint{parent: make([]int, n), size: make([]int, n), Count: n}
	for i := range d.parent {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}
func (d *Disjoint) Find(i int) int {
	if i < 0 || i >= len(d.parent) {
		panic("index")
	}
	for i != d.parent[i] {
		d.parent[i] = d.parent[d.parent[i]]
		i = d.parent[i]
	}
	return i
}
func (d *Disjoint) Union(a, b int) bool {
	a, b = d.Find(a), d.Find(b)
	if a == b {
		return false
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
	d.Count--
	return true
}
func (d *Disjoint) Size(i int) int { return d.size[d.Find(i)] }

type Fenwick struct{ tree []int64 }

func NewFenwick(n int) *Fenwick {
	if n < 0 {
		panic("negative size")
	}
	return &Fenwick{make([]int64, n+1)}
}
func (f *Fenwick) Add(i int, delta int64) {
	if i < 0 || i >= len(f.tree)-1 {
		panic("index")
	}
	for i++; i < len(f.tree); i += i & -i {
		f.tree[i] += delta
	}
}
func (f *Fenwick) Prefix(end int) int64 {
	if end < 0 || end >= len(f.tree) {
		panic("index")
	}
	var v int64
	for ; end > 0; end -= end & -end {
		v += f.tree[end]
	}
	return v
}
func (f *Fenwick) Sum(start, end int) int64 {
	if start < 0 || start > end {
		panic("range")
	}
	return f.Prefix(end) - f.Prefix(start)
}

type SegmentTree struct {
	n    int
	tree []int64
}

func NewSegment(n int) *SegmentTree {
	if n < 0 {
		panic("negative size")
	}
	return &SegmentTree{n, make([]int64, 2*n)}
}
func (s *SegmentTree) Set(i int, v int64) {
	if i < 0 || i >= s.n {
		panic("index")
	}
	i += s.n
	s.tree[i] = v
	for i > 1 {
		i /= 2
		s.tree[i] = s.tree[2*i] + s.tree[2*i+1]
	}
}
func (s *SegmentTree) Sum(a, b int) int64 {
	if a < 0 || b < a || b > s.n {
		panic("range")
	}
	var v int64
	for a, b = a+s.n, b+s.n; a < b; a, b = a/2, b/2 {
		if a&1 != 0 {
			v += s.tree[a]
			a++
		}
		if b&1 != 0 {
			b--
			v += s.tree[b]
		}
	}
	return v
}

type Bloom struct {
	bits   []byte
	hashes int
}

func NewBloom(bytes, hashes int) *Bloom {
	if bytes <= 0 || hashes <= 0 || hashes > 8 {
		panic("invalid bloom configuration")
	}
	return &Bloom{make([]byte, bytes), hashes}
}
func (b *Bloom) indices(key string) []uint64 {
	digest := sha256.Sum256([]byte(key))
	out := make([]uint64, b.hashes)
	for i := range out {
		out[i] = uint64(binary.BigEndian.Uint32(digest[4*i:])) % (uint64(len(b.bits)) * 8)
	}
	return out
}
func (b *Bloom) Add(key string) {
	for _, i := range b.indices(key) {
		b.bits[i/8] |= 1 << uint(i%8)
	}
}
func (b *Bloom) Has(key string) bool {
	for _, i := range b.indices(key) {
		if b.bits[i/8]&(1<<uint(i%8)) == 0 {
			return false
		}
	}
	return true
}
func Maxima(a []int64, width int) []int64 {
	if width <= 0 {
		panic("width")
	}
	q := []int{}
	out := []int64{}
	for i, v := range a {
		for len(q) > 0 && q[0] <= i-width {
			q = q[1:]
		}
		for len(q) > 0 && a[q[len(q)-1]] <= v {
			q = q[:len(q)-1]
		}
		q = append(q, i)
		if i+1 >= width {
			out = append(out, a[q[0]])
		}
	}
	return out
}

type minHeap []int64

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(v any)        { *h = append(*h, v.(int64)) }
func (h *minHeap) Pop() any          { v := (*h)[len(*h)-1]; *h = (*h)[:len(*h)-1]; return v }
func TopK(a []int64, k int) []int64 {
	if k < 0 {
		panic("k")
	}
	h := &minHeap{}
	for _, v := range a {
		if h.Len() < k {
			heap.Push(h, v)
		} else if k > 0 && v > (*h)[0] {
			(*h)[0] = v
			heap.Fix(h, 0)
		}
	}
	out := append([]int64{}, (*h)...)
	sort.Slice(out, func(i, j int) bool { return out[i] > out[j] })
	return out
}
func MergeSorted(streams [][]int64) []int64 {
	q := NewQueue()
	positions := make([]int, len(streams))
	keys := make(map[string]int)
	for i, s := range streams {
		if !sort.SliceIsSorted(s, func(a, b int) bool { return s[a] < s[b] }) {
			panic("unsorted")
		}
		if len(s) > 0 {
			key := strconv.Itoa(i)
			keys[key] = i
			q.Put(key, s[0])
		}
	}
	out := []int64{}
	for q.Len() > 0 {
		head, _ := q.Pop()
		out = append(out, head.Value)
		i := keys[head.Key]
		positions[i]++
		if positions[i] < len(streams[i]) {
			q.Put(head.Key, streams[i][positions[i]])
		}
	}
	return out
}

type Interval struct{ Start, End int64 }

func MergeIntervals(input []Interval) []Interval {
	a := append([]Interval{}, input...)
	for _, v := range a {
		if v.Start >= v.End {
			panic("interval")
		}
	}
	sort.Slice(a, func(i, j int) bool { return a[i].Start < a[j].Start })
	out := []Interval{}
	for _, v := range a {
		if len(out) == 0 || v.Start > out[len(out)-1].End {
			out = append(out, v)
		} else if v.End > out[len(out)-1].End {
			out[len(out)-1].End = v.End
		}
	}
	return out
}

// Traverse visits reachable integer vertices, breadth-first or depth-first.
func Traverse(graph [][]int, start int, depth bool) []int {
	if start < 0 || start >= len(graph) {
		panic("vertex")
	}
	for _, edges := range graph {
		for _, v := range edges {
			if v < 0 || v >= len(graph) {
				panic("vertex")
			}
		}
	}
	seen := make([]bool, len(graph))
	work := []int{start}
	out := []int{}
	for len(work) > 0 {
		var v int
		if depth {
			v = work[len(work)-1]
			work = work[:len(work)-1]
		} else {
			v = work[0]
			work = work[1:]
		}
		if seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
		work = append(work, graph[v]...)
	}
	return out
}
func Topological(g [][]int) ([]int, bool) {
	degree := make([]int, len(g))
	for _, edges := range g {
		for _, v := range edges {
			if v < 0 || v >= len(g) {
				panic("vertex")
			}
			degree[v]++
		}
	}
	q := []int{}
	for i, n := range degree {
		if n == 0 {
			q = append(q, i)
		}
	}
	out := []int{}
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		out = append(out, v)
		for _, w := range g[v] {
			degree[w]--
			if degree[w] == 0 {
				q = append(q, w)
			}
		}
	}
	return out, len(out) == len(g)
}
