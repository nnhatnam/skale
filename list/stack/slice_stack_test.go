package stack

import (
	"golang.org/x/exp/slices"
	"testing"
)

//test cases from https://github.com/emirpasic/gods

func TestStackInsertAfterClear(t *testing.T) {
	var stack = NewSStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	stack.Clear()

	if actualValue := stack.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := stack.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
	if actualValue, ok := stack.Top(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	stack.Push(1)

	if actualValue := stack.Empty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}

	if actualValue := stack.Len(); actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}

	if actualValue, ok := stack.Top(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestStackPush(t *testing.T) {
	var stack = NewSStack[int]()

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

func TestStackPeek(t *testing.T) {
	var stack = NewSStack[int]()
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

func TestStackPop(t *testing.T) {
	var stack = NewSStack[int]()
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

func benchmarkPush(b *testing.B, stack Stack[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			stack.Push(n)
		}
	}
}

func benchmarkPop(b *testing.B, stack Stack[int], size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			stack.Pop()
		}
	}
}

func BenchmarkSStackPop100(b *testing.B) {
	b.StopTimer()
	size := 100
	stack := NewSStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkSStackPop1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	stack := NewSStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkSStackPop10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	stack := NewSStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkSStackPop100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	stack := NewSStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkSStackPush100(b *testing.B) {
	b.StopTimer()
	size := 100
	stack := NewSStack[int]()
	b.StartTimer()
	benchmarkPush(b, stack, size)
}

func BenchmarkSStackPush1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	stack := NewSStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPush(b, stack, size)
}

func BenchmarkSStackPush10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	stack := NewSStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPush(b, stack, size)
}

func BenchmarkSStackPush100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	stack := NewSStack[int]()
	for n := 0; n < size; n++ {
		stack.Push(n)
	}
	b.StartTimer()
	benchmarkPush(b, stack, size)
}
