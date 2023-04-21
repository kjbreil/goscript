package goscript

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kjbreil/hass-ws/model"
	"sync"
)

type States struct {
	s map[string]*State
	m *sync.Mutex
}
type State struct {
	DomainEntity string
	Domain       string
	Entity       string
	State        string
	Attributes   map[string]interface{}
}

func (s *States) Combine(cs *States) {
	s.m.Lock()
	defer s.m.Unlock()

	cs.m.Lock()
	defer cs.m.Unlock()
	for k, v := range cs.s {
		s.s[k] = v
	}
}

func (s *States) Store(ps *State) {
	s.m.Lock()
	defer s.m.Unlock()
	// Update the state pointer so it follows to tasks
	if st, ok := s.s[ps.DomainEntity]; ok {
		*st = *ps
	} else {
		s.s[ps.DomainEntity] = ps
	}
}

func (s *States) Entities() []string {
	s.m.Lock()
	defer s.m.Unlock()

	en := make([]string, 0, len(s.s))
	for _, st := range s.s {
		en = append(en, st.DomainEntity)
	}
	return en
}

func (s *States) Get(key string) (*State, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	st, ok := s.s[key]
	if !ok {
		return nil, ok
	}
	return st, true
}

func (s *States) Find(keys []string) map[string]*State {
	s.m.Lock()
	defer s.m.Unlock()

	ss := make(map[string]*State)

	for _, k := range keys {
		if st, ok := s.s[k]; ok {
			ss[st.DomainEntity] = st
		}
	}

	return ss
}

func (s *States) Slice() []*State {
	s.m.Lock()
	defer s.m.Unlock()

	var ss []*State
	for _, st := range s.s {
		ss = append(ss, st)
	}

	return ss
}

func (s *States) Map() map[string]*State {
	s.m.Lock()
	defer s.m.Unlock()

	sts := make(map[string]*State)

	for k, v := range s.s {
		sts[k] = v
	}

	return sts
}

func (s *States) SubSet(keys []string) States {
	s.m.Lock()
	defer s.m.Unlock()

	sts := States{
		s: make(map[string]*State),
		m: &sync.Mutex{},
	}

	for _, k := range keys {
		if st, ok := s.s[k]; ok {
			sts.Store(st)
		}
	}

	return sts
}

func (s *States) FindDomain(keys []string) map[string]*State {
	s.m.Lock()
	defer s.m.Unlock()

	ss := make(map[string]*State)

	for _, st := range s.s {
		for _, k := range keys {
			if st.Domain == k {
				ss[st.DomainEntity] = st
			}
		}
	}

	return ss
}

func (gs *GoScript) GetState(domain, entityid string) *State {
	s, _ := gs.states.Get(fmt.Sprintf("%s%s", domain, entityid))
	return s
}

func (gs *GoScript) GetStates(domainentity []string) *States {
	rtn := States{
		s: gs.states.Find(domainentity),
		m: &sync.Mutex{},
	}
	return &rtn
}

func (gs *GoScript) GetDomainStates(domainentity []string) *States {
	rtn := States{
		s: gs.states.FindDomain(domainentity),
		m: &sync.Mutex{},
	}

	return &rtn
}

func (gs *GoScript) handleMessage(message model.Message) {
	if message.Type == model.MessageTypeEvent {
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
	for _, t := range statesFuncToRun {
		gs.taskToRun.add(t)
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
