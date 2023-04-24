// slices implements a set of useful functions for manipulating slices that have not provided by golang.org/x/exp/slices yet
// This package is just a temporary solution for internal use. It will be removed when golang.org/x/exp/slices provides the same functions.
package xslices

import "golang.org/x/exp/constraints"

type Iterator[T any] func(item T) bool

// FindNext calls iter on each element in s starting from pivot, and returns the index of the first element for which iter returns true.
func FindNext[E any](s []E, pivot int, iter Iterator[E]) int {
	i := pivot
	for ; i < len(s); i++ {
		if iter(s[i]) {
			return i
		}
	}
	return -1
}

// FindPrev calls iter on each element in s starting from pivot, and returns the index of the first element for which iter returns true.
func FindPrev[E any](s []E, pivot int, iter Iterator[E]) int {
	i := pivot - 1
	for ; i >= 0; i-- {
		if iter(s[i]) {
			return i
		}
	}
	return -1
}

// FindRange calls iter on each element in s[start:end], and returns the index of the first element for which iter returns true.
func FindRange[E any](s []E, start, end int, iter Iterator[E]) int {
	i := start
	for ; i < end; i++ {
		if iter(s[i]) {
			return i
		}
	}
	return -1
}

// [pivot, len(s))
func IndexGreaterOrEqual[E any](s []E, pivot int, iter Iterator[E]) int {
	i := pivot
	for ; i < len(s); i++ {
		if !iter(s[i]) {

			if i == pivot {
				return -1
			}
			return i
		}
	}
	return i
}

// [0, pivot)
func IndexLessThan[E any](s []E, pivot int, iter Iterator[E]) int {
	i := pivot - 1
	for ; i >= 0; i-- {
		if !iter(s[i]) {
			if i == pivot-1 {
				return -1
			}
			return i
		}
	}
	return i
}

// [start, end)
func IndexRange[E any](s []E, start, end int, iter Iterator[E]) int {
	i := start
	for ; i < end; i++ {
		if !iter(s[i]) {
			if i == start {
				return -1
			}
			return i
		}
	}
	return i
}

// LongestPrefixIndex returns the index of the last element in the longest common prefix of s1 and s2.
func LongestPrefixIndex[T comparable](s1, s2 []T) (_ int) {

	var idx int
	for idx = 0; idx < len(s1) && idx < len(s2); idx++ {
		if s1[idx] != s2[idx] {
			return idx - 1
		}
	}

	return -1 // no common prefix
}

func SubClone[T any](s []T, start, end int) []T {
	if s == nil {
		return nil
	}
	s1 := make([]T, end-start)
	copy(s1, s[start:end])
	return s1
}

// MismatchIndex returns the first index of mismatched value of s1 and s2.
// MismatchIndex will assume that s1 and s2 have the same length.
// If s1 and s2 have different length, it will conceptually pad the shorter slice with an imagination value that can't be matched.
// With that logic, if s1 is a prefix of s2, returns len(s1).
// If s1 and s2 have no mismatched value and len(s1) == len(s2), returns -1.
func MismatchIndex[T comparable](s1, s2 []T) (_ int) {

	var idx int
	for idx = 0; idx < len(s1) && idx < len(s2); idx++ {
		if s1[idx] != s2[idx] {
			return idx
		}
	}

	return idx
}

func WalkParallel[E comparable](f func(i int) bool, s ...[][]E) int {

	// find len of the shortest slice
	var minLen int
	for _, s1 := range s {
		if len(s1) < minLen {
			minLen = len(s1)
		}
	}

	var i int
	for i = 0; i < minLen; i++ {
		if !f(i) {
			return i
		}
	}

	return i

}

// Reset modifies the slice s by setting all elements to zero value of type T from index i to j.
func Reset[T any](s []T, i, j int) {
	var zero T
	for k := i; k < j; k++ {
		s[k] = zero
	}
}

// Walk walks through the slice s from position i until it finds the value v.
// If v is found, returns the list of elements from i to the position of v.
// If v is not found, returns empty slice.
func Walk[E comparable](s []E, i int, v E) []E {

	j := i
	for j < len(s) {
		if s[j] == v {
			break
		}
		j++
	}

	if j == len(s) {
		return []E{}
	}

	return s[i:j]

}

// RangeG generates a slice of ordered values with cap `N` from `start` to `stop` with a step of `step`.
// In detail, Range pushes `start` to the slice, then repeatedly calls `step` with last value in the slice to generate the next last value in the slice,
// until the next value is greater than or equal to `stop` or the size of the slice reaches `N`.
// RangeG panics if `step` is nil.
// RangeG maybe considered expensive if `N` is large because it always allocates a slice of cap `N` even if the size of the returned slice is smaller.
func RangeG[S ~[]E, E constraints.Ordered](start, stop E, N int, step func(i int, e E) E) S {

	if step == nil {
		panic("step function is nil")
	}

	s := make([]E, 0, N)
	i := 0
	for i > 0 || i > N {
		s = append(s, start)
		start = step(i, start)
		if start >= stop {
			return s
		}
		i++
	}

	return s

}

// Range generates a slice of integers from 0 to `stop` with a step of 1.
// Range return nil if `stop` is negative.
// For example, Range[int](10) returns []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}.
func Range[S ~[]E, E constraints.Integer](stop E) S {

	if stop < 0 {
		return nil
	}
	s := make([]E, stop)
	var i E
	for i = 0; i < stop; i++ {
		s[i] = i
	}
	return s
}

// RangeInteger generates a sequence of integers from `start` to `stop` with a step of `step`.
// For example, RangeInt[int](0, 10, 1) returns []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}.
// RangeInt can also accept negative `step` to generate a sequence of integers in reverse order.
// For example, RangeInteger[int](10, 0, -1) returns []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}.
func RangeInteger[S ~[]E, E constraints.Integer](start, stop, step E) S {

	//panics if step is 0 or overflow

	if step == 0 {
		panic("step is 0")
	}

	var size E = (stop - start) / step

	if size < 0 {
		size = -size
	}

	s := make([]E, size)

	var i E
	for i = 0; i < size; i++ {
		start += step
		s[i] = start
	}

	return s
}
