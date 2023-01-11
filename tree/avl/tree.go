package avl

import (
	"fmt"
	"github.com/nnhatnam/skale"
)

type Iterator[T any] func(item T) bool

type Node[T any] struct {
	Value T // value stored in the node
	Left  *Node[T]
	Right *Node[T]

	bFactor int8
}

func (n *Node[T]) String() string {
	return fmt.Sprintf("%v (%v)", n.Value, n.bFactor)
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
	r.Left = n
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

// fixOnInsert fix the tree at node n after insertion
// require: n != nil
func insertFix[T any](n *Node[T]) *Node[T] {

	if n.bFactor > 1 {
		n.bFactor = 0
		if n.Left.bFactor > 0 {
			n = rotateRight(n)
		} else { // n.Left.bFactor < 0 , there is no case n.Left.bFactor == 0
			n.Left.bFactor = 0
			n = rotateLeftRight(n)
		}
		n.bFactor = 0

	} else if n.bFactor < 1 {
		n.bFactor = 0
		if n.Right.bFactor < 0 {
			n = rotateLeft(n)
		} else { // n.Right.bFactor > 0 , there is no case n.Right.bFactor == 0
			n.Right.bFactor = 0
			n = rotateRightLeft(n)
		}
		n.bFactor = 0

	}

	return n
}

func deleteFix[T any](n *Node[T]) (_ *Node[T], hChanged int8) {
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

	return n, hChanged
}

// min returns the minimum node in the tree with root n
func min[T any](n *Node[T]) *Node[T] {
	for n.Left != nil {
		n = n.Left
	}

	return n
}

// max return the maximum node in the tree with root n
func max[T any](n *Node[T]) *Node[T] {
	for n.Right != nil {
		n = n.Right
	}

	return n
}

// find the node with value val in tree t, return the node
func (t *AVL[T]) find(val T) (_ *Node[T]) {
	curr := t.root
	for curr != nil {
		if t.less(val, curr.Value) {
			curr = curr.Left
		} else if t.less(curr.Value, val) {
			curr = curr.Right
		} else {
			return curr
		}
	}
	return nil
}

func (t *AVL[T]) insertNoReplace(n *Node[T], val T) *Node[T] {

	if n == nil {
		return NewNode(val)
	}

	if t.less(val, n.Value) {
		l := n.Left
		n.Left = t.insertNoReplace(n.Left, val)
		if n.Left.bFactor != 0 || l == nil {
			n.bFactor++
		}

	} else {
		r := n.Right
		n.Right = t.insertNoReplace(n.Right, val)

		if n.Right.bFactor != 0 || r == nil {
			n.bFactor--
		}

	}

	if n.bFactor < -1 || n.bFactor > 1 {

		n = insertFix(n)

	}

	return n

}

func (t *AVL[T]) replaceOrInsert(n *Node[T], val T) (_ *Node[T], _ T, insertCount int8) {

	if n == nil {
		return NewNode(val), *new(T), 1
	}

	var oldVal T

	if t.less(val, n.Value) {
		l := n.Left
		n.Left, oldVal, insertCount = t.replaceOrInsert(n.Left, val)
		if n.Left.bFactor != 0 || l == nil {
			n.bFactor += insertCount
		}

	} else if t.less(n.Value, val) {
		r := n.Right
		n.Right, oldVal, insertCount = t.replaceOrInsert(n.Right, val)
		if n.Right.bFactor != 0 || r == nil {
			n.bFactor -= insertCount
		}

	} else {
		oldVal = n.Value
		n.Value = val
		return n, oldVal, 0
	}

	if n.bFactor > 1 || n.bFactor < -1 {

		n = insertFix(n)

	}

	return n, oldVal, insertCount

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

	return n, hChanged, found
}

func (t *AVL[T]) print() {
	str := "AVLTree\n"
	output(t.root, "", true, &str)
	fmt.Println(str)
}

// -----------------------------------------public methods-----------------------------------------

// Len returns the number of nodes in the tree
func (t *AVL[T]) Len() int {
	return t.count
}

// Get returns the value of the node with value val
func (t *AVL[T]) Get(val T) (_ T, _ bool) {
	n := t.find(val)
	if n != nil {
		return n.Value, true
	}
	return
}

// Has returns true if the tree contains a node with value val
func (t *AVL[T]) Has(val T) bool {
	_, ok := t.Get(val)
	return ok
}

func (t *AVL[T]) InsertNoReplace(val T) {
	t.count++
	t.root = t.insertNoReplace(t.root, val)
}

func (t *AVL[T]) ReplaceOrInsert(val T) (T, bool) {
	var oldVal T
	var insertCount int8
	t.root, oldVal, insertCount = t.replaceOrInsert(t.root, val)

	if insertCount > 0 {
		t.count++
	}
	return oldVal, insertCount == 0
}

func (t *AVL[T]) Delete(val T) (T, bool) {
	var oldVal T
	var found bool
	t.root, _, found = t.delete(t.root, val)
	if found {
		t.count--
	}
	return oldVal, found
}

func output[T any](node *Node[T], prefix string, isTail bool, str *string) {
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
