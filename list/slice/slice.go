package slice

import (
	"errors"
	"github.com/nnhatnam/skale"
	"sort"
)

var (
	errIndexOutOfRange   = errors.New("index out of range")
	errIndexReachMaxSize = errors.New("index reach max size")
)

type Slice[T any] []T

func From[T any](values ...T) Slice[T] {
	return values
}

func (s Slice[T]) Len() int {
	return len(s)
}

func (s Slice[T]) Cap() int {
	return cap(s)
}

func (s Slice[T]) At(i int) T {
	if i < 0 || i >= len(s) {
		panic(errIndexOutOfRange.Error())
	}
	return s[i]
}

func (s Slice[T]) Set(i int, value T) {
	s[i] = value
}

func (s Slice[T]) Append(value ...T) Slice[T] {
	return append(s, value...)
}

func (s Slice[T]) Cut(i, j int) Slice[T] {
	var zero T
	copy(s[i:], s[j:])
	for k, n := len(s)-j+i, len(s); k < n; k++ {
		s[k] = zero
	}
	return s[:len(s)-j+i]
}

func (s Slice[T]) Delete(i int) Slice[T] {
	var zero T
	copy(s[i:], s[i+1:])
	s[len(s)-1] = zero
	return s[:len(s)-1]
}

func (s Slice[T]) DeleteUnordered(i int) Slice[T] {
	s[i] = s[len(s)-1]
	var zero T
	s[len(s)-1] = zero
	return s[:len(s)-1]
}

func (s Slice[T]) Expand(i int, value ...T) Slice[T] {
	s = append(s, value...)
	copy(s[i+len(value):], s[i:])
	copy(s[i:], value)
	return s
}

func (s Slice[T]) Insert(i int, value ...T) Slice[T] {
	s = append(s, value...)
	copy(s[i+len(value):], s[i:])
	copy(s[i:], value)
	return s
}

func (s Slice[T]) Clear() Slice[T] {
	return s[:0]
}

func (s Slice[T]) Copy() Slice[T] {
	return append(Slice[T]{}, s...)
}

func (s Slice[T]) Reverse() Slice[T] {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func (s Slice[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//
//func (s Slice[T]) Sort(less func(i, j int) bool) {
//	sort.Slice(s, less)
//}

func (s Slice[T]) SortBy(less skale.LessFunc[T]) {
	sort.Slice(s, func(i, j int) bool {
		return less(s[i], s[j])
	})
}

func (s Slice[T]) Contains(value T, less skale.LessFunc[T]) bool {
	for _, v := range s {
		if !less(v, value) && !less(value, v) {
			return true
		}
	}
	return false
}
