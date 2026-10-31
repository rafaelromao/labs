package reservation

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestBoundariesConflictCancelAndExpire(t *testing.T) {
	c := New(1)
	for id, v := range map[string]Interval{"a": {0, 10}, "b": {10, 20}} {
		if ok, err := c.Book(id, v); !ok || err != nil {
			t.Fatal(ok, err)
		}
	}
	if ok, err := c.Book("c", Interval{5, 15}); ok || err != nil {
		t.Fatal(ok, err)
	}
	if ok, err := c.Book("a", Interval{0, 10}); !ok || err != nil {
		t.Fatal(ok, err)
	}
	if _, err := c.Book("a", Interval{0, 11}); err == nil {
		t.Fatal("ID conflict not detected")
	}
	if !c.Cancel("a") || c.Cancel("a") {
		t.Fatal("cancellation contract")
	}
	if n := c.Expire(20); n != 1 {
		t.Fatal(n)
	}
}
func TestConcurrentBookingRespectsCapacity(t *testing.T) {
	c := New(3)
	var wg sync.WaitGroup
	var accepted atomic.Int32
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if ok, err := c.Book(fmt.Sprint(i), Interval{0, 10}); err == nil && ok {
				accepted.Add(1)
			}
		}(i)
	}
	wg.Wait()
	if accepted.Load() != 3 {
		t.Fatal(accepted.Load())
	}
}

func TestWeightedRescheduleRetainsOldOnFailure(t *testing.T) {
	c := New(3)
	if ok, err := c.BookUnits("a", Interval{0, 10}, 2); err != nil || !ok {
		t.Fatal(ok, err)
	}
	if ok, err := c.Book("b", Interval{10, 20}); err != nil || !ok {
		t.Fatal(ok, err)
	}
	if ok, err := c.Reschedule("a", Interval{10, 20}, 3); err != nil || ok {
		t.Fatal(ok, err)
	}
	if ok, err := c.BookUnits("c", Interval{0, 10}, 2); err != nil || ok {
		t.Fatal("old booking lost", ok, err)
	}
	if ok, err := c.Reschedule("a", Interval{20, 30}, 3); err != nil || !ok {
		t.Fatal(ok, err)
	}
}
