package queue

var _ Queue[any] = (*SQueue[any])(nil)

type SQueue[T any] struct {
	e []T
}

func SQueueFrom[T any](l []T) *SQueue[T] {
	return &SQueue[T]{e: l}
}

func NewSQueue[T any]() *SQueue[T] {
	return &SQueue[T]{e: make([]T, 0)}
}

func (s *SQueue[T]) Len() int {
	return len(s.e)
}

func (s *SQueue[T]) Enqueue(value T) {
	s.e = append(s.e, value)
}

func (s *SQueue[T]) Dequeue() (_ T, _ bool) {
	var zero T
	if len(s.e) == 0 {
		return zero, false
	}
	v := s.e[0]
	s.e = s.e[1:]
	return v, true
}

func (s *SQueue[T]) Peek() (_ T, _ bool) {
	var zero T
	if len(s.e) == 0 {
		return zero, false
	}
	return s.e[0], true
}

func (s *SQueue[T]) Empty() bool {
	return len(s.e) == 0
}

func (s *SQueue[T]) Clear() {
	s.e = make([]T, 0)
}

func (s *SQueue[T]) ToSlice() []T {
	return s.e
}

// Shrink copies the underlying slice with excess capacity to precisely sized one to avoid wasting memory.
// It should be called on queue with long static durations.
// Long-lived slices can waste memory on unused capacity, shrink them
func (s *SQueue[T]) Shrink() {
	//credit from https://about.sourcegraph.com/blog/zoekt-memory-optimizations-for-sourcegraph-cloud
	if cap(s.e)-len(s.e) < 32 {
		return
	}

	out := make([]T, len(s.e))
	copy(out, s.e)
	s.e = out

}
