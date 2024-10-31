// Package linkedlist implements a doubly linked list data structure.
// It is intended to be used internally by other packages.
//
// Structure is not thread safe.
//
// The implementation is based on the standard library's list package, but with several changes in the API:
//
// 1. The list is generic and can be used with any type.
//
// 2. *Element is renamed to *Node to avoid confusion with the standard library's list package.
//
// 3. Instead of using *Node as a way to iterate through the list, this package introduces a *Cursor type to do that, which is similar to c++ iter concept.
//
// Why Cursor?
//
// In the standard library's list package, internally, each `Element` has a pointer to the list it belongs to,
// which helps prevent mistakes like adding an `Element` to multiple lists or a list can delete an Element from another list.
// It offers a good level of foolproof. However, this comes at the cost of an extra pointer that increases the memory footprint by 8 bytes for 64-bit arch or 4 bytes for 32-bit arch.
// If you are implementing a data structure that uses a linked list internally, and you care about memory footprint, removing that the extra pointer and taking the responsibility of not making mistakes may make more sense.
// That's the reason for this package. Each `Node` in this package has no pointer to the list it belongs to. Without the extra pointer expose the risk of mistakes mention above.
// To avoid mistakes, this package introduces `Cursor`, which is similar, but not identical, to c++ iter concept.
//
// A Cursor is a list reader contains a pointer to the list it belongs to and a pointer to the current Node it points to. It can move freely in the list in both directions.
// To use a Cursor, you need to aware how the list is implemented internally.
// The list has a sentinel node that points to the first and last node in the list. It helps simplify certain list operations.
// So, in internal view, the list can be viewed as a ring with the sentinel node as the head and tail, and Cursor views the list in that manner.
// It is allowed a Cursor to point to the sentinel node, but it will return nil when you extract the Node if it is pointing to the sentinel.
// Technically, the Cursor will be able to move infinitely in both directions.
//
// To iterate over a list (where l is a *List):
//
//	cursor := l.Cursor() // create a cursor point to first node
//	for n := cursor.MoveNext(); n != nil; n = cursor.MoveNext() {
//		// do something with n
//	}
//
// Another way to iterate over a list:
//
//	cursor := l.Cursor() // create a cursor point to the sentinel node
//	for cursor.MoveNext() != nil {
//		n := cursor.Node()
//		// do something with n
//	}
package linkedlist

// Node is a node in a doubly linked list
type Node[T any] struct {
	next, prev *Node[T]

	Value T
}

// NewNode creates a new node with the given value
func NewNode[T any](value T) *Node[T] {
	return &Node[T]{Value: value}
}

// List represents a doubly linked list.
type List[T any] struct {
	root Node[T]
	len  int
}

// New returns an initialized list.
func New[T any]() *List[T] {
	l := &List[T]{}
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

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

// first returns the first node in the list
func (l *List[T]) front() *Node[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// last returns the last node in the list
func (l *List[T]) back() *Node[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// insert inserts a node after mark. The mask must not be nil.
func (l *List[T]) insert(n, mark *Node[T]) *Node[T] {

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

// remove removes n from the list. The node must not be nil.
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

// From returns an initialized list and add the given values, if any, to the list.
func From[T any](values ...T) *List[T] {
	l := New[T]()
	for _, v := range values {
		l.insertValue(v, l.root.prev)
	}
	return l
}

// Len returns the number of elements of list l. The complexity is O(1).
func (l *List[T]) Len() int {
	return l.len
}

// Front returns the first element of list l. Return nil if the list is empty.
// The complexity is O(1).
func (l *List[T]) Front() *Node[T] {
	return l.front()
}

// Back returns the last element of list l. Return nil if the list is empty.
// The complexity is O(1).
func (l *List[T]) Back() *Node[T] {
	return l.back()
}

// PushBack inserts a new value v at the back of list l.
// The complexity is O(1).
func (l *List[T]) PushBack(v T) {
	l.lazyInit()
	l.insertValue(v, l.root.prev)
}

// PushBackBulk inserts given values at the back of list l. PushBackBulk is slightly cheaper than calling PushBack in a loop.
// The complexity is O(len(values)).
func (l *List[T]) PushBackBulk(values ...T) {
	l.lazyInit()
	for _, v := range values {
		l.insertValue(v, l.root.prev)
	}
}

// PushFront inserts a new value v at the front of list l.
// The complexity is O(1).
func (l *List[T]) PushFront(v T) {
	l.lazyInit()
	l.insertValue(v, &l.root)
}

// PushFrontBulk inserts given values at the front of list l. PushFrontBulk is slightly cheaper than calling PushFront in a loop.
// The complexity is O(len(values)).
func (l *List[T]) PushFrontBulk(values ...T) {
	l.lazyInit()
	for _, v := range values {
		l.insertValue(v, &l.root)
	}
}

// PopFront removes the first element (front) from list l and returns it. Return nil if the list is empty.
// The complexity is O(1).
func (l *List[T]) PopFront() *Node[T] {
	l.lazyInit()

	n := l.front()
	if n != nil {
		return l.remove(n)
	}
	return nil

}

// PopBack removes the last element (back) from list l and returns it. Return nil if the list is empty.
// The complexity is O(1).
func (l *List[T]) PopBack() *Node[T] {
	l.lazyInit()

	if l.len == 0 {
		return nil
	}
	n := l.back()
	if n != nil {
		return l.remove(n)
	}

	return nil

}

// InsertBefore inserts a new value v before the cursor c, return the new node. cursor c stays at the same position after the insertion.
// If c is point to the sentinel node, InsertBefore inserts to the tail (same effect as PushBack).
// If c is not associated with l, InsertBefore returns nil.
// The complexity is O(1).
func (l *List[T]) InsertBefore(v T, c *Cursor[T]) *Node[T] {
	if c.list != l || c.current == &c.list.root {
		return nil
	}
	if c.IsValid() {
		return c.list.insertValue(v, c.current.prev)
	}

	return nil
}

// InsertAfter inserts a new value v after the cursor c, return the new node. cursor c stays at the same position after the insertion.
// If c is point to the sentinel node, InsertAfter inserts to the head (same effect as PushFront).
// If c is not associated with l or invalid, InsertAfter returns nil.
// The complexity is O(1).
func (l *List[T]) InsertAfter(v T, c *Cursor[T]) *Node[T] {
	if c.list != l || c.current == &c.list.root {
		return nil
	}
	if c.IsValid() {
		return c.list.insertValue(v, c.current)
	}

	return nil
}

// RemoveAt removes the node at the cursor c, return the removed node. Cursor c move to the next node after the removal.
// If c is point to the sentinel node, RemoveAt returns nil.
func (l *List[T]) RemoveAt(c *Cursor[T]) *Node[T] {

	if c.list != l || c.current == &c.list.root {
		return nil
	}

	if c.IsValid() {
		n := c.current
		c.current = c.current.next
		c.list.remove(n)
		return n
	}

	return nil
}

// RemoveAfter removes the node after the cursor c, return the removed node. Cursor c stays at the same position after the removal.
// If c is point to the sentinel node, RemoveAfter removes the first element of the list (same effect as RemoveFront).
// If c is not associated with l, RemoveAfter returns nil.
// The complexity is O(1).
func (l *List[T]) RemoveAfter(c *Cursor[T]) *Node[T] {
	if c.list != l || l.root.prev == c.current {
		return nil
	}

	if c.IsValid() {
		return c.list.remove(c.current.next)
	}

	return nil

}

// RemoveBefore removes the node before the cursor c, return the removed node. Cursor c stays at the same position after the removal.
// If c is point to the sentinel node, RemoveBefore removes the last element of the list (same effect as RemoveBack).
// If c is not associated with l, RemoveBefore returns nil.
// The complexity is O(1).
func (l *List[T]) RemoveBefore(c *Cursor[T]) *Node[T] {

	if c.list != l || l.root.next == c.current {
		return nil
	}

	if c.IsValid() {
		return c.list.remove(c.current.prev)
	}
	return nil
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
	l.lazyInit()
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
	l.lazyInit()
	back := other.BackCursor().Node()

	other.Cursor().WalkAscending(func(n *Node[T]) bool {
		l.insertValue(n.Value, l.root.prev)
		if n == back {
			return false
		}
		return true
	})
}

// PushFrontList inserts a copy of an `other` list at the front of `l`.
func (l *List[T]) PushFrontList(other *List[T]) {
	l.lazyInit()
	front := other.FrontCursor().Node()

	other.Cursor().WalkDescending(func(n *Node[T]) bool {

		l.insertValue(n.Value, &l.root)
		if n == front {
			return false
		}
		return true
	})
}
