package collections

import "iter"

type Set[T comparable] struct {
	m Map[T, struct{}]
}

func NewSetWithSizeHint[T comparable](hint int) *Set[T] {
	return &Set[T]{
		m: newMapWithSizeHint[T, struct{}](hint),
	}
}

func (s *Set[T]) Add(value T) {
	s.m.Set(value, struct{}{})
}

func (s *Set[T]) Has(value T) bool {
	return s.m.Has(value)
}

func (s *Set[T]) Delete(value T) bool {
	_, ok := s.m.Delete(value)
	return ok
}

func (s *Set[T]) Values() iter.Seq[T] {
	return s.m.Keys()
}

func (s *Set[T]) Size() int {
	return s.m.Size()
}

func (s *Set[T]) Clone() *Set[T] {
	return &Set[T]{
		m: s.m.clone(),
	}
}

func (s *Set[T]) Clear() {
	s.m.Clear()
}
