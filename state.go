package goscript

import (
	"fmt"
	"github.com/kjbreil/hass-ws/model"
	"sync"
)

type states struct {
	s sync.Map
}
type States map[string]*State

func (s *states) Store(ps *State) {
	s.s.Store(ps.DomainEntity, ps)
}

func (ss States) Entities() []string {
	var en []string
	for _, s := range ss {
		en = append(en, s.DomainEntity)
	}
	return en
}
func (s *states) Load(key string) (*State, bool) {
	st, ok := s.s.Load(key)
	if !ok {
		return nil, ok
	}
	return st.(*State), true
}

func (s *states) Find(keys []string) map[string]*State {
	ss := make(map[string]*State)
	s.s.Range(func(key, value any) bool {
		if len(keys) == 0 {
			return false
		}
		for i, k := range keys {
			if key == k {
				ss[k] = value.(*State)
				keys = append(keys[:i], keys[i+1:]...)
				break
			}
		}
		return true
	})
	return ss
}

func (gs *GoScript) GetState(domain, entityid string) *State {
	s, _ := gs.states.Load(fmt.Sprintf("%s%s", domain, entityid))
	return s
}

func (gs *GoScript) GetStates(domainentity []string) map[string]*State {
	rtn := gs.states.Find(domainentity)
	return rtn
}

func (gs *GoScript) handleMessage(message model.Message) {
	switch message.Type {
	case model.MessageTypeEvent:

		switch message.Event.EventType {
		case model.EventTypeStateChanged:

			s := &State{
				DomainEntity: message.DomainEntity(),
				Domain:       message.Domain(),
				Entity:       message.EntityID(),
				State:        message.State(),
				Attributes:   message.Attributes(),
			}

			gs.states.Store(s)

			gs.runTriggers(message)
		}
	}
}

func (gs *GoScript) handleGetStates(states []model.Result) {

	for _, sr := range states {
		s := &State{
			DomainEntity: sr.DomainEntity(),
			Domain:       sr.Domain(),
			Entity:       sr.EntityID(),
			State:        sr.State(),
			Attributes:   sr.Attributes,
		}

		gs.states.Store(s)

	}

}
