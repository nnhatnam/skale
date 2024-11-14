package linkedlist

// Cursor is a read-only object that points to a node in a list. It contains a reference to the list and the node it's currently pointing to.
type Cursor[T any] struct {
	list    *List[T]
	current *node[T]
}

// Equal returns true if the two cursors point to the same node in the same list.
// if either cursor is not valid, it returns false.
func (c *Cursor[T]) Equal(c2 *Cursor[T]) bool {
	if !c.IsValid() || !c2.IsValid() {
		return false
	}

	return c.list == c2.list && c.current == c2.current
}

// Value returns the value of the node that the cursor points to.
// If the cursor is not valid, it will panic.
func (c *Cursor[T]) Value() T {
	if !c.IsValid() {
		panic("cursor is not valid when calling Value()")
	}
	return c.current.value
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

// node returns the node that the cursor points to.
func (c *Cursor[T]) node() *node[T] {
	if c.current != &c.list.root && c.IsValid() {
		return c.current
	}
	return nil
}

// nodeNext returns the node after the node the cursor is currently pointing to.
// Return nil if the cursor is pointing to the last node in the list.
func (c *Cursor[T]) nodeNext() *node[T] {
	if c.current.next != &c.list.root && c.IsValid() {
		return c.current.next
	}
	return nil
}

// nodePrev returns the node before the node the cursor is currently pointing to.
// Return nil if the cursor is pointing to the first node in the list.
func (c *Cursor[T]) nodePrev() *node[T] {
	if c.current.prev != &c.list.root && c.IsValid() {
		return c.current.prev
	}
	return nil
}

// MoveNext moves the cursor to the next node in the list and return the node.
// Move to sentinel node and return nil if the cursor is pointing to the last node in the list.
func (c *Cursor[T]) MoveNext() *node[T] {

	c.current = c.current.next
	if c.current == &c.list.root {
		return nil
	}
	return c.current
}

// MovePrev moves the cursor to the previous node in the list and return the node.
// Move to sentinel node and return nil if the cursor is pointing to the first node in the list.
func (c *Cursor[T]) MovePrev() *node[T] {

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
func (c *Cursor[T]) WalkAscending(f func(n *node[T]) bool) {
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
func (c *Cursor[T]) WalkDescending(f func(n *node[T]) bool) {
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
