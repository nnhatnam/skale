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

func From[T any](less skale.LessFunc[T], e ...[]T) *Heap[T] {
	h := New[T](less)
	h.e = make([]T, 0, len(e))

	for i := len(e)/2 - 1; i >= 0; i-- {
		h.down(i, len(e)-1)
	}
	return h
}

func FromOrdered[T skale.Ordered](e ...[]T) *Heap[T] {
	return From[T](skale.Less[T](), e...)
}

// heapify
func (h *Heap[T]) Init(e []T) {
	for i := len(e)/2 - 1; i >= 0; i-- {
		h.down(i, len(e)-1)
	}
}

func (h *Heap[T]) up(i int) {
	for {

		if i == 0 {
			break
		}

		p := (i - 1) >> 1

		if !h.less(h.e[i], h.e[p]) {
			break
		}

		h.e[i], h.e[p] = h.e[p], h.e[i]
		i = p

	}
}

func (h *Heap[T]) down(i int, last int) {
	//current: e[i]
	//left: e[2i+1]
	//right: e[2i+2]
	//parent: e[(i-1)/2]

	for {
		l := i<<1 + 1

		if l > last || l < 0 {
			break
		}

		r := l + 1
		j := i

		if h.less(h.e[l], h.e[j]) {
			j = l
		}

		if r <= last && r > 0 && h.less(h.e[r], h.e[j]) {
			j = r
		}

		if j == i {
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
	if len(h.e) == 0 {
		return
	}
	var zero T

	n := len(h.e)
	v := h.e[0]
	h.e[0] = h.e[n-1]

	h.e[n-1] = zero //clear the last element
	h.e = h.e[:n-1]
	h.down(0, len(h.e)-1)
	return v, true
}

func (h *Heap[T]) Remove(i int) (_ T, _ bool) {
	//var zero T
	n := len(h.e) - 1

	if n != i {
		h.e[i], h.e[n] = h.e[n], h.e[i]
		h.down(i, n)
	}
	return h.Pop()
}

func (h *Heap[T]) Len() int {
	return len(h.e)
}
