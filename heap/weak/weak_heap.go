package weak

import (
	"github.com/nnhatnam/skale"
)

type WeakHeap[T any] struct {
	e []T
	r []int

	less skale.LessFunc[T]
}

func New[T any](less skale.LessFunc[T], size ...int) *WeakHeap[T] {
	h := &WeakHeap[T]{
		less: less,
	}
	if len(size) > 0 {
		h.e = make([]T, 0, size[0])
		h.r = make([]int, 0, size[0])
	}
	return h
}

func From[T any](less skale.LessFunc[T], e []T) *WeakHeap[T] {
	h := New[T](less)
	h.r = make([]int, len(e))
	h.e = make([]T, len(e))
	copy(h.e, e)

	for j := len(e) - 1; j >= 0; j-- {
		i := h.dAncestor(j)
		h.join(i, j)
	}
	return h
}

func (h *WeakHeap[T]) up(j int) {

	for j != 0 {
		i := h.dAncestor(j)

		if h.join(i, j) {
			break
		}
		j = i
	}
	//for j > 0 {
	//
	//	p := (j - 1) >> 1
	//
	//	if h.join(j, p) {
	//		break
	//	}
	//
	//	j = p
	//
	//}
}

func (h *WeakHeap[T]) down(j, n int) {

	k := 2*j + 1 - h.r[j]
	for 2*k+h.r[k] < n {
		k = 2*k + h.r[k]
	}

	for k != j {
		h.join(j, k)
		k = k / 2
	}
}

func (h *WeakHeap[T]) insert(v T) {
	h.e = append(h.e, v)
	h.r = append(h.r, 0)
	n := len(h.e) - 1
	for n&1 == 0 {
		h.r[n/2] = 0
	}
	h.up(n)

}

func (h *WeakHeap[T]) dAncestor(j int) int {

	for j&1 == h.r[j/2] {
		j = j / 2
	}

	return j / 2
}

// join conceptually combines two weak heaps into one weak heap
func (h *WeakHeap[T]) join(i, j int) bool {

	if h.less(h.e[i], h.e[j]) {
		h.e[i], h.e[j] = h.e[j], h.e[i]
		h.r[i] = 1 - h.r[i]
		return false
	}
	return true
}

//required: len(h.e) > 0
func (h *WeakHeap[T]) deleteMin() (_ T, _ bool) {

	v := h.e[0]
	n := len(h.e) - 1

	h.e[0] = h.e[n]

	if n > 1 {
		h.down(0, n)
	}
	return v, false
}

//func (h *WeakHeap[T]) insertBulk(e []T)  {
//	n := len(h.e)
//	k := len(e)
//	var left, right int
//	right = n + k - 2
//	left = n
//	if n < k/2 {
//		left = k/2
//	}
//
//	for k > 0 {
//		k--
//		h.e[n] = e[k]
//		h.r[n] = 0
//		n++
//	}
//
//	for right > left + 1 {
//		left  = left  / 2
//		right = right / 2
//		for i := left; i <= right; i++ {
//			h.down(i, n)
//		}
//
//	}
//
//
//}
