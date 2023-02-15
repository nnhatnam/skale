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

func UnionSets[T any](x, y *Element[T]) {
	link(FindSet(x), FindSet(y))
}

func link[T any](x, y *Element[T]) {
	if x.rank > y.rank {
		y.parent = x
	} else {
		x.parent = y
		if x.rank == y.rank {
			y.rank++
		}
	}
}
