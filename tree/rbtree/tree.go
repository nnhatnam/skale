package rbtree

import (
	"fmt"
	"github.com/nnhatnam/skale"
)

type Iterator[T any] func(item T) bool

const (
	negativeBlack      = -1
	red           int8 = 0
	black         int8 = 1
	doubleBlack   int8 = 2
)

type Node[T any] struct {
	Value               T // value stored in the node
	Left, Right, Parent *Node[T]

	Black int8 // 0: red, 1: black, 2: double black
	Color int8 // -1: negative black , 0: red, 1: black, 2: double black
}

func (n *Node[T]) String() string {
	if isBlack(n) {
		return fmt.Sprintf("%v (black)", n.Value)
	}

	return fmt.Sprintf("%v (red)", n.Value)

}

func newNode[T any](value T) *Node[T] {
	return &Node[T]{Value: value}
}

func isRed[T any](n *Node[T]) bool {
	if n == nil {
		return false
	}
	return n.Black == 0
}

func isBlack[T any](n *Node[T]) bool {
	if n == nil {
		return false
	}
	return n.Black == 1
}

func isDoubleBlack[T any](n *Node[T]) bool {
	if n == nil {
		return false
	}
	return n.Black == 2
}

func isNegativeBlack[T any](n *Node[T]) bool {
	if n == nil {
		return false
	}
	return n.Black == -1
}

type Rb[T any] struct {
	root  *Node[T]
	count int
	less  skale.LessFunc[T]
}

func New[T any](less skale.LessFunc[T]) *Rb[T] {
	return &Rb[T]{less: less}
}

func NewOrdered[T skale.Ordered]() *Rb[T] {
	return New(skale.Less[T]())
}

// -----------------------------------------internal methods-----------------------------------------

func rotateLeft[T any](n *Node[T]) *Node[T] {
	r := n.Right
	if isBlack(r) {
		panic("rotating a black link")
	}
	n.Right = r.Left
	r.Left = n
	r.Black, n.Black = n.Black, r.Black
	return r
}

func rotateRight[T any](n *Node[T]) *Node[T] {
	l := n.Left
	if isBlack(l) {
		panic("rotating a black link")
	}
	n.Left = l.Right
	l.Right = n
	l.Black, n.Black = n.Black, l.Black
	return l

}

// @require: g.Left != nil && g.Right != nil. g and its children can't be double black
func flip[T any](g *Node[T]) {
	g.Black = g.Black ^ 1
	g.Left.Black = g.Left.Black ^ 1
	g.Right.Black = g.Right.Black ^ 1
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

// @require: insertFixUp must be called on ancestor of inserted node. it can't be called on inserted node itself or its parent
func insertFixUp[T any](g *Node[T]) *Node[T] {

	if isRed(g.Left) && isRed(g.Right) {
		flip(g)
	}

	if isRed(g.Left) && isRed(g.Left.Left) {
		g = rotateRight(g)
	}

	if isRed(g.Left) && isRed(g.Left.Right) {
		g.Left = rotateLeft(g.Left)
		g = rotateRight(g)
	}

	if isRed(g.Right) && isRed(g.Right.Right) {
		g = rotateLeft(g)
	}

	if isRed(g.Right) && isRed(g.Right.Left) {
		g.Right = rotateRight(g.Right)
		g = rotateLeft(g)
	}

	return g
}

func bubbleUp[T any](n *Node[T]) {
	n.Black++
	if n.Left != nil {
		n.Left.Black--
	}

	if n.Right != nil {
		n.Right.Black--
	}

}

func bubbleDown[T any](n *Node[T]) {
	n.Black--
	if n.Left != nil {
		n.Left.Black++
	}

	if n.Right != nil {
		n.Right.Black++
	}

}

func bubbleFix[T any](g *Node[T]) *Node[T] {

	if isRed(g.Left) && isRed(g.Left.Left) {
		g = rotateRight(g)
		bubbleDown(g)
	}

	if isRed(g.Right) && isRed(g.Right.Right) {
		g = rotateLeft(g)
		bubbleDown(g)
	}

	if isRed(g.Left) && isRed(g.Left.Right) {
		g.Left = rotateLeft(g.Left)
		g = rotateRight(g)
		bubbleDown(g)
	}

	if isRed(g.Right) && isRed(g.Right.Left) {
		g.Right = rotateRight(g.Right)
		g = rotateLeft(g)
		bubbleDown(g)
	}

	if isNegativeBlack(g.Left) {
		g.Black--
		g.Left.Left.Black--
		g.Left.Left.Black += 2
		g.Left = rotateRight(g.Left)
		g = rotateLeft(g)
		g.Left = bubbleFix(g.Left)

	}

	if isNegativeBlack(g.Right) {
		g.Black--
		g.Right.Right.Black--
		g.Right.Right.Black += 2
		g.Right = rotateLeft(g.Right)
		g = rotateRight(g)
		g.Right = bubbleFix(g.Right)
	}

	return g
}

func fixNegativeBlackLeft[T any](p *Node[T]) *Node[T] {

	if isBlack(p.Right) {
		if isRed(p.Right.Right) {
			p = rotateRight(p)
		} else if isRed(p.Right.Left) {
			p.Right = rotateLeft(p.Right)
			p = rotateRight(p)
		} else {
			p.Black++ //red -> black, black -> double black
			p.Right.Black = red
		}
	}

	if isBlack(p) && isBlack(p.Right) {
		if isRed(p.Right.Right) {
			p = rotateRight(p)
		} else if isRed(p.Right.Left) {
			p.Right = rotateLeft(p.Right)
			p = rotateRight(p)
		} else {
			p.Black = doubleBlack
			p.Right.Black = red
		}
	}

	return p
}

func moveBlackLeft[T any](p *Node[T]) *Node[T] {
	if isBlack(p) && isBlack(p.Right) {
		if isRed(p.Right.Right) {
			p = rotateRight(p)
		} else if isRed(p.Right.Left) {
			p.Right = rotateLeft(p.Right)
			p = rotateRight(p)
		} else {
			p.Black = doubleBlack
			p.Right.Black = red
		}
	}
	return p
}

func fixNegativeBlackRight[T any](p *Node[T]) *Node[T] {

	if isRed(p) && isNegativeBlack(p.Left) {
		p.Left = nil // delete p.Left
		p.Black = black
	}

	if isRed(p) && isNegativeBlack(p.Right) {
		p.Right = nil // delete p.Right
		p.Black = black
	}

	return p
}

func deleteFixUp[T any](p *Node[T]) *Node[T] {

	if isRed(p) && isNegativeBlack(p.Left) {
		p.Left = nil // delete p.Left
		p.Black = black
	}

	if isRed(p) && isNegativeBlack(p.Right) {
		p.Right = nil // delete p.Right
		p.Black = black
	}

	if isBlack(p) && isDoubleBlack(p.Left) {
		p.Left = nil // delete p.Left
		p.Black = black
	}

	return p
}

// find the node with value val in tree t, return the node
func (t *Rb[T]) find(val T) (_ *Node[T]) {
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

func (t *Rb[T]) ascendLessThan(n *Node[T], pivot T, iterator Iterator[T]) bool {
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

func (t *Rb[T]) ascendGreaterOrEqual(n *Node[T], pivot T, iterator Iterator[T]) bool {
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

func (t *Rb[T]) ascendRange(n *Node[T], start, end T, iterator Iterator[T]) bool {
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

func (t *Rb[T]) descendGreaterThan(n *Node[T], pivot T, iterator Iterator[T]) bool {
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

func (t *Rb[T]) descendLessOrEqual(n *Node[T], pivot T, iterator Iterator[T]) bool {
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

func (t *Rb[T]) descendRange(n *Node[T], start, end T, iterator Iterator[T]) bool {
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

func (t *Rb[T]) insertNoReplace(n *Node[T], val T) *Node[T] {

	if n == nil {
		return newNode(val)
	}

	if t.less(val, n.Value) {
		n.Left = t.insertNoReplace(n.Left, val)
	} else {
		n.Right = t.insertNoReplace(n.Right, val)
	}

	return insertFixUp(n)

}

func (t *Rb[T]) replaceOrInsert(n *Node[T], val T) (_ *Node[T], replaced *T) {

	if n == nil {
		return newNode(val), nil
	}

	if t.less(val, n.Value) {
		n.Left, replaced = t.replaceOrInsert(n.Left, val)
	} else if t.less(n.Value, val) {
		n.Right, replaced = t.replaceOrInsert(n.Right, val)
	} else {
		replaced, n.Value = &n.Value, val
	}

	return insertFixUp(n), replaced

}

func (t *Rb[T]) remove(n *Node[T]) *Node[T] {

	if n.Left == nil && n.Right == nil {
		return nil
	}

	if n.Left == nil {
		return n.Right
	}

	if n.Right == nil {
		return n.Left
	}

	return n
}

func (t *Rb[T]) delete(n *Node[T], val T) (_ *Node[T], deleted *T) {

	if n == nil {
		return nil, nil
	}

	if t.less(val, n.Value) {
		if n.Left == nil {
			return n, nil // not found
		}
		isLBlack := isBlack(n.Left)
		n.Left, deleted = t.delete(n.Left, val)

		if n.Left == nil && isLBlack { // deleted
			//n.Right must be present because if it isn't, it will violate the red-black tree property
			//bubbling
			bubbleUp(n)
			n = bubbleFix(n)

		}

	} else if t.less(n.Value, val) {
		if n.Right == nil {
			return n, nil // not found
		}
		isRBlack := isBlack(n.Right)
		n.Right, deleted = t.delete(n.Right, val)
		if n.Left == nil && isRBlack { // deleted

			//bubbling
			bubbleUp(n)
			n = bubbleFix(n)
		}

	} else {
		deleted = &n.Value

		if n.Left == nil && n.Right == nil {
			return nil, deleted

		} else if n.Left == nil {
			n.Right.Black = n.Black + n.Right.Black
			return n.Right, deleted
		} else if n.Right == nil {
			n.Left.Black = n.Left.Black + n.Black
			return n.Left, deleted
		} else {
			s := findMin(n.Right) // find successor
			n.Value = s.Value
			n.Right, _ = t.delete(n.Right, s.Value)
		}

	}

	if isNegativeBlack(n.Left) || isNegativeBlack(n.Right) {
		bubbleUp(n)
		n = bubbleFix(n)
	}

	n = deleteFixUp(n)

	return n, deleted

	//return nil, deleted
}

func (t *Rb[T]) print() {

	str := "RbTree\n"
	if t.root == nil {
		fmt.Println("nil")
		return
	}

	output(t.root, "", true, &str)
	fmt.Println(str)
}

// -----------------------------------------public methods-----------------------------------------

// Len returns the number of nodes in the tree
func (t *Rb[T]) Len() int {
	return t.count
}

// Get returns the value of the node with value val
func (t *Rb[T]) Get(val T) (_ T, _ bool) {
	n := t.find(val)
	if n != nil {
		return n.Value, true
	}
	return
}

// Has returns true if the tree contains a node with value val
func (t *Rb[T]) Has(val T) bool {
	_, ok := t.Get(val)
	return ok
}

// Max returns the maximum value in the tree
func (t *Rb[T]) Max() (T, bool) {
	if t.root == nil {
		return *new(T), false
	}
	return findMax(t.root).Value, true
}

// Min returns the minimum value in the tree
func (t *Rb[T]) Min() (T, bool) {
	if t.root == nil {
		return *new(T), false
	}
	return findMin(t.root).Value, true
}

func (t *Rb[T]) AscendLessThan(pivot T, iterator Iterator[T]) {
	t.ascendLessThan(t.root, pivot, iterator)
}

func (t *Rb[T]) AscendGreaterOrEqual(pivot T, iterator Iterator[T]) {
	t.ascendGreaterOrEqual(t.root, pivot, iterator)
}

func (t *Rb[T]) AscendRange(start, end T, iterator Iterator[T]) {
	t.ascendRange(t.root, start, end, iterator)
}

func (t *Rb[T]) DescendGreaterThan(pivot T, iterator Iterator[T]) {
	t.descendGreaterThan(t.root, pivot, iterator)
}

func (t *Rb[T]) DescendLessOrEqual(pivot T, iterator Iterator[T]) {
	t.descendLessOrEqual(t.root, pivot, iterator)
}

func (t *Rb[T]) DescendRange(start, end T, iterator Iterator[T]) {
	t.descendRange(t.root, start, end, iterator)
}

func (t *Rb[T]) InsertNoReplace(val T) {
	t.count++
	t.root = t.insertNoReplace(t.root, val)
}

func (t *Rb[T]) ReplaceOrInsert(val T) (_ T, _ bool) {
	var replaced *T

	t.root, replaced = t.replaceOrInsert(t.root, val)
	t.root.Black = black
	if replaced == nil {
		t.count++
		return
	}
	return *replaced, true
}

func (t *Rb[T]) Delete(val T) (_ T, _ bool) {

	var replaced *T
	t.root, replaced = t.delete(t.root, val)
	t.root.Black = black

	if replaced != nil {
		t.count--
		return *replaced, true // found and deleted
	}
	return
}

func (t *Rb[T]) DeleteMax() (T, bool) {
	//temporary solution
	var zero T
	if t.root == nil {
		return zero, false
	}
	max := findMax(t.root)
	t.Delete(max.Value)
	return max.Value, true
}

func (t *Rb[T]) DeleteMin() (T, bool) {
	var zero T
	if t.root == nil {
		return zero, false
	}
	min := findMin(t.root)
	t.Delete(min.Value)
	return min.Value, true
}

func (t *Rb[T]) InOrder(iterator Iterator[T], reverse bool) {
	if reverse {
		inOrderReverse(t.root, iterator)
	} else {
		inOrder(t.root, iterator)
	}

}

func (t *Rb[T]) PreOrder(iterator Iterator[T], reverse bool) {
	if reverse {
		preOrderReverse(t.root, iterator)
	} else {
		preOrder(t.root, iterator)
	}
}

func (t *Rb[T]) PostOrder(iterator Iterator[T], reverse bool) {
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
