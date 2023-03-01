package queue

import (
	"github.com/nnhatnam/skale/list/linkedlist"
)

var (
	_ Queue[any] = (*LLQueue[any])(nil)
)

type LLQueue[T any] struct {
	e *linkedlist.List[T]
}

func LLQueueFrom[T any](l *linkedlist.List[T]) *LLQueue[T] {
	return &LLQueue[T]{e: l}
}

func NewLLQueue[T any]() *LLQueue[T] {
	return &LLQueue[T]{e: linkedlist.New[T]()}
}

func (s *LLQueue[T]) Len() int {
	return s.e.Len()
}

func (s *LLQueue[T]) Enqueue(value T) {
	s.e.PushBack(value)
}

func (s *LLQueue[T]) Dequeue() (_ T, _ bool) {
	var zero T
	n := s.e.PopFront()
	if n == nil {
		return zero, false
	}
	return n.Value, true
}

func (s *LLQueue[T]) Peek() (_ T, _ bool) {
	var zero T
	n := s.e.Front()
	if n == nil {
		return zero, false
	}
	return n.Value, true
}

func (s *LLQueue[T]) Empty() bool {
	return s.e.Len() == 0
}

func (s *LLQueue[T]) Clear() {
	s.e.Init()
}

func (s *LLQueue[T]) ToSlice() []T {
	var arr []T
	c := s.e.Cursor()

	c.WalkAscending(func(n *linkedlist.Node[T]) bool {
		arr = append(arr, n.Value)
		return true
	})
	return arr
}
