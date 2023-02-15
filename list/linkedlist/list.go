package linkedlist

type List[T any] struct {
	root Node[T]
	len  int
}

type IterFunc[T any] func(n *Node[T]) bool

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

// move moves e to next to at.
func (l *List[T]) move(e, at *Node[T]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

func (l *List[T]) remove(n *Node[T]) *Node[T] {

	//node before n is now before n.next
	n.prev.next = n.next

	//node after n is now after n.prev
	n.next.prev = n.prev
	n.next = nil // avoid memory leaks
	n.prev = nil // avoid memory leaks
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

// Front returns the first value in the linkedlist
func (l *List[T]) Front() T {
	return l.front().Value
}

// Back returns the last value in the linkedlist
func (l *List[T]) Back() T {

	return l.back().Value
}

func (l *List[T]) PushBack(v T) {
	l.lazyInit()
	l.insertValue(v, l.root.prev)
}

func (l *List[T]) PushFront(v T) {
	l.lazyInit()
	l.insertValue(v, &l.root)
}

func (l *List[T]) PopFront() (T, bool) {
	l.lazyInit()

	var zero T
	if l.len == 0 {
		return zero, false
	}
	n := l.front()
	l.remove(n)
	return n.Value, true
}

func (l *List[T]) PopBack() (T, bool) {
	l.lazyInit()

	var zero T
	if l.len == 0 {
		return zero, false
	}
	n := l.back()
	l.remove(n)
	return n.Value, true
}

func (l *List[T]) InsertBefore(v T, c *Cursor[T]) *Node[T] {
	if c.list != l || c.current == &c.list.root {
		return nil
	}
	return c.list.insertValue(v, c.current.prev)
}

func (l *List[T]) InsertAfter(v T, c *Cursor[T]) *Node[T] {
	if c.list != l || c.current == &c.list.root {
		return nil
	}
	return c.list.insertValue(v, c.current)
}

func (l *List[T]) RemoveCurrent(c *Cursor[T]) *Node[T] {

	if c.list != l || c.current == &c.list.root {
		return nil
	}
	n := c.current
	c.current = c.current.next
	c.list.remove(n)
	return n
}

func (l *List[T]) RemoveAfter(c *Cursor[T]) *Node[T] {
	if c.list != l || l.root.prev == c.current {
		return nil
	}
	return c.list.remove(c.current.next)

}

func (l *List[T]) RemoveBefore(c *Cursor[T]) *Node[T] {

	if c.list != l || l.root.next == c.current {
		return nil
	}
	return c.list.remove(c.current.prev)
}

func (l *List[T]) MoveToFront(c *Cursor[T]) {
	if c.list != l || l.root.next == c.current {
		return
	}

	l.move(c.current, &l.root)
}

func (l *List[T]) MoveToBack(c *Cursor[T]) {
	if c.list != l || l.root.prev == c.current {
		return
	}

	l.move(c.current, l.root.prev)
}

func (l *List[T]) MoveBefore(c, mark *Cursor[T]) {
	if c.list != l || c.current == mark.current || mark.list != l {
		return
	}
	l.move(c.current, mark.current.prev)
}

func (l *List[T]) MoveAfter(c, mark *Cursor[T]) {
	if c.list != l || c.current == mark.current || mark.list != l {
		return
	}
	l.move(c.current, mark.current)
}

func (l *List[T]) Traverse(f func(T)) {
	for n := l.front(); n != nil; n = n.next {
		f(n.Value)
	}
}

func (l *List[T]) PushBackList(other *List[T]) {
	l.lazyInit()

	back := other.BackCursor().Node()

	other.Cursor().Ascending(func(n *Node[T]) bool {

		l.insertValue(n.Value, l.root.prev)
		if n == back {
			return false
		}
		return true
	})

}

func (l *List[T]) PushFrontList(other *List[T]) {
	l.lazyInit()

	front := other.FrontCursor().Node()

	other.Cursor().Descending(func(n *Node[T]) bool {

		l.insertValue(n.Value, &l.root)
		if n == front {
			return false
		}
		return true
	})

}
