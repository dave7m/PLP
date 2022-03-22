package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	createFile      bool
	path            string
	reader          = bufio.NewReader(os.Stdin)
	TransitionTable = make(map[key]value)
	StateTable      = make(map[string]state)
)

type key struct {
	state  state
	action string
}
type value struct {
	state       state
	description string
}

// this is a new type
type stateType int

const (
	startState stateType = iota // iota 'enumerates' the states -> startState has val 0
	normalState
	endState
)

type state struct {
	stateType stateType
	stateName string
	stateText string
}

// init is called before main()
func init() {
	// create a custom flag accepting a string and saving it to path.
	flag.StringVar(&path, "path", "", "path of machine file")
	flag.BoolVar(&createFile, "c", false, "bool: create compressed file without comments")
}

func main() {
	flag.Parse()

	if len(path) == 0 {
		fmt.Println("Don't forget to specify the path to the .machine file!")
		flag.PrintDefaults()
		os.Exit(1)
	}

	executeStateMachineConfiguration()
}

func executeStateMachineConfiguration() {

	// parse the file and save it as data structures
	err = parse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(StateTable)

	var s state
	// validate internal representation of the State Machine and returns start state
	s, err = isValid(StateTable, TransitionTable)
	if err != nil {
		log.Fatal(err)
	}

	// flag
	if createFile {
		serialize()
	}

	// run internal State Machine

	for s.stateType != endState {
		checkSinkState(s)
		fmt.Println(s.stateText)
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		s = doTransition(s, input)
	}
}

// we are in a sink state, if there is no transition out of the state
func checkSinkState(s state) {
	counter := 0
	// check in the TransitionTable if there is a transition
	for k := range TransitionTable {
		if k.state == s {
			counter++
		}
	}
	if counter == 0 {
		fmt.Printf("\"%s\" is a sink state: there are no transitions out of this state, even though \"%s\" is not an end state."+
			"The program is terminated.", s.stateName, s.stateName)
		os.Exit(3)
	}
}

// changes the current state to the state specified in the transition
func doTransition(oldState state, input string) state {
	input = strings.Trim(input, "\r\n")
	t, found := TransitionTable[key{state: oldState, action: input}]
	if !found {
		fmt.Print("Invalid input")
		return oldState
	}
	fmt.Print(t.description)
	return t.state
}
