package trie

import "strings"

//TODO: Implement path trie based on https://codecatalog.org/articles/zookeeper-trie/

type PathConverter interface {
	Split() []string
	Merge(p []string) string
}

type path string

func (p path) Split() []string {
	return strings.Split(string(p), "/") // For now, use PathConverter if we need to support other separators
}

func (p path) Merge(s []string) string {
	return strings.Join(s, "/")
}

type PathTrie struct {
	trie *Trie[string]
	p    PathConverter
}

func NewPathTrie() *PathTrie {
	return &PathTrie{trie: New[string]()}
}

func (t *PathTrie) Insert(s string) {
	t.trie.Insert(path(s).Split())
}

func (t *PathTrie) Find(s string) bool {
	return t.trie.Get(path(s).Split())
}

func (t *PathTrie) Remove(s string) {
	t.trie.Delete(path(s).Split())
}

func (t *PathTrie) Size() int {
	return t.trie.ElemCount()
}

func (t *PathTrie) LongestPrefix(s string) []string {
	return t.trie.LongestPrefix(path(s).Split())
}

func (t *PathTrie) GetAll() []string {
	paths := make([]string, t.Size()+1, t.Size()+1)
	for _, p := range t.trie.GetAll() {
		paths = append(paths, path("").Merge(p))
	}
	return paths
}
