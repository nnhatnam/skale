package avl

import (
	"fmt"
	"github.com/nnhatnam/skale"
)

type Node[K, V any] struct {
	Key         K           // The key of the node
	Value       V           // The value of the node
	Parent      *Node[K, V] // The parent of the node
	Left, Right *Node[K, V] // The left and right children of the node

	bFactor int // balance factor
}

func newNode[K, V any](key K, value V) *Node[K, V] {
	return &Node[K, V]{Key: key, Value: value}
}

// NewNode creates a new node with the given key and value
func NewNode[K, V any](key K, value V) *Node[K, V] {
	return newNode(key, value)
}

type Tree[K, V any] struct {
	root *Node[K, V]       // The root of the tree
	less skale.LessFunc[K] // The less function
	size int
}

func newTree[K, V any](less skale.LessFunc[K]) *Tree[K, V] {
	return &Tree[K, V]{less: less}
}

func New[K, V any](less skale.LessFunc[K]) *Tree[K, V] {
	return newTree[K, V](less)
}

func NewOrdered[K, V skale.Ordered]() *Tree[K, V] {
	return newTree[K, V](skale.Less[K]())
}

// -----------------------------------------internal methods-----------------------------------------

// rotateLeft performs a left rotation on the given node
func (t *Tree[K, V]) rotateLeft(node *Node[K, V]) {
	right := node.Right
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Parent = node.Parent
	if node.Parent == nil {
		t.root = right
	} else if node == node.Parent.Left {
		node.Parent.Left = right
	} else {
		node.Parent.Right = right
	}
	right.Left = node
	node.Parent = right
}

// rotateRight performs a right rotation on the given node
func (t *Tree[K, V]) rotateRight(node *Node[K, V]) {
	left := node.Left
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Parent = node.Parent
	if node.Parent == nil {
		t.root = left
	} else if node == node.Parent.Right {
		node.Parent.Right = left
	} else {
		node.Parent.Left = left
	}
	left.Right = node
	node.Parent = left
}

// rebalance performs a rebalance on the given node
func (t *Tree[K, V]) rebalance(node *Node[K, V]) {
	for node != nil {
		if node.bFactor == -2 {
			if node.Left.bFactor == 1 {
				t.rotateLeft(node.Left)
			}
			t.rotateRight(node)
		} else if node.bFactor == 2 {
			if node.Right.bFactor == -1 {
				t.rotateRight(node.Right)
			}
			t.rotateLeft(node)
		}
		if node.bFactor == 0 {
			break
		}
		node = node.Parent
	}
}

// insert inserts the given node into the tree
func (t *Tree[K, V]) insert(node *Node[K, V]) {
	if t.root == nil {
		t.root = node
		t.size++
		return
	}
	var parent *Node[K, V]
	current := t.root
	for current != nil {
		parent = current
		if t.less(node.Key, current.Key) {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	if t.less(node.Key, parent.Key) {
		parent.Left = node
	} else {
		parent.Right = node
	}
	node.Parent = parent
	t.size++
}

// delete deletes the given node from the tree
func (t *Tree[K, V]) delete(node *Node[K, V]) {
	var child *Node[K, V]
	if node.Left == nil || node.Right == nil {
		if node.Parent == nil {
			t.root = nil
			return
		}
		child = node
	} else {
		child = node.Right
		for child.Left != nil {
			child = child.Left
		}
		node.Key = child.Key
		node.Value = child.Value
	}
	var parent *Node[K, V]
	if child.Left != nil {
		parent = child.Left
	} else {
		parent = child.Right
	}

	if parent != nil {
		parent.Parent = child.Parent
	}

	if child.Parent == nil {
		t.root = parent
	} else {
		if child == child.Parent.Left {
			child.Parent.Left = parent
		} else {
			child.Parent.Right = parent
		}
	}

	if child.bFactor == 0 {
		for parent != nil {
			if parent.Left == child {
				parent.bFactor++
			} else {
				parent.bFactor--
			}
			if parent.bFactor == 1 || parent.bFactor == -1 {
				break
			}
			if parent.bFactor == 2 || parent.bFactor == -2 {
				t.rebalance(parent)
			}
			child = parent
			parent = parent.Parent
		}
	}
	t.size--
}

func (t *Tree[K, V]) search(key K) *Node[K, V] {
	current := t.root
	for current != nil {
		if t.less(key, current.Key) {
			current = current.Left
		} else if t.less(current.Key, key) {
			current = current.Right
		} else {
			return current
		}
	}
	return nil
}

// inorder traverses the tree in order
func (t *Tree[K, V]) inorder(node *Node[K, V], f func(*Node[K, V])) {
	if node == nil {
		return
	}
	t.inorder(node.Left, f)
	f(node)
	t.inorder(node.Right, f)
}

// preorder traverses the tree in preorder
func (t *Tree[K, V]) preorder(node *Node[K, V], f func(*Node[K, V])) {
	if node == nil {
		return
	}
	f(node)
	t.preorder(node.Left, f)
	t.preorder(node.Right, f)
}

// postorder traverses the tree in postorder
func (t *Tree[K, V]) postorder(node *Node[K, V], f func(*Node[K, V])) {
	if node == nil {
		return
	}
	t.postorder(node.Left, f)
	t.postorder(node.Right, f)
	f(node)
}

// min returns the minimum node in the tree
func (t *Tree[K, V]) min(node *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}
	for node.Left != nil {
		node = node.Left
	}
	return node
}

// max returns the maximum node in the tree
func (t *Tree[K, V]) max(node *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

// successor returns the successor of the given node
func (t *Tree[K, V]) successor(node *Node[K, V]) *Node[K, V] {
	if node.Right != nil {
		return t.min(node.Right)
	}
	parent := node.Parent
	for parent != nil && node == parent.Right {
		node = parent
		parent = parent.Parent
	}
	return parent
}

// predecessor returns the predecessor of the given node
func (t *Tree[K, V]) predecessor(node *Node[K, V]) *Node[K, V] {
	if node.Left != nil {
		return t.max(node.Left)
	}
	parent := node.Parent
	for parent != nil && node == parent.Left {
		node = parent
		parent = parent.Parent
	}
	return parent
}

// draw draws the tree
func (t *Tree[K, V]) draw(node *Node[K, V], level int) {
	if node == nil {
		return
	}
	t.draw(node.Right, level+1)
	for i := 0; i < level; i++ {
		fmt.Print("    ")
	}
	fmt.Println(node.Key)
	t.draw(node.Left, level+1)
}

// -----------------------------------------public methods-----------------------------------------

// String returns a string representation of container
func (t *Tree[K, V]) String() string {
	str := "tree: "
	t.inorder(t.root, func(node *Node[K, V]) {
		str += fmt.Sprintf("%v ", node.Key)
	})
	return str
}

// Insert inserts the given key and value into the tree
func (t *Tree[K, V]) Put(key K, value V) {
	t.insert(newNode(key, value))
}

// Delete deletes the given key from the tree
func (t *Tree[K, V]) Delete(key K) {
	node := t.root
	for node != nil {
		if t.less(key, node.Key) {
			node = node.Left
		} else if t.less(node.Key, key) {
			node = node.Right
		} else {
			t.delete(node)
			return
		}
	}
}

// Get returns the value associated with the given key
func (t *Tree[K, V]) Get(key K) (V, bool) {
	node := t.root
	for node != nil {
		if t.less(key, node.Key) {
			node = node.Left
		} else if t.less(node.Key, key) {
			node = node.Right
		} else {
			return node.Value, true
		}
	}
	return *new(V), false
}

// Min returns the minimum key and value in the tree
func (t *Tree[K, V]) Min() (K, V, bool) {
	node := t.root
	if node == nil {
		return *new(K), *new(V), false
	}
	for node.Left != nil {
		node = node.Left
	}
	return node.Key, node.Value, true
}

// Max returns the maximum key and value in the tree
func (t *Tree[K, V]) Max() (K, V, bool) {
	node := t.root
	if node == nil {
		return *new(K), *new(V), false
	}
	for node.Right != nil {
		node = node.Right
	}
	return node.Key, node.Value, true
}

// Floor returns the largest key less than or equal to the given key
func (t *Tree[K, V]) Floor(key K) (K, V, bool) {
	node := t.root
	if node == nil {
		return *new(K), *new(V), false
	}
	var floor *Node[K, V]
	for node != nil {
		if t.less(key, node.Key) {
			node = node.Left
		} else if t.less(node.Key, key) {
			floor = node
			node = node.Right
		} else {
			return node.Key, node.Value, true
		}
	}
	if floor == nil {
		return *new(K), *new(V), false
	}
	return floor.Key, floor.Value, true
}

// Ceiling returns the smallest key greater than or equal to the given key
func (t *Tree[K, V]) Ceiling(key K) (K, V, bool) {
	node := t.root
	if node == nil {
		return *new(K), *new(V), false
	}
	var ceiling *Node[K, V]
	for node != nil {
		if t.less(key, node.Key) {
			ceiling = node
			node = node.Left
		} else if t.less(node.Key, key) {
			node = node.Right
		} else {
			return node.Key, node.Value, true
		}
	}

	if ceiling == nil {
		return *new(K), *new(V), false
	}

	return ceiling.Key, ceiling.Value, true
}

// Select returns the key and value of the given rank
func (t *Tree[K, V]) Select(rank int) (K, V, bool) {
	//implement later
	return *new(K), *new(V), false
}

// Rank returns the rank of the given key
func (t *Tree[K, V]) Rank(key K) int {
	//implement later
	return 0
}

// Size returns the number of keys in the tree
func (t *Tree[K, V]) Size() int {
	return t.size
}

// IsEmpty returns true if the tree is empty
func (t *Tree[K, V]) IsEmpty() bool {
	return t.size == 0
}

// Do calls the function f on every element of the tree in inorder
func (t *Tree[K, V]) Do(f func(K, V)) {
	t.inorder(t.root, func(node *Node[K, V]) {
		f(node.Key, node.Value)
	})
}

// DoPreorder calls the function f on every element of the tree in preorder
func (t *Tree[K, V]) DoPreorder(f func(K, V)) {
	t.preorder(t.root, func(node *Node[K, V]) {
		f(node.Key, node.Value)
	})
}

// DoPostorder calls the function f on every element of the tree in postorder
func (t *Tree[K, V]) DoPostorder(f func(K, V)) {
	t.postorder(t.root, func(node *Node[K, V]) {
		f(node.Key, node.Value)
	})
}

// Keys returns a slice of all keys in the tree
func (t *Tree[K, V]) Keys() []K {
	keys := make([]K, t.size)
	i := 0
	t.inorder(t.root, func(node *Node[K, V]) {
		keys[i] = node.Key
		i++
	})
	return keys
}

// KeysInRange returns a slice of all keys in the tree in the given range
func (t *Tree[K, V]) KeysInRange(lo, hi K) []K {
	keys := make([]K, 0)
	t.inorder(t.root, func(node *Node[K, V]) {
		if t.less(node.Key, lo) {
			return
		}
		if t.less(hi, node.Key) {
			return
		}
		keys = append(keys, node.Key)
	})
	return keys
}

// Values returns a slice of all values in the tree
func (t *Tree[K, V]) Values() []V {
	values := make([]V, t.size)
	i := 0
	t.inorder(t.root, func(node *Node[K, V]) {
		values[i] = node.Value
		i++
	})
	return values
}

// ValuesInRange returns a slice of all values in the tree in the given range
func (t *Tree[K, V]) ValuesInRange(lo, hi K) []V {
	values := make([]V, 0)
	t.inorder(t.root, func(node *Node[K, V]) {
		if t.less(node.Key, lo) {
			return
		}
		if t.less(hi, node.Key) {
			return
		}
		values = append(values, node.Value)
	})
	return values
}

// DeleteMin deletes the minimum key and value in the tree
func (t *Tree[K, V]) DeleteMin() {
	t.delete(t.min(t.root))
}

// DeleteMax deletes the maximum key and value in the tree
func (t *Tree[K, V]) DeleteMax() {
	t.delete(t.max(t.root))
}

// Draw prints the tree to stdout
func (t *Tree[K, V]) Draw() {
	t.draw(t.root, 0)
}
