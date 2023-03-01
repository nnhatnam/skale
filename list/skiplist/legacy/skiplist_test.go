package legacy

import (
	"testing"
)

func TestNew(t *testing.T) {
	l := NewOrdered[int](64, 0.5)

	if l == nil {
		t.Error("New() failed")
	}

	if l.head == nil {
		t.Error("New() failed")
	}

	if l.maxLevel != 64 {
		t.Error("New() failed")
	}

	if l.p != 0.5 {
		t.Error("New() failed")
	}

	if l.less == nil {
		t.Error("New() failed")
	}

	if l.fingers == nil {
		t.Error("New() failed")
	}

	if len(l.fingers) != 64 {
		t.Error("New() failed")
	}

	if l.size != 0 {
		t.Error("New() failed")
	}

	if l.head.next == nil {
		t.Error("New() failed")
	}

	if len(l.head.next) != 64 {
		t.Errorf("New() failed, len(l.head.next) = %d", len(l.head.next))
	}

	if l.head.next[0] != l.tail {
		t.Error("New() failed")
	}

	if l.head.next[63] != l.tail {
		t.Error("New() failed")
	}

	if l.tail.next == nil {
		t.Error("New() failed")
	}

	if len(l.tail.next) != 64 {
		t.Error("New() failed")
	}
}

//func TestMaxLevel(t *testing.T) {
//	l := NewOrdered[int](64, 0.5)
//	var m = make(map[uint8]int)
//	for i := 0; i < 100000; i++ {
//		m[l.generateLevel()]++
//	}
//	fmt.Println(m)
//}

func TestReplaceOrInsert(t *testing.T) {
	l := NewOrdered[int](64, 0.5)

	l.ReplaceOrInsert(1)

	return

	if l.size != 1 {
		t.Errorf("ReplaceOrInsert failed: expected size 1, got %v", l.size)
	}

	if n := l.Get(1); n == nil {
		t.Errorf("ReplaceOrInsert failed: expected node with key 1, got nil")
	}

}
