package mptrie

// double array item
type state struct {
	base, check int
	end         bool
}

type dArray struct {
	states []state
}

func newDArray(size int) *dArray {
	return &dArray{states: make([]state, size)}
}

func (da *dArray) base(s int) int {
	return da.states[s].base
}

func (da *dArray) setBase(s int, b int) {
	da.states[s].base = b
}

func (da *dArray) check(t int) int {
	return da.states[t].check
}

func (da *dArray) setCheck(t int, s int) {
	da.states[t].check = s
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
