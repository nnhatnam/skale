package skiplist

import (
	"github.com/nnhatnam/skale"
)

type item[K skale.Ordered, V any] struct {
	key   K
	value V
}

type SkipMap[K skale.Ordered, V any] struct {
	*SkipList[item[K, V]]
}

func NewSkipMap[K skale.Ordered, V any](maxLevel int, p float64) *SkipMap[K, V] {
	l := New[item[K, V]](maxLevel, p, func(a, b item[K, V]) bool {
		return a.key < b.key
	})
	return &SkipMap[K, V]{l}
}

func (l *SkipMap[K, V]) ReplaceOrInsert(key K, value V) (_ V, _ bool) {

	result, ok := l.SkipList.ReplaceOrInsert(item[K, V]{key, value})

	if ok {
		return result.value, ok
	}

	return
}

func (l *SkipMap[K, V]) InsertNoReplace(key K, value V) {

	l.SkipList.InsertNoReplace(item[K, V]{key, value})
}

func (l *SkipMap[K, V]) Get(key K) (_ V, _ bool) {

	result, ok := l.SkipList.Get(item[K, V]{key: key})

	if ok {
		return result.value, ok
	}

	return
}

// Delete removes an item equal to the passed in item from the list, returning it. If no such item exists, returns (zeroValue, false).
func (l *SkipMap[K, V]) Delete(key K) (_ V, _ bool) {

	result, ok := l.SkipList.Delete(item[K, V]{key: key})

	if ok {
		return result.value, ok
	}

	return
}

// DeleteMin deletes the minimum value in the list and returns it. If no such value exists, returns (zeroValue, false).
func (l *SkipMap[K, V]) DeleteMin() (_ K, _ V, _ bool) {

	result, ok := l.SkipList.DeleteMin()

	if ok {
		return result.key, result.value, ok
	}

	return
}

// DeleteMax deletes the maximum value in the list and returns it. If no such value exists, returns (zeroValue, false).
func (l *SkipMap[K, V]) DeleteMax() (_ K, _ V, _ bool) {

	result, ok := l.SkipList.DeleteMax()

	if ok {
		return result.key, result.value, ok
	}

	return
}

// Len returns the number of items currently in the list.
func (l *SkipMap[K, V]) Len() int {
	return l.SkipList.Len()
}

// Has returns true if the given value is in the list
func (l *SkipMap[K, V]) Has(key K) bool {
	return l.SkipList.Has(item[K, V]{key: key})
}

// Ascend calls the iterator for every value in the list within the range [first, last], until iterator returns false.
func (l *SkipMap[K, V]) Ascend(iter MapIterator[K, V]) {

	l.SkipList.Ascend(func(item item[K, V]) bool {
		return iter(item.key, item.value)
	})

}

// AscendGreaterOrEqual calls the iterator for every value in the list within the range [pivot , last], until iterator returns false.
func (l *SkipMap[K, V]) AscendGreaterOrEqual(pivot K, iter MapIterator[K, V]) {

	l.SkipList.AscendGreaterOrEqual(item[K, V]{key: pivot}, func(item item[K, V]) bool {
		return iter(item.key, item.value)
	})

}

// AscendLessThan calls the iterator for every value in the list within the range [first, pivot), until iterator returns false.
func (l *SkipMap[K, V]) AscendLessThan(pivot K, iter MapIterator[K, V]) {

	l.SkipList.AscendLessThan(item[K, V]{key: pivot}, func(item item[K, V]) bool {
		return iter(item.key, item.value)
	})

}

// AscendRange calls the iterator for every value in the list within the range [greaterOrEqual, lessThan) , until iterator returns false.
func (l *SkipMap[K, V]) AscendRange(greaterOrEqual, lessThan K, iter MapIterator[K, V]) {

	l.SkipList.AscendRange(item[K, V]{key: greaterOrEqual}, item[K, V]{key: lessThan}, func(item item[K, V]) bool {
		return iter(item.key, item.value)
	})

}

// Descend calls the iterator for every value in the list within the range [last, first], until iterator returns false.
func (l *SkipMap[K, V]) Descend(iter MapIterator[K, V]) {

	l.SkipList.Descend(func(item item[K, V]) bool {
		return iter(item.key, item.value)
	})

}

// DescendGreaterThan calls the iterator for every value in the list within the range [last, pivot), until iterator returns false.
func (l *SkipMap[K, V]) DescendGreaterThan(pivot K, iter MapIterator[K, V]) {

	l.SkipList.DescendGreaterThan(item[K, V]{key: pivot}, func(item item[K, V]) bool {
		return iter(item.key, item.value)
	})

}

// DescendLessOrEqual calls the iterator for every value in the list within the range [pivot, first], until iterator returns false.
func (l *SkipMap[K, V]) DescendLessOrEqual(pivot K, iter MapIterator[K, V]) {

	l.SkipList.DescendLessOrEqual(item[K, V]{key: pivot}, func(item item[K, V]) bool {
		return iter(item.key, item.value)
	})

}

// DescendRange calls the iterator for every value in the list within the range [lessOrEqual, greaterThan), until iterator returns false.
func (l *SkipMap[K, V]) DescendRange(lessOrEqual, greaterThan K, iter MapIterator[K, V]) {

	l.SkipList.DescendRange(item[K, V]{key: lessOrEqual}, item[K, V]{key: greaterThan}, func(item item[K, V]) bool {
		return iter(item.key, item.value)
	})

}

// Max returns the largest item in the list, or (zeroValue, false) if the list is empty.
func (l *SkipMap[K, V]) Max() (_ K, _ V, _ bool) {

	result, ok := l.SkipList.Max()

	if ok {
		return result.key, result.value, ok
	}

	return
}

// Min returns the smallest item in the list, or (zeroValue, false) if the list is empty.
func (l *SkipMap[K, V]) Min() (_ K, _ V, _ bool) {

	result, ok := l.SkipList.Min()

	if ok {
		return result.key, result.value, ok
	}

	return

}
