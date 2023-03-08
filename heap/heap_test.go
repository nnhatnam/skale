package heap

import (
	"container/heap"
	"math/rand"
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
		e = append(e, 0) // all elements are the same
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
		break
		h.verify(t, 0)
	}
}

func TestRemove1(t *testing.T) {
	h := NewOrdered[int]()
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	for i := 0; h.Len() > 0; i++ {
		x, success := h.Remove(0)
		if !success {
			t.Errorf("Remove(0) failed")
		}
		if x != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}
		h.verify(t, 0)
	}
}

func TestRemove2(t *testing.T) {
	N := 10

	var e []int
	for i := 0; i < N; i++ {
		e = append(e, i)
	}

	h := FromOrdered(e)

	h.verify(t, 0)

	m := make(map[int]bool)
	for h.Len() > 0 {
		x, _ := h.Remove((h.Len() - 1) / 2)
		m[x] = true
		h.verify(t, 0)
	}

	if len(m) != N {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}
	for i := 0; i < len(m); i++ {
		if !m[i] {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

//compare benchmark with container/heap

const benchmarkHeapSize = 1000000

//type IntHeap []int
//
//func (h IntHeap) Len() int           { return len(h) }
//func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
//func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
//
//func (h *IntHeap) Push(x any) {
//	// Push and Pop use pointer receivers because they modify the slice's length,
//	// not just its contents.
//	*h = append(*h, x.(int))
//}
//
//func (h *IntHeap) Pop() any {
//	old := *h
//	n := len(old)
//	x := old[n-1]
//	*h = old[0 : n-1]
//	return x
//}

type IntHeap []int

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Pop() (v any) {
	//*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]

	//var zero int

	n := len(*h) - 1

	v = (*h)[0]

	//delete
	(*h)[0] = (*h)[n]
	//(*h)[n] = zero
	*h = (*h)[:n]

	return v
}

func (h *IntHeap) Push(v any) {
	*h = append(*h, v.(int))
}

func BenchmarkDup(b *testing.B) {
	const n = 10000
	h := make(IntHeap, 0, n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			heap.Push(&h, 0) // all elements are the same
		}
		for h.Len() > 0 {
			heap.Pop(&h)
		}
	}
}

func BenchmarkDupSkale(b *testing.B) {
	const n = 10000

	h := NewOrdered[int](n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			h.Push(0) // all elements are the same
		}
		for h.Len() > 0 {
			h.Pop()
		}
	}
}

func BenchmarkHeapInsert(b *testing.B) {
	b.StopTimer()
	insertP := rand.Perm(benchmarkHeapSize)

	h := &IntHeap{}
	heap.Init(h)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, item := range insertP {
			heap.Push(h, item)
		}
	}
}

func BenchmarkHeapInsertSkale(b *testing.B) {
	b.StopTimer()
	insertP := rand.Perm(benchmarkHeapSize)

	h := NewOrdered[int]()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, item := range insertP {
			h.Push(item)
		}
	}
}

func BenchmarkHeapPop(b *testing.B) {
	const n = 10000

	b.StopTimer()
	insertP := rand.Perm(n)

	//h := &IntHeap{}
	h0 := make(IntHeap, 0, n)
	h := &h0
	heap.Init(h)

	j := 0
	for _, item := range insertP {
		heap.Push(h, item)
		j++
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {

		for h.Len() > 0 {
			heap.Pop(h)
		}
	}

}

func BenchmarkHeapPopSkale(b *testing.B) {
	const n = 10000

	b.StopTimer()
	insertP := rand.Perm(n)

	//h := &IntHeap{}
	h := NewOrdered[int](n)

	j := 0
	for _, item := range insertP {
		h.Push(item)
		j++
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {

		for h.Len() > 0 {
			h.Pop()
		}
	}
}
