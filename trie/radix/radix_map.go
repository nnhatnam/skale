package radix

import (
	"github.com/nnhatnam/skale/exp/xslices"
	"github.com/nnhatnam/skale/list/stack"
	"github.com/nnhatnam/skale/trie"
	"golang.org/x/exp/slices"
	"sort"
)

type block[K trie.Elem, V any] struct {
	label    []K            // label of the edge
	value    V              // value of the node
	lastElem bool           // true if this is the last element in the path
	next     []*block[K, V] // next blocks, next is always sorted
}

func newBlock[K trie.Elem, V any](label []K, value V, lastElem bool) *block[K, V] {
	return &block[K, V]{
		label:    label,
		value:    value,
		lastElem: lastElem,
	}
}

func (b *block[K, V]) nextBlock(e K) *block[K, V] {
	//perform binary search
	idx := sort.Search(len(b.next), func(i int) bool {
		return b.next[i].label[0] >= e
	})

	if idx < len(b.next) && b.next[idx].label[0] == e {
		return b.next[idx]
	}

	return nil
}

func (b *block[K, V]) insertBlock(label []K, value V, lastElem bool) *block[K, V] {

	if b.next == nil {
		b.next = make([]*block[K, V], 0, 1)
	}

	e := label[0]

	//perform binary search
	idx := sort.Search(len(b.next), func(i int) bool {
		return b.next[i].label[0] >= e
	})

	if idx < len(b.next) && b.next[idx].label[0] == e {
		panic("block already exists")
		return nil
	}

	//insert new block
	b.next = append(b.next, nil)
	copy(b.next[idx+1:], b.next[idx:])
	b.next[idx] = newBlock(label, value, lastElem)

	return b.next[idx]
}

func (b *block[K, V]) splitAt(idx int) *block[K, V] {

	//split label
	newLabel := xslices.SubClone(b.label, idx, len(b.label))
	b.label = xslices.SubClone(b.label, 0, idx)

	//split value
	nBlock := newBlock[K, V](newLabel, b.value, b.lastElem)

	//split next
	nBlock.next = b.next
	b.next = nil
	b.next = append(b.next, nBlock)

	return nBlock
}

func (b *block[K, V]) mergeChild() {

	if len(b.next) == 1 && !b.lastElem {
		child := b.next[0]

		//merge label
		newLabel := make([]K, len(b.label)+len(child.label))
		copy(newLabel, b.label)
		copy(newLabel[len(b.label):], child.label)
		b.label = newLabel
		//b.label = append(b.label, b.next[0].label...)

		//merge value
		b.value = child.value
		b.lastElem = child.lastElem

		//merge next
		b.next = child.next
		child.next = nil
	}

}

func (b *block[K, V]) removeBlock(e K) {

	//perform binary search
	idx := sort.Search(len(b.next), func(i int) bool {
		return b.next[i].label[0] >= e
	})

	if idx < len(b.next) && b.next[idx].label[0] == e {
		copy(b.next[idx:], b.next[idx+1:])
		b.next[len(b.next)-1] = nil
		b.next = b.next[:len(b.next)-1]
	}
}

type RadixMap[K trie.Elem, V any] struct {
	root *block[K, V]
	len  int // number of keys in the trie
}

func NewRadixMap[K trie.Elem, V any]() *RadixMap[K, V] {
	var zero V
	return &RadixMap[K, V]{
		root: newBlock[K, V]([]K{}, zero, false),
	}
}

func (t *RadixMap[K, V]) lazyInit() {
	if t.root == nil {
		var zero V
		t.root = newBlock[K, V]([]K{}, zero, false)
	}
}

// walkPath returns the path of blocks that contains the longest prefix of the key
// and the indexes of the last block element and the key element where the mismatch occurs.
func (t *RadixMap[K, V]) walkPath(key []K) (_ *stack.StackS[*block[K, V]], bIdx, kIdx int) {
	s := stack.NewStackS[*block[K, V]]()
	b := t.root
	i := 0
	s.Push(b)
	b = b.nextBlock(key[i])

	if b == nil {
		// root has no children
		return s, 0, 0
	}
	s.Push(b)

	if b != nil {

		// find the block that contains the longest prefix of the key

		for i < len(key) && len(key[i:]) >= len(b.label) && slices.Equal(key[i:i+len(b.label)], b.label) {
			i += len(b.label)
			if i >= len(key) {
				break
			}
			b = b.nextBlock(key[i])
			if b == nil {
				return s, len(s.Peek().label), i
			}
			s.Push(b)

		}
	}

	if i == len(key) {
		return s, len(b.label), i
	}
	bIdx = xslices.MismatchIndex(key[i:], b.label)
	kIdx = i + bIdx
	return s, bIdx, kIdx
}

func (t *RadixMap[K, V]) replaceOrInsert(s []K, value V) (_ V, _ bool) {

	path, bIdx, kIdx := t.walkPath(s)
	defer path.Clear()
	b := path.Peek()

	if kIdx == len(s) && bIdx == len(b.label) && b.lastElem {
		old := b.value
		b.value = value
		return old, true
	}

	if bIdx == len(b.label) {
		// insert new block
		b.insertBlock(s[kIdx:], value, true)
		return
	}

	// split block
	//nBlock := b.splitAt(bIdx)
	b.splitAt(bIdx)
	b.lastElem = false

	if kIdx == len(s) {
		//nBlock.insertBlock([]K{}, value, true)
		b.lastElem = true
		return
	}

	// insert new block
	b.insertBlock(s[kIdx:], value, true)
	return
}

func (t *RadixMap[K, V]) replaceOrInsertExtend(s []K, value V, replaceFunc ReplaceFunc[K, V]) (_ V, _ bool) {

	path, bIdx, kIdx := t.walkPath(s)
	defer path.Clear()
	b := path.Peek()

	if kIdx == len(s) && bIdx == len(b.label) && b.lastElem {
		old := b.value
		if replaceFunc != nil {
			b.value = replaceFunc(old)
			return old, true
		}
		b.value = value
		return old, true
	}

	if bIdx == len(b.label) {
		// insert new block
		b.insertBlock(s[kIdx:], value, true)
		return
	}

	// split block
	//nBlock := b.splitAt(bIdx)
	b.splitAt(bIdx)
	b.lastElem = false

	if kIdx == len(s) {
		//nBlock.insertBlock([]K{}, value, true)
		b.lastElem = true
		return
	}

	// insert new block
	b.insertBlock(s[kIdx:], value, true)
	return
}

// requires: len(s) > 0
func (t *RadixMap[K, V]) delete(s []K) (_ V, _ bool) {

	path, bIdx, kIdx := t.walkPath(s)
	defer path.Clear()

	b := path.Peek()
	var zero V

	if kIdx == len(s) && bIdx == len(b.label) && b.lastElem {
		old := b.value
		b.value = zero
		if len(b.next) == 0 {
			// delete block
			path.Pop()
			k := b.label[0]
			b = path.Peek()
			b.removeBlock(k)
			if len(b.next) == 1 && b != t.root {
				b.mergeChild()
			}

		} else if len(b.next) == 1 {

			b.lastElem = false
			b.mergeChild()
		} else {
			// mark block as non-last
			b.lastElem = false
		}

		return old, true
	}

	return
}

// requires: len(s) > 0
func (t *RadixMap[K, V]) markDelete(s []K) (_ V, _ bool) {

	path, bIdx, kIdx := t.walkPath(s)
	defer path.Clear()
	b := path.Peek()
	var zero V

	if kIdx == len(s) && bIdx == len(b.label) && b.lastElem {

		old := b.value
		b.value = zero
		b.lastElem = false

		return old, true
	}

	return
}

func (t *RadixMap[K, V]) shrink(b *block[K, V]) (b1 *block[K, V]) {

	if b == nil {
		return nil
	}

	for i, _ := range b.next {
		b.next[i] = t.shrink(b.next[i])
	}

	b2 := b.next[:0]

	//filter out nil
	for _, b3 := range b.next {
		if b3 != nil {
			b2 = append(b2, b3)
		}
	}

	b.next = b2

	if b != t.root && len(b.next) == 0 && !b.lastElem {
		return nil
	}

	if b != t.root && len(b.next) == 1 {
		b.mergeChild()
	}

	return b
}

func (t *RadixMap[K, V]) countChild(b *block[K, V]) int {
	if b == nil {
		return 0
	}
	return len(b.next)
}

// require: len(s) > 0
func (t *RadixMap[K, V]) deletePrefix(prefix []K) (t1 *RadixMap[K, V]) {

	path, bIdx, kIdx := t.walkPath(prefix)
	defer path.Clear()

	b := path.Peek()
	if kIdx != len(prefix) {
		return nil
	}

	var zero V

	// delete all children
	subLen := 0
	t.ascend(b, func(prefix []K, value V) bool {
		subLen++
		return false
	})

	k := b.label[0]
	b1 := b // save b
	path.Pop()
	b = path.Peek()
	b.removeBlock(k)
	t.len -= subLen

	// create new tree
	t1 = NewRadixMap[K, V]()
	t1.root = newBlock[K, V]([]K{}, zero, false)

	label := make([]K, len(prefix)+len(b1.label[bIdx:]))
	copy(label, prefix)
	copy(label[len(prefix):], b1.label[bIdx:])
	b1.label = label

	t1.root.next = append(t1.root.next, b1)
	t1.len = subLen

	if b != t.root {
		b.mergeChild()
	}

	return
}

func (t *RadixMap[K, V]) ascend(b *block[K, V], iterator ItemIterator[K, V]) {

	if b == nil {
		return // stop iteration
	}

	var inOrderRecursive func(b *block[K, V], prefix []K) bool
	inOrderRecursive = func(b *block[K, V], prefix []K) bool {

		if b.lastElem {

			if iterator(prefix, b.value) {
				return true // stop iteration
			}
		}

		for _, child := range b.next {

			if inOrderRecursive(child, append(prefix, child.label...)) {
				return true // stop iteration
			}

		}

		return false
	}

	inOrderRecursive(b, b.label)
}

func (t *RadixMap[K, V]) ascendGreaterOrEqual(b *block[K, V], greaterOrEqual []K, iterator ItemIterator[K, V]) {

	if b == nil {
		return // stop iteration
	}

	var inOrderRecursive func(b *block[K, V], prefix []K) bool

	inOrderRecursive = func(b *block[K, V], prefix []K) bool {

		if b.lastElem {

			if b != t.root || len(greaterOrEqual) == 0 {
				if iterator(prefix, b.value) {
					return true // stop iteration
				}
			}
		}

		for _, child := range b.next {

			j := len(prefix) + len(child.label)
			if j > len(greaterOrEqual) {
				j = len(greaterOrEqual)
			}

			if len(prefix) >= len(greaterOrEqual) || slices.Compare(child.label, greaterOrEqual[len(prefix):j]) >= 0 {

				if inOrderRecursive(child, append(prefix, child.label...)) {
					return true // stop
				}
			}

		}
		return false // continue iteration
	}

	inOrderRecursive(b, []K{}) // continue iteration
}

func (t *RadixMap[K, V]) ascendLessThan(b *block[K, V], lessThan []K, iterator ItemIterator[K, V]) {

	if b == nil || len(lessThan) == 0 {
		return // stop iteration
	}

	var inOrderRecursive func(b *block[K, V], prefix []K) bool

	inOrderRecursive = func(b *block[K, V], prefix []K) bool {

		if b.lastElem {
			if iterator(prefix, b.value) {
				return true // stop iteration
			}
		}

		for _, child := range b.next {

			j := len(prefix) + len(child.label)
			if j > len(lessThan) {
				j = len(lessThan)
			}

			if len(prefix) >= len(lessThan) || slices.Compare(child.label, lessThan[len(prefix):j]) >= 0 {
				return true // stop
			}

			if inOrderRecursive(child, append(prefix, child.label...)) {
				return true // stop iteration
			}

		}
		return false // continue iteration
	}

	inOrderRecursive(b, []K{})
}

func (t *RadixMap[K, V]) ascendRange(b *block[K, V], greaterOrEqual []K, lessThan []K, iterator ItemIterator[K, V]) {

	if b == nil || len(lessThan) == 0 || slices.Compare(greaterOrEqual, lessThan) >= 0 {
		return // stop iteration
	}

	var inOrderRecursive func(b *block[K, V], prefix []K) bool

	inOrderRecursive = func(b *block[K, V], prefix []K) bool {

		if b.lastElem {

			if b != t.root || len(greaterOrEqual) == 0 {
				if iterator(prefix, b.value) {
					return true // stop iteration
				}
			}
		}

		for _, child := range b.next {

			j := len(prefix) + len(child.label)
			if j > len(lessThan) {
				j = len(lessThan)
			}

			if len(prefix) >= len(lessThan) || slices.Compare(child.label, lessThan[len(prefix):j]) >= 0 {
				return true // stop
			}

			if len(prefix) >= len(greaterOrEqual) || slices.Compare(child.label, greaterOrEqual[len(prefix):j]) >= 0 {
				if inOrderRecursive(child, append(prefix, child.label...)) {
					return true // stop
				}
			}

		}
		return false // continue iteration
	}

	inOrderRecursive(b, []K{}) // continue iteration

}

func (t *RadixMap[K, V]) descend(b *block[K, V], iterator ItemIterator[K, V]) {

	if b == nil {
		return // stop iteration
	}

	var inOrderRecursive func(b *block[K, V], prefix []K) bool

	inOrderRecursive = func(b *block[K, V], prefix []K) bool {

		for i := len(b.next) - 1; i >= 0; i-- {

			if inOrderRecursive(b.next[i], append(prefix, b.next[i].label...)) {
				return true // stop iteration
			}
		}

		if b.lastElem {

			if iterator(prefix, b.value) {
				return true // stop iteration
			}
		}

		return false
	}

	inOrderRecursive(b, []K{})
}

func (t *RadixMap[K, V]) descendGreaterThan(b *block[K, V], greaterThan []K, iterator ItemIterator[K, V]) {

	if b == nil || len(greaterThan) == 0 {
		return // stop iteration
	}

	var inOrderRecursive func(b *block[K, V], prefix []K) bool

	inOrderRecursive = func(b *block[K, V], prefix []K) bool {

		//descendGreaterThan can't reach empty key, so we always ignore root value.
		if b.lastElem && b != t.root {

			if iterator(prefix, b.value) {
				return true // stop iteration
			}
		}

		for i := len(b.next) - 1; i >= 0; i-- {

			j := len(prefix) + len(b.next[i].label)
			if j > len(greaterThan) {
				j = len(greaterThan)
			}

			if len(prefix) >= len(greaterThan) || slices.Compare(b.next[i].label, greaterThan[len(prefix):j]) <= 0 {
				return true // stop
			}

			if inOrderRecursive(b.next[i], append(prefix, b.next[i].label...)) {
				return true // stop iteration
			}
		}

		return false
	}

	inOrderRecursive(b, []K{})

}

func (t *RadixMap[K, V]) descendLessOrEqual(b *block[K, V], lessOrEqual []K, iterator ItemIterator[K, V]) {

	if b == nil {
		return // stop iteration
	}

	var inOrderRecursive func(b *block[K, V], prefix []K) bool

	inOrderRecursive = func(b *block[K, V], prefix []K) bool {

		for _, child := range b.next {

			j := len(prefix) + len(child.label)
			if j > len(lessOrEqual) {
				j = len(lessOrEqual)
			}

			if len(prefix) >= len(lessOrEqual) || slices.Compare(child.label, lessOrEqual[len(prefix):j]) > 0 {
				return true // stop
			}

			if inOrderRecursive(child, append(prefix, child.label...)) {
				return true // stop iteration
			}

		}

		if b.lastElem {
			if iterator(prefix, b.value) {
				return true // stop iteration
			}
		}

		return false // continue iteration
	}

	inOrderRecursive(b, []K{})

}

func (t *RadixMap[K, V]) descendRange(b *block[K, V], greaterThan []K, lessOrEqual []K, iterator ItemIterator[K, V]) {

	if b == nil || slices.Compare(greaterThan, lessOrEqual) >= 0 {
		return // stop iteration
	}

	var inOrderRecursive func(b *block[K, V], prefix []K) bool

	inOrderRecursive = func(b *block[K, V], prefix []K) bool {

		for _, child := range b.next {

			j := len(prefix) + len(child.label)
			if j > len(lessOrEqual) {
				j = len(lessOrEqual)
			}

			if len(prefix) >= len(lessOrEqual) || slices.Compare(child.label, lessOrEqual[len(prefix):j]) > 0 {
				return true // stop
			}

			if len(prefix) >= len(greaterThan) || slices.Compare(child.label, greaterThan[len(prefix):j]) >= 0 {
				if inOrderRecursive(child, append(prefix, child.label...)) {
					return true // stop
				}
			}

			if b.lastElem && b != t.root {

				if iterator(prefix, b.value) {
					return true // stop iteration
				}
			}

		}
		return false // continue iteration
	}

	inOrderRecursive(b, []K{}) // continue iteration

}

func (t *RadixMap[K, V]) Get(key []K) (_ V, _ bool) {

	if t.root == nil || t.len == 0 {
		return
	}

	if len(key) == 0 {
		if t.root.lastElem {
			return t.root.value, true
		}
		return
	}

	path, bIdx, kIdx := t.walkPath(key)
	defer path.Clear()

	b := path.Peek()

	if kIdx == len(key) && bIdx == len(b.label) && b.lastElem {
		return b.value, true
	}

	return
}

func (t *RadixMap[K, V]) ReplaceOrInsert(key []K, value V) (_ V, _ bool) {
	var zero V
	t.lazyInit()

	if len(key) == 0 {

		if t.root.lastElem {
			old := t.root.value
			t.root.value = value
			return old, true
		}

		t.root.value = value
		t.root.lastElem = true
		t.len++

		return zero, false
	}
	v, ok := t.replaceOrInsert(key, value)
	if ok {
		return v, ok
	}

	t.len++
	return v, false
}

func (t *RadixMap[K, V]) ReplaceOrInsertExtend(key []K, value V, replaceFunc ReplaceFunc[K, V]) (_ V, _ bool) {
	var zero V
	t.lazyInit()

	if len(key) == 0 {

		if t.root.lastElem {
			old := t.root.value
			if replaceFunc != nil {
				t.root.value = replaceFunc(old)
				return old, true
			}
			t.root.value = value
			return old, true
		}

		t.root.value = value
		t.root.lastElem = true
		t.len++

		return zero, false
	}
	v, ok := t.replaceOrInsertExtend(key, value, replaceFunc)
	if ok {
		return v, ok
	}

	t.len++
	return v, false
}

func (t *RadixMap[K, V]) MarkDelete(key []K) (_ V, _ bool) {

	if t.root == nil || t.len == 0 {
		return
	}

	var zero V

	if len(key) == 0 {

		if old, found := t.root.value, t.root.lastElem; found {
			t.root.lastElem, t.root.value = false, zero
			t.len--
			return old, true
		}

		return
	}

	old, found := t.markDelete(key)

	if found {
		t.len--
	}

	return old, found
}

func (t *RadixMap[K, V]) Delete(key []K) (_ V, _ bool) {

	if t.root == nil || t.len == 0 {
		return
	}

	if len(key) == 0 {

		if t.root.lastElem {
			old := t.root.value
			var zero V
			t.root.lastElem = false
			t.root.value = zero
			t.len--
			return old, true
		}

		return
	}

	old, found := t.delete(key)

	if found {
		t.len--
	}

	return old, found
}

func (t *RadixMap[K, V]) DeletePrefix(prefix []K) (_ *RadixMap[K, V]) {

	if t.root == nil || t.len == 0 {
		return
	}

	if len(prefix) == 0 {
		t1 := NewRadixMap[K, V]()
		t1.lazyInit()
		t.root, t1.root = t1.root, t.root
		t.len, t1.len = t1.len, t.len
		return t1
	}

	return t.deletePrefix(prefix)
}

func (t *RadixMap[K, V]) Len() int {
	return t.len
}

func (t *RadixMap[K, V]) Ascend(iterator ItemIterator[K, V]) {

	if t.root == nil || t.len == 0 {
		return
	}

	t.ascend(t.root, iterator)
}

func (t *RadixMap[K, V]) AscendPrefix(prefix []K, iterator ItemIterator[K, V]) {

	path, bIdx, kIdx := t.walkPath(prefix)
	defer path.Clear()

	b := path.Peek()

	// there is no match for the prefix
	if b == t.root || kIdx < len(prefix) {
		return
	}

	// there is at least one match for the prefix
	matchedPrefix := prefix[0 : kIdx-bIdx]
	t.ascend(b, func(p []K, value V) bool {
		if iterator(append(matchedPrefix, p...), value) {
			return true
		}
		return false
	})

}

func (t *RadixMap[K, V]) AscendGreaterOrEqual(greaterOrEqual []K, iterator ItemIterator[K, V]) {

	if t.root == nil || t.len == 0 {
		return
	}

	t.ascendGreaterOrEqual(t.root, greaterOrEqual, iterator)

}

func (t *RadixMap[K, V]) AscendLessThan(lessThan []K, iterator ItemIterator[K, V]) {

	if t.root == nil || t.len == 0 {
		return
	}

	t.ascendLessThan(t.root, lessThan, iterator)

}

func (t *RadixMap[K, V]) AscendRange(greaterThan []K, lessOrEqual []K, iterator ItemIterator[K, V]) {

	if t.root == nil || t.len == 0 {
		return
	}

	t.ascendRange(t.root, greaterThan, lessOrEqual, iterator)

}

func (t *RadixMap[K, V]) Descend(iterator ItemIterator[K, V]) {

	if t.root == nil || t.len == 0 {
		return
	}

	t.descend(t.root, iterator)
}

func (t *RadixMap[K, V]) DescendLessOrEqual(lessOrEqual []K, iterator ItemIterator[K, V]) {

	if t.root == nil || t.len == 0 {
		return
	}

	t.descendLessOrEqual(t.root, lessOrEqual, iterator)

}

func (t *RadixMap[K, V]) DescendGreaterThan(greaterThan []K, iterator ItemIterator[K, V]) {

	if t.root == nil || t.len == 0 {
		return
	}

	t.descendGreaterThan(t.root, greaterThan, iterator)

}

func (t *RadixMap[K, V]) DescendRange(greaterThan []K, lessOrEqual []K, iterator ItemIterator[K, V]) {

	if t.root == nil || t.len == 0 {
		return
	}

	t.descendRange(t.root, greaterThan, lessOrEqual, iterator)

}

func (t *RadixMap[K, V]) LongestPrefix(prefix []K) (longest []K, value V, hasValue bool) {

	path, bIdx, kIdx := t.walkPath(prefix)
	defer path.Clear()

	b := path.Peek()

	// there is no match for the prefix
	if b == t.root {
		return
	}

	for kIdx >= 0 {
		// there is at least one match for the prefix

		if bIdx == len(b.label) {
			longest = prefix[0:kIdx]
			if b.lastElem {
				value = b.value
				hasValue = true
				return
			}
			kIdx -= len(b.label)
			path.Pop()
			b = path.Peek()
			continue
		}

		kIdx -= bIdx
		path.Pop()
		b = path.Peek()
		bIdx = len(b.label)
	}

	return
}

func (t *RadixMap[K, V]) Min() (_ []K, _ V, _ bool) {
	if t.root == nil || t.len == 0 {
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

func (t *RadixMap[K, V]) Max() (_ []K, _ V, _ bool) {

	if t.root == nil || t.len == 0 {
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

func (t *RadixMap[K, V]) Shrink() {
	t.root = t.shrink(t.root)
}
