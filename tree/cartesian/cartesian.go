package cartesian

import (
	"github.com/nnhatnam/skale"
	"github.com/nnhatnam/skale/list/stack"
)

type Node struct {
	left, right *Node
	value       int
}

func NewNode[T any](idx int) *Node {
	return &Node{value: idx}
}

type CartesianTree[T any] struct {
	root *Node
	less skale.LessFunc[T]
}

func New[T any](less skale.LessFunc[T]) *CartesianTree[T] {
	return &CartesianTree[T]{less: less}
}

func (t *CartesianTree[T]) Init(e []T) {
	s := stack.NewStackSWithSize[*Node[T]](len(e))

	for i := 0; i < len(e); i++ {

		//last := -1 // last node with value less than e[i]

		n := NewNode[T](i)

		for !s.Empty() {
			top, _ := s.Top()
			if !t.less(e[top.value], e[i]) {
				n.left = top
				//last = top.value
				s.Pop()
			} else {
				break
			}
		}

		if !s.Empty() {
			top, _ := s.Top()
			top.right = n
		}
		s.Push(n)
	}

}

//func merge[T any](t1, t2 *CartesianTree[T]) *CartesianTree[T] {
//
//}

func (t *CartesianTree[T]) Root() *Node {
	return t.root
}

func (t *CartesianTree[T]) merge(t1 *CartesianTree[T]) {

}
