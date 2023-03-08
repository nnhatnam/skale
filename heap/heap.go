package heap

import (
	"github.com/nnhatnam/skale"
)

type Heap[T any] struct {
	e    []T
	less skale.LessFunc[T]
}

func New[T any](less skale.LessFunc[T], size ...int) *Heap[T] {
	h := &Heap[T]{less: less}
	if len(size) > 0 {
		h.e = make([]T, 0, size[0])
	}
	return h
}

func NewOrdered[T skale.Ordered](size ...int) *Heap[T] {
	return New[T](skale.Less[T](), size...)
}

func From[T any](less skale.LessFunc[T], e []T) *Heap[T] {
	h := New[T](less)
	h.e = make([]T, len(e))

	copy(h.e, e)

	for i := len(h.e)/2 - 1; i >= 0; i-- {
		h.down(i)
	}
	return h
}

func FromOrdered[T skale.Ordered](e []T) *Heap[T] {
	return From[T](skale.Less[T](), e)
}

// heapify
func (h *Heap[T]) Init(e []T) {
	for i := len(e)/2 - 1; i >= 0; i-- {
		h.down(i)
	}
}

func (h *Heap[T]) up(i int) {

	for i > 0 {

		p := (i - 1) >> 1

		if !h.less(h.e[i], h.e[p]) {
			break
		}

		h.e[i], h.e[p] = h.e[p], h.e[i]
		i = p

	}
}

// sift-down procedure, O(log n). `i` is the index of the node to be sifted down.
func (h *Heap[T]) down(i int) {
	//current: e[i]
	//left: e[2i+1]
	//right: e[2i+2]
	//parent: e[(i-1)/2]

	last := len(h.e) - 1
	p := (last - 1) >> 1 // last internal node

	for i <= p { // p is the last internal node, so it's always have at least one child.

		j := i<<1 + 1 // left won't overflow, but right might overflow
		r := j + 1

		if r <= last && r > 0 && h.less(h.e[r], h.e[j]) { // r > 0 is to avoid overflow
			j = r
		}

		if !h.less(h.e[j], h.e[i]) {
			break
		}

		h.e[i], h.e[j] = h.e[j], h.e[i]
		i = j

	}
}

func (h *Heap[T]) Push(v T) {
	h.e = append(h.e, v)
	h.up(len(h.e) - 1)
}

func (h *Heap[T]) Pop() (_ T, _ bool) {

	//var zero T

	n := len(h.e) - 1

	v := h.e[0]

	//delete
	h.e[0] = h.e[n]
	//h.e[n] = zero //clear the last element
	h.e = h.e[:n]

	h.down(0)
	return v, true
}

func (h *Heap[T]) Remove(i int) (_ T, _ bool) {

	if i == 0 {
		return h.Pop()
	}

	var zero T

	n := len(h.e) - 1
	v := h.e[i]
	h.e[i] = h.e[n]
	h.e[n] = zero //clear the last element

	h.e = h.e[:n]
	h.down(i)

	return v, true
}

func (h *Heap[T]) Len() int {
	return len(h.e)
}
