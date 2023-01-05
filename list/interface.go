package list

type LinearList[V any] interface {
	Len() int
	PushBack(value V)
	PushFront(value V)
	PopBack() V
	PopFront() V
	Front() V
	Back() V
}
