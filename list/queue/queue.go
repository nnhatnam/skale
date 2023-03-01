package queue

type Queue[T any] interface {
	Peek() (_ T, _ bool)
	Len() int
	Enqueue(v T)
	Dequeue() (_ T, _ bool)
	Empty() bool
	Clear()
}
