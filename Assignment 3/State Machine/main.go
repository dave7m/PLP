package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/inancgumus/screen"
	"log"
	"os"
	"strings"
	"time"
)

var (
	input           string
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
type transitionType int

const (
	startState stateType = iota // iota 'enumerates' the states -> startState has val 0
	normalState
	endState
)

const (
	autoForward transitionType = iota
	defaultForward
)

type transitionBase struct {
	transitionType transitionType
	transitionTime uint64
}

// a state has now a transition type
type state struct {
	stateType      stateType
	transitionBase transitionBase
	stateName      string
	stateText      string
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

	// clear screen
	screen.Clear()
	screen.MoveTopLeft()

	// run internal State Machine
	for s.stateType != endState {
		checkSinkState(s)
		fmt.Println(s.stateText)

		if s.transitionBase.transitionType == autoForward {

			// sleep for t milliseconds
			t := time.Duration(int(s.transitionBase.transitionTime))
			time.Sleep(time.Millisecond * t)
			// input is stored globally, so default input for the next state is the same input as before (input = action)

		} else {
			input, err = reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			input = strings.Trim(input, "\r\n")

			if input == "quit" {
				os.Exit(0)
			}
		}
		// before doing the transition, we clear the screen (task A), for this, we use an external libraries to better
		// cope with different operating systems (I think it does not work with bash): github.com/inancgumus/screen
		screen.Clear()
		screen.MoveTopLeft()

		s = doTransition(s, input)
	}
	// we reached the end state, but want to show the stateText anyway
	fmt.Println(s.stateText)
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
	// because of issue (see Readme), redirect to doAutoTransition here
	if oldState.transitionBase.transitionType == autoForward {
		return doAutoTransition(oldState)
	}

	t, found := TransitionTable[key{state: oldState, action: input}]
	a := input[:len(input)-1]
	for !found && len(a) != 0 {
		// Algorithm: strip action by one character, append an asterisk manually and see if it matches a wildcard action
		t, found = TransitionTable[key{state: oldState, action: a + "*"}]
		a = a[:len(a)-1]
	}
	if !found {
		fmt.Println("Invalid input!")
		return oldState
	}
	fmt.Println(t.description)
	return t.state
}

func doAutoTransition(oldState state) state {
	for k, v := range TransitionTable {
		if k.state == oldState {
			fmt.Println(v.description)
			return v.state
		}
	}
	fmt.Println("If you can read this on the terminal, the validator failed analyzing the auto-forwarding states")
	os.Exit(3)
	return oldState
}
