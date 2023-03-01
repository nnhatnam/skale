package tournament

import (
	"fmt"
	"golang.org/x/exp/slices"
	"testing"
)

func TestNewTree(t *testing.T) {

	winTree := NewOrderedWT[int]()

	if winTree == nil {
		t.Errorf("New() = nil, want a tree")
	}

	winTree.Init([]int{1})

	if len(winTree.e) != 0 {
		t.Errorf("Init() = %v, want []", winTree.e)
	}

	winTree.Init([]int{1, 2})

	if !slices.Equal(winTree.e, []int{0}) {
		t.Errorf("expected [0], got %v", winTree.e)
	}

	winTree.Init([]int{1, 2, 3})

	if !slices.Equal(winTree.e, []int{0, 1}) {
		t.Errorf("expected [0, 1], got %v", winTree.e)
	}

	winTree.Init([]int{1, 2, 3, 4})

	if !slices.Equal(winTree.e, []int{0, 0, 2}) {
		t.Errorf("expected [0, 0, 2], got %v", winTree.e)
	}

	winTree.Init([]int{1, 2, 3, 4, 5})

	if !slices.Equal(winTree.e, []int{0, 0, 1, 3}) {
		t.Errorf("expected [0, 0, 1, 3], got %v", winTree.e)
	}

	winTree.Init([]int{3, 5, 6, 7, 20, 8, 2, 9})

	if !slices.Equal(winTree.e, []int{6, 0, 6, 0, 2, 5, 6}) {
		t.Errorf("expected [6, 0, 6, 0, 2, 5, 6], got %v", winTree.e)
	}

	fmt.Println(winTree.e)

}
