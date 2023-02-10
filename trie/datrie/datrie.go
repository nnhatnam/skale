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
	t.writeBase(1, 1)

	t.tail = make([]K, 1)
	t.pos = 1 // in most cases, pos = tail. Unless, we remove the last elements in the tail.

	t.alphaMap = arcMap
	return t
}

// writeTail inserts a slice of elements to the tail of the trie in the position pos.
func (dat *DATrie[K, V]) writeTail(pos int, key []K) {
	key = append(key, dat.alphaMap.StopElement())

	if pos == dat.pos {
		dat.tail = append(dat.tail, key...)
		dat.pos += len(key)
		return
	}

	i := pos
	j := i + len(key)
	var zero K
	for i < dat.pos {
		if i < j {
			dat.tail[i] = key[i-pos]
		} else {

			if dat.tail[i] == dat.alphaMap.StopElement() {
				dat.tail[i] = zero
				break
			}
			dat.tail[i] = zero

		}
		i++
	}

}

// readTail walks through the tail from position `pos` until it meets the first stop element.
// It returns elements starting from `pos` up to, but not including, the stop element.
// @require: pos < len(dat.tail)
func (dat *DATrie[K, V]) readTail(pos int) []K {

	j := pos
	for j < len(dat.tail) {
		if dat.tail[j] == dat.alphaMap.StopElement() {
			break
		}
		j++
	}
	return dat.tail[pos:j]

}

func (dat *DATrie[K, V]) xCheck(codes ...K) int {

	q := 1

	for idx := 0; idx < len(codes); idx++ {

		t := q + dat.alphaMap.Code(codes[idx])

		if dat.base(t) > 0 {
			q++
			idx = 0
		}
	}

	return q
}

func (dat *DATrie[K, V]) registerNextState(s int, c K) (int, bool) {
	return dat.dArray.registerNextState(s, dat.alphaMap.Code(c))
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

// This is a sub-task of insertOrReplace. It is used to update the base of `t` after `t` got successfully registered with no collision (base[s] = 0).,
// and is only called when there is no collision.
func (dat *DATrie[K, V]) writeBaseNoCollision(t int, key []K, value V) (success bool) {

	dat.len++

	if len(key) > 0 {
		dat.writeBase(t, -dat.pos)
		dat.writeTail(dat.pos, key)
	}

	dat.setEnd(t, true)
	dat.setValue(t, value)

	if dat.base(t) == 0 {
		dat.writeBase(t, 1)
	}

	return false
}

// This is a sub-task of insertOrReplace. It is used to update the base of `t` after `t` got successfully registered but the base is negative.,
// and is only called when there is collision.
func (dat *DATrie[K, V]) writeNegativeBase(t int, key []K, value V) (success bool) {
	temp := -dat.base(t)
	leaf := dat.readTail(temp)

	// if the leaf is the same as the key, then we just need to update the value
	if slices.Equal(leaf, key) {
		dat.setValue(t, value)
		return true
	}

	// if the leaf is not the same as the key, then we need to split the leaf
	// and insert the new key.

	// first, we need to find the first element that is different between the leaf and the key
	// we will split the leaf at this element.
	var idx int
	for idx = 0; idx < len(leaf) && idx < len(key); idx++ {
		if leaf[idx] != key[idx] {
			break
		}
	}

	// if the leaf is a prefix of the key, then we need to split the key at the first element that is different on the key.
	if idx == len(leaf) { //leaf is a prefix of the key
		for i := 0; i < idx; i++ {
			q := dat.xCheck(key[i])
			dat.writeBase(t, q)
			t, _ = dat.registerNextState(t, key[i])
		}

		dat.setEnd(t, true)
		dat.setValue(t, value)

	}

	// if the leaf is a prefix of the key, then we need to split the leaf at the first element that is different.
	// otherwise, we need to split the leaf at the first element that is different + 1.
	// this is because we need to make sure that the leaf is a prefix of the key.
	if len(leaf) < len(key) {
		idx++
	}

	// we need to find the first available state to split the leaf
	// we will use the first available state to split the leaf.
	// we will use the second available state to insert the new key.
	s := dat.xCheck(leaf[:idx]...)

}

func (dat *DATrie[K, V]) insertOrReplace(key []K, value V) (success bool) {

	s := 1 //start from root (state 1)
	t := 0

	//Essentially, in each iteration, we record the next state (t) based on the current state (s).
	//Most of the operations are focused on ensuring that the new state (t) is accurately recorded, assuming that s has already been properly recorded.
	for idx := 0; idx < len(key); idx++ {
		c := key[idx]

		//try to register for the next state
		t, success = dat.registerNextState(s, c)

		if success {

			// if the next state (t) is successfully registered, we need to update the base of t
			// since registerNextState only updates the check.
			// if base value of `t` is 0, then set it to -pos and the remaining part to the tail
			// if base value of `t` is negative, we may consider to split the tail
			if dat.base(t) == 0 {
				//Case 1: base[t] = 0

				dat.writeBaseNoCollision(t, key[idx+1:], value)

				break

			} else if dat.base(t) < 0 {

				//check if we need to split the tail
				temp := -dat.base(t)
				leaf := dat.readTail(temp)

				offset := xslices.LongestPrefixIndex(leaf, key[idx+1:])

				if offset == len(leaf)-1 && len(leaf) == len(key[idx+1:]) {
					//leaf == value[idx:]
					//the leaf is the same as the new value, do nothing
					break
				} else if offset >= 0 {

					for i := 0; i <= offset; i++ {
						q := dat.xCheck(leaf[i])
						dat.writeBase(t, q)

						//s--leaf[i]-->t
						t, _ = dat.registerNextState(t, leaf[i])
					}

					idx = idx + offset + 1
					//s--leaf[offset+1:]-->t
					q := dat.xCheck(leaf[offset+1], key[idx+1])
					dat.writeBase(t, q)

					//re-register the leaf
					t1, _ := dat.registerNextState(t, leaf[offset+1])
					dat.writeBase(t1, -temp)
					dat.writeTail(temp, leaf[offset+2:])

					//re-register the new value
					t2, _ := dat.registerNextState(t, key[idx+1])
					dat.writeBase(t2, -dat.pos)
					dat.writeTail(dat.pos, key[idx+2:])
					//dat.pos += len(key[idx+1:]) + 1
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
				dat.writeBase(s, q)

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
				dat.writeBase(parentNode, q)

				//move all arcs leaving t's parent to the new base
				for _, arc := range tPArcs {
					oldNode := parentOldBase + dat.alphaMap.Code(arc)
					newNode := dat.base(parentNode) + dat.alphaMap.Code(arc)
					dat.relocateState(newNode, oldNode)
				}

			}

			//re-register the new value
			t, _ = dat.registerNextState(s, c)
			dat.writeBase(t, -dat.pos)
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
				dat.writeBase(t, 0)
				dat.writeCheck(t, 0)
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
