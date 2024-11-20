package stack

import (
	"golang.org/x/exp/slices"
	"testing"
)

func TestStackLPush(t *testing.T) {
	stack := NewStackL[int]()
	if actualValue := stack.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if actualValue := stack.ToSlice(); slices.Equal(actualValue, []int{3, 2, 1}) == false {
		t.Errorf("Got %v expected %v", actualValue, []int{3, 2, 1})
	}
	if actualValue := stack.Empty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := stack.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := stack.Top(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestStackLTop(t *testing.T) {
	stack := NewStackL[int]()
	if actualValue, ok := stack.Top(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	if actualValue, ok := stack.Top(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}

func TestStackBottom(t *testing.T) {
	stack := NewStackL[int]()
	if actualValue, ok := stack.Bottom(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	if actualValue, ok := stack.Bottom(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestStackLPop(t *testing.T) {
	stack := NewStackL[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	stack.Pop()
	if actualValue, ok := stack.Top(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := stack.Pop(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := stack.Pop(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
	if actualValue, ok := stack.Pop(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := stack.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := stack.ToSlice(); len(actualValue) != 0 {
		t.Errorf("Got %v expected %v", actualValue, "[]")
	}
}

func BenchmarkStackLPop100(b *testing.B) {
	b.StopTimer()
	size := 100
	stack := NewStackL[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkStackLPop1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	stack := NewStackL[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkStackLPop10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	stack := NewStackL[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkStackLPop100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	stack := NewStackL[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkStackLPush100(b *testing.B) {
	b.StopTimer()
	size := 100
	stack := NewStackL[int]()
	b.StartTimer()
	benchmarkPush(b, stack, size)
}

func BenchmarkStackLPush1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	stack := NewStackL[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPush(b, stack, size)
}

func BenchmarkStackLPush10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	stack := NewStackL[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPush(b, stack, size)
}

func BenchmarkStackLPush100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	stack := NewStackL[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPush(b, stack, size)
}
