package linkedlist

import "testing"

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

			if c.Node() != nil {
				t.Errorf("Cursor.Node() = %v, want nil", c.Node())
			}
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
		t.Errorf("Cursor.current = %v, want %v", c.current, &l.root)
	}

	if c.Node() != nil {
		t.Errorf("Cursor.Node() = %v, want nil", c.Node())
	}

	if fc.current != l.root.next {
		t.Errorf("FrontCursor.current = %v, want %v", fc.current, l.root.next)
	}

	if fc.Node().Value != 1 {
		t.Errorf("FrontCursor.Node().Value = %v, want 1", fc.Node().Value)
	}

	if bc.current != l.root.prev {
		t.Errorf("BackCursor.current = %v, want %v", bc.current, l.root.prev)
	}

	if bc.Node().Value != 1 {
		t.Errorf("BackCursor.Node().Value = %v, want 1", bc.Node().Value)
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

	if c.Node() != nil {
		t.Errorf("Cursor.Node() = %v, want nil", c.Node())
	}

	if fc.current != l.root.next {
		t.Errorf("FrontCursor.current = %v, want %v", fc.current, l.root.next)
	}

	if fc.Node().Value != 1 {
		t.Errorf("FrontCursor.Node().Value = %v, want 1", fc.Node().Value)
	}

	if bc.current != l.root.prev {
		t.Errorf("BackCursor.current = %v, want %v", bc.current, l.root.prev)
	}

	if bc.Node().Value != 2 {
		t.Errorf("BackCursor.Node().Value = %v, want 2", bc.Node().Value)
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

	if c3.Node().Value != 1 {
		t.Errorf("Cursor.CloneNext().Node().Value = %v, want 1", c3.Node().Value)
	}

	if c4.current != c.current.prev || c.list != c4.list {
		t.Errorf("Cursor.ClonePrev() does not point to the same location as the original cursor")
	}

	if c4.Node().Value != 1 {
		t.Errorf("Cursor.ClonePrev().Node().Value = %v, want 1", c4.Node().Value)
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

	if c3.Node().Value != 1 {
		t.Errorf("Cursor.CloneNext().Node().Value = %v, want 1", c3.Node().Value)
	}

	if c4.current != c.current.prev || c.list != c4.list {
		t.Errorf("Cursor.ClonePrev() does not point to the same location as the original cursor")
	}

	if c4.Node().Value != 2 {
		t.Errorf("Cursor.ClonePrev().Node().Value = %v, want 2", c4.Node().Value)
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

	if c.Node().Value != 1 {
		t.Errorf("Cursor.MoveNext().Node().Value = %v, want 1", c.Node().Value)
	}

	c.MovePrev()
	checkCursor(t, l, []*Cursor[int]{c})

	c.MovePrev()

	if c.current != l.root.prev {
		t.Errorf("Cursor.MovePrev() = %v, want %v", c.current, l.root.prev)
	}

	if c.Node().Value != 1 {
		t.Errorf("Cursor.MovePrev().Node().Value = %v, want 1", c.Node().Value)
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

	if c.Node().Value != 1 {
		t.Errorf("Cursor.MoveToFront().Node().Value = %v, want 1", c.Node().Value)
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

	if c.Node().Value != 1 {
		t.Errorf("Cursor.MoveNext().Node().Value = %v, want 1", c.Node().Value)
	}

	c.MovePrev()
	checkCursor(t, l, []*Cursor[int]{c})

	c.MovePrev()

	if c.current != l.root.prev {
		t.Errorf("Cursor.MovePrev() = %v, want %v", c.current, l.root.prev)
	}

	if c.Node().Value != 2 {
		t.Errorf("Cursor.MovePrev().Node().Value = %v, want 2", c.Node().Value)
	}

	c.MoveNext()
	checkCursor(t, l, []*Cursor[int]{c})

	c.MoveToFront()

	if c.current != l.root.next {
		t.Errorf("Cursor.MoveToFront() = %v, want %v", c.current, l.root.next)
	}

	if c.Node().Value != 1 {
		t.Errorf("Cursor.MoveToFront().Node().Value = %v, want 1", c.Node().Value)
	}

	c.MoveToBack()

	if c.current != l.root.prev {
		t.Errorf("Cursor.MoveToBack() = %v, want %v", c.current, l.root.prev)
	}

	if c.Node().Value != 2 {
		t.Errorf("Cursor.MoveToBack().Node().Value = %v, want 2", c.Node().Value)
	}

}

func TestCursorWalk(t *testing.T) {

	l := New[string]()
	c := l.Cursor()

	// Test Walk on an empty list
	c.WalkAscending(func(c *Node[string]) bool {
		t.Errorf("Cursor.WalkAscending() called on an empty list")
		return true
	})

	c.WalkDescending(func(c *Node[string]) bool {
		t.Errorf("Cursor.WalkDescending() called on an empty list")
		return true
	})

	checkCursor(t, l, []*Cursor[string]{c})

	// Test Walk on a list with one element
	l.PushBack("1")

	c.WalkAscending(func(c *Node[string]) bool {
		if c.Value != "1" {
			t.Errorf("Cursor.WalkAscending() = %v, want 1", c.Value)
		}
		return true
	})

	c.WalkDescending(func(c *Node[string]) bool {
		if c.Value != "1" {
			t.Errorf("Cursor.WalkDescending() = %v, want 1", c.Value)
		}
		return true
	})

	checkCursor(t, l, []*Cursor[string]{c})

	if c.current != &l.root {
		t.Errorf("Cursor.WalkDescending() = %v, want %v", c.current, &l.root)
	}

	// Test Walk on a list with more than one element
	l.PushBack("2")

	c.WalkAscending(func(c *Node[string]) bool {
		if c.Value != "1" {
			t.Errorf("Cursor.WalkAscending() = %v, want 1", c.Value)
		}
		return false
	})

	if c.current != l.root.next {
		t.Errorf("Cursor.WalkAscending() = %v, want %v", c.current, l.root.next)
	}

	// The walker supposed to not move in this case
	c.WalkAscending(func(c *Node[string]) bool {
		if c.Value != "1" {
			t.Errorf("Cursor.WalkAscending() = %v, want 1", c.Value)
		}
		return false
	})

	// move to "2"
	c.WalkAscending(func(c *Node[string]) bool {
		if c.Value == "2" {
			return false // stop walking
		}
		return true
	})

	if c.current != l.root.prev {
		t.Errorf("Cursor.WalkAscending() = %v, want %v", c.current, l.root.prev)
	}

	// move to sentinel
	c.WalkAscending(func(c *Node[string]) bool {
		if c.Value != "2" {
			t.Errorf("Cursor.WalkAscending() = %v, want 2", c.Value)
		}
		return true
	})

	if c.current != &l.root {
		t.Errorf("Cursor.WalkAscending() = %v, want %v", c.current, &l.root)
	}

	c.WalkDescending(func(c *Node[string]) bool {
		if c.Value != "2" {
			t.Errorf("Cursor.WalkDescending() = %v, want 2", c.Value)
		}
		return false
	})

	if c.current != l.root.prev {
		t.Errorf("Cursor.WalkDescending() = %v, want %v", c.current, l.root.prev)
	}

	// The walker supposed to not move in this case
	c.WalkDescending(func(c *Node[string]) bool {
		if c.Value != "2" {
			t.Errorf("Cursor.WalkDescending() = %v, want 2", c.Value)
		}
		return false
	})

	// move to "1"
	c.WalkDescending(func(c *Node[string]) bool {
		if c.Value == "1" {
			return false // stop walking
		}
		return true
	})

	if c.current != l.root.next {
		t.Errorf("Cursor.WalkDescending() = %v, want %v", c.current, l.root.next)
	}

	// move to sentinel
	c.WalkDescending(func(c *Node[string]) bool {
		if c.Value != "1" {
			t.Errorf("Cursor.WalkDescending() = %v, want 1", c.Value)
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

	if c1.Node() != nil {
		t.Errorf("Cursor.Node() = %v, want nil", c1.Node())
	}

	if c2.Node() != nil {
		t.Errorf("Cursor.Node() = %v, want nil", c2.Node())
	}

	if c3.Node() != nil {
		t.Errorf("Cursor.Node() = %v, want nil", c3.Node())
	}

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

	if l.InsertAfter(5, c1) != nil {
		t.Errorf("Cursor.InsertAfter() = %v, want nil", l.InsertAfter(5, c1))
	}

}
