package collections

type SyncSet[T comparable] struct {
	m SyncMap[T, struct{}]
}

func (s *SyncSet[T]) Has(key T) bool {
	_, ok := s.m.Load(key)
	return ok
}

func (s *SyncSet[T]) Add(key T) {
	s.m.Store(key, struct{}{})
}

// AddIfAbsent adds the key to the set if it is not already present
// using LoadOrStore. It returns true if the key was not already present
// (opposite of the return value of LoadOrStore).
func (s *SyncSet[T]) AddIfAbsent(key T) bool {
	_, loaded := s.m.LoadOrStore(key, struct{}{})
	return !loaded
}

func (s *SyncSet[T]) Delete(key T) {
	s.m.Delete(key)
}

func (s *SyncSet[T]) Range(f func(key T) bool) {
	s.m.Range(func(key T, _ struct{}) bool {
		return f(key)
	})
}
