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

func TestCursor_Clone(t *testing.T) {
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

func TestCursor_Move(t *testing.T) {

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
