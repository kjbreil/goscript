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

// Insert only adds to the map if something does not exist already. Returns what is in the map whether added or not
func (s *States) Insert(ps *State) *State {
	s.m.Lock()
	defer s.m.Unlock()
	// Update the state pointer so it follows to tasks
	if st, ok := s.s[ps.DomainEntity]; ok {
		return st
	} else {
		s.s[ps.DomainEntity] = ps
		return ps
	}
}

// Upsert inserts a new record if one does not exist otherwise updates the data at the pointer so the update propagates
func (s *States) Upsert(ps *State) *State {
	s.m.Lock()
	defer s.m.Unlock()
	// Update the state pointer so it follows to tasks
	if st, ok := s.s[ps.DomainEntity]; ok {
		*st = *ps
	} else {
		s.s[ps.DomainEntity] = ps
	}
	return ps
}

// Combine takes two States objects and merges them, passed object will overwrite a state in current object
func (s *States) Combine(cs *States) {
	s.m.Lock()
	defer s.m.Unlock()

	cs.m.Lock()
	defer cs.m.Unlock()
	for k, v := range cs.s {
		s.s[k] = v
	}
}

// Entities returns a string of the entities contained in the States object
func (s *States) Entities() []string {
	s.m.Lock()
	defer s.m.Unlock()

	en := make([]string, 0, len(s.s))
	for _, st := range s.s {
		en = append(en, st.DomainEntity)
	}
	return en
}

// Get returns a single state record and a bool if found
func (s *States) Get(key string) (*State, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	st, ok := s.s[key]
	if !ok {
		return nil, ok
	}
	return st, true
}

// Find returns a new map of states of the passed entities
func (s *States) Find(entities []string) map[string]*State {
	s.m.Lock()
	defer s.m.Unlock()

	ss := make(map[string]*State)

	for _, k := range entities {
		if st, ok := s.s[k]; ok {
			ss[st.DomainEntity] = st
		}
	}

	return ss
}

// FindDomainMap returns a map of the states for the passed domain
func (s *States) FindDomainMap(keys []string) map[string]*State {
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

// Slice returns a slice of the states in no particular order
func (s *States) Slice() []*State {
	s.m.Lock()
	defer s.m.Unlock()

	var ss []*State
	for _, st := range s.s {
		ss = append(ss, st)
	}

	return ss
}

// Map returns a map of all the states
func (s *States) Map() map[string]*State {
	s.m.Lock()
	defer s.m.Unlock()

	sts := make(map[string]*State)

	for k, v := range s.s {
		sts[k] = v
	}

	return sts
}

// SubSet returns a new States which contains a subset of the current states based on entities passed
func (s *States) SubSet(entities []string) States {
	s.m.Lock()
	defer s.m.Unlock()

	return States{
		s: s.Find(entities),
		m: &sync.Mutex{},
	}
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
		s: gs.states.FindDomainMap(domainentity),
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

			gs.states.Upsert(s)

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

		gs.states.Upsert(s)
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
