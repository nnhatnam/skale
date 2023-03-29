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

// edgesMap is a map from a key to a slice of edges. key is the first character of the edge's label.
type edges[K trie.Elem, V any] []*edge[K, V]

func getEdgeByPrefix[K trie.Elem, V any](m edges[K, V], prefix K) *edge[K, V] {
	//binary search for key position in the map
	idx := sort.Search(len(m), func(i int) bool {
		return m[i].label[0] >= prefix
	})

	if idx == -1 {
		return nil
	}

	if m[idx].label[0] != prefix {
		return nil
	}

	return (*edge[K, V])(m[idx])
}

func setEdge[K trie.Elem, V any](m edges[K, V], e *edge[K, V]) edges[K, V] {
	//binary search for key position in the map
	idx := sort.Search(len(m), func(i int) bool {
		return (m)[i].label[0] >= e.label[0]
	})

	if idx == -1 {
		m = append(m, e)
	} else if m[idx].label[0] != e.label[0] {
		m = slices.Insert(m, idx, e)
	} else {
		m[idx] = e
	}
	return m
}
