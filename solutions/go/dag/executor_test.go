package dag

import (
	"context"
	"errors"
	"testing"
)

func TestDependenciesAndSkippedDescendants(t *testing.T) {
	tasks := map[string]Task[int]{
		"a": {Run: func(context.Context) (int, error) { return 1, nil }},
		"b": {DependsOn: []string{"a"}, Run: func(context.Context) (int, error) { return 0, errors.New("failed") }},
		"c": {DependsOn: []string{"b"}, Run: func(context.Context) (int, error) { t.Error("blocked task ran"); return 0, nil }},
	}
	result, err := Execute(context.Background(), tasks, 2, 3)
	if err != nil || result["a"].Value != 1 || result["b"].Err == nil || !result["c"].Skipped {
		t.Fatal(result, err)
	}
}
func TestCycleRejectedBeforeEffects(t *testing.T) {
	run := func(context.Context) (int, error) { t.Error("callback ran"); return 0, nil }
	_, err := Execute(context.Background(), map[string]Task[int]{"a": {[]string{"b"}, run}, "b": {[]string{"a"}, run}, "independent": {Run: run}}, 2, 3)
	if err == nil {
		t.Fatal("cycle accepted")
	}
}
