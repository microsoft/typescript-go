package core

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

// Id-indexed links store

const (
	idLinkPageShift = 8
	idLinkPageSize  = 1 << idLinkPageShift
	idLinkPageMask  = idLinkPageSize - 1
)

// IdLinkStore is a links store keyed by a dense uint64 id. It uses a paged array indexed by id
// instead of a map, which is significantly cheaper for the hottest stores.
type IdLinkStore[V any] struct {
	pages [][]*V // pages[id>>idLinkPageShift][id&idLinkPageMask]
	arena Arena[V]
}

func (s *IdLinkStore[V]) Get(id uint64) *V {
	page := id >> idLinkPageShift
	if page < uint64(len(s.pages)) {
		if p := s.pages[page]; p != nil {
			if value := p[id&idLinkPageMask]; value != nil {
				return value
			}
		}
	}
	return s.getSlow(id)
}

func (s *IdLinkStore[V]) getSlow(id uint64) *V {
	page := id >> idLinkPageShift
	if page >= uint64(len(s.pages)) {
		s.pages = append(s.pages, make([][]*V, int(page)+1-len(s.pages))...)
	}
	if s.pages[page] == nil {
		s.pages[page] = make([]*V, idLinkPageSize)
	}
	slot := &s.pages[page][id&idLinkPageMask]
	if *slot == nil {
		*slot = s.arena.New()
	}
	return *slot
}

func (s *IdLinkStore[V]) Has(id uint64) bool {
	return s.TryGet(id) != nil
}

func (s *IdLinkStore[V]) TryGet(id uint64) *V {
	page := id >> idLinkPageShift
	if page < uint64(len(s.pages)) {
		if p := s.pages[page]; p != nil {
			return p[id&idLinkPageMask]
		}
	}
	return nil
}
