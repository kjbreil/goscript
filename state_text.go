package goscript

import (
	"strconv"
	"strings"
)

type StateText string

// Float helps to convert state text to float if it is known to be a float. Will return 0 if it is not a float.
func (st *StateText) Float() float64 {
	f, err := strconv.ParseFloat(string(*st), 64)
	if err != nil {
		return 0
	}
	return f
}

// Int returns the int representation of the StateText, or 0 if the StateText cannot be converted to an int.
// If the StateText can be represented as a float, it is rounded down.
func (st *StateText) Int() int {
	i, err := strconv.Atoi(string(*st))
	if err != nil {
		f, err := strconv.ParseFloat(string(*st), 64)
		if err != nil {
			return 0
		}
		return int(f)
	}
	return i
}

// String returns the string representation of the StateText.
func (st *StateText) String() string {
	return string(*st)
}

// Equals returns true if the StateText is equal to the specified string, ignoring case.
func (st *StateText) Equals(cmp string) bool {
	return strings.EqualFold(string(*st), cmp)
}
