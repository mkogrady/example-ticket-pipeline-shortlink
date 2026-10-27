package store

import (
	"sync"
	"testing"
	"time"
)

type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time { return c.t }

func newFake(ttl time.Duration) (*Store, *fakeClock) {
	c := &fakeClock{t: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	s := NewWithTTL(ttl)
	s.now = c.now
	return s, c
}

func TestNewWithTTL(t *testing.T) {
	if New().ttl != DefaultTTL {
		t.Fatal("New should use DefaultTTL")
	}
	if NewWithTTL(time.Second).ttl != time.Second {
		t.Fatal("ttl not set")
	}
}

func TestGetBeforeExpiry(t *testing.T) {
	s, c := newFake(time.Hour)
	code := s.Put("https://example.com")
	c.t = c.t.Add(time.Hour - time.Nanosecond)
	if got, ok := s.Get(code); !ok || got != "https://example.com" {
		t.Fatalf("got %q %v", got, ok)
	}
}

func TestGetAfterExpiryMissesAndDeletes(t *testing.T) {
	s, c := newFake(time.Hour)
	code := s.Put("https://example.com")
	c.t = c.t.Add(time.Hour)
	if _, ok := s.Get(code); ok {
		t.Fatal("expected miss at expiry")
	}
	if _, present := s.links[code]; present {
		t.Fatal("expired entry not deleted")
	}
}

func TestPutSweepsExpired(t *testing.T) {
	s, c := newFake(time.Hour)
	old := s.Put("https://old.example")
	c.t = c.t.Add(2 * time.Hour)
	s.Put("https://new.example")
	if _, present := s.links[old]; present {
		t.Fatal("expired entry survived sweep")
	}
}

func TestPutSweepThrottled(t *testing.T) {
	s, c := newFake(time.Second)
	s.Put("https://a.example") // sweeps at t0
	old := s.Put("https://old.example")
	c.t = c.t.Add(10 * time.Second) // expired, but under a minute since sweep
	s.Put("https://new.example")
	if _, present := s.links[old]; !present {
		t.Fatal("sweep ran more than once per minute")
	}
}

func TestPutGet(t *testing.T) {
	s := New()
	code := s.Put("https://example.com")
	if len(code) != 6 {
		t.Fatalf("code %q", code)
	}
	got, ok := s.Get(code)
	if !ok || got != "https://example.com" {
		t.Fatalf("got %q %v", got, ok)
	}
}

func TestGetUnknown(t *testing.T) {
	if _, ok := New().Get("nothere"); ok {
		t.Fatal("expected miss")
	}
}

func TestPutDistinctCodes(t *testing.T) {
	s := New()
	if s.Put("https://a.example") == s.Put("https://b.example") {
		t.Fatal("same code twice")
	}
}

func TestPutConcurrent(t *testing.T) {
	s := New()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Put("https://example.com")
		}()
	}
	wg.Wait()
}
