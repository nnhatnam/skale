package heap

import (
	"github.com/nnhatnam/skale"
)

type Heap[T any] struct {
	less skale.LessFunc[T]
}

func New[T any](less skale.LessFunc[T]) *Heap[T] {
	return &Heap[T]{less: less}
}

func NewOrdered[T skale.Ordered]() *Heap[T] {
	return &Heap[T]{less: skale.Less[T]()}
}

// heapify
func (h *Heap[T]) Init(e []T) {
	for i := len(e)/2 - 1; i >= 0; i-- {
		h.down(e, i, len(e)-1)
	}
}

func (h *Heap[T]) up(e []T, i int) {
	for {
		if i == 0 {
			break
		}

		p := (i - 1) >> 1

		if !h.less(e[i], e[p]) {
			break
		}

		e[i], e[p] = e[p], e[i]
		i = p
	}
}

func (h *Heap[T]) down(e []T, i int, last int) {
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

		if h.less(e[l], e[j]) {
			j = l
		}

		if r <= last && r > 0 && h.less(e[r], e[j]) {
			j = r
		}

		if j == i {
			break
		}

		e[i], e[j] = e[j], e[i]
		i = j

	}

}
