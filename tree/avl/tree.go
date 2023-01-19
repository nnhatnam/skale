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

		if n.Left.bFactor > 0 {
			n.bFactor = 0
			n = rotateRight(n)

		} else {

			n.bFactor = -maxInt8(n.Left.Right.bFactor, 0) //n.bFactor = n.bFactor - 2 - maxInt8(n.Left.Right.bFactor, 0)
			n.Left.bFactor = n.Left.bFactor + 1 - minInt8(n.Left.Right.bFactor, 0)
			n = rotateLeftRight(n)
		}
		n.bFactor = 0

	} else if n.bFactor < 1 {

		if n.Right.bFactor < 0 {
			n.bFactor = 0
			n = rotateLeft(n)
		} else { // n.Right.bFactor == 1 , there is no case n.Right.bFactor == 0

			n.bFactor = -minInt8(n.Right.Left.bFactor, 0) //n.bFactor = n.bFactor + 2 - minInt8(n.Right.Left.bFactor, 0)
			n.Right.bFactor = n.Right.bFactor - 1 - maxInt8(n.Right.Left.bFactor, 0)
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
			hChanged = 1
		} else {
			//R-1 rotation
			n.bFactor = 0
			n.Left.bFactor = 0
			n = rotateLeftRight(n)
			n.bFactor = 0
			hChanged = 1 // height of the tree is changed
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
			hChanged = 1
		} else {
			//L1 rotation
			n.bFactor = 0
			n.Right.bFactor = 0
			n = rotateRightLeft(n)
			n.bFactor = 0
			hChanged = 1 // height of the tree is changed
		}
	}

	return n, hChanged
}

// min returns the minimum node in the tree with root n
func findMin[T any](n *Node[T]) *Node[T] {
	for n.Left != nil {
		n = n.Left
	}

	return n
}

// max return the maximum node in the tree with root n
func findMax[T any](n *Node[T]) *Node[T] {
	for n.Right != nil {
		n = n.Right
	}

	return n
}

func inOrder[T any](n *Node[T], iterator Iterator[T]) {
	if n == nil {
		return
	}

	inOrder(n.Left, iterator)
	iterator(n.Value)
	inOrder(n.Right, iterator)
}

func inOrderReverse[T any](n *Node[T], iterator Iterator[T]) {
	if n == nil {
		return
	}

	inOrderReverse(n.Right, iterator)
	iterator(n.Value)
	inOrderReverse(n.Left, iterator)
}

func preOrder[T any](n *Node[T], iterator Iterator[T]) {
	if n == nil {
		return
	}

	iterator(n.Value)
	preOrder(n.Left, iterator)
	preOrder(n.Right, iterator)
}

func preOrderReverse[T any](n *Node[T], iterator Iterator[T]) {
	if n == nil {
		return
	}

	iterator(n.Value)
	preOrderReverse(n.Right, iterator)
	preOrderReverse(n.Left, iterator)
}

func postOrder[T any](n *Node[T], iterator Iterator[T]) {
	if n == nil {
		return
	}

	postOrder(n.Left, iterator)
	postOrder(n.Right, iterator)
	iterator(n.Value)
}

func postOrderReverse[T any](n *Node[T], iterator Iterator[T]) {
	if n == nil {
		return
	}

	postOrderReverse(n.Right, iterator)
	postOrderReverse(n.Left, iterator)
	iterator(n.Value)
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

func (t *AVL[T]) ascendLessThan(n *Node[T], pivot T, iterator Iterator[T]) bool {
	if n == nil {
		return true
	}

	if !t.ascendLessThan(n.Left, pivot, iterator) {
		return false
	}
	if t.less(n.Value, pivot) {
		if !iterator(n.Value) {
			return false
		}
		return t.ascendLessThan(n.Right, pivot, iterator)
	}

	return true
}

func (t *AVL[T]) ascendGreaterOrEqual(n *Node[T], pivot T, iterator Iterator[T]) bool {
	if n == nil {
		return true
	}

	if t.less(n.Value, pivot) {
		if !t.ascendGreaterOrEqual(n.Left, pivot, iterator) {
			return false
		}

		if !iterator(n.Value) {
			return false
		}
	}

	return t.ascendGreaterOrEqual(n.Right, pivot, iterator)

}

func (t *AVL[T]) ascendRange(n *Node[T], start, end T, iterator Iterator[T]) bool {
	if n == nil {
		return true
	}
	if t.less(end, n.Value) {
		return t.ascendRange(n.Left, start, end, iterator)
	}

	if t.less(n.Value, start) {
		return t.ascendRange(n.Right, start, end, iterator)
	}

	if !t.ascendRange(n.Left, start, end, iterator) {
		return false
	}

	if !iterator(n.Value) {
		return false
	}
	return t.ascendRange(n.Right, start, end, iterator)
}

func (t *AVL[T]) descendGreaterThan(n *Node[T], pivot T, iterator Iterator[T]) bool {
	if n == nil {
		return true
	}

	if !t.descendGreaterThan(n.Right, pivot, iterator) {
		return false
	}

	if t.less(pivot, n.Value) {
		if !iterator(n.Value) {
			return false
		}
		return t.descendGreaterThan(n.Left, pivot, iterator)
	}

	return true
}

func (t *AVL[T]) descendLessOrEqual(n *Node[T], pivot T, iterator Iterator[T]) bool {
	if n == nil {
		return true
	}

	if t.less(n.Value, pivot) {
		if !t.descendLessOrEqual(n.Right, pivot, iterator) {
			return false
		}

		if !iterator(n.Value) {
			return false
		}
	}

	return t.descendLessOrEqual(n.Left, pivot, iterator)
}

func (t *AVL[T]) descendRange(n *Node[T], start, end T, iterator Iterator[T]) bool {
	if n == nil {
		return true
	}

	if t.less(n.Value, start) {
		return t.descendRange(n.Right, start, end, iterator)
	}

	if t.less(end, n.Value) {
		return t.descendRange(n.Left, start, end, iterator)
	}

	if !t.descendRange(n.Right, start, end, iterator) {
		return false
	}

	if !iterator(n.Value) {
		return false
	}

	return t.descendRange(n.Left, start, end, iterator)
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

func (t *AVL[T]) replaceOrInsert(n *Node[T], val T, oldVal **T) (_ *Node[T], hChanged int8) {

	if n == nil {
		return NewNode(val), 1
	}

	if t.less(val, n.Value) {

		n.Left, hChanged = t.replaceOrInsert(n.Left, val, oldVal)
		n.bFactor += hChanged

	} else if t.less(n.Value, val) {

		n.Right, hChanged = t.replaceOrInsert(n.Right, val, oldVal)
		n.bFactor -= hChanged

	} else {
		*oldVal = &n.Value
		n.Value = val
		return n, 0
	}

	if n.bFactor > 1 || n.bFactor < -1 {

		n = insertFix(n)
		hChanged = 0

	} else if hChanged == 1 && n.bFactor == 0 {

		hChanged = 0

	}

	return n, hChanged

}

func (t *AVL[T]) delete(n *Node[T], val T, oldVal **T) (_ *Node[T], hChanged int8) {

	if n == nil {
		return nil, 0
	}

	if t.less(val, n.Value) {
		n.Left, hChanged = t.delete(n.Left, val, oldVal)
		n.bFactor -= hChanged

	} else if t.less(n.Value, val) {
		n.Right, hChanged = t.delete(n.Right, val, oldVal)
		n.bFactor += hChanged

	} else {
		*oldVal = &n.Value //save the pointer to the value to oldVal
		if n.Left == nil && n.Right == nil {
			return nil, 1
		} else if n.Left == nil {
			return n.Right, 1
		} else if n.Right == nil {
			return n.Left, 1
		} else {
			s := findMin(n.Right)
			n.Value = s.Value
			n.Right, hChanged = t.delete(n.Right, s.Value, oldVal)
			n.bFactor += hChanged
		}
	}

	//the tree only has one node, height is change from 1 -> 0 after deletion
	if n == nil {
		return nil, 1
	}

	if n.bFactor > 1 || n.bFactor < -1 {
		return deleteFix(n)
	} else if hChanged == 1 && n.bFactor != 0 {
		return n, 0
	}

	return n, hChanged
}

//
//func (t *AVL[T]) delete(n *Node[T], val T) (_ *Node[T], hChanged int8, oldVal T, found bool) {
//	if n == nil {
//		return nil, 0, oldVal, false // not found
//	}
//
//	if t.less(val, n.Value) {
//		n.Left, hChanged, oldVal, found = t.delete(n.Left, val)
//		n.bFactor -= hChanged
//	} else if t.less(n.Value, val) {
//		n.Right, hChanged, oldVal, found = t.delete(n.Right, val)
//		n.bFactor += hChanged
//		fmt.Println("nRight: ", n.Right.Value, n.Right.bFactor)
//	} else {
//		fmt.Println("found ", n.Value)
//		t.print()
//		found = true
//		oldVal = n.Value
//		if n.Left == nil && n.Right == nil {
//			return nil, 1, oldVal, found // found and deleted
//		} else if n.Right == nil {
//			return n.Left, 1, oldVal, found
//		} else if n.Left == nil {
//			return n.Right, 1, oldVal, found
//		} else {
//			succ := min(n.Right)
//			n.Value = succ.Value
//			n.Right, hChanged, _, _ = t.delete(n.Right, succ.Value)
//			n.bFactor += hChanged
//			fmt.Println("after delete succ ", hChanged)
//		}
//
//	}
//
//	if n == nil {
//		return nil, 1, oldVal, found
//	}
//
//	if n.bFactor > 1 || n.bFactor < -1 {
//		n, hChanged = deleteFix(n)
//	}
//
//	return n, hChanged, oldVal, found
//}

func (t *AVL[T]) print() {

	str := "AVLTree\n"
	if t.root == nil {
		fmt.Println("nil")
		return
	}

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

// Max returns the maximum value in the tree
func (t *AVL[T]) Max() (T, bool) {
	if t.root == nil {
		return *new(T), false
	}
	return findMax(t.root).Value, true
}

// Min returns the minimum value in the tree
func (t *AVL[T]) Min() (T, bool) {
	if t.root == nil {
		return *new(T), false
	}
	return findMin(t.root).Value, true
}

func (t *AVL[T]) AscendLessThan(pivot T, iterator Iterator[T]) {
	t.ascendLessThan(t.root, pivot, iterator)
}

func (t *AVL[T]) AscendGreaterOrEqual(pivot T, iterator Iterator[T]) {
	t.ascendGreaterOrEqual(t.root, pivot, iterator)
}

func (t *AVL[T]) AscendRange(start, end T, iterator Iterator[T]) {
	t.ascendRange(t.root, start, end, iterator)
}

func (t *AVL[T]) DescendGreaterThan(pivot T, iterator Iterator[T]) {
	t.descendGreaterThan(t.root, pivot, iterator)
}

func (t *AVL[T]) DescendLessOrEqual(pivot T, iterator Iterator[T]) {
	t.descendLessOrEqual(t.root, pivot, iterator)
}

func (t *AVL[T]) DescendRange(start, end T, iterator Iterator[T]) {
	t.descendRange(t.root, start, end, iterator)
}

func (t *AVL[T]) InsertNoReplace(val T) {
	t.count++
	t.root = t.insertNoReplace(t.root, val)
}

func (t *AVL[T]) ReplaceOrInsert(val T) (_ T, _ bool) {
	var oldVal *T

	t.root, _ = t.replaceOrInsert(t.root, val, &oldVal)

	if oldVal == nil {
		t.count++
		return
	}
	return *oldVal, true
}

func (t *AVL[T]) Delete(val T) (_ T, _ bool) {

	var oldVal *T
	t.root, _ = t.delete(t.root, val, &oldVal)

	if oldVal != nil {
		t.count--
		return *oldVal, true // found and deleted
	}
	return
}

func (t *AVL[T]) DeleteMax() (T, bool) {
	//temporary solution
	var zero T
	if t.root == nil {
		return zero, false
	}
	max := findMax(t.root)
	t.Delete(max.Value)
	return max.Value, true
}

func (t *AVL[T]) DeleteMin() (T, bool) {
	var zero T
	if t.root == nil {
		return zero, false
	}
	min := findMin(t.root)
	t.Delete(min.Value)
	return min.Value, true
}

func (t *AVL[T]) InOrder(iterator Iterator[T], reverse bool) {
	if reverse {
		inOrderReverse(t.root, iterator)
	} else {
		inOrder(t.root, iterator)
	}

}

func (t *AVL[T]) PreOrder(iterator Iterator[T], reverse bool) {
	if reverse {
		preOrderReverse(t.root, iterator)
	} else {
		preOrder(t.root, iterator)
	}
}

func (t *AVL[T]) PostOrder(iterator Iterator[T], reverse bool) {
	if reverse {
		postOrderReverse(t.root, iterator)
	} else {
		postOrder(t.root, iterator)
	}
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
