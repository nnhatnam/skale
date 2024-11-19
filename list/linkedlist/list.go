// Package linkedlist implements a doubly linked list data structure.
// It is intended to be used internally by other packages.
//
// Structure is not thread safe.
// To iterate over a list (where l is a *List):
//
//	cursor := l.Cursor() // create a cursor point to first node
//	for cursor.MoveNext() {
//		v := cursor.Value()
//		// do something with v
//	}
//
// A better way to iterate over a list:
//
//	cursor := l.Cursor() // create a cursor point to the sentinel node
//	cursor.WalkAscending(func(v T) bool {
//		// do something with v
//		return true
//	})
package linkedlist

// node is a node in a doubly linked list
type node[T any] struct {
	next, prev *node[T]
	value      T
}

// newNode creates a new node with the given value
func newNode[T any](value T) *node[T] {
	return &node[T]{value: value}
}

// List represents a doubly linked list.
type List[T any] struct {
	root node[T]
	len  int
}

// New returns an initialized list.
func New[T any]() *List[T] {
	l := &List[T]{}
	l.root.next = &l.root
	l.root.prev = &l.root
	return l
}

// Init resets the list to an empty state
func (l *List[T]) Init() *List[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// isEmpty checks if the list is empty.
func (l *List[T]) isEmpty() bool { return l.len == 0 }

// first returns the first node in the list
func (l *List[T]) front() *node[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// last returns the last node in the list
func (l *List[T]) back() *node[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// insert inserts a node after mark. The mask must not be nil.
func (l *List[T]) insert(n, mark *node[T]) *node[T] {

	//n after mark, n before mark.next
	n.prev = mark
	n.next = mark.next

	//mark before n
	mark.next = n

	//n before mark.next
	n.next.prev = n
	l.len++
	return n
}

// insertValue is a convenience wrapper for insert(&node{Value: v}, at)
func (l *List[T]) insertValue(v T, mark *node[T]) *node[T] {
	return l.insert(newNode(v), mark)
}

// move moves e to next to at.
func (l *List[T]) move(e, at *node[T]) {
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

// remove removes n from the list. The node must not be nil.
func (l *List[T]) remove(n *node[T]) *node[T] {

	//node before n is now before n.next
	n.prev.next = n.next

	//node after n is now after n.prev
	n.next.prev = n.prev
	n.next, n.prev = nil, nil // Avoid memory leaks
	l.len--
	return n
}

// From returns an initialized list and add the given values, if any, to the list.
func From[T any](values ...T) *List[T] {
	l := New[T]()
	for _, v := range values {
		l.insertValue(v, l.root.prev)
	}
	return l
}

// Len returns the number of elements of list l. The complexity is O(1).
func (l *List[T]) Len() int { return l.len }

// Front returns the value of the first element or zero value if empty.
// The complexity is O(1).
func (l *List[T]) Front() (T, bool) {
	if l.isEmpty() {
		var zero T
		return zero, false
	}
	return l.root.next.value, true
}

// Back returns the value of the last element or zero value if empty.
// The complexity is O(1).
func (l *List[T]) Back() (T, bool) {
	if l.isEmpty() {
		var zero T
		return zero, false
	}
	return l.root.prev.value, true
}

// PushBack inserts a new value v at the back of list l.
// The complexity is O(1).
func (l *List[T]) PushBack(v T) {
	l.insertValue(v, l.root.prev)
}

// PushFront inserts a new value v at the front of list l.
// The complexity is O(1).
func (l *List[T]) PushFront(v T) {
	l.insertValue(v, &l.root)
}

// PopFront removes the first element (front) from list l and returns it. Return nil if the list is empty.
// The complexity is O(1).
func (l *List[T]) PopFront() (T, bool) {

	if l.isEmpty() {
		var zero T
		return zero, false
	}
	return l.remove(l.root.next).value, true
}

// PopBack removes the last element (back) from list l and returns it. Return nil if the list is empty.
// The complexity is O(1).
func (l *List[T]) PopBack() (T, bool) {
	if l.isEmpty() {
		var zero T
		return zero, false
	}
	return l.remove(l.root.prev).value, true

}

// InsertBefore inserts a new value v before the cursor c, return the new node. cursor c stays at the same position after the insertion.
// If c is point to the sentinel node, InsertBefore inserts to the tail (same effect as PushBack).
// If c is not associated with l, InsertBefore returns nil.
// The complexity is O(1).
func (l *List[T]) InsertBefore(v T, c *Cursor[T]) {
	if c.list != l || c.current == &c.list.root {
		return
	}
	if c.IsValid() {
		c.list.insertValue(v, c.current.prev)
	}

}

// InsertAfter inserts a new value v after the cursor c, cursor c stays at the same position after the insertion.
// If c is point to the sentinel node, InsertAfter inserts to the head (same effect as PushFront).
// If c is not associated with l or invalid, InsertAfter returns nil.
// The complexity is O(1).
func (l *List[T]) InsertAfter(v T, c *Cursor[T]) bool {
	if c.list != l || c.current == &c.list.root {
		return false

	}
	if c.IsValid() {
		c.list.insertValue(v, c.current)
		return true
	}

	return false

}

// RemoveAt removes the node at the cursor c, return the removed node. Cursor c move to the next node after the removal.
// If c is point to the sentinel node, RemoveAt returns nil.
func (l *List[T]) RemoveAt(c *Cursor[T]) T {

	var zero T
	if c.list != l || c.current == &c.list.root {
		return zero
	}

	if c.IsValid() {
		n := c.current
		c.current = c.current.next
		c.list.remove(n)
		return n.value
	}

	return zero
}

// RemoveAfter removes the node after the cursor c, return the removed node. Cursor c stays at the same position after the removal.
// If c is point to the sentinel node, RemoveAfter removes the first element of the list (same effect as RemoveFront).
// If c is not associated with l, RemoveAfter returns nil.
// The complexity is O(1).
func (l *List[T]) RemoveAfter(c *Cursor[T]) T {
	var zero T
	if c.list != l || l.root.prev == c.current {
		return zero
	}

	if c.IsValid() {
		n := c.current.next
		c.list.remove(n)
		return n.value
	}

	return zero

}

// RemoveBefore removes the node before the cursor c, return the removed node. Cursor c stays at the same position after the removal.
// If c is point to the sentinel node, RemoveBefore removes the last element of the list (same effect as RemoveBack).
// If c is not associated with l, RemoveBefore returns nil.
// The complexity is O(1).
func (l *List[T]) RemoveBefore(c *Cursor[T]) T {
	var zero T
	if c.list != l || l.root.next == c.current {
		return zero
	}

	if c.IsValid() {
		n := c.current.prev
		c.list.remove(n)
		return n.value
	}

	return zero
}

// MoveToFront moves the node at the cursor c to the front of the list.
// It does nothing if c is point to the sentinel node, invalid or not associated with l.
// The complexity is O(1).
func (l *List[T]) MoveToFront(c *Cursor[T]) {
	if c.list != l || l.root.next == c.current {
		return
	}

	if c.IsValid() {
		l.move(c.current, &l.root)
	}

}

// MoveToBack moves the node at the cursor c to the back of the list.
// It does nothing if c is point to the sentinel node, invalid or not associated with l.
// The complexity is O(1).
func (l *List[T]) MoveToBack(c *Cursor[T]) {
	if c.list != l || l.root.prev == c.current {
		return
	}

	if c.IsValid() {
		l.move(c.current, l.root.prev)
	}

}

// MoveBefore moves the node at the cursor c to the position before the cursor mark.
// It does nothing if c is point to the sentinel node, or c or mark are not associated with l.
// The complexity is O(1).
func (l *List[T]) MoveBefore(c, mark *Cursor[T]) {
	if c.list != l || c.current == mark.current || mark.list != l {
		return
	}

	if c.IsValid() && mark.IsValid() {
		l.move(c.current, mark.current.prev)
	}

}

// MoveAfter moves the node at the cursor c to the position after the cursor mark.
// It does nothing if c is point to the sentinel node, or c or mark are not associated with l.
// The complexity is O(1).
func (l *List[T]) MoveAfter(c, mark *Cursor[T]) {
	if c.list != l || c.current == mark.current || mark.list != l {
		return
	}

	if c.IsValid() && mark.IsValid() {
		l.move(c.current, mark.current)
	}
}

// Cursor returns a cursor pointing to the sentinel node of the list.
func (l *List[T]) Cursor() *Cursor[T] {
	return &Cursor[T]{list: l, current: &l.root}
}

// FrontCursor returns a cursor pointing to the first node of the list.
func (l *List[T]) FrontCursor() *Cursor[T] {
	return &Cursor[T]{list: l, current: l.root.next}
}

// BackCursor returns a cursor pointing to the last node of the list.
func (l *List[T]) BackCursor() *Cursor[T] {
	return &Cursor[T]{list: l, current: l.root.prev}
}

// PushBackList inserts a copy of an `other` list at the back of `l`.
func (l *List[T]) PushBackList(other *List[T]) {

	other.Cursor().WalkAscending(func(v T) bool {
		l.insertValue(v, l.root.prev)
		return true
	})
}

// PushFrontList inserts a copy of an `other` list at the front of `l`.
func (l *List[T]) PushFrontList(other *List[T]) {
	other.Cursor().WalkDescending(func(v T) bool {
		l.insertValue(v, &l.root)
		return true
	})
}
