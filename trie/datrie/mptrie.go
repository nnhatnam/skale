package mptrie

import (
	"github.com/nnhatnam/skale/trie"
	"golang.org/x/exp/slices"
)

type MPTrie[T trie.Elem] struct {
	base  []int
	check []int
	tail  []T

	pos int //current position in tail

	alphaMap ArcLabels[T] //a collection of valid arc labels, and their corresponding codes.

}

func New[T trie.Elem]() *MPTrie[T] {
	da := &MPTrie[T]{}
	da.base = make([]int, 1000)
	da.base[1] = 1
	da.check = make([]int, 1000)

	da.tail = make([]T, 1000)
	da.pos = 1
	return da
}

func (mp *MPTrie[T]) prevState(t int) int {
	return mp.check[t]
}

// nextState identified the next state of state s with input c
func (mp *MPTrie[T]) nextState(s int, c T) int {
	return mp.base[s] + mp.alphaMap.Code(c)
}

// setPrevState sets the previous state of the current state i to p.
func (mp *MPTrie[T]) setPrevState(i, p int) {
	//In MP-tries, the state is identified by array indices. If we take an arbitrary state i, check[i] is the previous state of node i, if check[i] != 0
	mp.check[i] = p
}

func (mp *MPTrie[T]) term(pos int) int {
	return slices.Index(mp.tail[pos:], mp.alphaMap.StopElement())
}

// xCheckCode returns the state that is reached by following the path from the root to the state s, and then following the arc labeled by c.
func (mp *MPTrie[T]) xCheckCodeAt(s int, c T) int {

	for mp.nextState(s, c) > 0 {
		s++
	}

	return s

}

func (mp *MPTrie[T]) xCheckCodes(l []T) int {

	q := 1
	for _, v := range l {
		q = mp.xCheckCodeAt(q, v)

	}

	return q
}

func (mp *MPTrie[T]) checkTailCollision(q int, l []T) bool {

	for _, v := range l {
		if mp.nextState(q, v) > 0 {
			return true
		}
	}

	return false

}

// checkCollision checks if there is a collision between two slices. It's pretty much similar to finding the longest common prefix of two slices.
// It returns the index of the first element that is different, and a boolean value indicating if there is a collision.
func checkCollision[T trie.Elem](s1, s2 []T) (_ int, col bool) {
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
		col = true
	}

	for i, v := range s1 {
		if v != s2[i] {
			return i, true //collision
		}
	}
	return len(s1) - 1, col
}

// getNextState find the next state of state s with input c.
// if the new state is available, it registers the state in the check array and returns the new state
// if the new state is not available, it returns the state and false
func (mp *MPTrie[T]) registerNextState(s int, c T) (t int, success bool) {
	// (s) --c--> (t)
	t = mp.base[s] + mp.alphaMap.Code(c)
	if mp.check[t] == 0 {
		mp.check[t] = s
		return t, true
	} else if mp.check[t] == s {
		return t, true
	}
	return
}

// insertTail inserts a slice of elements to the tail of the trie in the position pos, then updates the pos accordingly.
func (mp *MPTrie[T]) insertTail(l []T) (success bool) {

	slices.Insert(mp.tail, mp.pos, l...)
	mp.pos += len(l)
	mp.tail[mp.pos] = mp.alphaMap.StopElement()
	mp.pos++
	return true

}

func (mp *MPTrie[T]) insert(value []T) {
	//TODO: make sure enough space
	s := 1 //start from root (state 1)
	for idx := 1; idx < len(value); idx++ {
		c := value[idx]
		if t, success := mp.registerNextState(s, c); success {
			// if the next state is successfully registered, check the base
			// if base value is 0, then set it to -pos and add the tail
			// if base value is negative, we may consider to split the tail

			if mp.base[t] == 0 {
				mp.base[t] = -mp.pos
				mp.insertTail(value[idx:])
				break
			} else if mp.base[t] < 0 {
				//consider to split the tail
				temp := -mp.base[t]
				leaf := mp.tail[temp:mp.term(temp)]
				//check if there is a collision between the leaf and the new value
				offset := longestPrefixIndex(leaf, value[idx:])

				if offset == len(leaf)-1 && len(leaf) == len(value[idx:]) {
					//leaf == value[idx:]
					//the leaf is the same as the new value, do nothing
					break
				} else if offset >= 0 {
					for i := 0; i <= offset; i++ {
						q := mp.xCheckCodeAt(t, leaf[i])
						mp.base[t] = q
						s = t
						//s--leaf[i]-->t
						t, _ = mp.registerNextState(t, leaf[i])

					}
				}

			}
			s = t
			continue
		}

	}
}
