package linkedlist

import (
	"fmt"
	"testing"
)

func traverse[T any](t *testing.T, l *List[T]) {
	c := l.Cursor()
	for c.MoveNext() != nil {
		t.Logf("%v", c.Node().Value)
	}
}

func debugPrinter[T any](n1, n2 string, es []*Node[T]) {
	fmt.Printf("%v -> %v is: [%p](%v) -> [%p](%v)\n", n1, n2, es[0], *es[0], es[1], *es[1])
}

func checkListLen[T any](t *testing.T, l *List[T], len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkListPointers[T any](t *testing.T, l *List[T], es []*Node[T]) {

	root := &l.root
	if !checkListLen(t, l, len(es)) {
		return
	}

	// zero length lists must be the zero value or properly initialized (sentinel circle)
	if len(es) == 0 {
		if l.root.next != nil && l.root.next != root || l.root.prev != nil && l.root.prev != root {
			t.Errorf("l.root.next = %p, l.root.prev = %p; both should both be nil or %p", l.root.next, l.root.prev, root)
		}
		return
	}
	//len(es) > 0

	//check internal and external prev/next connections
	for i, e := range es {
		prev := root
		Prev := (*Node[T])(nil)
		if i > 0 {
			prev = es[i-1]
			Prev = prev
		}
		if p := e.prev; p != prev {
			t.Errorf("elt[%d](%p).prev = %p, want %p ??", i, e, p, prev)
		}

		if p := e.prev; i > 0 && p != Prev {
			t.Errorf("elt[%d](%p).Prev() = %p, want %p", i, e, p, Prev)
		}

		next := root
		Next := (*Node[T])(nil)
		if i < len(es)-1 {
			next = es[i+1]
			Next = next
		}
		if n := e.next; n != next {
			t.Errorf("elt[%d](%p).next = %p, want %p", i, e, n, next)
		}
		if n := e.next; i < len(es)-1 && n != Next {
			t.Errorf("elt[%d](%p).Next() = %p, want %p", i, e, n, Next)
		}
	}
}

func checkList[T comparable](t *testing.T, l *List[T], es []T) {
	if !checkListLen(t, l, len(es)) {
		return
	}

	i := 0
	for e := l.front(); e != &l.root; e = e.next {
		le := e.Value
		if le != es[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, le, es[i])
		}
		i++
	}
}

func TestNew(t *testing.T) {
	l := New[any]()
	if l == nil {
		t.Errorf("New[any]() = nil, want non-nil")
	}

	if l.Len() != 0 {
		t.Errorf("l.Len() = %d, want 0", l.Len())
	}

	if l.Front() != nil {
		t.Errorf("l.Front() = %v, want nil", l.Front())
	}

	if l.Back() != nil {
		t.Errorf("l.Back() = %v, want nil", l.Back())
	}

}

func TestFrom(t *testing.T) {

	l := From[int]([]int{1, 2, 3, 4, 5}...)

	if l == nil {
		t.Errorf("From[int]([]int{1, 2, 3, 4, 5}...) = nil, want non-nil")
	}

	if l.Len() != 5 {
		t.Errorf("l.Len() = %d, want 5", l.Len())
	}

	if l.Front().Value != 1 {
		t.Errorf("l.Front() = %v, want 1", l.Front())
	}

	if l.Back().Value != 5 {
		t.Errorf("l.Back() = %v, want 5", l.Back())
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
	l.RemoveAt(e)
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

	l.RemoveAt(e2)
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
	e2 = e1.ClonePrev()
	checkListPointers(t, l, []*Node[any]{e2.Node(), e1.Node(), e4.Node(), e3.Node()})

	l.RemoveAt(e2)
	l.InsertBefore(2, e4) // insert before middle
	e2 = e4.ClonePrev()

	checkListPointers(t, l, []*Node[any]{e1.Node(), e2.Node(), e4.Node(), e3.Node()})
	l.RemoveAt(e2)
	l.InsertBefore(2, e3) // insert before back
	e2 = e3.ClonePrev()
	checkListPointers(t, l, []*Node[any]{e1.Node(), e4.Node(), e2.Node(), e3.Node()})
	l.RemoveAt(e2)

	l.InsertAfter(2, e1) // insert after front
	e2 = e1.CloneNext()
	checkListPointers(t, l, []*Node[any]{e1.Node(), e2.Node(), e4.Node(), e3.Node()})
	l.RemoveAt(e2)
	l.InsertAfter(2, e4) // insert after middle
	e2 = e4.CloneNext()

	checkListPointers(t, l, []*Node[any]{e1.Node(), e4.Node(), e2.Node(), e3.Node()})
	l.RemoveAt(e2)
	l.InsertAfter(2, e3) // insert after back
	e2 = e3.CloneNext()
	checkListPointers(t, l, []*Node[any]{e1.Node(), e4.Node(), e3.Node(), e2.Node()})
	l.RemoveAt(e2)

	// Check standard iteration.
	sum := 0
	l.Cursor().WalkAscending(func(n *Node[any]) bool {
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
	l2.RemoveAt(e) // l2 should not change because e is not a cursor of l2
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
	l.RemoveAt(e) // e moves to next element
	checkListPointers(t, l, []*Node[int]{e2.Node()})

	l.RemoveAt(e) // e moves to "dummy" element
	checkListPointers(t, l, []*Node[int]{})

}

func TestIssue6349V2(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)

	e := l.FrontCursor()
	l.RemoveAt(e)
	if e.Node().Value != 2 {
		t.Errorf("e.value = %d, want 1", e.Node().Value)
	}
	if e.NodeNext() != nil {
		t.Errorf("e.Next() != nil")
	}
	if e.NodePrev() != nil {
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

// Test that a linked list l is not modified when calling InsertBefore with a mark that is not an Node of l.
func TestInsertBeforeUnknownMark(t *testing.T) {
	var l List[int]
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.InsertBefore(1, new(Cursor[int]))
	checkList(t, &l, []int{1, 2, 3})
}

// Test that a linked list l is not modified when calling InsertAfter with a mark that is not an Node of l.
func TestInsertAfterUnknownMark(t *testing.T) {
	var l List[int]
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.InsertAfter(1, new(Cursor[int]))
	checkList(t, &l, []int{1, 2, 3})
}

// Test that a linked list l is not modified when calling MoveAfter or MoveBefore with a mark that is not an Node of l.
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

func TestPopFront(t *testing.T) {
	var l = &List[int]{}

	fmt.Println(l.PopFront())

	if n := l.PopFront(); n != nil {
		t.Errorf("PopFront() = %v, want nil", n)
	}

	l = New[int]()
	if n := l.PopFront(); n != nil {
		t.Errorf("PopFront() = %v, want nil", n)
	}

	l.PushBack(1)
	if n := l.PopFront(); n.Value != 1 {
		t.Errorf("PopFront() = %v, want 1", n)
	}

	l.PushBack(1)
	l.PushBack(2)
	if n := l.PopFront(); n.Value != 1 {
		t.Errorf("PopFront() = %v, want 1", n)
	}

	if n := l.PopFront(); n.Value != 2 {
		t.Errorf("PopFront() = %v, want 2", n)
	}

	if n := l.PopFront(); n != nil {
		t.Errorf("PopFront() = %v, want nil", n)
	}
}

func TestPopBack(t *testing.T) {

	var l = &List[int]{}
	if n := l.PopBack(); n != nil {
		t.Errorf("PopBack() = %v, want nil", n)
	}

	l = New[int]()
	if n := l.PopBack(); n != nil {
		t.Errorf("PopBack() = %v, want nil", n)
	}

	l.PushBack(1)
	if n := l.PopBack(); n.Value != 1 {
		t.Errorf("PopBack() = %v, want 1", n)
	}

	l.PushBack(1)
	l.PushBack(2)

	if n := l.PopBack(); n.Value != 2 {
		t.Errorf("PopBack() = %v, want 2", n)
	}

	if n := l.PopBack(); n.Value != 1 {
		t.Errorf("PopBack() = %v, want 1", n)
	}

	if n := l.PopBack(); n != nil {
		t.Errorf("PopBack() = %v, want nil", n)
	}
}
