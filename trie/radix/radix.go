package radix

import (
	"fmt"
	"github.com/nnhatnam/skale/exp/xslices"
	"github.com/nnhatnam/skale/trie"
	"golang.org/x/exp/slices"
)

type edge[K trie.Elem, V any] struct {
	label []K         // label of the edge
	node  *node[K, V] // node at the end of the edge
}

func newEdge[K trie.Elem, V any](label []K, node *node[K, V]) *edge[K, V] {

	return &edge[K, V]{label: label, node: node}

}

type node[K trie.Elem, V any] struct {
	edges edges[K, V] // edges from this node, sorted by label

	lastElem bool // true if this node is a leaf

	//prefix []K // prefix of the word represented by this node

	value V // value of the word represented by this node

}

// newInternalNode creates a new internal node with the given prefix.
// internal nodes won't have a value.
//func newInternalNode[K trie.Elem, V any]() *node[K, V] {
//
//	return &node[K, V]{}
//
//}

// newLeafNode creates a new leaf node with the given value.
// leaf nodes won't have any edges, and their prefix will be the same as the word they represent.
// leaf nodes will store values.
func newLeafNode[K trie.Elem, V any](value V) *node[K, V] {

	return &node[K, V]{value: value, lastElem: true}

}

func newNode[K trie.Elem, V any](value V, lastElem bool) *node[K, V] {

	return &node[K, V]{value: value, lastElem: lastElem}

}

func (n *node[K, V]) findEdge(label []K) *edge[K, V] {

	for _, e := range n.edges {

		if slices.Equal(e.label, label) {

			return e

		}

	}

	return nil

}

func (n *node[K, V]) findEdgeWithPrefix(prefix []K) *edge[K, V] {

	for _, e := range n.edges {

		if len(prefix) <= len(e.label) && slices.Equal(e.label[:len(prefix)], prefix) {

			return e

		}

	}
	return nil
}

type RadixTrieMap[K trie.Elem, V any] struct {
	root *node[K, V]

	len int // number of keys in the trie
}

func NewRadixTrieMap[K trie.Elem, V any]() *RadixTrieMap[K, V] {
	var zero V
	return &RadixTrieMap[K, V]{root: newNode[K, V](zero, false)}

}

func (t *RadixTrieMap[K, V]) lazyInit() {

	if t.root == nil {

		var zero V

		t.root = newNode[K, V](zero, false)

	}

}

// walk conceptually walks through every element of the trie, starting from the given node, and reports the location when it stops.
// While walking, it will compare the element at step i of the key with the label of the edge at step i.
// It will stop when it reaches the end of the key, or there is a mismatch between the key and the label.
// walk will report the location where it stops (the edge or the node at that location). There are several cases:
// 1. walk stops at a node, meaning that s still has some elements left, and it can't find the next edge. It will report the node and the index `kIdx` of the key (reported edge will be null)
// 2. walk stops on an edge. In this case, there are two subcases:
// 2.1. walk stops in the middle of the edge. It will report the edge, the index `eIdx` of the edge's label (reported node will be nil)
// 2.2. walk stops at the end of the edge. It will report the edge, the index `eIdx` of the edge's label and the index `kIdx` of the key (reported node will be null)
// In this case, i == len(key) - 1, because if it's not, it will be case 1 or case 2.1
//
// Also, in case 2.1 and 2.2, the index eIdx and kIdx are the index of the last matched element in key and edge.label
func walk[K trie.Elem, V any](n *node[K, V], key []K) (e *edge[K, V], n1 *node[K, V], kIdx, eIdx int) {

	if n == nil {

		return

	}

	i, j := 0, 0

	for i < len(key) {

		e = getEdgeByPrefix[K, V](n.edges, key[i])

		if e == nil {

			return nil, n, i - 1, -1
		}

		for j = 0; j < len(e.label) && i < len(key); j++ {
			if key[i] != e.label[j] {
				return e, nil, i - 1, j - 1
			}
			i++
		}

		n = e.node
	}

	return e, nil, i - 1, j - 1 // i - 1, j - 1 is the index of the last matched element in key and edge.label

}

func replaceOrInsert[K trie.Elem, V any](root *node[K, V], s []K, value V) (_ V, _ bool) {
	// To be considered: this function may need to write in a way such that s will be garbage collected at the end of the function
	// Due to the design of Go slice, the underline data of s may not release if we set trie's key reference to a slice of s

	if root == nil || len(s) == 0 {
		//n = newInternalNode[K, V](s)
		panic("n is nil")
		return
	}

	e, n, kIdx, eIdx := walk[K, V](root, s)

	if n != nil {
		// walk stops at a node, meaning that s still has some elements left, and it can't find the next edge.
		// => the leftover elements are not in the trie, and we need to insert them.

		// create a new leaf node with the leftover elements
		n1 := newLeafNode[K, V](value)
		// create a new edge with the leftover elements
		e1 := newEdge(xslices.SubClone(s, kIdx+1, len(s)), n1)
		// add the new edge to the current node
		n.edges = setEdge(n.edges, e1)
		return

	}

	// walk stop on an edge

	if eIdx == len(e.label)-1 {
		// walk stops at the end of the edge -> s has no more elements left (because if it does, walk would have stopped at a node or in the middle of the edge)

		e.node.value = value
		return
	}

	// walk stops in the middle of the edge -> split at the mismatched indexes (eIdx & kIdx)

	eSplitIndex := eIdx + 1
	kSplitIndex := kIdx + 1
	var zero V
	// split the edge to two edges at the mismatched index (eIdx)
	// (n0) --e--> (n4)  => (n0) --e --> (n1) --e1--> (n4)
	e1 := newEdge(xslices.SubClone(e.label, eSplitIndex, len(e.label)), e.node)
	n1 := newNode[K, V](zero, false)
	n1.edges = append(n1.edges, e1)

	e.label = xslices.SubClone(e.label, 0, eSplitIndex)
	e.node = n1

	// insert the leftover elements into the trie
	// (n0) --e --> (n1) --e1--> (n4)
	// => (n0) --e --> (n1) --e1--> (n4)
	//                  |
	//                  +--e2--> (n2)

	n2 := newLeafNode[K, V](value)
	e2 := newEdge(xslices.SubClone(s, kSplitIndex, len(s)), n2)
	n1.edges = setEdge(n1.edges, e2)

	return
}

func insert[K trie.Elem, V any](n *node[K, V], s []K, value V) {
	// To be considered: this function may need to write in a way such that s will be garbage collected at the end of the function
	// Due to the design of Go slice, the underline data of s may not release if we set trie's key reference to a slice of s

	if n == nil || len(s) == 0 {
		//n = newInternalNode[K, V](s)
		panic("n is nil")
		return
	}

	e := getEdgeByPrefix[K, V](n.edges, s[0])

	if e == nil {
		e1 := newEdge(xslices.SubClone(s, 0, len(s)), newLeafNode[K, V](value))
		n.edges = setEdge(n.edges, e1)
		return

	}

	if slices.Equal(e.label, s) {
		// s is already in the trie
		// update the value
		e.node.value = value
		return
	}

	splitIdx := xslices.MismatchIndex(s, e.label)
	//var zero V

	if splitIdx == -1 {
		// no mismatch and len(s) == len(n.edges[idx].label)

		e.node.value = value
		return

	} else if splitIdx == len(s) {
		// s is a prefix of n.edges[idx].label => split the current edge at splitIdx.

		e1 := newEdge(
			xslices.SubClone(e.label, 0, splitIdx),
			newLeafNode[K, V](value),
		)

		e.label = xslices.SubClone(e.label, splitIdx, len(e.label))

		setEdge(n.edges, e1)
		setEdge(e1.node.edges, e)
		return
	} else if splitIdx == len(e.label) {
		// n.edges[idx].label is a prefix of s => add new edge after splitIdx.

		e1 := newEdge(
			xslices.SubClone(s, splitIdx, len(s)),
			newLeafNode[K, V](value),
		)

		setEdge(e.node.edges, e1)

		return
	}

	// split the current edge at splitIdx.

	e1 := newEdge(
		xslices.SubClone(e.label, 0, splitIdx),
		newLeafNode[K, V](value),
	)

	e.label = xslices.SubClone(e.label, splitIdx, len(e.label))

	e2 := newEdge(
		xslices.SubClone(s, splitIdx, len(s)),
		newLeafNode[K, V](value),
	)

	setEdge(n.edges, e1)
	setEdge(e1.node.edges, e)
	setEdge(e1.node.edges, e2)
	return
}

func (t *RadixTrieMap[K, V]) insert(key []K, value V) {
	var zero V
	if t.root == nil {

		t.root = newNode[K, V](zero, false)

	}
	fmt.Println("toi day: ", key)
	insert(t.root, key, value)

}

func (t *RadixTrieMap[K, V]) replaceOrInsert(key []K, value V) (_ V, _ bool) {
	var zero V
	if t.root == nil {

		t.root = newNode[K, V](zero, false)

	}

	v, ok := replaceOrInsert(t.root, key, value)
	if ok {
		return v, ok
	}

	t.len++
	return v, false

}

func (t *RadixTrieMap[K, V]) ReplaceOrInsert(key []K, value V) (_ V, _ bool) {

	return t.replaceOrInsert(key, value)

}

func (t *RadixTrieMap[K, V]) Len() int {

	return t.len

}
