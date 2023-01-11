package avl

import (
	"fmt"
	"github.com/nnhatnam/skale"
)

type Node[T any] struct {
	Value T // value stored in the node
	Left  *Node[T]
	Right *Node[T]

	bFactor int8
}

func (n *Node[T]) String() string {
	return fmt.Sprintf("%v", n.Value)
}

func NewNode[T any](value T) *Node[T] {
	return &Node[T]{Value: value}
}

type AVL[T any] struct {
	root  *Node[T]
	count int
	less  skale.LessFunc[T]
}

func New[T any](less skale.LessFunc[T]) *AVL[T] {
	return &AVL[T]{less: less}
}

func NewOrdered[T skale.Ordered]() *AVL[T] {
	return New(skale.Less[T]())
}

// -----------------------------------------internal methods-----------------------------------------

// rotate left = everything move left
// rotate left brings the right child to the root, and the root to the left child. Return the new root
func rotateLeft[T any](n *Node[T]) *Node[T] {
	r := n.Right
	n.Right = r.Left
	n.Left = r
	return r

}

func rotateRight[T any](n *Node[T]) *Node[T] {
	l := n.Left
	n.Left = l.Right
	l.Right = n
	return l

}

func rotateLeftRight[T any](n *Node[T]) *Node[T] {
	n.Left = rotateLeft(n.Left)
	return rotateRight(n)
}

func rotateRightLeft[T any](n *Node[T]) *Node[T] {
	n.Right = rotateRight(n.Right)
	return rotateLeft(n)
}

func (t *AVL[T]) insertNoReplace(n *Node[T], val T) (_ *Node[T], hChanged int8) {

	if n == nil {
		return NewNode(val), 1
	}

	if t.less(val, n.Value) {
		n.Left, hChanged = t.insertNoReplace(n.Left, val)
		n.bFactor += hChanged
	} else {
		n.Right, hChanged = t.insertNoReplace(n.Right, val)
		n.bFactor -= hChanged
	}

	if n.bFactor > 1 {
		if n.Left.bFactor > 0 {
			n.bFactor = 0
			n = rotateRight(n)
			n.bFactor = 0
		} else { // n.Left.bFactor < 0 , there is no case n.Left.bFactor == 0
			n.bFactor = 0
			n.Left.bFactor = 0
			n = rotateLeftRight(n)
			n.bFactor = 0
		}
		hChanged = 0
	} else if n.bFactor < 1 {
		if n.Right.bFactor < 0 {
			n.bFactor = 0
			n = rotateLeft(n)
			n.bFactor = 0
		} else { // n.Right.bFactor > 0 , there is no case n.Right.bFactor == 0
			n.bFactor = 0
			n.Right.bFactor = 0
			n = rotateRightLeft(n)
			n.bFactor = 0
		}
		hChanged = 0
	}

	return n, hChanged

}

func min[T any](n *Node[T]) *Node[T] {
	for n.Left != nil {
		n = n.Left
	}

	return n
}

func max[T any](n *Node[T]) *Node[T] {
	for n.Right != nil {
		n = n.Right
	}

	return n
}

func (t *AVL[T]) delete(n *Node[T], val T) (_ *Node[T], hChanged int8, found bool) {
	if n == nil {
		return nil, 0, false // not found
	}

	if t.less(val, n.Value) {
		n.Left, hChanged, found = t.delete(n.Left, val)
		n.bFactor -= hChanged
	} else if t.less(n.Value, val) {
		n.Right, hChanged, found = t.delete(n.Right, val)
		n.bFactor += hChanged
	} else {
		found = true
		if n.Left == nil && n.Right == nil {
			return nil, 1, found // found and deleted
		} else if n.Right == nil {
			return n.Left, 1, found
		} else if n.Left == nil {
			return n.Right, 1, found
		} else {
			succ := min(n.Right)
			n.Value = succ.Value
			n.Right, hChanged, _ = t.delete(n.Right, succ.Value)
		}

	}

	if n.bFactor > 1 {
		if n.Left.bFactor == 0 {
			//R0 rotation
			n.bFactor = 1
			n = rotateRight(n)
			n.bFactor = -1

		} else if n.Left.bFactor == 1 {
			//R1 rotation
			n.bFactor = 0
			n = rotateRight(n)
			n.bFactor = 0
		} else {
			//R-1 rotation
			n.bFactor = 0
			n.Left.bFactor = 0
			n = rotateLeftRight(n)
			n.bFactor = 0
			hChanged = 1
		}

	} else if n.bFactor < -1 {
		if n.Right.bFactor == 0 {
			//L0 rotation
			n.bFactor = -1
			n = rotateLeft(n)
			n.bFactor = 1
		} else if n.Right.bFactor == -1 {
			//L-1 rotation
			n.bFactor = 0
			n = rotateLeft(n)
			n.bFactor = 0
		} else {
			//L1 rotation
			n.bFactor = 0
			n.Right.bFactor = 0
			n = rotateRightLeft(n)
			n.bFactor = 0
			hChanged = 1
		}
	}

	return n, hChanged, found
}

// -----------------------------------------public methods-----------------------------------------

func (t *AVL[V]) InsertNoReplace(val V) {
	t.insertNoReplace(t.root, val)
}

func output[V any](node *Node[V], prefix string, isTail bool, str *string) {
	if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.Right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.Left, newPrefix, true, str)
	}
}
