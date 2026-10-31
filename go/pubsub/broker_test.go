package pubsub

import "testing"

func TestSlowConsumerDropAndDrain(t *testing.T) {
	b := New[int](2, 1)
	a, _ := b.Subscribe()
	c, _ := b.Subscribe()
	if _, err := b.Subscribe(); err == nil {
		t.Fatal("subscriber bound")
	}
	d, err := b.Publish(1)
	if err != nil || d.Delivered != 2 {
		t.Fatal(d, err)
	}
	if v := <-a.Messages(); v != 1 {
		t.Fatal(v)
	}
	d, err = b.Publish(2)
	if err != nil || d.Delivered != 1 || d.Dropped != 1 {
		t.Fatal(d, err)
	}
	c.Unsubscribe()
	c.Unsubscribe()
	if v, ok := <-c.Messages(); !ok || v != 1 {
		t.Fatal(v, ok)
	}
	if _, ok := <-c.Messages(); ok {
		t.Fatal("not closed")
	}
	b.Close()
	b.Close()
	if v := <-a.Messages(); v != 2 {
		t.Fatal(v)
	}
	if _, ok := <-a.Messages(); ok {
		t.Fatal("not closed")
	}
	if _, err := b.Publish(3); err == nil {
		t.Fatal("publish after close")
	}
}
