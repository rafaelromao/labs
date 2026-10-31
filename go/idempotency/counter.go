// Package idempotency commits a counter mutation and retained response together.
package idempotency

import (
	"context"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

var ErrConflict = errors.New("idempotency: key reused with different amount")

type Counter struct{}

func Open(filename string) (*Counter, error) {
	panic("TODO: implement")
}

func (c *Counter) Close() error {
	panic("TODO: implement")
}

// Add permits positive amounts only. The scoped key and exact integer amount
// identify one operation; stored responses have no automatic expiry.
func (c *Counter) Add(ctx context.Context, scope, key string, amount int64) (int64, error) {
	panic("TODO: implement")
}
