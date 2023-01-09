package redblacktree

import (
	"fmt"
	"github.com/nnhatnam/skale"
)

const (
	RED uint8 = iota
	BLACK
)

type Node[K, V any] struct {
	Key   K
	Value V
	color uint8

	Left, Right, Parent *Node[K, V]
}

func (t *Node[K, V]) GetKey() K {
	return t.Key
}

func NewNode[K, V any](key K, value V) *Node[K, V] {
	return &Node[K, V]{Key: key, Value: value, color: RED}
}

type Tree[K, V any] struct {
	root *Node[K, V]
	size int
	less skale.LessFunc[K]
}

func NewTree[K, V any](less skale.LessFunc[K]) *Tree[K, V] {
	return &Tree[K, V]{less: less}
}

func NewTreeOrdered[K, V skale.Ordered]() *Tree[K, V] {
	return NewTree[K, V](skale.Less[K]())
}

// -------------------------private methods-------------------------

// rotateLeft rotates the node x to the left
func (t *Tree[K, V]) rotateLeft(x *Node[K, V]) {
	y := x.Right
	x.Right = y.Left
	if y.Left != nil {
		y.Left.Parent = x
	}
	y.Parent = x.Parent
	if x.Parent == nil {
		t.root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}
	y.Left = x
	x.Parent = y
}

// rotateRight rotates the node x to the right
func (t *Tree[K, V]) rotateRight(x *Node[K, V]) {
	y := x.Left
	x.Left = y.Right
	if y.Right != nil {
		y.Right.Parent = x
	}
	y.Parent = x.Parent
	if x.Parent == nil {
		t.root = y
	} else if x == x.Parent.Right {
		x.Parent.Right = y
	} else {
		x.Parent.Left = y
	}
	y.Right = x
	x.Parent = y
}

// insertFixup fixes the tree after inserting a node
func (t *Tree[K, V]) insertFixup(z *Node[K, V]) {
	for z.Parent != nil && z.Parent.color == RED {
		if z.Parent == z.Parent.Parent.Left {
			y := z.Parent.Parent.Right
			if y != nil && y.color == RED {
				z.Parent.color = BLACK
				y.color = BLACK
				z.Parent.Parent.color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Right {
					z = z.Parent
					t.rotateLeft(z)
				}
				z.Parent.color = BLACK
				z.Parent.Parent.color = RED
				t.rotateRight(z.Parent.Parent)
			}
		} else {
			y := z.Parent.Parent.Left
			if y != nil && y.color == RED {
				z.Parent.color = BLACK
				y.color = BLACK
				z.Parent.Parent.color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Left {
					z = z.Parent
					t.rotateRight(z)
				}
				z.Parent.color = BLACK
				z.Parent.Parent.color = RED
				t.rotateLeft(z.Parent.Parent)
			}
		}
	}
	t.root.color = BLACK
}

// transplant replaces the subtree rooted at u with the subtree rooted at v
func (t *Tree[K, V]) transplant(u, v *Node[K, V]) {
	if u.Parent == nil {
		t.root = v
	} else if u == u.Parent.Left {
		u.Parent.Left = v
	} else {
		u.Parent.Right = v
	}
	if v != nil {
		v.Parent = u.Parent
	}
}

// minimum returns the node with the minimum key in the subtree rooted at x
func (t *Tree[K, V]) minimum(x *Node[K, V]) *Node[K, V] {
	if x == nil {
		return nil
	}
	for x.Left != nil {
		x = x.Left
	}
	return x
}

// maximum returns the node with the maximum key in the subtree rooted at x
func (t *Tree[K, V]) maximum(x *Node[K, V]) *Node[K, V] {
	if x == nil {
		return nil
	}
	for x.Right != nil {
		x = x.Right
	}
	return x
}

// successor returns the node with the smallest key greater than x.key
func (t *Tree[K, V]) successor(x *Node[K, V]) *Node[K, V] {
	if x == nil {
		return nil
	}
	if x.Right != nil {
		return t.minimum(x.Right)
	}
	y := x.Parent
	for y != nil && x == y.Right {
		x = y
		y = y.Parent
	}
	return y
}

// predecessor returns the node with the largest key smaller than x.key
func (t *Tree[K, V]) predecessor(x *Node[K, V]) *Node[K, V] {
	if x == nil {
		return nil
	}
	if x.Left != nil {
		return t.maximum(x.Left)
	}
	y := x.Parent
	for y != nil && x == y.Left {
		x = y
		y = y.Parent
	}
	return y
}

// deleteFixup fixes the tree after deleting a node
func (t *Tree[K, V]) deleteFixup(x *Node[K, V]) {
	if x == nil {
		return
	}
	for x != t.root && x.color == BLACK {
		if x == x.Parent.Left {
			w := x.Parent.Right
			if w.color == RED {
				w.color = BLACK
				x.Parent.color = RED
				t.rotateLeft(x.Parent)
				w = x.Parent.Right
			}
			if w.Left.color == BLACK && w.Right.color == BLACK {
				w.color = RED
				x = x.Parent
			} else {
				if w.Right.color == BLACK {
					w.Left.color = BLACK
					w.color = RED
					t.rotateRight(w)
					w = x.Parent.Right
				}
				w.color = x.Parent.color
				x.Parent.color = BLACK
				w.Right.color = BLACK
				t.rotateLeft(x.Parent)
				x = t.root
			}
		} else {
			w := x.Parent.Left
			if w.color == RED {
				w.color = BLACK
				x.Parent.color = RED
				t.rotateRight(x.Parent)
				w = x.Parent.Left
			}
			if w.Right.color == BLACK && w.Left.color == BLACK {
				w.color = RED
				x = x.Parent
			} else {
				if w.Left.color == BLACK {
					w.Right.color = BLACK
					w.color = RED
					t.rotateLeft(w)
					w = x.Parent.Left
				}
				w.color = x.Parent.color
				x.Parent.color = BLACK
				w.Left.color = BLACK
				t.rotateRight(x.Parent)
				x = t.root
			}
		}
	}

	x.color = BLACK

}

// insert inserts a new node into the tree
func (t *Tree[K, V]) insert(n *Node[K, V]) {
	if n == nil {
		return
	}
	var y *Node[K, V]
	x := t.root
	for x != nil {
		y = x
		if t.less(n.Key, x.Key) {
			x = x.Left
		} else {
			x = x.Right
		}
	}
	n.Parent = y
	if y == nil {
		t.root = n
	} else if t.less(n.Key, y.Key) {
		y.Left = n
	} else {
		y.Right = n
	}
	n.Left = nil
	n.Right = nil
	n.color = RED
	t.insertFixup(n)
}

// delete deletes a node from the tree
func (t *Tree[K, V]) delete(z *Node[K, V]) {
	//implement later
}

// search searches for a node with the given key
func (t *Tree[K, V]) search(key K) *Node[K, V] {
	x := t.root
	for x != nil && !t.less(x.Key, key) {
		if t.less(key, x.Key) {
			x = x.Left
		} else {
			x = x.Right
		}
	}
	return x
}

// inorderTraverse traverses the tree in order
func (t *Tree[K, V]) inorderTraverse(x *Node[K, V], f func(*Node[K, V])) {
	if x == nil {
		return
	}
	t.inorderTraverse(x.Left, f)
	f(x)
	t.inorderTraverse(x.Right, f)
}

// preorderTraverse traverses the tree in preorder
func (t *Tree[K, V]) preorderTraverse(x *Node[K, V], f func(*Node[K, V])) {
	if x == nil {
		return
	}
	f(x)
	t.preorderTraverse(x.Left, f)
	t.preorderTraverse(x.Right, f)
}

// postorderTraverse traverses the tree in postorder
func (t *Tree[K, V]) postorderTraverse(x *Node[K, V], f func(*Node[K, V])) {
	if x == nil {
		return
	}
	t.postorderTraverse(x.Left, f)
	t.postorderTraverse(x.Right, f)
	f(x)
}

// print prints the tree
func (t *Tree[K, V]) print() {
	t.inorderTraverse(t.root, func(x *Node[K, V]) {
		fmt.Println(x)
	})
}

//-------------------------public methods-------------------------

func (t *Tree[K, V]) Size() int {
	return t.size
}

// Insert inserts a new node into the tree
func (t *Tree[K, V]) Insert(key K, value V) {
	t.insert(NewNode(key, value))
	t.size++
}

// Delete deletes a node from the tree
func (t *Tree[K, V]) Delete(key K) {
	z := t.search(key)
	if z != nil {
		t.delete(z)
		t.size--
	}
}

// Search searches for a node with the given key
func (t *Tree[K, V]) Search(key K) *Node[K, V] {
	return t.search(key)
}

// Minimum returns the node with the smallest key
func (t *Tree[K, V]) Minimum() *Node[K, V] {
	return t.minimum(t.root)
}

// Maximum returns the node with the largest key
func (t *Tree[K, V]) Maximum() *Node[K, V] {
	return t.maximum(t.root)
}

// Successor returns the node with the smallest key larger than x.key
func (t *Tree[K, V]) Successor(x *Node[K, V]) *Node[K, V] {
	return t.successor(x)
}

// Predecessor returns the node with the largest key smaller than x.key
func (t *Tree[K, V]) Predecessor(x *Node[K, V]) *Node[K, V] {
	return t.predecessor(x)
}

// InorderTraverse traverses the tree in inorder
func (t *Tree[K, V]) InorderTraverse() {
	t.inorderTraverse(t.root, func(n *Node[K, V]) {
		fmt.Println(n.Key, n.Value)
	})
}

// PreorderTraverse traverses the tree in preorder
func (t *Tree[K, V]) PreorderTraverse() {
	t.preorderTraverse(t.root, func(n *Node[K, V]) {
		fmt.Println(n.Key, n.Value)
	})
}

// PostorderTraverse traverses the tree in postorder
func (t *Tree[K, V]) PostorderTraverse() {
	t.postorderTraverse(t.root, func(n *Node[K, V]) {
		fmt.Println(n.Key, n.Value)
	})
}

// Print prints the tree
func (t *Tree[K, V]) Print() {
	t.print()
}
