package queue

var _ Queue[any] = (*SliceQueue[any])(nil)

type SliceQueue[T any] struct {
	e []T
}

func SliceQueueFrom[T any](l []T) *SliceQueue[T] {
	return &SliceQueue[T]{e: l}
}

func NewSliceQueue[T any]() *SliceQueue[T] {
	return &SliceQueue[T]{e: make([]T, 0)}
}

func (s *SliceQueue[T]) Len() int {
	return len(s.e)
}

func (s *SliceQueue[T]) Enqueue(value T) {
	s.e = append(s.e, value)
}

func (s *SliceQueue[T]) Dequeue() (_ T, _ bool) {
	var zero T
	if len(s.e) == 0 {
		return zero, false
	}
	v := s.e[0]
	s.e = s.e[1:]
	return v, true
}

func (s *SliceQueue[T]) Peek() (_ T, _ bool) {
	var zero T
	if len(s.e) == 0 {
		return zero, false
	}
	return s.e[0], true
}

func (s *SliceQueue[T]) Empty() bool {
	return len(s.e) == 0
}

func (s *SliceQueue[T]) Clear() {
	s.e = make([]T, 0)
}

func (s *SliceQueue[T]) ToSlice() []T {
	return s.e
}
