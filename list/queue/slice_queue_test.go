package queue

import (
	"testing"
)

func TestSliceQueueEnqueue(t *testing.T) {
	queue := NewSliceQueue[int]()
	if actualValue := queue.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	if actualValue := queue.ToSlice(); actualValue[0] != 1 || actualValue[1] != 2 || actualValue[2] != 3 {
		t.Errorf("Got %v expected %v", actualValue, "[1,2,3]")
	}
	if actualValue := queue.Empty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := queue.Len(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := queue.Peek(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestSliceQueuePeek(t *testing.T) {
	queue := NewSliceQueue[int]()
	if actualValue, ok := queue.Peek(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)
	if actualValue, ok := queue.Peek(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestSliceQueueDequeue(t *testing.T) {
	queue := NewSliceQueue[int]()
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)
	queue.Dequeue()
	if actualValue, ok := queue.Peek(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := queue.Dequeue(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := queue.Dequeue(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := queue.Dequeue(); actualValue != 0 || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := queue.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := queue.ToSlice(); len(actualValue) != 0 {
		t.Errorf("Got %v expected %v", actualValue, "[]")
	}
}

func BenchmarkSliceQueueDequeue100(b *testing.B) {
	b.StopTimer()
	size := 100
	queue := NewSliceQueue[int]()
	for n := 0; n < size; n++ {
		queue.Enqueue(n)
	}
	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkSliceQueueDequeue1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	queue := NewSliceQueue[int]()
	for n := 0; n < size; n++ {
		queue.Enqueue(n)
	}
	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkSliceQueueDequeue10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	queue := NewSliceQueue[int]()
	for n := 0; n < size; n++ {
		queue.Enqueue(n)
	}
	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkSliceQueueDequeue100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	queue := NewSliceQueue[int]()
	for n := 0; n < size; n++ {
		queue.Enqueue(n)
	}
	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkSliceQueueEnqueue100(b *testing.B) {
	b.StopTimer()
	size := 100
	queue := NewSliceQueue[int]()
	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkSliceQueueEnqueue1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	queue := NewSliceQueue[int]()
	for n := 0; n < size; n++ {
		queue.Enqueue(n)
	}
	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkSliceQueueEnqueue10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	queue := NewSliceQueue[int]()
	for n := 0; n < size; n++ {
		queue.Enqueue(n)
	}
	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkSliceQueueEnqueue100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	queue := NewSliceQueue[int]()
	for n := 0; n < size; n++ {
		queue.Enqueue(n)
	}
	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}
