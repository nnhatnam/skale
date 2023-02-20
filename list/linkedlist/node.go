package linkedlist

// Node is a node in a doubly linked list
type Node[T any] struct {
	next, prev *Node[T]

	Value T
}

// NewNode creates a new node with the given value
func NewNode[T any](value T) *Node[T] {
	return &Node[T]{Value: value}
}
