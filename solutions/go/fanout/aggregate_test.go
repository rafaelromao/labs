package fanout

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

func TestBoundAndPartialFailure(t *testing.T) {
	var active, peak atomic.Int32
	tasks := map[string]func(context.Context) (int, error){}
	for _, key := range []string{"a", "b", "c", "d"} {
		tasks[key] = func(context.Context) (int, error) {
			n := active.Add(1)
			defer active.Add(-1)
			for old := peak.Load(); n > old; old = peak.Load() {
				if peak.CompareAndSwap(old, n) {
					break
				}
			}
			return 7, nil
		}
	}
	tasks["failure"] = func(context.Context) (int, error) { return 0, errors.New("unavailable") }
	result, err := Run(context.Background(), tasks, 2, 5)
	if err != nil || len(result) != 5 || result["a"].Value != 7 || result["failure"].Err == nil || peak.Load() > 2 {
		t.Fatal(result, err, peak.Load())
	}
}
func TestCancellationSkipsUnstartedCalls(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result, err := Run(ctx, map[string]func(context.Context) (int, error){"a": func(context.Context) (int, error) { t.Error("task started"); return 0, nil }}, 1, 1)
	if err != nil || !errors.Is(result["a"].Err, context.Canceled) {
		t.Fatal(result, err)
	}
}
