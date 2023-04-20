package nary

import (
	"github.com/nnhatnam/skale"
	"github.com/nnhatnam/skale/list/skiplist"
)

type Node[K skale.Ordered, V any] struct {
	Key   K
	Value V
	Child skiplist.SkipList[K]
}
