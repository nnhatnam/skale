package unionfind

type DSU[T comparable] struct {
	sets map[T]*element[T] // TODO: replace it to a BIMap
}

func NewDisjointSet[T comparable]() *DSU[T] {
	return &DSU[T]{sets: make(map[T]*element[T])}
}

func (ds *DSU[T]) MakeSet(v T) {
	ds.sets[v] = newElement(v)
}

func (ds *DSU[T]) FindSet(v T) T {
	return findRep(ds.sets[v]).value
}

func (ds *DSU[T]) UnionSets(x, y T) {
	link(findRep(ds.sets[x]), findRep(ds.sets[y]))
}

type element[T any] struct {
	parent *element[T]
	rank   int

	value T
}

// newElement creates a new element with the given value.
// An element is a set by itself.
func newElement[T any](v T) *element[T] {
	e := &element[T]{value: v}
	e.parent = e
	return e
}

// findRep returns the representative element of the set containing the given element.
func findRep[T any](e *element[T]) *element[T] {
	if e.parent != e {
		e.parent = findRep(e.parent)
	}
	return e.parent
}

func link[T any](x, y *element[T]) {
	if x.rank > y.rank {
		y.parent = x
	} else {
		x.parent = y
		if x.rank == y.rank {
			y.rank++
		}
	}
}
