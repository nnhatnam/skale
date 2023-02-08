package datrie

import (
	"golang.org/x/exp/slices"
	"testing"
)

type testingArcDomain []rune

func (d testingArcDomain) Code(s rune) int {
	for i, c := range d {
		if c == s {
			return i
		}
	}

	panic("invalid arc label")
}

func (d testingArcDomain) Label(i int) rune {
	return d[i]
}

func (d testingArcDomain) StopElement() rune {
	return d[1]
}

var alphabetRuneTesting testingArcDomain = []rune("_#abcdefghijklmnopqrstuvwxyz1234567890")

//
//func TestInitialState(t *testing.T) {
//
//	trie := New[rune, int](alphabetRuneTesting)
//
//	if trie.Len() != 0 {
//		t.Errorf("expecting len 0")
//	}
//
//	if trie.Contain([]rune("test")) {
//		t.Errorf("not expecting to find key=test")
//	}
//
//	expectedBase := []int{0, 1}
//
//	if !slices.Equal(trie.dArray._base(), expectedBase) {
//		t.Errorf("expecting base=%v, got %v", expectedBase, trie.dArray._base())
//	}
//
//	expectedCheck := []int{0, 0}
//
//	if !slices.Equal(trie.dArray._check(), expectedCheck) {
//		t.Errorf("expecting check=%v, got %v", expectedCheck, trie.dArray._check())
//	}
//
//	if len(trie.tail) != 2 {
//		t.Errorf("expecting tail len 2, got %v", len(trie.tail))
//	}
//
//	if trie.pos != 1 {
//		t.Errorf("expecting pos 1")
//	}
//
//}

func TestInsert(t *testing.T) {
	trie := New[rune, int](alphabetRuneTesting)
	trie.Insert([]rune("bachelor"), 100)

	if trie.Len() != 1 {
		t.Errorf("expecting len 1, got %v", trie.Len())
	}

	//Case 1: the insertion of the new word when the double-array is empty
	if !trie.Contain([]rune("bachelor")) {
		t.Errorf("expecting to find key=bachelor")
	}

	if value, found := trie.Get([]rune("bachelor")); value != 100 || !found {
		t.Errorf("expecting value 100, got %v", value)
	}

	expectedBase := []int{0, 1, 0, 0, -1}

	if !slices.Equal(trie.dArray._base(), expectedBase) {
		t.Errorf("expecting base=%v, got %v", expectedBase, trie.dArray._base())
	}

	expectedCheck := []int{0, 0, 0, 0, 1}

	if !slices.Equal(trie.dArray._check(), expectedCheck) {
		t.Errorf("expecting check=%v, got %v", expectedCheck, trie.dArray._check())
	}

	var zero rune
	expectedTail := []rune{zero, 'a', 'c', 'h', 'e', 'l', 'o', 'r', '#'}

	if !slices.Equal(trie.tail, expectedTail) {
		t.Errorf("expecting tail=%v, got %v", expectedTail, trie.tail)
	}

	if trie.pos != 9 {
		t.Errorf("expecting pos 9 , got %v", trie.pos)
	}

	trie.Insert([]rune("jar"), 101)

	if trie.Len() != 2 {
		t.Errorf("expecting len 2, got %v", trie.Len())
	}

	if !trie.Contain([]rune("jar")) {
		t.Errorf("expecting to find key=jar")
	}

	//Case 2: insertion, when the new word is inserted without collisions
	if value, found := trie.Get([]rune("jar")); value != 101 || !found {
		t.Errorf("expecting value 101, got %v", value)
	}

	expectedBase = []int{0, 1, 0, 0, -1, 0, 0, 0, 0, 0, 0, 0, -9}

	if !slices.Equal(trie.dArray._base(), expectedBase) {
		t.Errorf("expecting base=%v, got %v", expectedBase, trie.dArray._base())
	}

	expectedCheck = []int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1}

	if !slices.Equal(trie.dArray._check(), expectedCheck) {
		t.Errorf("expecting check=%v, got %v", expectedCheck, trie.dArray._check())
	}

	expectedTail = []rune{zero, 'a', 'c', 'h', 'e', 'l', 'o', 'r', '#', 'a', 'r', '#'}

	if !slices.Equal(trie.tail, expectedTail) {
		t.Errorf("expecting tail=%v, got %v", expectedTail, trie.tail)
	}

	if trie.pos != 12 {
		t.Errorf("expecting pos 12 , got %v", trie.pos)
	}

	//Case 3: insertion, when a collision occurs
	trie.Insert([]rune("badge"), 102)

	//if trie.Len() != 3 {
	//	t.Errorf("expecting len 3, got %v", trie.Len())
	//}
	//
	//if !trie.Contain([]rune("badge")) {
	//	t.Errorf("expecting to find key=badge")
	//}
	//
	//if value, found := trie.Get([]rune("badge")); value != 102 || !found {
	//	t.Errorf("expecting value 102, got %v", value)
	//}

	expectedBase = []int{0, 1, 0, 1, 1, -1, -12, 0, 0, 0, 0, 0, -9}

	if !slices.Equal(trie.dArray._base(), expectedBase) {
		t.Errorf("expecting base=%v, got %v", expectedBase, trie.dArray._base())
	}

	expectedCheck = []int{0, 0, 0, 4, 1, 3, 3, 0, 0, 0, 0, 0, 1}

	if !slices.Equal(trie.dArray._check(), expectedCheck) {
		t.Errorf("expecting check=%v, got %v", expectedCheck, trie.dArray._check())
	}

	expectedTail = []rune{zero, 'h', 'e', 'l', 'o', 'r', '#', zero, zero, 'a', 'r', '#', 'g', 'e', '#'}

	if !slices.Equal(trie.tail, expectedTail) {
		t.Errorf("expecting tail=%v, got %v", expectedTail, trie.tail)
	}
}
