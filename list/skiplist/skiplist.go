package skiplist

import (
	"github.com/nnhatnam/skale"
	"math/rand"
	"time"
)

type node[T any] struct {
	next []*node[T]

	value *T // value stored in the node
}

// must supply level
func newNode[T any](value T, level uint8) *node[T] {
	return &node[T]{value: &value, next: make([]*node[T], level, level)}
}

func (n *node[T]) Next() *node[T] {
	return n.next[0]
}

func (n *node[T]) NextAt(level int) *node[T] {
	return n.next[level]
}

type SkipList[T any] struct {
	head *node[T]

	less skale.LessFunc[T]

	maxLevel uint8
	p        float64

	fingers []*node[T] // for faster search

	size int // number of elements in the list
}

func New[T any](maxLevel uint8, p float64, less skale.LessFunc[T]) *SkipList[T] {
	l := &SkipList[T]{
		head:     newNode[T](nil, maxLevel),
		maxLevel: maxLevel,
		p:        p,
		less:     less,
		fingers:  make([]*node[T], maxLevel, maxLevel),
	}

	var i uint8
	for i = 0; i < l.maxLevel; i++ {
		l.fingers[i] = l.head
	}

	return l
}

var generator = rand.New(rand.NewSource(time.Now().UnixNano()))

func (l *SkipList[T]) generateLevel() (level uint8) {

	for level = uint8(0); level < uint8(l.maxLevel); level++ {
		if generator.Float64() > l.p {
			return
		}
	}

	return
}

// cachePrev will cache the prev node list in the fingers
func (l *SkipList[T]) getPrevAndCache(value T, from *node[T]) {

	prevNode, currNode := from, from // from is the starting point

	for i := l.maxLevel; i >= 0; i-- {

		if currNode.next[i] != nil && l.less(*currNode.value, value) {
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

	if !l.less(*l.fingers[0].value, value) {
		// replace
		l.fingers[0].value = &value
		return
	}

	l.size++
	// insert
	for i := level; i >= 0; i-- {
		n.next[i] = l.fingers[i].next[i]
		l.fingers[i].next[i] = n
	}

}

func (l *SkipList[T]) delete(value T) {

	if l.less(*l.fingers[0].value, value) {
		l.getPrevAndCache(value, l.fingers[0])
	} else {
		l.getPrevAndCache(value, l.head)
	}

	if l.less(value, *l.fingers[0].value) {
		return
	}

	l.size--

	for i := l.maxLevel; i >= 0; i-- {
		if l.fingers[i].next[i] != nil && !l.less(value, *l.fingers[i].next[i].value) {
			l.fingers[i].next[i] = l.fingers[i].next[i].next[i]
		}
	}
}

func (l *SkipList[T]) Insert(value T) {

	level := l.generateLevel()

	l.replaceOrInsert(value, level)
}
