// Package httpio owns bounded HTTP response bodies; callers own and reuse Client.
package httpio

import (
	"context"
	"errors"
	"io"
	"net/http"
)

func Fetch(ctx context.Context, client *http.Client, url string, maxBytes int64) ([]byte, error) {
	if client == nil || maxBytes <= 0 || maxBytes > 1<<20 {
		return nil, errors.New("invalid client or limit")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, errors.New("HTTP status rejected")
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > maxBytes {
		return nil, errors.New("body limit")
	}
	return body, nil
}
