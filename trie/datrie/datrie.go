package datrie

import (
	"github.com/nnhatnam/skale/trie"
	"golang.org/x/exp/slices"
)

type DATrie[T trie.Elem] struct {
	base  []int
	check []int
	tail  []rune

	pos int //current position in tail

}

func New[T trie.Elem]() *DATrie[T] {
	da := &DATrie[T]{}
	da.base = make([]int, 1000)
	da.base[1] = 1
	da.check = make([]int, 1000)

	da.tail = make([]rune, 1000)
	da.pos = 1
	return da

}

func (da *DATrie[T]) nextState(state int, label rune) int {
	return da.base[state] + alphabet.Code(label)
}

func (da *DATrie[T]) xCheckCode(c rune, q int) int {

	for {
		if da.check[q+alphabet.Code(c)] == q {
			return q //found
		}
		q++
	}

}

func (da *DATrie[T]) xCheckCodes(l []rune) int {

	q := 1
	for _, v := range l {
		q = da.xCheckCode(v, q)

	}

	return q
}

// getLeaf returns the leaf at position pos of tail
func (da *DATrie[T]) getLeaf(pos int) []rune {
	result := make([]rune, 0)
	for da.tail[pos] != alphabet.StopElement() {
		result = append(result, da.tail[pos])
		pos++
	}

	return result
}

func (da *DATrie[T]) insert(value []rune) {
	value = append(value, alphabet.StopElement())
	//TODO: make sure enough space
	s := 1 //start from node 1
	for idx, c := range value {

		t := da.base[s] + alphabet.Code(c) //base[s] + c = t

		// no collision
		if da.check[t] == 0 {
			da.check[t] = s //check[t] = s
			da.base[t] = -da.pos

			copy(da.tail[da.pos:], value[idx:])
			da.pos += len(value[idx:])
			break
		} else {
			// there is already a node at t
			// check if the node is a leaf
			if da.base[t] < 0 {

				temp := -da.base[t]
				leaf := da.getLeaf(temp)

				if !slices.Equal(leaf, value[idx:]) {
					// different word
					// split the leaf
					// find the first different character
					for i := 0; i < len(leaf); i++ {
						if leaf[i] == value[idx+i] {
							q := da.xCheckCode(leaf[i], 1)
							da.base[t] = q
							// continue
						} else if leaf[i] != value[idx+i] {
							// split the leaf
							// create a new node
							// move the leaf to the new node
							// move the current word to the new node
							// update the base and check of the parent node
							// update the base and check of the new node
							// update the base and check of the current node
							// update the tail
							// update the pos
							// break

						}
					}
				} else {
					// same word
					break
				}

			}
		}
	}
}
