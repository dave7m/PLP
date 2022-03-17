package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/**
- file names are lowercase with underscores
- variable declaration: var %name %type = %value
- optionally inside functions: %name := %value
-
*/
// global reader
var reader = bufio.NewReader(os.Stdin)

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
	answer, err := reader.ReadString('\n')

	if err != nil {
		handleError(err)
	}

	for !isInputSyntaxValid(answer) {
		fmt.Println("Invalid input, please try again!")
		answer, err = reader.ReadString('\n')
		if err != nil {
			handleError(err)
		}
	}
	if wantsToQuit(answer) {
		return false
	}
	year, month, day := getYearMonthDay(answer)
	fmt.Printf("year: %d, month: %d, day: %d\n", year, month, day)
	return true
}

func getYearMonthDay(answer string) (int64, int64, int64) {
	// this creates an array with the words seperated by a whitespace
	// similar to split() int python
	wordArray := strings.Fields(answer)
	if len(wordArray) != 3 {
		fmt.Println("Something went horribly wrong, therefore, the program is terminated")
	}

	// convert the strings in the array to
	year, _ := strconv.ParseInt(wordArray[0], 10, 64)
	month, _ := strconv.ParseInt(wordArray[1], 10, 64)
	day, _ := strconv.ParseInt(wordArray[2], 10, 64)
	return year, month, day
}

func wantsToQuit(answer string) bool {
	res, err := regexp.MatchString("(\\s*)quit(\\s*)", answer)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return res
}

func isInputSyntaxValid(answer string) bool {
	fmt.Println(answer)
	// first compile the regular expression function, this makes it faster at run time
	re, err := regexp.Compile("[0-9]{1,4}(\\s*)([0-9]{1,2}(\\s*)){2}|quit(\\s*)")
	if err != nil {
		fmt.Println(err)
	}
	return re.MatchString(answer)
}

func handleError(err error) {
	fmt.Println("Something went wrong while scanning the input: " + error.Error(err))
	os.Exit(0)
}
