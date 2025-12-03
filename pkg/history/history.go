// Package history provides types and methods for managing entity state history.
package history

import (
	"time"

	"github.com/kjbreil/goscript/pkg/state"
)

// Histories is a collection of History records.
type Histories []*History

// History represents the state history for a single entity over a time range.
type History struct {
	Entity       string
	Domain       string
	DomainEntity string
	Start        time.Time
	End          time.Time
	States       []*state.State
}

// GetHistories represents a request for entity histories over a time range.
type GetHistories struct {
	Start    time.Time
	End      time.Time
	Entities []string
}

// Get retrieves a History for the specified domain and entity.
func (hs *Histories) Get(domain, entity string) *History {
	for _, h := range *hs {
		if h.Entity == entity && h.Domain == domain {
			return h
		}
	}
	return nil
}

// At returns the two states surrounding the given time.
func (h *History) At(t time.Time) (before, after *state.State) {
	for _, s := range h.States {
		if s.LastChanged.After(t) {
			return before, s
		}
		before = s
	}
	return nil, nil
}
