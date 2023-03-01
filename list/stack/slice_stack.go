package stack

var (
	_ Stack[any] = (*SliceStack[any])(nil)
)

type SliceStack[T any] struct {
	e []T
}

func (s *SliceStack[T]) Empty() bool {
	return len(s.e) == 0
}

func NewSliceStack[T any]() *SliceStack[T] {
	return &SliceStack[T]{e: make([]T, 0)}
}

func (s *SliceStack[T]) Top() (_ T, _ bool) {
	if len(s.e) == 0 {
		return
	}
	return s.e[len(s.e)-1], true
}

func (s *SliceStack[T]) Len() int {
	return len(s.e)
}

func (s *SliceStack[T]) Push(v T) {
	s.e = append(s.e, v)
}

func (s *SliceStack[T]) Pop() (_ T, _ bool) {
	if len(s.e) == 0 {
		return
	}
	var zero T
	x := s.e[len(s.e)-1]
	s.e[len(s.e)-1] = zero // clear reference
	s.e = s.e[:len(s.e)-1]
	return x, true
}

func (s *SliceStack[T]) Clear() {
	s.e = nil
}

func (s *SliceStack[T]) ToSlice() []T {
	arr := make([]T, len(s.e))
	for i, j := 0, len(s.e)-1; j >= 0; i, j = i+1, j-1 {
		arr[i] = s.e[j]
	}

	return arr
}
