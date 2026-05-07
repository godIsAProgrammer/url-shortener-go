package shortener

import (
	"crypto/rand"
	"errors"
	"sync"
)

const (
	codeLen  = 8
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// Store keeps short codes in a process-local map; the data is lost when the
// container restarts so reviewers can always start from a clean slate.
type Store struct {
	mu    sync.RWMutex
	links map[string]string
}

func NewStore() *Store {
	return &Store{links: make(map[string]string)}
}

// Save generates a unique code for url and persists the mapping. It retries
// a few times on the unlikely event of a code collision.
func (s *Store) Save(url string) (string, error) {
	for i := 0; i < 5; i++ {
		code, err := newCode(codeLen)
		if err != nil {
			return "", err
		}
		s.mu.Lock()
		if _, exists := s.links[code]; !exists {
			s.links[code] = url
			s.mu.Unlock()
			return code, nil
		}
		s.mu.Unlock()
	}
	return "", errors.New("could not generate unique code after retries")
}

// Get resolves a code back to its original URL.
func (s *Store) Get(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.links[code]
	return url, ok
}

func newCode(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i := range buf {
		buf[i] = alphabet[int(buf[i])%len(alphabet)]
	}
	return string(buf), nil
}
