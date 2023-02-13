package linkedlist

type CursorMut[T any] struct {
	Cursor[T]
}

func (l *List[T]) CursorMut() *CursorMut[T] {
	return &CursorMut[T]{Cursor: Cursor[T]{list: l, current: &l.root}}
}

func (c *CursorMut[T]) InsertBefore(v T) *Node[T] {
	if c.current == &c.list.root {
		return nil
	}
	c.current = c.list.insertValue(v, c.current.prev)
	return c.current
}

func (c *CursorMut[T]) InsertAfter(v T) *Node[T] {
	if c.current == &c.list.root {
		return nil
	}
	c.current = c.list.insertValue(v, c.current)
	return c.current
}

func (c *CursorMut[T]) Remove() *Node[T] {
	if c.current == &c.list.root {
		return nil
	}
	n := c.current
	c.current = c.current.prev
	c.list.remove(n)
	return n
}
