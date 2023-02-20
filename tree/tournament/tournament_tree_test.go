package tournament

import (
	"testing"
)

func TestNew(t *testing.T) {
	tt := NewWithOrdered[int]()
	if tt == nil {
		t.Errorf("New() returns nil")
	}

	tt.insert(10)
	tt.print()
	tt.insert(9)
	tt.print()
	tt.insert(8)
	tt.print()
	tt.insert(7)
	tt.print()
	tt.insert(6)
	tt.print()
	tt.insert(5)
	tt.print()
	tt.insert(4)
	tt.print()
	tt.insert(3)
	tt.print()
	tt.insert(2)
	tt.print()
	tt.insert(1)
	tt.print()
	tt.insert(0)
	tt.print()
	tt.insert(-1)
	tt.print()
	tt.insert(-2)
	tt.print()
	tt.insert(-3)
	tt.print()
	tt.insert(-4)
	tt.print()
}
