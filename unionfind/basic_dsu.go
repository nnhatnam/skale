package unionfind

type Parent[T comparable] interface {
	Get(v T) T
	Set(i, v T)
}

type Rank[T comparable] interface {
	Get(r T) int
	Set(i T, r int)
}

type BasicDSU[T comparable] struct {
	parent Parent[T]
	rank   Rank[T]
}

func NewBasicDSU[T comparable](parent Parent[T], rank Rank[T]) *BasicDSU[T] {
	return &BasicDSU[T]{parent: parent, rank: rank}
}

func (uf *BasicDSU[T]) MakeSet(v T) {
	uf.parent.Set(v, v)
	uf.rank.Set(v, 0)
}

func (uf *BasicDSU[T]) FindSet(v T) T {
	p := uf.parent.Get(v)
	if p != v {
		p = uf.FindSet(p)
		uf.parent.Set(v, p)
	}
	return p
}

func (uf *BasicDSU[T]) UnionSets(v1, v2 T) {
	uf.link(uf.FindSet(v1), uf.FindSet(v2))
}

func (uf *BasicDSU[T]) link(v1, v2 T) {
	if uf.rank.Get(v1) > uf.rank.Get(v2) {
		uf.parent.Set(v2, v1)
	} else {
		uf.parent.Set(v1, v2)
		if uf.rank.Get(v1) == uf.rank.Get(v2) {
			uf.rank.Set(v2, uf.rank.Get(v1)+1)
		}
	}
}
