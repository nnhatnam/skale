package mptrie

import "github.com/nnhatnam/skale/trie"

type ArcLabels[T trie.Elem] interface {
	Code(s T) int
	Get(i int) T
	StopElement() T
	Size() int
}

type defaultArc[T trie.Elem] []T

func (d defaultArc[T]) Code(s T) int {
	for i, c := range d {
		if c == s {
			return i
		}
	}
	return -1
}

func (d defaultArc[T]) Get(i int) T {
	return d[i]
}

func (d defaultArc[T]) StopElement() T {
	return d[0]
}

func (d defaultArc[T]) Size() int {
	return len(d) - 1 //exclude the stop element
}

var (
	alphabet defaultArc[rune] = []rune("#abcdefghijklmnopqrstuvwxyz1234567890")
)
