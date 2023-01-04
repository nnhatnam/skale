package list

type Cursor struct {
	list    *List
	current *Node
}

type MulFunc func(n *Node)

func NewCursor(l *List) *Cursor {
	return &Cursor{list: l, current: l.root.next}
}

// Node gives access to the current node.
func (i *Cursor) Node() *Node {
	if i.list != nil && i.current != &i.list.root {
		return i.current
	}
	return nil
}

func (i *Cursor) Value() any {
	if i.current == &i.list.root {
		return nil
	}
	return i.current.Value
}

func (i *Cursor) HasNext() bool {
	return i.current.next != &i.list.root
}

func (i *Cursor) Next() any {
	//if current is root (in case of empty list), return false
	//if current is the last node, return false
	if i.current == &i.list.root || i.current.next == &i.list.root {
		return nil
	}
	i.current = i.current.next
	return i.current.Value
}

func (i *Cursor) HasPrev() bool {
	return i.current.prev != &i.list.root
}

func (i *Cursor) Prev() any {
	//if current is root (in case of empty list), return false
	//if current is the first node, return false
	if i.current == &i.list.root || i.current.prev == &i.list.root {
		return nil
	}
	i.current = i.current.prev
	return i.current
}

// First moves the Cursor to the first node in the list
func (i *Cursor) First() {
	i.current = i.list.root.next
}

// Last moves the Cursor to the last node in the list
func (i *Cursor) Last() {
	i.current = i.list.root.prev
}

// IsLast check if the current node is the last node in the list
func (i *Cursor) IsLast() bool {
	return i.current.next == &i.list.root
}
