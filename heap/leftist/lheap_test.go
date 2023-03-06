package leftist

import (
	"golang.org/x/exp/slices"
	"testing"
)

func walkInOrder[T any](n *Node[T], f func(*Node[T])) {
	if n == nil {
		return
	}

	walkInOrder(n.left, f)
	f(n)
	walkInOrder(n.right, f)
}

func treeItems[T any](n *Node[T]) []T {
	var result []T
	walkInOrder(n, func(n *Node[T]) {
		result = append(result, n.Value)
	})

	return result
}

func TestLHeap(t *testing.T) {
	h := NewOrdered[int]()

	h.Insert(1)

	if h.root.Value != 1 {
		t.Errorf("root value is %d, want 1", h.root.Value)
	}

	if h.Len() != 1 {
		t.Errorf("heap length is %d, want 1", h.Len())
	}

	h.Insert(2)

	if h.Len() != 2 {
		t.Errorf("heap length is %d, want 2", h.Len())
	}

	if h.root.Value != 1 {
		t.Errorf("root value is %d, want 1", h.root.Value)
	}

	if h.root.left.Value != 2 {
		t.Errorf("root.left value is %d, want 2", h.root.left.Value)
	}

	h = NewOrdered[int]()
	// test insertion steps

	// 50
	h.Insert(50)

	if h.root.Value != 50 || h.root.npl != 0 {
		t.Errorf("root value is incorrect, want 50 , 0 but got %v %v", h.root.Value, h.root.npl)
	}

	// Normal insertion of a new node into the tree (value = 75)
	// First, insert as far right as possible, then swing left. Result
	//			50
	//    75
	h.Insert(75)

	if h.root.left == nil || h.root.left.Value != 75 || h.root.left.npl != 0 {
		t.Errorf("root.left value is %d, want 75", h.root.left.Value)
	}

	if h.root.npl != 0 {
		t.Errorf("root.npl is %d, want 0", h.root.npl)
	}

	// Normal insertion of a new node into the tree (value = 25)
	// As this is smaller than the root, it becomes the new root
	// Then swing left to satisfy the npl rule. Result
	//		  25
	//    	50
	//    75
	h.Insert(25)

	if h.root.Value != 25 || h.root.npl != 0 {
		t.Errorf("root value is incorrect, want 25 , 0 but got %v %v", h.root.Value, h.root.npl)
	}

	if h.root.left.Value != 50 || h.root.left.npl != 0 {
		t.Errorf("root.left value is incorrect, want 50, 0 but got %v %v", h.root.left.Value, h.root.left.npl)
	}

	if h.root.left.left.Value != 75 || h.root.left.left.npl != 0 {
		t.Errorf("root.left.left value is incorrect, want 75, 0 but got %v %v", h.root.left.left.Value, h.root.left.left.npl)
	}

	// Normal insertion of a new node into the tree (value = 55)
	// No swing required as the npl rule is already satisfied. Result
	//		  25
	//    	50	55
	//    75
	h.Insert(55)

	if h.root.Value != 25 || h.root.npl != 1 {
		t.Errorf("root value is incorrect, want 25 , 1 but got %v %v", h.root.Value, h.root.npl)
	}

	if h.root.left.Value != 50 || h.root.left.npl != 0 {
		t.Errorf("root.left value is incorrect, want 50, 0 but got %v %v", h.root.left.Value, h.root.left.npl)
	}

	if h.root.left.left.Value != 75 || h.root.left.left.npl != 0 {
		t.Errorf("root.left.left value is incorrect, want 75, 0 but got %v %v", h.root.left.left.Value, h.root.left.left.npl)
	}

	if h.root.right.Value != 55 || h.root.right.npl != 0 {
		t.Errorf("root.right value is incorrect, want 55, 0 but got %v %v", h.root.right.Value, h.root.right.npl)
	}

	// Normal insertion of a new node into the tree (value = 40). Result
	//  		25
	//    	   /   \
	//    	  50   40
	//       /     /
	//      75	  55
	h.Insert(40)

	if h.root.Value != 25 || h.root.npl != 1 {
		t.Errorf("root value is incorrect, want 25 , 1 but got %v %v", h.root.Value, h.root.npl)
	}

	if h.root.left.Value != 50 || h.root.left.npl != 0 {
		t.Errorf("root.left value is incorrect, want 50, 0 but got %v %v", h.root.left.Value, h.root.left.npl)
	}

	if h.root.left.left.Value != 75 || h.root.left.left.npl != 0 {
		t.Errorf("root.left.left value is incorrect, want 75, 0 but got %v %v", h.root.left.left.Value, h.root.left.left.npl)
	}

	if h.root.right.Value != 40 || h.root.right.npl != 0 {
		t.Errorf("root.right value is incorrect, want 40, 0 but got %v %v", h.root.right.Value, h.root.right.npl)
	}

	if h.root.right.left.Value != 55 || h.root.right.left.npl != 0 {
		t.Errorf("root.right.left value is incorrect, want 55, 0 but got %v %v", h.root.right.left.Value, h.root.right.left.npl)
	}

	// Normal insertion of a new node into the tree (value = 65). First, we will get the tree on the left, but then we will swing left to satisfy the npl rule. R
	//  		25						  25
	//    	   /   \				 	/   \
	//    	  50   40       -> 		   40   50
	//       /     / \				  / \   /
	//      75	  55 65				 55 65 75

	h.Insert(65)

	if h.root.Value != 25 || h.root.npl != 1 {
		t.Errorf("root value is incorrect, want 25 , 1 but got %v %v", h.root.Value, h.root.npl)
	}

	if h.root.left.Value != 40 || h.root.left.npl != 1 {
		t.Errorf("root.left value is incorrect, want 40, 0 but got %v %v", h.root.left.Value, h.root.left.npl)
	}

	if h.root.left.left.Value != 55 || h.root.left.left.npl != 0 {
		t.Errorf("root.left.left value is incorrect, want 55, 0 but got %v %v", h.root.left.left.Value, h.root.left.left.npl)
	}

	if h.root.left.right.Value != 65 || h.root.left.right.npl != 0 {
		t.Errorf("root.left.right value is incorrect, want 65, 0 but got %v %v", h.root.left.right.Value, h.root.left.right.npl)
	}

	if h.root.right.Value != 50 || h.root.right.npl != 0 {
		t.Errorf("root.right value is incorrect, want 50, 0 but got %v %v", h.root.right.Value, h.root.right.npl)
	}

	if h.root.right.left.Value != 75 || h.root.right.left.npl != 0 {
		t.Errorf("root.right.left value is incorrect, want 75, 0 but got %v %v", h.root.right.left.Value, h.root.right.left.npl)
	}

	h = NewOrdered[int]()

	h.InsertBulk(21, 14, 17, 10, 3, 23, 26, 8)

	items := treeItems(h.root)
	if !slices.Equal(items, []int{21, 14, 17, 10, 3, 26, 23, 8}) {
		t.Errorf("items are not in the right order, want %v but got %v", []int{21, 14, 17, 10, 3, 26, 23, 8}, items)
	}

	if h.Len() != 8 {
		t.Errorf("heap length is incorrect, want 8 but got %v", h.Len())
	}

}

func TestIsEmpty(t *testing.T) {
	h := NewOrdered[int]()

	if !h.IsEmpty() {
		t.Errorf("heap should be empty")
	}

	h.Insert(1)

	if h.IsEmpty() {
		t.Errorf("heap should not be empty")
	}
}

func TestMerge(t *testing.T) {

	h1 := NewOrdered[int]()

	h1.InsertBulk(10, 5, 15, 1, 50, 20, 99, 7, 25)
	items := treeItems(h1.root)

	if !slices.Equal(items, []int{10, 5, 15, 1, 50, 20, 99, 7, 25}) {
		t.Errorf("items are not in the right order, want %v but got %v", []int{10, 5, 15, 1, 50, 20, 99, 7, 25}, items)
	}

	h2 := NewOrdered[int]()
	h2.InsertBulk(75, 22)

	items = treeItems(h2.root)

	if !slices.Equal(items, []int{75, 22}) {
		t.Errorf("items are not in the right order, want %v but got %v", []int{75, 22}, items)
	}
	h1.Merge(h2)

	items = treeItems(h1.root)

	if !slices.Equal(items, []int{50, 20, 99, 7, 75, 22, 25, 1, 10, 5, 15}) {
		t.Errorf("items are not in the right order, want %v but got %v", []int{50, 20, 99, 7, 75, 22, 25, 1, 10, 5, 15}, items)
	}

	if h1.Len() != 11 {
		t.Errorf("heap length is incorrect, want 11 but got %v", h1.Len())
	}

}

func TestDeleteMin(t *testing.T) {

	h := NewOrdered[int]()

	h.InsertBulk(14, 8, 23, 3, 21, 10, 26, 17)

	var popItems []int

	//use Len as a exit condition to test h.Len()
	for h.Len() > 0 {
		item, found := h.DeleteMin()
		if !found {
			t.Errorf("error: item not found")
		}
		popItems = append(popItems, item)
	}

	if !slices.Equal(popItems, []int{3, 8, 10, 14, 17, 21, 23, 26}) {
		t.Errorf("items are not in the right order, want %v but got %v", []int{3, 8, 10, 14, 17, 21, 23, 26}, popItems)
	}

}
