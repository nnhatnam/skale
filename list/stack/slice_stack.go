package stack

var (
	_ Stack[any] = (*SStack[any])(nil)
)

type SStack[T any] struct {
	e []T
}

func SStackFrom[T any](s []T) *SStack[T] {
	return &SStack[T]{e: s}
}

func NewSStack[T any]() *SStack[T] {
	return &SStack[T]{e: make([]T, 0)}
}

func (s *SStack[T]) Empty() bool {
	return len(s.e) == 0
}

func (s *SStack[T]) Top() (_ T, _ bool) {
	if len(s.e) == 0 {
		return
	}
	return s.e[len(s.e)-1], true
}

func (s *SStack[T]) Len() int {
	return len(s.e)
}

func (s *SStack[T]) Push(v T) {
	s.e = append(s.e, v)
}

func (s *SStack[T]) Pop() (_ T, _ bool) {
	if len(s.e) == 0 {
		return
	}
	var zero T
	x := s.e[len(s.e)-1]
	s.e[len(s.e)-1] = zero // clear reference
	s.e = s.e[:len(s.e)-1]
	return x, true
}

func (s *SStack[T]) Clear() {
	s.e = nil
}

func (s *SStack[T]) ToSlice() []T {
	arr := make([]T, len(s.e))
	for i, j := 0, len(s.e)-1; j >= 0; i, j = i+1, j-1 {
		arr[i] = s.e[j]
	}

	return arr
}

// Shrink copies the underlying slice with excess capacity to precisely sized one to avoid wasting memory.
// It should be called on stack with long static durations.
// Long-lived slices can waste memory on unused capacity, shrink them
func (s *SStack[T]) Shrink() {
	//credit from https://about.sourcegraph.com/blog/zoekt-memory-optimizations-for-sourcegraph-cloud
	if cap(s.e)-len(s.e) < 32 {
		return
	}

	out := make([]T, len(s.e))
	copy(out, s.e)
	s.e = out

}
