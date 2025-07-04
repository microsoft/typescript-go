package collections

type ManyToManyMap[K comparable, V comparable] struct {
	keyToValueSet map[K]*Set[V]
	valueToKeySet map[V]*Set[K]
}

func (m *ManyToManyMap[K, V]) GetKeys(value V) (*Set[K], bool) {
	keys, present := m.valueToKeySet[value]
	return keys, present
}

func (m *ManyToManyMap[K, V]) GetValues(key K) (*Set[V], bool) {
	values, present := m.keyToValueSet[key]
	return values, present
}

func (m *ManyToManyMap[K, V]) Len() int {
	return len(m.keyToValueSet)
}

func (m *ManyToManyMap[K, V]) Keys() map[K]*Set[V] {
	return m.keyToValueSet
}

func (m *ManyToManyMap[K, V]) Add(key K, valueSet *Set[V]) {
	existingValueSet, hasExisting := m.keyToValueSet[key]
	if m.keyToValueSet == nil {
		m.keyToValueSet = make(map[K]*Set[V])
	}
	m.keyToValueSet[key] = valueSet
	for value := range valueSet.Keys() {
		if !hasExisting || !existingValueSet.Has(value) {
			// Add to valueToKeySet
			if m.valueToKeySet == nil {
				m.valueToKeySet = make(map[V]*Set[K])
			}
			addToMapOfSet(m.valueToKeySet, value, key)
		}
	}

	if hasExisting {
		for value := range existingValueSet.Keys() {
			if !valueSet.Has(value) {
				// Remove from valueToKeySet
				defer deleteFromMapOfSet(m.valueToKeySet, value, key)
			}
		}
	}
}

func addToMapOfSet[K comparable, V comparable](mapKSetV map[K]*Set[V], key K, value V) {
	set, exists := mapKSetV[key]
	if !exists {
		set = &Set[V]{}
		mapKSetV[key] = set
	}
	set.Add(value)
}

func deleteFromMapOfSet[K comparable, V comparable](mapKSetV map[K]*Set[V], key K, value V) {
	if set, exists := mapKSetV[key]; exists {
		set.Delete(value)
		if set.Len() == 0 {
			delete(mapKSetV, key)
		}
	}
}
