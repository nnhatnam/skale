// slices implements a set of useful functions for manipulating slices that have not provided by golang.org/x/exp/slices yet
// This package is just a temporary solution for internal use. It will be removed when golang.org/x/exp/slices provides the same functions.
package xslices

// LongestPrefixIndex returns the index of the last element in the longest common prefix of s1 and s2.
func LongestPrefixIndex[T comparable](s1, s2 []T) (_ int) {
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}

	//if len(s1) == 0 || s1[0] != s2[0] {
	//	return -1
	//}

	for i, v := range s1 {
		if v != s2[i] {
			return i - 1 //collision
		}
	}
	return len(s1) - 1
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
