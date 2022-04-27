package main

import (
	"os"
	"strconv"
)

var newFile *os.File
var err error

func serialize() {
	// create new file
	newFile, err = os.Create("serialized_machine.machine")
	if err != nil {
		panic(err)
	}
	defer func(newFile *os.File) {
		err = newFile.Close()
		if err != nil {
			panic(err)
		}
	}(newFile)

	// write file
	writeStates()
	writeTransitions()

}

func writeStates() {
	// iterate through StateTable
	for _, st := range StateTable {
		var typ = "@"
		var t = ""
		switch st.stateType {
		case startState:
			typ += "*"
			break
		case endState:
			typ += "+"
		}
		if st.transitionBase.transitionType == autoForward {
			typ += "!"
			t = strconv.FormatUint(st.transitionBase.transitionTime, 10)
		}

		stateString := typ + st.stateName + t + "{" + st.stateText + "}\n"
		_, err = newFile.WriteString(stateString)
		if err != nil {
			panic(err)
		}

	}
}

func writeTransitions() {
	for k, v := range TransitionTable {
		transitionString := "> " + k.state.stateName + " (" + k.action + ") " + v.state.stateName + ": " + v.description + "\n"

		_, err = newFile.WriteString(transitionString)
		if err != nil {
			panic(err)
		}
	}
}
