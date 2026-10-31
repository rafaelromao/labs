package outbox

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
)

func TestReplayAfterLostAckAndReopen(t *testing.T) {
	ctx := context.Background()
	p, err := Open(filepath.Join(t.TempDir(), "producer.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	file := filepath.Join(t.TempDir(), "consumer.db")
	c, err := Open(file)
	if err != nil {
		t.Fatal(err)
	}
	e := Event{"event-1", 7}
	if added, err := p.Produce(ctx, e); err != nil || !added {
		t.Fatal(added, err)
	}
	if added, err := p.Produce(ctx, e); err != nil || added {
		t.Fatal(added, err)
	}
	pending, err := p.Pending(ctx, 10)
	if err != nil || len(pending) != 1 {
		t.Fatal(pending, err)
	}
	if added, err := c.Consume(ctx, pending[0]); err != nil || !added {
		t.Fatal(added, err)
	}
	c.Close()
	// Delivery succeeded but producer has not acknowledged it: replay after reopen.
	c, err = Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	pending, err = p.Pending(ctx, 10)
	if err != nil || len(pending) != 1 {
		t.Fatal(pending, err)
	}
	if added, err := c.Consume(ctx, pending[0]); err != nil || added {
		t.Fatal(added, err)
	}
	if err := p.Ack(ctx, e.ID); err != nil {
		t.Fatal(err)
	}
	if n, err := c.EffectCount(ctx); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	if pending, err := p.Pending(ctx, 10); err != nil || len(pending) != 0 {
		t.Fatal(pending, err)
	}
	if _, err := c.Consume(ctx, Event{e.ID, 8}); err == nil {
		t.Fatal("payload conflict")
	}
}
func TestOutboxFailureRollsBackBusinessRow(t *testing.T) {
	file := filepath.Join(t.TempDir(), "producer.db")
	s, err := Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	db := rawDB(t, file)
	_, err = db.Exec(`CREATE TRIGGER reject_event BEFORE INSERT ON outgoing BEGIN SELECT RAISE(ABORT,'injected'); END;`)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Produce(context.Background(), Event{"event", 1}); err == nil {
		t.Fatal("injection failed")
	}
	var n int
	if err := db.QueryRow("SELECT count(*) FROM commands").Scan(&n); err != nil || n != 0 {
		t.Fatal(n, err)
	}
}

// rawDB opens a second connection to the same SQLite file so a test can inject failures.
func rawDB(t *testing.T, file string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}
