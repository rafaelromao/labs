// Package ratelimiter provides a bounded-cardinality token bucket per client.
package ratelimiter

import (
	"errors"
	"time"
)

var ErrClients = errors.New("ratelimiter: client limit reached")

type Limiter struct{}

// New configures tokens/second, burst, maximum retained clients and idle expiry.
// The clock must not go backwards. Recreated clients start with a full burst.
func New(rate float64, burst, maxClients int, idle time.Duration, now func() time.Time) (*Limiter, error) {
	panic("TODO: implement")
}

// Allow consumes one token, or returns false when that client has no token.
// A distinct error signals that a new client could not be retained.
func (l *Limiter) Allow(client string) (bool, error) {
	panic("TODO: implement")
}
