package skiplist

import (
	"reflect"
	"testing"
)

func TestAscendGreaterOrEqual(t *testing.T) {
	l := NewOrdered[int](8, 0.5)
	l.InsertNoReplace(4)
	l.InsertNoReplace(6)
	l.InsertNoReplace(1)
	l.InsertNoReplace(3)
	var ary []int
	l.AscendGreaterOrEqual(-1, func(i int) bool {
		ary = append(ary, i)
		return true
	})
	expected := []int{1, 3, 4, 6}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	l.AscendGreaterOrEqual(3, func(i int) bool {
		ary = append(ary, i)
		return true
	})
	expected = []int{3, 4, 6}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	l.AscendGreaterOrEqual(2, func(i int) bool {
		ary = append(ary, i)
		return true
	})
	expected = []int{3, 4, 6}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
}

func TestDescendLessOrEqual(t *testing.T) {
	l := NewOrdered[int](8, 0.5)
	l.InsertNoReplace(4)
	l.InsertNoReplace(6)
	l.InsertNoReplace(1)
	l.InsertNoReplace(3)
	var ary []int
	l.DescendLessOrEqual(10, func(i int) bool {
		ary = append(ary, i)
		return true
	})
	expected := []int{6, 4, 3, 1}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil

	l.DescendLessOrEqual(4, func(i int) bool {
		ary = append(ary, i)
		return true
	})
	expected = []int{4, 3, 1}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	l.DescendLessOrEqual(5, func(i int) bool {
		ary = append(ary, i)
		return true
	})
	expected = []int{4, 3, 1}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
}

func TestAscendLessThan(t *testing.T) {
	l := NewOrdered[int](8, 0.5)
	l.InsertNoReplace(4)
	l.InsertNoReplace(6)
	l.InsertNoReplace(1)
	l.InsertNoReplace(3)
	var ary []int
	l.AscendLessThan(10, func(i int) bool {
		ary = append(ary, i)
		return true
	})
	expected := []int{1, 3, 4, 6}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	l.AscendLessThan(4, func(i int) bool {
		ary = append(ary, i)
		return true
	})
	expected = []int{1, 3}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	l.AscendLessThan(5, func(i int) bool {
		ary = append(ary, i)
		return true
	})
	expected = []int{1, 3, 4}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
}
