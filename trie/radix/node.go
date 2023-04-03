package radix

import (
	"github.com/nnhatnam/skale/trie"
	"golang.org/x/exp/slices"
	"sort"
)

//// insertionSortEdges sorts the edges in place using insertion sort traveling from last to first.
//func insertionSortEdges[K trie.Elem, V any](edges []*edge[K, V]) {
//	// In this radix tree implementation, the edges are sorted by their labels, and when we insert a new edge, we need to insert it in the correct position.
//	// Because the edges are sorted, so when we insert a new edge, we can use binary search to find the correct position to insert the new edge.
//	// Another approach is inserting the new edge at the end of the slice, and then use insertion sort to sort the slice.
//	// In the second option, the array is almost sorted, so the insertion sort will be very fast.
//	// Using the binary search approach may be faster, but it may complicate the  base, so we choose the insertion sort approach instead.
//
//	// insertion sort
//	for i := len(edges) - 1; i > 0; i-- {
//		for j := i; j < len(edges) && edges[j].label[0] < edges[j-1].label[0]; j++ {
//			edges[j], edges[j-1] = edges[j-1], edges[j]
//		}
//	}
//
//}

type edge[K trie.Elem, V any] struct {
	label []K         // label of the edge
	node  *node[K, V] // node at the end of the edge
}

// edges is a weak map from a key to an edge, and it is sorted by the key.
// key in edges is defined as the first element of the edge's label.
// With that definition, two different keys with the same prefix will be considered as the same key.
type edges[K trie.Elem, V any] []*edge[K, V]

func newEdge[K trie.Elem, V any](label []K, node *node[K, V]) *edge[K, V] {

	return &edge[K, V]{label: label, node: node}

}

type node[K trie.Elem, V any] struct {
	edges edges[K, V] // edges from this node, sorted by label

	lastElem bool // true if this node is a leaf

	//prefix []K // prefix of the word represented by this node

	value V // value of the word represented by this node

}

func newNode[K trie.Elem, V any](value V, lastElem bool) *node[K, V] {

	return &node[K, V]{value: value, lastElem: lastElem}

}

// newLeafNode creates a new leaf node with the given value.
// leaf nodes won't have any edges, and their prefix will be the same as the word they represent.
// leaf nodes will store values.
func newLeafNode[K trie.Elem, V any](value V) *node[K, V] {

	return newNode[K, V](value, true)

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

func (n *node[K, V]) getEdge(prefix []K) *edge[K, V] {
	//binary search for key position in the map

	if len(n.edges) == 0 {
		return nil
	}

	idx := sort.Search(len(n.edges), func(i int) bool {
		return n.edges[i].label[0] >= prefix[0]
	})

	if idx >= len(n.edges) {
		return nil
	}

	if n.edges[idx].label[0] != prefix[0] {
		return nil
	}

	return (*edge[K, V])(n.edges[idx])
}

func (n *node[K, V]) setEdge(e *edge[K, V]) {
	//binary search for key position in the map

	if len(n.edges) == 0 {
		n.edges = append(n.edges, e)
	}

	idx := sort.Search(len(n.edges), func(i int) bool {
		return n.edges[i].label[0] >= e.label[0]
	})

	if idx >= len(n.edges) {
		n.edges = append(n.edges, e)
	} else if n.edges[idx].label[0] != e.label[0] {
		n.edges = slices.Insert(n.edges, idx, e)
	} else {
		n.edges[idx] = e
	}

}

func (n *node[K, V]) deleteEdge(e *edge[K, V]) bool {
	//binary search for key position in the map

	if len(n.edges) == 0 {
		return false
	}

	idx := sort.Search(len(n.edges), func(i int) bool {
		return n.edges[i].label[0] >= e.label[0]
	})

	if idx >= len(n.edges) {
		return false
	}

	if n.edges[idx].label[0] != e.label[0] {
		return false
	}

	copy(n.edges[idx:], n.edges[idx+1:])
	n.edges[len(n.edges)-1] = nil // or the zero value of T
	n.edges = n.edges[:len(n.edges)-1]
	return true
}
