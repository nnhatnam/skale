package trie

import "testing"

func TestPut(t *testing.T) {
	trie := New()
	trie.Insert("apple")
	if !trie.Search("apple") {
		t.Error("Failed to Put word")
	}
	trie.Insert("app")
	if !trie.Search("app") {
		t.Error("Failed to Put prefix")
	}
	trie.Insert("application")
	if !trie.Search("application") {
		t.Error("Failed to Put word")
	}
}

func TestSearch(t *testing.T) {
	trie := New()
	trie.Insert("apple")
	if !trie.Search("apple") {
		t.Error("Failed to Search word")
	}
	if trie.Search("appl") {
		t.Error("Searching for prefix should return false")
	}
	if !trie.Search("app") {
		t.Error("Failed to Search prefix")
	}
	if !trie.Search("application") {
		t.Error("Failed to Search word")
	}
}

func TestDelete(t *testing.T) {
	trie := New()
	trie.Insert("apple")
	trie.Insert("application")
	trie.Delete("apple")
	if trie.Search("apple") {
		t.Error("Failed to delete word")
	}
	if !trie.Search("app") {
		t.Error("Deleting word should not delete prefix")
	}
	trie.Delete("application")
	if trie.Search("application") {
		t.Error("Failed to delete word")
	}
	if !trie.Search("app") {
		t.Error("Deleting word should not delete prefix")
	}
}

func TestKeysWithPrefix(t *testing.T) {
	trie := New()
	trie.Insert("application")
	trie.Insert("apple")
	trie.Insert("app")
	trie.Insert("apply")
	trie.Insert("boy")
	trie.Insert("bat")
	trie.Insert("batman")
	keys := trie.WordsWithPrefix("ap")
	if len(keys) != 3 {
		t.Error("Incorrect number of keys with prefix")
	}
	if !(keys[0] == "app" && keys[1] == "apple" && keys[2] == "apply") {
		t.Error("Incorrect keys with prefix")
	}
}

//
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
	trie := New()
	trie.Insert("application")
	trie.Insert("apple")
	trie.Insert("app")
	trie.Insert("apply")
	prefix := trie.LongestPrefixOf("applicable")
	if prefix != "app" {
		t.Error("Incorrect longest prefix")
	}
	prefix = trie.LongestPrefixOf("boy")
	if prefix != "b" {
		t.Error("Incorrect longest prefix")
	}
}
