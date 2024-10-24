package collections

import (
	"iter"
	"maps"
	"slices"
)

// InsertionOrderedMap is an insertion ordered map.
type InsertionOrderedMap[K comparable, V any] struct {
	keys []K
	mp   map[K]V
}

// NewMapWithSizeHint creates a new Map with a hint for the number of elements it will contain.
func NewMapWithSizeHint[K comparable, V any](hint int) *InsertionOrderedMap[K, V] {
	return &InsertionOrderedMap[K, V]{
		keys: make([]K, 0, hint),
		mp:   make(map[K]V, hint),
	}
}

// Set sets a key-value pair in the map.
func (m *InsertionOrderedMap[K, V]) Set(key K, value V) {
	if m.mp == nil {
		m.mp = make(map[K]V)
	}

	if _, ok := m.mp[key]; !ok {
		m.keys = append(m.keys, key)
	}
	m.mp[key] = value
}

// Get retrieves a value from the map.
func (m *InsertionOrderedMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.mp[key]
	return v, ok
}

// Has returns true if the map contains the key.
func (m *InsertionOrderedMap[K, V]) Has(key K) bool {
	_, ok := m.mp[key]
	return ok
}

// Delete removes a key-value pair from the map.
func (m *InsertionOrderedMap[K, V]) Delete(key K) (V, bool) {
	v, ok := m.mp[key]
	if !ok {
		var zero V
		return zero, false
	}

	delete(m.mp, key)
	i := slices.Index(m.keys, key)
	// If we're just removing the first or last element, avoid shifting everything around.
	if i == 0 {
		var zero K
		m.keys[0] = zero
		m.keys = m.keys[1:]
	} else if end := len(m.keys) - 1; i == end {
		var zero K
		m.keys[end] = zero
		m.keys = m.keys[:end]
	} else {
		m.keys = slices.Delete(m.keys, i, i+1)
	}

	return v, true
}

// Keys returns an iterator over the keys in the map.
// A slice of the keys can be obtained by calling `slices.Collect`.
func (m *InsertionOrderedMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, key := range m.keys {
			if !yield(key) {
				break
			}
		}
	}
}

// Values returns an iterator over the values in the map.
// A slice of the values can be obtained by calling `slices.Collect`.
func (m *InsertionOrderedMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, key := range m.keys {
			if !yield(m.mp[key]) {
				break
			}
		}
	}
}

// Entries returns an iterator over the key-value pairs in the map.
func (m *InsertionOrderedMap[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, key := range m.keys {
			if !yield(key, m.mp[key]) {
				break
			}
		}
	}
}

// Clear removes all key-value pairs from the map.
// The space allocated for the map will be reused.
func (m *InsertionOrderedMap[K, V]) Clear() {
	clear(m.keys)
	m.keys = m.keys[:0]
	clear(m.mp)
}

// Size returns the number of key-value pairs in the map.
func (m *InsertionOrderedMap[K, V]) Size() int {
	return len(m.keys)
}

// Clone returns a copy of the map. Entries are shallow-copied.
func (m *InsertionOrderedMap[K, V]) Clone() *InsertionOrderedMap[K, V] {
	return &InsertionOrderedMap[K, V]{
		keys: slices.Clone(m.keys),
		mp:   maps.Clone(m.mp),
	}
}
