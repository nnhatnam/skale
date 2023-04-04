package stack

var (
	_ Stack[any] = (*StackS[any])(nil)
)

// StackS is a stack implemented using slice. StackS[T] is literally []T.
// So one can easily convert back and forth between StackS[T] and []T.
// For example:
// var a = []int{1, 2, 3}
// s := StackS[int](a)
// a = []int(s)
type StackS[T any] []T

// NewStackSWithSize creates a new stack with given size
func NewStackSWithSize[T any](size int) *StackS[T] {
	var s StackS[T]
	s = make([]T, 0, size)
	return &s
}

// NewStackS creates a new stack
func NewStackS[T any]() *StackS[T] {
	return NewStackSWithSize[T](0)
}

// Empty returns true if the stack is empty
func (s *StackS[T]) Empty() bool {
	return len(*s) == 0
}

// Top returns the top element of the stack without removing it,
// and a boolean indicating if there is an element
func (s *StackS[T]) Top() (_ T, _ bool) {

	if len(*s) == 0 {
		return
	}

	return (*s)[len(*s)-1], true
}

// Peek simply returns the top element of the stack if there is one. Return zero-value if not.
func (s *StackS[T]) Peek() (_ T) {

	if len(*s) == 0 {
		return
	}

	return (*s)[len(*s)-1]
}

// Bottom returns the bottom element of the stack
func (s *StackS[T]) Bottom() (_ T, _ bool) {

	if len(*s) == 0 {
		return
	}

	return (*s)[0], true
}

// Len returns the length of the stack
func (s *StackS[T]) Len() int {
	return len(*s)
}

// Push adds an element to the stack
func (s *StackS[T]) Push(v T) {
	*s = append(*s, v)
}

// Pop removes an element from the stack
func (s *StackS[T]) Pop() (_ T, _ bool) {
	if len(*s) == 0 {
		return
	}
	var zero T
	n := len(*s) - 1
	x := (*s)[n]
	(*s)[n] = zero // clear reference
	*s = (*s)[:n]
	return x, true

}

// Clear clears the stack
func (s *StackS[T]) Clear() {
	//clear the underlying slice
	*s = make([]T, 0, cap(*s))
}

// ToSlice returns the slice of the stack
func (s *StackS[T]) ToSlice() []T {

	arr := make([]T, len(*s))

	//reverse the slice
	for i, j := 0, len(*s)-1; j >= 0; i, j = i+1, j-1 {
		arr[i] = (*s)[j]
	}

	return arr
}

// Shrink copies the underlying slice with excess capacity to precisely sized one to avoid wasting memory.
// It should be called on stack with long static durations.
// Long-lived slices can waste memory on unused capacity, shrink them
func (s *StackS[T]) Shrink() {
	//credit from https://about.sourcegraph.com/blog/zoekt-memory-optimizations-for-sourcegraph-cloud

	if cap(*s)-len(*s) < 32 {
		return
	}

	out := make([]T, len(*s))
	copy(out, *s)
	*s = out

}
