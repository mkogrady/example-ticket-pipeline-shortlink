package store

import (
	"crypto/rand"
	"sync"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// DefaultTTL is how long a link lives when no TTL is configured.
const DefaultTTL = 24 * time.Hour

const sweepInterval = time.Minute

type entry struct {
	target  string
	expires time.Time
}

type Store struct {
	mu        sync.Mutex
	links     map[string]entry
	now       func() time.Time
	ttl       time.Duration
	lastSweep time.Time
}

// Option configures a Store.
type Option func(*Store)

// WithTTL sets how long a link lives before it expires.
func WithTTL(ttl time.Duration) Option {
	return func(s *Store) { s.ttl = ttl }
}

func New(opts ...Option) *Store {
	s := &Store{links: map[string]entry{}, now: time.Now, ttl: DefaultTTL}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Store) Put(target string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	s.sweep(now)
	for {
		code := newCode()
		if _, taken := s.links[code]; taken {
			continue
		}
		s.links[code] = entry{target: target, expires: now.Add(s.ttl)}
		return code
	}
}

func (s *Store) Get(code string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.links[code]
	if !ok {
		return "", false
	}
	if !s.now().Before(e.expires) {
		delete(s.links, code)
		return "", false
	}
	return e.target, true
}

// sweep deletes expired entries, at most once per sweepInterval.
// The caller must hold s.mu.
func (s *Store) sweep(now time.Time) {
	if now.Sub(s.lastSweep) < sweepInterval {
		return
	}
	s.lastSweep = now
	for code, e := range s.links {
		if !now.Before(e.expires) {
			delete(s.links, code)
		}
	}
}

func newCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b)
}
