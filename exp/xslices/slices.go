// slices implements a set of useful functions for manipulating slices that have not provided by golang.org/x/exp/slices yet
// This package is just a temporary solution for internal use. It will be removed when golang.org/x/exp/slices provides the same functions.
package xslices

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

	return -1 // no common prefix
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
