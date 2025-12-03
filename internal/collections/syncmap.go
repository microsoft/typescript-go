package collections

import (
	"iter"

	"github.com/go4org/hashtriemap"
)

type SyncMap[K comparable, V any] struct {
	hashtriemap.HashTrieMap[K, V]
}

// Size returns the approximate number of items in the map.
// Note that this is not a precise count, as the map may be modified
// concurrently while this method is running.
func (s *SyncMap[K, V]) Size() int {
	count := 0
	for range s.All() {
		count++
	}
	return count
}

func (s *SyncMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range s.All() {
			if !yield(k) {
				return
			}
		}
	}
}

func (s *SyncMap[K, V]) Clone() *SyncMap[K, V] {
	clone := &SyncMap[K, V]{}
	for k, v := range s.All() {
		clone.Store(k, v)
	}
	return clone
}
