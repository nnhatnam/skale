// Package datrie implements a double-array trie inspired by the algorithm described in the paper:
// An Efficient Implementation of Trie Structures (https://doi.org/10.1002/spe.4380220902)
package datrie

import (
	"github.com/nnhatnam/skale/exp/xslices"
	"github.com/nnhatnam/skale/trie"
	"golang.org/x/exp/slices"
)

// DATrie implements a double-array trie.
type DATrie[K trie.Elem, V any] struct {

	//dArray is a double-array that stores the trie.
	*dArray[V]

	//tail is used to store non-branching suffixes. (leaf nodes)
	tail []K

	//position of the tail
	pos int

	//a collection of valid arc labels, and their corresponding codes.
	alphaMap ArcDomain[K]

	//number of words in the trie
	len int
}

// New returns a new empty trie. The trie will use the given arcMap to map arc labels to their corresponding codes.
func New[K trie.Elem, V any](arcMap ArcDomain[K]) *DATrie[K, V] {
	t := &DATrie[K, V]{}
	t.dArray = newDArray[V](2)
	t.setBase(1, 1)

	t.tail = make([]K, 2)
	t.pos = 1

	t.alphaMap = arcMap
	return t
}

// writeTail inserts a slice of elements to the tail of the trie in the position pos, then updates the pos accordingly.
func (dat *DATrie[K, V]) writeTail(pos int, l []K) (success bool) {
	l = append(l, dat.alphaMap.StopElement())
	dat.tail = xslices.Insert(dat.tail, pos, l...)

	return true

}

// readTail returns the slice of elements that ends at the given position in the tail.
func (dat *DATrie[K, V]) readTail(pos int) []K {
	return xslices.Walk(dat.tail, pos, dat.alphaMap.StopElement())
}

func (dat *DATrie[K, V]) xCheckAt(s int, c K) int {

	for dat.nextState(s, dat.alphaMap.Code(c)) > 0 {
		s++
	}

	return s

}

func (dat *DATrie[K, V]) xCheck(codes ...K) int {

	q := 1

	for idx := 0; idx < len(codes); idx++ {
		if q != dat.xCheckAt(q, codes[idx]) {
			q++
			idx = 0
		}
	}

	return q
}

func (dat *DATrie[K, V]) findAllArcsLeaving(s int) []K {

	var children []K

	for i := 1; i <= len(dat.states); i++ {
		if dat.check(i) == s {
			children = append(children, dat.alphaMap.Label(i-dat.base(s)))
		}
	}
	return children
}

func (dat *DATrie[K, V]) lookup(key []K) (value V, found bool) {
	s := 1 //start from root (state 1)
	for idx := 0; idx < len(key); idx++ {
		c := dat.alphaMap.Code(key[idx])

		if c > len(dat.states) {
			return value, false
		}

		t := dat.nextState(s, c)
		if dat.check(t) != s {
			return value, false
		}

		if dat.base(t) < 0 {
			//if base is negative, then it is a leaf
			//we need to check the tail

			pos := -dat.base(t)

			if !slices.Equal(dat.readTail(pos), key[idx+1:]) {
				return value, false
			}

			return dat.value(t), true
		}
		s = t
	}
	return dat.value(s), true
}

func (dat *DATrie[K, V]) insertOrReplace(key []K, value V) (success bool) {

	s := 1 //start from root (state 1)
	t := 0
	for idx := 0; idx < len(key); idx++ {
		c := dat.alphaMap.Code(key[idx])

		//try to register for the next state
		t, success = dat.registerNextState(s, c)

		if success {

			// if the next state is successfully registered, we need to update the base
			//since registerNextState only updates the check.
			// if base value is 0 and the word is not end, then set it to -pos and the remaining part to the tail
			// if base value is negative, we may consider to split the tail
			// if base value is 0 and the word is end, we know that it has no children, so we can pick any base value > 0.
			// in this case, we set it to 1
			if dat.base(t) == 0 {
				//Case 1: base = 0

				dat.len++

				if idx == len(key)-1 {
					//if we are at the end of the key, then we can set the value directly
					dat.setValue(t, value)
					dat.setEnd(t, true)
					dat.setBase(t, 1)
					return true
				}

				dat.setBase(t, -dat.pos) //base = -pos
				dat.setValue(t, value)

				dat.writeTail(dat.pos, key[idx+1:])
				dat.pos += len(key[idx+1:])

				break

			} else if dat.base(t) < 0 {

				//check if we need to split the tail
				temp := -dat.base(t)
				leaf := dat.readTail(temp)

				offset := xslices.LongestPrefixIndex(leaf, key[idx:])
				if offset == len(leaf)-1 && len(leaf) == len(key[idx:]) {
					//leaf == value[idx:]
					//the leaf is the same as the new value, do nothing
					break
				} else if offset >= 0 {
					for i := 0; i <= offset; i++ {
						q := dat.xCheckAt(t, leaf[i])
						dat.setBase(t, q)

						//s--leaf[i]-->t
						t, _ = dat.registerNextState(t, dat.alphaMap.Code(leaf[i]))
					}

					idx = idx + offset
					//s--leaf[offset+1:]-->t
					q := dat.xCheck(leaf[offset+1], key[idx+1])
					dat.setBase(t, q)

					//re-register the leaf
					t, _ = dat.registerNextState(t, dat.alphaMap.Code(leaf[offset+1]))
					dat.setBase(t, -temp)
					dat.writeTail(temp, leaf[offset+1:])
					//re-register the new value
					t, _ = dat.registerNextState(t, dat.alphaMap.Code(key[idx+1]))
					dat.setBase(t, -dat.pos)
					dat.writeTail(dat.pos, key[idx+1:])
					dat.pos += len(key[idx+1:]) + 1
					break //done
				}

			}

			s = t
			continue
		} else {
			//collision occurs. s or t's parent must be moved

			//s--c-->t (t is not available)
			//find all arcs leaving s
			sArcs := dat.findAllArcsLeaving(s)

			//find all arcs leaving t's parent
			tPArcs := dat.findAllArcsLeaving(dat.check(t))

			//we move what ever has fewer branches
			if len(sArcs)+1 < len(tPArcs) {

				sOldBase := dat.base(s)
				q := dat.xCheck(sArcs...)
				dat.setBase(s, q)

				//move all arcs leaving s to the new base
				for _, arc := range sArcs {
					oldNode := sOldBase + dat.alphaMap.Code(arc)
					newNode := dat.base(s) + dat.alphaMap.Code(arc)
					dat.relocateState(newNode, oldNode)
				}

			} else {
				parentNode := dat.check(t)
				parentOldBase := dat.base(parentNode)

				q := dat.xCheck(tPArcs...)
				dat.setBase(parentNode, q)

				//move all arcs leaving t's parent to the new base
				for _, arc := range tPArcs {
					oldNode := parentOldBase + dat.alphaMap.Code(arc)
					newNode := dat.base(parentNode) + dat.alphaMap.Code(arc)
					dat.relocateState(newNode, oldNode)
				}

			}

			//re-register the new value
			t, _ = dat.registerNextState(s, c)
			dat.setBase(t, -dat.pos)
			dat.writeTail(dat.pos, key[idx:])
			dat.pos += len(key[idx:]) + 1
			break

		}

	}
	return
}

func (dat *DATrie[K, V]) delete(value []K) bool {
	s := 1 //start from root (state 1)
	for idx := 1; idx < len(value); idx++ {
		c := dat.alphaMap.Code(value[idx])
		t := dat.nextState(s, c)

		if dat.base(t) < 0 {
			//if base is negative, then it is a leaf
			//we need to check the tail
			pos := -dat.base(t)

			leaf := dat.readTail(pos)
			if !slices.Equal(leaf, value[idx:]) {
				return false
			} else {
				//delete the leaf
				xslices.Reset(dat.tail, pos, len(leaf))
				if len(leaf) == dat.pos {
					dat.pos = pos - 1
				}
				dat.setBase(t, 0)
				dat.setCheck(t, 0)
			}
			break
		}
		s = t
	}
	return true
}

func (dat *DATrie[K, V]) Len() int {
	return dat.len
}

func (dat *DATrie[K, V]) Delete(value []K) bool {
	return dat.delete(value)
}

func (dat *DATrie[K, V]) Insert(key []K, value V) {
	dat.insertOrReplace(key, value)
}

func (dat *DATrie[K, V]) Contain(value []K) (found bool) {
	_, found = dat.lookup(value)
	return
}

func (dat *DATrie[K, V]) Get(key []K) (value V, found bool) {
	return dat.lookup(key)
}
