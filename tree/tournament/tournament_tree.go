package tournament

import (
	"fmt"
	"github.com/nnhatnam/skale"
)

type TTree[T any] struct {
	internal []int // internal[0] will not be used

	leaf []T

	less skale.LessFunc[T]

	h int // height of the tree

	maxLeaf int // max number of leaf nodes in the tree = 2^h

	root int // index of the root node in the internal array
}

func New[T any](less skale.LessFunc[T]) *TTree[T] {

	return &TTree[T]{
		less: less,
	}
}

func NewWithOrdered[T skale.Ordered]() *TTree[T] {

	return &TTree[T]{
		less: skale.Less[T](),
	}
}

func (t *TTree[T]) init(n int) {
	t.leaf = make([]T, 1, n)
	t.internal = make([]int, 1, n)
	t.maxLeaf = 1
}

func (t *TTree[T]) lazyInit() {
	if len(t.leaf) == 0 {
		fmt.Println("lazyInit")
		t.init(1)
	}

}

func (t *TTree[T]) leafCount() int {
	return len(t.leaf) - 1
}

func (t *TTree[T]) lastLeafIdx() int {
	return len(t.leaf) - 1
}

func (t *TTree[T]) lastInternalIdx() int {
	return len(t.internal) - 1
}

func (t *TTree[T]) updateH(h int) {
	t.h = h
	t.maxLeaf = 1 << h
}

// i : t.leaf[i] is the newly added node, i = len(t.leaf)
// j : t.internal[j] is the match that i is going to against
// h : height of the match
func (t *TTree[T]) matchDown(i, j, h int) int {
	fmt.Println("matchDown: ", i, j, h, t.lastInternalIdx())

	if j > t.lastInternalIdx() {
		return t.matchDown(i, j-(1<<(h-2)), h-1)
	}

	//fmt.Println("j , internal: ", j, len(t.internal))

	if j < t.lastInternalIdx() && h > 1 {
		t.internal[j] = t.matchDown(i, j+(1<<(h-2)), h-1)
		return t.internal[j]
	}

	var r, l int

	// j == t.lastInternalIdx()
	if h == 1 {

		l = j
		r = j + 1

		if r > i {
			r = i
		}

		if t.less(t.leaf[l], t.leaf[r]) {
			t.internal[j] = l
			return l
		}

		t.internal[j] = r
		return r

	}

	l = j - (1 << (h - 2))
	r = j + (1 << (h - 2))

	if r > i {
		r = i
	}

	if t.less(t.leaf[l], t.leaf[r]) {
		t.internal[j] = l
		return l
	}
	t.internal[j] = r
	return r
}

func (t *TTree[T]) internalValue(i int) T {
	return t.leaf[t.internal[i]]
}

func (t *TTree[T]) print() {
	fmt.Println(t.internal[1:])
	fmt.Println(t.leaf[1:])
	fmt.Println("root:", t.root)
	fmt.Println("height:", t.h)
	fmt.Println("maxLeaf:", t.maxLeaf)
	fmt.Println("----------------------------------------------------")
}

func (t *TTree[T]) insert(v T) {

	t.lazyInit()

	// i := 1 => t.internal[i] = 1 and t.leaf[ 2*i - len(t.internal) ] = 1
	// i := 2 => t.internal[ i / 2 ] = 2 and t.leaf[ 2*i - len(t.internal) ] = 2

	t.leaf = append(t.leaf, v)

	if t.leafCount() == 1 {

		t.root = 1
		t.maxLeaf = 1
		t.h = 0
		return
	}

	if t.leafCount() > t.maxLeaf {

		if t.less(v, t.leaf[t.root]) {
			t.internal = append(t.internal, t.lastLeafIdx())

		} else {
			t.internal = append(t.internal, t.root)
		}

		t.h++
		t.root = t.maxLeaf // t.root = 2^(h - 1)
		t.maxLeaf <<= 1    // t.maxLeaf = 2^h

		return

	}

	t.internal = append(t.internal, 0)
	t.internal[t.root] = t.matchDown(t.lastLeafIdx(), t.root, t.h)

}

func (t *TTree[T]) matchUp(v T) {

}
