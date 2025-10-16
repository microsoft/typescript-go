package pnp

type Trie[T any] struct {
	root *node[T]
}

type node[T any] struct {
	children map[rune]*node[T]
	value    T
	hasValue bool
}

func New[T any]() *Trie[T] {
	return &Trie[T]{root: &node[T]{children: make(map[rune]*node[T])}}
}

func (t *Trie[T]) Set(key string, v T) {
	n := t.root
	for _, r := range key {
		child, ok := n.children[r]
		if !ok {
			child = &node[T]{children: make(map[rune]*node[T])}
			n.children[r] = child
		}
		n = child
	}
	n.value = v
	n.hasValue = true
}

func (t *Trie[T]) Get(key string) (T, bool) {
	n := t.root
	for _, r := range key {
		child, ok := n.children[r]
		if !ok {
			var zero T
			return zero, false
		}
		n = child
	}
	if n.hasValue {
		return n.value, true
	}
	var zero T
	return zero, false
}

func (t *Trie[T]) GetAncestorValue(key string) (T, bool) {
	n := t.root

	var best T
	var okBest bool
	if n.hasValue {
		best, okBest = n.value, true
	}

	for _, r := range key {
		child, ok := n.children[r]
		if !ok {
			var zero T
			if okBest {
				return best, true
			}
			return zero, false
		}
		n = child
		if n.hasValue {
			best, okBest = n.value, true
		}
	}

	if n.hasValue {
		return n.value, true
	}
	var zero T
	if okBest {
		return best, true
	}
	return zero, false
}
