// Package outbox demonstrates transactional production and idempotent consumption.
package outbox

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"net/url"
	"path/filepath"
)

type Event struct {
	ID     string
	Amount int64
}
type Store struct{ db *sql.DB }

func Open(filename string) (*Store, error) {
	absolute, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	u := url.URL{Scheme: "file", Path: absolute, RawQuery: "_busy_timeout=5000&_txlock=immediate&_synchronous=FULL"}
	db, err := sql.Open("sqlite3", u.String())
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS commands(id TEXT PRIMARY KEY,amount INTEGER NOT NULL);
CREATE TABLE IF NOT EXISTS outgoing(id TEXT PRIMARY KEY,amount INTEGER NOT NULL,sent INTEGER NOT NULL DEFAULT 0);
CREATE TABLE IF NOT EXISTS inbox(id TEXT PRIMARY KEY,amount INTEGER NOT NULL);
CREATE TABLE IF NOT EXISTS effects(id TEXT PRIMARY KEY,amount INTEGER NOT NULL);`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db}, nil
}
func (s *Store) Close() error { return s.db.Close() }
func valid(e Event) error {
	if e.ID == "" || e.Amount <= 0 {
		return errors.New("outbox: invalid event")
	}
	return nil
}

// transact deduplicates the ID and rejects mismatched payloads, then commits
// two writes together. All SQL identifiers below are fixed internal constants.
func (s *Store) transact(ctx context.Context, e Event, first, second string) (bool, error) {
	if err := valid(e); err != nil {
		return false, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	var prior int64
	err = tx.QueryRowContext(ctx, "SELECT amount FROM "+first+" WHERE id=?", e.ID).Scan(&prior)
	if err == nil {
		if prior != e.Amount {
			return false, errors.New("outbox: payload conflict")
		}
		return false, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO "+first+"(id,amount) VALUES(?,?)", e.ID, e.Amount); err != nil {
		return false, err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO "+second+"(id,amount) VALUES(?,?)", e.ID, e.Amount); err != nil {
		return false, err
	}
	return true, tx.Commit()
}

// Produce atomically writes a business ledger entry and an outgoing event.
func (s *Store) Produce(ctx context.Context, e Event) (bool, error) {
	return s.transact(ctx, e, "commands", "outgoing")
}

// Consume atomically records receipt and one local projection effect.
func (s *Store) Consume(ctx context.Context, e Event) (bool, error) {
	return s.transact(ctx, e, "inbox", "effects")
}
func (s *Store) Pending(ctx context.Context, limit int) ([]Event, error) {
	if limit <= 0 {
		return nil, errors.New("outbox: positive batch size required")
	}
	rows, err := s.db.QueryContext(ctx, "SELECT id,amount FROM outgoing WHERE sent=0 ORDER BY id LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := []Event{}
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Amount); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// Ack is called only after confirmed delivery; it does not coordinate with a broker.
func (s *Store) Ack(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE outgoing SET sent=1 WHERE id=?", id)
	return err
}
func (s *Store) EffectCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx, "SELECT count(*) FROM effects").Scan(&count)
	return count, err
}
