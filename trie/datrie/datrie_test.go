package datrie

import "testing"

func TestInitialState(t *testing.T) {

	trie := New[rune](alphabetRune)
	if trie.Size() != 0 {
		t.Errorf("expecting len 0")
	}

	if trie.Contain([]rune("test")) {
		t.Errorf("not expecting to find key=test")
	}

}
