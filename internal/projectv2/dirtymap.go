package projectv2

import "maps"

type cloneable[T any] interface {
	Clone() T
}

type dirtyMapEntry[K comparable, V cloneable[V]] struct {
	m        *dirtyMap[K, V]
	key      K
	original V
	value    V
	dirty    bool
	delete   bool
}

func (e *dirtyMapEntry[K, V]) Change(apply func(V)) {
	if e.delete {
		panic("tried to change a deleted entry")
	}
	if !e.dirty {
		e.value = e.value.Clone()
		e.dirty = true
		e.m.dirty[e.key] = e
	}
	apply(e.value)
}

func (e *dirtyMapEntry[K, V]) Delete() {
	if !e.dirty {
		e.m.dirty[e.key] = e
	}
	e.delete = true
}

type dirtyMap[K comparable, V cloneable[V]] struct {
	base  map[K]V
	dirty map[K]*dirtyMapEntry[K, V]
}

func newDirtyMap[K comparable, V cloneable[V]](base map[K]V) *dirtyMap[K, V] {
	return &dirtyMap[K, V]{
		base:  base,
		dirty: make(map[K]*dirtyMapEntry[K, V]),
	}
}

func (m *dirtyMap[K, V]) Get(key K) (*dirtyMapEntry[K, V], bool) {
	if entry, ok := m.dirty[key]; ok {
		return entry, true
	}
	value, ok := m.base[key]
	if !ok {
		return nil, false
	}
	return &dirtyMapEntry[K, V]{
		m:        m,
		key:      key,
		original: value,
		value:    value,
		dirty:    false,
	}, true
}

// Add sets a new entry in the dirty map without checking if it exists
// in the base map. The entry added is considered dirty, so it should
// be a fresh value, mutable until finalized (i.e., it will not be cloned
// before changing if a change is made). If modifying an entry that may
// exist in the base map, use `Change` instead.
func (m *dirtyMap[K, V]) Add(key K, value V) {
	m.dirty[key] = &dirtyMapEntry[K, V]{
		m:     m,
		key:   key,
		value: value,
		dirty: true,
	}
}

// !!! Decide whether this, entry.Change(), or both should exist
func (m *dirtyMap[K, V]) Change(key K, apply func(V)) {
	if entry, ok := m.Get(key); ok {
		entry.Change(apply)
	} else {
		panic("tried to change a non-existent entry")
	}
}

func (m *dirtyMap[K, V]) Delete(key K) {
	if entry, ok := m.Get(key); ok {
		entry.Delete()
	} else {
		panic("tried to delete a non-existent entry")
	}
}

func (m *dirtyMap[K, V]) Range(fn func(*dirtyMapEntry[K, V]) bool) {
	seenInDirty := make(map[K]struct{})
	for _, entry := range m.dirty {
		seenInDirty[entry.key] = struct{}{}
		if !entry.delete && !fn(entry) {
			break
		}
	}
	for key, value := range m.base {
		if _, ok := seenInDirty[key]; ok {
			continue // already processed in dirty entries
		}
		if !fn(&dirtyMapEntry[K, V]{m: m, key: key, original: value, value: value, dirty: false}) {
			break
		}
	}
}

func (m *dirtyMap[K, V]) Finalize() (result map[K]V, changed bool) {
	if len(m.dirty) == 0 {
		return m.base, false // no changes, return base map
	}
	if m.base == nil {
		result = make(map[K]V, len(m.dirty))
	} else {
		result = maps.Clone(m.base)
	}
	for key, entry := range m.dirty {
		if entry.delete {
			delete(result, key)
		} else {
			result[key] = entry.value
		}
	}
	return result, true
}
