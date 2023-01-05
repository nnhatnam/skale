package stack

import (
	"github.com/nnhatnam/skale/list"
)

type Stack[V any] struct {
	elements list.LinearList[V]
}

//func New[V any]() *Stack[V] {
//	return &Stack[V]{elements: new(Li)}
//}

func NewWithContainer[V any](l list.LinearList[V]) *Stack[V] {
	return &Stack[V]{elements: l}
}

func (s *Stack[V]) Top() V {
	return s.elements.Back()
}

func (s *Stack[V]) Len() int {
	return s.Len()
}

func (s *Stack[V]) PushBack(value V) {
	s.elements.PushBack(value)
}

func (s *Stack[V]) PopBack() V {
	return s.elements.PopBack()
}

func (s *Stack[V]) Front() V {
	return s.elements.Front()
}

func (s *Stack[V]) Back() V {
	return s.elements.Back()
}
