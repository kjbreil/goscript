package goscript

import (
	"fmt"
	"time"
)

type Histories []*History

type History struct {
	Entity       string
	Domain       string
	DomainEntity string
	Start        time.Time
	End          time.Time
	States       []*State
}

func (gs *GoScript) GetHistory(start, end time.Time, entities ...string) (Histories, error) {
	hs, err := gs.ws.GetHistory(start, end, entities...)
	if err != nil {
		return nil, err
	}
	var gsHs Histories
	for _, h := range *hs {
		gsh := &History{
			Entity:       h.Entity,
			Domain:       h.Domain,
			DomainEntity: fmt.Sprintf("%s.%s", h.Domain, h.Entity),
			Start:        h.Start,
			End:          h.End,
			States:       make([]*State, 0, len(h.States)),
		}
		for _, s := range h.States {
			gsh.States = append(gsh.States, StateFromWS(s))
		}
		gsHs = append(gsHs, gsh)
	}
	return gsHs, nil
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
func (h *History) At(t time.Time) (before, after *State) {

	for _, s := range h.States {
		if s.LastChanged.After(t) {
			return before, s
		}
		before = s
	}
	return nil, nil
}
