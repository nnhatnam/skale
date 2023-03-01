package tournament

import (
	"testing"
)

//func TestNew(t *testing.T) {
//	tt := NewWithOrdered[int]()
//	if tt == nil {
//		t.Errorf("New() returns nil")
//	}
//
//	tt.Insert(10)
//	tt.print()
//	tt.Insert(9)
//	tt.print()
//	tt.Insert(8)
//	tt.print()
//	tt.Insert(7)
//	tt.print()
//	tt.Insert(6)
//	tt.print()
//	tt.Insert(5)
//	tt.print()
//	tt.Insert(4)
//	tt.print()
//	tt.Insert(3)
//	tt.print()
//	tt.Insert(2)
//	tt.print()
//	tt.Insert(1)
//	tt.print()
//	tt.Insert(0)
//	tt.print()
//	tt.Insert(-1)
//	tt.print()
//	tt.Insert(-2)
//	tt.print()
//	tt.Insert(-3)
//	tt.print()
//	tt.Insert(-4)
//	tt.print()
//}

func TestDelete(t *testing.T) {

	tt := NewWithOrdered[int]()
	if tt == nil {
		t.Errorf("New() returns nil")
	}

	tt.Insert(1)
	tt.Insert(2)
	tt.Insert(3)

	tt.print()

}
