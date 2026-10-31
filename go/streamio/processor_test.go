package streamio

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestBoundariesAndFailures(t *testing.T) {
	for _, tt := range []struct {
		input string
		max   int
		want  Result
		bad   bool
	}{
		{"", 4, Result{}, false}, {"1\r\n-2\n+3", 4, Result{3, 2}, false},
		{"123\n", 3, Result{1, 123}, false}, {"1234", 3, Result{}, true},
		{"\n", 3, Result{}, true}, {"1\nno", 3, Result{}, true},
		{"\xff", 3, Result{}, true}, {"9223372036854775807\n1", 20, Result{}, true},
		{"-9223372036854775808\n-1", 20, Result{}, true},
	} {
		got, err := Process(context.Background(), strings.NewReader(tt.input), tt.max)
		if (err != nil) != tt.bad || got != tt.want {
			t.Fatalf("%q: got %+v, %v", tt.input, got, err)
		}
	}
}

func TestCancelledBeforeRead(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Process(ctx, strings.NewReader("1"), 4); err != context.Canceled {
		t.Fatal(err)
	}
}

type shortReader struct{ input io.Reader }

func (r shortReader) Read(p []byte) (int, error) {
	if len(p) > 1 {
		p = p[:1]
	}
	return r.input.Read(p)
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("injected transport failure") }
func TestShortReadsAndMidstreamFailure(t *testing.T) {
	input := "1\r\n-2\n+3"
	result, err := Process(context.Background(), shortReader{strings.NewReader(input)}, 4)
	if err != nil || result != (Result{3, 2}) {
		t.Fatal(result, err)
	}
	result, err = Process(context.Background(), io.MultiReader(strings.NewReader("1\n"), failingReader{}), 4)
	if err == nil || result != (Result{}) {
		t.Fatal("partial aggregate escaped", result, err)
	}
}
