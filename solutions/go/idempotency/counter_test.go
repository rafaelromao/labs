package idempotency

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
)

func TestConcurrentDuplicatesConflictAndReopen(t *testing.T) {
	file := filepath.Join(t.TempDir(), "counter.db")
	a, err := Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	b, err := Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c := a
			if i%2 == 1 {
				c = b
			}
			value, err := c.Add(context.Background(), "tenant", "key", 7)
			if err != nil || value != 7 {
				t.Errorf("got %d %v", value, err)
			}
		}(i)
	}
	wg.Wait()
	if _, err := a.Add(context.Background(), "tenant", "key", 8); !errors.Is(err, ErrConflict) {
		t.Fatal(err)
	}
	reopened, err := Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	if value, err := reopened.Add(context.Background(), "tenant", "key", 7); err != nil || value != 7 {
		t.Fatal(value, err)
	}
	if value, err := reopened.Add(context.Background(), "tenant", "next", 1); err != nil || value != 8 {
		t.Fatal(value, err)
	}
}
func TestRollbackAfterBusinessMutation(t *testing.T) {
	c, err := Open(filepath.Join(t.TempDir(), "counter.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	_, err = c.db.Exec(`CREATE TRIGGER fail_insert BEFORE INSERT ON requests WHEN NEW.key='fail' BEGIN SELECT RAISE(ABORT,'injected'); END;`)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Add(context.Background(), "tenant", "fail", 7); err == nil {
		t.Fatal("injection failed")
	}
	if value, err := c.Add(context.Background(), "tenant", "ok", 1); err != nil || value != 1 {
		t.Fatal(value, err)
	}
}
