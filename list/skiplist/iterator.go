package skiplist

import "github.com/nnhatnam/skale"

type ItemIterator[T any] func(item T) bool

type MapIterator[K skale.Ordered, V any] func(key K, value V) bool
