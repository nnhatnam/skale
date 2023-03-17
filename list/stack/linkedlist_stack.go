package stack

import (
	"github.com/nnhatnam/skale/list/linkedlist"
)

var (
	_ Stack[any] = (*StackL[any])(nil)
)

// StackL is a stack implemented using a linked list.
type StackL[T any] linkedlist.List[T]

// NewStackL creates a new stack.
func NewStackL[T any]() *StackL[T] {
	l := linkedlist.New[T]()
	return (*StackL[T])(l)
}

// Top returns the top element of the stack.
func (s *StackL[T]) Top() (_ T, _ bool) {
	var zero T

	n := (*linkedlist.List[T])(s).Back()
	if n == nil {
		return zero, false
	}
	return n.Value, true
}

// Len returns the length of the stack.
func (s *StackL[T]) Len() int {
	return (*linkedlist.List[T])(s).Len()
}

// Push adds an element to the stack.
func (s *StackL[T]) Push(value T) {
	(*linkedlist.List[T])(s).PushBack(value)
}

// Pop removes an element from the stack.
func (s *StackL[T]) Pop() (_ T, _ bool) {
	var zero T
	n := (*linkedlist.List[T])(s).PopBack()
	if n == nil {
		return zero, false
	}
	return n.Value, true
}

// Empty returns true if the stack is empty.
func (s *StackL[T]) Empty() bool {
	return s.Len() == 0
}

// Clear clears the stack.
func (s *StackL[T]) Clear() {
	(*linkedlist.List[T])(s).Init()
}

// ToSlice returns the stack as a slice.
func (s *StackL[T]) ToSlice() []T {
	arr := make([]T, s.Len())
	c := (*linkedlist.List[T])(s).Cursor()
	i := 0
	c.WalkDescending(func(n *linkedlist.Node[T]) bool {
		arr[i] = n.Value
		i++
		return true
	})
	return arr
}
