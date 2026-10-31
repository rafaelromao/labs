// Package postgres adds the PostgreSQL transaction/outbox adapter. The singleton
// counter row intentionally serializes all commands; this is not a throughput design.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"math"
)

type Store struct{ DB *sql.DB }
type Event struct {
	Scope, Key string
	Amount     int64
}

var ErrConflict = errors.New("postgres: payload conflict")

// Add commits business value, retained response, and outgoing event together.
// fault is a teaching failure hook before commit; nil is normal operation.
func (s Store) Add(ctx context.Context, e Event, fault func() error) (int64, error) {
	if e.Scope == "" || e.Key == "" || e.Amount <= 0 {
		return 0, errors.New("invalid request")
	}
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, "SET LOCAL lock_timeout='3s'; SET LOCAL statement_timeout='5s'"); err != nil {
		return 0, err
	}
	var current int64
	if err = tx.QueryRowContext(ctx, "SELECT value FROM business WHERE id=1 FOR UPDATE").Scan(&current); err != nil {
		return 0, err
	}
	var amount, response int64
	err = tx.QueryRowContext(ctx, "SELECT amount,response FROM requests WHERE scope=$1 AND key=$2", e.Scope, e.Key).Scan(&amount, &response)
	if err == nil {
		if amount != e.Amount {
			return 0, ErrConflict
		}
		return response, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	if current > math.MaxInt64-e.Amount {
		return 0, errors.New("overflow")
	}
	response = current + e.Amount
	if _, err = tx.ExecContext(ctx, "UPDATE business SET value=$1 WHERE id=1", response); err != nil {
		return 0, err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO requests VALUES($1,$2,$3,$4)", e.Scope, e.Key, e.Amount, response); err != nil {
		return 0, err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO outgoing(scope,key,amount) VALUES($1,$2,$3)", e.Scope, e.Key, e.Amount); err != nil {
		return 0, err
	}
	if fault != nil {
		if err = fault(); err != nil {
			return 0, err
		}
	}
	return response, tx.Commit()
}
func (s Store) Consume(ctx context.Context, e Event) (bool, error) {
	if e.Scope == "" || e.Key == "" || e.Amount <= 0 {
		return false, errors.New("invalid event")
	}
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, "SET LOCAL lock_timeout='3s'; SET LOCAL statement_timeout='5s'"); err != nil {
		return false, err
	}
	var amount int64
	err = tx.QueryRowContext(ctx, "INSERT INTO inbox VALUES($1,$2,$3) ON CONFLICT DO NOTHING RETURNING amount", e.Scope, e.Key, e.Amount).Scan(&amount)
	if errors.Is(err, sql.ErrNoRows) {
		if err = tx.QueryRowContext(ctx, "SELECT amount FROM inbox WHERE scope=$1 AND key=$2", e.Scope, e.Key).Scan(&amount); err != nil {
			return false, err
		}
		if amount != e.Amount {
			return false, ErrConflict
		}
		return false, tx.Commit()
	}
	if err != nil {
		return false, err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO effects VALUES($1,$2,$3)", e.Scope, e.Key, e.Amount); err != nil {
		return false, err
	}
	return true, tx.Commit()
}

// Relay attempts one bounded batch, acknowledging only successful deliveries.
// Concurrent relays can duplicate events; the destination inbox handles them.
func (s Store) Relay(ctx context.Context, limit int, deliver func(Event) error) error {
	if limit <= 0 || deliver == nil {
		return errors.New("invalid relay")
	}
	rows, err := s.DB.QueryContext(ctx, "SELECT scope,key,amount FROM outgoing WHERE NOT sent ORDER BY scope,key LIMIT $1", limit)
	if err != nil {
		return err
	}
	events := []Event{}
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.Scope, &e.Key, &e.Amount); err != nil {
			rows.Close()
			return err
		}
		events = append(events, e)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	for _, e := range events {
		if err := deliver(e); err != nil {
			return err
		}
		if _, err := s.DB.ExecContext(ctx, "UPDATE outgoing SET sent=true WHERE scope=$1 AND key=$2", e.Scope, e.Key); err != nil {
			return err
		}
	}
	return nil
}
