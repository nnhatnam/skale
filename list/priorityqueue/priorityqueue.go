package priorityqueue

type PriorityQueue[T any] interface {
	Push(x T)
	Pop() T
}
