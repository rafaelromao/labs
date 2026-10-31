// Package fencing illustrates resource-enforced owner tokens and revisions.
package fencing

import (
	"math"
	"sync"
)

type Resource struct {
	mu              sync.Mutex
	token, revision int64
	value           string
}

func (r *Resource) Write(token, expected int64, value string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if token < r.token || token < 0 || expected != r.revision || r.revision == math.MaxInt64 {
		return false
	}
	r.token = token
	r.revision++
	r.value = value
	return true
}
func (r *Resource) Read() (string, int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.value, r.revision
}
