package linkedlist

type List[T any] struct {
	root Node[T]
	len  int
}

type IterFunc func(v any)
type IterFuncWithIndex func(i int, v any)

func New[T any]() *List[T] {
	l := &List[T]{}
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// insert inserts a node after mark. The mask must not be nil.
func (l *List[T]) insert(n, at *Node[T]) *Node[T] {

	//n after at, n before at.next
	n.prev = at
	n.next = at.next

	//at before n
	at.next = n

	//n before at.next
	n.next.prev = n
	l.len++
	return n
}

// insertValue is a convenience wrapper for insert(&Node{Value: v}, at)
func (l *List[T]) insertValue(v T, mark *Node[T]) *Node[T] {
	return l.insert(NewNode(v), mark)
}

func (l *List[T]) remove(n *Node[T]) *Node[T] {

	//node before n is now before n.next
	n.prev.next = n.next

	//node after n is now after n.prev
	n.next.prev = n.prev
	l.len--
	return n
}

func From[T any](values ...T) *List[T] {
	l := New[T]()
	for _, v := range values {
		l.insertValue(v, l.root.prev)
	}
	return l
}

// Len returns the length of the linkedlist
func (l *List[T]) Len() int {
	return l.len
}

// first returns the first node in the linkedlist
func (l *List[T]) front() *Node[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// last returns the last node in the linkedlist
func (l *List[T]) back() *Node[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// First returns the first value in the linkedlist
func (l *List[T]) Front() T {
	return l.front().Value
}

// Last returns the last value in the linkedlist
func (l *List[T]) Back() T {

	return l.back().Value
}

func (l *List[T]) PushBack(v T) {
	l.insertValue(v, l.root.prev)
}

func (l *List[T]) PushFront(v T) {
	l.insertValue(v, &l.root)
}

func (l *List[T]) PopFront() (T, bool) {
	var zero T
	if l.len == 0 {
		return zero, false
	}
	n := l.front()
	l.remove(n)
	return n.Value, true
}

func (l *List[T]) PopBack() (T, bool) {
	var zero T
	if l.len == 0 {
		return zero, false
	}
	n := l.back()
	l.remove(n)
	return n.Value, true
}

//
//func (l *List) Begin() *Cursor {
//	return &Cursor{list: l, current: l.root.next}
//}
//
//func (l *List) End() *Cursor {
//	return &Cursor{list: l, current: l.root.prev}
//}
//
//func (l *List) Traverse(f IterFunc) {
//	for n := l.first(); n != l.last(); n = n.next {
//		f(n.Value)
//	}
//}
//
//func (l *List) RTraverse(f IterFunc) {
//	for n := l.last(); n != l.first(); n = n.prev {
//		f(n.Value)
//	}
//}
//
//func (l *List) TraverseWithIndex(f IterFuncWithIndex) {
//	i := 0
//	for n := l.first(); n != l.last(); n, i = n.next, i+1 {
//		f(i, n.Value)
//	}
//}
//
//func (l *List) RTraverseWithIndex(f IterFuncWithIndex) {
//	i := l.length - 1
//	for n := l.first(); n != l.last(); n, i = n.next, i-1 {
//		f(i, n.Value)
//	}
//}
