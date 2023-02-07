package datrie

import "github.com/nnhatnam/skale/trie"

// ArcDomain represents the mapping between arc labels and their corresponding codes.
// The "Code" and "Get" methods allows for the conversion between the labels and codes,
// while the "StopElement" method returns a special element that indicates the end of a word.
// The "StopElement()" must not be a valid arc label.
type ArcDomain[T trie.Elem] interface {
	//Code converts an arc label to its corresponding code.
	//It will panic if the label is not valid.
	Code(s T) int

	//Get converts a code to its corresponding arc label .
	Get(i int) T

	//StopElement returns a special element of type trie.Elem that indicates the end of a word.
	//The returned element must not be a valid arc label.
	StopElement() T
}

// DefaultArcDomain is a default implementation of ArcDomain.
// It is a slice of English alphabet from "a" to "z" and numbers from "0" to "9", with each element's position in the slice serving as its code
type DefaultArcDomain[T trie.Elem] []T

// Code converts an arc label to its corresponding code. Panic if the label is not valid.
func (d DefaultArcDomain[T]) Code(s T) int {
	for i, c := range d {
		if c == s {
			return i
		}
	}

	panic("invalid arc label")
}

// Get converts a code to its corresponding arc label.
func (d DefaultArcDomain[T]) Get(i int) T {
	return d[i]
}

// StopElement returns "#", which is a special element that indicates the end of a word.
func (d DefaultArcDomain[T]) StopElement() T {
	return d[0]
}

var (
	alphabetRune DefaultArcDomain[rune] = []rune("#abcdefghijklmnopqrstuvwxyz1234567890")
	alphabetByte DefaultArcDomain[byte] = []byte("#abcdefghijklmnopqrstuvwxyz1234567890")
)
