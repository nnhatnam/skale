package tournament

import (
	"github.com/nnhatnam/skale"
)

type WinnerTree[T any] struct {
	e    []int // index of the winner of each level
	less skale.LessFunc[T]
}

func NewWT[T any](less skale.LessFunc[T]) *WinnerTree[T] {
	return &WinnerTree[T]{less: less}
}

func NewOrderedWT[T skale.Ordered]() *WinnerTree[T] {
	return &WinnerTree[T]{less: skale.Less[T]()}
}

func (t *WinnerTree[T]) Init(v []T) {

	if len(v) < 2 {
		return
	}

	//a winner tree has 2n+1 nodes (n internal nodes, n + 1 leaves)
	n := len(v) - 1 // number of internal nodes (also last leaf index)
	t.e = make([]int, n)

	// in 1 array, internal nodes are at [0..n-1], leaves are at [n..2n]
	// in 2 arrays, internal nodes are at [0..n-1], leaves are at [0..n]
	i := len(v) - 1

	unMatch := false // unMatch is a flag to indicate if the last match is unmatched
	for i >= 0 {     // equal to for i := 2n; i >= n; i--

		if i == 0 {
			unMatch = true
			break
		}

		l, r := i-1, i

		p := int(uint(l+n) >> 1) // avoid overflow
		t.e[p] = t.match(v, l, r)

		i -= 2

	}

	j := len(t.e) - 1

	if unMatch {
		// match i and j
		p := (j - 1) >> 1
		t.e[p] = t.match(v, i, t.e[j])
		j--
	}

	// [0........j] internal nodes
	for j > 0 {
		p := (j - 1) >> 1
		t.e[p] = t.match(v, t.e[j], t.e[j-1])
		j -= 2
	}

}

func (t *WinnerTree[T]) internalRematch(v []T, i int) {

	//TODO: boundary check

	for {

		if i == 0 {
			break
		}

		p := (i - 1) >> 1
		l := p<<1 + 1
		r := l + 1

		t.e[p] = t.match(v, t.e[l], t.e[r])

		i = p
	}

}

func (t *WinnerTree[T]) reMatch(v []T, i int) {

	if len(v) < 2 {
		return
	}

	n := len(v) - 1            // number of internal nodes (also last leaf index)
	p := int(uint(i+n-1)) >> 1 // avoid overflow

	m := len(t.e) - 1 // last internal node index

	if (m-1)>>1 == p>>1 {
		// last internal node is at the same level as the last leaf
		// so we need to match the last leaf with the last internal node
		t.e[m] = t.match(v, t.e[m], i)

		t.internalRematch(v, m)
		return

	}

	l := 2*p - n + 1
	r := l + 1

	t.e[p] = t.match(v, t.e[l], t.e[r])

	t.internalRematch(v, p)

}

func (t *WinnerTree[T]) Push(v []T, x T) []T {

	var zero T
	if len(t.e) == 0 {
		v = append(v, x)
		return v
	}

	v = append(v, v[0], x)
	v[0] = zero // avoid memory leak
	v = v[1:]
	t.e = append(t.e, 0)

	t.reMatch(v, len(v)-1)
	return v
}

func (t *WinnerTree[T]) Fix(v []T, i int) {
	t.reMatch(v, i)
}

func (t *WinnerTree[T]) Pop(v []T) T {
	return v[t.e[0]]
}

func (t *WinnerTree[T]) match(e []T, i, j int) int {
	if t.less(e[i], e[j]) {
		return i
	}
	return j
}
