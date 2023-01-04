package queue

import "github.com/nnhatnam/skale/list"

type Queue[V any] struct {
	elements list.LinearList[V]
}

func NewQueueWithContainer[V any](l list.LinearList[V]) *Queue[V] {
	return &Queue[V]{elements: l}
}

func (q *Queue[V]) Len() int {
	return q.elements.Len()
}

func (q *Queue[V]) PushBack(value V) {
	q.elements.PushBack(value)
}

func (q *Queue[V]) PopFront() V {
	return q.elements.PopFront()
}

func (q *Queue[V]) Front() V {
	return q.elements.Front()
}

func (q *Queue[V]) Back() V {
	return q.elements.Back()
}
