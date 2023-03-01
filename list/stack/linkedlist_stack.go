package stack

import (
	"github.com/nnhatnam/skale/list/linkedlist"
)

var (
	_ Stack[any] = (*LinkedListStack[any])(nil)
)

type LinkedListStack[T any] struct {
	e *linkedlist.List[T]
}

func LinkedListStackFrom[T any](l *linkedlist.List[T]) *LinkedListStack[T] {
	return &LinkedListStack[T]{e: l}
}

func NewLinkedListStack[T any]() *LinkedListStack[T] {
	return &LinkedListStack[T]{e: linkedlist.New[T]()}
}

func (s *LinkedListStack[T]) Top() (_ T, _ bool) {
	var zero T
	n := s.e.Back()
	if n == nil {
		return zero, false
	}
	return n.Value, true
}

func (s *LinkedListStack[T]) Len() int {
	return s.e.Len()
}

func (s *LinkedListStack[T]) Push(value T) {
	s.e.PushBack(value)
}

func (s *LinkedListStack[T]) Pop() (_ T, _ bool) {
	var zero T
	n := s.e.PopBack()
	if n == nil {
		return zero, false
	}
	return n.Value, true
}

func (s *LinkedListStack[T]) Empty() bool {
	return s.e.Len() == 0
}

func (s *LinkedListStack[T]) Clear() {
	s.e.Init()
}

func (s *LinkedListStack[T]) ToSlice() []T {
	arr := make([]T, s.e.Len())
	c := s.e.Cursor()
	i := 0
	c.WalkDescending(func(n *linkedlist.Node[T]) bool {
		arr[i] = n.Value
		i++
		return true
	})
	return arr
}
