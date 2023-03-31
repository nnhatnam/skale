package radix

import "github.com/nnhatnam/skale/trie"

type ItemIterator[K trie.Elem, V any] func(key []K, v V) bool
