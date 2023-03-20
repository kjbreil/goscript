package goscript

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kjbreil/hass-ws/model"
	"sync"
)

type states struct {
	s sync.Map
}
type State struct {
	DomainEntity string
	Domain       string
	Entity       string
	State        string
	Attributes   map[string]interface{}
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
	newKeys := make([]string, len(keys))
	for i := range keys {
		newKeys[i] = keys[i]
	}

	ss := make(map[string]*State)
	s.s.Range(func(key, value any) bool {
		if len(newKeys) == 0 {
			return false
		}
		for i, k := range newKeys {
			if key == k {
				ss[key.(string)] = value.(*State)
				newKeys = append(newKeys[:i], newKeys[i+1:]...)
				break
			}
		}
		return true
	})
	return ss
}

func (s *states) FindDomain(keys []string) map[string]*State {

	ss := make(map[string]*State)
	for _, k := range keys {
		s.s.Range(func(key, value any) bool {
			switch value.(type) {
			case *State:
				if value.(*State).Domain == k {
					ss[key.(string)] = value.(*State)
				}
			}
			return true
		})
	}
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

func (gs *GoScript) GetDomainStates(domainentity []string) map[string]*State {
	rtn := gs.states.FindDomain(domainentity)
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
	statesFuncToRun := make(map[uuid.UUID]*Task)

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

	for _, sr := range states {
		domainEntity := sr.DomainEntity()
		entityState := sr.State()
		message := &model.Message{
			Type: model.MessageTypeEvent,
			Event: &model.Event{
				Data: &model.Data{
					EntityId: &domainEntity,
					NewState: &model.State{
						EntityId:    &domainEntity,
						LastChanged: sr.LastChanged,
						State:       &entityState,
						Attributes:  sr.Attributes,
						LastUpdated: sr.LastUpdated,
						Context:     sr.Context,
					},
					OldState: nil,
				},
				EventType: model.EventTypeStateChanged,
				Context:   sr.Context,
			},
		}

		gs.runTriggers(*message)
	}
	for k, t := range statesFuncToRun {
		gs.funcToRun[k] = t
	}

}

func MessageState(message *model.Message) *State {
	return &State{
		DomainEntity: message.DomainEntity(),
		Domain:       message.Domain(),
		Entity:       message.EntityID(),
		State:        message.State(),
		Attributes:   message.Attributes(),
	}
}
