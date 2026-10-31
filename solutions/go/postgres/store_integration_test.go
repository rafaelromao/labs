//go:build integration

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCrashHelper(t *testing.T) {
	dsn := os.Getenv("DEPTH_CRASH_DSN")
	if dsn == "" {
		return
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		os.Exit(90)
	}
	s := Store{db}
	phase := os.Getenv("DEPTH_CRASH_PHASE")
	var hook func() error
	if phase == "before" {
		hook = func() error { os.Exit(17); return nil }
	}
	if _, err := s.Add(context.Background(), Event{"tenant", phase, 3}, hook); err != nil {
		os.Exit(91)
	}
	os.Exit(18)
}
func TestPostgresRecoveryAndPostGIS(t *testing.T) {
	dsn := os.Getenv("DEPTH_POSTGRES_DSN")
	if dsn == "" {
		t.Fatal("integration requires DEPTH_POSTGRES_DSN pointing at a dedicated PostgreSQL/PostGIS database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	admin, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := fmt.Sprintf("depth_go_%d", time.Now().UnixNano())
	if _, err := admin.ExecContext(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP SCHEMA " + schema + " CASCADE")
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Set("search_path", schema+",public")
	q.Set("connect_timeout", "5")
	u.RawQuery = q.Encode()
	db, err := sql.Open("postgres", u.String())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(8)
	migration, err := os.ReadFile("../../../fixtures/postgres.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, string(migration)); err != nil {
		t.Fatal(err)
	}
	s := Store{db}
	e := Event{"tenant", "duplicate", 7}
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			value, err := s.Add(ctx, e, nil)
			if err != nil || value != 7 {
				t.Errorf("duplicate %d %v", value, err)
			}
		}()
	}
	wg.Wait()
	if _, err := s.Add(ctx, Event{e.Scope, e.Key, 8}, nil); !errors.Is(err, ErrConflict) {
		t.Fatal(err)
	}
	if _, err := s.Add(ctx, Event{"tenant", "rollback", 4}, func() error { return errors.New("injected") }); err == nil {
		t.Fatal("fault not observed")
	}
	var value int64
	if err := db.QueryRowContext(ctx, "SELECT value FROM business WHERE id=1").Scan(&value); err != nil || value != 7 {
		t.Fatal(value, err)
	}
	for _, phase := range []string{"before", "after"} {
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestCrashHelper$")
		cmd.Env = append(os.Environ(), "DEPTH_CRASH_DSN="+u.String(), "DEPTH_CRASH_PHASE="+phase)
		err := cmd.Run()
		var exit *exec.ExitError
		if !errors.As(err, &exit) || (exit.ExitCode() != 17 && exit.ExitCode() != 18) {
			t.Fatal(phase, err)
		}
		if _, err := s.Add(ctx, Event{"tenant", phase, 3}, nil); err != nil {
			t.Fatal(err)
		}
	}
	if err := db.QueryRowContext(ctx, "SELECT value FROM business WHERE id=1").Scan(&value); err != nil || value != 13 {
		t.Fatal("recovery", value, err)
	}
	// Relay loss after consumer commit but before ack must replay without a second effect.
	if err := s.Relay(ctx, 10, func(e Event) error {
		_, err := s.Consume(ctx, e)
		if err != nil {
			return err
		}
		return errors.New("lost ack")
	}); err == nil {
		t.Fatal("injection")
	}
	if err := s.Relay(ctx, 10, func(e Event) error { _, err := s.Consume(ctx, e); return err }); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM effects").Scan(&count); err != nil || count != 3 {
		t.Fatal(count, err)
	}
	geo, err := os.ReadFile("../../../fixtures/postgis.sql")
	if err != nil {
		t.Fatal(err)
	}
	var correct bool
	if err := db.QueryRowContext(ctx, strings.TrimSpace(string(geo))).Scan(&correct); err != nil || !correct {
		t.Fatal("PostGIS", correct, err)
	}
}
