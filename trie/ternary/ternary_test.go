package ternary

import (
	"strings"
	"testing"

	. "github.com/nnhatnam/skale/list/slice"
)

func TestPut(t *testing.T) {
	trie := New[rune]()
	trie.Insert([]rune("apple"))
	if !trie.Get([]rune("apple")) {
		t.Error("Failed to Put word")
	}
	trie.Insert([]rune("app"))
	if !trie.Get([]rune("app")) {
		t.Error("Failed to Put prefix")
	}
	trie.Insert([]rune("application"))
	if !trie.Get([]rune("application")) {
		t.Error("Failed to Put word")
	}
}

func TestSearch(t *testing.T) {
	trie := New[rune]()
	trie.Insert([]rune("apple"))
	if !trie.Get([]rune("apple")) {
		t.Error("Failed to Search word")
	}
	if trie.Get([]rune("appl")) {
		t.Error("Searching for prefix should return false")
	}
	if trie.Get([]rune("app")) {
		t.Error("Searching for prefix should return false")
	}
	if trie.Get([]rune("application")) {
		t.Error("Failed to Search word")
	}
}

func TestDelete(t *testing.T) {
	trie := New[rune]()
	trie.Insert([]rune("apple"))

	trie.Insert([]rune("application"))
	trie.Delete([]rune("apple"))

	if trie.Get([]rune("apple")) {
		t.Error("Failed to delete word")
	}

	if !trie.findPrefix([]rune("app")) {
		t.Error("Deleting word should not delete prefix")
	}
	trie.Delete([]rune("application"))
	if trie.Get([]rune("application")) {
		t.Error("Failed to delete word")
	}
	if trie.findPrefix([]rune("app")) {
		t.Error("Deleting word should delete prefix")
	}
}

func TestValuesWithPrefix(t *testing.T) {
	trie := New[rune]()

	trie.Insert([]rune("application"))
	trie.Insert([]rune("apple"))
	trie.Insert([]rune("app"))
	trie.Insert([]rune("apply"))
	trie.Insert([]rune("boy"))
	trie.Insert([]rune("bat"))
	trie.Insert([]rune("batman"))
	elems := trie.ValuesWithPrefix([]rune("ap"))

	if len(elems) != 4 {
		t.Error("Incorrect number of keys with prefix")
	}

	Slice[[]rune](elems).SortBy(func(a, b []rune) bool {
		return string(a) < string(b)
	})

	if !(string(elems[0]) == "app" && string(elems[1]) == "apple" && string(elems[2]) == "application" && string(elems[3]) == "apply") {
		t.Error("Incorrect keys with prefix")
	}
}

//func TestKeysThatMatch(t *testing.T) {
//	trie := New()
//	trie.Insert("application")
//	trie.Insert("apple")
//	trie.Insert("app")
//	trie.Insert("apply")
//	trie.Insert("boy")
//	trie.Insert("bat")
//	trie.Insert("batman")
//	keys := trie.KeysThatMatch(".ppl.")
//	if len(keys) != 2 {
//		t.Error("Incorrect number of keys that match pattern")
//	}
//	if !(keys[0] == "apple" && keys[1] == "apply") {
//		t.Error("Incorrect keys that match pattern")
//	}
//}

func TestLongestPrefixOf(t *testing.T) {
	trie := New[rune]()
	trie.Insert([]rune("application"))
	trie.Insert([]rune("apple"))
	trie.Insert([]rune("app"))
	trie.Insert([]rune("apply"))
	prefix := trie.LongestPrefixOf([]rune("applicable"))

	if string(prefix) != "applica" {
		t.Error("Incorrect longest prefix")
	}
	prefix = trie.LongestPrefixOf([]rune("boy"))

	if string(prefix) != "" {
		t.Error("Incorrect longest prefix")
	}

	IPList := []string{
		"128",
		"128.112",
		"128.112.055",
		"128.112.055.15",
		"128.112.136",
		"128.112.155.11",
		"128.112.155.13",
		"128.222",
		"128.222.136",
	}

	IPTrie := New[string]()
	for _, ip := range IPList {
		IPTrie.Insert(strings.Split(ip, "."))
	}

	longestPrefixTable := map[string]string{
		"128.112.136.11": "128.112.136",
		"128.112.100.16": "128.112",
		"128.166.123.45": "128",
	}

	for query, result := range longestPrefixTable {
		longestPrefix := IPTrie.LongestPrefixOf(strings.Split(query, "."))
		if strings.Join(longestPrefix, ".") != result {
			t.Errorf("Incorrect longest prefix. Expected %s, got %s", result, strings.Join(longestPrefix, "."))
		}

	}

}
