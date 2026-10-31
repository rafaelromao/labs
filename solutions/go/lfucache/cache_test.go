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

func TestRecencyTiesAndBucketAgreement(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Get("a")
	c.Get("b")
	c.Put("c", 3)
	if _, ok := c.Get("a"); ok {
		t.Fatal("older frequency-two a should be evicted")
	}
	count := 0
	for frequency, bucket := range c.buckets {
		for p := bucket.Front(); p != nil; p = p.Next() {
			e := p.Value.(*entry[string, int])
			count++
			if c.entries[e.key] != e || e.frequency != frequency || e.position != p {
				t.Fatal("bucket and key map disagree")
			}
		}
	}
	if count != len(c.entries) {
		t.Fatal("bucket size mismatch")
	}
}

func TestSaturatedFrequencyStaysInOneBucket(t *testing.T) {
	c := New[string, int](1)
	c.Put("a", 1)
	e := c.entries["a"]
	e.frequency = ^uint64(0)
	c.buckets[e.frequency] = c.buckets[1]
	delete(c.buckets, 1)
	c.minimum = e.frequency
	if value, ok := c.Get("a"); !ok || value != 1 {
		t.Fatal(value, ok)
	}
	if e.frequency != ^uint64(0) || len(c.buckets) != 1 || c.buckets[e.frequency].Len() != 1 {
		t.Fatal("saturated promotion corrupted buckets")
	}
}
