package avl

import (
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
	//test cases from https://github.com/petar/GoLLRB

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

func TestTreeGet(t *testing.T) {

	// test cases from Gods(https://github.com/emirpasic/gods)
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
	tree.print()

	//	AVLTree
	//	│       ┌── 6 (0)
	//	│   ┌── 5 (-1)
	//	└── 4 (0)
	//		│   ┌── 3 (0)
	//		└── 2 (0)
	//			└── 1 (0)

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

	if count := countFrom(n); count != 3 {
		t.Errorf("Got %v expected %v", count, 3)
	}

	n = tree.find(4)
	if n == nil {
		t.Errorf("Got nil expected node with value 4")
	}

	if n.Value != 4 {
		t.Errorf("Got %v expected %v", n.Value, 4)
	}

	if count := countFrom(n); count != 6 {
		t.Errorf("Got %v expected %v", count, 6)
	}

	n = tree.find(7)
	if n != nil {
		t.Errorf("Got node expected nil")
	}

	if count := countFrom(n); count != 0 {
		t.Errorf("Got %v expected %v", count, 0)
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
