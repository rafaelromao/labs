// Package streamio parses bounded UTF-8 decimal records without buffering a file.
package streamio

import (
	"context"
	"io"
)

type Result struct {
	Count int64
	Sum   int64
}

// Process reads LF-delimited signed decimal integers. CRLF is also accepted.
// maxBytes excludes LF but includes CR. A final unterminated record is accepted.
// It leaves input open and returns no partial result on error. Cancellation cannot
// interrupt an arbitrary blocked Reader; the owner must close or deadline it.
func Process(ctx context.Context, input io.Reader, maxBytes int) (Result, error) {
	panic("TODO: implement")
}
