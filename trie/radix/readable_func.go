package radix

import "github.com/nnhatnam/skale/trie"

func Keys[K trie.Elem, V any](m *RadixMap[K, V]) [][]K {

	result := make([][]K, 0, m.Len())

	m.Ascend(func(k []K, v V) bool {
		result = append(result, k)
		return false
	})

	return result
}

func Values[K trie.Elem, V any](m *RadixMap[K, V]) []V {

	result := make([]V, 0, m.Len())

	m.Ascend(func(k []K, v V) bool {
		result = append(result, v)
		return false
	})

	return result
}
