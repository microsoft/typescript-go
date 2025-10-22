package dirty

import "maps"

type SetEntry[K comparable] struct {
	m *Set[K]
	mapEntry[K, struct{}]
}

func (e *SetEntry[K]) Delete() {
	if !e.dirty {
		e.m.dirty[e.key] = e
	}
	e.delete = true
}

type Set[K comparable] struct {
	base  map[K]struct{}
	dirty map[K]*SetEntry[K]
}

func NewSet[K comparable](base map[K]struct{}) *Set[K] {
	return &Set[K]{
		base:  base,
		dirty: make(map[K]*SetEntry[K]),
	}
}

func (m *Set[K]) Get(key K) (*SetEntry[K], bool) {
	if entry, ok := m.dirty[key]; ok {
		if entry.delete {
			return nil, false
		}
		return entry, true
	}
	value, ok := m.base[key]
	if !ok {
		return nil, false
	}
	return &SetEntry[K]{
		m: m,
		mapEntry: mapEntry[K, struct{}]{
			key:      key,
			original: value,
			value:    value,
			dirty:    false,
		},
	}, true
}

func (m *Set[K]) Has(key K) bool {
	if entry, ok := m.dirty[key]; ok {
		if entry.delete {
			return false
		}
		return true
	}
	_, ok := m.base[key]
	return ok
}

// Add sets a new entry in the dirty map without checking if it exists
// in the base map. The entry added is considered dirty, so it should
// be a fresh value, mutable until finalized (i.e., it will not be cloned
// before changing if a change is made). If modifying an entry that may
// exist in the base map, use `Change` instead.
func (m *Set[K]) Add(key K) {
	m.dirty[key] = &SetEntry[K]{
		m: m,
		mapEntry: mapEntry[K, struct{}]{
			key:   key,
			value: struct{}{},
			dirty: true,
		},
	}
}

func (m *Set[K]) Delete(key K) {
	if entry, ok := m.Get(key); ok {
		entry.Delete()
	} else {
		panic("tried to delete a non-existent entry")
	}
}

func (m *Set[K]) Range(fn func(*SetEntry[K]) bool) {
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
		if !fn(&SetEntry[K]{m: m, mapEntry: mapEntry[K, struct{}]{
			key:      key,
			original: value,
			value:    value,
			dirty:    false,
		}}) {
			break
		}
	}
}

func (m *Set[K]) Finalize() (result map[K]struct{}, changed bool) {
	if len(m.dirty) == 0 {
		return m.base, false // no changes, return base map
	}
	if m.base == nil {
		result = make(map[K]struct{}, len(m.dirty))
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
