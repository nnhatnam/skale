package mptrie

// double array item
type state struct {
	base, check int
	end         bool
}

type dArray struct {
	states []state

	biggestIdx int
}

func newDArray(size int) *dArray {
	return &dArray{states: make([]state, size)}
}

func (da *dArray) base(s int) int {
	return da.states[s].base
}

func (da *dArray) setBase(s int, b int) {
	da.states[s].base = b
	if s > da.biggestIdx {
		da.biggestIdx = s
	}
}

func (da *dArray) check(t int) int {
	return da.states[t].check
}

func (da *dArray) setCheck(t int, s int) {
	da.states[t].check = s
	if t > da.biggestIdx {
		da.biggestIdx = t
	}
}

func (da *dArray) prevState(t int) int {
	return da.states[t].check
}

func (da *dArray) nextState(s int, c int) int {
	return da.base(s) + c
}

func (da *dArray) isEnd(t int) bool {
	return da.states[t].end
}

func (da *dArray) setEnd(t int) {
	da.states[t].end = true
}

// getNextState find the next state of state s with input c.
// if the new state is available, it registers the state in the check array and returns the new state
// if the new state is not available, it returns the state and false
func (da *dArray) registerNextState(s int, c int) (t int, success bool) {
	// (s) --c--> (t)

	t = da.base(s) + c
	if da.check(t) == 0 {
		da.setCheck(t, s)
		//da.states[t].check = s
		return t, true
	} else if da.check(t) == s {
		return t, true
	}
	return t, false
}

// copyState copies the state s to t
func (da *dArray) copyState(t int, s int) {
	// (s) --c--> (t)
	da.setCheck(t, da.check(s))
	da.setBase(t, da.base(s))
}

// copyState move the state s to t. Update its children's check value
func (da *dArray) moveState(t int, s int) {
	// (s) --c--> (t)
	da.setCheck(t, da.check(s))
	da.setBase(t, da.base(s))

	for i := 0; i < len(da.states); i++ {
		if da.check(i) == s {
			da.setCheck(i, t)
		}
	}

	da.setBase(s, 0)
	da.setCheck(s, 0)
}