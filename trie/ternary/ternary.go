package ternary

import (
	"github.com/nnhatnam/skale/trie"
)

type Node[T trie.Elem] struct {
	Value            T
	Left, Mid, Right *Node[T]
	last             bool
}

func NewNode[T trie.Elem](value T) *Node[T] {
	return &Node[T]{Value: value}
}

type TernaryST[T trie.Elem] struct {
	root *Node[T]
	size int
	//less skale.LessFunc[T]
}

func New[T trie.Elem]() *TernaryST[T] {
	return &TernaryST[T]{root: nil}
}

func replaceOrInsert[T trie.Elem](node *Node[T], value []T, idx int) (_ *Node[T], replaced bool) {
	if node == nil {
		node = NewNode(value[idx])
	}

	if value[idx] < node.Value {
		node.Left, replaced = replaceOrInsert(node.Left, value, idx)
	} else if value[idx] > node.Value {
		node.Right, replaced = replaceOrInsert(node.Right, value, idx)
	} else if idx < len(value)-1 {
		node.Mid, replaced = replaceOrInsert(node.Mid, value, idx+1)
	} else {
		replaced = node.last
		node.Value = value[idx]
		node.last = true
	}

	return node, replaced
}

func get[T trie.Elem](node *Node[T], value []T, idx int) *Node[T] {
	if node == nil {
		return nil
	}

	if value[idx] < node.Value {
		return get(node.Left, value, idx)
	} else if value[idx] > node.Value {
		return get(node.Right, value, idx)
	} else if idx < len(value)-1 {
		return get(node.Mid, value, idx+1)
	} else {
		return node
	}
}

func remove[T trie.Elem](node *Node[T], value []T, idx int) (_ *Node[T], removed bool) {
	if node == nil {
		return
	}

	if value[idx] < node.Value {
		node.Left, removed = remove(node.Left, value, idx)
	} else if value[idx] > node.Value {
		node.Right, removed = remove(node.Right, value, idx)
	} else if idx < len(value)-1 {
		node.Mid, removed = remove(node.Mid, value, idx+1)
	} else {
		node.last = false
	}

	if node.Left == nil && node.Right == nil && node.Mid == nil && !node.last {
		return nil, true // remove node
	}

	return node, removed
}

func search[T trie.Elem](node *Node[T], query []T, idx int, length int) int {
	if node == nil {
		return length
	}

	if query[idx] < node.Value {
		return search(node.Left, query, idx, length)
	} else if query[idx] > node.Value {
		return search(node.Right, query, idx, length)
	} else if idx < len(query)-1 {
		return search(node.Mid, query, idx+1, length+1)
	} else {
		return length + 1
	}

}

func (t *TernaryST[T]) findPrefix(query []T) bool {
	return search(t.root, query, 0, 0) == len(query)
}

func (t *TernaryST[T]) valueWithPrefix(node *Node[T], prefix []T) [][]T {
	//view prefix as a list of legit visited nodes. When visit a node, a node is considered legit if it is the last node of a word, or it has a middle node.
	//if a node is legit, it is appended to the prefix. Otherwise, it is not appended to the prefix.

	if node == nil {
		return nil
	}

	var values [][]T

	if node.last { // legit node
		values = append(values, append(prefix, node.Value))
	}

	values = append(values, t.valueWithPrefix(node.Left, prefix)...)
	values = append(values, t.valueWithPrefix(node.Mid, append(prefix, node.Value))...) // legit node
	values = append(values, t.valueWithPrefix(node.Right, prefix)...)

	return values
}

func (t *TernaryST[T]) longestPrefixOfNode(node *Node[T], value []T, idx int, prefix []T) []T {

	if node == nil {
		return prefix
	}

	if value[idx] < node.Value {
		return t.longestPrefixOfNode(node.Left, value, idx, prefix)
	} else if value[idx] > node.Value {
		return t.longestPrefixOfNode(node.Right, value, idx, prefix)
	} else if idx < len(value)-1 {
		prefix = append(prefix, node.Value)
		return t.longestPrefixOfNode(node.Mid, value, idx+1, prefix)
	} else {
		return prefix
	}
}

func (t *TernaryST[T]) longestPrefixOf(value []T) []T {
	return t.longestPrefixOfNode(t.root, value, 0, []T{})
}

func (t *TernaryST[T]) Values() [][]T {
	return t.valueWithPrefix(t.root, []T{})
}

func (t *TernaryST[T]) Size() int {
	return t.size
}

func (t *TernaryST[T]) Contains(value []T) bool {
	return get(t.root, value, 0) != nil
}

func (t *TernaryST[T]) Get(value []T) bool {
	node := get(t.root, value, 0)
	return node != nil && node.last
}

func (t *TernaryST[T]) Insert(value []T) {
	var replaced bool
	t.root, replaced = replaceOrInsert(t.root, value, 0)
	if !replaced {
		t.size++
	}
}

func (t *TernaryST[T]) Delete(value []T) {
	var removed bool
	t.root, removed = remove(t.root, value, 0)
	if removed {
		t.size--
	}
}

func (t *TernaryST[T]) LongestPrefixOf(value []T) []T {
	return t.longestPrefixOf(value)
}

func (t *TernaryST[T]) ValuesWithPrefix(prefix []T) [][]T {
	node := get(t.root, prefix, 0)

	if len(prefix) > 0 {
		prefix = prefix[:len(prefix)-1]
	}

	return t.valueWithPrefix(node, prefix)
}
