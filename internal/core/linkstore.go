package core

import "slices"

// Links store

type LinkStore[K comparable, V any] struct {
	entries map[K]*V
	arena   Arena[V]
}

func (s *LinkStore[K, V]) Get(key K) *V {
	value := s.entries[key]
	if value != nil {
		return value
	}
	if s.entries == nil {
		s.entries = make(map[K]*V)
	}
	value = s.arena.New()
	s.entries[key] = value
	return value
}

func (s *LinkStore[K, V]) Has(key K) bool {
	_, ok := s.entries[key]
	return ok
}

func (s *LinkStore[K, V]) TryGet(key K) *V {
	return s.entries[key]
}

const (
	pageShift = 8
	pageSize  = 1 << pageShift
	pageMask  = pageSize - 1
)

// Implements a sparse-array-like structure for storing elements keyed by dense uint64 keys. Elements are
// stored in fixed-size pages of 256 entries and an index of pages is maintained in either an array or a map.
type PagedLinkStore[V any] struct {
	pageMap  map[int]*[pageSize]V // Page table as a map
	pageList []*[pageSize]V       // Page table as an array
}

// Initialize the link store. If useArrayPageTable is true, the store will use an array for the page table.
// Otherwise, the store defaults to using a map for the page table. An array-based page table is very fast,
// but requires key values to range densely from 0 to the maximum possible value. A map-based page table
// places no restrictions on the key values.
func (s *PagedLinkStore[V]) Initialize(useArrayPageTable bool) {
	s.pageMap = nil
	s.pageList = nil
	if useArrayPageTable {
		s.pageList = []*[pageSize]V{}
	}
}

func (s *PagedLinkStore[V]) Get(key uint64) *V {
	var page *[pageSize]V
	pageIndex := int(key >> pageShift)
	if s.pageList != nil {
		if pageIndex >= len(s.pageList) {
			// Grow the length of the list to pageIndex+1
			s.pageList = slices.Grow(s.pageList, pageIndex-len(s.pageList)+1)[:pageIndex+1]
		}
		page = s.pageList[pageIndex]
		if page == nil {
			page = new([pageSize]V)
			s.pageList[pageIndex] = page
		}
	} else {
		page = s.pageMap[pageIndex]
		if page == nil {
			page = new([pageSize]V)
			if s.pageMap == nil {
				s.pageMap = make(map[int]*[pageSize]V)
			}
			s.pageMap[pageIndex] = page
		}
	}
	return &page[key&pageMask]
}

func (s *PagedLinkStore[V]) Has(key uint64) bool {
	return s.TryGet(key) != nil
}

func (s *PagedLinkStore[V]) TryGet(key uint64) *V {
	var page *[pageSize]V
	pageIndex := int(key >> pageShift)
	if s.pageList != nil {
		if pageIndex < len(s.pageList) {
			page = s.pageList[pageIndex]
		}
	} else {
		page = s.pageMap[pageIndex]
	}
	if page != nil {
		return &page[key&pageMask]
	}
	return nil
}
