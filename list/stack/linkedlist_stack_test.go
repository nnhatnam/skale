package stack

import (
	"golang.org/x/exp/slices"
	"testing"
)

func TestLLStackPush(t *testing.T) {
	stack := NewLLStack[int]()
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

func TestLLStackPeek(t *testing.T) {
	stack := NewLLStack[int]()
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

func TestLLStackPop(t *testing.T) {
	stack := NewLLStack[int]()
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

func BenchmarkLLStackPop100(b *testing.B) {
	b.StopTimer()
	size := 100
	stack := NewLLStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkLLStackPop1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	stack := NewLLStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkLLStackPop10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	stack := NewLLStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkLLStackPop100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	stack := NewLLStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkLLStackPush100(b *testing.B) {
	b.StopTimer()
	size := 100
	stack := NewLLStack[int]()
	b.StartTimer()
	benchmarkPush(b, stack, size)
}

func BenchmarkLLStackPush1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	stack := NewLLStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPush(b, stack, size)
}

func BenchmarkLLStackPush10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	stack := NewLLStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPush(b, stack, size)
}

func BenchmarkLLStackPush100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	stack := NewLLStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPush(b, stack, size)
}
