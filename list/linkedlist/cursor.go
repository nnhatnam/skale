package linkedlist

type Cursor[T any] struct {
	list    *List[T]
	current *Node[T]
}

func (c *Cursor[T]) CloneNext() *Cursor[T] {

	return &Cursor[T]{list: c.list, current: c.current.next}
}

func (c *Cursor[T]) ClonePrev() *Cursor[T] {

	return &Cursor[T]{list: c.list, current: c.current.prev}
}

func (c *Cursor[T]) Node() *Node[T] {
	if c.current != &c.list.root {
		return c.current
	}
	return nil
}

func (c *Cursor[T]) NodeNext() *Node[T] {
	if c.current.next != &c.list.root {
		return c.current.next
	}
	return nil
}

func (c *Cursor[T]) NodePrev() *Node[T] {
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

func (c *Cursor[T]) Clone() *Cursor[T] {
	return &Cursor[T]{list: c.list, current: c.current}
}

func (c *Cursor[T]) Close() {
	c.list = nil
	c.current = nil
}
