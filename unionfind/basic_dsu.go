package unionfind

type BasicDSU struct {
	parent []int
	rank   []int
}

func NewBasicDSU(n int) *BasicDSU {
	dsu := &BasicDSU{
		parent: make([]int, n),
		rank:   make([]int, n),
	}
	return dsu
}

func (dsu *BasicDSU) MakeSet(v int) {
	dsu.parent[v] = v
	dsu.rank[v] = 0
}

func (dsu *BasicDSU) FindSet(v int) int {
	if dsu.parent[v] == v {
		return v
	}
	dsu.parent[v] = dsu.FindSet(dsu.parent[v])
	return dsu.parent[v]
}

func (dsu *BasicDSU) UnionSets(a, b int) {
	a = dsu.FindSet(a)
	b = dsu.FindSet(b)
	if a != b {
		if dsu.rank[a] < dsu.rank[b] {
			a, b = b, a
		}
		dsu.parent[b] = a
		if dsu.rank[a] == dsu.rank[b] {
			dsu.rank[a]++
		}
	}
}
