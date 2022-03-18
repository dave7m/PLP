package main

import (
	"bufio"
	"errors"
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
var firstRound = true

func main() {
	// variables are written in camelCase
	var inProgress = true
	// there is no while-loop in Go
	for inProgress {
		inProgress = play()
	}
}

func play() bool {
	if firstRound {
		fmt.Println("Please enter three positive integer numbers (year month day) seperated by one or more blank spaces or type quit.")
		firstRound = false
	}

	dateIsValid := false
	var year, month, day int64

	// only accept valid dates
	for !dateIsValid {

		answer, err := reader.ReadString('\n')
		if err != nil {
			handleError(err)
		}

		if !isInputSyntaxValid(answer) {
			fmt.Println("Invalid input, please try again!")
			continue
		}
		if wantsToQuit(answer) {
			return false
		}
		year, month, day = getYearMonthDay(answer)

		if !doesDateExist(year, month, day) {
			fmt.Println("Input date does not exist, please try again!")
			continue
		}

		// all sanity checks passed, thus the date must exit and can be processed
		dateIsValid = true
	}

	// convert the valid date to the IFC date format
	newYear, newMonth, newDay := convertGregorianToIFC(year, month, day)
	IFCString := getIFCString(newYear, newMonth, newDay)
	fmt.Println(IFCString)
	return true
}

// accepts three int (year, month, day) and returns a corresponding string
func getIFCString(year int64, month int64, day int64) string {
	var res string
	// per definition, Year Day has month and day as -1, Leap Day -2 each
	if month == -1 && day == -1 {
		res = fmt.Sprintf("Year Day\n")
	} else if month == -2 && day == -2 {
		res = fmt.Sprintf("Leap Day\n")
	} else {
		months := [13]string{"January", "February", "March", "April", "Mai", "June", "Sol", "July", "August",
			"September", "October", "November", "December"}
		weekdays := [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
		res = fmt.Sprintf("%04d %s %02d (%s)\n", year, months[month-1], day, weekdays[(day-1)%7])
	}
	return res
}

/**
returned values are analogue to gregorian calendar, except we return year, -1, -1 for year day and year, -2, -2
for leap day (which must be taken care of by the caller)
*/
func convertGregorianToIFC(year int64, month int64, day int64) (int64, int64, int64) {
	isLeapYear := isGivenYearALeapYear(year)
	var leapDay int64 = 0
	if isLeapYear {
		leapDay = 1
	}
	daysPerMonth := [12]int64{31, 28 + leapDay, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	var nThDay int64 = 0
	// add all preceding months as full
	for i := 0; int64(i) < month-1; i++ {
		nThDay += daysPerMonth[i]
	}
	// add the days of the current month
	nThDay += day

	// check for year day and leap day
	if month == 12 && day == 31 {
		return year, -1, -1
	}
	if isLeapYear && month == 6 && day == 17 {
		return year, -2, -2
	}

	// for calculation, we subtract one day from nThDay if it is after the leap day
	// 6 * 28 because the leap day is inserted after the sixth month
	// nThDay might be a misleading name, as we manipulate it to get the day of the month, so first it is the
	// n-th day of the year, and after the manipulation it is the n-th day of a month
	if isLeapYear && nThDay > 6*28 {
		nThDay--
	}

	// calculate new month
	var newMonth int64 = 1
	for ; nThDay > 28; nThDay -= 28 {
		newMonth += 1
	}

	return year, newMonth, nThDay

}

func isGivenYearALeapYear(year int64) bool {
	return year%4 == 0 && year%100 != 0 || year%400 == 0
}

func doesDateExist(year int64, month int64, day int64) bool {
	isLeapYear := isGivenYearALeapYear(year)
	var leapDay int64 = 0
	if isLeapYear {
		leapDay = 1
	}
	daysPerMonth := [12]int64{31, 28 + leapDay, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	// because it is first checked if month is <= 12 there can't be an index error in the last comparison
	return year >= 0 && 1 <= month && month <= 12 && day >= 1 && day <= daysPerMonth[month-1]
}

// out of a valid string (which is not 'quit'), extract the year, month and day and return it as int64
func getYearMonthDay(answer string) (int64, int64, int64) {
	// this creates an array with the words seperated by a whitespace
	// similar to split() int python
	wordArray := strings.Fields(answer)
	if len(wordArray) != 3 {
		err := errors.New("index Error: tried to access element outside of array")
		handleError(err)
	}

	// convert the strings in the array to int64
	year, _ := strconv.ParseInt(wordArray[0], 10, 64)
	month, _ := strconv.ParseInt(wordArray[1], 10, 64)
	day, _ := strconv.ParseInt(wordArray[2], 10, 64)
	return year, month, day
}

func wantsToQuit(answer string) bool {
	res, err := regexp.MatchString("(\\s*)quit(\\s*)", answer)
	if err != nil {
		handleError(err)
	}
	return res
}

func isInputSyntaxValid(answer string) bool {

	// first compile the regular expression function, this makes it faster at run time
	re, err := regexp.Compile("^(\\s*)(0*[0-9]{1,4}(\\s+)(0*[0-9]{1,2}(\\s+))(0*[0-9]{1,2}(\\s*))|quit(\\s*))$")
	if err != nil {
		handleError(err)
	}
	return re.MatchString(answer)
}

// print error and terminate program
func handleError(err error) {
	fmt.Println("Something went wrong while scanning the input: " + error.Error(err))
	os.Exit(0)
}
