package streamio

import (
	"context"
	"strings"
	"testing"
)

func FuzzBoundedRecords(f *testing.F) {
	for _, s := range []string{"", "1\n-2\r\n3", "\xff", "9223372036854775807\n1", "123456789"} {
		f.Add(s, uint8(4))
	}
	f.Fuzz(func(t *testing.T, input string, limit uint8) {
		if len(input) > 4096 {
			t.Skip()
		}
		max := 1 + int(limit%64)
		result, err := Process(context.Background(), strings.NewReader(input), max)
		if err != nil && result != (Result{}) {
			t.Fatal("partial result on error")
		}
		if result.Count < 0 {
			t.Fatal("negative count")
		}
	})
}
