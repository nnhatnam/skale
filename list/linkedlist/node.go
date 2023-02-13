package linkedlist

// Node is a node in a double linked linkedlist
type Node[T any] struct {
	next, prev *Node[T]

	Value T
}

func NewNode[T any](value T) *Node[T] {
	return &Node[T]{Value: value}
}
