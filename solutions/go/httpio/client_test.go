package httpio

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestChunkedBodiesLimitsAndCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/blocked" {
			<-r.Context().Done()
			return
		}
		for _, chunk := range []string{"1\n", "2\n"} {
			w.Write([]byte(chunk))
			w.(http.Flusher).Flush()
		}
	}))
	defer server.Close()
	client := server.Client()
	body, err := Fetch(context.Background(), client, server.URL, 4)
	if err != nil || string(body) != "1\n2\n" {
		t.Fatal(string(body), err)
	}
	if _, err := Fetch(context.Background(), client, server.URL, 3); err == nil {
		t.Fatal("size")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if _, err := Fetch(ctx, client, server.URL+"/blocked", 4); err == nil {
		t.Fatal("deadline")
	}
}
