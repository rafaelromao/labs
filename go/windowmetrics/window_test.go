package windowmetrics

import "testing"

func TestWindowBoundsLateDataAndCapacity(t *testing.T) {
	w := New(10, 2)
	if err := w.Add(Event{0, 3}); err != nil {
		t.Fatal(err)
	}
	_ = w.Add(Event{5, 9})
	if err := w.Add(Event{5, 1}); err == nil {
		t.Fatal("capacity")
	}
	s, err := w.Query(9, 2)
	if err != nil || s.Count != 2 || s.Maximum != 9 || len(s.Top) != 2 || s.Top[1] != 3 {
		t.Fatal(s, err)
	}
	s, err = w.Query(10, 2)
	if err != nil || s.Count != 1 || s.Maximum != 9 {
		t.Fatal(s, err)
	}
	if err := w.Add(Event{9, 7}); err == nil {
		t.Fatal("late event")
	}
	if err := w.Add(Event{10, 7}); err != nil {
		t.Fatal(err)
	}
	s, err = w.Query(20, 2)
	if err != nil || s.Count != 0 || s.HasMaximum {
		t.Fatal(s, err)
	}
}
