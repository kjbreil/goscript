package goscript

import (
	"strconv"
	"strings"
)

type StateText string

func (st *StateText) Float() float64 {
	f, err := strconv.ParseFloat(string(*st), 64)
	if err != nil {
		return 0
	}
	return f
}

func (st *StateText) Int() int {
	i, err := strconv.Atoi(string(*st))
	if err != nil {
		return 0
	}
	return i
}

func (st *StateText) String() string {
	return string(*st)
}

func (st *StateText) Equals(cmp string) bool {

	return strings.EqualFold(string(*st), cmp)
}
