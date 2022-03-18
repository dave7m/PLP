package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertGregorianToIFC(t *testing.T) {

	// for writing less
	assertions := assert.New(t)

	// our test cases have these parameters
	var cases = []struct {
		year          int64
		month         int64
		day           int64
		expectedYear  int64
		expectedMonth int64
		expectedDay   int64
	}{
		// cases from the assignment sheet
		{2021, 1, 1, 2021, 1, 1},
		{2021, 1, 28, 2021, 1, 28},
		{2021, 1, 29, 2021, 2, 1},
		{2021, 3, 1, 2021, 3, 4},
		{2020, 2, 29, 2020, 3, 4},
		{2020, 6, 17, 2020, -2, -2},
		{2021, 12, 30, 2021, 13, 28},
		{2020, 12, 30, 2020, 13, 28},
		{2021, 12, 31, 2021, -1, -1},
		{2020, 12, 31, 2020, -1, -1},
	}

	// iterate through test cases and assert equalities
	for _, c := range cases {
		actualYear, actualMonth, actualDay := convertGregorianToIFC(c.year, c.month, c.day)
		assertions.Equal(c.expectedYear, actualYear, fmt.Sprintf("When converting %d %d %d, "+
			"a wrong result was yielded: Year is not correct", c.year, c.month, c.day))
		assertions.Equal(c.expectedMonth, actualMonth, fmt.Sprintf("When converting %d %d %d, "+
			"a wrong result was yielded: Month is not correct", c.year, c.month, c.day))
		assertions.Equal(c.expectedDay, actualDay, fmt.Sprintf("When converting %d %d %d, "+
			"a wrong result was yielded: Day is not correct", c.year, c.month, c.day))
	}
}

func TestGetIFCString(t *testing.T) {

	assertions := assert.New(t)

	var cases = []struct {
		year           int64
		month          int64
		day            int64
		expectedString string
	}{
		// cases from the assignment sheet
		{2021, 1, 1, "2021 January 01 (Sunday)"},
		{2021, 1, 28, "2021 January 28 (Saturday)"},
		{2021, 2, 1, "2021 February 01 (Sunday)"},
		{2021, 3, 4, "2021 March 04 (Wednesday)"},
		{2020, 3, 4, "2020 March 04 (Wednesday)"},
		{2020, -2, -2, "Leap Day"},
		{2021, 13, 28, "2021 December 28 (Saturday)"},
		{2020, 13, 28, "2020 December 28 (Saturday)"},
		{2021, -1, -1, "Year Day"},
		{2020, -1, -1, "Year Day"},
	}

	for _, c := range cases {
		actualString := getIFCString(c.year, c.month, c.day)
		assertions.Equal(c.expectedString, actualString, fmt.Sprintf("Wrong result for %d %d %d: ", c.year, c.month, c.day))
	}
}

func TestIsGivenYearALeapYear(t *testing.T) {
	assertions := assert.New(t)
	var cases = []struct {
		year     int64
		expected bool
	}{
		{0, true},
		{100, false},
		{104, true},
		{400, true},
		{1000, false},
		{1600, true},
		{1960, true},
		{1996, true},
		{2000, true},
		{2003, false},
		{2004, true},
		{2005, false},
		{40000000, true},
	}

	for _, c := range cases {
		actual := isGivenYearALeapYear(c.year)
		var m string
		if c.expected {
			m = fmt.Sprintf("%d is a leap year, but returned %t\n", c.year, actual)
		} else {
			m = fmt.Sprintf("%d is not a leap year, but returned %t\n", c.year, actual)
		}
		assertions.Equal(c.expected, actual, m)
	}
}

func TestIsInputSyntaxValid(t *testing.T) {

	var cases = []struct {
		inputString string
		expected    bool
	}{
		{"2021 12 25", true},
		{"     2021      12           25       ", true},
		{"quit", true},
		{"         quit          ", true},
		{"000000000001 000000000001 000000000001", true},
		{"-5 12 21", false},
		{"1234b 12 21", false},
		{"2021 1225", false},
		{"", false},
		{"20 21 12 25", false},
	}
	for _, c := range cases {
		actual := isInputSyntaxValid(c.inputString)
		assert.Equal(t, c.expected, actual, fmt.Sprintf("%s should return %t, got %t", c.inputString, c.expected, actual))
	}
}

func TestGetYearMonthDay(t *testing.T) {
	assertions := assert.New(t)

	var cases = []struct {
		validInputString string
		expectedYear     int64
		expectedMonth    int64
		expectedDay      int64
	}{
		{"2021 12 25", 2021, 12, 25},
		{"    2021    12     25    ", 2021, 12, 25},
		{"000000000001 00000000001 00000000000001", 1, 1, 1},
		{"9999 12 21", 9999, 12, 21},
	}
	for _, c := range cases {
		actualYear, actualMonth, actualDay := getYearMonthDay(c.validInputString)
		assertions.Equal(c.expectedYear, actualYear)
		assertions.Equal(c.expectedMonth, actualMonth)
		assertions.Equal(c.expectedDay, actualDay)
	}
}

func TestDoesDateExist(t *testing.T) {
	assertions := assert.New(t)

	var cases = []struct {
		year     int64
		month    int64
		day      int64
		expected bool
	}{
		{2021, 12, 25, true},
		{2012, 2, 29, true},
		{2100, 2, 29, false},
		{-1, 3, 16, false},
		{2000, -1, 1, false},
		{2000, 1, -1, false},
		{2000, 13, 1, false},
		{2000, 1, 32, false},
		{2022, 4, 31, false},
	}

	for _, c := range cases {
		actual := doesDateExist(c.year, c.month, c.day)
		assertions.Equal(c.expected, actual)
	}

}
