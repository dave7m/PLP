package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var machinePath = "machine.machine"
var reader = bufio.NewReader(os.Stdin)
var (
	path string
)
var transitionTable = make(map[key]value)
var stateTable = make(map[string]state)

type key struct {
	state  state
	action string
}
type value struct {
	state       state
	description string
}

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

type transition struct {
	oldState       state
	action         string
	newState       state
	transitionText string
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
	fmt.Println("path: ", path)

	executeMachineAt(path)
}

func executeMachineAt(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// prepend newline to content

	simplifyStateMachine(append(append([]byte{10}, content...), 10))
	initializeStateTable()
	fmt.Printf("%q\n", stateTable)

	initializeTransitionTable()
	println(transitionTable)
	s := state{stateName: "hello", stateType: startState, stateText: "hello"}
	for s.stateType != endState {
		_, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		//s = doTransition(s, input)
	}
}

func initializeStateTable() {
	content, err := ioutil.ReadFile(machinePath)
	if err != nil {
		panic(err)
	}
	startStateRE, _ := regexp.Compile("\n@\\*[a-zA-Z0-9]*{[^}]*}")
	normalStateRE, _ := regexp.Compile("\n@[a-zA-Z0-9]*{[^}]*}")
	endStateRE, _ := regexp.Compile("\n@\\+[a-zA-Z0-9]*{[^}]*}")

	sState := string(startStateRE.Find(content))
	normalStates := normalStateRE.FindAll(content, -1)
	endStates := endStateRE.FindAll(content, -1)

	startStateName := strings.Trim(sState[strings.Index(sState, "*")+1:strings.Index(sState, "{")], " ")
	startStateText := getStateText(sState)
	fmt.Println(startStateName)
	startStateState := state{
		stateName: startStateName,
		stateText: startStateText,
		stateType: startState,
	}
	stateTable[startStateName] = startStateState

	for i := 0; i < len(normalStates); i++ {
		st := string(normalStates[i])
		stName := strings.Trim(st[strings.Index(st, "@")+1:strings.Index(st, "{")], " ")
		fmt.Printf("normal state name: %s\n", stName)
		stText := getStateText(st)
		stState := state{
			stateName: stName,
			stateText: stText,
			stateType: normalState,
		}
		stateTable[stName] = stState
	}

	for i := 0; i < len(endStates); i++ {
		st := string(endStates[i])
		stName := strings.Trim(st[strings.Index(st, "+")+1:strings.Index(st, "{")], " ")
		fmt.Printf("end state name: %s\n", stName)
		stText := getStateText(st)
		stState := state{
			stateName: stName,
			stateText: stText,
			stateType: endState,
		}
		stateTable[stName] = stState
	}
}

func getStateText(s string) string {
	return strings.Trim(s[strings.Index(s, "{")+1:len(s)-1], " ")
}

func initializeTransitionTable() {
	content, err := ioutil.ReadFile(machinePath)
	if err != nil {
		panic(err)
	}
	//fmt.Println(content)
	transitionRE, _ := regexp.Compile("(?m)^>\\s*\\w*\\s*\\(\\s*\\w*\\s*\\)\\s*\\w*\\s*:\\s*[^\\n]*$")
	allTransitions := transitionRE.FindAll(content, -1)
	for i := 0; i < len(allTransitions); i++ {
		s := string(allTransitions[i])
		// fmt.Printf("transition %d: %q \n", i, s)
		oldStateName := strings.Trim(s[strings.Index(s, ">")+1:strings.Index(s, "(")], " ")
		// fmt.Printf("old state name: %q\n", oldStateName)

		actionName := strings.Trim(s[strings.Index(s, "(")+1:strings.Index(s, ")")], " ")
		// fmt.Println(actionName)
		newStateName := strings.Trim(s[strings.Index(s, ")")+1:strings.Index(s, ":")], " ")
		fmt.Println(newStateName)
		tText := strings.Trim(s[strings.Index(s, ":")+1:], " ")

		oldState, exists := stateTable[oldStateName]
		if !exists {
			panic("old state seems not to exist")
		}
		newState, exists := stateTable[newStateName]
		if !exists {
			panic(fmt.Sprintf("newState: %s seems not to exist in state table!", newStateName))
		}
		k := key{
			state:  oldState,
			action: actionName,
		}
		v := value{
			state:       newState,
			description: tText,
		}
		transitionTable[k] = v
	}
}

/*
func doTransition(oldState state, input string) state {
	t := getTransition(oldState, input)
}

func getTransition(oldState state, input string) transition {

}*/

// puts start state to the top of a new file, gets rid of comments (only selects valid syntax)
// catches: null or multiple start states
func simplifyStateMachine(content []byte) {
	// fmt.Println(string(content))

	file, err := os.Create(machinePath)
	if err != nil {
		panic(err)
	}

	startStateRE, _ := regexp.Compile("\n@\\*[a-zA-Z0-9]*{[^}]*}")
	allStartStates := startStateRE.FindAll(content, -1)
	//fmt.Printf("%q\n", allStartStates)
	if len(allStartStates) != 1 {
		panic("found null or more than one start state!")
	}
	_, err = file.WriteString(string(allStartStates[0]) + "\n")
	if err != nil {
		panic(err)
	}

	normalStateRE, _ := regexp.Compile("\n@[a-zA-Z0-9]*{[^}]*}")
	allNormalStates := normalStateRE.FindAll(content, -1)
	//fmt.Printf("%q\n", allStartStates)
	for i := 0; i < len(allNormalStates); i++ {
		_, err = file.WriteString(string(allNormalStates[i]) + "\n")
		if err != nil {
			panic(err)
		}
	}

	endStateRE, _ := regexp.Compile("\n@\\+[a-zA-Z0-9]*{[^}]*}")
	allEndStates := endStateRE.FindAll(content, -1)
	//fmt.Printf("%q\n", allStartStates)
	for i := 0; i < len(allEndStates); i++ {
		_, err = file.WriteString(string(allEndStates[i]) + "\n")
		if err != nil {
			panic(err)
		}
	}

	transitionRE, _ := regexp.Compile("(?m)>\\s*\\w*\\s*\\(\\s*\\w*\\s*\\)\\s*\\w*\\s*:\\s*[\\s\\S]*$")
	allTransitions := transitionRE.FindAll(content, -1)
	for i := 0; i < len(allTransitions); i++ {
		_, err = file.WriteString(string(allTransitions[i]))
		if err != nil {
			panic(err)
		}
	}
}
