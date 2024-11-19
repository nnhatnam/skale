package linkedlist

// Cursor is a read-only object that points to a node in a list. It contains a reference to the list and the node it's currently pointing to.
type Cursor[T any] struct {
	list    *List[T]
	current *node[T]
}

// Value returns the value of the current node if valid.
func (c *Cursor[T]) Value() (T, bool) {
	if !c.IsValid() {
		var zero T
		return zero, false
	}
	return c.current.value, true
}

// Clone creates a new cursor that points to the same node as the current cursor.
// Return nil if the current cursor is not valid.
func (c *Cursor[T]) Clone() *Cursor[T] {
	if !c.IsValid() {
		return nil
	}
	return &Cursor[T]{list: c.list, current: c.current}
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

// MoveNext advances the cursor to the next node and return true for each successful move.
// Calling MoveNext at the end of the list will return false and move the cursor to the sentinel node.
// Calling MoveNext on an invalid cursor will return false.
func (c *Cursor[T]) MoveNext() bool {

	if !c.IsValid() {
		return false
	}

	c.current = c.current.next
	return c.current != &c.list.root

}

// MovePrev moves the cursor to the previous node and return true for each successful move.
// Calling MovePrev at the beginning of the list will return false and move the cursor to the sentinel node.
// Calling MovePrev on an invalid cursor will return false.
func (c *Cursor[T]) MovePrev() bool {

	if !c.IsValid() {
		return false
	}

	c.current = c.current.prev
	return c.current != &c.list.root

}

// MoveToFront moves the valid cursor to the first node in the list. If the list is empty or the cursor is invalid, return false.
func (c *Cursor[T]) MoveToFront() bool {
	if c.list.isEmpty() {
		return false
	}

	c.current = c.list.root.next
	return true
}

// MoveToBack moves the valid cursor to the last node in the list. If the list is empty or the cursor is invalid, return false.
func (c *Cursor[T]) MoveToBack() bool {
	if c.list.isEmpty() {
		return false
	}
	c.current = c.list.root.prev
	return true
}

// WalkAscending moves the cursor to the next node in the list and calls the function f with the node.
// It continues until f returns false or the cursor reaches the sentinel node.
func (c *Cursor[T]) WalkAscending(f func(v T) bool) {

	if c.list.len == 0 || !c.IsValid() {
		return // Early exit for empty list or invalid cursor
	}

	if c.current == &c.list.root {
		c.current = c.current.next
	}

	for {

		// Stop if we reach the sentinel node
		if c.current == &c.list.root {
			return
		}

		// Call the function with the current node
		if !f(c.current.value) {
			return // Stop walking if f returns false
		}

		// Move to the next node
		c.current = c.current.next
	}
}

// WalkDescending moves the cursor to the previous node in the list and call the function f with the node.
// Keep walking until f returns false or the cursor reach the head of the list.
func (c *Cursor[T]) WalkDescending(f func(v T) bool) {

	if c.list.len == 0 || !c.IsValid() {
		return // Early exit for empty list or invalid cursor
	}

	if c.current == &c.list.root {
		c.current = c.current.prev
	}

	for {

		// Stop if we reach the head of the list
		if c.current == &c.list.root {
			return
		}

		// Call the function with the current node
		if !f(c.current.value) {
			return // Stop walking if f returns false
		}

		// Move to the previous node
		c.current = c.current.prev

	}
}

// Close releases the cursor's reference to the list. The cursor can no longer be used.
// defer cursor.Close() is recommended.
func (c *Cursor[T]) Close() {
	c.list = nil
	c.current = nil
}

// IsValid detects if the cursor is valid.
// A cursor is not valid if it is closed, or it is pointing to a node that no longer exists.
func (c *Cursor[T]) IsValid() bool {
	// Invalid happens in two cases:
	// 1. The cursor has no reference to the list or the node.
	// 2. The cursor is pointing to a node that has been removed from the list.
	// Note that a cursor pointing to the sentinel node (root cursor) is still valid.

	if c.list == nil || c.current == nil {
		// Case 1
		return false
	} else if c.current.next == nil {
		// Case 2
		c.Close()
		return false
	}

	return true
}
