package list

// Node is a node in a double linked list
type Node struct {
	next, prev *Node

	Value any
}

func NewNode(value any) *Node {
	return &Node{Value: value}
}

type List struct {
	root   Node
	length int
}

type IterFunc func(v any)
type IterFuncWithIndex func(i int, v any)

func New() *List {
	l := new(List)
	l.root.next = &l.root
	l.root.prev = &l.root
	l.length = 0
	return l
}

func From(values ...any) *List {
	l := New()
	for _, v := range values {
		l.insertValue(v, l.root.prev)
	}
	return l
}

func (l *List) increment() {
	l.length++
}

func (l *List) decrement() {
	l.length--
}

// Len returns the length of the list
func (l *List) Len() int {
	return l.length
}

// first returns the first node in the list
func (l *List) first() *Node {
	if l.length == 0 {
		return nil
	}
	return l.root.next
}

// last returns the last node in the list
func (l *List) last() *Node {
	if l.length == 0 {
		return nil
	}
	return l.root.prev
}

// First returns the first value in the list
func (l *List) First() any {
	return l.first().Value
}

// Last returns the last value in the list
func (l *List) Last() any {

	return l.last().Value
}

// insert inserts a node after mark. The mask must not be nil.
func (l *List) insert(n, at *Node) *Node {

	//n after at, n before at.next
	n.prev = at
	n.next = at.next

	//at before n
	at.next = n

	//n before at.next
	n.next.prev = n
	l.increment()
	return n
}

// insertValue is a convenience wrapper for insert(&Node{Value: v}, at)
func (l *List) insertValue(v any, mark *Node) *Node {
	return l.insert(NewNode(v), mark)
}

func (l *List) remove(n *Node) *Node {

	//node before n is now before n.next
	n.prev.next = n.next

	//node after n is now after n.prev
	n.next.prev = n.prev
	l.decrement()
	return n
}

func (l *List) PushBack(v any) {
	l.insertValue(v, l.root.prev)
}

func (l *List) PushFront(v any) {
	l.insertValue(v, &l.root)
}

func (l *List) PopFront() any {
	if l.length == 0 {
		return nil
	}
	n := l.first()
	l.remove(n)
	return n.Value
}

func (l *List) PopBack() any {
	if l.length == 0 {
		return nil
	}
	n := l.last()
	l.remove(n)
	return n.Value
}

func (l *List) Begin() *Cursor {
	return &Cursor{list: l, current: l.root.next}
}

func (l *List) End() *Cursor {
	return &Cursor{list: l, current: l.root.prev}
}

func (l *List) Traverse(f IterFunc) {
	for n := l.first(); n != l.last(); n = n.next {
		f(n.Value)
	}
}

func (l *List) RTraverse(f IterFunc) {
	for n := l.last(); n != l.first(); n = n.prev {
		f(n.Value)
	}
}

func (l *List) TraverseWithIndex(f IterFuncWithIndex) {
	i := 0
	for n := l.first(); n != l.last(); n, i = n.next, i+1 {
		f(i, n.Value)
	}
}

func (l *List) RTraverseWithIndex(f IterFuncWithIndex) {
	i := l.length - 1
	for n := l.first(); n != l.last(); n, i = n.next, i-1 {
		f(i, n.Value)
	}
}
