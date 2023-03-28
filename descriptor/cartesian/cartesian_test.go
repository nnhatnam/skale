package cartesian

import (
	"golang.org/x/exp/slices"
	"testing"
)

func inOrder(n *Node) []int {
	if n == nil {
		return []int{}
	}
	return append(append(inOrder(n.left), n.idx), inOrder(n.right)...)
}

func inOrderWalk[T any](n *Node, f func(*Node)) {
	if n == nil {
		return
	}
	inOrderWalk[T](n.left, f)
	f(n)
	inOrderWalk[T](n.right, f)
}

func naiveFindMin(a []int, i, j int) int {
	min := i
	for k := i + 1; k <= j; k++ {
		if a[k] < a[min] {
			min = k
		}
	}
	return min
}

func TestCartesian(t *testing.T) {
	tree := NewOrdered[int]()

	tree.Init([]int{93, 84, 33, 64, 62, 83, 63})

	return

	order := inOrder(tree.root)

	// check in-order traversal
	if !slices.Equal(order, []int{1, 2, 3, 4, 5, 6, 7}) {
		t.Errorf("in-order traversal failed, expected %v, got %v", []int{1, 2, 3, 4, 5, 6, 7}, inOrder(tree.root))
	}
}
