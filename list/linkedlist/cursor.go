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
//
// Multiple cursors can be created from a single list. However, aware of situation where multiple cursors are pointing to the same node,
// and the list removes the node with one of the cursors. This will cause the other cursors to point to a node that no longer exists,
// but they are unaware of it and still think they are pointing to a valid node (still have a reference to the node and the list).
// It may cause unexpected behavior. It's the caller's responsibility to manage the cursors properly.
type Cursor[T any] struct {
	list    *List[T]
	current *Node[T]
}

// NewCursor creates a new cursor that points to the given node in the list.
// NewCursor returns nil if the given node is nil.
// NewCursor doesn't check if the given node is actually in the list, and it can make unexpected bug if the given node is not in the list.
// It is the caller's responsibility to ensure that the node is in the list. Use with caution. It's recommended to use List.Cursor() instead.
func NewCursor[T any](list *List[T], current *Node[T]) *Cursor[T] {
	if current == nil || list == nil {
		return nil
	}
	return &Cursor[T]{list: list, current: current}
}

// Equal returns true if the two cursors point to the same node in the same list.
// if either cursor is not valid, it returns false.
func (c *Cursor[T]) Equal(c2 *Cursor[T]) bool {
	if !c.IsValid() || !c2.IsValid() {
		return false
	}

	return c.list == c2.list && c.current == c2.current
}

// Clone creates a new cursor that points to the same node as the current cursor.
// Return nil if the current cursor is not valid.
func (c *Cursor[T]) Clone() *Cursor[T] {
	if c.IsValid() {
		return &Cursor[T]{list: c.list, current: c.current}
	}
	return nil
}

// CloneNext creates a new cursor that points to the next node in the list.
// If the current node is the last node in the list, the new cursor will point to the sentinel node.
// If the current node is the sentinel node, the new cursor will point to the first node in the list.
// Return nil if the current cursor is not valid.
func (c *Cursor[T]) CloneNext() *Cursor[T] {

	if c.IsValid() {
		return &Cursor[T]{list: c.list, current: c.current.next}
	}
	return nil // invalid cursor
}

// ClonePrev creates a new cursor that points to the previous node in the list.
// If the current node is the first node in the list, the new cursor will point to the sentinel node.
// If the current node is the sentinel node, the new cursor will point to the last node in the list.
// Return nil if the current cursor is not valid.
func (c *Cursor[T]) ClonePrev() *Cursor[T] {

	if c.IsValid() {
		return &Cursor[T]{list: c.list, current: c.current.prev}
	}
	return nil // invalid cursor
}

// Node returns the node that the cursor points to.
func (c *Cursor[T]) Node() *Node[T] {
	if c.current != &c.list.root && c.IsValid() {
		return c.current
	}
	return nil
}

// NodeNext returns the node after the node the cursor is currently pointing to.
// Return nil if the cursor is pointing to the last node in the list.
func (c *Cursor[T]) NodeNext() *Node[T] {
	if c.current.next != &c.list.root && c.IsValid() {
		return c.current.next
	}
	return nil
}

// NodePrev returns the node before the node the cursor is currently pointing to.
// Return nil if the cursor is pointing to the first node in the list.
func (c *Cursor[T]) NodePrev() *Node[T] {
	if c.current.prev != &c.list.root && c.IsValid() {
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

// MoveToFront moves the valid cursor to the first node in the list. If the list is empty or the cursor is invalid, return false.
func (c *Cursor[T]) MoveToFront() bool {
	if c.list.len == 0 {
		return false
	}

	if c.IsValid() {
		c.current = c.list.root.next
		return true
	}

	return false
}

// MoveToBack moves the valid cursor to the last node in the list. If the list is empty or the cursor is invalid, return false.
func (c *Cursor[T]) MoveToBack() bool {
	if c.list.len == 0 {
		return false
	}

	if c.IsValid() {
		c.current = c.list.root.prev
		return true
	}

	return false
}

// WalkAscending moves the cursor to the next node in the list and call the function f with the node.
// Keep walking until f returns false or the cursor reach the sentinel node.
func (c *Cursor[T]) WalkAscending(f func(n *Node[T]) bool) {
	if c.list.len > 0 && c.IsValid() {

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
	if c.list.len > 0 && c.IsValid() {

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

// IsValid detects if the cursor is valid.
// A cursor is not valid if it is closed, or it is pointing to a node that no longer exists.
// If the cursor is not valid, Close() will be called automatically.
func (c *Cursor[T]) IsValid() bool {
	if c.list == nil || c.current == nil || c.current.next == nil {
		c.Close()
		return false
	}
	return true
}
