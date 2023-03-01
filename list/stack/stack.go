package stack

type Stack[T any] interface {
	Top() (_ T, _ bool)
	Len() int
	Push(v T)
	Pop() (_ T, _ bool)
	Empty() bool
	Clear() // clear all elements
}
