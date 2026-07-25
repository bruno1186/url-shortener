package shortener

import (
	"errors"
	"sync"
)

// ErrNotFound e retornado quando um codigo curto nao existe.
var ErrNotFound = errors.New("shortener: code not found")

// Store define a persistencia de mapeamentos codigo -> URL.
// A interface permite trocar a implementacao (memoria, Redis, SQL)
// sem alterar a logica de negocio.
type Store interface {
	Save(code, url string) error
	Load(code string) (string, error)
}

// MemoryStore e uma implementacao thread-safe em memoria.
type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewMemoryStore cria um MemoryStore vazio.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]string)}
}

func (s *MemoryStore) Save(code, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[code] = url
	return nil
}

func (s *MemoryStore) Load(code string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.data[code]
	if !ok {
		return "", ErrNotFound
	}
	return url, nil
}
