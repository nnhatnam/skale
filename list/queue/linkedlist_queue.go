package queue

import (
	"github.com/nnhatnam/skale/list/linkedlist"
)

var (
	_ Queue[any] = (*LinkedListQueue[any])(nil)
)

type LinkedListQueue[T any] struct {
	e *linkedlist.List[T]
}

func LinkedListQueueFrom[T any](l *linkedlist.List[T]) *LinkedListQueue[T] {
	return &LinkedListQueue[T]{e: l}
}

func NewLinkedListQueue[T any]() *LinkedListQueue[T] {
	return &LinkedListQueue[T]{e: linkedlist.New[T]()}
}

func (s *LinkedListQueue[T]) Len() int {
	return s.e.Len()
}

func (s *LinkedListQueue[T]) Enqueue(value T) {
	s.e.PushBack(value)
}

func (s *LinkedListQueue[T]) Dequeue() (_ T, _ bool) {
	var zero T
	n := s.e.PopFront()
	if n == nil {
		return zero, false
	}
	return n.Value, true
}

func (s *LinkedListQueue[T]) Peek() (_ T, _ bool) {
	var zero T
	n := s.e.Front()
	if n == nil {
		return zero, false
	}
	return n.Value, true
}

func (s *LinkedListQueue[T]) Empty() bool {
	return s.e.Len() == 0
}

func (s *LinkedListQueue[T]) Clear() {
	s.e.Init()
}

func (s *LinkedListQueue[T]) ToSlice() []T {
	var arr []T
	c := s.e.Cursor()

	c.WalkAscending(func(n *linkedlist.Node[T]) bool {
		arr = append(arr, n.Value)
		return true
	})
	return arr
}
