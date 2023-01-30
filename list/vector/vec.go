package vector

import "errors"

var (
	ErrIndexOutOfRange   = errors.New("index out of range")
	ErrIndexReachMaxSize = errors.New("index reach max size")
	MaxArraySize         = 1<<31 - 1
	//
	//maxUint = ^uint(0)
	//maxInt  = int(maxUint >> 1)
)

// GrowFUnc will be call when the capacity is not enough. minCapacity is the min capacity that the vector need.
var GrowFunc = func(minCapacity int) int {
	newCap := minCapacity + minCapacity>>1

	//over flow
	if newCap < minCapacity || newCap > MaxArraySize {
		return MaxArraySize
	}

	return newCap
}

// https://github.com/golang/go/wiki/SliceTricks
type Vector[V any] struct {
	elements []V
}

func NewVector[V any](size int) *Vector[V] {
	return &Vector[V]{elements: make([]V, 0, size)}
}

func From[V any](values ...V) *Vector[V] {
	return &Vector[V]{elements: values}
}

func (v *Vector[V]) Len() int {
	return len(v.elements)
}

func (v *Vector[V]) Cap() int {
	return cap(v.elements)
}

func (v *Vector[V]) At(i int) V {
	if i < 0 || i >= len(v.elements) {
		panic(ErrIndexOutOfRange.Error())
	}
	return v.elements[i]
}

func (v *Vector[V]) Set(i int, value V) {
	v.elements[i] = value
}

func (v *Vector[V]) resize(newCap int) {
	newElements := make([]V, len(v.elements), newCap)
	copy(newElements, v.elements)
	v.elements = newElements
}

func (v *Vector[V]) ensureEnoughInternalCap(extraLen int) {

	//extraLength is 0, no need to grow the vector capacity
	if extraLen == 0 {
		return
	}

	minCapacity := len(v.elements) + extraLen

	//overflow happens, we can't expand the vector anymore
	if minCapacity < extraLen || int(minCapacity) > MaxArraySize {
		panic(ErrIndexReachMaxSize.Error())
		return
	}

	if minCapacity > cap(v.elements) {
		v.resize(GrowFunc(minCapacity))
	}

}

func (v *Vector[V]) Append(value ...V) {
	v.ensureEnoughInternalCap(len(value))
	v.elements = append(v.elements, value...)
}

func (v *Vector[V]) Cut(i, j int) {
	copy(v.elements[i:], v.elements[j:])
	for k, n := len(v.elements)-j+i, len(v.elements); k < n; k++ {
		v.elements[k] = *new(V)
	}
	v.elements = v.elements[:len(v.elements)-j+i]
}

func (v *Vector[V]) Delete(i int) {
	copy(v.elements[i:], v.elements[i+1:])
	v.elements[len(v.elements)-1] = interface{}(nil)
	v.elements = v.elements[:len(v.elements)-1]
}

func (v *Vector[V]) DeleteUnordered(i int) {
	v.elements[i] = v.elements[len(v.elements)-1]
	v.elements[len(v.elements)-1] = interface{}(nil)
	v.elements = v.elements[:len(v.elements)-1]
}

func (v *Vector[V]) Expand(i int, value ...V) {
	v.ensureEnoughInternalCap(len(value))
	v.elements = append(v.elements[:i], append(value, v.elements[i:]...)...)
}

func (v *Vector[V]) Insert(i int, value ...V) {
	v.ensureEnoughInternalCap(len(value))
	v.elements = append(v.elements[:i], append(value, v.elements[i:]...)...)
}

func (v *Vector[V]) Clear() {
	v.elements = v.elements[:0]
}
