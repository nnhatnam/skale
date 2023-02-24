package skiplist

import (
	"golang.org/x/exp/slices"
	"math/rand"
	"testing"
)

func TestNew(t *testing.T) {
	l := NewOrdered[int](63, 0.5)

	if l == nil {
		t.Error("New() failed")
	}

	if len(l.root.next) != 64 {
		t.Error("New() failed")
	}

	if l.root.next[0] == nil {
		t.Error("New() failed")
	}

	if l.root.next[0] != &l.root {
		t.Error("New() failed")
	}

	if l.root.prev == nil {
		t.Error("New() failed")
	}

	if l.root.prev != &l.root {
		t.Error("New() failed")
	}

	if l.maxLevel != 63 {
		t.Error("New() failed")
	}

	if l.p != 0.5 {
		t.Error("New() failed")
	}

	if l.less == nil {
		t.Error("New() failed")
	}

	if l.fingers == nil {
		t.Error("New() failed")
	}

	if len(l.fingers) != 64 {
		t.Error("New() failed")
	}

	if l.len != 0 {
		t.Error("New() failed")
	}

}

//func TestMaxLevel(t *testing.T) {
//	l := NewOrdered[int](64, 0.5)
//	var m = make(map[uint8]int)
//	for i := 0; i < 100000; i++ {
//		m[l.generateLevel()]++
//	}
//	fmt.Println(m)
//}

func getListItems[T any](l *SkipList[T]) []T {
	var items []T
	l.Ascend(func(item T) bool {
		items = append(items, item)
		return true
	})
	return items
}

func TestReplaceOrInsert(t *testing.T) {

	l := NewOrdered[int](63, 0.5)

	l.ReplaceOrInsert(1)

	if l.len != 1 {
		t.Errorf("ReplaceOrInsert failed: expected size 1, got %v", l.len)
	}

	if n, found := l.Get(1); !found {
		t.Errorf("ReplaceOrInsert failed: expected node with key 1, got %v and %v", n, found)
	}

	l.ReplaceOrInsert(2)

	if l.len != 2 {
		t.Errorf("ReplaceOrInsert failed: expected size 2, got %v", l.len)
	}

	if n, found := l.Get(2); !found {
		t.Errorf("ReplaceOrInsert failed: expected node with key 2, got %v and %v", n, found)
	}

	l.ReplaceOrInsert(3)
	l.ReplaceOrInsert(4)
	l.ReplaceOrInsert(5)

	items := getListItems(l)

	if !slices.Equal(items, []int{1, 2, 3, 4, 5}) {
		t.Errorf("ReplaceOrInsert failed: expected [1, 2, 3, 4, 5], got %v", items)
	}

	l.ReplaceOrInsert(10)
	l.ReplaceOrInsert(9)
	l.ReplaceOrInsert(8)
	l.ReplaceOrInsert(7)
	l.ReplaceOrInsert(6)

	items = getListItems(l)

	if !slices.Equal(items, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Errorf("ReplaceOrInsert failed: expected [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], got %v", items)
	}

	l.ReplaceOrInsert(1)

	items = getListItems(l)

	if !slices.Equal(items, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Errorf("ReplaceOrInsert failed: expected [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], got %v", items)
	}

	l2 := NewOrdered[int](63, 0.5)
	l2.ReplaceOrInsert(1)
	l2.ReplaceOrInsert(1)

	if l2.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !l2.Has(1) {
		t.Errorf("expecting to find key=1")
	}

	l2.Delete(1)
	if l2.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if l2.Has(1) {
		t.Errorf("not expecting to find key=1")
	}

	l2.Delete(1)
	if l2.Len() != 0 {
		t.Errorf("expecting len 0")
	}

	if l2.Has(1) {
		t.Errorf("not expecting to find key=1")
	}

}

func TestReverseReplaceOrInsert(t *testing.T) {
	l := NewOrdered[int](63, 0.5)
	n := 100
	for i := 0; i < n; i++ {
		l.ReplaceOrInsert(n - i)
	}
	i := 0
	l.AscendGreaterOrEqual(0, func(item int) bool {
		i++
		if item != i {
			t.Errorf("bad order: got %d, expect %d", item, i)
		}
		return true
	})
}

func TestRange(t *testing.T) {
	l := NewOrdered[string](63, 0.5)
	order := []string{
		"ab", "aba", "abc", "a", "aa", "aaa", "b", "a-", "a!",
	}
	for _, i := range order {
		l.ReplaceOrInsert(i)
	}
	k := 0
	l.AscendRange("ab", "ac", func(item string) bool {
		if k > 3 {
			t.Fatalf("returned more items than expected")
		}
		i1 := order[k]
		i2 := item
		if i1 != i2 {
			t.Errorf("expecting %s, got %s", i1, i2)
		}
		k++
		return true
	})
}

func TestRandomInsertOrder(t *testing.T) {
	l := NewOrdered[int](16, 0.5)
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		l.ReplaceOrInsert(perm[i])
	}
	j := 0
	l.AscendGreaterOrEqual(0, func(item int) bool {

		if item != j {
			t.Fatalf("bad order %d, %d", item, j)
		}
		j++
		return true
	})
}

func TestRandomReplace(t *testing.T) {
	l := NewOrdered[int](48, 0.5)
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		l.ReplaceOrInsert(perm[i])
	}
	perm = rand.Perm(n)
	for i := 0; i < n; i++ {
		if replaced, isReplace := l.ReplaceOrInsert(perm[i]); !isReplace || replaced != perm[i] {
			t.Errorf("error replacing")
		}
	}
}

func TestRandomInsertSequentialDelete(t *testing.T) {
	l := NewOrdered[int](16, 0.5)

	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		l.ReplaceOrInsert(perm[i])
	}

	for i := 0; i < n; i++ {
		l.Delete(i)
	}

}

func TestRandomInsertDeleteNonExistent(t *testing.T) {
	l := NewOrdered[int](8, 0.5)
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		l.ReplaceOrInsert(perm[i])
	}

	for i := 0; i < n; i++ {
		if u, success := l.Delete(i); !success || u != i {
			t.Errorf("delete failed")
		}
	}

	if _, success := l.Delete(200); success {
		t.Errorf("deleted non-existent item")
	}
	if _, success := l.Delete(-2); success {
		t.Errorf("deleted non-existent item")
	}

	if _, success := l.Delete(200); success {
		t.Errorf("deleted non-existent item")
	}
	if _, success := l.Delete(-2); success {
		t.Errorf("deleted non-existent item")
	}
}

func TestRandomInsertPartialDeleteOrder(t *testing.T) {
	l := NewOrdered[int](8, 0.5)
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		l.ReplaceOrInsert(perm[i])
	}

	for i := 1; i < n-1; i++ {
		l.Delete(i)

	}

	j := 0
	l.AscendGreaterOrEqual(0, func(item int) bool {
		switch j {
		case 0:
			if item != 0 {
				t.Errorf("expecting 0")
			}
		case 1:
			if item != n-1 {
				t.Errorf("expecting %d", n-1)
			}
		}
		j++
		return true
	})
}

func TestInsertNoReplace(t *testing.T) {
	l := NewOrdered[int](8, 0.5)
	n := 1000
	for q := 0; q < 2; q++ {
		perm := rand.Perm(n)
		for i := 0; i < n; i++ {
			l.InsertNoReplace(perm[i])
		}
	}

	j := 0
	l.AscendGreaterOrEqual(0, func(item int) bool {
		if item != j/2 {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}
