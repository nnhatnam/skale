package datrie

// state represents a node in a double array trie.
// It stores the base and check values of the node, as well as a flag indicating whether the node is an end node (i.e. represents the end of a word)
type state[T any] struct {

	//The "base" value is a way to map a node to its children in a more efficient manner than traditional trie data structures.
	//It can be thought of as a pointer to the children of a node. The state number of a child is found by adding the base value to the child's character code.
	//If the base value is negative, then the remaining value is stored in the tail array. In this case, -base is the index of the first character of the tail.
	//The "check" value is used to validate the "base" value, and it's also used find the parent of a node if we know the node state number.
	// Example:
	// (s) ---c---> (t)  where s is the state number of the parent node, t is the state number of the child node, and c is the alphabet code of the arc label for the transition from s -> t. Then we hvae
	// base[s] + c = t
	// check[t] = s
	base, check int

	//end is a flag indicates if the node is an end node, meaning it represents the end of a word.
	end bool

	// value is the value associated with the key that ends at this node
	value T
}

// dArray represents a double array trie data structure
// It consists of two one-dimensional arrays (base and check) that store the nodes of the trie.
type dArray[T any] struct {

	//`states` stores all the nodes of the trie compressed into base and check arrays.
	//The `base` array is used to map nodes to its children, serving as pointers to them.
	//The `check` array is used to map nodes to its parent, serving as pointers to the parent node.
	//The index of the array is the state number of the node, and the value of the array is the node itself.
	//To keep it simple, instead of saying a node with state number `s`, we will just say node `s` (or just simply `s`).
	states []state[T]
}

// newDArray creates a new double array trie data structure with the given size
func newDArray[T any](size int) *dArray[T] {
	return &dArray[T]{states: make([]state[T], size)}
}

// for testing purpose only
func (da *dArray[T]) _base() []int {
	baseArr := make([]int, len(da.states))
	for k, v := range da.states {
		baseArr[k] = v.base
	}
	return baseArr
}

// for testing purpose only
func (da *dArray[T]) _check() []int {
	checkArr := make([]int, len(da.states))
	for k, v := range da.states {
		checkArr[k] = v.check
	}
	return checkArr
}

// base returns the base value of the node with the given state number
func (da *dArray[T]) base(s int) int {
	return da.states[s].base
}

// setBase sets the base value of the node s to b
// If s > biggestIdx, then the biggestIdx is updated to s
func (da *dArray[T]) setBase(s int, b int) {

	if s >= len(da.states) {
		da.states = append(da.states, make([]state[T], s-len(da.states)+1)...)
	}

	da.states[s].base = b

}

// check returns the check value of the node t
func (da *dArray[T]) check(t int) int {
	return da.states[t].check
}

// setCheck sets the check value of the node t to the given value s
// If the given state number is greater than the biggestIdx, then the biggestIdx is updated to the given state number
// It can also be thought of as setting node s as the parent of the node t (s -> t)
func (da *dArray[T]) setCheck(t int, s int) {

	if t >= len(da.states) {
		da.states = append(da.states, make([]state[T], t-len(da.states)+1)...)
	}

	da.states[t].check = s //check[t] = s

}

// value returns the value associated with the key that ends at node t
func (da *dArray[T]) value(t int) T {
	return da.states[t].value
}

// setValue sets the value associated with the key that ends at node t to v
func (da *dArray[T]) setValue(t int, v T) {
	da.states[t].value = v
}

// prevState returns the parent node s of node t
func (da *dArray[T]) prevState(t int) (s int) {
	return da.states[t].check
}

// nextState returns the child node t of node s with input c
func (da *dArray[T]) nextState(s int, c int) (t int) {
	return da.base(s) + c
}

// isEnd returns true if the node t is an end node, false otherwise
func (da *dArray[T]) isEnd(t int) bool {
	return da.states[t].end
}

// setEnd sets the end flag of the node t to b
func (da *dArray[T]) setEnd(t int, b bool) {
	da.states[t].end = b
}

// registerNextState tries to add the next node `t` to the existing node `s` with input `c`.
// If `t` is not registered yet, it registers `t` by updating the `check` array and returns `t` and `true`.
// the `base` array is not altered because it has no information about the next input (i.e. next child) of `t`.
// If `t` has already been added, the function checks if `t` was previously registered by node `s`. If so, return `t` and `true`.
// Otherwise, there is a conflict, return `t` and `false`.
// Resolving the conflict or altering the `base` array is the responsibility of the caller.
// @require: s is a valid node (check[s] != 0)
func (da *dArray[T]) registerNextState(s int, c int) (t int, success bool) {
	// (s) --c--> (t)

	t = da.base(s) + c

	//make sure that the states array is big enough to hold the new state
	if t >= len(da.states) {
		da.states = append(da.states, make([]state[T], t-len(da.states)+1)...)
	}

	if da.check(t) == 0 {
		// t is not registered yet
		da.setCheck(t, s)
		return t, true
	} else if da.check(t) == s {
		// t is already registered by s
		return t, true
	}

	// t is registered by another node
	return t, false
}

// relocateState move node s to node t. Reset node s and update its children's check value
func (da *dArray[T]) relocateState(t int, s int) {

	//copy base and check from s to t
	da.setCheck(t, da.check(s)) //check[t] = check[s]
	da.setBase(t, da.base(s))   //base[t] = base[s]
	da.setEnd(t, da.isEnd(s))   //end[t] = end[s]

	//update children's check value to point to t
	for i := 0; i < len(da.states); i++ {
		if da.check(i) == s {
			da.setCheck(i, t)
		}
	}

	//reset node s
	da.setBase(s, 0)
	da.setCheck(s, 0)
	da.setEnd(s, false)
}
