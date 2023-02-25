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

// leafIdx : t.leaf[leafIdx] is the newly added node, leafIdx = len(t.leaf)
// j : t.internal[j] is the match that i is going to against
// h : height of the match
func (t *TTree[T]) matchDown(leafIdx, internalIdx, h int) int {

	if internalIdx > t.lastInternalIdx() {
		return t.matchDown(leafIdx, internalIdx-(1<<(h-2)), h-1)
	}

	//fmt.Println("j , internal: ", j, len(t.internal))

	if internalIdx < t.lastInternalIdx() && h > 1 {
		t.internal[internalIdx] = t.matchDown(leafIdx, internalIdx+(1<<(h-2)), h-1)
		return t.internal[internalIdx]
	}

	var r, l int

	// j == t.lastInternalIdx()
	if h == 1 {

		l = internalIdx
		r = internalIdx + 1

		if r > leafIdx {
			r = leafIdx
		}

		if t.less(t.leaf[l], t.leaf[r]) {
			t.internal[internalIdx] = l
			return l
		}

		t.internal[internalIdx] = r
		return r

	}

	l = internalIdx - (1 << (h - 2))
	r = internalIdx + (1 << (h - 2))

	if r > leafIdx {
		r = leafIdx
	}

	if t.less(t.leaf[l], t.leaf[r]) {
		t.internal[internalIdx] = l
		return l
	}
	t.internal[internalIdx] = r
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

// matchUp matches internal[i] with internal[j]. i is the replaced node
func (t *TTree[T]) matchUp(i, h int) {

	if i == t.root {
		return
	}

	if h > 0 {
		space := 1 << h    // space between i and j, 2^h
		iH := i<<(h-1) + 1 // index of i at level h

		if iH%2 != 0 { // i is the left child
			j := i + space
			parent := (i + j) / 2 // parent of i and j

			if t.less(t.leaf[t.internal[i]], t.leaf[t.internal[j]]) {
				t.internal[parent] = i

				t.matchUp(parent, h+1)
				return
			} else {

				if !t.less(t.leaf[t.internal[j]], t.leaf[t.internal[parent]]) && !t.less(t.leaf[t.internal[parent]], t.leaf[t.internal[j]]) {
					return
				}

				t.internal[parent] = j
				t.matchUp(parent, h+1)
				return
			}

		} else { // i is the right child
			j := i - space
			parent := (i + j) / 2 // parent of i and j

			if t.less(t.leaf[t.internal[j]], t.leaf[t.internal[i]]) {
				t.internal[parent] = i

				t.matchUp(parent, h+1)
				return
			} else {

				if !t.less(t.leaf[t.internal[j]], t.leaf[t.internal[parent]]) && !t.less(t.leaf[t.internal[parent]], t.leaf[t.internal[j]]) {
					return
				}

				t.internal[parent] = j
				t.matchUp(parent, h+1)
				return

			}
		}
	}

}

func (t *TTree[T]) deleteMin() {
	//delte t.root
	//replace t.root with t.lastLeafIdx()
	//matchUp t.root

	var zero T

	if t.leafCount() == 1 {
		t.root = 0
		t.h = 0
		t.maxLeaf = 0
		t.leaf = t.leaf[:1]
		t.internal = t.internal[:1]
		return
	}

	last := t.lastLeafIdx()

	t.leaf[t.root] = t.leaf[last]

	t.leaf[t.lastLeafIdx()] = zero
	t.leaf = t.leaf[:t.lastLeafIdx()]
	t.internal[t.lastInternalIdx()] = 0
	t.internal = t.internal[:t.lastInternalIdx()]

	var l, r, p int

	if t.root%2 == 0 { // root is the right child
		r = t.root
		r = t.root - 1
	} else {
		l = t.root
		if t.root+1 >= len(t.leaf) {
			return
		}
		r = t.root + 1
	}

	p = t.internal[t.root]

	if t.less(t.leaf[t.root], t.leaf[t.root+1]) {
		t.internal[p] = t.root
		t.matchUp(p, 1)
	} else {
		t.internal[p] = t.root + 1
		t.matchUp(p, 1)
	}

	t.matchUp(t.root, 1)

	t.matchUp(t.lastLeafIdx(), 1)

}
