package queue

import (
	"github.com/nnhatnam/skale/list/linkedlist"
)

var (
	_ Queue[any] = (*QueueL[any])(nil)
)

type QueueL[T any] linkedlist.List[T]

// NewQueueL returns a new QueueL[T]
func NewQueueL[T any]() *QueueL[T] {
	l := linkedlist.New[T]()
	return (*QueueL[T])(l)
}

// Len returns the number of elements in the queue.
func (q *QueueL[T]) Len() int {
	return (*linkedlist.List[T])(q).Len()
}

// Enqueue adds an element to the queue.
func (q *QueueL[T]) Enqueue(value T) {
	(*linkedlist.List[T])(q).PushBack(value)
}

// Dequeue removes an element from the queue.
func (q *QueueL[T]) Dequeue() (_ T, _ bool) {

	return (*linkedlist.List[T])(q).PopFront()

}

// Peek returns the first element of the queue.
func (q *QueueL[T]) Peek() (_ T, _ bool) {
	return (*linkedlist.List[T])(q).Front()
}

// Empty returns true if the queue is empty.
func (q *QueueL[T]) Empty() bool {
	return q.Len() == 0
}

// Clear clears the queue.
func (q *QueueL[T]) Clear() {
	(*linkedlist.List[T])(q).Init()
}

// ToSlice returns the underlying slice.
func (q *QueueL[T]) ToSlice() []T {
	var arr []T
	c := (*linkedlist.List[T])(q).Cursor()

	c.WalkAscending(func(v T) bool {
		arr = append(arr, v)
		return true
	})
	return arr
}
