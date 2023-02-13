package linkedlist

type Cursor[T any] struct {
	list    *List[T]
	current *Node[T]
}

func (l *List[T]) Cursor() *Cursor[T] {
	return &Cursor[T]{list: l, current: &l.root}
}

func (l *List[T]) CursorFront() *Cursor[T] {
	return &Cursor[T]{list: l, current: l.root.next}
}

func (l *List[T]) CursorBack() *Cursor[T] {
	return &Cursor[T]{list: l, current: l.root.prev}
}

func (c *Cursor[T]) Current() *Node[T] {
	if c.current != &c.list.root {
		return c.current
	}
	return nil
}

func (c *Cursor[T]) Value() (T, bool) {
	var zero T
	if c.current != &c.list.root {
		return c.current.Value, true
	}
	return zero, false
}

func (c *Cursor[T]) Next() *Node[T] {

	if c.current.next == &c.list.root {
		return nil
	}
	return c.current.next
}

func (c *Cursor[T]) Prev() *Node[T] {

	if c.current.prev == &c.list.root {
		return nil
	}
	return c.current.prev
}

func (c *Cursor[T]) MoveNext() *Node[T] {

	c.current = c.current.next
	if c.current == &c.list.root {
		return nil
	}
	return c.current
}

func (c *Cursor[T]) MovePrev() *Node[T] {

	c.current = c.current.prev
	if c.current == &c.list.root {
		return nil
	}
	return c.current
}

func (c *Cursor[T]) MoveToFront() bool {
	if c.list.len == 0 {
		return false
	}
	c.current = c.list.root.next
	return true
}

func (c *Cursor[T]) MoveToBack() bool {
	if c.list.len == 0 {
		return false
	}
	c.current = c.list.root.prev
	return true
}

//func (c *Cursor[T]) MoveNext() bool {
//	if c.current.next == &c.list.root {
//		return false
//	}
//	c.current = c.current.next
//	return true
//}
//
//func (c *Cursor[T]) MovePrev() bool {
//	if c.current.prev == &c.list.root {
//		return false
//	}
//	c.current = c.current.prev
//	return true
//}
