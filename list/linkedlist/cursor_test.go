package linkedlist

import (
	"testing"
)

func checkCursor[T any](t *testing.T, l *List[T], cursors []*Cursor[T]) {

	for _, c := range cursors {

		if c == nil {
			t.Errorf("Cursor = nil, want non-nil")
		}

		if c.current == nil {
			t.Errorf("Cursor.current = nil, want non-nil")
		}

		if c.list != l {
			t.Errorf("Cursor.list = %v, want %v", c.list, l)
		}

		if l.len == 0 {
			//check cursor at sentinel node
			if c.current != &l.root {
				t.Errorf("Cursor.current = %v, want %v", c.current, &l.root)
			}

			//if c.current != nil {
			//	t.Errorf("cursor.current = %v, want nil", c.current)
			//}
		}
	}
}

func TestNewCursor(t *testing.T) {
	//Test cursor on an empty list
	l := New[int]()
	c := l.Cursor()
	fc := l.FrontCursor()
	bc := l.BackCursor()

	checkCursor(t, l, []*Cursor[int]{c, fc, bc})

	//Test cursor on a non-empty list
	l.PushBack(1)
	c = l.Cursor()
	fc = l.FrontCursor()
	bc = l.BackCursor()

	checkCursor(t, l, []*Cursor[int]{c, fc, bc})

	if c.current != &l.root {
		t.Errorf("cursor.current = %v, want %v", c.current, &l.root)
	}

	if fc.current != l.root.next {
		t.Errorf("FrontCursor.current = %v, want %v", fc.current, l.root.next)
	}

	if fc.current.value != 1 {
		t.Errorf("FrontCursor.current.value = %v, want 1", fc.current.value)
	}

	if bc.current != l.root.prev {
		t.Errorf("BackCursor.current = %v, want %v", bc.current, l.root.prev)
	}

	if bc.current.value != 1 {
		t.Errorf("BackCursor.current.value = %v, want 1", bc.current.value)
	}

	// Test cursor on a list with more than one element
	l.PushBack(2)
	c = l.Cursor()
	fc = l.FrontCursor()
	bc = l.BackCursor()

	checkCursor(t, l, []*Cursor[int]{c, fc, bc})

	if c.current != &l.root {
		t.Errorf("Cursor.current = %v, want %v", c.current, &l.root)
	}

	if fc.current != l.root.next {
		t.Errorf("FrontCursor.current = %v, want %v", fc.current, l.root.next)
	}

	if fc.current.value != 1 {
		t.Errorf("FrontCursor.current.value = %v, want 1", fc.current.value)
	}

	if bc.current != l.root.prev {
		t.Errorf("BackCursor.current = %v, want %v", bc.current, l.root.prev)
	}

	if bc.current.value != 2 {
		t.Errorf("BackCursor.current.value = %v, want 2", bc.current.value)
	}

	//check cursor on list.Init()
	l.Init()
	c = l.Cursor()
	fc = l.FrontCursor()
	bc = l.BackCursor()

	checkCursor(t, l, []*Cursor[int]{c, fc, bc})

}

func TestCursorClone(t *testing.T) {
	l := New[any]()
	c := l.Cursor()
	c2 := c.Clone()
	c3 := c.CloneNext()
	c4 := c.ClonePrev()

	checkCursor(t, l, []*Cursor[any]{c2, c3, c4})

	if c2.current != c.current || c.list != c2.list {
		t.Errorf("Cursor.Clone() does not point to the same location as the original cursor")
	}

	if c3.current != c.current.next || c.list != c3.list {
		t.Errorf("Cursor.CloneNext() does not point to the same location as the original cursor")
	}

	if c4.current != c.current.prev || c.list != c4.list {
		t.Errorf("Cursor.ClonePrev() does not point to the same location as the original cursor")
	}

	//Test cursor on a non-empty list (1 element)
	l.PushBack(1)
	c = l.Cursor()
	c2 = c.Clone()
	c3 = c.CloneNext()
	c4 = c.ClonePrev()

	checkCursor(t, l, []*Cursor[any]{c2, c3, c4})

	if c2.current != c.current || c.list != c2.list {
		t.Errorf("Cursor.Clone() does not point to the same location as the original cursor")
	}

	if c3.current != c.current.next || c.list != c3.list {
		t.Errorf("Cursor.CloneNext() does not point to the same location as the original cursor")
	}

	if c3.current.value != 1 {
		t.Errorf("Cursor.CloneNext().current.value = %v, want 1", c3.current.value)
	}

	if c4.current != c.current.prev || c.list != c4.list {
		t.Errorf("Cursor.ClonePrev() does not point to the same location as the original cursor")
	}

	if c4.current.value != 1 {
		t.Errorf("Cursor.ClonePrev().current.value = %v, want 1", c4.current.value)
	}

	//Test cursor on a non-empty list (2 elements)
	l.PushBack(2)
	c = l.Cursor()
	c2 = c.Clone()
	c3 = c.CloneNext()
	c4 = c.ClonePrev()

	checkCursor(t, l, []*Cursor[any]{c2, c3, c4})

	if c2.current != c.current || c.list != c2.list {
		t.Errorf("Cursor.Clone() does not point to the same location as the original cursor")
	}

	if c3.current != c.current.next || c.list != c3.list {
		t.Errorf("Cursor.CloneNext() does not point to the same location as the original cursor")
	}

	if c3.current.value != 1 {
		t.Errorf("Cursor.CloneNext().current.value = %v, want 1", c3.current.value)
	}

	if c4.current != c.current.prev || c.list != c4.list {
		t.Errorf("Cursor.ClonePrev() does not point to the same location as the original cursor")
	}

	if c4.current.value != 2 {
		t.Errorf("Cursor.ClonePrev().current.value = %v, want 2", c4.current.value)
	}

	//Test cursor on list.Init()
	l.Init()
	c = l.Cursor()
	c2 = c.Clone()
	c3 = c.CloneNext()
	c4 = c.ClonePrev()

	checkCursor(t, l, []*Cursor[any]{c2, c3, c4})

}

func TestCursorMove(t *testing.T) {

	l := New[int]()
	c := l.Cursor()

	// Test Move on an empty list
	c.MoveToFront()
	checkCursor(t, l, []*Cursor[int]{c})
	c.MoveToBack()
	checkCursor(t, l, []*Cursor[int]{c})
	c.MoveNext()
	checkCursor(t, l, []*Cursor[int]{c})
	c.MoveToBack()
	checkCursor(t, l, []*Cursor[int]{c})

	// Test Move on a list with one element
	l.PushBack(1)
	c = l.Cursor()

	c.MoveNext()

	checkCursor(t, l, []*Cursor[int]{c})

	if c.current != l.root.next {
		t.Errorf("Cursor.MoveNext() = %v, want %v", c.current, l.root.next)
	}

	if c.current.value != 1 {
		t.Errorf("Cursor.MoveNext().current.value = %v, want 1", c.current.value)
	}

	c.MovePrev()
	checkCursor(t, l, []*Cursor[int]{c})

	c.MovePrev()

	if c.current != l.root.prev {
		t.Errorf("Cursor.MovePrev() = %v, want %v", c.current, l.root.prev)
	}

	if c.current.value != 1 {
		t.Errorf("Cursor.MovePrev().current.value = %v, want 1", c.current.value)
	}

	c.MoveNext()
	checkCursor(t, l, []*Cursor[int]{c})

	c.MoveToFront()
	checkCursor(t, l, []*Cursor[int]{c})

	if c.current == &l.root {
		t.Errorf("Cursor.MoveToFront() = %v, want %v", c.current, l.root.next)
	}

	if c.current != l.root.next {
		t.Errorf("Cursor.MoveToFront() = %v, want %v", c.current, l.root.next)
	}

	if c.current.value != 1 {
		t.Errorf("Cursor.MoveToFront().current.value = %v, want 1", c.current.value)
	}

	c.MoveToBack()

	if c.current == &l.root {
		t.Errorf("Cursor.MoveToBack() = %v, want %v", c.current, l.root.next)
	}

	if c.current != l.root.prev {
		t.Errorf("Cursor.MoveToBack() = %v, want %v", c.current, l.root.next)
	}

	// Test Move on a list with more than one element
	l.PushBack(2)
	c = l.Cursor()

	c.MoveNext()

	if c.current != l.root.next {
		t.Errorf("Cursor.MoveNext() = %v, want %v", c.current, l.root.next)
	}

	if c.current.value != 1 {
		t.Errorf("Cursor.MoveNext().current.value = %v, want 1", c.current.value)
	}

	c.MovePrev()
	checkCursor(t, l, []*Cursor[int]{c})

	c.MovePrev()

	if c.current != l.root.prev {
		t.Errorf("Cursor.MovePrev() = %v, want %v", c.current, l.root.prev)
	}

	if c.current.value != 2 {
		t.Errorf("Cursor.MovePrev().current.value = %v, want 2", c.current.value)
	}

	c.MoveNext()
	checkCursor(t, l, []*Cursor[int]{c})

	c.MoveToFront()

	if c.current != l.root.next {
		t.Errorf("Cursor.MoveToFront() = %v, want %v", c.current, l.root.next)
	}

	if c.current.value != 1 {
		t.Errorf("Cursor.MoveToFront().current.value = %v, want 1", c.current.value)
	}

	c.MoveToBack()

	if c.current != l.root.prev {
		t.Errorf("Cursor.MoveToBack() = %v, want %v", c.current, l.root.prev)
	}

	if c.current.value != 2 {
		t.Errorf("Cursor.MoveToBack().current.value = %v, want 2", c.current.value)
	}

}

func TestCursorWalk(t *testing.T) {

	l := New[string]()
	c := l.Cursor()

	// Test Walk on an empty list
	c.WalkAscending(func(c string) bool {
		t.Errorf("Cursor.WalkAscending() called on an empty list")
		return true
	})

	c.WalkDescending(func(c string) bool {
		t.Errorf("Cursor.WalkDescending() called on an empty list")
		return true
	})

	checkCursor(t, l, []*Cursor[string]{c})

	// Test Walk on a list with one element
	l.PushBack("1")

	c.WalkAscending(func(v string) bool {
		if v != "1" {
			t.Errorf("Cursor.WalkAscending() = %v, want 1", v)
		}
		return true
	})

	c.WalkDescending(func(v string) bool {

		if v != "1" {
			t.Errorf("Cursor.WalkDescending() = %v, want 1", v)
		}
		return true
	})

	checkCursor(t, l, []*Cursor[string]{c})

	if c.current != &l.root {
		t.Errorf("Cursor.WalkDescending() = %v, want %v", c.current, &l.root)
	}

	// Test Walk on a list with more than one element
	l.PushBack("2")

	c.WalkAscending(func(v string) bool {
		if v != "1" {
			t.Errorf("Cursor.WalkAscending() = %v, want 1", v)
		}
		return false
	})

	if c.current != l.root.next {
		t.Errorf("Cursor.WalkAscending() = %v, want %v", c.current, l.root.next)
	}

	// The walker supposed to not move in this case
	c.WalkAscending(func(v string) bool {
		if v != "1" {
			t.Errorf("Cursor.WalkAscending() = %v, want 1", v)
		}
		return false
	})

	// move to "2"
	c.WalkAscending(func(v string) bool {
		if v == "2" {
			return false // stop walking
		}
		return true
	})

	if c.current != l.root.prev {
		t.Errorf("Cursor.WalkAscending() = %v, want %v", c.current, l.root.prev)
	}

	// move to sentinel
	c.WalkAscending(func(v string) bool {
		if v != "2" {
			t.Errorf("Cursor.WalkAscending() = %v, want 2", v)
		}
		return true
	})

	if c.current != &l.root {
		t.Errorf("Cursor.WalkAscending() = %v, want %v", c.current, &l.root)
	}

	c.WalkDescending(func(v string) bool {
		if v != "2" {
			t.Errorf("Cursor.WalkDescending() = %v, want 2", v)
		}
		return false
	})

	if c.current != l.root.prev {
		t.Errorf("Cursor.WalkDescending() = %v, want %v", c.current, l.root.prev)
	}

	// The walker supposed to not move in this case
	c.WalkDescending(func(v string) bool {
		if v != "2" {
			t.Errorf("Cursor.WalkDescending() = %v, want 2", v)
		}
		return false
	})

	// move to "1"
	c.WalkDescending(func(v string) bool {
		if v == "1" {
			return false // stop walking
		}
		return true
	})

	if c.current != l.root.next {
		t.Errorf("Cursor.WalkDescending() = %v, want %v", c.current, l.root.next)
	}

	// move to sentinel
	c.WalkDescending(func(v string) bool {
		if v != "1" {
			t.Errorf("Cursor.WalkDescending() = %v, want 1", v)
		}
		return true
	})

	if c.current != &l.root {
		t.Errorf("Cursor.WalkDescending() = %v, want %v", c.current, &l.root)
	}

}

func TestCursorClose(t *testing.T) {

	l := New[string]()
	c := l.Cursor()

	// Test Close on an empty list
	c.Close()

	if c.current != nil || c.list != nil {
		t.Errorf("Cursor.Close() = %v, %v, want nil, nil", c.current, c.list)
	}

	// Test Close on a list with one element
	l.PushBack("1")

	c = l.FrontCursor()
	c.Close()

	if c.current != nil || c.list != nil {
		t.Errorf("Cursor.Close() = %v, %v, want nil, nil", c.current, c.list)
	}

}

func TestCursorInvalid(t *testing.T) {

	l := New[int]()

	l.PushBack(1)
	c1 := l.Cursor()
	c2 := l.FrontCursor()
	c3 := l.BackCursor()

	l.PopFront()

	if c1.current != &l.root {
		t.Errorf("c1.current = %v, want %v", c1.current, &l.root)
	}

	//if c2.current != &l.root {
	//	t.Errorf("c2.current = %v, want %v", c2.current, &l.root)
	//}
	//
	//if c3.current != &l.root {
	//	t.Errorf("c3.current = %v, want %v", c3.current, &l.root)
	//}

	if !c1.IsValid() {
		t.Errorf("Cursor.IsValid() = %v, want true", c1.IsValid())
	}

	if c2.IsValid() {
		t.Errorf("Cursor.IsValid() = %v, want false", c2.IsValid())
	}

	if c3.IsValid() {
		t.Errorf("Cursor.IsValid() = %v, want false", c3.IsValid())
	}

	//clone invalid cursor
	invalid1 := c2.Clone()
	invalid2 := c2.CloneNext()
	invalid3 := c2.ClonePrev()

	if invalid1 != nil || invalid2 != nil || invalid3 != nil {
		t.Errorf("Cursor.Clone() = %v, %v, %v, want nil, nil, nil", invalid1, invalid2, invalid3)
	}

	l.PushBack(2)
	l.PushBack(3)

	c1 = l.FrontCursor()
	c2 = l.FrontCursor()

	l.RemoveAt(c1) // <-- this should invalidate c2

	if c2.IsValid() {
		t.Errorf("Cursor.IsValid() = %v, want false", c2.IsValid())
	}

	if !c1.IsValid() {
		t.Errorf("Cursor.IsValid() = %v, want true", c1.IsValid())
	}

	l.PushBack(4)
	c3 = l.BackCursor()
	l.RemoveBefore(c3) // <-- this should invalidate c1

	if c1.IsValid() {
		t.Errorf("Cursor.IsValid() = %v, want false", c1.IsValid())
	}

	if l.InsertAfter(5, c1) != false {
		t.Errorf("Cursor.InsertAfter() = %v, want false", l.InsertAfter(5, c1))
	}

}
