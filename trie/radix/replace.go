package radix

import "github.com/nnhatnam/skale/trie"

type ReplaceFunc[K trie.Elem, V any] func(old V) V

