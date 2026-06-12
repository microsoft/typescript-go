package ast

import (
	"github.com/microsoft/typescript-go/internal/core"
)

type entry[T any] struct {
	key  uint64
	data *T
}

type intHashTable[T any] struct {
	buckets []entry[T] // Hash table entries (length is power of 2)
	size    int        // Number of entries in hash table
}

// MurmurHash3-style bit mixer for balanced index distribution
func hashInt(x uint64) uint64 {
	x ^= x >> 33
	x *= 0xff51afd7ed558ccd
	x ^= x >> 33
	x *= 0xc4ceb9fe1a85ec53
	x ^= x >> 33
	return x
}

// Assumes the given key doesn't already exist in the hash table
func (h *intHashTable[T]) insert(key uint64, data *T) {
	if data == nil {
		panic("IntHashTable cannot have nil entries")
	}
	if h.size == 0 {
		// Initial length must be a power of 2
		h.buckets = make([]entry[T], 256)
	} else if h.size*4 >= len(h.buckets)*3 {
		// Resize when 0.75 load factor is crossed
		h.resize()
	}
	hash := hashInt(key)
	mask := uint64(len(h.buckets) - 1)
	i := hash & mask
	for h.buckets[i].data != nil {
		i = (i + 1) & mask
	}
	h.buckets[i] = entry[T]{key, data}
	h.size++
}

func (h *intHashTable[T]) get(key uint64) *T {
	if h.size == 0 {
		return nil
	}
	hash := hashInt(key)
	mask := uint64(len(h.buckets) - 1)
	i := hash & mask
	for {
		data := h.buckets[i].data
		if data == nil || h.buckets[i].key == key {
			return data
		}
		i = (i + 1) & mask
	}
}

func (h *intHashTable[T]) resize() {
	newBuckets := make([]entry[T], len(h.buckets)*2)
	mask := uint64(len(newBuckets) - 1)
	for i := range h.buckets {
		if h.buckets[i].data != nil {
			hash := hashInt(h.buckets[i].key)
			j := hash & mask
			for newBuckets[j].data != nil {
				j = (j + 1) & mask
			}
			newBuckets[j] = h.buckets[i]
		}
	}
	h.buckets = newBuckets
}

type NodeLinkStore[T any] struct {
	entries intHashTable[T]
	arena   core.Arena[T]
}

func (s *NodeLinkStore[T]) Get(node *Node) *T {
	key := uint64(GetNodeId(node))
	value := s.entries.get(key)
	if value == nil {
		value = s.arena.New()
		s.entries.insert(key, value)
	}
	return value
}

func (s *NodeLinkStore[T]) Has(node *Node) bool {
	return s.entries.get(uint64(GetNodeId(node))) != nil
}

func (s *NodeLinkStore[T]) TryGet(node *Node) *T {
	return s.entries.get(uint64(GetNodeId(node)))
}

type SymbolLinkStore[T any] struct {
	entries intHashTable[T]
	arena   core.Arena[T]
}

func (s *SymbolLinkStore[T]) Get(symbol *Symbol) *T {
	key := uint64(GetSymbolId(symbol))
	value := s.entries.get(key)
	if value == nil {
		value = s.arena.New()
		s.entries.insert(key, value)
	}
	return value
}

func (s *SymbolLinkStore[T]) Has(symbol *Symbol) bool {
	return s.entries.get(uint64(GetSymbolId(symbol))) != nil
}

func (s *SymbolLinkStore[T]) TryGet(symbol *Symbol) *T {
	return s.entries.get(uint64(GetSymbolId(symbol)))
}
