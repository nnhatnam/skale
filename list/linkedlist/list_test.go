package linkedlist

import (
	"testing"
)

func traverse[T any](t *testing.T, l *List[T]) {
	c := l.Cursor()
	for c.MoveNext() != nil {
		t.Logf("%v", c.Current().Value)
	}
}

func TestIssue4103(t *testing.T) {
	l1 := New[int]()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := New[int]()
	l2.PushBack(3)
	l2.PushBack(4)

	e := l1.CursorFront()
	l2.RemoveCurrent(e) // l2 should not change because e is not an Node of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(8, e)

	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

//func TestZeroListV2(t *testing.T) {
//	var l1 = new(List[int])
//	l1.PushFront(1)
//	checkList(t, l1, []int{1})
//
//	var l2 = new(List[int])
//	l2.PushBack(1)
//	checkList(t, l2, []int{1})
//
//	var l3 = new(List[int])
//	l3.PushFrontList(l1)
//	checkList(t, l3, []int{1})
//
//	var l4 = new(List[int])
//	l4.PushBackList(l2)
//	checkList(t, l4, []int{1})
//}

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
	c1 := l1.CursorBack()

	var l2 List[int]
	l2.PushBack(2)
	c2 := l2.CursorBack()

	l1.MoveAfter(c1, c2)
	checkList(t, &l1, []int{1})
	checkList(t, &l2, []int{2})

	l1.MoveBefore(c1, c2)
	checkList(t, &l1, []int{1})
	checkList(t, &l2, []int{2})
}
