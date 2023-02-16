package linkedlist

import (
	"fmt"
	"testing"
)

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

func TestOriginalList(t *testing.T) {

	l := New[any]()
	checkListPointers(t, l, []*Node[any]{})

	// Single Node linked list
	e := l.pushFront("a")

	//debugPrinter("e", "l root", []*Node[any]{e, &l.root})

	checkListPointers(t, l, []*Node[any]{e})
	l.moveToFront(e)
	checkListPointers(t, l, []*Node[any]{e})
	l.moveToBack(e)
	checkListPointers(t, l, []*Node[any]{e})
	l.remove(e)
	checkListPointers(t, l, []*Node[any]{})

	// Bigger linkedlist
	e2 := l.pushFront(2)
	e1 := l.pushFront(1)
	e3 := l.pushBack(3)
	e4 := l.pushBack("banana")
	checkListPointers(t, l, []*Node[any]{e1, e2, e3, e4})

	l.remove(e2)
	checkListPointers(t, l, []*Node[any]{e1, e3, e4})

	l.moveToFront(e3) // move from middle
	checkListPointers(t, l, []*Node[any]{e3, e1, e4})

	l.moveToFront(e1)
	l.moveToBack(e3) // move from middle
	checkListPointers(t, l, []*Node[any]{e1, e4, e3})

	l.moveToFront(e3) // move from back
	checkListPointers(t, l, []*Node[any]{e3, e1, e4})
	l.moveToFront(e3) // should be no-op
	checkListPointers(t, l, []*Node[any]{e3, e1, e4})

	l.moveToBack(e3) // move from front
	checkListPointers(t, l, []*Node[any]{e1, e4, e3})
	l.moveToBack(e3) // should be no-op
	checkListPointers(t, l, []*Node[any]{e1, e4, e3})

	e2 = l.insertBefore(2, e1) // insert before front
	checkListPointers(t, l, []*Node[any]{e2, e1, e4, e3})
	l.remove(e2)
	e2 = l.insertBefore(2, e4) // insert before middle
	checkListPointers(t, l, []*Node[any]{e1, e2, e4, e3})
	l.remove(e2)
	e2 = l.insertBefore(2, e3) // insert before back
	checkListPointers(t, l, []*Node[any]{e1, e4, e2, e3})
	l.remove(e2)

	e2 = l.insertAfter(2, e1) // insert after front
	checkListPointers(t, l, []*Node[any]{e1, e2, e4, e3})
	l.remove(e2)
	e2 = l.insertAfter(2, e4) // insert after middle
	checkListPointers(t, l, []*Node[any]{e1, e4, e2, e3})
	l.remove(e2)
	e2 = l.insertAfter(2, e3) // insert after back
	checkListPointers(t, l, []*Node[any]{e1, e4, e3, e2})
	l.remove(e2)

	// Check standard iteration.
	sum := 0
	c := l.Cursor()
	for c.MoveNext() != nil {
		if i, ok := c.Node().Value.(int); ok {
			sum += i
		}
	}

	if sum != 4 {
		t.Errorf("sum over l = %d, want 4", sum)
	}

	// Clear all Nodes
	c = l.Cursor()
	for l.RemoveAfter(c) != nil {
	}

	checkListPointers(t, l, []*Node[any]{})
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

func TestExtending(t *testing.T) {
	l1 := New[int]()
	l2 := New[int]()

	l1.pushBack(1)
	l1.pushBack(2)
	l1.pushBack(3)

	l2.pushBack(4)
	l2.pushBack(5)

	l3 := New[int]()
	l3.pushBackList(l1)
	checkList(t, l3, []int{1, 2, 3})

	l3.pushBackList(l2)
	checkList(t, l3, []int{1, 2, 3, 4, 5})

	l3 = New[int]()
	l3.pushFrontList(l2)
	checkList(t, l3, []int{4, 5})
	l3.pushFrontList(l1)
	checkList(t, l3, []int{1, 2, 3, 4, 5})

	checkList(t, l1, []int{1, 2, 3})
	checkList(t, l2, []int{4, 5})

	l3 = New[int]()
	l3.pushBackList(l1)
	checkList(t, l3, []int{1, 2, 3})
	l3.pushBackList(l3)
	checkList(t, l3, []int{1, 2, 3, 1, 2, 3})

	l3 = New[int]()
	l3.pushFrontList(l1)
	checkList(t, l3, []int{1, 2, 3})
	l3.pushFrontList(l3)
	checkList(t, l3, []int{1, 2, 3, 1, 2, 3})

	l3 = New[int]()
	l1.pushBackList(l3)
	checkList(t, l1, []int{1, 2, 3})
	l1.pushFrontList(l3)
	checkList(t, l1, []int{1, 2, 3})
}

func TestRemove(t *testing.T) {
	l := New[int]()
	e1 := l.pushBack(1)
	e2 := l.pushBack(2)
	checkListPointers(t, l, []*Node[int]{e1, e2})

	e := l.front()
	l.remove(e)
	checkListPointers(t, l, []*Node[int]{e2})
	//l.remove(e)
	//checkListPointers(t, l, []*Node[int]{e2})

}

func TestIssue6349(t *testing.T) {
	l := New[int]()
	l.pushBack(1)
	l.pushBack(2)

	e := l.front()
	l.remove(e)
	if e.Value != 1 {
		t.Errorf("e.value = %d, want 1", e.Value)
	}
	if e.next != nil {
		t.Errorf("e.Next() != nil")
	}
	if e.prev != nil {
		t.Errorf("e.Prev() != nil")
	}
}

func TestMove(t *testing.T) {
	l := New[int]()
	e1 := l.pushBack(1)
	e2 := l.pushBack(2)
	e3 := l.pushBack(3)
	e4 := l.pushBack(4)

	l.moveAfter(e3, e3)
	checkListPointers(t, l, []*Node[int]{e1, e2, e3, e4})
	l.moveBefore(e2, e2)
	checkListPointers(t, l, []*Node[int]{e1, e2, e3, e4})

	l.moveAfter(e3, e2)
	checkListPointers(t, l, []*Node[int]{e1, e2, e3, e4})
	l.moveBefore(e2, e3)
	checkListPointers(t, l, []*Node[int]{e1, e2, e3, e4})

	l.moveBefore(e2, e4)
	checkListPointers(t, l, []*Node[int]{e1, e3, e2, e4})
	e2, e3 = e3, e2

	l.moveBefore(e4, e1)
	checkListPointers(t, l, []*Node[int]{e4, e1, e2, e3})
	e1, e2, e3, e4 = e4, e1, e2, e3

	l.moveAfter(e4, e1)
	checkListPointers(t, l, []*Node[int]{e1, e4, e2, e3})
	e2, e3, e4 = e4, e2, e3

	l.moveAfter(e2, e3)
	checkListPointers(t, l, []*Node[int]{e1, e3, e2, e4})
}

// Test PushFront, PushBack, PushFrontList, PushBackList with uninitialized List
func TestZeroList(t *testing.T) {
	var l1 = new(List[int])
	l1.pushFront(1)
	checkList(t, l1, []int{1})

	var l2 = new(List[int])
	l2.pushBack(1)
	checkList(t, l2, []int{1})

	var l3 = new(List[int])
	l3.pushFrontList(l1)
	checkList(t, l3, []int{1})

	var l4 = new(List[int])
	l4.pushBackList(l2)
	checkList(t, l4, []int{1})
}
