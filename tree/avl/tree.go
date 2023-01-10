package avl

import (
	"fmt"
	"github.com/nnhatnam/skale"
)

type direction uint8

var (
	LEFT  direction = 0
	RIGHT direction = 1
)

// the idea is we can create a Node with anytype, but the node maynot be a tree node if we can't compare it with other nodes
type Node[V any] struct {
	Value  V        // The value of the node
	Parent *Node[V] // The parent of the node
	//Left, Right *Node[V]    // The left and right children of the node
	Child [2]*Node[V] // The left and right children of the node

	bFactor int8 // balance factor
}

func (v Node[V]) String() string {
	return fmt.Sprintf("%v (%v)", v.Value, v.bFactor)
}

func NewNode[V any](value V) *Node[V] {
	return &Node[V]{Value: value}
}

type Tree[V any] struct {
	root *Node[V]          // The root of the tree
	less skale.LessFunc[V] // The less function
	size int
}

func New[V any](less skale.LessFunc[V]) *Tree[V] {
	return &Tree[V]{less: less}
}

func NewOrdered[V skale.Ordered]() *Tree[V] {
	return New(skale.Less[V]())
}

// -----------------------------------------internal methods-----------------------------------------

func rotate[V any](node *Node[V], dir direction) *Node[V] {
	//dir = 0: left, dir = 1: right

	//Example: dir = 0, which is left rotation
	y := node.Child[1-dir] // y := node.Right
	t := y.Child[dir]      // t := y.Left

	//perform rotation
	y.Child[dir] = node   // y.Left = node
	node.Child[1-dir] = t // node.Right = t

	if t != nil {
		t.Parent = node // t.Parent = node
	}

	y.Parent = node.Parent // y.Parent = node.Parent

	node.Parent = y // node.Parent = y
	str := "AVLTree\n"
	output[V](node, "", true, &str)
	//fmt.Println("node:", str)

	return y

}

// rotateLeft performs a left rotation on the given node
func (t *Tree[V]) rotateLeft(node *Node[V]) *Node[V] {
	return rotate(node, LEFT)
}

func (t *Tree[V]) rotateRightLeft(node *Node[V]) *Node[V] {
	node.Child[RIGHT] = t.rotateRight(node.Child[RIGHT])
	return t.rotateLeft(node)
}

// rotateRight performs a right rotation on the given node
func (t *Tree[V]) rotateRight(node *Node[V]) *Node[V] {

	return rotate(node, RIGHT)

}

func (t *Tree[V]) rotateLeftRight(node *Node[V]) *Node[V] {
	node.Child[LEFT] = t.rotateLeft(node.Child[LEFT])
	return t.rotateRight(node)
}

// insert insert val into a tree t
// 1. Find the spot
// 2. Insert the node
// 3. Balance the tree
func (t *Tree[V]) insert(node *Node[V], val V) (*Node[V], bool) {

	// 0 : false, 1: true
	var isHeightChanged bool
	if node == nil {
		return NewNode(val), true
	}

	if t.less(val, node.Value) {

		node.Child[LEFT], isHeightChanged = t.insert(node.Child[LEFT], val)
		node.Child[LEFT].Parent = node

		// after insertion, if the height of the left subtree is changed, we need to change to balance factor of the parent.
		//In insertion, insert a node either make the tree "more balance" or grow the tree. After insertion, if the balance factor
		//is 0, the height is not changed, otherwise, the height is changed.
		if isHeightChanged {
			node.bFactor++
			isHeightChanged = !(node.bFactor == 0)
		}

	} else {

		node.Child[RIGHT], isHeightChanged = t.insert(node.Child[RIGHT], val)
		node.Child[RIGHT].Parent = node

		if isHeightChanged {

			node.bFactor--
			isHeightChanged = !(node.bFactor == 0)
		}

	}

	//bFactor = 0 => balance
	//bFactor = 1 => left heavy
	//bFactor = -1 => right heavy

	if node.bFactor > 1 {
		if node.Child[LEFT].bFactor > 0 {
			//left left case
			node.bFactor = 0
			node = t.rotateRight(node)
			node.bFactor = 0
		} else {
			//left right case
			node.bFactor = 0
			node.Child[LEFT].bFactor = 0
			node = t.rotateLeftRight(node)
		}
		return node, false
	} else if node.bFactor < -1 {
		if node.Child[RIGHT].bFactor < 0 {
			//right right case
			node.bFactor = 0
			node = t.rotateLeft(node)
			node.bFactor = 0
		} else {
			//right left case
			node.bFactor = 0
			node.Child[RIGHT].bFactor = 0
			node = t.rotateRightLeft(node)

		}
		return node, false
	}

	return node, isHeightChanged

}

// successor returns the successor of the node
func (t *Tree[V]) successor(node *Node[V]) *Node[V] {
	if node == nil {
		return nil
	}

	node = node.Child[RIGHT]
	for node.Child[LEFT] != nil {
		node = node.Child[LEFT]
	}

	return node
}

// predecessor returns the predecessor of the node
func (t *Tree[V]) predecessor(node *Node[V]) *Node[V] {
	if node == nil {
		return nil
	}

	node = node.Child[LEFT]
	for node.Child[RIGHT] != nil {
		node = node.Child[RIGHT]
	}

	return node
}

// delete delete val from a tree t
func (t *Tree[V]) delete(node *Node[V], val V) (*Node[V], int8) {
	//Before deletion
	var hChanged int8
	if node == nil {
		return node, 0
	}

	//Delete the node

	if t.less(val, node.Value) {
		node.Child[LEFT], hChanged = t.delete(node.Child[LEFT], val)
		node.bFactor -= hChanged

	} else if t.less(node.Value, val) {
		node.Child[RIGHT], hChanged = t.delete(node.Child[RIGHT], val)
		node.bFactor += hChanged
	} else {

		//the node has no child
		if node.Child[LEFT] == nil && node.Child[RIGHT] == nil {
			node.Parent = nil //delete the node
			node = nil
			return node, 1
		} else if node.Child[RIGHT] == nil {
			//the node has only left child
			node.Child[LEFT].Parent = node.Parent
			node.Parent = nil
			node = node.Child[LEFT]
			hChanged = 1
		} else if node.Child[LEFT] == nil {
			//the node has only right child
			node.Child[RIGHT].Parent = node.Parent
			node.Parent = nil
			node = node.Child[RIGHT]
			hChanged = 1
		} else {
			//the node has both left and right child
			//find the successor
			successor := t.successor(node)
			node.Value = successor.Value
			node.Child[RIGHT], hChanged = t.delete(node.Child[RIGHT], successor.Value)

		}

	}

	//After deletion

	//if the node has only one child, then the node is the successor
	if node == nil {
		return node, 1
	}

	//update the bFactor
	if node.bFactor > 1 {
		if node.Child[LEFT].bFactor == 0 {
			//R0 case
			node.bFactor = 1
			node = t.rotateRight(node)
			node.bFactor = -1
		} else if node.Child[LEFT].bFactor == 1 {
			//R1 case
			node.bFactor = 0
			node = t.rotateRight(node)
			node.bFactor = 0
		} else if node.Child[LEFT].bFactor == -1 {
			//R-1 case
			node.bFactor = 0
			node.Child[LEFT].bFactor = 0
			node = t.rotateLeftRight(node)
			hChanged = 1
		}

	} else if node.bFactor < -1 {
		if node.Child[RIGHT].bFactor == 0 {
			//L0 case
			node.bFactor = -1
			node = t.rotateLeft(node)
			node.bFactor = 1
		} else if node.Child[RIGHT].bFactor == -1 {
			//L-1 case
			node.bFactor = 0
			node = t.rotateLeft(node)
			node.bFactor = 0
		} else if node.Child[RIGHT].bFactor == 1 {
			//L1 case
			node.bFactor = 0
			node.Child[RIGHT].bFactor = 0
			node = t.rotateRightLeft(node)
			hChanged = 1
		}

	}

	return node, hChanged

}

func (t *Tree[V]) insertFromRoot(val V) {
	t.size++
	t.root, _ = t.insert(t.root, val)
}

func (t *Tree[V]) output() {
	str := "AVLTree\n"
	output[V](t.root, "", true, &str)

	fmt.Println(str)
}

// -----------------------------------------public methods-----------------------------------------

func (t *Tree[V]) Insert(val V) {
	t.insertFromRoot(val)
}

func output[V any](node *Node[V], prefix string, isTail bool, str *string) {
	if node.Child[1] != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.Child[1], newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.Child[0] != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.Child[0], newPrefix, true, str)
	}
}
