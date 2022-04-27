package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// simplify input file, initialize StateTable and TransitionTable
func parse() error {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err1 := initializeStateTable(content)
	if err1 != nil {
		return err1
	}
	err2 := initializeTransitionTable(content)
	return err2 // nil or error
}

// reads the simplified file and creates entries in a dictionary. Key is the stateName, Value is a state.
func initializeStateTable(content []byte) error {

	// parse all states
	sStates := startStateRE.FindAll(content, -1)
	normalStates := normalStateRE.FindAll(content, -1)
	endStates := endStateRE.FindAll(content, -1)

	// get name and text of start state
	for i := 0; i < len(sStates); i++ {
		st := string(sStates[i])

		// get if state is an auto-forwarding state and the stateName and transitionBase
		stBase, stName, err := isAutoForwarding(st)
		if err != nil {
			return err
		}

		// get name and text of states
		// stName := strings.Trim(st[strings.Index(st, "*")+1:strings.Index(st, "{")], " ")

		stText := getStateText(st)
		stState := state{
			stateName:      stName,
			transitionBase: stBase,
			stateText:      stText,
			stateType:      startState,
		}
		// create entry at stateName
		err = setState(stState)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(normalStates); i++ {
		st := string(normalStates[i])
		// get name and text of states
		// get if state is an auto-forwarding state and the stateName and transitionBase
		stBase, stName, err := isAutoForwarding(st)
		if err != nil {
			return err
		}
		stText := getStateText(st)
		stState := state{
			stateName:      stName,
			transitionBase: stBase,
			stateText:      stText,
			stateType:      normalState,
		}
		// create entry at stateName
		err = setState(stState)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(endStates); i++ {
		st := string(endStates[i])

		// get name and text for end states
		// get if state is an auto-forwarding state and the stateName and transitionBase
		stBase, stName, err := isAutoForwarding(st)
		if err != nil {
			return err
		}
		stText := getStateText(st)
		stState := state{
			stateName:      stName,
			transitionBase: stBase,
			stateText:      stText,
			stateType:      endState,
		}
		// create entry at stateName
		err = setState(stState)
		if err != nil {
			return err
		}
	}
	return nil
}

// precondition: st is the string of a state (scanned by regex here)
// in order to work, it can only have one "!", the state name cannot have digits, and the time is only allowed if it is a
// auto-forwarding state
func isAutoForwarding(st string) (transitionBase, string, error) {
	// get the substring from (excluding) "@" to (excluding) "{" and eliminate white spaces
	var err error = nil
	subString := strings.Trim(st[strings.Index(st, "@")+1:strings.Index(st, "{")], " ")
	// the first two characters now determine the state type and if it is auto-forwarding
	autoForwarding := false
	if subString[0] == '!' || subString[1] == '!' {
		autoForwarding = true
	}
	s := strings.Trim(subString, "*+!")
	// the first digit splits s into the stateName and the transitionTime
	var stateName string
	var tT string
	// s must have the form abcd392, ab8c9d3 would lead to unpredicted behaviour (but we feed the function only with valid values anyway)
	for _, char := range s {
		if isAlpha(char) {
			stateName += string(char)
		} else {
			tT += string(char)
		}
	}
	tT = strings.Trim(tT, " ")
	var transitionTime uint64
	if len(tT) == 0 {
		transitionTime = 0
	} else {
		transitionTime, err = strconv.ParseUint(tT, 10, 64)
	}

	stBase := transitionBase{transitionTime: transitionTime}
	if autoForwarding {
		stBase.transitionType = autoForward
	} else {
		stBase.transitionType = defaultForward
	}

	return stBase, stateName, err
}

func isAlpha(charVariable rune) bool {
	return (charVariable >= 'a' && charVariable <= 'z') || (charVariable >= 'A' && charVariable <= 'Z')
}

func setState(s state) error {
	if _, notOk := StateTable[s.stateName]; notOk {
		return fmt.Errorf("multiple state definitions for state %s", s.stateName)
	}
	StateTable[s.stateName] = s
	return nil
}

func getStateText(s string) string {
	// everything should be printed, including newlines!
	return s[strings.Index(s, "{")+1 : len(s)-1]
}

// reads simplified file and adds entries to the dictionary TransitionTable. key is a struct with old state and the action,
// value is a struct with new state and transition text. If a transition is invalid, because a state is missing in the
// StateTable, the program stops with an error.
func initializeTransitionTable(content []byte) error {

	//get all transitions from content
	allTransitions := transitionRE.FindAll(content, -1)

	for i := 0; i < len(allTransitions); i++ {
		s := string(allTransitions[i])

		// get components from transition
		oldStateName := strings.Trim(s[strings.Index(s, ">")+1:strings.Index(s, "(")], " ")
		actionName := strings.Trim(s[strings.Index(s, "(")+1:strings.Index(s, ")")], " ")
		newStateName := strings.Trim(s[strings.Index(s, ")")+1:strings.Index(s, ":")], " ")
		tText := strings.Trim(s[strings.Index(s, ":")+1:], " ")

		// lookups in StateTable
		oldState, exists := StateTable[oldStateName]
		if !exists {
			return fmt.Errorf("An Error occured while parsing the file. See for more details: \n"+
				"Source state \"%s\" does not exist! -> State machine is invalid! \nProgram terminated", oldStateName)
		}
		newState, exists := StateTable[newStateName]
		if !exists {
			return fmt.Errorf("An Error occured while parsing the file. See for more details: \n"+
				"Destintation state \"%s\" does not exist! -> State machine is invalid! \nProgram terminated", newStateName)
		}

		// create new transition
		k := key{
			state:  oldState,
			action: actionName,
		}
		v := value{
			state:       newState,
			description: tText,
		}

		// create new entry in the TransitionTable at k
		if _, notOk := TransitionTable[k]; notOk {
			return fmt.Errorf("the key %q appears twice in the input file", k)
		}
		TransitionTable[k] = v
	}
	return nil
}
