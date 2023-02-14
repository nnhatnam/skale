package unionfind

type Element[T any] struct {
	parent *Element[T]
	rank   int

	Value T
}

func MakeSet[T any](v T) *Element[T] {
	e := &Element[T]{Value: v}
	e.parent = e
	return e
}

func FindSet[T any](e *Element[T]) *Element[T] {
	if e.parent != e {
		e.parent = FindSet(e.parent)
	}
	return e.parent
}

func UnionSets[T any](a, b *Element[T]) {
	a = FindSet(a)
	b = FindSet(b)
	if a != b {
		if a.rank < b.rank {
			a, b = b, a
		}
		b.parent = a
		if a.rank == b.rank {
			a.rank++
		}
	}
}
