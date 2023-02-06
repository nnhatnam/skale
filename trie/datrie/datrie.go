package mptrie

import (
	"github.com/nnhatnam/skale/exp/xslices"
	"github.com/nnhatnam/skale/trie"
	"golang.org/x/exp/slices"
)

type DATrie[T trie.Elem] struct {
	*dArray
	tail []T

	pos int //current position in tail

	alphaMap ArcLabels[T] //a collection of valid arc labels, and their corresponding codes.

}

func NewDATrie[T trie.Elem]() *DATrie[T] {
	t := &DATrie[T]{}
	t.dArray = newDArray(1000)

	t.tail = make([]T, 1000)
	t.pos = 1
	return t
}

//func (da *DATrie[T]) prevState(t int) int {
//	return da.states[t].check
//}

// getNextState find the next state of state s with input c.
// if the new state is available, it registers the state in the check array and returns the new state
// if the new state is not available, it returns the state and false
//func (da *DATrie[T]) registerNextState(s int, c T) (t int, success bool) {
//	// (s) --c--> (t)
//	t = da.states[s].base + da.alphaMap.Code(c)
//	if da.states[t].check == 0 {
//		da.states[t].check = s
//		return t, true
//	} else if da.states[t].check == s {
//		return t, true
//	}
//	return t, false
//}

// insertTail inserts a slice of elements to the tail of the trie in the position pos, then updates the pos accordingly.
func (dat *DATrie[T]) insertTail(pos int, l []T) (success bool) {

	slices.Insert(dat.tail, pos, l...)
	dat.tail[pos+len(l)+1] = dat.alphaMap.StopElement()
	return true

}

func (dat *DATrie[T]) term(pos int) int {
	return slices.Index(dat.tail[pos:], dat.alphaMap.StopElement())
}

func (dat *DATrie[T]) xCheckAt(s int, c T) int {

	for dat.nextState(s, dat.alphaMap.Code(c)) > 0 {
		s++
	}

	return s

}

func (dat *DATrie[T]) xCheck(codes ...T) int {

	q := 1

	for idx := 0; idx < len(codes); idx++ {
		if q != dat.xCheckAt(q, codes[idx]) {
			q++
			idx = 0
		}
	}

	return q
}

func (dat *DATrie[T]) findAllArcsLeaving(s int) []T {

	var children []T

	for i := 1; i <= dat.biggestIdx; i++ {
		if dat.check(i) == s {
			children = append(children, dat.alphaMap.Get(i-dat.base(s)))
		}
	}
	return children
}

func (dat *DATrie[T]) lookup(value []T) bool {
	s := 1 //start from root (state 1)
	for idx := 1; idx < len(value); idx++ {
		c := dat.alphaMap.Code(value[idx])
		t := dat.nextState(s, c)
		if dat.check(t) != s {
			return false
		}

		if dat.base(t) < 0 {
			//if base is negative, then it is a leaf
			//we need to check the tail
			pos := -dat.base(t)
			term := dat.term(pos)
			if term == 0 {
				return false
			}

			leaf := dat.tail[pos:term]
			if !slices.Equal(leaf, value[idx:]) {
				return false
			}
			break
		}
		s = t
	}
	return true
}

func (dat *DATrie[T]) insert(value []T) {
	//TODO: make sure enough space
	s := 1 //start from root (state 1)
	for idx := 1; idx < len(value); idx++ {
		c := dat.alphaMap.Code(value[idx])
		t, success := dat.registerNextState(s, c)
		if success {

			// if the next state is successfully registered, check the base
			// if base value is 0, then set it to -pos and add the tail
			// if base value is negative, we may consider to split the tail
			if dat.base(t) == 0 {
				dat.setBase(t, -dat.pos)
				dat.insertTail(dat.pos, value[idx:])
				dat.pos += len(value[idx:]) + 1
				break
			} else if dat.base(t) < 0 {

				//check if we need to split the tail
				temp := -dat.base(t)
				leaf := dat.tail[temp:dat.term(temp)]

				offset := xslices.LongestPrefixIndex(leaf, value[idx:])
				if offset == len(leaf)-1 && len(leaf) == len(value[idx:]) {
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
					q := dat.xCheck(leaf[offset+1], value[idx+1])
					dat.setBase(t, q)

					//re-register the leaf
					t, _ = dat.registerNextState(t, dat.alphaMap.Code(leaf[offset+1]))
					dat.setBase(t, -temp)
					dat.insertTail(temp, leaf[offset+1:])
					//re-register the new value
					t, _ = dat.registerNextState(t, dat.alphaMap.Code(value[idx+1]))
					dat.setBase(t, -dat.pos)
					dat.insertTail(dat.pos, value[idx+1:])
					dat.pos += len(value[idx+1:]) + 1
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
					dat.moveState(newNode, oldNode)
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
					dat.moveState(newNode, oldNode)
				}

			}

			//re-register the new value
			t, _ = dat.registerNextState(s, c)
			dat.setBase(t, -dat.pos)
			dat.insertTail(dat.pos, value[idx:])
			dat.pos += len(value[idx:]) + 1
			break

		}

	}
}

func (dat *DATrie[T]) delete(value []T) bool {
	s := 1 //start from root (state 1)
	for idx := 1; idx < len(value); idx++ {
		c := dat.alphaMap.Code(value[idx])
		t := dat.nextState(s, c)

		if dat.base(t) < 0 {
			//if base is negative, then it is a leaf
			//we need to check the tail
			pos := -dat.base(t)
			term := dat.term(pos)
			if term == 0 {
				return false
			}

			leaf := dat.tail[pos:term]
			if !slices.Equal(leaf, value[idx:]) {
				return false
			} else {
				//delete the leaf
				xslices.Reset(dat.tail, pos, term)
				if term == dat.pos {
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
