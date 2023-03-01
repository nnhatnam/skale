package skew

import "github.com/nnhatnam/skale"

type Node[T any] struct {
	left, right *Node[T]
	Value       T
}

func NewNode[T any](v T) *Node[T] {
	return &Node[T]{Value: v}
}

type SkewHeap[T any] struct {
	root *Node[T]

	less skale.LessFunc[T]
}

func New[T any]() *SkewHeap[T] {
	return &SkewHeap[T]{}
}

func (h *SkewHeap[T]) merge(a, b *Node[T]) *Node[T] {
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

	return a
}

func (h *SkewHeap[T]) insert(v T) {
	if h.root == nil {
		h.root = NewNode(v)
		return
	}

	h.root = h.merge(h.root, NewNode(v))

	h.root.left, h.root.right = h.root.right, h.root.left
}

func (h *SkewHeap[T]) findMin() T {
	var zero T
	if h.root == nil {
		return zero
	}

	return h.root.Value
}

func (h *SkewHeap[T]) isEmpty() bool {
	return h.root == nil
}

func (h *SkewHeap[T]) deleteMin() T {
	var zero T
	if h.root == nil {
		return zero
	}

	v := h.root.Value
	h.root = h.merge(h.root.left, h.root.right)

	return v
}