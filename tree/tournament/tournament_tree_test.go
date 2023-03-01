package tournament

import (
	"golang.org/x/exp/slices"
	"testing"
)

func TestNew(t *testing.T) {

	// sorted ascending
	tree := NewOrdered[int]()
	if tree == nil {
		t.Errorf("New() returns nil")
	}

	tree.InsertNoReplace(1)

	if !slices.Equal(tree.e, []int{1}) {
		t.Errorf("Insert failed, expected %v, got %v", []int{1}, tree.e)
	}

	tree.InsertNoReplace(2)

	if !slices.Equal(tree.e, []int{1, 1, 2}) {
		t.Errorf("Insert failed, expected %v, got %v", []int{1, 1, 2}, tree.e)
	}

	tree.InsertNoReplace(3)

	if !slices.Equal(tree.e, []int{1, 1, 2, 1, 3}) {
		t.Errorf("Insert failed, expected %v, got %v", []int{1, 1, 2, 1, 3}, tree.e)
	}

	// sorted descending
	tree = NewOrdered[int]()

	tree.InsertNoReplace(3)
	tree.InsertNoReplace(2)

	if !slices.Equal(tree.e, []int{2, 3, 2}) {
		t.Errorf("Insert failed, expected %v, got %v", []int{3, 3, 2}, tree.e)
	}

	tree.InsertNoReplace(1)

	if !slices.Equal(tree.e, []int{1, 1, 2, 3, 1}) {
		t.Errorf("Insert failed, expected %v, got %v", []int{1, 1, 2, 3, 1}, tree.e)
	}

	tree = NewOrdered[int]()

	tree.InsertNoReplaceBulk(3, 5, 6, 7, 20, 8, 2, 9)

	if !slices.Equal(tree.e, []int{2, 3, 2, 3, 6, 2, 7, 3, 20, 6, 8, 5, 2, 7, 9}) {
		t.Errorf("Insert failed, expected %v, got %v", []int{2, 3, 2, 3, 6, 2, 7, 3, 20, 6, 8, 5, 2, 7, 9}, tree.e)
	}

	tree = NewOrdered[int]()

	tree.InsertNoReplaceBulk(9, 5, 10, 7)

	if !slices.Equal(tree.e, []int{5, 9, 5, 9, 10, 5, 7}) {
		t.Errorf("Insert failed, expected %v, got %v", []int{5, 9, 5, 9, 10, 5, 7}, tree.e)
	}

	tree.InsertNoReplace(8)

	if !slices.Equal(tree.e, []int{5, 8, 5, 8, 10, 5, 7, 9, 8}) {
		t.Errorf("Insert failed, expected %v, got %v", []int{5, 8, 5, 8, 10, 5, 7, 9, 8}, tree.e)
	}

}

func TestDelete(t *testing.T) {

	tree := NewOrdered[int]()
	if tree == nil {
		t.Errorf("New() returns nil")
	}

	tree.InsertNoReplaceBulk(1, 2, 3)

	v, success := tree.DeleteMin()

	if !slices.Equal(tree.e, []int{2, 3, 2}) {
		t.Errorf("DeleteMin failed, expected %v, got %v", []int{2, 3, 2}, tree.e)
	}

	if v != 1 || success != true {
		t.Errorf("DeleteMin failed, expected %v, got %v", 1, v)
	}

	v, success = tree.DeleteMin()

	if !slices.Equal(tree.e, []int{3}) {
		t.Errorf("DeleteMin failed, expected %v, got %v", []int{3}, tree.e)
	}

	if v != 2 || success != true {
		t.Errorf("DeleteMin failed, expected %v, got %v", 2, v)
	}

	v, success = tree.DeleteMin()

	if !slices.Equal(tree.e, []int{}) {
		t.Errorf("DeleteMin failed, expected %v, got %v", []int{}, tree.e)
	}

	if v != 3 || success != true {
		t.Errorf("DeleteMin failed, expected %v, got %v", 3, v)
	}

	v, success = tree.DeleteMin()

	if !slices.Equal(tree.e, []int{}) {
		t.Errorf("DeleteMin failed, expected %v, got %v", []int{}, tree.e)
	}

	if v != 0 || success != false {
		t.Errorf("DeleteMin failed, expected %v, got %v", 0, v)
	}

}
