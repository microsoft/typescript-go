package core

type Stack[T any] struct {
	data []T
}

// Push adds an item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.data = append(s.data, item)
}

// Pop removes and returns the top item from the stack.
// Returns the item and a boolean indicating success.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		var zero T // Declare zero value for type T
		return zero, false
	}
	l := len(s.data) - 1
	item := s.data[l]
	s.data = s.data[:l] // Shrink the slice to remove the last element
	return item, true
}

// Peek returns the top item without removing it.
// Returns the item and a boolean indicating success.
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

// Len returns the number of items in the stack.
func (s *Stack[T]) Len() int {
	return len(s.data)
}

// Clear removes all elements from the stack.
func (s *Stack[T]) Clear() {
	s.data = nil
}
