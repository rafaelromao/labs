package lfucache

import "testing"

func TestFrequencyAndRecency(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Get("a")
	c.Put("c", 3)
	if _, ok := c.Get("b"); ok {
		t.Fatal("less frequent b should be evicted")
	}
	c.Put("d", 4)
	if _, ok := c.Get("c"); ok {
		t.Fatal("older frequency-one c should be evicted")
	}
	if value, ok := c.Get("a"); !ok || value != 1 {
		t.Fatal("a missing")
	}
}

func TestUpdateCountsAsAccess(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("a", 9)
	c.Put("c", 3)
	if _, ok := c.Get("b"); ok {
		t.Fatal("b should be evicted")
	}
	if value, ok := c.Get("a"); !ok || value != 9 {
		t.Fatal("updated a missing")
	}
}

func TestRecencyTies(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Get("a")
	c.Get("b")
	c.Put("c", 3)
	if _, ok := c.Get("a"); ok {
		t.Fatal("older frequency-two a should be evicted")
	}
}
