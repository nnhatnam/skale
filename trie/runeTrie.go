package trie

type RuneTrie struct {
	trie *Trie[rune]
}

func NewRuneTrie() *RuneTrie {
	return &RuneTrie{trie: New[rune]()}
}

func (t *RuneTrie) Insert(s string) {
	t.trie.Insert([]rune(s))
}

func (t *RuneTrie) Get(s string) bool {
	return t.trie.Get([]rune(s))
}

func (t *RuneTrie) Delete(s string) {
	t.trie.Delete([]rune(s))
}

func (t *RuneTrie) Size() int {
	return t.trie.ElemCount()
}

func (t *RuneTrie) LongestPrefix(s string) string {
	return string(t.trie.LongestPrefix([]rune(s)))
}

func (t *RuneTrie) GetAll() []string {
	elems := t.trie.GetAll()
	s := make([]string, len(elems))
	for i, e := range elems {
		s[i] = string(e)
	}
	return s
}
