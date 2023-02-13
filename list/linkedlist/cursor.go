package linkedlist

type Cursor[T any] struct {
	list    *List[T]
	current *Node[T]
}

func (l *List[T]) Cursor() *Cursor[T] {
	return &Cursor[T]{list: l, current: &l.root}
}

func (c *Cursor[T]) Current() *Node[T] {
	if c.list != nil && c.current != &c.list.root {
		return c.current
	}
	return nil
}

func (c *Cursor[T]) Next() *Node[T] {
	if c.current.next == &c.list.root {
		return nil
	}
	c.current = c.current.next
	return c.current
}

func (c *Cursor[T]) Prev() *Node[T] {
	if c.current.prev == &c.list.root {
		return nil
	}
	c.current = c.current.prev
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

func (c *Cursor[T]) MoveNext() bool {
	if c.current.next == &c.list.root {
		return false
	}
	c.current = c.current.next
	return true
}

func (c *Cursor[T]) MovePrev() bool {
	if c.current.prev == &c.list.root {
		return false
	}
	c.current = c.current.prev
	return true
}
