package tournament

import (
	"github.com/nnhatnam/skale"
)

type TTree[T any] struct {
	less skale.LessFunc[T]

	e []T
}

func New[T any](less skale.LessFunc[T]) *TTree[T] {

	return &TTree[T]{
		less: less,
	}
}

func NewOrdered[T skale.Ordered]() *TTree[T] {

	return New(skale.Less[T]())
}

func (t *TTree[T]) up(i int) {

	//TODO: boundary check

	for {

		if i == 0 {
			break
		}

		p := (i - 1) >> 1

		if t.less(t.e[i], t.e[p]) {
			t.e[p] = t.e[i]
		} else if t.less(t.e[p], t.e[i]) {
			l := p<<1 + 1
			r := l + 1

			if t.less(t.e[l], t.e[r]) {
				t.e[p] = t.e[l]
			} else {
				t.e[p] = t.e[r]
			}
		} else {
			break
		}

		i = p
	}

}

func (t *TTree[T]) insert(v T) {

	if len(t.e) == 0 {
		t.e = append(t.e, v)
		return
	}

	p := (len(t.e) - 1) >> 1

	t.e = append(t.e, t.e[p], v)

	t.up(len(t.e) - 1)

}

func (t *TTree[T]) deleteMin() (_ T, _ bool) {

	var zero T

	if len(t.e) == 0 {
		return
	}

	min := t.e[0]

	if len(t.e) == 1 {
		t.e[0] = zero // avoid memory leak
		t.e = t.e[:0]
		return min, true
	}

	i := 0
	last := len(t.e) - 1
	prevLast := last - 1

	var replaced T // min will be replaced by this node

	// set replaced node to which is bigger between last and prevLast
	if t.less(t.e[last], t.e[prevLast]) {
		replaced = t.e[prevLast]

	} else {
		replaced = t.e[last]
	}

	// find min node location
	for {

		r := i<<1 + 2
		if r > last {
			break
		}
		l := r - 1

		if t.less(t.e[i], t.e[r]) {
			//min is left child
			i = l
		} else {
			//min is right child
			i = r
		}

	}

	t.e[i] = replaced

	t.up(i)
	t.e[last] = zero     // avoid memory leak
	t.e[prevLast] = zero // avoid memory leak
	t.e = t.e[:prevLast]

	return min, true

}

func (t *TTree[T]) InsertNoReplace(v T) {

	t.insert(v)

}

// InsertNoReplaceBulk inserts multiple elements into the tree. It does not guarantee the order of the elements.
func (t *TTree[T]) InsertNoReplaceBulk(v ...T) {

	for _, e := range v {
		t.InsertNoReplace(e)
	}

}

func (t *TTree[T]) DeleteMin() (T, bool) {

	return t.deleteMin()

}
