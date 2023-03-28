package cartesian

import (
	"github.com/nnhatnam/skale"
	"github.com/nnhatnam/skale/list/stack"
)

type Node struct {
	left, right *Node
	idx         int
}

func NewNode(idx int) *Node {
	return &Node{idx: idx}
}

type Cartesian[T any] struct {
	s *[]T

	root *Node
	less skale.LessFunc[T]
}

func New[T any](less skale.LessFunc[T]) *Cartesian[T] {
	return &Cartesian[T]{less: less}
}

func NewOrdered[T skale.Ordered]() *Cartesian[T] {
	t := New[T](skale.Less[T]())
	return t
}

func (t *Cartesian[T]) Init(e []T) {

	if len(e) == 0 {
		arr := make([]T, 0)
		t.s = &arr // assign to empty array to avoid nil and also clear t.s in case of re-init with empty array
		return
	}

	s := stack.NewStackSWithSize[*Node](len(e))
	t.s = &e

	for i := 0; i < len(e); i++ {
		n := NewNode(i)

		for !s.Empty() {
			top, _ := s.Top()
			if t.less(e[top.idx], e[i]) {
				// value at top < e[i]
				break
			}

			// value at top >= e[i]
			n.left = top
			s.Pop()
		}

		if !s.Empty() {
			top, _ := s.Top()
			top.right = n
		}
		s.Push(n)
	}

	t.root, _ = s.Bottom()

}

// merge merges two cartesian trees t and other into one tree.
// The resulting tree will contain all elements from both trees.
// Merge will use less function from t to determine the order of elements.
// Merge is destructive, t will contain all elements from both trees and
// other will be empty after the merge.
func merge[T any](t *Cartesian[T], other *Cartesian[T]) *Cartesian[T] {

	*t.s = append(*t.s, *other.s...)

	l := t.root
	r := other.root

	var subMerge func(l, r *Node) *Node

	subMerge = func(l, r *Node) *Node {
		// https://habr.com/en/post/101818/
		if l == nil {
			return r
		}

		if r == nil {
			return l
		}

		if t.less((*t.s)[l.idx], (*other.s)[r.idx]) {
			l.right = subMerge(l.right, r)
			return l
		}

		r.left = subMerge(l, r.left)
		return r
	}

	newT := New[T](t.less)
	*newT.s = make([]T, 0, len(*t.s)+len(*other.s))
	copy(*newT.s, *t.s)
	copy((*newT.s)[len(*t.s):], *other.s)

	newT.root = subMerge(l, r)

	t.s = nil
	t.root = nil
	other.s = nil
	other.root = nil
	return newT
}

// split splits the cartesian tree t into two trees l and r such that all keys in l are less than x and all keys in r are greater than or equal to x.
func split(t *Node, x int) (l *Node, r *Node) {
	// https://habr.com/en/post/101818/
	// http://e-maxx.ru/algo/treap

	if t == nil {
		return nil, nil
	}

	if t.idx < x {
		l, r = split(t.right, x)
		t.right = l
		return t, r
	}

	l, r = split(t.left, x)
	t.left = r
	return l, t
}

func (t *Cartesian[T]) Get(i int) T {
	return (*t.s)[i]
}
