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
	n := len(e)
	h := New[T](less)
	h.e = make([]T, n)

	copy(h.e, e)

	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
	return h
}

func FromOrdered[T skale.Ordered](e []T) *Heap[T] {
	return From[T](skale.Less[T](), e)
}

// heapify
func (h *Heap[T]) Init(e []T) {
	for i := len(e)/2 - 1; i >= 0; i-- {
		h.down(i, len(e))
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
// i : the index of the node to be sifted down
// n : the length of the heap
func (h *Heap[T]) down(i int, n int) {
	//current: e[i]
	//left: e[2i+1]
	//right: e[2i+2]
	//parent: e[(i-1)/2]

	for {

		j := i<<1 + 1 // i : parent, j : left child

		if j >= n || j < 0 { // j < 0 is to avoid overflow
			break
		}

		if r := j + 1; r < n && h.less(h.e[r], h.e[j]) { // r won't overflow, because j < n => j + 1 <= n
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

	var zero T

	n := len(h.e) - 1

	v := h.e[0]

	//delete
	h.e[0] = h.e[n]
	h.e[n] = zero //clear the last element
	h.e = h.e[:n]

	h.down(0, n)
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
	h.down(i, n)

	return v, true
}

func (h *Heap[T]) Len() int {
	return len(h.e)
}
