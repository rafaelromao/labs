package scheduler

import "testing"

func TestDueRescheduleCancelAndReentrant(t *testing.T) {
	s := New(2)
	order := []int{}
	if err := s.Schedule("a", 10, func() { order = append(order, 1) }); err != nil {
		t.Fatal(err)
	}
	_ = s.Schedule("b", 20, func() { order = append(order, 2) })
	_ = s.Schedule("a", 30, func() { order = append(order, 3) })
	if due, ok := s.Next(); !ok || due != 20 {
		t.Fatal(due, ok)
	}
	if n := s.RunDue(19); n != 0 {
		t.Fatal(n)
	}
	if !s.Cancel("b") || s.Cancel("b") {
		t.Fatal("cancel")
	}
	_ = s.Schedule("b", 30, func() { order = append(order, 4); _ = s.Schedule("c", 30, func() { order = append(order, 5) }) })
	if n := s.RunDue(30); n != 2 || len(order) != 2 || order[0] != 3 {
		t.Fatal(n, order)
	}
	if n := s.RunDue(30); n != 1 || order[2] != 5 {
		t.Fatal(n, order)
	}
}
