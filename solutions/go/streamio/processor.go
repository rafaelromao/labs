// Package streamio parses bounded UTF-8 decimal records without buffering a file.
package streamio

import (
	"bufio"
	"context"
	"errors"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
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
	if maxBytes <= 0 {
		return Result{}, errors.New("streamio: positive record limit required")
	}
	r := bufio.NewReader(input)
	line := make([]byte, 0, min(maxBytes, 4096))
	result := Result{}
	consume := func() error {
		if !utf8.Valid(line) {
			return errors.New("streamio: invalid UTF-8")
		}
		text := strings.TrimSuffix(string(line), "\r")
		if text == "" {
			return errors.New("streamio: empty record")
		}
		start := 0
		if text[0] == '+' || text[0] == '-' {
			start = 1
		}
		if start == len(text) {
			return errors.New("streamio: missing digits")
		}
		for _, b := range []byte(text[start:]) {
			if b < '0' || b > '9' {
				return errors.New("streamio: decimal digits required")
			}
		}
		value, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return err
		}
		if (value > 0 && result.Sum > math.MaxInt64-value) || (value < 0 && result.Sum < math.MinInt64-value) {
			return errors.New("streamio: sum overflow")
		}
		if result.Count == math.MaxInt64 {
			return errors.New("streamio: count overflow")
		}
		result.Count++
		result.Sum += value
		line = line[:0]
		return nil
	}
	for {
		if err := ctx.Err(); err != nil {
			return Result{}, err
		}
		b, err := r.ReadByte()
		if err == io.EOF {
			if len(line) > 0 {
				if err := consume(); err != nil {
					return Result{}, err
				}
			}
			return result, nil
		}
		if err != nil {
			return Result{}, err
		}
		if b == '\n' {
			if err := consume(); err != nil {
				return Result{}, err
			}
			continue
		}
		if len(line) == maxBytes {
			return Result{}, errors.New("streamio: record too large")
		}
		line = append(line, b)
	}
}
