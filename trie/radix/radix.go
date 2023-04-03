package radix

import (
	"github.com/nnhatnam/skale/exp/xslices"
	"github.com/nnhatnam/skale/trie"
	"golang.org/x/exp/slices"
)

type RadixTrieMap[K trie.Elem, V any] struct {
	root *node[K, V]

	eKey   bool // true if the trie contains an empty key
	eValue V    // value of the empty key

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

	if len(key) == 0 {
		return nil, n, -1, -1
	}

	i, j := 0, 0

	for i < len(key) {

		e = n.getEdge(key[i:])

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

	if root == nil {
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
		n.setEdge(e1)
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
	n1.setEdge(e2)

	return
}

func (t *RadixTrieMap[K, V]) ascendGreaterOrEqual(root *node[K, V], key []K, iterator ItemIterator[K, V]) bool {

	if root == nil {
		return true // stop
	}

	var inOrderRecursive func(n *node[K, V], prefix []K) bool // return true to stop

	inOrderRecursive = func(n *node[K, V], prefix []K) bool {

		if n.lastElem {
			if iterator(prefix, n.value) {
				return true // stop
			}

		}

		for _, e := range n.edges {

			if len(prefix) >= len(key) || slices.Compare(e.label, key[len(prefix):]) >= 0 {
				if inOrderRecursive(e.node, append(prefix, e.label...)) {
					return true // stop
				}
			}
			//inOrderRecursive(e.node, append(prefix, e.label...))
		}
		return false
	}

	return inOrderRecursive(root, []K{})

}

func (t *RadixTrieMap[K, V]) ascendLessThan(root *node[K, V], key []K, iterator ItemIterator[K, V]) bool {

	if root == nil {
		return true // stop
	}

	var inOrderRecursive func(n *node[K, V], prefix []K) bool // return true to stop

	inOrderRecursive = func(n *node[K, V], prefix []K) bool {

		if n.lastElem {
			if iterator(prefix, n.value) {
				return true // stop
			}

		}

		for _, e := range n.edges {

			if len(prefix) >= len(key) || slices.Compare(e.label, key[len(prefix):]) >= 0 {
				return false // stop
			}

			if inOrderRecursive(e.node, append(prefix, e.label...)) {
				return true // stop
			}

		}
		return false
	}

	return inOrderRecursive(root, []K{})

}

func (t *RadixTrieMap[K, V]) ascendRange(root *node[K, V], greaterOrEqual []K, lessThan []K, iterator ItemIterator[K, V]) bool {
	if root == nil {
		return true // stop
	}

	var inOrderRecursive func(n *node[K, V], prefix []K) bool // return true to stop

	inOrderRecursive = func(n *node[K, V], prefix []K) bool {

		if n.lastElem {
			if iterator(prefix, n.value) {
				return true // stop
			}

		}

		for _, e := range n.edges {

			if len(prefix) >= len(lessThan) || slices.Compare(e.label, lessThan[len(prefix):]) >= 0 {
				return true // stop
			}

			if len(prefix) >= len(greaterOrEqual) || slices.Compare(e.label, greaterOrEqual[len(prefix):]) >= 0 {
				if inOrderRecursive(e.node, append(prefix, e.label...)) {
					return true // stop
				}
			}
			//inOrderRecursive(e.node, append(prefix, e.label...))
		}
		return false
	}

	return inOrderRecursive(root, []K{})
}

func (t *RadixTrieMap[K, V]) ascend(root *node[K, V], iterator ItemIterator[K, V]) bool {
	if root == nil {
		return true // stop
	}

	var inOrderRecursive func(n *node[K, V], prefix []K) bool // return true to stop

	inOrderRecursive = func(n *node[K, V], prefix []K) bool {

		if n.lastElem {
			if iterator(prefix, n.value) {
				return true // stop
			}

		}

		for _, e := range n.edges {

			if inOrderRecursive(e.node, append(prefix, e.label...)) {
				return true // stop
			}
		}
		return false
	}

	return inOrderRecursive(root, []K{})
}

func (t *RadixTrieMap[K, V]) descend(root *node[K, V], iterator ItemIterator[K, V]) bool {

	if root == nil {
		return true // stop
	}

	var inOrderRecursive func(n *node[K, V], prefix []K) bool // return true to stop

	inOrderRecursive = func(n *node[K, V], prefix []K) bool {

		if n.lastElem {
			if iterator(prefix, n.value) {
				return true // stop
			}

		}

		for i := len(n.edges) - 1; i >= 0; i-- {

			if inOrderRecursive(n.edges[i].node, append(prefix, n.edges[i].label...)) {
				return true // stop
			}
		}
		return false
	}

	return inOrderRecursive(root, []K{})

}

func (t *RadixTrieMap[K, V]) descendRange(root *node[K, V], lessOrEqual []K, greaterThan []K, iterator ItemIterator[K, V]) bool {
	if root == nil {
		return true // stop
	}

	var inOrderRecursive func(n *node[K, V], prefix []K) bool // return true to stop

	inOrderRecursive = func(n *node[K, V], prefix []K) bool {

		for _, e := range n.edges {

			if len(prefix) > len(lessOrEqual) || slices.Compare(e.label, lessOrEqual[len(prefix):]) > 0 {
				return true // stop
			}

			if len(prefix) > len(greaterThan) || slices.Compare(e.label, greaterThan[len(prefix):]) > 0 {
				if inOrderRecursive(e.node, append(prefix, e.label...)) {
					return true // stop
				}
			}
			//inOrderRecursive(e.node, append(prefix, e.label...))
		}

		if n.lastElem {
			if iterator(prefix, n.value) {
				return true // stop
			}

		}

		return false
	}

	return inOrderRecursive(root, []K{})
}

func (t *RadixTrieMap[K, V]) descendLessOrEqual(root *node[K, V], key []K, iterator ItemIterator[K, V]) bool {
	if root == nil {
		return true // stop
	}

	var inOrderRecursive func(n *node[K, V], prefix []K) bool // return true to stop

	inOrderRecursive = func(n *node[K, V], prefix []K) bool {

		for _, e := range n.edges {

			if len(prefix) > len(key) || slices.Compare(e.label, key[len(prefix):]) > 0 {
				return false // stop
			}

			if inOrderRecursive(e.node, append(prefix, e.label...)) {
				return true // stop
			}

		}

		if n.lastElem {
			if iterator(prefix, n.value) {
				return true // stop
			}

		}

		return false
	}

	return inOrderRecursive(root, []K{})
}

func (t *RadixTrieMap[K, V]) descendGreaterThan(root *node[K, V], key []K, iterator ItemIterator[K, V]) bool {
	if root == nil {
		return true // stop
	}

	//var ret []K

	var inOrderRecursive func(n *node[K, V], prefix []K) bool // return true to stop

	inOrderRecursive = func(n *node[K, V], prefix []K) bool {

		for _, e := range n.edges {

			if len(prefix) > len(key) || slices.Compare(e.label, key[len(prefix):]) > 0 {
				if inOrderRecursive(e.node, append(prefix, e.label...)) {
					return true // stop
				}
			}
			//inOrderRecursive(e.node, append(prefix, e.label...))
		}

		if n.lastElem {
			if iterator(prefix, n.value) {
				return true // stop
				//ret = append(ret, prefix...)
			}

		}

		return false
	}

	return inOrderRecursive(root, []K{})
}

func (t *RadixTrieMap[K, V]) get(key []K) (_ V, _ bool) {

	e, _, _, eIdx := walk[K, V](t.root, key)

	if e == nil {
		return
	}

	if eIdx == len(e.label)-1 {
		return e.node.value, true
	}

	return
}

func (t *RadixTrieMap[K, V]) delete(key []K) (_ V, _ bool) {

	var e, prevEdge *edge[K, V]
	currNode := t.root
	prevNode := t.root
	i, j := 0, 0

	for i < len(key) {
		prevEdge = e
		e = currNode.getEdge(key[i:])

		if e == nil {

			return
		}

		for j = 0; j < len(e.label) && i < len(key); j++ {
			if key[i] != e.label[j] {
				return
			}
			i++
		}

		prevNode = currNode
		currNode = e.node
	}

	if currNode.lastElem {

		currNode.lastElem = false

		if len(currNode.edges) == 0 {
			prevNode.deleteEdge(e)

			if len(prevNode.edges) == 1 && prevNode != t.root && !prevNode.lastElem {
				prevEdge.node = prevNode.edges[0].node
				prevEdge.label = append(prevEdge.label, prevNode.edges[0].label...)
				prevNode.edges[0].node = nil
				prevNode.deleteEdge(prevNode.edges[0])
			}
			return currNode.value, true
		} else if len(currNode.edges) == 1 {

			e.node = currNode.edges[0].node
			e.label = append(e.label, currNode.edges[0].label...)
			currNode.edges[0].node = nil
			currNode.deleteEdge(currNode.edges[0])
			return currNode.value, true
		} else {
			return currNode.value, true
		}
	}

	return

}

//func (t *RadixTrieMap[K, V]) deletePrefix(prefix []K) (_ []K, _ []V) {
//
//	var e, prevEdge *edge[K, V]
//	currNode := t.root
//	prevNode := t.root
//	i, j := 0, 0
//
//	for i < len(prefix) {
//		prevEdge = e
//		e = currNode.getEdge(prefix[i:])
//
//		if e == nil {
//
//			return
//		}
//
//		for j = 0; j < len(e.label) && i < len(prefix); j++ {
//			if prefix[i] != e.label[j] {
//				return
//			}
//			i++
//		}
//
//		prevNode = currNode
//		currNode = e.node
//	}
//
//	var keys []K
//	var values []V
//
//	var inOrderRecursive func(n *node[K, V], prefix []K) bool // return true to stop
//
//	inOrderRecursive = func(n *node[K, V], prefix []K) bool {
//
//		for _, e := range n.edges {
//
//			if inOrderRecursive(e.node, append(prefix, e.label...)) {
//				return true // stop
//			}
//
//		}
//
//		if n.lastElem {
//			keys = append(keys, prefix)
//			values = append(values, n.value)
//		}
//
//		return false
//	}
//
//	inOrderRecursive(currNode, prefix)
//
//	if len(currNode.edges) == 1 {
//
//		e.node = currNode.edges[0].node
//		e.label = append(e.label, currNode.edges[0].label...)
//		currNode.edges[0].node = nil
//		currNode.deleteEdge(currNode.edges[0])
//	} else {
//		prevNode.deleteEdge(e)
//	}
//
//	return keys, values
//}

func (t *RadixTrieMap[K, V]) Delete(key []K) (_ V, _ bool) {

	if t == nil || t.root == nil || t.len == 0 {
		return
	}

	if len(key) == 0 {

		if t.eKey {
			t.eKey = false
			t.len--
			return t.eValue, true
		}

		return
	}

	old, found := t.delete(key)

	if found {
		t.len--
	}

	return old, found
}

func (t *RadixTrieMap[K, V]) Min() (_ []K, _ V, _ bool) {

	if t == nil || t.root == nil || t.len == 0 {
		return
	}

	var key []K
	var min V

	t.Ascend(func(k []K, v V) bool {
		key, min = k, v
		return true
	})

	return key, min, true
}

func (t *RadixTrieMap[K, V]) Max() (_ []K, _ V, _ bool) {

	if t == nil || t.root == nil || t.len == 0 {
		return
	}

	var key []K
	var max V

	t.Descend(func(k []K, v V) bool {
		key, max = k, v
		return true
	})

	return key, max, true
}

func (t *RadixTrieMap[K, V]) Get(key []K) (_ V, _ bool) {

	if t == nil || t.root == nil || t.len == 0 {
		return
	}

	if len(key) == 0 {
		return t.eValue, t.eKey
	}

	return t.get(key)
}

func (t *RadixTrieMap[K, V]) Descend(iterator ItemIterator[K, V]) {

	t.descend(t.root, iterator)
}

func (t *RadixTrieMap[K, V]) DescendRange(start []K, end []K, iterator ItemIterator[K, V]) {

	t.descendRange(t.root, start, end, iterator)

}

func (t *RadixTrieMap[K, V]) DescendLessOrEqual(key []K, iterator ItemIterator[K, V]) {

	t.descendLessOrEqual(t.root, key, iterator)

}

func (t *RadixTrieMap[K, V]) DescendGreaterThan(key []K, iterator ItemIterator[K, V]) {

	t.descendGreaterThan(t.root, key, iterator)

}

func (t *RadixTrieMap[K, V]) Ascend(iterator ItemIterator[K, V]) {

	t.ascend(t.root, iterator)

}

func (t *RadixTrieMap[K, V]) AscendRange(start []K, end []K, iterator ItemIterator[K, V]) {

	t.ascendRange(t.root, start, end, iterator)

}

func (t *RadixTrieMap[K, V]) AscendLessThan(key []K, iterator ItemIterator[K, V]) {

	t.ascendLessThan(t.root, key, iterator)
}

func (t *RadixTrieMap[K, V]) AscendGreaterOrEqual(key []K, iterator ItemIterator[K, V]) {

	if t == nil || t.root == nil || t.len == 0 {
		return
	}

	if len(key) == 0 && t.eKey {
		if iterator([]K{}, t.eValue) {
			return
		}
	}

	t.ascendGreaterOrEqual(t.root, key, iterator)

}

func (t *RadixTrieMap[K, V]) ReplaceOrInsert(key []K, value V) (_ V, _ bool) {
	var zero V
	if t.root == nil {

		t.root = newNode[K, V](zero, false)

	}

	if len(key) == 0 {

		if t.eKey {
			v := t.eValue
			t.eValue = value
			return v, true
		}

		t.eKey = true
		t.eValue = value
		t.len++

		return zero, false
	}

	v, ok := replaceOrInsert(t.root, key, value)
	if ok {
		return v, ok
	}

	t.len++
	return v, false

}

func (t *RadixTrieMap[K, V]) Len() int {

	return t.len

}
