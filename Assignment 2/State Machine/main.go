package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	path            string
	machinePath     = "machine.machine"
	reader          = bufio.NewReader(os.Stdin)
	TransitionTable = make(map[key]value)
	StateTable      = make(map[string]state)
	Start           state
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
	startState stateType = iota
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
	// todo get default value for car.machine path or make it required
	flag.StringVar(&path, "path", "", "absolute path of machine file")
}

func main() {
	flag.Parse()

	if len(path) == 0 {
		fmt.Println("Usage: main.go -path")
		flag.PrintDefaults()
		os.Exit(1)
	}
	// fmt.Println("path: ", path)

	executeMachineAt(path)
}

func executeMachineAt(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// prepend newline to content
	contentToParse := append(append([]byte{10}, content...), 10)
	doParse(contentToParse)
	serialize()

	// fmt.Printf("%q\n", TransitionTable)
	/*
		s := Start
		for s.stateType != endState {
			fmt.Println(s.stateText)
			input, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			s = doTransition(s, input)
		}*/
}

func doTransition(oldState state, input string) state {
	input = strings.Trim(input, "\r\n")
	t, found := TransitionTable[key{state: oldState, action: input}]
	// fmt.Printf("state: %s, action: %s", oldState.stateName, input)
	if !found {
		fmt.Print("Invalid input")
		return oldState
	}
	fmt.Print(t.description)
	return t.state
}
