package trie

type Node struct {
	children map[rune]*Node
	last     bool
}

type Trie struct {
	root *Node
	size int
}

func New() *Trie {
	return &Trie{root: &Node{children: make(map[rune]*Node)}}
}

func (t *Trie) insert(r []rune) {
	node := t.root
	for _, c := range r {
		if _, ok := node.children[c]; !ok {
			node.children[c] = &Node{children: make(map[rune]*Node)}
		}
		node = node.children[c]
	}
	if !node.last {
		node.last = true
		t.size++
	}
}

func (t *Trie) get(r []rune) *Node {
	node := t.root
	for _, c := range r {
		if _, ok := node.children[c]; !ok {
			return nil
		}
		node = node.children[c]
	}
	return node
}

func (t *Trie) search(r []rune) bool {
	node := t.root
	for _, c := range r {
		if _, ok := node.children[c]; !ok {
			return false
		}
		node = node.children[c]
	}
	return node.last
}

func (t *Trie) delete(r []rune) {
	node := t.root
	for _, c := range r {
		if _, ok := node.children[c]; !ok {
			return
		}
		node = node.children[c]
	}
	if node.last {
		node.last = false
		t.size--
	}
}

func (t *Trie) allWordsFrom(node *Node, prefix string) []string {
	var words []string
	if node.last {
		words = append(words, prefix)
	}
	for c, n := range node.children {
		words = append(words, t.allWordsFrom(n, prefix+string(c))...)
	}
	return words
}

func (t *Trie) allWords() []string {
	return t.allWordsFrom(t.root, "")
}

func (t *Trie) contains(s string) bool {
	node := t.root
	for _, c := range []rune(s) {
		if _, ok := node.children[c]; !ok {
			return false
		}
		node = node.children[c]
	}
	return node.last
}

func wordsWithPrefix(node *Node, prefix string) []string {
	var words []string
	if node.last {
		words = append(words, prefix)
	}
	for c, n := range node.children {
		words = append(words, wordsWithPrefix(n, prefix+string(c))...)
	}
	return words
}

func (t *Trie) longestPrefixOf(s string) string {
	node := t.root
	var prefix string
	for _, c := range []rune(s) {
		if _, ok := node.children[c]; !ok {
			return prefix
		}
		prefix += string(c)
		node = node.children[c]
	}
	return prefix
}

func (t *Trie) Insert(s string) {
	t.insert([]rune(s))
}

func (t *Trie) Search(s string) bool {
	return t.search([]rune(s))
}

func (t *Trie) Delete(s string) {
	t.delete([]rune(s))
}

func (t *Trie) WordsWithPrefix(s string) []string {
	return wordsWithPrefix(t.get([]rune(s)), s)
}

func (t *Trie) AllWords() []string {
	return t.allWords()
}

func (t *Trie) Contains(s string) bool {
	return t.contains(s)
}

func (t *Trie) LongestPrefixOf(s string) string {
	return t.longestPrefixOf(s)
}
