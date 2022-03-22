package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// simplify input file, initialize StateTable and TransitionTable. Because of the order of the operations, the tables
// must be correct and not scanned for invalidity.
func parse(content []byte) {
	// note: the simplifyStateMachine is not at all necessary, we could do just fine without them, but it produces a file
	// which is ordered and without comments and invalid state / transition definitions, which is nice
	simplifyStateMachine(content)
	initializeStateTable()
	// fmt.Printf("%q\n", StateTable)
	initializeTransitionTable()
}

// reads the simplified file and creates entries in a dictionary. Key is the stateName, Value is a state.
func initializeStateTable() {
	content, err := ioutil.ReadFile(machinePath)
	if err != nil {
		panic(err)
	}

	sState := string(startStateRE.Find(content))
	normalStates := normalStateRE.FindAll(content, -1)
	endStates := endStateRE.FindAll(content, -1)

	startStateName := strings.Trim(sState[strings.Index(sState, "*")+1:strings.Index(sState, "{")], " ")
	startStateText := getStateText(sState)
	// fmt.Println(startStateName)

	// initialize global Start state.
	Start = state{
		stateName: startStateName,
		stateText: startStateText,
		stateType: startState,
	}
	// create entry at stateName
	StateTable[startStateName] = Start

	for i := 0; i < len(normalStates); i++ {
		st := string(normalStates[i])
		stName := strings.Trim(st[strings.Index(st, "@")+1:strings.Index(st, "{")], " ")
		// fmt.Printf("normal state name: %s\n", stName)
		stText := getStateText(st)
		stState := state{
			stateName: stName,
			stateText: stText,
			stateType: normalState,
		}
		// create entry at stateName
		StateTable[stName] = stState
	}

	for i := 0; i < len(endStates); i++ {
		st := string(endStates[i])
		stName := strings.Trim(st[strings.Index(st, "+")+1:strings.Index(st, "{")], " ")
		// fmt.Printf("end state name: %s\n", stName)
		stText := getStateText(st)
		stState := state{
			stateName: stName,
			stateText: stText,
			stateType: endState,
		}
		// create entry at stateName
		StateTable[stName] = stState
	}
}

func getStateText(s string) string {
	// everything should be printed, including newlines!
	a := s[strings.Index(s, "{")+1 : len(s)-1]
	// fmt.Printf("%q\n", a)
	return a
}

// reads simplified file and adds entries to the dictionary TransitionTable. key is a struct with old state and the action,
// value is a struct with new state and transition text. If a transition is invalid, because a state is missing in the
// StateTable, the parsing stops, and we panic.
func initializeTransitionTable() {
	content, err := ioutil.ReadFile(machinePath)
	if err != nil {
		panic(err)
	}
	//fmt.Println(content)
	//transitionRE, _ := regexp.Compile("(?m)^>\\s*\\w*\\s*\\([\\s*\\w]+\\)\\s*\\w*\\s*:\\s*[^\\n]*$")
	allTransitions := transitionRE.FindAll(content, -1)
	for i := 0; i < len(allTransitions); i++ {
		s := string(allTransitions[i])
		// fmt.Printf("transition %d: %q \n", i, s)
		oldStateName := strings.Trim(s[strings.Index(s, ">")+1:strings.Index(s, "(")], " ")
		// fmt.Printf("old state name: %q\n", oldStateName)

		actionName := strings.Trim(s[strings.Index(s, "(")+1:strings.Index(s, ")")], " ")
		// fmt.Println(actionName)
		newStateName := strings.Trim(s[strings.Index(s, ")")+1:strings.Index(s, ":")], " ")
		// fmt.Println(newStateName)
		tText := strings.Trim(s[strings.Index(s, ":")+1:], " ")

		// lookup in StateTable
		oldState, exists := StateTable[oldStateName]
		if !exists {
			log.Fatal(fmt.Sprintf("An Error occured while parsing the file. See for more details: \n"+
				"Source state \"%s\" does not exist! -> State machine is invalid! \nProgram terminated", newStateName))
		}
		newState, exists := StateTable[newStateName]
		if !exists {
			log.Fatal(fmt.Sprintf("An Error occured while parsing the file. See for more details: \n"+
				"Destintation state \"%s\" does not exist! -> State machine is invalid! \nProgram terminated", newStateName))
		}
		k := key{
			state:  oldState,
			action: actionName,
		}
		v := value{
			state:       newState,
			description: tText,
		}
		// create new entry in the TransitionTable at k
		TransitionTable[k] = v
	}
}

// puts start state to the top of a new file, gets rid of comments (only selects valid syntax)
// catches: null or multiple start states
func simplifyStateMachine(content []byte) {
	// fmt.Println(string(content))

	// first, create new file (machine.machine)
	file, err := os.Create(machinePath)
	if err != nil {
		panic(err)
	}

	// match all start states and write them to file
	allStartStates := startStateRE.FindAll(content, -1)
	//fmt.Printf("%q\n", allStartStates)
	if len(allStartStates) != 1 {
		panic(fmt.Sprintf("Expected exactly one start state, got %d!", len(allStartStates)))
	}
	_, err = file.WriteString(string(allStartStates[0]) + "\n")
	if err != nil {
		panic(err)
	}

	// match all normal states and write them to file
	allNormalStates := normalStateRE.FindAll(content, -1)
	//fmt.Printf("%q\n", allStartStates)
	for i := 0; i < len(allNormalStates); i++ {
		_, err = file.WriteString(string(allNormalStates[i]) + "\n")
		if err != nil {
			panic(err)
		}
	}

	// match all end states and write them to file
	allEndStates := endStateRE.FindAll(content, -1)
	//fmt.Printf("%q\n", allStartStates)
	for i := 0; i < len(allEndStates); i++ {
		_, err = file.WriteString(string(allEndStates[i]) + "\n")
		if err != nil {
			panic(err)
		}
	}

	// match all transitions and write them to file
	allTransitions := transitionRE.FindAll(content, -1)
	for i := 0; i < len(allTransitions); i++ {
		_, err = file.WriteString(string(allTransitions[i]) + "\n")
		if err != nil {
			panic(err)
		}
	}
}
