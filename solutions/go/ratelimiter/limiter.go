// Package ratelimiter provides a bounded-cardinality token bucket per client.
package ratelimiter

import (
	"errors"
	"math"
	"sync"
	"time"
)

var ErrClients = errors.New("ratelimiter: client limit reached")

type bucket struct {
	tokens float64
	last   time.Time
}

type Limiter struct {
	mu         sync.Mutex
	burst      float64
	rate       float64
	maxClients int
	idle       time.Duration
	now        func() time.Time
	clients    map[string]*bucket
}

// New configures tokens/second, burst, maximum retained clients and idle expiry.
// The clock must not go backwards. Recreated clients start with a full burst.
func New(rate float64, burst, maxClients int, idle time.Duration, now func() time.Time) (*Limiter, error) {
	if rate <= 0 || math.IsNaN(rate) || math.IsInf(rate, 0) || burst <= 0 || maxClients <= 0 || idle <= 0 || now == nil {
		return nil, errors.New("ratelimiter: invalid configuration")
	}
	return &Limiter{rate: rate, burst: float64(burst), maxClients: maxClients, idle: idle, now: now, clients: make(map[string]*bucket)}, nil
}

// Allow consumes one token, or returns false when that client has no token.
// A distinct error signals that a new client could not be retained.
func (l *Limiter) Allow(client string) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	for key, b := range l.clients {
		if now.Sub(b.last) >= l.idle {
			delete(l.clients, key)
		}
	}
	b := l.clients[client]
	if b == nil {
		if len(l.clients) == l.maxClients {
			return false, ErrClients
		}
		b = &bucket{tokens: l.burst, last: now}
		l.clients[client] = b
	}
	elapsed := now.Sub(b.last).Seconds()
	if elapsed < 0 {
		return false, errors.New("ratelimiter: clock moved backwards")
	}
	b.tokens = math.Min(l.burst, b.tokens+elapsed*l.rate)
	b.last = now
	if b.tokens < 1 {
		return false, nil
	}
	b.tokens--
	return true, nil
}
