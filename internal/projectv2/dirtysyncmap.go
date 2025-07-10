package projectv2

import (
	"maps"
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
)

type dirtySyncMapEntry[K comparable, V cloneable[V]] struct {
	m        *dirtySyncMap[K, V]
	mu       sync.Mutex
	key      K
	original V
	value    V
	dirty    bool
	delete   bool
}

func (e *dirtySyncMapEntry[K, V]) Change(apply func(V)) *dirtySyncMapEntry[K, V] {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.changeLocked(apply)
}

func (e *dirtySyncMapEntry[K, V]) changeLocked(apply func(V)) *dirtySyncMapEntry[K, V] {
	if e.dirty {
		apply(e.value)
		return e
	}

	entry, loaded := e.m.dirty.LoadOrStore(e.key, e)
	if loaded {
		entry.mu.Lock()
		defer entry.mu.Unlock()
	}
	if !entry.dirty {
		entry.value = entry.value.Clone()
		entry.dirty = true
	}
	apply(entry.value)
	return entry
}

func (e *dirtySyncMapEntry[K, V]) ChangeIf(cond func(V) bool, apply func(V)) *dirtySyncMapEntry[K, V] {
	e.mu.Lock()
	defer e.mu.Unlock()
	if cond(e.value) {
		return e.changeLocked(apply)
	}
	return e
}

type dirtySyncMap[K comparable, V cloneable[V]] struct {
	base          map[K]V
	dirty         collections.SyncMap[K, *dirtySyncMapEntry[K, V]]
	finalizeValue func(dirty V, original V) V
}

func newDirtySyncMap[K comparable, V cloneable[V]](base map[K]V, finalizeValue func(dirty V, original V) V) *dirtySyncMap[K, V] {
	return &dirtySyncMap[K, V]{
		base:          base,
		dirty:         collections.SyncMap[K, *dirtySyncMapEntry[K, V]]{},
		finalizeValue: finalizeValue,
	}
}

func (m *dirtySyncMap[K, V]) Load(key K) (*dirtySyncMapEntry[K, V], bool) {
	if entry, ok := m.dirty.Load(key); ok {
		return entry, true
	}
	if val, ok := m.base[key]; ok {
		return &dirtySyncMapEntry[K, V]{
			m:        m,
			key:      key,
			original: val,
			value:    val,
			dirty:    false,
			delete:   false,
		}, true
	}
	return nil, false
}

func (m *dirtySyncMap[K, V]) LoadOrStore(key K, value V) (*dirtySyncMapEntry[K, V], bool) {
	// Check for existence in the base map first so the sync map access is atomic.
	if value, ok := m.base[key]; ok {
		if dirty, ok := m.dirty.Load(key); ok {
			return dirty, true
		}
		return &dirtySyncMapEntry[K, V]{
			m:        m,
			key:      key,
			original: value,
			value:    value,
			dirty:    false,
			delete:   false,
		}, true
	}
	entry, loaded := m.dirty.LoadOrStore(key, &dirtySyncMapEntry[K, V]{
		m:     m,
		key:   key,
		value: value,
		dirty: true,
	})
	return entry, loaded
}

func (m *dirtySyncMap[K, V]) Range(fn func(*dirtySyncMapEntry[K, V]) bool) {
	seenInDirty := make(map[K]struct{})
	m.dirty.Range(func(key K, entry *dirtySyncMapEntry[K, V]) bool {
		seenInDirty[key] = struct{}{}
		if !entry.delete && !fn(entry) {
			return false
		}
		return true
	})
	for key, value := range m.base {
		if _, ok := seenInDirty[key]; ok {
			continue // already processed in dirty entries
		}
		if !fn(&dirtySyncMapEntry[K, V]{m: m, key: key, original: value, value: value, dirty: false}) {
			break
		}
	}
}

func (m *dirtySyncMap[K, V]) Finalize() (map[K]V, bool) {
	var changed bool
	result := m.base
	ensureCloned := func() {
		if !changed {
			if m.base == nil {
				result = make(map[K]V)
			} else {
				result = maps.Clone(m.base)
			}
			changed = true
		}
	}

	m.dirty.Range(func(key K, entry *dirtySyncMapEntry[K, V]) bool {
		if entry.delete {
			ensureCloned()
			delete(result, key)
		} else if entry.dirty {
			ensureCloned()
			if m.finalizeValue != nil {
				result[key] = m.finalizeValue(entry.value, entry.original)
			} else {
				result[key] = entry.value
			}
		}
		return true
	})
	return result, changed
}

func cloneMapIfNil[K comparable, V any, T any](dirty *T, original *T, getMap func(*T) map[K]V) map[K]V {
	dirtyMap := getMap(dirty)
	if dirtyMap == nil {
		if original == nil {
			return make(map[K]V)
		}
		originalMap := getMap(original)
		if originalMap == nil {
			return make(map[K]V)
		}
		return maps.Clone(originalMap)
	}
	return dirtyMap
}
