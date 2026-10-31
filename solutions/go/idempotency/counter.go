// Package idempotency commits a counter mutation and retained response together.
package idempotency

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"net/url"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var ErrConflict = errors.New("idempotency: key reused with different amount")

type Counter struct{ db *sql.DB }

func Open(filename string) (*Counter, error) {
	absolute, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	u := url.URL{Scheme: "file", Path: absolute, RawQuery: "_busy_timeout=5000&_txlock=immediate&_synchronous=FULL"}
	db, err := sql.Open("sqlite3", u.String())
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS business(id INTEGER PRIMARY KEY CHECK(id=1), value INTEGER NOT NULL);
INSERT OR IGNORE INTO business VALUES(1,0);
CREATE TABLE IF NOT EXISTS requests(scope TEXT NOT NULL, key TEXT NOT NULL, amount INTEGER NOT NULL, response INTEGER NOT NULL, PRIMARY KEY(scope,key));`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Counter{db}, nil
}
func (c *Counter) Close() error { return c.db.Close() }

// Add permits positive amounts only. The scoped key and exact integer amount
// identify one operation; stored responses have no automatic expiry.
func (c *Counter) Add(ctx context.Context, scope, key string, amount int64) (int64, error) {
	if scope == "" || key == "" || amount <= 0 {
		return 0, errors.New("idempotency: invalid input")
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var prior, response int64
	err = tx.QueryRowContext(ctx, "SELECT amount,response FROM requests WHERE scope=? AND key=?", scope, key).Scan(&prior, &response)
	if err == nil {
		if prior != amount {
			return 0, ErrConflict
		}
		return response, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	var current int64
	if err = tx.QueryRowContext(ctx, "SELECT value FROM business WHERE id=1").Scan(&current); err != nil {
		return 0, err
	}
	if current > math.MaxInt64-amount {
		return 0, errors.New("idempotency: overflow")
	}
	response = current + amount
	if _, err = tx.ExecContext(ctx, "UPDATE business SET value=? WHERE id=1", response); err != nil {
		return 0, err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO requests(scope,key,amount,response) VALUES(?,?,?,?)", scope, key, amount, response); err != nil {
		return 0, err
	}
	return response, tx.Commit()
}
