package heap

type IHeap[T any] interface {
	// Len returns the number of elements in the heap.
	Len() int

	// Push pushes the value v onto the heap.
	Push(v T)

	// Pop removes the root element from the heap and returns it.
	Pop() (_ T, _ bool)

	// Peek returns the root element of the heap without removing it.
	Peek() (_ T, _ bool)
}
