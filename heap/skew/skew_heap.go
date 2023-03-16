package skew

import "github.com/nnhatnam/skale"

type Node[T any] struct {
	//left, right, up *Node[T]

	up, down *Node[T]
	value    T
}

func NewNode[T any](v T) *Node[T] {
	return &Node[T]{value: v}
}

type SkewHeap[T any] struct {
	root *Node[T]

	less skale.LessFunc[T]
}

func New[T any]() *SkewHeap[T] {
	return &SkewHeap[T]{}
}

func (h *SkewHeap[T]) meld(h1, h2 *Node[T]) *Node[T] {
	//base case
	if h1 == nil {
		return h2
	}

	if h2 == nil {
		return h1
	}

	if h.less(h1.up.value, h2.up.value) {
		h1, h2 = h2, h1
	}

	//remove from h1 its bottom right most node
	h3 := h1.up   // h3 is the bottom right most node of h1
	h1.up = h3.up // h1.next is the parent node of h3
	h3.up = h3

	var x *Node[T]

	for h1 != h3 {

		// h1 will be the heap with the bigger right most node
		if h.less(h1.up.value, h2.up.value) {
			h1, h2 = h2, h1
		}

		// remove the bottom right most node of h1 and store it in x
		x = h1.up
		h1.up = x.up

		// add x to the top of h3 and swap its children
		x.up = x.down // swap left and right child
		x.down = h3.up
		h3.up = x

		h3 = h3.up // h3 = x

	}

	// attach h3 to the right of h2
	h2.up, h3.up = h3.up, h2.up

	return h2
}

func (h *SkewHeap[T]) insert(v T) {

	if h.root == nil {
		h.root = NewNode(v)
		h.root.up = h.root
		h.root.down = h.root
		return
	}

	x := NewNode(v)

	// v <= h.root.value => x will be the new root of the heap, and h.root become the left child of x
	if !h.less(h.root.value, v) {

		x.down = h.root.up // x.down is the bottom right most node of the left subtree of x
		h.root.up = x
		x.up = x
		return
	}

	x.up = h.root.up
	x.down = x
	h.root.up = x

	// v > h.root.value => x becomes the new lowest node on the major right path of the heap.
	// In the general case, x will be inserted somewhere in the middle of the major right path.

	y, z := h.root.up, h.root.up // y is the node before x, z is the node after x
	for h.less(x.value, y.up.value) {

		y = y.up              // move y up
		z, y.down = y.down, z // swap left and right child of y, and make z points to the lowest node on the major right path of y

	}

	// put y to the left of x
	x.up = y.up
	x.down = z
	y.up = x
	h.root.up = x // update the root of the heap

}

func (h *SkewHeap[T]) findMin() T {
	var zero T
	if h.root == nil {
		return zero
	}

	return h.root.value
}

func (h *SkewHeap[T]) isEmpty() bool {
	return h.root == nil
}

// @require: h.root != nil
func (h *SkewHeap[T]) deleteMin() T {
	var x, y1, y2, h3 *Node[T]

	y1 = h.root.up
	y2 = h.root.down

	if y1 == y2 {
		h.root = nil
		return y1.value
	}

	// make sure y is farther from the root than y2 in terms of value ordering
	// if y1.value == y2.value => we can't determine which one is farther from the root.
	// In that case, we will check
	if h.less(y1.value, y2.value) {
		y1, y2 = y2, y1
	}

	h3 = y1
	y1 = y1.up
	h3.up = h3

	for y1.up != y2.up {

		if h.less(y1.value, y2.value) { // swap if y2.value > y1.value, so that y1.value <= y2.value
			y1, y2 = y2, y1
		}

		if y1 == h.root {

			for {
				if y2 == h.root {
					h.root = h3
					return y2.value
				}
				x = y2
				y2 = y2.up

				x.up = x.down
				x.down = h3
				h3.up = x
				h3 = h3.up
			}

		}

		x = y1
		x.up = x.down
		x.down = h3.up
		h3.up = x
		h3 = h3.up

	}

	h.root = h3
	return y1.value
}

// Push inserts a value into the heap.
// Panics if `h` is nil.
func (h *SkewHeap[T]) Push(v T) {
	h.insert(v)
}
