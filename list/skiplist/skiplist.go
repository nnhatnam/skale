package skiplist

import (
	"fmt"
	"github.com/nnhatnam/skale"
	"math/rand"
	"time"
)

// Node is an element of a skip list.
type Node[T any] struct {
	next []*Node[T]

	prev *Node[T]

	value T // value stored in the node
}

// must supply level
func newNode[T any](value T, level int) *Node[T] {
	return &Node[T]{value: value, next: make([]*Node[T], level+1, level+1)}
}

// Value returns the value stored in the node.
func (n *Node[T]) Value() T {
	return n.value
}

// SkipList represents a skip list.
type SkipList[T any] struct {
	root Node[T]

	less skale.LessFunc[T]

	maxLevel int
	p        float64

	modifiedFingers []*Node[T] // for faster insertion/deletion

	searchFingers []*Node[T] // for faster search

	len int // number of elements in the list

}

func New[T any](maxLevel int, p float64, less skale.LessFunc[T]) *SkipList[T] {

	l := &SkipList[T]{
		maxLevel:        maxLevel - 1,
		p:               p,
		less:            less,
		modifiedFingers: make([]*Node[T], maxLevel, maxLevel),
		searchFingers:   make([]*Node[T], maxLevel, maxLevel),
	}

	l.root.next = make([]*Node[T], maxLevel, maxLevel) // [0, maxLevel]
	l.root.prev = &l.root

	var i int
	for i = 0; i <= l.maxLevel; i++ {
		l.root.next[i] = &l.root
		l.modifiedFingers[i] = &l.root
		l.searchFingers[i] = &l.root
	}

	return l
}

// lessThan compares the value of two nodes x and y by following rule:
// if x is root, x = -inf, so x is always less than y
// if y is root, y = +inf, so x is always greater than y
// if x and y are not root, compare their values
// If x or y is the root node, the order is important.
func (l *SkipList[T]) lessThan(x, y *Node[T]) bool {
	if x == &l.root {
		return true
	}
	if y == &l.root {
		return false
	}
	return l.less(x.value, y.value)
}

func (l *SkipList[T]) lessThanL(x *Node[T], y T) bool {
	if x == &l.root {
		return true
	}
	return l.less(x.value, y)
}

func (l *SkipList[T]) lessThanR(x T, y *Node[T]) bool {
	if y == &l.root {
		return true
	}
	return l.less(x, y.value)
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

// searchPrevAndCache will update the SEARCH finger list and return the prev node.
func (l *SkipList[T]) searchPrevAndCache(v T) *Node[T] {

	var j int // j is where the search starts

	for i := 0; i <= l.maxLevel; i++ {
		j = i

		if l.lessThanL(l.searchFingers[i], v) && l.lessThanR(v, l.searchFingers[i].next[i]) {
			// finger < i and i < finger.next[i] (we don't want i <= finger.next[i])
			break
		}

	}

	if j == l.maxLevel && !l.lessThanL(l.searchFingers[j], v) {
		// reset if j at top and finger is NOT on v's right (j == maxLevel and v <= finger[maxLevel])
		l.searchFingers[j] = &l.root
	}

	prev := l.searchFingers[j]

	for i := j; i >= 0; i-- {

		for prev.next[i] != &l.root && l.less(prev.next[i].value, v) {
			prev = prev.next[i]
		}

		l.searchFingers[i] = prev

	}

	return l.searchFingers[0]
}

// getPrevAndCache will go up from level 0 until it finds the prev node less than v, and looks pass v.
func (l *SkipList[T]) getPrevAndCache(v T) *Node[T] {

	return l.getPrevAndCacheAtLevel(v, 0)

}

func (l *SkipList[T]) getPrevAndCacheAtLevel(v T, level int) *Node[T] {
	var j int // j is where the search starts

	for i := level; i <= l.maxLevel; i++ {
		j = i

		if l.lessThanL(l.modifiedFingers[i], v) && l.lessThanR(v, l.modifiedFingers[i].next[i]) {
			// finger < i and i < finger.next[i] (we don't want i <= finger.next[i])
			break
		}

	}

	if j == l.maxLevel && !l.lessThanL(l.modifiedFingers[j], v) {
		// reset if j at top and finger is NOT on v's right (j == maxLevel and v <= finger[maxLevel])
		l.modifiedFingers[j] = &l.root
	}

	prev := l.modifiedFingers[j]

	for i := j; i >= 0; i-- {

		for prev.next[i] != &l.root && l.less(prev.next[i].value, v) {
			prev = prev.next[i]
		}

		l.modifiedFingers[i] = prev

	}

	return l.modifiedFingers[0]

}

func (l *SkipList[T]) insertNoReplace(v T, level int) {

	prev := l.getPrevAndCacheAtLevel(v, level)

	n := newNode[T](v, level)

	n.prev = prev
	n.next[0] = prev.next[0]
	n.next[0].prev = n
	prev.next[0] = n

	l.len++

}

func (l *SkipList[T]) replaceOrInsert(v T, level int) (_ T, _ bool) {
	prev := l.getPrevAndCacheAtLevel(v, level)

	next := prev.next[0]

	if next == &l.root || l.less(v, next.value) {
		//v is not in the list, insert it

		n := newNode[T](v, level)

		n.next[0] = prev.next[0]
		n.next[0].prev = n
		n.prev = prev
		prev.next[0] = n

		for i := len(n.next) - 1; i > 0; i-- {

			n.next[i] = l.modifiedFingers[i].next[i]
			l.modifiedFingers[i].next[i] = n
		}

		l.len++

		return
	}

	//v is in the list, replace it
	old := next.value
	next.value = v
	return old, true
}

func (l *SkipList[T]) get(value T) *Node[T] {

	prev := l.searchPrevAndCache(value)

	next := prev.next[0]

	if next == &l.root || l.less(value, next.value) {
		return nil
	}

	return next
}

func (l *SkipList[T]) delete(value T) *Node[T] {

	prev := l.getPrevAndCache(value)
	curr := prev.next[0]

	if curr == &l.root || l.less(value, curr.value) {
		return nil // not found
	}

	l.len--

	prev.next[0] = curr.next[0]

	curr.next[0].prev = prev

	curr.prev = nil
	curr.next[0] = nil

	for i := len(curr.next) - 1; i > 0; i-- {

		prev = l.modifiedFingers[i]
		//fmt.Printf("at level %v, we link %v to %v\n", i, prev.value, curr.next[i].value)
		prev.next[i] = curr.next[i]

		curr.next[i] = nil

	}

	return curr // found
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

// Delete removes an item equal to the passed in item from the list, returning it. If no such item exists, returns (zeroValue, false).
func (l *SkipList[T]) Delete(value T) (_ T, _ bool) {

	n := l.delete(value)

	if n == nil {
		return
	}

	return n.value, true
}

// DeleteMin deletes the minimum value in the list and returns it. If no such value exists, returns (zeroValue, false).
func (l *SkipList[T]) DeleteMin() (_ T, _ bool) {

	n := l.root.next[0]

	if n == &l.root {
		return
	}

	return l.delete(n.value).value, true
}

// DeleteMax deletes the maximum value in the list and returns it. If no such value exists, returns (zeroValue, false).
func (l *SkipList[T]) DeleteMax() (_ T, _ bool) {

	n := l.root.prev

	if n == &l.root {
		return
	}

	return l.delete(n.value).value, true
}

// Len returns the number of items currently in the list.
func (l *SkipList[T]) Len() int {
	return l.len
}

// Has returns true if the given value is in the list
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

// AscendGreaterOrEqual calls the iterator for every value in the list within the range [pivot , last], until iterator returns false.
func (l *SkipList[T]) AscendGreaterOrEqual(pivot T, iter ItemIterator[T]) {
	for n := l.searchPrevAndCache(pivot).next[0]; n != &l.root; n = n.next[0] {
		if !iter(n.value) {
			return
		}
	}
}

// AscendLessThan calls the iterator for every value in the list within the range [first, pivot), until iterator returns false.
func (l *SkipList[T]) AscendLessThan(pivot T, iter ItemIterator[T]) {
	for n := l.root.next[0]; n != &l.root && l.less(n.value, pivot); n = n.next[0] {
		if !iter(n.value) {
			return
		}
	}
}

// AscendRange calls the iterator for every value in the list within the range [greaterOrEqual, lessThan) , until iterator returns false.
func (l *SkipList[T]) AscendRange(greaterOrEqual, lessThan T, iter ItemIterator[T]) {
	for n := l.searchPrevAndCache(greaterOrEqual).next[0]; n != &l.root && l.less(n.value, lessThan); n = n.next[0] {
		if !iter(n.value) {
			return
		}
	}
}

// Descend calls the iterator for every value in the list within the range [last, first], until iterator returns false.
func (l *SkipList[T]) Descend(iter ItemIterator[T]) {
	for n := l.root.prev; n != &l.root; n = n.prev {
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

// DescendLessOrEqual calls the iterator for every value in the list within the range [pivot, first], until iterator returns false.
func (l *SkipList[T]) DescendLessOrEqual(pivot T, iter ItemIterator[T]) {

	n := l.searchPrevAndCache(pivot)

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

// DescendRange calls the iterator for every value in the list within the range [lessOrEqual, greaterThan), until iterator returns false.
func (l *SkipList[T]) DescendRange(lessOrEqual, greaterThan T, iter ItemIterator[T]) {

	n := l.searchPrevAndCache(lessOrEqual)

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

// Max returns the largest item in the list, or (zeroValue, false) if the list is empty.
func (l *SkipList[T]) Max() (_ T, _ bool) {

	n := l.root.prev

	if n == &l.root {
		return
	}

	return n.value, true
}

// Min returns the smallest item in the list, or (zeroValue, false) if the list is empty.
func (l *SkipList[T]) Min() (_ T, _ bool) {

	n := l.root.next[0]

	if n == &l.root {
		return
	}

	return n.value, true
}

func (l *SkipList[T]) print() {
	for i := 0; i <= l.maxLevel; i++ {

		for n := l.root.next[i]; n != &l.root; n = n.next[i] {
			fmt.Printf("%v -> ", n.value)
		}

		fmt.Println()

	}
}
