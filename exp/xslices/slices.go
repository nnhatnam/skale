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

// Reset modifies the slice s by setting all elements to zero value of type T from index i to j.
func Reset[T any](s []T, i, j int) {
	var zero T
	for k := i; k < j; k++ {
		s[k] = zero
	}
}

// Insert inserts the values v... into s at index i,
// returning the modified slice.
// In the returned slice r, r[i] == v[0].
// Insert panics if i is out of range.
// This function is O(len(s) + len(v)).
// This is a temporary fix for golang.org/x/exp/slices#Insert
// slices.Insert add extra `i` zero-value elements to the end of the slice causing the length to be longer than expected.
//func Insert[S ~[]E, E any](s S, i int, v ...E) S {
//	tot := len(s) + len(v) // 2 + 7 = 9
//	if tot <= cap(s) {
//		s2 := s[:tot]
//		copy(s2[i+len(v):], s[i:])
//		copy(s2[i:], v)
//		return s2
//	}
//	s2 := make(S, len(v)+i) // <-- change to make(S, len(v)+i, tot)
//	copy(s2, s[:i])
//	copy(s2[i:], v)
//	copy(s2[i+len(v):], s[i:])
//	return s2
//}

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
