package algorithms

import "testing"

func FuzzRangeModels(f *testing.F) {
	f.Add([]byte{0, 2, 1, 5, 2, 255})
	f.Fuzz(func(t *testing.T, ops []byte) {
		if len(ops) > 512 {
			t.Skip()
		}
		a := make([]int64, 16)
		tree := NewFenwick(16)
		segment := NewSegment(16)
		for i := 0; i+1 < len(ops); i += 2 {
			index := int(ops[i] % 16)
			delta := int64(int8(ops[i+1]))
			a[index] += delta
			tree.Add(index, delta)
			segment.Set(index, a[index])
			var want int64
			for j := range a {
				want += a[j]
				if tree.Sum(0, j+1) != want || segment.Sum(0, j+1) != want {
					t.Fatal("array disagreement")
				}
			}
		}
	})
}
