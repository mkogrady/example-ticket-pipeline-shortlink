package store

import (
	"crypto/rand"
	"sync"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Store struct {
	mu    sync.Mutex
	links map[string]string
}

func New() *Store {
	return &Store{links: map[string]string{}}
}

func (s *Store) Put(target string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	for {
		code := newCode()
		if _, taken := s.links[code]; taken {
			continue
		}
		s.links[code] = target
		return code
	}
}

func (s *Store) Get(code string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	target, ok := s.links[code]
	return target, ok
}

func newCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b)
}
