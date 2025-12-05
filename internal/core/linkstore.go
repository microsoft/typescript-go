package core

import "sync"

type LinkStore[K comparable, V any] struct {
	mu      sync.RWMutex
	entries map[K]*V
	pool    Pool[V]
}

func (s *LinkStore[K, V]) Get(key K) *V {
	s.mu.RLock()
	value := s.entries[key]
	s.mu.RUnlock()
	if value != nil {
		return value
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if value := s.entries[key]; value != nil {
		return value
	}

	if s.entries == nil {
		s.entries = make(map[K]*V)
	}

	value = s.pool.New()
	s.entries[key] = value
	return value
}

func (s *LinkStore[K, V]) Has(key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.entries[key]
	return ok
}

func (s *LinkStore[K, V]) TryGet(key K) *V {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.entries[key]
}
