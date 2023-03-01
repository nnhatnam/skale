package legacy

import (
	"fmt"
	"github.com/nnhatnam/skale"
	"math/rand"
	"time"
)

type Node[T any] struct {
	next []*Node[T]

	value T // value stored in the node
}

// must supply level
func newNode[T any](value T, level uint8) *Node[T] {
	return &Node[T]{value: value, next: make([]*Node[T], level, level)}
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
	head *Node[T]

	tail *Node[T]

	less skale.LessFunc[T]

	maxLevel uint8
	p        float64

	fingers []*Node[T] // for faster search

	size int // number of elements in the list
}

func New[T any](maxLevel uint8, p float64, less skale.LessFunc[T]) *SkipList[T] {
	var zero T

	l := &SkipList[T]{
		head:     newNode[T](zero, maxLevel),
		tail:     newNode[T](zero, maxLevel),
		maxLevel: maxLevel,
		p:        p,
		less:     less,
		fingers:  make([]*Node[T], maxLevel, maxLevel),
	}

	var i uint8
	for i = 0; i < l.maxLevel; i++ {
		l.head.next[i] = l.tail
		l.fingers[i] = l.head
	}

	return l
}

func NewOrdered[T skale.Ordered](maxLevel uint8, p float64) *SkipList[T] {
	return New[T](maxLevel, p, skale.Less[T]())
}

var generator = rand.New(rand.NewSource(time.Now().UnixNano()))

func (l *SkipList[T]) generateLevel() (level uint8) {

	for level = uint8(0); level < l.maxLevel; level++ {
		if generator.Float64() > l.p {
			return
		}
	}

	return
}

// getPrevAndCache will cache the prev node list in the fingers, starting from the given node `from`.
// `form` will be the starting point for the search. If `from` is nil, the search will start from the head.
func (l *SkipList[T]) getPrevAndCache(value T, from *Node[T]) {

	prevNode, currNode := from, from // from is the starting point

	for i := l.maxLevel - 1; i >= 0; i-- {

		if currNode.next[i] != nil && l.less(currNode.value, value) {
			prevNode = currNode
			currNode = currNode.Next()
		}

		l.fingers[i] = prevNode
		currNode = prevNode
	}

}

func (l *SkipList[T]) replaceOrInsert(value T, level uint8) {

	curr := l.head
	n := newNode[T](value, level)

	if curr == nil {

		for i := level; i >= 0; i-- {
			l.head.next[i] = n
		}

		l.size++
		return
	}

	l.getPrevAndCache(value, l.fingers[0])

	if !l.less(l.fingers[0].value, value) {
		// replace
		l.fingers[0].value = value
		return
	}

	l.size++
	// insert
	for i := level; i >= 0; i-- {
		n.next[i] = l.fingers[i].next[i]
		l.fingers[i].next[i] = n
	}

}

func (l *SkipList[T]) get(value T) *Node[T] {

	if l.less(l.fingers[0].value, value) {
		l.getPrevAndCache(value, l.fingers[0])
	} else {
		l.getPrevAndCache(value, l.head)
	}

	if l.less(value, l.fingers[0].value) {
		return nil
	}

	return l.fingers[0]
}

func (l *SkipList[T]) delete(value T) {

	if l.less(l.fingers[0].value, value) {
		l.getPrevAndCache(value, l.fingers[0])
	} else {
		l.getPrevAndCache(value, l.head)
	}

	if l.less(value, l.fingers[0].value) {
		return
	}

	l.size--

	for i := l.maxLevel; i >= 0; i-- {
		if l.fingers[i].next[i] != nil && !l.less(value, l.fingers[i].next[i].value) {
			l.fingers[i].next[i] = l.fingers[i].next[i].next[i]
		}
	}
}

func (l *SkipList[T]) ReplaceOrInsert(value T) {

	level := l.generateLevel()

	fmt.Println("will insert at level", level)

	l.replaceOrInsert(value, level)
}

func (l *SkipList[T]) Get(value T) *Node[T] {
	return l.get(value)
}
