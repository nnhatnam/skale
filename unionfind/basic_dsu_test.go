package unionfind

import (
	"golang.org/x/exp/slices"
	"testing"
)

type ints []int // implements Parent[int]

func (p ints) Get(v int) int {
	return p[v]
}

func (p ints) Set(i, v int) {
	p[i] = v
}

func TestBasicDSU(t *testing.T) {

	parent := make(ints, 10)
	rank := make(ints, 10)

	ds := NewBasicDSU[int](parent, rank)

	ds.MakeSet(1)
	ds.MakeSet(2)
	ds.MakeSet(3)
	ds.MakeSet(4)
	ds.MakeSet(5)
	ds.MakeSet(6)

	ds.UnionSets(1, 2)
	ds.UnionSets(1, 3)
	ds.UnionSets(1, 4)

	if slices.Index([]int{1, 2, 3, 4}, ds.FindSet(1)) == -1 {
		t.Errorf("FindSet(1) must be 1 or 2 or 3 or 4")
	}

	if slices.Index([]int{1, 2, 3, 4}, ds.FindSet(2)) == -1 {
		t.Errorf("FindSet(2) must be 1 or 2 or 3 or 4")
	}

	if slices.Index([]int{1, 2, 3, 4}, ds.FindSet(3)) == -1 {
		t.Errorf("FindSet(3) must be 1 or 2 or 3 or 4")
	}

	if slices.Index([]int{1, 2, 3, 4}, ds.FindSet(4)) == -1 {
		t.Errorf("FindSet(4) must be 1 or 2 or 3 or 4")
	}

	if ds.FindSet(5) != 5 {
		t.Errorf("FindSet(5) must be 5")
	}

	if ds.FindSet(6) != 6 {
		t.Errorf("FindSet(6) must be 6")
	}
}
