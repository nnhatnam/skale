package leftist

import (
	"github.com/nnhatnam/skale"
)

// Node is a node in a leftist heap.
type Node[T any] struct {
	left, right *Node[T]
	Value       T

	// npl is the distance between the node and the nearest leaf in the subtree of the node.
	npl int
}

// NewNode returns a new node with value v.
func NewNode[T any](v T) *Node[T] {
	return &Node[T]{Value: v}
}

// npl_ returns the npl of the node. if n is nil, return -1.
func (n *Node[T]) npl_() int {
	//*Node(nil) has npl = 0
	if n == nil {
		return -1
	}

	return n.npl
}

// LHeap represents a leftist heap.
type LHeap[T any] struct {
	root *Node[T]

	less skale.LessFunc[T]
	len  int
}

// len_ returns the number of items in the heap. if h is nil, return 0.
func (h *LHeap[T]) len_() int {
	if h == nil {
		return 0
	}
	return h.len
}

func New[T any](less skale.LessFunc[T]) *LHeap[T] {
	return &LHeap[T]{less: less}
}

func NewOrdered[T skale.Ordered]() *LHeap[T] {
	return &LHeap[T]{less: skale.Less[T]()}
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
	a.npl = a.right.npl_() + 1

	return a
}

func (h *LHeap[T]) insert(v T) {

	h.len++

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

func (h *LHeap[T]) deleteMin() (_ T, _ bool) {

	if h.root == nil {
		return
	}
	h.len--
	min := h.root.Value
	h.root = h.merge(h.root.left, h.root.right)
	return min, true
}

// Push inserts a value into the heap.
// Panics if `h` is nil.
// Time complexity: O(log(n))
func (h *LHeap[T]) Push(v T) {
	h.insert(v)
}

// PushBulk inserts multiple values into the heap.
// Panics if `h` is nil.
// Time complexity: O(n log(n))
func (h *LHeap[T]) PushBulk(vs ...T) {
	for _, v := range vs {
		h.insert(v)
	}
}

// IsEmpty returns true if the heap is empty.
// Panics if `h` is nil.
// Time complexity: O(1)
func (h *LHeap[T]) IsEmpty() bool {
	return h.isEmpty()
}

// Peek returns the minimum value in the heap according to h.less. If the heap is empty, it returns the zero value of the type.
// Panics if `h` is nil.
// Time complexity: O(1)
func (h *LHeap[T]) Peek() T {
	return h.findMin()
}

// Pop deletes the minimum value in the heap according to h.less and returns it. If the heap is empty, (zero value, false) is returned.
// Panics if `h` is nil.
// Time complexity: O(log n)
func (h *LHeap[T]) Pop() (_ T, _ bool) {
	return h.deleteMin()
}

// Merge merges heap `other` into `h`. Merge will use the less function of `h` to compare values.
// After merging, `other` will be empty.
// If you want to use less function from `other`, you should call `other.Merge(h)`.
// Panics if `h` is nil.
func (h *LHeap[T]) Merge(other *LHeap[T]) {

	h.len += other.len_() // panic if h is nil
	h.root = h.merge(h.root, other.root)

	other.root = nil
	other.len = 0
	other.less = nil
}

// Len returns the number of items in the heap.
func (h *LHeap[T]) Len() int {
	return h.len
}
