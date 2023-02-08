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

	if trie.pos != 8 {
		t.Errorf("expecting pos 8 , got %v", trie.pos)
	}

}
