package list

type LinearList[L any] interface {
	Len() int
	PushBack(value L)
	Front() L
	Back() L
}
