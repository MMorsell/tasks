package internal

import "fmt"

func NewStack[T any](v []T) *Stack[T] {
	return &Stack[T]{
		items: v,
	}
}

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(f T) {
	s.items = append(s.items, f)
}

func (s *Stack[T]) IsEmpty() bool {
	if len(s.items) == 0 {
		return true
	}
	return false
}

func (s *Stack[T]) Top() (T, error) {
	if s.IsEmpty() {
		var zero T
		return zero, fmt.Errorf("stack is empty")
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top, nil
}
