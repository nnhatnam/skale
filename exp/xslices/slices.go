// slices implements a set of useful functions for manipulating slices that have not provided by golang.org/x/exp/slices yet
// This package is just a temporary solution for internal use. It will be removed when golang.org/x/exp/slices provides the same functions.
package xslices

func LongestPrefixIndex[T comparable](s1, s2 []T) (_ int) {
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}

	if len(s1) == 0 || s1[0] != s2[0] {
		return -1
	}

	for i, v := range s1 {
		if v != s2[i] {
			return i //collision
		}
	}
	return len(s1) - 1
}
