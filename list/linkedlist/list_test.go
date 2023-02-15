package linkedlist

import (
	"testing"
)

func traverse[T any](t *testing.T, l *List[T]) {
	c := l.Cursor()
	for c.MoveNext() != nil {
		t.Logf("%v", c.Node().Value)
	}
}

func TestList(t *testing.T) {

	l := New[any]()
	//checkListPointers(t, l, []*Node[any]{})

	// Single Node linked list
	l.PushFront("a")

	e := l.FrontCursor()

	checkListPointers(t, l, []*Node[any]{e.Node()})
	l.MoveToFront(e)
	checkListPointers(t, l, []*Node[any]{e.Node()})
	l.MoveToBack(e)
	checkListPointers(t, l, []*Node[any]{e.Node()})
	l.RemoveCurrent(e)
	checkListPointers(t, l, []*Node[any]{})

	// Bigger linked list
	l.PushFront(2)
	e2 := l.FrontCursor()
	l.PushFront(1)
	e1 := l.FrontCursor()
	l.PushBack(3)
	e3 := l.BackCursor()
	l.PushBack("banana")
	e4 := l.BackCursor()
	checkListPointers(t, l, []*Node[any]{e1.Node(), e2.Node(), e3.Node(), e4.Node()})

	l.RemoveCurrent(e2)
	checkListPointers(t, l, []*Node[any]{e1.Node(), e3.Node(), e4.Node()})

	l.MoveToFront(e3) // move from middle
	checkListPointers(t, l, []*Node[any]{e3.Node(), e1.Node(), e4.Node()})

	l.MoveToFront(e1)
	l.MoveToBack(e3) // move from middle
	checkListPointers(t, l, []*Node[any]{e1.Node(), e4.Node(), e3.Node()})

	l.MoveToFront(e3) // move from back
	checkListPointers(t, l, []*Node[any]{e3.Node(), e1.Node(), e4.Node()})
	l.MoveToFront(e3) // should be no-op
	checkListPointers(t, l, []*Node[any]{e3.Node(), e1.Node(), e4.Node()})

	l.MoveToBack(e3) // move from front
	checkListPointers(t, l, []*Node[any]{e1.Node(), e4.Node(), e3.Node()})
	l.MoveToBack(e3) // should be no-op
	checkListPointers(t, l, []*Node[any]{e1.Node(), e4.Node(), e3.Node()})

	l.InsertBefore(2, e1) // insert before front
	e2 = e1.PrevCursor()
	checkListPointers(t, l, []*Node[any]{e2.Node(), e1.Node(), e4.Node(), e3.Node()})

	l.RemoveCurrent(e2)
	l.InsertBefore(2, e4) // insert before middle
	e2 = e4.PrevCursor()

	checkListPointers(t, l, []*Node[any]{e1.Node(), e2.Node(), e4.Node(), e3.Node()})
	l.RemoveCurrent(e2)
	l.InsertBefore(2, e3) // insert before back
	e2 = e3.PrevCursor()
	checkListPointers(t, l, []*Node[any]{e1.Node(), e4.Node(), e2.Node(), e3.Node()})
	l.RemoveCurrent(e2)

	l.InsertAfter(2, e1) // insert after front
	e2 = e1.NextCursor()
	checkListPointers(t, l, []*Node[any]{e1.Node(), e2.Node(), e4.Node(), e3.Node()})
	l.RemoveCurrent(e2)
	l.InsertAfter(2, e4) // insert after middle
	e2 = e4.NextCursor()

	checkListPointers(t, l, []*Node[any]{e1.Node(), e4.Node(), e2.Node(), e3.Node()})
	l.RemoveCurrent(e2)
	l.InsertAfter(2, e3) // insert after back
	e2 = e3.NextCursor()
	checkListPointers(t, l, []*Node[any]{e1.Node(), e4.Node(), e3.Node(), e2.Node()})
	l.RemoveCurrent(e2)

	// Check standard iteration.
	sum := 0
	l.Cursor().Ascending(func(n *Node[any]) bool {
		if i, ok := n.Value.(int); ok {
			sum += i
		}
		return true
	})

	if sum != 4 {
		t.Errorf("sum over l = %d, want 4", sum)
	}

	// Clear all Nodes
	c := l.Cursor()
	for l.RemoveAfter(c) != nil {
	}

	checkListPointers(t, l, []*Node[any]{})
}

func TestExtendingV2(t *testing.T) {
	l1 := New[int]()
	l2 := New[int]()

	l1.PushBack(1)
	l1.PushBack(2)
	l1.PushBack(3)

	l2.PushBack(4)
	l2.PushBack(5)

	l3 := New[int]()
	l3.PushBackList(l1)

	checkList(t, l3, []int{1, 2, 3})

	l3.PushBackList(l2)
	checkList(t, l3, []int{1, 2, 3, 4, 5})

	l3 = New[int]()
	l3.PushFrontList(l2)

	checkList(t, l3, []int{4, 5})

	l3.PushFrontList(l1)
	checkList(t, l3, []int{1, 2, 3, 4, 5})

	checkList(t, l1, []int{1, 2, 3})
	checkList(t, l2, []int{4, 5})

	l3 = New[int]()
	l3.PushBackList(l1)
	checkList(t, l3, []int{1, 2, 3})

	l3.PushBackList(l3)
	checkList(t, l3, []int{1, 2, 3, 1, 2, 3})

	l3 = New[int]()
	l3.PushFrontList(l1)
	checkList(t, l3, []int{1, 2, 3})
	l3.PushFrontList(l3)
	checkList(t, l3, []int{1, 2, 3, 1, 2, 3})

	l3 = New[int]()
	l1.PushBackList(l3)
	checkList(t, l1, []int{1, 2, 3})
	l1.PushFrontList(l3)
	checkList(t, l1, []int{1, 2, 3})
}

func TestIssue4103(t *testing.T) {
	l1 := New[int]()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := New[int]()
	l2.PushBack(3)
	l2.PushBack(4)

	e := l1.FrontCursor()
	l2.RemoveCurrent(e) // l2 should not change because e is not a cursor of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(8, e)

	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestRemoveV2(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	e1 := l.BackCursor()
	l.PushBack(2)
	e2 := l.BackCursor()
	checkListPointers(t, l, []*Node[int]{e1.Node(), e2.Node()})

	e := l.FrontCursor()
	l.RemoveCurrent(e) // e moves to next element
	checkListPointers(t, l, []*Node[int]{e2.Node()})

	l.RemoveCurrent(e) // e moves to "dummy" element
	checkListPointers(t, l, []*Node[int]{})

}

func TestIssue6349V2(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)

	e := l.FrontCursor()
	l.RemoveCurrent(e)
	if e.Node().Value != 2 {
		t.Errorf("e.value = %d, want 1", e.Node().Value)
	}
	if e.NextNode() != nil {
		t.Errorf("e.Next() != nil")
	}
	if e.PrevNode() != nil {
		t.Errorf("e.Prev() != nil")
	}
}

func TestMoveV2(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	e1 := l.BackCursor()
	l.PushBack(2)
	e2 := l.BackCursor()
	l.PushBack(3)
	e3 := l.BackCursor()
	l.PushBack(4)
	e4 := l.BackCursor()

	l.MoveAfter(e3, e3)
	checkListPointers(t, l, []*Node[int]{e1.Node(), e2.Node(), e3.Node(), e4.Node()})
	l.MoveBefore(e2, e2)
	checkListPointers(t, l, []*Node[int]{e1.Node(), e2.Node(), e3.Node(), e4.Node()})

	l.MoveAfter(e3, e2)
	checkListPointers(t, l, []*Node[int]{e1.Node(), e2.Node(), e3.Node(), e4.Node()})
	l.MoveBefore(e2, e3)
	checkListPointers(t, l, []*Node[int]{e1.Node(), e2.Node(), e3.Node(), e4.Node()})

	l.MoveBefore(e2, e4)
	checkListPointers(t, l, []*Node[int]{e1.Node(), e3.Node(), e2.Node(), e4.Node()})
	e2, e3 = e3, e2

	l.MoveBefore(e4, e1)
	checkListPointers(t, l, []*Node[int]{e4.Node(), e1.Node(), e2.Node(), e3.Node()})
	e1, e2, e3, e4 = e4, e1, e2, e3

	l.MoveAfter(e4, e1)
	checkListPointers(t, l, []*Node[int]{e1.Node(), e4.Node(), e2.Node(), e3.Node()})
	e2, e3, e4 = e4, e2, e3

	l.MoveAfter(e2, e3)
	checkListPointers(t, l, []*Node[int]{e1.Node(), e3.Node(), e2.Node(), e4.Node()})
}

func TestZeroListV2(t *testing.T) {
	var l1 = new(List[int])
	l1.PushFront(1)
	checkList(t, l1, []int{1})

	var l2 = new(List[int])
	l2.PushBack(1)
	checkList(t, l2, []int{1})

	var l3 = new(List[int])
	l3.PushFrontList(l1)
	checkList(t, l3, []int{1})

	var l4 = new(List[int])
	l4.PushBackList(l2)
	checkList(t, l4, []int{1})
}

// Test that a linkedlist l is not modified when calling InsertBefore with a mark that is not an Node of l.
func TestInsertBeforeUnknownMark(t *testing.T) {
	var l List[int]
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.InsertBefore(1, new(Cursor[int]))
	checkList(t, &l, []int{1, 2, 3})
}

// Test that a linkedlist l is not modified when calling InsertAfter with a mark that is not an Node of l.
func TestInsertAfterUnknownMark(t *testing.T) {
	var l List[int]
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.InsertAfter(1, new(Cursor[int]))
	checkList(t, &l, []int{1, 2, 3})
}

// Test that a linkedlist l is not modified when calling MoveAfter or MoveBefore with a mark that is not an Node of l.
func TestMoveUnknownMark(t *testing.T) {
	var l1 List[int]
	l1.PushBack(1)
	c1 := l1.BackCursor()

	var l2 List[int]
	l2.PushBack(2)
	c2 := l2.BackCursor()

	l1.MoveAfter(c1, c2)
	checkList(t, &l1, []int{1})
	checkList(t, &l2, []int{2})

	l1.MoveBefore(c1, c2)
	checkList(t, &l1, []int{1})
	checkList(t, &l2, []int{2})
}
