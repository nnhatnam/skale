package slice

import (
	"github.com/nnhatnam/skale"
	"sort"
)

type OrderedSlice[T skale.Ordered] []T //consider to replace with constraints.Ordered

func FromOrdered[T skale.Ordered](values ...T) OrderedSlice[T] {
	return values
}

func (s OrderedSlice[T]) Len() int {
	return len(s)
}

func (s OrderedSlice[T]) Cap() int {
	return cap(s)
}

func (s OrderedSlice[T]) At(i int) T {
	if i < 0 || i >= len(s) {
		panic(errIndexOutOfRange.Error())
	}
	return s[i]
}

func (s OrderedSlice[T]) Set(i int, value T) {
	s[i] = value
}

func (s OrderedSlice[T]) Append(value ...T) OrderedSlice[T] {
	return append(s, value...)
}

func (s OrderedSlice[T]) Cut(i, j int) OrderedSlice[T] {
	var zero T
	copy(s[i:], s[j:])
	for k, n := len(s)-j+i, len(s); k < n; k++ {
		s[k] = zero
	}
	return s[:len(s)-j+i]
}

func (s OrderedSlice[T]) Delete(i int) OrderedSlice[T] {
	var zero T
	copy(s[i:], s[i+1:])
	s[len(s)-1] = zero
	return s[:len(s)-1]
}

func (s OrderedSlice[T]) DeleteUnordered(i int) OrderedSlice[T] {
	s[i] = s[len(s)-1]
	var zero T
	s[len(s)-1] = zero
	return s[:len(s)-1]
}

func (s OrderedSlice[T]) Expand(i int, value ...T) OrderedSlice[T] {
	s = append(s, value...)
	copy(s[i+len(value):], s[i:])
	copy(s[i:], value)
	return s
}

func (s OrderedSlice[T]) Insert(i int, value ...T) OrderedSlice[T] {
	s = append(s, value...)
	copy(s[i+len(value):], s[i:])
	copy(s[i:], value)
	return s
}

func (s OrderedSlice[T]) Clear() OrderedSlice[T] {
	return s[:0]
}

func (s OrderedSlice[T]) Copy() OrderedSlice[T] {
	return append(OrderedSlice[T]{}, s...)
}

func (s OrderedSlice[T]) Reverse() OrderedSlice[T] {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func (s OrderedSlice[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s OrderedSlice[T]) Sort() {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

func (s OrderedSlice[T]) SortBy(less func(i, j int) bool) {
	sort.Slice(s, less)
}

func (s OrderedSlice[T]) Contains(value T, less skale.LessFunc[T]) bool {
	for _, v := range s {
		if !less(v, value) && !less(value, v) {
			return true
		}
	}
	return false
}
