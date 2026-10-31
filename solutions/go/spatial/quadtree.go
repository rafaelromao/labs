package spatial

// QuadTree is a sequential bounded-domain planar index with depth-limited leaves.
// Coincident points stay in a leaf; empty branches are pruned after deletion.
type Box struct{ Min, Max Point }

func (b Box) Contains(p Point) bool {
	return p.X >= b.Min.X && p.X <= b.Max.X && p.Y >= b.Min.Y && p.Y <= b.Max.Y
}
func (b Box) Intersects(c Box) bool {
	return b.Min.X <= c.Max.X && c.Min.X <= b.Max.X && b.Min.Y <= c.Max.Y && c.Min.Y <= b.Max.Y
}

type quadNode struct {
	box      Box
	ids      map[string]Point
	children [4]*quadNode
}
type QuadTree struct {
	root   *quadNode
	points map[string]Point
}

func NewQuadTree() *QuadTree {
	return &QuadTree{root: &quadNode{box: Box{Point{-Limit, -Limit}, Point{Limit, Limit}}, ids: make(map[string]Point)}, points: make(map[string]Point)}
}
func childIndex(b Box, p Point) int {
	i := 0
	if p.X > (b.Min.X+b.Max.X)/2 {
		i++
	}
	if p.Y > (b.Min.Y+b.Max.Y)/2 {
		i += 2
	}
	return i
}
func childBox(b Box, i int) Box {
	x, y := (b.Min.X+b.Max.X)/2, (b.Min.Y+b.Max.Y)/2
	out := b
	if i&1 == 0 {
		out.Max.X = x
	} else {
		out.Min.X = x + 1
	}
	if i&2 == 0 {
		out.Max.Y = y
	} else {
		out.Min.Y = y + 1
	}
	return out
}
func insert(n *quadNode, id string, p Point, depth int) {
	if n.children[0] == nil && (len(n.ids) < 8 || depth == 16) {
		n.ids[id] = p
		return
	}
	if n.children[0] == nil {
		for i := 0; i < 4; i++ {
			n.children[i] = &quadNode{box: childBox(n.box, i), ids: make(map[string]Point)}
		}
		for k, v := range n.ids {
			insert(n.children[childIndex(n.box, v)], k, v, depth+1)
		}
		n.ids = nil
	}
	insert(n.children[childIndex(n.box, p)], id, p, depth+1)
}
func remove(n *quadNode, id string, p Point) bool {
	if n.children[0] == nil {
		delete(n.ids, id)
		return len(n.ids) == 0
	}
	remove(n.children[childIndex(n.box, p)], id, p)
	empty := true
	for _, c := range n.children {
		if c.children[0] != nil || len(c.ids) > 0 {
			empty = false
		}
	}
	if empty {
		n.children = [4]*quadNode{}
		n.ids = make(map[string]Point)
	}
	return empty
}
func (q *QuadTree) Put(id string, p Point) {
	if !valid(p) {
		panic("coordinate domain")
	}
	q.Delete(id)
	q.points[id] = p
	insert(q.root, id, p, 0)
}
func (q *QuadTree) Delete(id string) bool {
	p, ok := q.points[id]
	if !ok {
		return false
	}
	remove(q.root, id, p)
	delete(q.points, id)
	return true
}
func (q *QuadTree) Range(box Box) map[string]Point {
	if !valid(box.Min) || !valid(box.Max) || box.Min.X > box.Max.X || box.Min.Y > box.Max.Y {
		panic("box")
	}
	out := make(map[string]Point)
	var visit func(*quadNode)
	visit = func(n *quadNode) {
		if !n.box.Intersects(box) {
			return
		}
		if n.children[0] != nil {
			for _, c := range n.children {
				visit(c)
			}
		} else {
			for id, p := range n.ids {
				if box.Contains(p) {
					out[id] = p
				}
			}
		}
	}
	visit(q.root)
	return out
}
