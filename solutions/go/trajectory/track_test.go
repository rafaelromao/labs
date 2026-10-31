package trajectory

import (
	"github.com/rafaelromao/labs/solutions/go/nearby"
	"testing"
)

func TestOrderingReplacementRetentionAndJumps(t *testing.T) {
	track := New(3)
	for _, o := range []Observation{{2000, nearby.Point{Longitude: 0.01}}, {0, nearby.Point{}}, {1000, nearby.Point{}}, {1000, nearby.Point{Longitude: 0.005}}} {
		if err := track.Put(o); err != nil {
			t.Fatal(err)
		}
	}
	s, err := track.Segments(100)
	if err != nil || len(s) != 2 || s[0].FromMillis != 0 || !s[0].Implausible {
		t.Fatal(s, err)
	}
	_ = track.Put(Observation{3000, nearby.Point{Longitude: 0.01}})
	s, _ = track.Segments(100)
	if len(s) != 2 || s[0].FromMillis != 1000 || s[1].SpeedMetersPerSecond != 0 {
		t.Fatal(s)
	}
	_ = track.Put(Observation{0, nearby.Point{Longitude: 1}})
	s, _ = track.Segments(100)
	if s[0].FromMillis != 1000 {
		t.Fatal("late observation broke retention")
	}
}
