package core

import (
	"fmt"
	"time"

	"github.com/kjbreil/goscript/pkg/history"
	"github.com/kjbreil/goscript/pkg/state"
)

func (gs *GoScript) GetHistory(start, end time.Time, entities ...string) (history.Histories, error) {
	hs, err := gs.ws.GetHistory(start, end, entities...)
	if err != nil {
		return nil, fmt.Errorf("failed to get history from websocket: %w", err)
	}
	var gsHs history.Histories
	for _, h := range *hs {
		gsh := &history.History{
			Entity:       h.Entity,
			Domain:       h.Domain,
			DomainEntity: fmt.Sprintf("%s.%s", h.Domain, h.Entity),
			Start:        h.Start,
			End:          h.End,
			States:       make([]*state.State, 0, len(h.States)),
		}
		for _, s := range h.States {
			gsh.States = append(gsh.States, state.StateFromWS(s))
		}
		gsHs = append(gsHs, gsh)
	}
	return gsHs, nil
}
