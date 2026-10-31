// Package outbox demonstrates transactional production and idempotent consumption.
package outbox

import (
	"context"
	_ "github.com/mattn/go-sqlite3"
)

type Event struct {
	ID     string
	Amount int64
}

type Store struct{}

func Open(filename string) (*Store, error) {
	panic("TODO: implement")
}

func (s *Store) Close() error {
	panic("TODO: implement")
}

// Produce atomically writes a business ledger entry and an outgoing event.
func (s *Store) Produce(ctx context.Context, e Event) (bool, error) {
	panic("TODO: implement")
}

// Consume atomically records receipt and one local projection effect.
func (s *Store) Consume(ctx context.Context, e Event) (bool, error) {
	panic("TODO: implement")
}

func (s *Store) Pending(ctx context.Context, limit int) ([]Event, error) {
	panic("TODO: implement")
}

// Ack is called only after confirmed delivery; it does not coordinate with a broker.
func (s *Store) Ack(ctx context.Context, id string) error {
	panic("TODO: implement")
}

func (s *Store) EffectCount(ctx context.Context) (int64, error) {
	panic("TODO: implement")
}
