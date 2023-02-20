package linkedlist

// Cursor is a read-only object that points to a node in a list. It contains a reference to the list and the node it's currently pointing to.
// To modify the list, caller must use the methods provided by the list struct. List is the only place we can modify the list.
// The List struct provides several methods that modify the list based on information obtained from the Cursor, such as MoveToFront, MoveToBack, RemoveAt, etc.
//
// While a cursor is similar to an iterator, but it has a few important differences.
// A cursor can freely move back-and-forth in the list and views the list as the way it is implemented.
// Which mean, in a cursor perspective, the list is a circular doubly linked list.
// The result is that if a cursor is pointing at the end of the list, call MoveNext() will point it to the sentinel node.
// If MoveNext() is called again, it will point to the first node of the list. The same applies to MovePrev(), but in the opposite direction.
//
// It's important to use caution when a cursor is pointing to the sentinel node, as this can lead to unexpected behavior.
// It's similar to being aware of an out-of-range index when working with a slice.
type Cursor[T any] struct {
	list    *List[T]
	current *Node[T]
}

// Clone creates a new cursor that points to the same node as the current cursor.
func (c *Cursor[T]) Clone() *Cursor[T] {
	return &Cursor[T]{list: c.list, current: c.current}
}

// CloneNext creates a new cursor that points to the next node in the list.
// If the current node is the last node in the list, the new cursor will point to the sentinel node.
// If the current node is the sentinel node, the new cursor will point to the first node in the list.
func (c *Cursor[T]) CloneNext() *Cursor[T] {

	return &Cursor[T]{list: c.list, current: c.current.next}
}

// ClonePrev creates a new cursor that points to the previous node in the list.
// If the current node is the first node in the list, the new cursor will point to the sentinel node.
// If the current node is the sentinel node, the new cursor will point to the last node in the list.
func (c *Cursor[T]) ClonePrev() *Cursor[T] {

	return &Cursor[T]{list: c.list, current: c.current.prev}
}

// Node returns the node that the cursor points to.
func (c *Cursor[T]) Node() *Node[T] {
	if c.current != &c.list.root {
		return c.current
	}
	return nil
}

// NodeNext returns the node after the node the cursor is currently pointing to.
// Return nil if the cursor is pointing to the last node in the list.
func (c *Cursor[T]) NodeNext() *Node[T] {
	if c.current.next != &c.list.root {
		return c.current.next
	}
	return nil
}

// NodePrev returns the node before the node the cursor is currently pointing to.
// Return nil if the cursor is pointing to the first node in the list.
func (c *Cursor[T]) NodePrev() *Node[T] {
	if c.current.prev != &c.list.root {
		return c.current.prev
	}
	return nil
}

// MoveNext moves the cursor to the next node in the list and return the node.
// Move to sentinel node and return nil if the cursor is pointing to the last node in the list.
func (c *Cursor[T]) MoveNext() *Node[T] {

	c.current = c.current.next
	if c.current == &c.list.root {
		return nil
	}
	return c.current
}

// MovePrev moves the cursor to the previous node in the list and return the node.
// Move to sentinel node and return nil if the cursor is pointing to the first node in the list.
func (c *Cursor[T]) MovePrev() *Node[T] {

	c.current = c.current.prev
	if c.current == &c.list.root {
		return nil
	}
	return c.current
}

// MoveToFront moves the cursor to the first node in the list. If the list is empty, return false.
func (c *Cursor[T]) MoveToFront() bool {
	if c.list.len == 0 {
		return false
	}
	c.current = c.list.root.next
	return true
}

// MoveToBack moves the cursor to the last node in the list. If the list is empty, return false.
func (c *Cursor[T]) MoveToBack() bool {
	if c.list.len == 0 {
		return false
	}
	c.current = c.list.root.prev
	return true
}

// WalkAscending moves the cursor to the next node in the list and call the function f with the node.
// Keep walking until f returns false or the cursor reach the sentinel node.
func (c *Cursor[T]) WalkAscending(f func(n *Node[T]) bool) {
	if c.list.len > 0 {

		if c.current != &c.list.root {
			if !f(c.current) {
				return
			}
		}

		for c.MoveNext() != nil {
			if !f(c.current) {
				return
			}
		}
	}
}

// WalkDescending moves the cursor to the previous node in the list and call the function f with the node.
// Keep walking until f returns false or the cursor reach the sentinel node.
func (c *Cursor[T]) WalkDescending(f func(n *Node[T]) bool) {
	if c.list.len > 0 {

		if c.current != &c.list.root {
			if !f(c.current) {
				return
			}
		}

		for c.MovePrev() != nil {
			if !f(c.current) {
				return
			}
		}

	}

}

// Close closes the cursor and release the reference to the list. The cursor can no longer be used.
// defer cursor.Close() is recommended.
func (c *Cursor[T]) Close() {
	c.list = nil
	c.current = nil
}
