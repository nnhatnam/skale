package skiplist

import (
	"github.com/nnhatnam/skale"
	"math/rand"
	"time"
)

type Node[T any] struct {
	next []*Node[T]

	prev *Node[T]

	value T // value stored in the node
}

// must supply level
func newNode[T any](value T, level int) *Node[T] {
	return &Node[T]{value: value, next: make([]*Node[T], level+1, level+1)}
}

func (n *Node[T]) Next() *Node[T] {
	return n.next[0]
}

func (n *Node[T]) NextAt(level int) *Node[T] {
	return n.next[level]
}

func (n *Node[T]) Value() T {
	return n.value
}

type SkipList[T any] struct {
	root Node[T]

	less skale.LessFunc[T]

	maxLevel int
	p        float64

	fingers []*Node[T] // for faster search

	len int // number of elements in the list

}

func New[T any](maxLevel int, p float64, less skale.LessFunc[T]) *SkipList[T] {

	l := &SkipList[T]{
		maxLevel: maxLevel,
		p:        p,
		less:     less,
		fingers:  make([]*Node[T], maxLevel+1, maxLevel+1),
	}

	l.root.next = make([]*Node[T], maxLevel+1, maxLevel+1) // [0, maxLevel]
	l.root.prev = &l.root

	var i int
	for i = 0; i <= l.maxLevel; i++ {
		l.root.next[i] = &l.root
		l.fingers[i] = &l.root
	}

	return l
}

func NewOrdered[T skale.Ordered](maxLevel int, p float64) *SkipList[T] {
	return New[T](maxLevel, p, skale.Less[T]())
}

var generator = rand.New(rand.NewSource(time.Now().UnixNano()))

func (l *SkipList[T]) generateLevel() (level int) {

	for level = 0; level < l.maxLevel; level++ {

		if generator.Float64() > l.p {
			return
		}
	}

	return
}

// getPrevAndCache will cache the prev node list in the fingers, starting from the given node `from`.
// `form` will be the starting point for the search. If `from` is nil, the search will start from the head.
func (l *SkipList[T]) getPrevAndCacheFrom(value T, from *Node[T]) {

	prevNode, currNode := from, from // from is the starting point

	for i := l.maxLevel; i >= 0; i-- {

		if currNode.next[i] != nil && l.less(currNode.value, value) {
			prevNode = currNode
			currNode = currNode.Next()
		}

		l.fingers[i] = prevNode
		currNode = prevNode
	}

}

//func (l *SkipList[T]) getPrevAndCachev1(value T) (curr *Node[T]) {
//
//	var j int
//
//	if l.fingers[l.maxLevel] != &l.root && l.less(value, l.fingers[l.maxLevel].value) {
//		j = l.maxLevel
//	} else {
//		j = sort.Search(l.maxLevel, func(i int) bool {
//
//			if l.fingers[i].next[i] == &l.root || l.less(value, l.fingers[i].next[i].value) {
//				return true
//			}
//			return false
//
//		})
//
//		if j > 0 {
//			j--
//		}
//	}
//
//	curr = l.fingers[j]
//
//	// prevNode, currNode := from, from // from is the starting point
//
//	for i := j; i >= 0; i-- {
//
//		for curr.next[i] != &l.root && l.less(curr.next[i].value, value) {
//			curr = curr.next[i]
//		}
//
//		l.fingers[i] = curr
//
//	}
//
//	return
//
//}

func (l *SkipList[T]) getPrevAndCache(value T) (curr *Node[T]) {
	var j int
	for i := 0; i <= l.maxLevel; i++ {
		if l.fingers[i] == &l.root || l.less(l.fingers[i].value, value) {
			j = i
			break
		}
	}

	curr = l.fingers[j]

	// prevNode, currNode := from, from // from is the starting point

	for i := j; i >= 0; i-- {

		for curr.next[i] != &l.root && l.less(curr.next[i].value, value) {
			curr = curr.next[i]
		}

		l.fingers[i] = curr

	}

	return

}

func (l *SkipList[T]) insertNoReplace(v T, level int) {

	l.getPrevAndCache(v)

	n := newNode[T](v, level)

	n.prev = l.getPrevAndCache(v)
	n.next[0] = l.fingers[0].next[0]
	n.next[0].prev = n
	l.fingers[0].next[0] = n

	for i := 1; i < len(n.next); i++ {
		n.next[i] = l.fingers[i].next[i]
		l.fingers[i].next[i] = n
	}
	l.len++

}

func (l *SkipList[T]) replaceOrInsert(v T, level int) (_ T, _ bool) {

	curr := l.get(v)

	if curr != nil {
		old := curr.value
		curr.value = v
		return old, true
	}

	n := newNode[T](v, level)

	n.prev = l.fingers[0]
	n.next[0] = l.fingers[0].next[0]
	n.next[0].prev = n
	l.fingers[0].next[0] = n

	for i := 1; i < len(n.next); i++ {
		n.next[i] = l.fingers[i].next[i]
		l.fingers[i].next[i] = n
	}

	l.len++

	return

}

func (l *SkipList[T]) get(value T) *Node[T] {

	prev := l.getPrevAndCache(value)

	next := prev.next[0]

	if next == &l.root || l.less(value, next.value) {
		return nil
	}

	return next
}

func (l *SkipList[T]) delete(value T) *Node[T] {

	prev := l.getPrevAndCache(value)
	next := prev.next[0]

	if next == &l.root || l.less(value, next.value) {
		return nil // not found
	}

	l.len--

	l.fingers[0].next[0] = next.next[0]
	next.next[0].prev = l.fingers[0]
	next.prev = nil
	next.next[0] = nil

	for i := len(next.next) - 1; i > 0; i-- {
		l.fingers[i].next[i] = next.next[i]
		next.next[i] = nil
	}

	return next // found
}

func (l *SkipList[T]) ReplaceOrInsert(value T) (_ T, _ bool) {

	level := l.generateLevel()

	return l.replaceOrInsert(value, level)
}

func (l *SkipList[T]) InsertNoReplace(value T) {

	level := l.generateLevel()

	l.insertNoReplace(value, level)
}

func (l *SkipList[T]) Get(value T) (_ T, _ bool) {

	n := l.get(value)

	if n == nil {
		return
	}

	return n.value, true
}

func (l *SkipList[T]) Delete(value T) (_ T, _ bool) {

	n := l.delete(value)

	if n == nil {
		return
	}

	return n.value, true
}

// DeleteMin deletes the minimum value in the list and returns it. If no such value exists, returns (zero-value, false).
func (l *SkipList[T]) DeleteMin() (_ T, _ bool) {

	n := l.root.next[0]

	if n == &l.root {
		return
	}

	return l.delete(n.value).value, true
}

// DeleteMax deletes the maximum value in the list and returns it. If no such value exists, returns (zero-value, false).
func (l *SkipList[T]) DeleteMax() (_ T, _ bool) {

	n := l.root.prev

	if n == &l.root {
		return
	}

	return l.delete(n.value).value, true
}

func (l *SkipList[T]) Len() int {
	return l.len
}

func (l *SkipList[T]) Has(value T) bool {
	return l.get(value) != nil
}

// Ascend calls the iterator for every value in the list within the range [first, last], until iterator returns false.
func (l *SkipList[T]) Ascend(iter ItemIterator[T]) {
	for n := l.root.next[0]; n != &l.root; n = n.next[0] {
		if !iter(n.value) {
			return
		}
	}
}

// AscendRange calls the iterator for every value in the list within the range [first, last], until iterator returns false.
func (l *SkipList[T]) AscendRange(first, last T, iter ItemIterator[T]) {
	for n := l.getPrevAndCache(first).next[0]; n != &l.root && l.less(n.value, last); n = n.next[0] {
		if !iter(n.value) {
			return
		}
	}
}

// AscendGreaterOrEqual calls the iterator for every value in the list within the range [first, last], until iterator returns false.
func (l *SkipList[T]) AscendGreaterOrEqual(first T, iter ItemIterator[T]) {
	for n := l.getPrevAndCache(first).next[0]; n != &l.root; n = n.next[0] {
		if !iter(n.value) {
			return
		}
	}
}

// AscendLessThan calls the iterator for every value in the list within the range [first, last], until iterator returns false.
func (l *SkipList[T]) AscendLessThan(last T, iter ItemIterator[T]) {
	for n := l.root.next[0]; n != &l.root && l.less(n.value, last); n = n.next[0] {
		if !iter(n.value) {
			return
		}
	}
}

// Descend calls the iterator for every value in the list within the range [first, last], until iterator returns false.
func (l *SkipList[T]) Descend(iter ItemIterator[T]) {
	for n := l.root.prev; n != &l.root; n = n.prev {
		if !iter(n.value) {
			return
		}
	}
}

// DescendLessOrEqual calls the iterator for every value in the list within the range [pivot, first], until iterator returns false.
func (l *SkipList[T]) DescendLessOrEqual(pivot T, iter ItemIterator[T]) {

	n := l.getPrevAndCache(pivot)

	//in case of duplicate values, we need to start from the first one
	for n.next[0] != &l.root && !l.less(pivot, n.next[0].value) {
		n = n.next[0]
	}

	for ; n != &l.root; n = n.prev {
		if !iter(n.value) {
			return
		}
	}
}

// DescendGreaterThan calls the iterator for every value in the list within the range [last, pivot), until iterator returns false.
func (l *SkipList[T]) DescendGreaterThan(pivot T, iter ItemIterator[T]) {

	for n := l.root.prev; n != &l.root && l.less(pivot, n.value); n = n.prev {
		if !iter(n.value) {
			return
		}
	}

}

// DescendRange calls the iterator for every value in the list within the range [lessOrEqual, greaterThan), until iterator returns false.
func (l *SkipList[T]) DescendRange(lessOrEqual, greaterThan T, iter ItemIterator[T]) {

	n := l.getPrevAndCache(lessOrEqual)

	// in case we have multiple values equal to lessOrEqual, we need to find the last one
	for n.next[0] != &l.root && !l.less(lessOrEqual, n.next[0].value) {
		n = n.next[0]
	}

	for ; n != &l.root && l.less(greaterThan, n.value); n = n.prev {
		if !iter(n.value) {
			return
		}
	}
}

func (l *SkipList[T]) Max() (_ T, _ bool) {

	n := l.root.prev

	if n == &l.root {
		return
	}

	return n.value, true
}

func (l *SkipList[T]) Min() (_ T, _ bool) {

	n := l.root.next[0]

	if n == &l.root {
		return
	}

	return n.value, true
}
