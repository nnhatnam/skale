package queue

var _ Queue[any] = (*QueueS[any])(nil)

// QueueS is a queue implemented using a slice. QueueS[T] is a type alias of []T.
// So one can easily convert back and forth between QueueS[T] and []T.
// For example:
// var a = []int{1, 2, 3}
// q := QueueS[int](a)
// a = []int(q)
type QueueS[T any] []T

// NewQueueSWithSize creates a new queue with given size
func NewQueueSWithSize[T any](size int) *QueueS[T] {
	var q QueueS[T] = make([]T, 0, size)
	return &q
}

// NewQueueS creates a new queue
func NewQueueS[T any]() *QueueS[T] {
	return NewQueueSWithSize[T](0)
}

// Len returns the length of the queue
func (s *QueueS[T]) Len() int {
	return len(*s)
}

// Enqueue adds an element to the queue
func (s *QueueS[T]) Enqueue(value T) {
	*s = append(*s, value)
}

// Dequeue removes an element from the queue
func (s *QueueS[T]) Dequeue() (_ T, _ bool) {
	var zero T
	if len(*s) == 0 {
		return zero, false
	}
	v := (*s)[0]
	*s = (*s)[1:]
	return v, true
}

// Peek returns the first element of the queue
func (s *QueueS[T]) Peek() (_ T, _ bool) {
	var zero T
	if len(*s) == 0 {
		return zero, false
	}
	return (*s)[0], true
}

// Empty returns true if the queue is empty
func (s *QueueS[T]) Empty() bool {
	return len(*s) == 0
}

// Clear clears the queue
func (s *QueueS[T]) Clear() {
	*s = make([]T, 0)
}

// ToSlice returns the underlying slice
func (s *QueueS[T]) ToSlice() []T {
	return *s
}

// Shrink copies the underlying slice with excess capacity to precisely sized one to avoid wasting memory.
// It should be called on queue with long static durations.
// Long-lived slices can waste memory on unused capacity, shrink them
func (s *QueueS[T]) Shrink() {
	//credit from https://about.sourcegraph.com/blog/zoekt-memory-optimizations-for-sourcegraph-cloud
	if cap(*s)-len(*s) < 32 {
		return
	}

	out := make([]T, len(*s))
	copy(out, *s)
	*s = out
}
