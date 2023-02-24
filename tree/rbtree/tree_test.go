package rbtree

import (
	"math/rand"
	"testing"
)

// countFrom counts the number of nodes in the tree with root n
func countFrom[T any](n *Node[T]) int {
	if n == nil {
		return 0
	}

	return 1 + countFrom(n.Left) + countFrom(n.Right)
}

func TestCases(t *testing.T) {

	tree := NewOrdered[int]()
	tree.ReplaceOrInsert(1)
	tree.ReplaceOrInsert(1)

	if tree.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !tree.Has(1) {
		t.Errorf("expecting to find key=1")
	}

	tree.Delete(1)
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(1) {
		t.Errorf("not expecting to find key=1")
	}

	tree.Delete(1)
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}

	if tree.Has(1) {
		t.Errorf("not expecting to find key=1")
	}

}

func TestReverseInsertOrder(t *testing.T) {
	tree := NewOrdered[int]()
	n := 100
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(n - i)
	}
	i := 0
	tree.AscendGreaterOrEqual(0, func(item int) bool {
		i++
		if item != i {
			t.Errorf("bad order: got %d, expect %d", item, i)
		}
		return true
	})
}

func TestRange(t *testing.T) {
	tree := NewOrdered[string]()
	order := []string{
		"ab", "aba", "abc", "a", "aa", "aaa", "b", "a-", "a!",
	}
	for _, i := range order {
		tree.ReplaceOrInsert(i)
	}
	k := 0
	tree.AscendRange("ab", "ac", func(item string) bool {
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
	tree := NewOrdered[int]()
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(perm[i])
	}
	j := 0
	tree.AscendGreaterOrEqual(0, func(item int) bool {

		if item != j {
			t.Fatalf("bad order %d, %d", item, j)
		}
		j++
		return true
	})
}

func TestRandomReplace(t *testing.T) {
	tree := NewOrdered[int]()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(perm[i])
	}
	perm = rand.Perm(n)
	for i := 0; i < n; i++ {
		if replaced, isReplace := tree.ReplaceOrInsert(perm[i]); !isReplace || replaced != perm[i] {
			t.Errorf("error replacing")
		}
	}
}

func TestRandomInsertSequentialDelete(t *testing.T) {
	tree := NewOrdered[int]()

	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(perm[i])
	}

	for i := 0; i < n; i++ {
		tree.Delete(i)
	}

}

func TestRandomInsertDeleteNonExistent(t *testing.T) {
	tree := NewOrdered[int]()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(perm[i])
	}

	for i := 0; i < n; i++ {
		if u, success := tree.Delete(i); !success || u != i {
			t.Errorf("delete failed")
		}
	}

	if _, success := tree.Delete(200); success {
		t.Errorf("deleted non-existent item")
	}
	if _, success := tree.Delete(-2); success {
		t.Errorf("deleted non-existent item")
	}

	if _, success := tree.Delete(200); success {
		t.Errorf("deleted non-existent item")
	}
	if _, success := tree.Delete(-2); success {
		t.Errorf("deleted non-existent item")
	}
}

func TestRandomInsertPartialDeleteOrder(t *testing.T) {
	tree := NewOrdered[int]()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(perm[i])
	}

	for i := 1; i < n-1; i++ {
		tree.Delete(i)

	}

	j := 0
	tree.AscendGreaterOrEqual(0, func(item int) bool {
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

////func TestRandomInsertStats(t *testing.T) {
////	tree := NewOrdered[int]()
////	n := 100000
////	perm := rand.Perm(n)
////	for i := 0; i < n; i++ {
////		tree.ReplaceOrInsert(perm[i])
////	}
////	avg, _ := tree.HeightStats()
////	expAvg := math.Log2(float64(n)) - 1.5
////	if math.Abs(avg-expAvg) >= 2.0 {
////		t.Errorf("too much deviation from expected average height")
////	}
////}

func BenchmarkInsert(b *testing.B) {
	tree := NewOrdered[int]()
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(b.N - i)
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	tree := NewOrdered[int]()
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(b.N - i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Delete(i)
	}
}

func BenchmarkDeleteMin(b *testing.B) {
	b.StopTimer()
	tree := NewOrdered[int]()
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(b.N - i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.DeleteMin()
	}
}

func TestInsertNoReplace(t *testing.T) {
	tree := NewOrdered[int]()
	n := 1000
	for q := 0; q < 2; q++ {
		perm := rand.Perm(n)
		for i := 0; i < n; i++ {
			tree.InsertNoReplace(perm[i])
		}
	}
	j := 0
	tree.AscendGreaterOrEqual(0, func(item int) bool {
		if item != j/2 {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}

func TestTreeGet(t *testing.T) {

	tree := NewOrdered[int]()

	if count := tree.Len(); count != 0 {
		t.Errorf("Expected count 0, got %v", count)
	}

	if v, found := tree.Get(2); v != 0 || found {
		t.Errorf("Expected 0 and false, got %v and %v", v, found)
	}

	tree.ReplaceOrInsert(1)
	tree.ReplaceOrInsert(2)
	tree.ReplaceOrInsert(1)
	tree.ReplaceOrInsert(3)
	tree.ReplaceOrInsert(4)
	tree.ReplaceOrInsert(5)
	tree.ReplaceOrInsert(6)

	if count := tree.Len(); count != 6 {
		t.Errorf("Got %v expected %v", count, 6)
	}

	n := tree.find(2)
	if n == nil {
		t.Errorf("Got nil expected node with value 2")
	}

	if n.Value != 2 {
		t.Errorf("Got %v expected %v", n.Value, 2)
	}

	n = tree.find(4)
	if n == nil {
		t.Errorf("Got nil expected node with value 4")
	}

	if n.Value != 4 {
		t.Errorf("Got %v expected %v", n.Value, 4)
	}

	n = tree.find(7)
	if n != nil {
		t.Errorf("Got node expected nil")
	}

	if v, found := tree.Get(2); v != 2 || !found {
		t.Errorf("Expected 3 and true, got %v and %v", v, found)
	}

	if v, found := tree.Get(4); v != 4 || !found {
		t.Errorf("Expected 4 and true, got %v and %v", v, found)
	}

	if v, found := tree.Get(7); v != 0 || found {
		t.Errorf("Expected 0 and false, got %v and %v", v, found)
	}

}
