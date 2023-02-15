package linkedlist

type Cursor[T any] struct {
	list    *List[T]
	current *Node[T]
}

func (l *List[T]) Cursor() *Cursor[T] {
	return &Cursor[T]{list: l, current: &l.root}
}

func (l *List[T]) FrontCursor() *Cursor[T] {
	return &Cursor[T]{list: l, current: l.root.next}
}

func (l *List[T]) BackCursor() *Cursor[T] {
	return &Cursor[T]{list: l, current: l.root.prev}
}

func (c *Cursor[T]) NextCursor() *Cursor[T] {
	
	return &Cursor[T]{list: c.list, current: c.current.next}
}

func (c *Cursor[T]) PrevCursor() *Cursor[T] {

	return &Cursor[T]{list: c.list, current: c.current.prev}
}

func (c *Cursor[T]) Node() *Node[T] {
	if c.current != &c.list.root {
		return c.current
	}
	return nil
}

func (c *Cursor[T]) NextNode() *Node[T] {
	if c.current.next != &c.list.root {
		return c.current.next
	}
	return nil
}

func (c *Cursor[T]) PrevNode() *Node[T] {
	if c.current.prev != &c.list.root {
		return c.current.prev
	}
	return nil
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

func (c *Cursor[T]) Ascending(iter IterFunc[T]) {
	if c.list.len > 0 {

		if c.current != &c.list.root {
			if !iter(c.current) {
				return
			}
		}

		for c.MoveNext() != nil {
			if !iter(c.current) {
				return
			}
		}

	}

}

func (c *Cursor[T]) Descending(iter IterFunc[T]) {
	if c.list.len > 0 {

		if c.current != &c.list.root {
			if !iter(c.current) {
				return
			}
		}

		for c.MovePrev() != nil {
			if !iter(c.current) {
				return
			}
		}

	}

}

func (c *Cursor[T]) Clone() *Cursor[T] {
	return &Cursor[T]{list: c.list, current: c.current}
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
