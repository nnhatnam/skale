package heap

import (
	"testing"
)

func (h *Heap[T]) verify(t *testing.T, i int) {
	t.Helper()
	n := len(h.e)
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h.less(h.e[j1], h.e[i]) {
			t.Errorf("heap invariant invalidated [%v] = %v > [%v] = %v", i, h.e[i], j1, h.e[j1])
			return
		}
		h.verify(t, j1)
	}
	if j2 < n {
		if h.less(h.e[j2], h.e[i]) {
			t.Errorf("heap invariant invalidated [%v] = %v > [%v] = %v", i, h.e[i], j1, h.e[j2])
			return
		}
		h.verify(t, j2)
	}
}

func TestInit0(t *testing.T) {
	var e []int
	for i := 20; i > 0; i-- {
		e = append(e, i) // all elements are the same
	}
	h := FromOrdered(e)

	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x, success := h.Pop()
		if !success {
			t.Errorf("Pop failed")
		}
		h.verify(t, 0)
		if x != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, x, 0)
		}
	}
}

func TestInit1(t *testing.T) {
	h := NewOrdered[int]()
	for i := 20; i > 0; i-- {
		h.Push(i) // all elements are different
	}
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x, success := h.Pop()
		if !success {
			t.Errorf("Pop failed")
		}
		h.verify(t, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func Test(t *testing.T) {
	h := NewOrdered[int]()
	h.verify(t, 0)

	for i := 20; i > 10; i-- {
		h.Push(i)
	}

	h.verify(t, 0)

	for i := 10; i > 0; i-- {
		h.Push(i)
		h.verify(t, 0)
	}

	for i := 1; h.Len() > 0; i++ {
		x, success := h.Pop()

		if !success {
			t.Errorf("Pop failed")
		}
		if i < 20 {
			h.Push(20 + i)
		}

		h.verify(t, 0)

		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
		//if x == 2 {
		//	fmt.Println("x is ", x)
		//	fmt.Println("h.e is ", h.e)
		//	return
		//}

	}
}

func TestRemove0(t *testing.T) {
	h := NewOrdered[int]()
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	for h.Len() > 0 {
		i := h.Len() - 1
		x, success := h.Remove(i)
		if !success {
			t.Errorf("Remove(%d) failed", i)
		}
		if x != i {
			t.Errorf("Remove(%d) got %d; want %d", i, x, i)
		}
		h.verify(t, 0)
	}
}

//
//func TestRemove1(t *testing.T) {
//	h := NewOrdered[int]()
//	for i := 0; i < 10; i++ {
//		h.Push(i)
//	}
//	h.verify(t, 0)
//
//	for i := 0; h.Len() > 0; i++ {
//		x, success := h.Remove(0)
//		if !success {
//			t.Errorf("Remove(0) failed")
//		}
//		if x != i {
//			t.Errorf("Remove(0) got %d; want %d", x, i)
//		}
//		h.verify(t, 0)
//	}
//}
//
//func TestRemove2(t *testing.T) {
//	N := 10
//
//	var e []int
//	for i := 0; i < N; i++ {
//		e = append(e, i)
//	}
//	h := FromOrdered(e)
//
//	h.verify(t, 0)
//
//	m := make(map[int]bool)
//	for h.Len() > 0 {
//		x, _ := h.Remove((h.Len() - 1) / 2)
//		m[x] = true
//		h.verify(t, 0)
//	}
//
//	if len(m) != N {
//		t.Errorf("len(m) = %d; want %d", len(m), N)
//	}
//	for i := 0; i < len(m); i++ {
//		if !m[i] {
//			t.Errorf("m[%d] doesn't exist", i)
//		}
//	}
//}
//
//func BenchmarkDup(b *testing.B) {
//	const n = 10000
//
//	h := NewOrdered[int](n)
//	for i := 0; i < b.N; i++ {
//		for j := 0; j < n; j++ {
//			h.Push(0) // all elements are the same
//		}
//		for h.Len() > 0 {
//			h.Pop()
//		}
//	}
//}

//
//func TestFix(t *testing.T) {
//	h := new(myHeap)
//	h.verify(t, 0)
//
//	for i := 200; i > 0; i -= 10 {
//		Push(h, i)
//	}
//	h.verify(t, 0)
//
//	if (*h)[0] != 10 {
//		t.Fatalf("Expected head to be 10, was %d", (*h)[0])
//	}
//	(*h)[0] = 210
//	Fix(h, 0)
//	h.verify(t, 0)
//
//	for i := 100; i > 0; i-- {
//		elem := rand.Intn(h.Len())
//		if i&1 == 0 {
//			(*h)[elem] *= 2
//		} else {
//			(*h)[elem] /= 2
//		}
//		Fix(h, elem)
//		h.verify(t, 0)
//	}
//}
