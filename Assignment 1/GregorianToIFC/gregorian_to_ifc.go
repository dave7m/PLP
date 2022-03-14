package main

import (
	"fmt"
	"os"
)

/**
- file names are lowercase with underscores
- variable declaration: var %name %type = %value
- optionally inside functions: %name := %value
-
*/
func main() {
	// variables are written in camelCase
	var inProgress = true
	// there is no while-loop in Go
	for inProgress {
		inProgress = play(true)
	}
}

func play(firstRound bool) bool {
	if firstRound {
		fmt.Println("Please enter three positive integer numbers (year month day) seperated by one or more blank spaces or type quit.")
	}
	var answer string
	i, err := fmt.Scanln(&answer)
	if i != 1 {
		handleError(err)
	}
	for !isInputValid(answer) {
		fmt.Println("Invalid input, please try again!")
		i, err := fmt.Scanln(&answer)
		if i != 1 {
			handleError(err)
		}
	}
	return false
}

func isInputValid(answer string) bool {
	return false
}

func handleError(err error) {
	fmt.Println("Something went wrong while scanning the input: " + error.Error(err))
	os.Exit(0)
}
