package store

import (
	"sync"
	"testing"
)

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
