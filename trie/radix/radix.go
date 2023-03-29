package radix

import (
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
func newInternalNode[K trie.Elem, V any]() *node[K, V] {

	return &node[K, V]{}

}

// newLeafNode creates a new leaf node with the given value.
// leaf nodes won't have any edges, and their prefix will be the same as the word they represent.
// leaf nodes will store values.
func newLeafNode[K trie.Elem, V any](value V) *node[K, V] {

	return &node[K, V]{value: value, lastElem: true}

}

func newNode[K trie.Elem, V any](value V) *node[K, V] {

	return &node[K, V]{value: value, lastElem: true}

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

	size int // number of keys in the trie
}

func NewRadixTrieMap[K trie.Elem, V any]() *RadixTrieMap[K, V] {

	return &RadixTrieMap[K, V]{root: newInternalNode[K, V]()}

}

func (t *RadixTrieMap[K, V]) lazyInit() {

	if t.root == nil {

		t.root = newInternalNode[K, V]()

	}

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

		setEdge(n.edges, e1)

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

		t.root = newNode[K, V](zero)

	}

	insert(t.root, key, value)

}
