package routing

import (
	"reflect"
	"testing"
)

func TestForbiddenTurnRequiresDifferentArrivalState(t *testing.T) {
	roads := []Road{{"ab", "a", "b", 1}, {"ac", "a", "c", 2}, {"cb", "c", "b", 1}, {"bd", "b", "d", 1}}
	r, err := FindWithTurns(roads, []Turn{{"ab", "bd"}}, "a", "d")
	if err != nil || r.Cost != 4 || !reflect.DeepEqual(r.Nodes, []string{"ac", "cb", "bd"}) {
		t.Fatal(r, err)
	}
}
