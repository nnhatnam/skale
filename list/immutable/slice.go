package immutable

type Slice[T any] []T

func NewImmutableSlice[T any](elems ...T) Slice[T] {
	s := Slice[T](elems)
	return s
}

func (s Slice[T]) Len() int {
	return len(s)
}

func (s Slice[T]) Get(i int) (_ T, _ bool) {
	if i < 0 || i >= len(s) {
		return
	}
	return s[i], true
}

func (s Slice[T]) Set(i int, v T) (_ Slice[T]) {
	if i < 0 || i >= len(s) {
		panic("index out of range")

	}
	s2 := make([]T, len(s))
	copy(s2, s)
	s2[i] = v
	return s2
}

func (s Slice[T]) Insert(i int, v T) (_ Slice[T]) {
	if i < 0 || i > len(s) {
		panic("index out of range")
	}
	s2 := make([]T, len(s)+1)
	copy(s2, s[:i])
	s2[i] = v
	copy(s2[i+1:], s[i:])
	return s2
}
