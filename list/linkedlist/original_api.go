package linkedlist

// Init initializes or clears list l.
func (l *List[T]) Init() *List[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// lazyInit lazily initializes a zero List value.
func (l *List[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// pushFront inserts a new element e with value v at the front of list l and returns e.
func (l *List[T]) pushFront(v T) *Node[T] {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// pushBack inserts a new element e with value v at the back of list l and returns e.
func (l *List[T]) pushBack(v T) *Node[T] {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// insertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *List[T]) insertBefore(v T, mark *Node[T]) *Node[T] {

	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// insertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *List[T]) insertAfter(v T, mark *Node[T]) *Node[T] {

	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// moveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *List[T]) moveToFront(e *Node[T]) {
	if l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, &l.root)
}

// moveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *List[T]) moveToBack(e *Node[T]) {
	if l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, l.root.prev)
}

// moveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *List[T]) moveBefore(e, mark *Node[T]) {
	if e == mark {
		return
	}
	l.move(e, mark.prev)
}

// moveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *List[T]) moveAfter(e, mark *Node[T]) {
	if e == mark {
		return
	}
	l.move(e, mark)
}

// pushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List[T]) pushBackList(other *List[T]) {
	l.lazyInit()

	for i, e := other.Len(), other.front(); i > 0; i, e = i-1, e.next {
		l.insertValue(e.Value, l.root.prev)
	}
}

// pushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List[T]) pushFrontList(other *List[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.back(); i > 0; i, e = i-1, e.prev {
		l.insertValue(e.Value, &l.root)
	}
}
