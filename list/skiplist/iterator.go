package skiplist

type ItemIterator[T any] func(item T) bool

type NodeIterator[T any] func(node *Node[T]) bool
