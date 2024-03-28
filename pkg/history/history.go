package history

import (
	"github.com/kjbreil/goscript/pkg/state"
	"time"
)

type Histories []*History

type History struct {
	Entity       string
	Domain       string
	DomainEntity string
	Start        time.Time
	End          time.Time
	States       []*state.State
}

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
