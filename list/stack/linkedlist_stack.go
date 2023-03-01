package stack

import (
	"github.com/nnhatnam/skale/list/linkedlist"
)

var (
	_ Stack[any] = (*LLStack[any])(nil)
)

type LLStack[T any] struct {
	e *linkedlist.List[T]
}

func LLStackFrom[T any](l *linkedlist.List[T]) *LLStack[T] {
	return &LLStack[T]{e: l}
}

func NewLLStack[T any]() *LLStack[T] {
	return &LLStack[T]{e: linkedlist.New[T]()}
}

func (s *LLStack[T]) Top() (_ T, _ bool) {
	var zero T
	n := s.e.Back()
	if n == nil {
		return zero, false
	}
	return n.Value, true
}

func (s *LLStack[T]) Len() int {
	return s.e.Len()
}

func (s *LLStack[T]) Push(value T) {
	s.e.PushBack(value)
}

func (s *LLStack[T]) Pop() (_ T, _ bool) {
	var zero T
	n := s.e.PopBack()
	if n == nil {
		return zero, false
	}
	return n.Value, true
}

func (s *LLStack[T]) Empty() bool {
	return s.e.Len() == 0
}

func (s *LLStack[T]) Clear() {
	s.e.Init()
}

func (s *LLStack[T]) ToSlice() []T {
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
