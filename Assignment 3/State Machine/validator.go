package main

import (
	"fmt"
)

func isValid(originalStateTable map[string]state, originalTransitionTable map[key]value) (state, error) {
	stateTable := make(map[string]state)
	transitionTable := make(map[key]value)

	// first copy maps in order not to manipulate original content
	for k, v := range originalStateTable {
		stateTable[k] = v
	}
	for k, v := range originalTransitionTable {
		transitionTable[k] = v
	}

	// optional:
	start, err1 := assertOneStart(stateTable)
	if err1 != nil {
		return state{}, err1
	}

	// optional:
	err2 := assertMoreThanOneEnd(stateTable)
	if err2 != nil {
		return state{}, err2
	}
	// check if end states don't are auto-forwarding
	err3 := assertEndStatesAreNotForwarding(stateTable)
	if err3 != nil {
		return state{}, err3
	}

	// check if each auto-forwarding state only has one transition
	err4 := assertOneCorrectTransitionPerForwardingState(stateTable, transitionTable)
	if err4 != nil {
		return state{}, err4
	}
	isReachable := solveEmptinessProblem(stateTable, transitionTable)
	if !isReachable {
		return state{}, fmt.Errorf("state machine is empty, that is, there is no run from the start state to any end. " +
			"This might be because of an invalid syntax or missing transisions")
	}

	return start, nil
}

func assertOneCorrectTransitionPerForwardingState(stTable map[string]state, trTable map[key]value) error {
	for _, st := range stTable {
		// get an auto-forwarding state
		if st.transitionBase.transitionType == autoForward {
			counter := 0
			for k, v := range trTable {
				// autoForward state reaches
				if st == k.state {
					counter++
				}
				// autoForward state is being reached
				if st == v.state {
					action := k.action
					_, found := trTable[key{state: st, action: action}]
					if !found {
						return fmt.Errorf("State \"" + st.stateName + "\" which is auto-forwarding is reached by " +
							"action \"" + action + "\" has no transition to another state with the same action.")
					}
				}
			}
			// see if it has exactly one transition
			if counter != 1 {
				return fmt.Errorf("State \""+st.stateName+"\" must have exactly one transition because it is auto "+
					"forwarding but %d were given.", counter)
			}

		}
	}
	return nil
}

func assertEndStatesAreNotForwarding(table map[string]state) error {
	for _, st := range table {
		if st.stateType == endState && st.transitionBase.transitionType == autoForward {
			return fmt.Errorf("End State \"" + st.stateName + "\" is an auto forwarding state, but should not be one!")
		}
	}
	return nil
}

func solveEmptinessProblem(stateTable map[string]state, transitionTable map[key]value) bool {
	// create a target state finalState
	finalState := state{
		stateName: "finalState",
		stateType: endState,
		stateText: "",
	}
	// stateTable[finalState.stateName] = finalState // should not change old stateTable?
	var start state
	for _, st := range stateTable {
		if st.stateType == startState {
			start = st
		}
	}

	// insert transitions from all end states to finalState
	for _, v := range stateTable {
		if v.stateType == endState {
			nK := key{
				state:  v,
				action: "added link",
			}
			nV := value{
				state:       finalState,
				description: "This link is used to solve the Emptiness Problem",
			}
			transitionTable[nK] = nV
		}
	}
	return IDDFS(start, finalState, transitionTable, len(stateTable))
	// test if there is a path from the start state to finalState
}

// IDDFS = Iterative Deepening DFS
// note: BFS would probably be more suited, but I wanted to try out IDDFS. IDDFS has slower performance, but uses less space
func IDDFS(root state, target state, table map[key]value, maxDepth int) bool {
	for limit := 0; limit <= 1+maxDepth; limit++ { // >= ?
		if DLS(root, target, table, limit) {
			return true
		}
	}
	return false
}

func DLS(node state, target state, table map[key]value, limit int) bool {
	if node == target {
		return true
	}

	// stop recursion if maximum depth is reached
	if limit <= 0 {
		return false
	}

	// for each adjacent state (to node) recurse
	for k, v := range table {
		if k.state == node {
			adj := v.state
			if DLS(adj, target, table, limit-1) {
				return true
			}
		}
	}
	return false
}

func assertOneStart(table map[string]state) (state, error) {
	counter := 0
	err = nil
	var s state
	for _, v := range table {
		if v.stateType == startState {
			s = v
			counter++
		}
	}
	if counter != 1 {
		return s, fmt.Errorf("A state machine must have exactly one start state, but %d were given", counter)
	}
	return s, nil
}

func assertMoreThanOneEnd(table map[string]state) error {
	counter := 0
	err = nil
	for _, v := range table {
		if v.stateType == endState {
			counter++
		}
	}
	if counter < 1 {
		return fmt.Errorf("A state machine must have one or more end states, but %d were given", counter)
	}
	return nil
}
