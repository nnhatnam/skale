package example

import (
	"fmt"
	"github.com/nnhatnam/skale/list/linkedlist"
	"testing"
)

// https://leetcode.com/problems/add-two-numbers/
// You are given two non-empty linked lists representing two non-negative integers. The digits are stored in reverse order, and each of their nodes contains a single digit. Add the two numbers and return the sum as a linked linkedlist.
// You may assume the two numbers do not contain any leading zero, except the number 0 itself.
func TestAddTwoNumbers(t *testing.T) {

	//Example 1:
	l1 := linkedlist.From(2, 4, 3)
	l2 := linkedlist.From(5, 6, 4)

	//Output: 7 -> 0 -> 8 -> nil (342 + 465 = 807)
	l3 := addTwoNumbers(l1, l2)
	sum := 0
	l3.RTraverse(func(v any) {
		sum = sum*10 + v.(int)
	})
	fmt.Println("aa")
	if sum != 807 {
		t.Errorf("Example 1, expected %d, got %d", sum, 807)
	}

}

// addTwoNumbers adds two numbers represented by linked linkedlist. l1 and l2 are linked linkedlist representing two non-negative integers. The digits are stored in reverse order, and each of their nodes contains a single digit. Add the two numbers and return the sum as a linked linkedlist.
// l1 and l2 must not be nil.
func addTwoNumbers(l1 *linkedlist.List, l2 *linkedlist.List) *linkedlist.List {
	l3 := linkedlist.New()
	if l1.Len() == 0 && l2.Len() == 0 {
		return l3
	} else if l1.Len() == 0 {
		return l2
	} else if l2.Len() == 0 {
		return l1
	}

	it1 := l1.Begin()
	it2 := l2.Begin()

	diff := l1.Len() - l2.Len()

	//l1.Len() > l2.Len() or l1.Len() < l2.Len() => it1 and it2 is in different position. Move one of them to the same position.
	if diff > 0 {
		l3.PushBack(it1.Value())

		for i := 1; i < diff; i++ {
			l3.PushBack(it1.Next())
		}

		it1.Next()
	} else if diff < 0 {
		l3.PushBack(it2.Value())

		for i := 1; i < -diff; i++ {
			l3.PushBack(it2.Next())
		}

		it2.Next()
	}

	//it1 and it2 has the same position
	sum := it1.Value().(int) + it2.Value().(int)

	l3.PushBack(sum % 10)

	for {
		sum = it1.Value().(int) + it2.Value().(int) + sum/10
		l3.PushBack(sum % 10)
		it1.Next()
		it2.Next()

		if !it1.HasNext() {
			if sum/10 > 0 {
				l3.PushBack(sum / 10)
			}
			break
		}
	}

	return l3

}
