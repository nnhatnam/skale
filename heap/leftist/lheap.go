package leftist

import "github.com/nnhatnam/skale"

type Node[T any] struct {
	left, right *Node[T]
	Value       T

	// s is the distance between the node and the nearest leaf in the subtree of the node.
	npl int
}

func NewNode[T any](v T) *Node[T] {
	return &Node[T]{Value: v}
}

type LHeap[T any] struct {
	root *Node[T]

	less skale.LessFunc[T]
}

func New[T any](less skale.LessFunc[T]) *LHeap[T] {
	return &LHeap[T]{less: less}
}

func (h *LHeap[T]) merge(a, b *Node[T]) *Node[T] {
	//base case
	if a == nil {
		return b
	}

	if b == nil {
		return a
	}

	//make `a` a heap with smaller root value
	if h.less(b.Value, a.Value) {
		a, b = b, a
	}

	a.right = h.merge(a.right, b)

	//swap left and right subtree if needed
	if a.left == nil || a.left.npl < a.right.npl {
		a.left, a.right = a.right, a.left
	}

	//update npl
	a.npl = b.right.npl + 1

	return nil
}

func (h *LHeap[T]) insert(v T) {
	if h.root == nil {
		h.root = NewNode(v)
		return
	}

	h.root = h.merge(h.root, NewNode(v))

}

func (h *LHeap[T]) isEmpty() bool {
	return h.root == nil
}

func (h *LHeap[T]) findMin() T {
	var zero T
	if h.root == nil {
		return zero
	}

	return h.root.Value
}

func (h *LHeap[T]) deleteMin() T {
	var zero T
	if h.root == nil {
		return zero
	}

	min := h.root.Value
	h.root = h.merge(h.root.left, h.root.right)
	return min
}
