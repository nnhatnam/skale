package list

type Iterator struct {
	list    *List
	current *Node
}

func NewIterator(l *List) *Iterator {
	return &Iterator{list: l, current: l.root.Next}
}

// Node gives access to the current node. Be careful not to modify the node
// In skale, Node.Next and Node.Prev are public, so you're allowed to modify it but be aware that it may break the underline list
// Also, skale implements using a sentinel node, so if you extract the tail node, Node.Next will be the sentinel node
// If you extract the head node, Node.Prev will be the sentinel node
func (i *Iterator) Node() *Node {
	if i.list != nil && i.current != &i.list.root {
		return i.current
	}
	return nil
}

func (i *Iterator) Value() any {
	if i.current == &i.list.root {
		return nil
	}
	return i.current.Value
}

func (i *Iterator) HasNext() bool {
	return i.current.Next != &i.list.root
}

func (i *Iterator) Next() any {
	//if current is root (in case of empty list), return false
	//if current is the last node, return false
	if i.current == &i.list.root || i.current.Next == &i.list.root {
		return nil
	}
	i.current = i.current.Next
	return i.current.Value
}

func (i *Iterator) HasPrev() bool {
	return i.current.Prev != &i.list.root
}

func (i *Iterator) Prev() any {
	//if current is root (in case of empty list), return false
	//if current is the first node, return false
	if i.current == &i.list.root || i.current.Prev == &i.list.root {
		return nil
	}
	i.current = i.current.Prev
	return i.current
}

// First moves the iterator to the first node in the list
func (i *Iterator) First() {
	i.current = i.list.root.Next
}

// Last moves the iterator to the last node in the list
func (i *Iterator) Last() {
	i.current = i.list.root.Prev
}

// IsLast check if the current node is the last node in the list
func (i *Iterator) IsLast() bool {
	return i.current.Next == &i.list.root
}
